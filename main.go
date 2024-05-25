package main

import (
	"fmt"

	k "github.com/solumD/go-tg-bot-movie-saver/clients/kinopoisk"
)

func main() {
	client := k.New()
	m, err := client.RandomMovie()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Random")
	fmt.Println("Age", m.Age)
	fmt.Println("Title", m.Title)
	fmt.Println("Description", m.Description)
	fmt.Println("Length", m.Length)
	fmt.Println("Id", m.Id)
	fmt.Println("Year", m.Year)
	fmt.Println("Rating", m.Rating.KpRating)
	m, err = client.RandomMovieWithGosling()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println()
	fmt.Println("Gosling")
	fmt.Println("Age", m.Age)
	fmt.Println("Title", m.Title)
	fmt.Println("Description", m.Description)
	fmt.Println("Length", m.Length)
	fmt.Println("Id", m.Id)
	fmt.Println("Year", m.Year)
	fmt.Println("Rating", m.Rating.KpRating)
	m, err = client.MovieById("1309570")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println()
	fmt.Println("ById")
	fmt.Println("Age", m.Age)
	fmt.Println("Title", m.Title)
	fmt.Println("Description", m.Description)
	fmt.Println("Length", m.Length)
	fmt.Println("Id", m.Id)
	fmt.Println("Year", m.Year)
	fmt.Println("Rating", m.Rating.KpRating)
}
