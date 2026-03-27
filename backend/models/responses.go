package models

type ParticipantSummary struct {
	PUUID        string `json:"puuid"`
	RiotId       string `json:"riotId"`
	ChampionName string `json:"championName"`
	KDA          string `json:"kda"`
	Win          bool   `json:"win"`
	IsDuoPartner bool   `json:"isDuoPartner"`
}

type MatchSummary struct {
	MatchId    string               `json:"matchId"`
	GameDate   string               `json:"gameDate"`
	Duration   string               `json:"duration"`
	Win        bool                 `json:"win"`
	Champion   string               `json:"champion"`
	KDA        string               `json:"kda"`
	Teammates  []ParticipantSummary `json:"teammates"`  // same team, excluding target
	DuoPartner *ParticipantSummary  `json:"duoPartner"` // most frequent, if in this game
}

type FrequentPartner struct {
	PUUID       string `json:"puuid"`
	RiotId      string `json:"riotId"`
	GamesPlayed int    `json:"gamesPlayed"`
}

type DuoResponse struct {
	TargetRiotId     string            `json:"targetRiotId"`
	FrequentPartners []FrequentPartner `json:"frequentPartners"`
	Matches          []MatchSummary    `json:"matches"`
}
