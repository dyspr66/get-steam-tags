package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
)

type Record struct {
	ID         int    `json:"sid"`
	Name       string `json:"name"`
	Categories string `json:"tags"`
}

func main() {
	// Read file to get game records
	data, err := os.ReadFile("steamdb.json")
	if err != nil {
		fmt.Println(err)
	}

	var records []Record
	err = json.Unmarshal(data, &records)
	if err != nil {
		fmt.Println(err)
	}

	// Make excel sheet
	main_sheet := "Sheet1"
	column_game_name := "A"
	column_game_id := "B"
	column_game_tags := "C"

	f := excelize.NewFile()
	f.NewSheet(main_sheet)

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Place records into file
	row := 0
	for i, record := range records {
		fmt.Printf("Reading %dth record: %s\n", i, record.Name)

		tags := strings.Split(record.Categories, ",")
		for _, tag := range tags {
			game_name_cell := fmt.Sprintf("%s%d", column_game_name, row)
			game_id_cell := fmt.Sprintf("%s%d", column_game_id, row)
			game_tag_cell := fmt.Sprintf("%s%d", column_game_tags, row)

			f.SetCellValue(main_sheet, game_name_cell, record.Name)
			f.SetCellValue(main_sheet, game_id_cell, record.ID)
			f.SetCellValue(main_sheet, game_tag_cell, tag)

			row += 1 // go to next row for next tag
		}
	}

	// Save file as "all_steam_games_data.xlsx"
	if err := f.SaveAs("all_steam_games_data.xlsx"); err != nil {
		fmt.Println(err)
	}
}
