package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/chromedp/chromedp"
)

func scrapeUserTags(gameUrl string) ([]string, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var result string
	err := chromedp.Run(ctx,
		chromedp.Navigate(gameUrl),
		chromedp.Click(`div.app_tag`, chromedp.NodeVisible),
		chromedp.Text(`.app_tags`, &result, chromedp.NodeVisible),
	)

	if err != nil {
		return nil, fmt.Errorf("getting tags: %w", err)
	}

	if result == "" {
		return nil, fmt.Errorf("no tags found.")
	}

	tags := strings.Split(result, "\n")
	return tags, nil
}
