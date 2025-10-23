/*
main.go contains functions related to generating the adjacency matrix.
*/

package main

import (
	"fmt"
	"log/slog"

	"github.com/xuri/excelize/v2"
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

	// Make excel file
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			slog.Error("Closing excel file", "err", err)
			return
		}
	}()

	mainSheet := "Sheet1" // name of sheet data is stored in
	currentMaxRow := 2    // track last edded row
	currentMaxCol := "B"  // track last added column

	tagToCol := make(map[string]string) // maps a tag to its column id

	// Go through each game and record their tags
	for i, game := range rs.AllGames.Response.Apps {
		// Add game as new row
		cell := fmt.Sprintf("A%d", currentMaxRow)
		f.SetCellValue(mainSheet, cell, game.Name)

		// Process game's tags
		var steamSpyResponse SteamSpyResponse
		err := getTagsForGame(game, &steamSpyResponse)
		if err != nil {
			slog.Error("Getting tags for game", "id", game.AppID, "name", game.Name, "err", err)

			// Comment out the return if the whole program shouldn't
			// stop if an error with tags is detected:
			// return
		}

		for tag := range steamSpyResponse.Tags {
			_, tagExists := tagToCol[tag]
			if tagExists {
				tagColumn := tagToCol[tag] // Get column for existing tag

				gameTagPairCell := fmt.Sprintf("%s%d", tagColumn, currentMaxRow) // Get cell for game:tag pair
				f.SetCellValue(mainSheet, gameTagPairCell, "1")                  // Mark game:tag pair as 1
			} else {
				// Add new column to map
				tagToCol[tag] = currentMaxCol

				columnCell := fmt.Sprintf("%s1", currentMaxCol) // Get cell for tag's new column
				f.SetCellValue(mainSheet, columnCell, tag)      // Add tag as new column

				gameTagPairCell := fmt.Sprintf("%s%d", currentMaxCol, currentMaxRow) // Get cell for game:tag pair
				f.SetCellValue(mainSheet, gameTagPairCell, "1")                      // Mark game:tag pair as 1

				// Shift to next tag
				currentMaxCol = getNextColumn(currentMaxCol)
			}
		}

		currentMaxRow += 1 // Shift to next game
		slog.Info("Obtained tags for game", "count", i, "id", game.AppID, "name", game.Name)
	}

	if err := f.SaveAs("game_tags_adjacency_matrix.xlsx"); err != nil {
		slog.Error("Saving excel file", "err", err)
		return
	}

	slog.Info("Program finished.")
}

func getIDsForAllGames() error {
	url := fmt.Sprintf("https://api.steampowered.com/IStoreService/GetAppList/v1/?include_games=true&include_dlc=false&include_software=false&include_videos=false&include_hardware=false&max_results=50000&key=%s", rs.SteamWebAPIKey)
	var base GetAppListResponse

	err := unmarshalJsonApiData(url, &base)
	if err != nil {
		return fmt.Errorf("getting id for all games: %w", err)
	}

	rs.AllGames.Response.Apps = append(rs.AllGames.Response.Apps, base.Response.Apps...)
	slog.Info("Obtained all games up to a certain ID", "last obtained ID", base.Response.LastAppID)

	for base.Response.HaveMoreResults {
		url := fmt.Sprintf("https://api.steampowered.com/IStoreService/GetAppList/v1/?include_games=true&include_dlc=false&include_software=false&include_videos=false&include_hardware=false&max_results=50000&last_appid=%d&key=%s", base.Response.LastAppID, rs.SteamWebAPIKey)

		var r GetAppListResponse
		err := unmarshalJsonApiData(url, &r)
		if err != nil {
			return fmt.Errorf("getting id for all games: %w", err)
		}

		base = r

		rs.AllGames.Response.Apps = append(rs.AllGames.Response.Apps, base.Response.Apps...)
		slog.Info("Obtained all games up to a certain ID", "last obtained ID", base.Response.LastAppID)
	}

	return nil
}

func getTagsForGame(game App, tags *SteamSpyResponse) error {
	url := fmt.Sprintf("https://steamspy.com/api.php?request=appdetails&appid=%d", game.AppID)

	err := unmarshalJsonApiData(url, tags)
	if err != nil {
		return fmt.Errorf("getting tags for game: %w", err)
	}

	return nil
}
