/*
utils.go contains functions that aren't based on game data.
*/

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func initialize() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("loading .env vars: %w", err)
	}

	rs.SteamWebAPIKey = os.Getenv("STEAM_WEB_API_KEY")

	return nil
}

// unmarshalJsonApiData sends a request to url, obtains response data,
// then unmarshals that data into v.
// v must be a pointer.
func unmarshalJsonApiData(url string, v any) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("getting data: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		return fmt.Errorf("umarshaling response into v: %w", err)
	}

	resp.Body.Close()

	return nil
}

func excelToNumber(col string) int {
	result := 0
	for i := 0; i < len(col); i++ {
		result = result*26 + int(col[i]-'A'+1)
	}
	return result
}

func numberToExcel(num int) string {
	var result strings.Builder
	for num > 0 {
		num--
		result.WriteByte(byte(num%26 + 'A'))
		num /= 26
	}

	runes := []rune(result.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func getNextColumn(column string) string {
	columnNumber := excelToNumber(column)
	nextColumnNumber := columnNumber + 1
	return numberToExcel(nextColumnNumber)
}
