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
	"time"
)

type BookController struct {
}

func (c *BookController) GetDate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		utils.SendSuccess(w, time.Now())
	}
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

func (c *BookController) PostBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book model.Book
		br := repository.BookRepository{}
		var error model.ErrorMessage
		err := json.NewDecoder(r.Body).Decode(&book)

		if err != nil || book.Author == "" || book.Title == "" || book.Year == "" {
			error.Message = "One or more fields are missing!"
			utils.SendError(w, http.StatusForbidden, error) /// 400
			return
		}

		s, err := br.CreateABook(db, book)
		if err != nil {
			error.Message = error.Error()
			utils.SendError(w, http.StatusInternalServerError, error) /// 500
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		utils.SendSuccess(w, s)
	}
}

func (c *BookController) UpdateBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book model.Book
		rb := repository.BookRepository{}
		var error model.ErrorMessage
		err := json.NewDecoder(r.Body).Decode(&book)
		if err != nil {
			error.Message = err.Error()
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}
		s, err := rb.UpdateBook(db, &book)
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
		w.Header().Set("Content-Type", "text/plain")
		utils.SendSuccess(w, s)
	}
}

func (c *BookController) DeleteBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		br := repository.BookRepository{}
		var error model.ErrorMessage
		parmas := mux.Vars(r)
		s, err := br.DeleteABook(db, parmas["id"])
		if err != nil {
			if s != "" {
				error.Message = s
				utils.SendError(w, http.StatusNotFound, error) // 404
			} else {
				error.Message = error.Error()
				utils.SendError(w, http.StatusInternalServerError, error) // 500
			}
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		utils.SendSuccess(w, s)
	}
}
