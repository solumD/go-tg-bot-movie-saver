package storage

// Интерфейс для работы с хранилищем фильмов
type Storage interface {
	Save(movieId, title, username string) error
	Pick(username string) ([]Movie, error)
	Remove(title, username string) error
	IsExistByTitle(title, username string) (bool, error)
	IsExistById(id, username string) (bool, error)
}

// Тип для фильма из БД
type Movie struct {
	Id    string
	Title string
}
