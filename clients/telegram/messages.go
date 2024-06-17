package telegram

const msgGreeting = `Привет, я бот для сохранения фильмов с Кинопоиска и не только

Доступные команды:

/random - посоветую случайный фильм
/randomwithgosling - посоветую случайный фильм с Райаном Гослингом ;)
/savemovie - сохранить фильм
/removemovie - удалить фильм из сохраненного
/mymovies - вывести все сохраненные фильмы`

const msgInvalidSaveCommand = `Неверный формат ссылки!
Отправьте ссылку в указанном формате:
/savemovie https://www.kinopoisk.ru/film/0000`

const msgSave = `Для сохранения отправьте ссылку фильма на Кинопоиске после вызова команды.

Пример: /savemovie https://www.kinopoisk.ru/film/0000`

const msgRemove = `Для удаления отправьте название фильма без кавычек, как оно записано в сохраненном.
/mymovies - показать сохраненные фильмы

Пример: /removemovie Драйв`

const msgRemoveNotFound = `Фильм не найден в сохраненном. Проверьте название фильма.

Пример: /removemovie Драйв`

const msgUnrecognized = "Неизвестная команда"

const msgInDevelopment = "В разработке..."

const kinopoiskMovieLink = "https://www.kinopoisk.ru/film/"
