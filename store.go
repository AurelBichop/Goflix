package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Store interface {
	Open() error
	Close() error

	GetMovies() ([]*Movie, error)
	GetMovieById(id int64) (*Movie, error)
	CreateMovie(m *Movie) error
	FindUser(username string, password string) (bool, error)
}

type dbStore struct {
	db *sqlx.DB
}

var schemaMovie = `
CREATE TABLE IF NOT EXISTS movie 
(
	id INT NOT NULL AUTO_INCREMENT,
	title VARCHAR(255) NOT NULL ,
	release_date VARCHAR(255) NOT NULL ,
	duration INT NOT NULL ,
	trailer_url VARCHAR(255) NOT NULL ,
	PRIMARY KEY (id)) ENGINE = InnoDB;
`
var userSchema = `
CREATE TABLE IF NOT EXISTS user 
(
	id INT NOT NULL AUTO_INCREMENT,
	username VARCHAR(30) NOT NULL ,
	password VARCHAR(255) NOT NULL ,
	PRIMARY KEY (id)) ENGINE = InnoDB;
`

func (store *dbStore) Open() error {
	db, err := sqlx.Connect("mysql", "root:@tcp(127.0.0.1:3306)/goflix")
	if err != nil {
		return err
	}
	log.Println("Connected to DB")

	db.MustExec(schemaMovie)
	db.MustExec(userSchema)
	store.db = db
	return nil
}

func (store *dbStore) Close() error {
	return store.db.Close()
}

func (store *dbStore) GetMovies() ([]*Movie, error) {
	var movies []*Movie
	err := store.db.Select(&movies, "SELECT * FROM movie")
	if err != nil {
		return movies, err
	}
	return movies, nil
}

func (store *dbStore) GetMovieById(id int64) (*Movie, error) {
	var movie = &Movie{}
	err := store.db.Get(movie, "SELECT * FROM movie WHERE id=?", id)
	if err != nil {
		return movie, err
	}
	return movie, nil
}

func (store *dbStore) CreateMovie(m *Movie) error {
	res, err := store.db.Exec("INSERT INTO movie (title, release_date, duration, trailer_url) VALUES (?,?,?,?)",
		m.Title, m.ReleaseDate, m.Duration, m.TrailerURL)
	if err != nil {
		return err
	}

	m.ID, err = res.LastInsertId()
	return err
}

func (store *dbStore) FindUser(username string, password string) (bool, error) {
	var count int

	err := store.db.Get(&count, "SELECT COUNT(id) FROM user WHERE username=? AND password=?", username, password)
	if err != nil {
		return false, err
	}

	return count == 1, nil
}
