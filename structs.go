package main

type readSteamDB struct {
	SteamWebAPIKey string
	AllGames       GetAppList
}

// Models the JSON data obtained from https://api.steampowered.com/IStoreService/GetAppList/v1/
type GetAppList struct {
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
