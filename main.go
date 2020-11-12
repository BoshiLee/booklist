package main

import (
	"bookList/controller"
	"bookList/dirver"
	"database/sql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
	"log"
	"net/http"
)

var db *sql.DB

func init() {
	gotenv.Load()
}

func main() {
	db = dirver.ConnectToDB()
	router := mux.NewRouter()
	bc := controller.BookController{}
	router.HandleFunc("/books", bc.GetBooks(db)).Methods("GET")
	router.HandleFunc("/books/{id}", bc.GetBook(db)).Methods("GET")
	router.HandleFunc("/books", bc.PostBook(db)).Methods("POST")
	router.HandleFunc("/books", bc.UpdateBook(db)).Methods("PUT")
	router.HandleFunc("/books/{id}", bc.DeleteBook(db)).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}))(router)))
}
