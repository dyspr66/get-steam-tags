package main

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/xuri/excelize/v2"
)

type Game struct {
	ID    int
	Title string

	ScrapedOn        time.Time
	ReleaseDate      string
	TotalReviewCount string
	ReviewPositivity string
	Tags             []string
}

// Games is safe to use concurrently.
type Games struct {
	mu sync.Mutex
	// Games map[int]Game
	ProcessedGameCount int

	TagToCol map[string]string

	CurrentMaxRow int    // tracks last edded row
	CurrentMaxCol string // tracks last added column
	Spreadsheet   excelize.File
}

func (gs *Games) Update(key int, game Game) error {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	gs.ProcessedGameCount += 1

	// Update metadata
	f.SetCellValue(mainSheet, fmt.Sprintf("A%d", gs.CurrentMaxRow), game.Title)
	f.SetCellValue(mainSheet, fmt.Sprintf("B%d", gs.CurrentMaxRow), game.ScrapedOn)
	f.SetCellValue(mainSheet, fmt.Sprintf("C%d", gs.CurrentMaxRow), game.ReleaseDate)
	f.SetCellValue(mainSheet, fmt.Sprintf("D%d", gs.CurrentMaxRow), game.TotalReviewCount)
	f.SetCellValue(mainSheet, fmt.Sprintf("E%d", gs.CurrentMaxRow), game.ReviewPositivity)

	// Update tags
	for _, tag := range game.Tags {
		_, tagExists := gs.TagToCol[tag]
		if tagExists {
			tagColumn := gs.TagToCol[tag] // Get column for existing tag

			gameTagPairCell := fmt.Sprintf("%s%d", tagColumn, gs.CurrentMaxRow) // Get cell for game:tag pair
			f.SetCellValue(mainSheet, gameTagPairCell, "1")                     // Mark game:tag pair as 1
		} else {
			// Add new column to map
			gs.TagToCol[tag] = gs.CurrentMaxCol

			columnCell := fmt.Sprintf("%s1", gs.CurrentMaxCol) // Get cell for tag's new column
			f.SetCellValue(mainSheet, columnCell, tag)         // Add tag as new column

			gameTagPairCell := fmt.Sprintf("%s%d", gs.CurrentMaxCol, gs.CurrentMaxRow) // Get cell for game:tag pair
			f.SetCellValue(mainSheet, gameTagPairCell, "1")                            // Mark game:tag pair as 1

			// Shift to next tag
			gs.CurrentMaxCol = getNextColumn(gs.CurrentMaxCol)
		}
	}

	// Go to next row
	gs.CurrentMaxRow += 1

	// Save every n games
	if gs.ProcessedGameCount%10 == 0 { // TODO
		slog.Info("Saving progress", "games processed", gs.ProcessedGameCount, "last processed game id", game.ID, "last processed game title", game.Title)
		if err := f.SaveAs("game_tags_adjacency_matrix.xlsx"); err != nil {
			return fmt.Errorf("saving excel file: %w", err)
		}
	}

	return nil
}

func (gs *Games) GetAndSetData(key int, title string, wg *sync.WaitGroup) {
	defer wg.Done()
	slog.Info("Getting data", "key", key)

	// Get metadata
	url := fmt.Sprintf("https://store.steampowered.com/app/%d/", key)

	scrapedOn := time.Now() // Get time scraped

	releaseDate, totalReviewCount, reviewPositivity, err := scrapeMetadata(url)
	if err != nil {
		slog.Error("Getting metadata for game", "id", key, "title", title, "err", err)
		// return // NOTE - uncomment if err should make break program
	}

	// Get tags
	tags, err := scrapeUserTags(url)
	if err != nil {
		slog.Error("Getting tags for game", "id", key, "title", title, "err", err)
		// return // NOTE - uncomment if err should make break program
	}

	// Update data
	game := Game{
		ID:    key,
		Title: title,

		ScrapedOn:        scrapedOn,
		ReleaseDate:      releaseDate,
		TotalReviewCount: totalReviewCount,
		ReviewPositivity: reviewPositivity,
		Tags:             tags,
	}

	err = gs.Update(key, game)
	if err != nil {
		slog.Warn("Updating data", "game id", game.ID, "game title", game.Title, "err", err)
	} else {
		slog.Info("Obtained data", "key", key)
	}
}
