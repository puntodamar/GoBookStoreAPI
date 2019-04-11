package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/subosito/gotenv"
	"log"
	"net/http"
	"os"

	_ "database/sql"
	_ "github.com/lib/pq"
	_ "github.com/subosito/gotenv"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   string `json:"year"`
}

var books []Book
var db *sql.DB

func init() {
	gotenv.Load()
}

func logFatal(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {

	pgURL, err := pq.ParseURL(os.Getenv("ELEPHANTSQL_URL"))
	logFatal(err)
	db, err = sql.Open("postgres", pgURL)

	logFatal(err)
	db.Ping()
	logFatal(err)
	log.Println(pgURL)

	router := mux.NewRouter()

	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books", addBook).Methods("POST")
	router.HandleFunc("/books", updateBook).Methods("PUT")
	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", removeBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	var book Book
	books = []Book{}

	rows, err := db.Query("SELECT * FROM books")
	logFatal(err)

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		logFatal(err)

		books = append(books, book)
	}

	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
}

func addBook(w http.ResponseWriter, r *http.Request) {
}

func updateBook(w http.ResponseWriter, r *http.Request) {
}

func removeBook(w http.ResponseWriter, r *http.Request) {
}
