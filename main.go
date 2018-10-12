package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bmizerany/pat"
	_ "github.com/mattn/go-sqlite3"
)

type Song struct {
	Artist string `json:"artist"`
	Song   string `json:"song"`
	Genre  string `json:"genre"`
	Length int64  `json:"length"`
}

<<<<<<< HEAD
type TotalGenre struct {
	Genre string `json:"genre"`
	Total int64  `json:"total"`
}

type Songs []Song
type TotalsGenre []TotalGenre
=======
type Songs []Song
>>>>>>> development

var mainDB *sql.DB

func main() {

	db, errOpenDB := sql.Open("sqlite3", "sources/jrdd.db")
	checkErr(errOpenDB)
	mainDB = db

	r := pat.New()
	r.Get("/songs/", http.HandlerFunc(getAll))
	r.Get("/artist/:artist", http.HandlerFunc(getByArtist))
	r.Get("/song/:song", http.HandlerFunc(getBySong))
	r.Get("/genre/:genre", http.HandlerFunc(getByGenre))
	r.Get("/length/:min/:max", http.HandlerFunc(getByLength))
	r.Get("/totals/", http.HandlerFunc(getTotal))

	http.Handle("/", r)

	log.Print(" Running on 12345")
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func getTotal(w http.ResponseWriter, r *http.Request) {
	rows, err := mainDB.Query("SELECT count(songs.song) as total, genres.name as genre FROM songs inner join genres on genres.id = songs.genre group by songs.genre")
	checkErr(err)
	var totalsGenre TotalsGenre
	for rows.Next() {
		var totalGenre TotalGenre
		err = rows.Scan(&totalGenre.Total, &totalGenre.Genre)
		checkErr(err)
		totalsGenre = append(totalsGenre, totalGenre)
	}
	jsonB, errMarshal := json.Marshal(totalsGenre)
	checkErr(errMarshal)
	fmt.Fprintf(w, "%s", string(jsonB))
}

func getAll(w http.ResponseWriter, r *http.Request) {
	rows, err := mainDB.Query("SELECT songs.artist, songs.song, genres.name as genre, songs.length FROM songs inner join genres on genres.id = songs.genre")
	checkErr(err)
	var songs Songs
	for rows.Next() {
		var song Song
		err = rows.Scan(&song.Artist, &song.Song, &song.Genre, &song.Length)
		checkErr(err)
		songs = append(songs, song)
	}
	jsonB, errMarshal := json.Marshal(songs)
	checkErr(errMarshal)
	fmt.Fprintf(w, "%s", string(jsonB))
}

func getByArtist(w http.ResponseWriter, r *http.Request) {
	artist := r.URL.Query().Get(":artist")
	rows, err := mainDB.Query("SELECT songs.artist, songs.song, genres.name as genre, songs.length FROM songs inner join genres on genres.id = songs.genre where songs.artist like ?", "%"+artist+"%")
	checkErr(err)
	var songs Songs
	for rows.Next() {
		var song Song
		err = rows.Scan(&song.Artist, &song.Song, &song.Genre, &song.Length)
		checkErr(err)
		songs = append(songs, song)
	}
	jsonB, errMarshal := json.Marshal(songs)
	checkErr(errMarshal)
	fmt.Fprintf(w, "%s", string(jsonB))
}

func getBySong(w http.ResponseWriter, r *http.Request) {
	searchSong := r.URL.Query().Get(":song")
	rows, err := mainDB.Query("SELECT songs.artist, songs.song, genres.name as genre, songs.length FROM songs inner join genres on genres.id = songs.genre where songs.song like ?", "%"+searchSong+"%")
	checkErr(err)
	var songs Songs
	for rows.Next() {
		var song Song
		err = rows.Scan(&song.Artist, &song.Song, &song.Genre, &song.Length)
		checkErr(err)
		songs = append(songs, song)
	}
	jsonB, errMarshal := json.Marshal(songs)
	checkErr(errMarshal)
	fmt.Fprintf(w, "%s", string(jsonB))
}

func getByGenre(w http.ResponseWriter, r *http.Request) {
	genre := r.URL.Query().Get(":genre")
	rows, err := mainDB.Query("SELECT songs.artist, songs.song, genres.name as genre, songs.length FROM songs inner join genres on genres.id = songs.genre where genres.name like ?", "%"+genre+"%")
	checkErr(err)
	var songs Songs
	for rows.Next() {
		var song Song
		err = rows.Scan(&song.Artist, &song.Song, &song.Genre, &song.Length)
		checkErr(err)
		songs = append(songs, song)
	}
	jsonB, errMarshal := json.Marshal(songs)
	checkErr(errMarshal)
	fmt.Fprintf(w, "%s", string(jsonB))
}

func getByLength(w http.ResponseWriter, r *http.Request) {
	min := r.URL.Query().Get(":min")
	max := r.URL.Query().Get(":max")
	rows, err := mainDB.Query("SELECT songs.artist, songs.song, genres.name as genre, songs.length FROM songs inner join genres on genres.id = songs.genre where songs.length > ? AND songs.length  < ?", min, max)
	checkErr(err)
	var songs Songs
	for rows.Next() {
		var song Song
		err = rows.Scan(&song.Artist, &song.Song, &song.Genre, &song.Length)
		checkErr(err)
		songs = append(songs, song)
	}
	jsonB, errMarshal := json.Marshal(songs)
	checkErr(errMarshal)
	fmt.Fprintf(w, "%s", string(jsonB))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
