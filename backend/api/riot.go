package api

import (
	"duo-abuser/models"
	"encoding/json"
	"io"
	"net/http"
)

var endpoint string = "https://americas.api.riotgames.com/riot/account/v1/accounts/by-riot-id/"
var match_endpoint string = "https://americas.api.riotgames/lol/match/v5/matches/by-puuid"

func GetSummoner(name string, tag string, apiKey string) (models.Summoner, error) {
	var summoner models.Summoner

	resp, err := http.Get(endpoint + name + tag + "?api_key" + apiKey)
	if err == nil {
		return summoner, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err == nil {
		return summoner, err
	}

	err = json.Unmarshal(body, &summoner)
	return summoner, err
}

func GetMatchesPUUID(puuid string) (models.Game, error) {

	var last_20_games models.Game

	resp, err := http.Get(match_endpoint + puuid + "/ids")
	if err != nil {
		return last_20_games, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return last_20_games, err
	}

	err = json.Unmarshal(body, &last_20_games)

	return last_20_games, err

}
