package main

import (
	"./driver"
	"./models"

	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"

	_ "database/sql"
	_ "github.com/lib/pq"
	_ "github.com/subosito/gotenv"
)

var books []models.Book
var db *sql.DB

func init() {
	gotenv.Load()
}

func main() {
	db = driver.ConnectDB()
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
	var book models.Book
	books = []models.Book{}

	rows, err := db.Query("SELECT * FROM books")
	driver.LogFatal(err)

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		driver.LogFatal(err)

		books = append(books, book)
	}

	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	params := mux.Vars(r)

	rows := db.QueryRow(
		"SELECT * FROM books WHERE id=$1",
		params["id"])

	err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
	driver.LogFatal(err)

	json.NewEncoder(w).Encode(book)
}

func addBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	var bookID int

	json.NewDecoder(r.Body).Decode(&book)
	err := db.QueryRow(
		"INSERT INTO books (title, author, year) "+
			"VALUES ($1, $2, $3) RETURNING id;",
		book.Title, book.Author, book.Year).
		Scan(&bookID)

	driver.LogFatal(err)

	json.NewEncoder(w).Encode(bookID)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book

	json.NewDecoder(r.Body).Decode(&book)

	result, err := db.Exec(
		"UPDATE books SET title=$1, author=$2, year=$3 "+
			"WHERE id=$4 RETURNING ID",
		book.Title, book.Author, book.Year, book.ID)

	driver.LogFatal(err)
	rowsUpdated, err := result.RowsAffected()
	driver.LogFatal(err)

	json.NewEncoder(w).Encode(rowsUpdated)

}

func removeBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	result, err := db.Exec("DELETE FROM books WHERE id = $1", params["id"])
	driver.LogFatal(err)

	rowsDeleted, err := result.RowsAffected()
	driver.LogFatal(err)

	json.NewEncoder(w).Encode(rowsDeleted)

}
