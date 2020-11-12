package controller

import (
	"bookList/model"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type BookController struct {
}

func (c *BookController) GetBooks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		books, err := getAllBooks(db)
		if err != nil {
			json.NewEncoder(w).Encode(model.ErrorMessage{
				err.Error(),
			})
			return
		}
		json.NewEncoder(w).Encode(books)
	}
}

func (c *BookController) GetBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book model.Book
		params := mux.Vars(r)
		row := db.QueryRow("select * from books where id=$1", params["id"])
		err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		if err != nil {
			json.NewEncoder(w).Encode(model.ErrorMessage{
				err.Error(),
			})
			return
		}
		json.NewEncoder(w).Encode(book)
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

func getAllBooks(db *sql.DB) ([]model.Book, error) {
	books := []model.Book{}
	rows, err := db.Query("select * from books")
	if err != nil {
		return books, err
	}
	defer rows.Close()
	for rows.Next() {
		var book model.Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		if err != nil {
			return books, err
		}
		books = append(books, book)
	}
	return books, nil
}
