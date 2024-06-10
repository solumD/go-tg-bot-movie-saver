package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

// New создает объект БД
func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("can't open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to db: %w", err)
	}

	return &Storage{db: db}, nil

}

// Save сохраняет фильм в БД
func (s *Storage) Save(movieId, title, username string) error {
	q := `insert into saved_movies (movie_id, title, user_name) values (?, ?, ?)`

	_, err := s.db.Exec(q, movieId, title, username)
	if err != nil {
		return fmt.Errorf("can't save movie to db: %w", err)
	}

	return nil
}

// Тип для фильма из БД
type movie struct {
	Id    string
	Title string
}

// Pick получает список всех фильмов пользователя из БД
func (s *Storage) Pick(username string) ([]movie, error) {
	q := `select movie_id, title from saved_movies where user_name = ?`

	rows, err := s.db.Query(q, username)
	if err != nil {
		return nil, fmt.Errorf("can't pick movies from db: %w", err)
	}
	defer rows.Close()

	var movies []movie
	for rows.Next() {
		var m movie
		if err := rows.Scan(&m.Id, &m.Title); err != nil {
			return nil, fmt.Errorf("can't select movie from db: %w", err)
		}
		movies = append(movies, m)
	}

	return movies, nil
}

// Remove удаляет конкретный фильм конкретного пользователя из БД по названию фильма
func (s *Storage) Remove(title, username string) error {
	q := `delete from saved_movies where title = ? and user_name = ?`

	_, err := s.db.Exec(q, title, username)
	if err != nil {
		return fmt.Errorf("can't delete movie from db for user %s: %w", username, err)
	}

	return nil
}

// IsExist проверяет, существует ли конкретный фильм у конкретного пользователя по названию фильма
func (s *Storage) IsExistByTitle(title, username string) (bool, error) {
	q := `select count(*) from saved_movies where title = ? and user_name = ?`

	var count int

	if err := s.db.QueryRow(q, title, username).Scan(&count); err != nil {
		return false, fmt.Errorf("can't check if movie exist: %w", err)
	}

	return count > 0, nil
}

// IsExist проверяет, существует ли конкретный фильм у конкретного пользователя по id фильма
func (s *Storage) IsExistById(id, username string) (bool, error) {
	q := `select count(*) from saved_movies where movie_id = ? and user_name = ?`

	var count int

	if err := s.db.QueryRow(q, id, username).Scan(&count); err != nil {
		return false, fmt.Errorf("can't check if movie exist: %w", err)
	}

	return count > 0, nil
}

// Init создает в БД таблицу `saved_movies`, если она не создана
func (s *Storage) Init() error {
	q := `create table if not exists saved_movies (
		id integer primary key autoincrement,
		movie_id text,
		title text,
		user_name text
	)`
	_, err := s.db.Exec(q)
	if err != nil {
		return fmt.Errorf("can't create table saved_movies: %w", err)
	}

	return nil
}
