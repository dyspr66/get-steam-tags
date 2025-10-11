package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var rs readSteamDB

func main() {
	initialize()

	getIDsForAllGames()

	// var tagToCol map[string]string // maps a tag to its column id
	// var currentMaxCol string       // track last added column
	// var currentMaxRow string       // track last edded row

	for _, game := range rs.AllGames.Response.Apps {
		fmt.Println("Finding tags for", game.AppID, ":", game.Name)
		// TODO
	}

}

func initialize() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("loading env: ", err)
	}

	rs.SteamWebAPIKey = os.Getenv("STEAM_WEB_API_KEY")
}

func getIDsForAllGames() {
	url := fmt.Sprintf("https://api.steampowered.com/IStoreService/GetAppList/v1/?include_games=true&include_dlc=true&max_results=50000&key=%s", rs.SteamWebAPIKey)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("getting data: ", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("reading response body: ", err)
	}

	var base GetAppList
	json.Unmarshal(body, &base)
	if err != nil {
		log.Fatal("unmarshaling body", err)
	}
	resp.Body.Close()

	rs.AllGames.Response.Apps = append(rs.AllGames.Response.Apps, base.Response.Apps...)
	fmt.Println("Sucessfully obtained all games up to ID", base.Response.LastAppID)

	for base.Response.HaveMoreResults {
		url := fmt.Sprintf("https://api.steampowered.com/IStoreService/GetAppList/v1/?include_games=true&include_dlc=true&last_appid=%d&max_results=50000&key=%s", base.Response.LastAppID, rs.SteamWebAPIKey)
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal("getting data: ", err)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("reading response body: ", err)
		}

		var r GetAppList
		json.Unmarshal(body, &r)
		if err != nil {
			log.Fatal("unmarshaling body", err)
		}
		resp.Body.Close()

		base = r
		rs.AllGames.Response.Apps = append(rs.AllGames.Response.Apps, base.Response.Apps...)
		fmt.Println("Sucessfully obtained all games up to ID", base.Response.LastAppID)
	}
}
