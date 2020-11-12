package controller

import (
	"bookList/model"
	"bookList/repository"
	"bookList/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type BookController struct {
}

func (c *BookController) GetBooks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		br := repository.BookRepository{}
		var error model.ErrorMessage
		books, err := br.GetBooks(db)
		if err != nil {
			fmt.Println(err.Error())
			error.Message = "Fetch Books Failed"
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		utils.SendSuccess(w, books)
	}
}

func (c *BookController) GetBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book model.Book
		var error model.ErrorMessage
		rb := repository.BookRepository{}
		params := mux.Vars(r)
		book, err := rb.GetBook(db, params["id"])
		if err != nil {
			if err == sql.ErrNoRows {
				error.Message = "Resource not found"
				utils.SendError(w, http.StatusNotFound, error)
			} else {
				error.Message = error.Error()
				utils.SendError(w, http.StatusInternalServerError, error)
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		utils.SendSuccess(w, book)
	}
}

func (c *BookController) UpdateBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book model.Book
		err := json.NewDecoder(r.Body).Decode(&book)
		if err != nil {
			json.NewEncoder(w).Encode(model.ErrorMessage{
				err.Error(),
			})
			return
		}
		result, err := db.Exec("update books set title=$1, author=$2, year=$3 where id=$4 RETURNING id;", &book.Title, &book.Author, &book.Year, &book.ID)
		if err != nil {
			json.NewEncoder(w).Encode(model.ErrorMessage{
				err.Error(),
			})
			return
		}
		rowsUpdated, err := result.RowsAffected()
		if err != nil {
			json.NewEncoder(w).Encode(model.ErrorMessage{
				err.Error(),
			})
			return
		}
		json.NewEncoder(w).Encode(rowsUpdated)
	}
}

func (c *BookController) PostBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book model.Book
		var id int
		err := json.NewDecoder(r.Body).Decode(&book)
		if err != nil {
			json.NewEncoder(w).Encode(model.ErrorMessage{
				err.Error(),
			})
			return
		}
		err = db.QueryRow("insert into books (title, author, year) values($1, $2, $3) RETURNING id;", &book.Title, &book.Author, &book.Year).Scan(&id)
		if err != nil {
			json.NewEncoder(w).Encode(model.ErrorMessage{
				err.Error(),
			})
			return
		}
		json.NewEncoder(w).Encode(&id)
	}
}

func (c *BookController) DeleteBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parmas := mux.Vars(r)
		result, err := db.Exec("delete from books where id=$1", parmas["id"])
		if err != nil {
			json.NewEncoder(w).Encode(model.ErrorMessage{
				err.Error(),
			})
			return
		}
		rowAffected, err := result.RowsAffected()

		if err != nil {
			json.NewEncoder(w).Encode(model.ErrorMessage{
				err.Error(),
			})
			return
		}
		json.NewEncoder(w).Encode(rowAffected)
	}
}
