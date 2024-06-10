package main

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
		log.Fatal(err)
	}
	log.Println("Started telegram bot client ✔")

	// Инициализируем хранилище
	storage, err := s.New(cfg.DatabasePath)
	if err != nil {
		log.Fatal(err)
	}
	err = storage.Init()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to movies.db ✔")

	// Канал для получения обновлений от пользователя
	updatesChan := tgBot.Update()

	for update := range updatesChan {
		if update.Message == nil {
			continue
		}
		go CommandsHandler(update, client, storage, tgBot)
	}

}

func CommandsHandler(update tgbotapi.Update, client *k.KinopoiskClient, storage *s.Storage, tgBot *tg.TgBotClient) {
	msgText := update.Message.Text
	username := update.Message.From.UserName
	chatID := update.Message.Chat.ID
	command := strings.Fields(msgText)[0] // отделяем команду от остального содержания сообщения

	switch command {
	case `/start`:
		err := tgBot.Greeting(chatID)
		if err != nil {
			log.Fatalf("can't send a greeting to user: %s", err)
		}

	case `/random`:
		err := tgBot.Random(client, chatID)
		if err != nil {
			log.Fatalf("can't get a random movie: %s", err)
		}

	case `/randomwithgosling`:
		err := tgBot.RandomWithGosling(client, chatID)
		if err != nil {
			log.Fatalf("can't get a random movie with Gosling: %s", err)
		}

	case `/savemovie`:
		err := tgBot.Save(client, storage, msgText, username, chatID)
		if err != nil {
			log.Fatalf("can't save a movie: %s", err)
		}

	case `/removemovie`:
		err := tgBot.Remove(storage, msgText, username, chatID)
		if err != nil {
			log.Fatalf("can't remove a movie: %s", err)
		}

	case `/mymovies`:
		err := tgBot.Movies(chatID, username, storage)
		if err != nil {
			log.Fatalf("can't show all saved movies: %s", err)
		}

	default:
		err := tgBot.Unrecognized(chatID)
		if err != nil {
			log.Fatalf("can't show `unrecognized` message: %s", err)
		}
	}
}
