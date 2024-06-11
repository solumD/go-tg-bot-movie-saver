package telegram

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	k "github.com/solumD/go-tg-bot-movie-saver/clients/kinopoisk"
	s "github.com/solumD/go-tg-bot-movie-saver/storage"
)

// Клиент тг-бота
type TgBotClient struct {
	Bot *tgbotapi.BotAPI
}

// New создает объект клиента тг-бота
func New(token string) (*TgBotClient, error) {
	b, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("can't connect to TG by token: %s", err)
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
	m := "Привет, я бот для сохранения фильмов с Кинопоиска и не только\n\nДоступные команды:\n\n/random - посоветую случайный фильм\n/randomwithgosling - посоветую случайный фильм с Райаном Гослингом ;)\n/savemovie - сохранить фильм\n/removemovie - удалить фильм из сохраненного\n/mymovies - вывести все сохраненные фильмы"

	mConfig := tgbotapi.NewMessage(chatID, m)
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
	link := "https://www.kinopoisk.ru/film/" + strconv.Itoa(movie.Id)

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
	link := "https://www.kinopoisk.ru/film/" + strconv.Itoa(movie.Id)

	m := fmt.Sprintf("Конечно, вот хороший фильм с Райаном Гослингом\n\nНазвание: \"%s\"\n\nО чем: %s\n\nРейтинг на КП: %.2f\nОграничение по возрасту: %d+\nГод выхода: %d\nДлительность: %d минут\nСсылка на КП: %s",
		movie.Title, movie.Description, movie.Rating.KpRating, movie.Age, movie.Year, movie.Length, link)

	mConfig := tgbotapi.NewMessage(chatID, m)
	if _, err := t.Bot.Send(mConfig); err != nil {
		return err
	}
	return nil
}

// SaveMovie сохраняет фильм в БД
func (t *TgBotClient) Save(c *k.KinopoiskClient, storage *s.Storage, msgText string, username string, chatID int64) error {
	msgText = strings.ToLower(strings.TrimSpace(msgText))
	fields := strings.Fields(msgText)
	if len(fields) < 2 {
		m := "Для сохранения отправьте ссылку фильма на Кинопоиске после вызова команды.\n\nПример: /savemovie https://www.kinopoisk.ru/film/0000/"
		mConfig := tgbotapi.NewMessage(chatID, m)
		if _, err := t.Bot.Send(mConfig); err != nil {
			return err
		}
		return nil
	}

	if len(fields) > 2 {
		m := "Неверный формат команды!\n\nОтправьте ссылку в указанном формате:\n/savemovie https://www.kinopoisk.ru/film/0000/"
		mConfig := tgbotapi.NewMessage(chatID, m)
		if _, err := t.Bot.Send(mConfig); err != nil {
			return err
		}
		return nil
	}

	uri := fields[1]
	if strings.Contains(uri, "kinopoisk.ru") && strings.Contains(uri, "film") && strings.Contains(uri, "https://") {
		uriFields := strings.Split(uri, "/")

		if len(uriFields) > 5 {
			m := "Неверный формат команды!\n\nОтправьте ссылку в указанном формате:\n/savemovie https://www.kinopoisk.ru/film/0000/"
			mConfig := tgbotapi.NewMessage(chatID, m)
			if _, err := t.Bot.Send(mConfig); err != nil {
				return err
			}
			return nil
		}

		movieId := uriFields[4]
		for _, v := range movieId {
			if !unicode.IsDigit(v) {
				m := fmt.Sprintf("Неверное id фильма: %s", movieId)
				mConfig := tgbotapi.NewMessage(chatID, m)
				if _, err := t.Bot.Send(mConfig); err != nil {
					return err
				}
				return nil
			}
		}

		if string(movieId[0]) == "0" {
			m := fmt.Sprintf("Неверное id фильма: %s", movieId)
			mConfig := tgbotapi.NewMessage(chatID, m)
			if _, err := t.Bot.Send(mConfig); err != nil {
				return err
			}
			return nil
		}

		movie, err := c.ById(movieId)
		if err != nil {
			return err
		}
		if len(movie.Title) == 0 {
			m := "Неизвестный фильм. Отсутсвует название"
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

		if err = storage.Save(movieId, movie.Title, username); err != nil {
			return err
		}

		m := fmt.Sprintf("Фильм \"%s\" был успешно сохранен. ", movie.Title)
		mConfig := tgbotapi.NewMessage(chatID, m)
		if _, err := t.Bot.Send(mConfig); err != nil {
			return err
		}

	} else {
		m := "Неверный формат команды!\n\nОтправьте ссылку в указанном формате:\n/savemovie https://www.kinopoisk.ru/film/0000/"
		mConfig := tgbotapi.NewMessage(chatID, m)
		if _, err := t.Bot.Send(mConfig); err != nil {
			return err
		}
	}

	return nil
}

// Remove удаляет фильм из БД по названию
func (t *TgBotClient) Remove(storage *s.Storage, msgText string, username string, chatID int64) error {
	fields := strings.Fields(msgText)
	if len(fields) < 2 {
		m := "Для удаления отправьте название фильма без кавычек, как оно записано в сохраненном.\n/mymovies - показать сохраненные фильмы\n\nПример: /removemovie Драйв"
		mConfig := tgbotapi.NewMessage(chatID, m)
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
		m := "Фильм не найден в сохраненном. Проверьте название фильма.\nПример: /removemovie Драйв"
		mConfig := tgbotapi.NewMessage(chatID, m)
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
func (t *TgBotClient) Movies(chatID int64, username string, storage *s.Storage) error {
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
	uri := "https://www.kinopoisk.ru/film/"
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
	m := "Неизвестная команда"

	mConfig := tgbotapi.NewMessage(chatID, m)
	if _, err := t.Bot.Send(mConfig); err != nil {
		return err
	}
	return nil
}

// InDevelopment сообщает пользователю, что команда в разработке
func (t *TgBotClient) InDevelopment(chatID int64) error {
	m := "В разработке..."

	mConfig := tgbotapi.NewMessage(chatID, m)
	if _, err := t.Bot.Send(mConfig); err != nil {
		return err
	}
	return nil
}
