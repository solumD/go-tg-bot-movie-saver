package main

import (
	"log"

	k "github.com/solumD/go-tg-bot-movie-saver/clients/kinopoisk"
	config "github.com/solumD/go-tg-bot-movie-saver/internal"
	s "github.com/solumD/go-tg-bot-movie-saver/storage"
	tg "github.com/solumD/go-tg-bot-movie-saver/telegram"
)

func main() {

	// Загружаем конфиг
	cfg := config.MustLoad()

	// Инициализируем клиент кинопоиска
	client := k.New(cfg.Timeout, cfg.Uri, cfg.ApiToken)
	log.Println("Started kinopoisk client ✔")

	// Инициализируем клиент тг-бота
	tgBot, err := tg.New(cfg.BotToken)
	if err != nil {
		log.Fatal()
	}
	log.Println("Started telegram bot client ✔")

	//
	storage, err := s.New(cfg.DatabasePath)
	if err != nil {
		log.Fatal()
	}

	// Канал для получения обновлений от пользователя
	updatatesChan := tgBot.Update()

	for update := range updatatesChan {
		switch update.Message.Text {

		case `/start`:
			chatId := update.Message.Chat.ID
			err := tgBot.Greeting(chatId)
			if err != nil {
				log.Fatalf("can't greet a user: %s", err)
			}

		case `/random`:
			chatId := update.Message.Chat.ID
			err := tgBot.Random(client, chatId)
			if err != nil {
				log.Fatalf("can't get a random movie: %s", err)
			}

		case `/randomwithgosling`:
			chatId := update.Message.Chat.ID
			err := tgBot.RandomWithGosling(client, chatId)
			if err != nil {
				log.Fatalf("can't get a random movie: %s", err)
			}

		case `/savemovie`:
			chatId := update.Message.Chat.ID
			err := tgBot.SaveMovie(client, storage, chatId, updatatesChan)
			if err != nil {
				log.Fatalf("can't save a movie: %s", err)
			}

		case `/removemovie`:
			chatId := update.Message.Chat.ID
			err := tgBot.InDevelopment(chatId)
			if err != nil {
				log.Fatalf("can't remove a movie: %s", err)
			}

		case `/mymovies`:
			chatId := update.Message.Chat.ID
			err := tgBot.InDevelopment(chatId)
			if err != nil {
				log.Fatalf("can't show all saved movies: %s", err)
			}

		default:
			chatId := update.Message.Chat.ID
			err := tgBot.Unrecognized(chatId)
			if err != nil {
				log.Fatalf("can't show all saved movies: %s", err)
			}
		}
	}

}
