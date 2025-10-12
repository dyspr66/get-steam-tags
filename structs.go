package main

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
	Name  string `json:"name"`
}

type SteamSpyResponse struct {
	Tags map[string]int `json:"tags"`
}
