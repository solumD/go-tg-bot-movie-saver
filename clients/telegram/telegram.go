package telegram

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	k "github.com/solumD/go-tg-bot-movie-saver/clients/kinopoisk"
	"github.com/solumD/go-tg-bot-movie-saver/storage"
)

// Клиент тг-бота
type TgBotClient struct {
	Bot *tgbotapi.BotAPI
}

// New создает объект клиента тг-бота
func New(token string) (*TgBotClient, error) {
	b, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("can't connect to TG by token: %s", token)
		return nil, err
	}
	return &TgBotClient{Bot: b}, nil
}

// Update получает обновления от пользователя
func (t *TgBotClient) Update() tgbotapi.UpdatesChannel {
	updateConfig := tgbotapi.NewUpdate(0)
	return t.Bot.GetUpdatesChan(updateConfig)
}

// Greeting отправляет сообщение приветствия
func (t *TgBotClient) Greeting(chatID int64) error {
	mConfig := tgbotapi.NewMessage(chatID, msgGreeting)
	if _, err := t.Bot.Send(mConfig); err != nil {
		return err
	}
	return nil
}

// Random получает случайный фильм через клиент Кинопоиска и отправляет его
func (t *TgBotClient) Random(c *k.KinopoiskClient, chatID int64) error {
	movie, err := c.Random()
	if err != nil {
		return err
	}
	link := kinopoiskMovieLink + strconv.Itoa(movie.Id)

	m := fmt.Sprintf("Название: \"%s\"\n\nО чем: %s\n\nРейтинг на КП: %.2f\nОграничение по возрасту: %d+\nГод выхода: %d\nДлительность: %d минут\nСсылка на КП: %s",
		movie.Title, movie.Description, movie.Rating.KpRating, movie.Age, movie.Year, movie.Length, link)

	mConfig := tgbotapi.NewMessage(chatID, m)
	if _, err := t.Bot.Send(mConfig); err != nil {
		return err
	}
	return nil
}

// RandomWithGosling получает случайный фильм с Райаном Гослингом через клиент Кинопоиска и отправляет его
func (t *TgBotClient) RandomWithGosling(c *k.KinopoiskClient, chatID int64) error {
	movie, err := c.RandomWithGosling()
	if err != nil {
		return err
	}
	link := kinopoiskMovieLink + strconv.Itoa(movie.Id)

	m := fmt.Sprintf("Конечно, вот хороший фильм с Райаном Гослингом\n\nНазвание: \"%s\"\n\nО чем: %s\n\nРейтинг на КП: %.2f\nОграничение по возрасту: %d+\nГод выхода: %d\nДлительность: %d минут\nСсылка на КП: %s",
		movie.Title, movie.Description, movie.Rating.KpRating, movie.Age, movie.Year, movie.Length, link)

	mConfig := tgbotapi.NewMessage(chatID, m)
	if _, err := t.Bot.Send(mConfig); err != nil {
		return err
	}
	return nil
}

// SaveMovie сохраняет фильм в БД
func (t *TgBotClient) Save(c *k.KinopoiskClient, storage storage.Storage, msgText string, username string, chatID int64) error {
	msgText = strings.ToLower(strings.TrimSpace(msgText))
	fields := strings.SplitN(msgText, " ", 2)
	if len(fields) < 2 {
		mConfig := tgbotapi.NewMessage(chatID, msgSave)
		if _, err := t.Bot.Send(mConfig); err != nil {
			return err
		}
		return nil
	}

	uri := fields[1]
	if matched, err := regexp.MatchString(`^https:\/\/www\.kinopoisk\.ru\/film\/\d+$`, uri); matched {
		if err != nil {
			mConfig := tgbotapi.NewMessage(chatID, msgInvalidSaveCommand)
			if _, err := t.Bot.Send(mConfig); err != nil {
				return err
			}
			return err
		}

		movieId := strings.Split(uri, "/")[4]

		movie, err := c.ById(movieId)
		if err != nil {
			return err
		}

		if len(movie.Title) == 0 {
			m := "Неизвестный фильм"
			mConfig := tgbotapi.NewMessage(chatID, m)
			if _, err := t.Bot.Send(mConfig); err != nil {
				return err
			}
			return nil
		}

		exist, err := storage.IsExistById(movieId, username)
		if err != nil {
			return err
		}

		if exist {
			m := fmt.Sprintf("Фильм \"%s\" уже был сохранен!", movie.Title)
			mConfig := tgbotapi.NewMessage(chatID, m)
			if _, err := t.Bot.Send(mConfig); err != nil {
				return err
			}
			return nil
		}

		exist, err = storage.IsExistByTitle(movie.Title, username)
		if err != nil {
			return err
		}

		if exist {
			m := fmt.Sprintf("Фильм \"%s\" уже был сохранен!", movie.Title)
			mConfig := tgbotapi.NewMessage(chatID, m)
			if _, err := t.Bot.Send(mConfig); err != nil {
				return err
			}
			return nil
		}

		if err = storage.Save(movieId, movie.Title, username); err != nil {
			return err
		}

		m := fmt.Sprintf("Фильм \"%s\" был успешно сохранен. ", movie.Title)
		mConfig := tgbotapi.NewMessage(chatID, m)
		if _, err := t.Bot.Send(mConfig); err != nil {
			return err
		}

	} else {
		mConfig := tgbotapi.NewMessage(chatID, msgInvalidSaveCommand)
		if _, err := t.Bot.Send(mConfig); err != nil {
			return err
		}
	}

	return nil
}

// Remove удаляет фильм из БД по названию
func (t *TgBotClient) Remove(storage storage.Storage, msgText string, username string, chatID int64) error {
	fields := strings.Fields(msgText)
	if len(fields) < 2 {
		mConfig := tgbotapi.NewMessage(chatID, msgRemove)
		if _, err := t.Bot.Send(mConfig); err != nil {
			return err
		}
		return nil
	}

	title := ""
	for i := 1; i < len(fields); i++ {
		title = title + " " + fields[i]
	}
	title = title[1:]

	exist, err := storage.IsExistByTitle(title, username)
	if err != nil {
		return err
	}
	if !exist {
		mConfig := tgbotapi.NewMessage(chatID, msgRemoveNotFound)
		if _, err := t.Bot.Send(mConfig); err != nil {
			return err
		}
		return nil
	}

	if err := storage.Remove(title, username); err != nil {
		return err
	}
	m := fmt.Sprintf("Фильм \"%s\" был удален. ", title)
	mConfig := tgbotapi.NewMessage(chatID, m)
	if _, err := t.Bot.Send(mConfig); err != nil {
		return err
	}

	return nil
}

// Movies выводит все сохраненные пользователем фильмы
func (t *TgBotClient) Movies(chatID int64, username string, storage storage.Storage) error {
	allMovies, err := storage.Pick(username)
	if err != nil {
		return err
	}
	if len(allMovies) == 0 {
		m := "Вы пока не сохранили ни одного фильма"
		mConfig := tgbotapi.NewMessage(chatID, m)
		if _, err := t.Bot.Send(mConfig); err != nil {
			return err
		}
		return nil
	}

	m := "Сохраненные фильмы\n\n"
	uri := kinopoiskMovieLink
	for idx, movie := range allMovies {
		row := fmt.Sprintf("%d) %s\nСсылка: %s\n\n", idx+1, movie.Title, uri+movie.Id)
		m = m + row
	}
	m = m + "Для удаления фильма воспользуйтесь командой /removemovie"

	mConfig := tgbotapi.NewMessage(chatID, m)
	if _, err := t.Bot.Send(mConfig); err != nil {
		return err
	}

	return nil
}

// Unrecognized отправляет сообщение о неизвестной команде
func (t *TgBotClient) Unrecognized(chatID int64) error {
	mConfig := tgbotapi.NewMessage(chatID, msgUnrecognized)
	if _, err := t.Bot.Send(mConfig); err != nil {
		return err
	}
	return nil
}

// InDevelopment сообщает пользователю, что команда в разработке
func (t *TgBotClient) InDevelopment(chatID int64) error {
	mConfig := tgbotapi.NewMessage(chatID, msgInDevelopment)
	if _, err := t.Bot.Send(mConfig); err != nil {
		return err
	}
	return nil
}
