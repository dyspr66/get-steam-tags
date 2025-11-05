/*
main.go contains functions related to generating the adjacency matrix.
*/

package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/xuri/excelize/v2"
)

var rs readSteamDB

func main() {
	// Open or create a log file
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("Opening log file: " + err.Error())
	}
	defer file.Close()

	fileHandler := slog.NewTextHandler(file, nil)
	logger := slog.New(fileHandler)
	slog.SetDefault(logger)

	// Actually start the program
	slog.Info("Starting program...")
	fmt.Println(time.Now(), "INFO: Starting program...") // TODO - Have 2 slog outputs instead of this

	err = initialize()
	if err != nil {
		slog.Error("Initializing", "err", err)
		fmt.Println(time.Now(), "ERROR: Initializing:", "err =", err)
		return
	}

	err = getIDsForAllGames()
	if err != nil {
		slog.Error("Getting all game IDs", "err", err)
		fmt.Println(time.Now(), "ERROR: Getting all game IDs:", "err =", err)
		return
	}

	// Make excel file
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			slog.Error("Closing excel file", "err", err)
			fmt.Println(time.Now(), "ERROR: Closing excel file:", "err =", err)
			return
		}
	}()

	mainSheet := "Sheet1" // name of sheet data is stored in
	currentMaxRow := 2    // track last edded row
	currentMaxCol := "F"  // track last added column

	// Set default column titles
	f.SetCellValue(mainSheet, "A1", "Title")
	f.SetCellValue(mainSheet, "B1", "Scraped On")
	f.SetCellValue(mainSheet, "C1", "Release Date")
	f.SetCellValue(mainSheet, "D1", "Total Review Count")
	f.SetCellValue(mainSheet, "E1", "Review Positivity")

	tagToCol := make(map[string]string) // maps a tag to its column id

	// Go through each game and record their tags
	for i, game := range rs.AllGames.Response.Apps {
		url := fmt.Sprintf("https://store.steampowered.com/app/%d/", game.AppID)

		// Get time scraped
		scrapedOn := time.Now()

		// Get metadata
		releaseDate, totalReviewCount, reviewPositivity, err := scrapeMetadata(url)
		if err != nil {
			slog.Error("Getting metadata for game", "id", game.AppID, "title", game.Title, "err", err)
			fmt.Println(time.Now(), "ERROR: Getting metadata for game:", "id =", game.AppID, "title =", game.Title, "err =", err)

			// NOTE - Comment out the "return" if the whole program
			// should stop when tag scrape error is detected:
			// return
		}

		// Set title, time scraped, and metadata
		f.SetCellValue(mainSheet, fmt.Sprintf("A%d", currentMaxRow), game.Title)
		f.SetCellValue(mainSheet, fmt.Sprintf("B%d", currentMaxRow), scrapedOn)
		f.SetCellValue(mainSheet, fmt.Sprintf("C%d", currentMaxRow), releaseDate)
		f.SetCellValue(mainSheet, fmt.Sprintf("D%d", currentMaxRow), totalReviewCount)
		f.SetCellValue(mainSheet, fmt.Sprintf("E%d", currentMaxRow), reviewPositivity)

		// Get and set tags
		tags, err := scrapeUserTags(url)
		if err != nil {
			slog.Error("Getting tags for game", "id", game.AppID, "title", game.Title, "err", err)
			fmt.Println(time.Now(), "ERROR: Getting tags for game:", "id =", game.AppID, "title =", game.Title, "err =", err)

			// NOTE - Comment out the "return" if the whole program
			// should stop when tag scrape error is detected:
			// return
		}

		for _, tag := range tags {
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

		if i%100 == 0 { // Save every 100 records
			slog.Info("Saving data:", "processed count =", i)
			fmt.Println(time.Now(), "INFO: Saving data:", "processed count =", i)

			if err := f.SaveAs("game_tags_adjacency_matrix.xlsx"); err != nil {
				slog.Error("Saving excel file", "err", err)
				fmt.Println(time.Now(), "ERROR: Saving excel file:", "err =", err)
			}
		}

		currentMaxRow += 1 // Shift to next game
		slog.Info("Finished processing game", "count", i, "id", game.AppID, "title", game.Title)
		fmt.Println(time.Now(), "INFO: Finished processing game:", "count =", i, "id =", game.AppID, "title =", game.Title)
	}

	if err := f.SaveAs("game_tags_adjacency_matrix.xlsx"); err != nil {
		slog.Error("Saving excel file", "err", err)
		fmt.Println(time.Now(), "ERROR: Saving excel file:", "err =", err)
		return
	}

	slog.Info("Program finished.")
	fmt.Println(time.Now(), "INFO: Program finished.")
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
	fmt.Println(time.Now(), "INFO: Obtained all games up to a certain ID:", "last obtained ID =", base.Response.LastAppID)

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
		fmt.Println(time.Now(), "INFO: Obtained all games up to a certain ID:", "last obtained ID =", base.Response.LastAppID)
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
