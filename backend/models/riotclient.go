package models

import (
	"net/http"
)

type RiotClient struct {
	apiKey     string
	httpClient *http.Client
}
