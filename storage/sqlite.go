package storage

import (
	"database/sql"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

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

func (s *Storage) Save(movieId int, title, username string) {
	query := `insert into movies(movie_id, title, user_name) values(?, ?, ?)`

}
