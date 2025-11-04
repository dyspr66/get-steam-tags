package main

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func scrapeUserTags(gameUrl string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	var tags string
	err := chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			// TODO - Set cookies for age verification
			expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))

			cookies := make(map[string]string)
			cookies["birthtime"] = "946706401"
			cookies["lastagecheckage"] = "1-January-2000"

			for k, v := range cookies {
				err := network.SetCookie(k, v).
					WithExpires(&expr).
					WithDomain("store.steampowered.com").
					WithHTTPOnly(true).
					Do(ctx)

				if err != nil {
					return err
				}
			}

			return nil
		}),

		// Go to URL for game and get tags
		chromedp.Navigate(gameUrl),
		chromedp.Click(`div.app_tag`, chromedp.NodeVisible), // click button to show all tags
		chromedp.Text(`.app_tags`, &tags, chromedp.NodeVisible),
	)

	if err != nil {
		return nil, fmt.Errorf("getting tags: %w", err)
	}

	if tags == "" {
		return nil, fmt.Errorf("no tags found.")
	}

	tagSlice := strings.Split(tags, "\n")
	return tagSlice, nil
}

// scrapeMetadata scrapes release date, total review count,
// and review positivity
func scrapeMetadata(gameUrl string) (string, string, string, error) { // TODO - Make a game struct for all this data
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	var releaseDate string
	var totalReviewCount string
	var reviewPositivity string
	err := chromedp.Run(ctx,
		// Go to URL for game and get release date
		chromedp.Navigate(gameUrl),
		chromedp.Text(`div.release_date div.date`, &releaseDate, chromedp.NodeVisible),
		chromedp.Text(`div.summary_text span.app_reviews_count`, &totalReviewCount, chromedp.NodeVisible),
		chromedp.Text(`div.summary_text span.game_review_summary`, &reviewPositivity, chromedp.NodeVisible),
	)

	if err != nil {
		return "", "", "", fmt.Errorf("scraping data: %w", err)
	}

	if releaseDate == "" {
		slog.Warn("No release date found", "url", gameUrl)
	}
	if totalReviewCount == "" {
		slog.Warn("No total review count found", "url", gameUrl)
	}
	if reviewPositivity == "" {
		slog.Warn("No review positivity found", "url", gameUrl)
	}

	totalReviewCount = strings.TrimPrefix(totalReviewCount, "(")
	totalReviewCount = strings.TrimSuffix(totalReviewCount, " reviews)")

	return releaseDate, totalReviewCount, reviewPositivity, nil
}
