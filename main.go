package main

import (
	"fmt"
	"log"

	"github.com/solumD/go-tg-bot-movie-saver/storage"
)

func main() {
	s, err := storage.New("./storage/database/movies.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println("connected to db")

	err = s.Init()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = s.Remove("12313", "sassyba")
	if err != nil {
		fmt.Println(err)
		return
	}
	m, err := s.Pick("sassyba")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(m)
	flag, err := s.IsExist("12313", "sassyba")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(flag)
}
