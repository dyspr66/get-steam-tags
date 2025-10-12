/*
main.go contains functions related to generating the adjacency matrix.
*/

package main

import (
	"fmt"
	"log/slog"
)

var rs readSteamDB

func main() {
	slog.Info("Starting program...")

	err := initialize()
	if err != nil {
		slog.Error("Initializing", "err", err)
		return
	}

	err = getIDsForAllGames()
	if err != nil {
		slog.Error("Getting all game IDs", "err", err)
		return
	}

	// var tagToCol map[string]string // maps a tag to its column id
	// var currentMaxCol string       // track last added column
	// var currentMaxRow string       // track last edded row

	for _, game := range rs.AllGames.Response.Apps {
		var steamSpyResponse SteamSpy
		err := getTagsForGame(game, &steamSpyResponse)
		if err != nil {
			slog.Error("Getting tags for game", "id", game.AppID, "name", game.Name, "err", err)
			return
		}

		// for tag, _ := range steamSpyResponse.Tags {
		// }

		slog.Info("Obtained tags for game", "id", game.AppID, "name", game.Name)
	}
}

func getIDsForAllGames() error {
	url := fmt.Sprintf("https://api.steampowered.com/IStoreService/GetAppList/v1/?include_games=true&include_dlc=true&max_results=50000&key=%s", rs.SteamWebAPIKey)
	var base GetAppListResponse

	err := unmarshalJsonApiData(url, &base)
	if err != nil {
		return fmt.Errorf("getting id for all games: %w", err)
	}

	rs.AllGames.Response.Apps = append(rs.AllGames.Response.Apps, base.Response.Apps...)
	slog.Info("Sucessfully obtained all games up to a certain ID", "last obtained ID", base.Response.LastAppID)

	// for base.Response.HaveMoreResults {
	// 	url := fmt.Sprintf("https://api.steampowered.com/IStoreService/GetAppList/v1/?include_games=true&include_dlc=true&last_appid=%d&max_results=50000&key=%s", base.Response.LastAppID, rs.SteamWebAPIKey)

	// 	var r GetAppList
	// 	err := unmarshalJsonApiData(url, &r)
	// 	if err != nil {
	// 		return fmt.Errorf("getting id for all games: %w", err)
	// 	}

	// 	base = r

	// 	rs.AllGames.Response.Apps = append(rs.AllGames.Response.Apps, base.Response.Apps...)
	// 	slog.Info("Sucessfully obtained all games up to a certain ID", "last obtained ID", base.Response.LastAppID)
	// }

	return nil
}

func getTagsForGame(game App, tags *SteamSpy) error {
	url := fmt.Sprintf("https://steamspy.com/api.php?request=appdetails&appid=%d", game.AppID)

	err := unmarshalJsonApiData(url, tags)
	if err != nil {
		return fmt.Errorf("getting tags for game: %w", err)
	}

	return nil
}
