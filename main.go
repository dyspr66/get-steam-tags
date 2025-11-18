/*
main.go contains functions related to generating the adjacency matrix.
*/

package main

import (
	"log/slog"
	"os"
	"sync"

	"github.com/xuri/excelize/v2"
)

var rs readSteamDB
var f *excelize.File
var mainSheet = "Sheet1"

func main() {
	var err error

	// Open or create a log file
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("Opening log file: " + err.Error())
	}
	defer file.Close()

	slog.SetDefault(slog.New(NewCopyHandler(slog.NewTextHandler(os.Stdout, nil), slog.NewTextHandler(file, nil))))

	// Start program
	slog.Info("Starting program...")

	err = initialize()
	if err != nil {
		slog.Error("Initializing", "err", err)
		return
	}

	// Make excel file
	f = excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			slog.Error("Closing excel file", "err", err)
			return
		}
	}()

	// Set default excel column titles
	f.SetCellValue(mainSheet, "A1", "Title")
	f.SetCellValue(mainSheet, "B1", "ID")
	f.SetCellValue(mainSheet, "C1", "Scraped On")
	f.SetCellValue(mainSheet, "D1", "Release Date")
	f.SetCellValue(mainSheet, "E1", "Total Review Count")
	f.SetCellValue(mainSheet, "F1", "Review Positivity")

	// Get IDs of all steam games
	err = getIDsForAllGames()
	if err != nil {
		slog.Error("Getting all game IDs", "err", err)
		return
	}

	games := Games{ProcessedGameCount: 0, TagToCol: make(map[string]string)}
	games.CurrentMaxCol = "G"
	games.CurrentMaxRow = 2

	// Get and set data for each game
	var wg sync.WaitGroup
	wg.Add(len(rs.AllGames.Response.Apps))
	maxConcurrentRequests := 4
	limiter := make(chan struct{}, maxConcurrentRequests)

	for i, game := range rs.AllGames.Response.Apps {
		limiter <- struct{}{}

		go func() {
			defer func() { <-limiter }()
			slog.Info("Processing nth game", "n", i, "id", game.AppID, "title", game.Title)
			games.GetAndSetData(game.AppID, game.Title, &wg)
		}()
	}

	wg.Wait()

	// Finalize program
	if err := f.SaveAs("game_tags_adjacency_matrix.xlsx"); err != nil {
		slog.Error("Saving excel file", "err", err)
		return
	}

	slog.Info("Program finished.")
}
