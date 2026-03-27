package models

type AccountDTO struct {
	PUUID    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
}

type Participant struct {
	PUUID          string `json:"puuid"`
	RiotIdGameName string `json:"riotIdGameName"`
	RiotIdTagLine  string `json:"riotIdTagLine"`
	TeamId         int    `json:"teamId"`
	ChampionName   string `json:"championName"`
	Kills          int    `json:"kills"`
	Deaths         int    `json:"deaths"`
	Assists        int    `json:"assists"`
	Win            bool   `json:"win"`
	TeamPosition   string `json:"teamPosition"`
}

type MatchInfo struct {
	GameStartTimestamp int64         `json:"gameStartTimestamp"`
	GameDuration       int           `json:"gameDuration"`
	Participants       []Participant `json:"participants"`
}

type MatchDTO struct {
	Metadata struct {
		MatchId string `json:"matchId"`
	} `json:"metadata"`
	Info MatchInfo `json:"info"`
}
