package main

//imports
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bmizerany/pat"
	_ "github.com/mattn/go-sqlite3"
)

//song structure
type Song struct {
	Artist string `json:"artist"`
	Song   string `json:"song"`
	Genre  string `json:"genre"`
	Length int64  `json:"length"`
}

//total by genre structure
type TotalGenre struct {
	Genre string `json:"genre"`
	Total int64  `json:"total"`
}

//Declarations
type Songs []Song
type TotalsGenre []TotalGenre

var mainDB *sql.DB

func main() {
	// database conecct
	db, errOpenDB := sql.Open("sqlite3", "sources/jrdd.db")
	checkErr(errOpenDB)
	mainDB = db

	//url path
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

// get total of songs by genre
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

//get all songs
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

// get songs by artist name
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

//get songs by song name
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

//get songs by genre
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

// get songs between max and min length
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

//error
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
