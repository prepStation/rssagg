package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/prepStation/rssagg/internal/database"
)

// concurrency is to indicatee the startScraping function how many go routines
// we want to use to fetch all of these different feeds. The point is that we can fetch
// them at the same time.
func startScrapping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequest time.Duration,
) {
	log.Printf("Scrapping on  %v goroutines every %v duration \n", concurrency, timeBetweenRequest)

	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Printf("error fetching feeds :%v\n", err)
			continue
		}

		//synchronisation mechanism
		wg := &sync.WaitGroup{}

		//we are iterating over the same go routine as same as the startscarpping.
		for _, feed := range feeds {
			//on the main go routine we are adding one to the waitgroup for every feed.
			//so if we have a concurrency of 30  we are going to add 30 go routines to the waitgroup.
			//
			wg.Add(1)

			go scrapFeed(db, wg, feed)
		}
		//waiting for waitgroups for 30 distincts calls to done.
		wg.Wait()
	}
}

func scrapFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Error marking feed as fetched %v\n", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Printf("Error fetching feed %v\n", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {

		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("couldn't parse date %v eith error %v\n", item.PubDate, err)
			continue
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: description,
			PublishedAt: pubAt,
			Url:         item.Link,
			FeedID:      feed.ID,
		})

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Printf("failed to create post %v\n", err)
		}

	}
	log.Printf("feedd %v collected, %v posts found\n", feed.Name, len(rssFeed.Channel.Item))

}
