package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
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

func main() {
	books = append(books,
		Book{ID: 1, Title: "Golang pointers", Author: "Mr. Golang", Year: "2010"},
		Book{ID: 2, Title: "Goroutines", Author: "Mr. Goroutine", Year: "2011"},
		Book{ID: 3, Title: "Golang routers", Author: "Mr. Router", Year: "2012"},
		Book{ID: 4, Title: "Golang concurrency", Author: "Mr. Currency", Year: "2013"},
		Book{ID: 5, Title: "Golang good parts", Author: "Mr. Good", Year: "2014"})
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
	prams := mux.Vars(r)
	bookId, error := strconv.Atoi(prams["id"])
	if error != nil {
		json.NewEncoder(w).Encode(ErrorMessage{
			Message: "Please put correct id.",
		})
		return
	}
	book, errMsg := checkBookIdContainsInBooks(bookId)
	if errMsg != nil {
		json.NewEncoder(w).Encode(errMsg)
		return
	}
	json.NewEncoder(w).Encode(book)
}

func postBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorMessage{Message: "Encode book failed."})
		return
	}
	books = append(books, book)
	json.NewEncoder(w).Encode(books)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorMessage{Message: "Encode book failed."})
		return
	}
	prams := mux.Vars(r)
	bookId, error := strconv.Atoi(prams["id"])
	if error != nil {
		json.NewEncoder(w).Encode(ErrorMessage{
			Message: "Please put correct id.",
		})
		return
	}
	_, errMsg := checkBookIdContainsInBooks(bookId)
	if errMsg != nil {
		json.NewEncoder(w).Encode(errMsg)
		return
	}
	books[bookId] = book
	json.NewEncoder(w).Encode(&book)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	var params = mux.Vars(r)
	bookId, err := strconv.Atoi(params["id"])
	if err != nil {
		json.NewEncoder(w).Encode(ErrorMessage{
			"Please put correct id.",
		})
		return
	}
	index, errMsg := checkBookIdContainsInBooks(bookId)
	if errMsg != nil {
		json.NewEncoder(w).Encode(errMsg)
		return
	}
	books = remove(books, index)
	json.NewEncoder(w).Encode(books)
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
