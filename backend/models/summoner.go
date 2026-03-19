package models

import ()

type Summoner struct {
	Puuid    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
}

type SummonerGames struct {
	Summoner Summoner `json:"summoner"`
	Games    []Game   `json:"games"`
}
