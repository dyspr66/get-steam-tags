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
