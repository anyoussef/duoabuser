package api

import (
	"duo-abuser/models"
	"encoding/json"
	"io"
	"net/http"
)

var account_endpoint string = "https://americas.api.riotgames.com/riot/account/v1/accounts/by-riot-id/"
var match_endpoint string = "https://americas.api.riotgames/lol/match/v5/matches/by-puuid"

/* Provided with user input: name, tag map returned api request to structure summoner */
func GetSummoner(name string, tag string, apiKey string) (models.Summoner, error) {
	var summoner models.Summoner

	resp, err := http.Get(account_endpoint + name + "/" + tag + "?api_key=" + apiKey)
	if err != nil {
		return summoner, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return summoner, err
	}

	err = json.Unmarshal(body, &summoner)
	return summoner, err
}

/* Based off a puuid get the last_20_games for that player */
func GetMatchesPUUID(puuid string) ([]models.Game, error) {
	resp, err := http.Get(match_endpoint + puuid + "/ids")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Step 1: unmarshal into []string
	var matchIDs []string
	err = json.Unmarshal(body, &matchIDs)
	if err != nil {
		return nil, err
	}

	// Step 2: convert to []Game
	var games []models.Game
	for _, id := range matchIDs {
		games = append(games, models.Game{GameID: id})
	}

	return games, nil
}
