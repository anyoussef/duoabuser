package api

import (
	"duo-abuser/models"
	"encoding/json"
	"io"
	"net/http"
)

var endpoint string = "https://americas.api.riotganes.com/riot/account/v1/accounts/by-riot-id/"

func GetSummoner(name, tag, apiKey string) (models.Summoner, error) {
	var summoner models.Summoner

	resp, err := http.Get(endpoint + name + tag + "?api_key" + apiKey)
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
