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
	"strings"

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

func printItem(i rss.RSSItem, l int) {
	props := []string{
		i.Title,
		fmt.Sprintf("Published: %s", i.PubDate),
		fmt.Sprintf("Description: %s", i.Description),
		fmt.Sprintf("Link: %s", i.Link),
	}

	for _, prop := range props {
		for _, line := range wrapLine(prop, l-1) {
			line = strings.ReplaceAll(line, "\n", " ")
			fmt.Println(line)
		}
	}
}

func printFeed(f *rss.RSSFeed, l int) {
	props := []string{
		fmt.Sprintf("Title: %s", f.Channel.Title),
		fmt.Sprintf("Description: %s", f.Channel.Description),
	}

	if f.Channel.Link != "" {
		li := fmt.Sprintf("Link: %s", f.Channel.Link)
		props = append(props, li)
	}

	for _, prop := range props {
		for _, line := range wrapLine(prop, l) {
			line = strings.ReplaceAll(line, "\n", " ")
			fmt.Println(line)
		}
	}

	fmt.Println("Recent Posts:")
	ic := len(f.Channel.Item)
	for i, item := range f.Channel.Item {
		printItem(item, l)
		if i < ic {
			fmt.Println(strings.Repeat("-", l))
		}
	}
}

func handleAgg(s *state, cmd command) error {
	fmt.Println(" gator> fetching rss feed...")
	url := "https://www.wagslane.dev/index.xml"
	ctx := context.Background()
	feed, err := fetchFeed(ctx, url)
	if err != nil {
		fmt.Printf("handleAgg err: %v", err)
		os.Exit(1)
	}

	// printFeed(feed, 80)
	fmt.Printf("%v\n", feed)
	return nil
}
