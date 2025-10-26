package main

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/sjadczak/gator/internal/rss"
)

func processRSSFeed(r io.Reader) (*rss.RSSFeed, error) {
	var feed rss.RSSFeed
	decoder := xml.NewDecoder(r)
	err := decoder.Decode(&feed)
	if err != nil {
		msg := fmt.Sprintf("could not unmarshal rss feed: %v", err)
		return nil, errors.New(msg)
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i, item := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(item.Title)
		feed.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}

	return &feed, nil
}

func fetchFeed(ctx context.Context, feedUrl string) (*rss.RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedUrl, nil)
	if err != nil {
		msg := fmt.Sprintf("could not create request: %v", err)
		return nil, errors.New(msg)
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		msg := fmt.Sprintf("could not request rss feed: %v", err)
		return nil, errors.New(msg)
	}
	defer res.Body.Close()

	feed, err := processRSSFeed(res.Body)
	if err != nil {
		return nil, err
	}

	return feed, nil
}

func scrapeFeeds(s *state) {
	ctx := context.Background()
	toFetch, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		fmt.Printf("GetNextFeedToFetch error: %v\n", err)
		os.Exit(1)
	}

	toFetch, err = s.db.MarkFeedFetched(ctx, toFetch.ID)
	if err != nil {
		fmt.Printf("MarkFeedFetched error: %v\n", err)
		os.Exit(1)
	}

	feed, err := fetchFeed(ctx, toFetch.Url)
	if err != nil {
		fmt.Printf("fetchFeed error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf(" > gator: Fetching posts from %s\n", toFetch.Name)
	for _, p := range feed.Channel.Item {
		fmt.Printf("    - Post: %s\n", p.Title)
	}
	fmt.Printf("    Found %d posts\n", len(feed.Channel.Item))
}

// func printItem(i rss.RSSItem, l int) {
// 	props := []string{
// 		i.Title,
// 		fmt.Sprintf("Published: %s", i.PubDate),
// 		fmt.Sprintf("Description: %s", i.Description),
// 		fmt.Sprintf("Link: %s", i.Link),
// 	}

// 	for _, prop := range props {
// 		for _, line := range wrapLine(prop, l-1) {
// 			line = strings.ReplaceAll(line, "\n", " ")
// 			fmt.Println(line)
// 		}
// 	}
// }

// func printFeed(f *rss.RSSFeed, l int) {
// 	props := []string{
// 		fmt.Sprintf("Title: %s", f.Channel.Title),
// 		fmt.Sprintf("Description: %s", f.Channel.Description),
// 	}

// 	if f.Channel.Link != "" {
// 		li := fmt.Sprintf("Link: %s", f.Channel.Link)
// 		props = append(props, li)
// 	}

// 	for _, prop := range props {
// 		for _, line := range wrapLine(prop, l) {
// 			line = strings.ReplaceAll(line, "\n", " ")
// 			fmt.Println(line)
// 		}
// 	}

// 	fmt.Println("Recent Posts:")
// 	ic := len(f.Channel.Item)
// 	for i, item := range f.Channel.Item {
// 		printItem(item, l)
// 		if i < ic {
// 			fmt.Println(strings.Repeat("-", l))
// 		}
// 	}
// }

func handleAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		msg := " gator> Usage:\n" +
			" gator agg <time_between_requests>\n" +
			" example: gator agg 15m"
		fmt.Println(msg)
		return ErrInvalidArgs
	}

	step, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		fmt.Printf(" gator> failed to parse `%s` as time.Duration.", cmd.args[0])
		return ErrInvalidArgs
	}

	fmt.Printf(" > gator: collecting feeds every %s...\n", step)
	ticker := time.NewTicker(step)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}
