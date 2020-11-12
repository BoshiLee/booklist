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
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   string `json:"year"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

func (e ErrorMessage) Error() string {
	return e.Message
}

var books []Book
var db *sql.DB

func init() {
	gotenv.Load()
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	pqUrl, err := pq.ParseURL(os.Getenv("SQL_URL"))
	logFatal(err)
	db, err = sql.Open("postgres", pqUrl)
	logFatal(err)
	err = db.Ping()

	router := mux.NewRouter()
	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books", postBook).Methods("POST")
	router.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
}

func postBook(w http.ResponseWriter, r *http.Request) {
}

func updateBook(w http.ResponseWriter, r *http.Request) {
}

func deleteBook(w http.ResponseWriter, r *http.Request) {

}

func checkBookIdContainsInBooks(id int) (int, error) {
	for i, book := range books {
		if book.ID == id {
			return i, nil
		}
	}
	return 0, ErrorMessage{
		"Your Book is Not in shelf, please check another id.",
	}
}

func remove(slice []Book, s int) []Book {
	return append(slice[:s], slice[s+1:]...)
}
