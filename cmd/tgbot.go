package cmd

import (
	"log"

	c "github.com/solumD/go-tg-bot-movie-saver/clients/kinopoisk"
	s "github.com/solumD/go-tg-bot-movie-saver/storage"
)

func main() {
	client := c.New()
	log.Println("Created client")

	storage, err := s.New("./storage/movies.db")
	if err != nil {
		log.Fatalf("can't open storage: %w", err)

	}
	if err = storage.Init(); err != nil {
		log.Fatalf("can't init storage: %w", err)
	}
	log.Println("Connected to database")

}
