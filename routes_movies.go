package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type jsonMovie struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	ReleaseDate string `json:"release_date"`
	Duration    int    `json:"duration"`
	TrailerURL  string `json:"trailer_url"`
}

func (s *server) handleMovieList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		movies, err := s.store.GetMovies()
		if err != nil {
			log.Printf("Cannot load movies err=%v\n", err)

			//TODO handle response to the client
			s.respond(w, r, nil, http.StatusInternalServerError)
		}

		var resp = make([]jsonMovie, len(movies))
		for i, m := range movies {
			resp[i] = mapMovieToJson(m)
		}

		//TODO response JSON Format
		s.respond(w, r, resp, http.StatusOK)
	}
}

func (s *server) handleMovieDetail() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			log.Printf("Cannot parse id to int. err=%v", err)
			s.respond(w, r, nil, http.StatusBadRequest)
			return
		}

		movie, err := s.store.GetMovieById(id)
		if err != nil {
			log.Printf("Cannot load movie err=%v\n", err)

			s.respond(w, r, nil, http.StatusInternalServerError)
			return
		}

		var resp = mapMovieToJson(movie)

		s.respond(w, r, resp, http.StatusOK)
	}
}

func (s *server) handleMovieCreate() http.HandlerFunc {
	type request struct {
		Title       string `json:"title"`
		ReleaseDate string `json:"release_date"`
		Duration    int    `json:"duration"`
		TrailerURL  string `json:"trailer_url"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}
		err := s.decode(w,r,&req)
		if err !=nil{
			log.Printf("Cannot parse movie. err=%v", err)
			s.respond(w,r,nil,http.StatusBadRequest)
			return
		}

		//Create a movie
		m := &Movie{
			ID:0,
			Title:req.Title,
			ReleaseDate:req.ReleaseDate,
			Duration:req.Duration,
			TrailerURL:req.TrailerURL,
		}

		//Store the movi in the database
		err = s.store.CreateMovie(m)
		if err != nil{
			log.Printf("Cannot create movie. err=%v", err)
			s.respond(w,r,nil,http.StatusInternalServerError)
			return
		}

		var resp = mapMovieToJson(m)
		s.respond(w,r,resp,http.StatusOK)
	}
}

func mapMovieToJson(m *Movie) jsonMovie {
	return jsonMovie{
		ID:          m.ID,
		Title:       m.Title,
		ReleaseDate: m.ReleaseDate,
		Duration:    m.Duration,
		TrailerURL:  m.TrailerURL,
	}
}

