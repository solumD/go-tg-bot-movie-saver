package main

import (
	"log"

	k "github.com/solumD/go-tg-bot-movie-saver/clients/kinopoisk"
	config "github.com/solumD/go-tg-bot-movie-saver/internal"
	s "github.com/solumD/go-tg-bot-movie-saver/storage"
)

func main() {
	cfg := config.MustLoad()

	client := k.New(
		cfg.Timeout,
		cfg.Uri,
		cfg.ApiToken,
	)

	log.Println("Created client")

	storage, err := s.New(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("can't open storage: %s", err)

	}
	if err = storage.Init(); err != nil {
		log.Fatalf("can't init storage: %s", err)
	}
	log.Println("Connected to database")
}
