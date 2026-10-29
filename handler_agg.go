package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return errors.New("usage: agg <time_between_requests>")
	}

	timeBetweenRequest, err := time.ParseDuration(cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("parse duration: %w", err)
	}

	log.Printf("Collecting feeds every %v\n", timeBetweenRequest)

	ticker := time.NewTicker(timeBetweenRequest)
	defer ticker.Stop()

	ctx := context.Background()

	for {
		if err := scrapeFeeds(s, ctx); err != nil {
			log.Println(err)
		} else {
			log.Println("Collected feed")
		}
		<-ticker.C

	}

}
