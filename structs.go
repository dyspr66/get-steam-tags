package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type readSteamDB struct {
	SteamWebAPIKey string
	AllGames       GetAppListResponse
}

// Models the JSON data obtained from https://api.steampowered.com/IStoreService/GetAppListResponse/v1/
type GetAppListResponse struct {
	Response response `json:"response"`
}

type response struct {
	Apps            []App `json:"apps"`
	HaveMoreResults bool  `json:"have_more_results"`
	LastAppID       int   `json:"last_appid"`
}

type AppList struct {
	Apps []App `json:"apps"`
}

type App struct {
	AppID int    `json:"appid"`
	Title string `json:"name"`
}

type SteamSpyResponse struct {
	Tags map[string]int `json:"tags"`
}

func (s *SteamSpyResponse) UnmarshalJSON(data []byte) error {
	var temp struct {
		Tags map[string]int `json:"tags"`
	}

	err := json.Unmarshal(data, &temp)

	s.Tags = temp.Tags

	if err != nil && strings.Contains(err.Error(), "cannot unmarshal array") {
		return nil
	}

	if err == nil {
		return nil
	}

	return fmt.Errorf("unmarshaling SteamSpyResponse: %w", err)
}
