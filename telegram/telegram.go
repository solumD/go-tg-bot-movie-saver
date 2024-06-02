package telegram

import (
	"fmt"
	"log"
	"strconv"
	"strings"

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
func (t *TgBotClient) SaveMovie(c *k.KinopoiskClient, storage *s.Storage, chatID int64, updateChan tgbotapi.UpdatesChannel) error {
	m := "Отравьте ссылку фильма на Кинопоиске без вспомогательных символов.\nПример: https://www.kinopoisk.ru/film/0000/"
	mConfig := tgbotapi.NewMessage(chatID, m)
	if _, err := t.Bot.Send(mConfig); err != nil {
		return err
	}

	for update := range updateChan {
		msgText := update.Message.Text

		if strings.Contains(msgText, "kinopoisk.ru") && strings.Contains(msgText, "film") {

			movieId := strings.Split(msgText, "/")[4]
			username := update.Message.From.UserName

			movie, err := c.ById(movieId)
			if err != nil {
				log.Fatalf("can't find movie by id: %s", err)
				return err
			}

			exist, err := storage.IsExist(movieId, username)
			if err != nil {
				return err
			}
			if exist {
				m = fmt.Sprintf("Фильм \"%s\" уже был сохранен!", movie.Title)
				mConfig := tgbotapi.NewMessage(chatID, m)
				if _, err := t.Bot.Send(mConfig); err != nil {
					return err
				}
				break
			}

			if err = storage.Save(movieId, movie.Title, username); err != nil {
				return err
			}

			m = fmt.Sprintf("Фильм \"%s\" был успешно сохранен. ", movie.Title)
			mConfig := tgbotapi.NewMessage(chatID, m)
			if _, err := t.Bot.Send(mConfig); err != nil {
				return err
			}
			break

		} else if strings.TrimSpace(strings.ToLower(msgText)) == "отмена" {
			m = "Действие отменено."
			mConfig := tgbotapi.NewMessage(chatID, m)
			if _, err := t.Bot.Send(mConfig); err != nil {
				return err
			}
			break

		} else {
			m = "Ошибка!\nОтправьте ссылку в указанном формате: https://www.kinopoisk.ru/film/0000/\nДля отмены команды отправьте слово \"Отмена\""
			mConfig := tgbotapi.NewMessage(chatID, m)
			if _, err := t.Bot.Send(mConfig); err != nil {
				return err
			}
		}
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
