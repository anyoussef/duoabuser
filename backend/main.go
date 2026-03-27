package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"duo-abuser/models"
	"github.com/joho/godotenv"
)

type RiotClient struct {
	apiKey     string
	httpClient *http.Client
}

func NewRiotClient(apiKey string) *RiotClient {
	return &RiotClient{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *RiotClient) get(url string, out interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Riot-Token", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("riot API returned %d for %s", resp.StatusCode, url)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func (c *RiotClient) GetAccount(gameName, tagLine string) (*models.AccountDTO, error) {
	url := fmt.Sprintf("https://americas.api.riotgames.com/riot/account/v1/accounts/by-riot-id/%s/%s",
		gameName, tagLine)
	var acc models.AccountDTO
	return &acc, c.get(url, &acc)
}

func (c *RiotClient) GetMatchIDs(puuid string, count int) ([]string, error) {
	url := fmt.Sprintf(
		"https://americas.api.riotgames.com/lol/match/v5/matches/by-puuid/%s/ids?queue=420&count=%d",
		puuid, count)
	var ids []string
	return ids, c.get(url, &ids)
}

func (c *RiotClient) GetMatch(matchID string) (*models.MatchDTO, error) {
	url := fmt.Sprintf("https://americas.api.riotgames.com/lol/match/v5/matches/%s", matchID)
	var m models.MatchDTO
	return &m, c.get(url, &m)
}

// ── Handler ───────────────────────────────────────────────────────────────────

type Server struct {
	riot *RiotClient
}

func (s *Server) handleDuo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	riotId := r.URL.Query().Get("riotId") // e.g. "Faker#NA1"
	if riotId == "" {
		http.Error(w, `{"error":"riotId query param required (e.g. Name%23TAG)"}`, 400)
		return
	}

	parts := strings.SplitN(riotId, "#", 2)
	if len(parts) != 2 {
		http.Error(w, `{"error":"riotId must be in format Name#TAG"}`, 400)
		return
	}
	gameName, tagLine := parts[0], parts[1]

	// 1. Resolve PUUID
	acc, err := s.riot.GetAccount(gameName, tagLine)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"could not find account: %v"}`, err), 404)
		return
	}

	// 2. Fetch recent ranked match IDs (last 20)
	matchIDs, err := s.riot.GetMatchIDs(acc.PUUID, 20)
	if err != nil || len(matchIDs) == 0 {
		http.Error(w, `{"error":"no ranked matches found"}`, 404)
		return
	}

	// 3. Fetch all matches concurrently
	type result struct {
		match *models.MatchDTO
		err   error
	}
	results := make([]result, len(matchIDs))
	var wg sync.WaitGroup
	// Rate-limit: Riot allows 20 req/sec on dev keys, use a simple semaphore
	sem := make(chan struct{}, 5)
	for i, id := range matchIDs {
		wg.Add(1)
		go func(idx int, matchId string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			m, err := s.riot.GetMatch(matchId)
			results[idx] = result{m, err}
		}(i, id)
	}
	wg.Wait()

	// 4. Build teammate frequency map and match summaries
	teammateCounts := map[string]int{}    // puuid -> games together
	teammateRiotId := map[string]string{} // puuid -> riotId string

	var matches []models.MatchSummary

	for _, res := range results {
		if res.err != nil || res.match == nil {
			continue
		}
		info := res.match.Info

		// Find the target player in this match
		var targetParticipant *models.Participant
		for i := range info.Participants {
			if info.Participants[i].PUUID == acc.PUUID {
				targetParticipant = &info.Participants[i]
				break
			}
		}
		if targetParticipant == nil {
			continue
		}

		// Build teammate list (same team, not target)
		var teammates []models.ParticipantSummary
		for _, p := range info.Participants {
			if p.PUUID == acc.PUUID || p.TeamId != targetParticipant.TeamId {
				continue
			}
			rid := p.RiotIdGameName
			if p.RiotIdTagLine != "" {
				rid += "#" + p.RiotIdTagLine
			}
			teammateCounts[p.PUUID]++
			teammateRiotId[p.PUUID] = rid

			teammates = append(teammates, models.ParticipantSummary{
				PUUID:        p.PUUID,
				RiotId:       rid,
				ChampionName: p.ChampionName,
				KDA:          fmt.Sprintf("%d/%d/%d", p.Kills, p.Deaths, p.Assists),
				Win:          p.Win,
			})
		}

		ms := models.MatchSummary{
			MatchId:   res.match.Metadata.MatchId,
			GameDate:  time.Unix(info.GameStartTimestamp/1000, 0).Format("Jan 2, 2006"),
			Duration:  fmt.Sprintf("%dm %ds", info.GameDuration/60, info.GameDuration%60),
			Win:       targetParticipant.Win,
			Champion:  targetParticipant.ChampionName,
			KDA:       fmt.Sprintf("%d/%d/%d", targetParticipant.Kills, targetParticipant.Deaths, targetParticipant.Assists),
			Teammates: teammates,
		}
		matches = append(matches, ms)
	}

	// 5. Sort frequent partners (min 2 games to count as a duo)
	type kv struct {
		puuid string
		count int
	}
	var sorted []kv
	for puuid, count := range teammateCounts {
		if count >= 2 {
			sorted = append(sorted, kv{puuid, count})
		}
	}
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].count > sorted[j].count })

	var frequentPartners []models.FrequentPartner
	duoPUUIDs := map[string]bool{}
	for _, kv := range sorted {
		frequentPartners = append(frequentPartners, models.FrequentPartner{
			PUUID:       kv.puuid,
			RiotId:      teammateRiotId[kv.puuid],
			GamesPlayed: kv.count,
		})
		duoPUUIDs[kv.puuid] = true
	}

	// 6. Annotate each match: mark duo partners and attach the top duo in this game
	for i := range matches {
		for j := range matches[i].Teammates {
			if duoPUUIDs[matches[i].Teammates[j].PUUID] {
				matches[i].Teammates[j].IsDuoPartner = true
				if matches[i].DuoPartner == nil {
					// pick the most frequent duo in this game
					best := -1
					var bestSummary *models.ParticipantSummary
					for k := range matches[i].Teammates {
						if matches[i].Teammates[k].IsDuoPartner {
							count := teammateCounts[matches[i].Teammates[k].PUUID]
							if count > best {
								best = count
								copy := matches[i].Teammates[k]
								bestSummary = &copy
							}
						}
					}
					matches[i].DuoPartner = bestSummary
				}
			}
		}
	}

	json.NewEncoder(w).Encode(models.DuoResponse{
		TargetRiotId:     gameName + "#" + tagLine,
		FrequentPartners: frequentPartners,
		Matches:          matches,
	})
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, falling back to environment variables")
	}

	apiKey := os.Getenv("RIOT_API_KEY")
	if apiKey == "" {
		log.Fatal("RIOT_API_KEY environment variable not set")
	}

	srv := &Server{riot: NewRiotClient(apiKey)}
	http.HandleFunc("/api/duo", srv.handleDuo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
