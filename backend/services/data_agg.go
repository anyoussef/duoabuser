package services

import (
	"duo-abuser/api"
	"duo-abuser/models"
)

func GetSummonerWithGames(name, tag, apiKey string) (models.SummonerGames, error) {
	summoner, err := api.GetSummoner(name, tag, apiKey)
	if err != nil {
		return models.SummonerGames{}, err
	}

	games, err := api.GetMatchesPUUID(summoner.Puuid)
	if err != nil {
		return models.SummonerGames{}, err
	}

	return models.SummonerGames{
		Summoner: summoner,
		Games:    games,
	}, nil
}
