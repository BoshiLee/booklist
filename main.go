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
	router.HandleFunc("/books", updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	books, err := getAllBooks()
	if err != nil {
		json.NewEncoder(w).Encode(ErrorMessage{
			err.Error(),
		})
		return
	}
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	params := mux.Vars(r)
	row := db.QueryRow("select * from books where id=$1", params["id"])
	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorMessage{
			err.Error(),
		})
		return
	}
	json.NewEncoder(w).Encode(book)
}

func postBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	var id int
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorMessage{
			err.Error(),
		})
		return
	}
	err = db.QueryRow("insert into books (title, author, year) values($1, $2, $3) RETURNING id;", &book.Title, &book.Author, &book.Year).Scan(&id)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorMessage{
			err.Error(),
		})
		return
	}
	json.NewEncoder(w).Encode(&id)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorMessage{
			err.Error(),
		})
		return
	}
	result, err := db.Exec("update books set title=$1, author=$2, year=$3 where id=$4 RETURNING id;", &book.Title, &book.Author, &book.Year, &book.ID)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorMessage{
			err.Error(),
		})
		return
	}
	rowsUpdated, err := result.RowsAffected()
	if err != nil {
		json.NewEncoder(w).Encode(ErrorMessage{
			err.Error(),
		})
		return
	}
	json.NewEncoder(w).Encode(rowsUpdated)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	parmas := mux.Vars(r)
	result, err := db.Exec("delete from books where id=$1", parmas["id"])
	if err != nil {
		json.NewEncoder(w).Encode(ErrorMessage{
			err.Error(),
		})
		return
	}
	rowAffected, err := result.RowsAffected()

	if err != nil {
		json.NewEncoder(w).Encode(ErrorMessage{
			err.Error(),
		})
		return
	}
	json.NewEncoder(w).Encode(rowAffected)
}

func getAllBooks() ([]Book, error) {
	books := []Book{}
	rows, err := db.Query("select * from books")
	if err != nil {
		return books, err
	}
	defer rows.Close()
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		if err != nil {
			return books, err
		}
		books = append(books, book)
	}
	return books, nil
}

func checkBookIdContainsInBooks(id int, books []Book) (int, error) {
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
