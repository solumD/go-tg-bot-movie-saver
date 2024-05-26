package client

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	r "github.com/solumD/go-tg-bot-movie-saver/clients/kinopoisk/request"
	m "github.com/solumD/go-tg-bot-movie-saver/clients/kinopoisk/types"
)

const (
	RandomMovieEP            = "v1.4/movie/random"
	RandomMovieWithGoslingEP = "v1.4/movie/random"
	MovieByIdEP              = "v1.4/movie/"
)

type KinopoiskClient struct {
	Client   *http.Client
	Uri      string
	APIToken string
}

func New() *KinopoiskClient {

	return &KinopoiskClient{
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
		Uri:      "https://api.kinopoisk.dev/",
		APIToken: "C9PJMWY-QH34X1R-KXKTE6C-C11YG86",
	}
}

func (k KinopoiskClient) RandomMovie() (*m.Movie, error) {
	params := map[string]string{"limit": "1", "type": "movie", "rating.kp": "7-10", "lists": "top250"}
	req, err := r.CreateRequest(k.Uri, RandomMovieEP, k.APIToken, params)
	if err != nil {
		log.Fatalf("Endpoint: %s, error: %s", RandomMovieEP, err)
		return nil, err
	}

	res, err := k.Client.Do(req)
	if err != nil {
		log.Fatalf("Endpoint: %s, error: %s", RandomMovieEP, err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Endpoint: %s, error: %s", RandomMovieEP, err)
		return nil, err
	}

	var movie m.Movie
	err = json.Unmarshal(body, &movie)
	if err != nil {
		log.Fatalf("Gosling, Endpoint: %s, error: %s", RandomMovieEP, err)
		return nil, err
	}

	return &movie, nil
}

func (k KinopoiskClient) RandomMovieWithGosling() (*m.Movie, error) {
	params := map[string]string{"limit": "1", "type": "movie", "rating.kp": "7-10", "persons.id": "10143"}
	req, err := r.CreateRequest(k.Uri, RandomMovieWithGoslingEP, k.APIToken, params)
	if err != nil {
		log.Fatalf("Gosling, Endpoint: %s, error: %s", RandomMovieWithGoslingEP, err)
		return nil, err
	}

	res, err := k.Client.Do(req)
	if err != nil {
		log.Fatalf("Gosling, Endpoint: %s, error: %s", RandomMovieWithGoslingEP, err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Gosling, Endpoint: %s, error: %s", RandomMovieWithGoslingEP, err)
		return nil, err
	}

	var movie m.Movie
	err = json.Unmarshal(body, &movie)
	if err != nil {
		log.Fatalf("Gosling, Endpoint: %s, error: %s", RandomMovieWithGoslingEP, err)
		return nil, err
	}

	return &movie, nil
}

func (k KinopoiskClient) MovieById(id string) (*m.Movie, error) {
	req, err := r.CreateRequest(k.Uri, MovieByIdEP+id, k.APIToken, nil)
	if err != nil {
		log.Fatalf("Endpoint: %s, error: %s", MovieByIdEP, err)
		return nil, err
	}
	res, err := k.Client.Do(req)
	if err != nil {
		log.Fatalf("Endpoint: %s, error: %s", MovieByIdEP, err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Endpoint: %s, error: %s", MovieByIdEP, err)
		return nil, err
	}

	var movie m.Movie
	err = json.Unmarshal(body, &movie)
	if err != nil {
		log.Fatalf("Endpoint: %s, error: %s", MovieByIdEP, err)
		return nil, err
	}

	return &movie, nil
}
