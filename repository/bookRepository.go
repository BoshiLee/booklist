package repository

import (
	"bookList/model"
	"database/sql"
)

type BookRepository struct {
}

func (r *BookRepository) GetBooks(db *sql.DB) ([]model.Book, error) {
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

func (r *BookRepository) GetBook(db *sql.DB, id string) (model.Book, error) {
	var book model.Book
	row := db.QueryRow("select * from books where id=$1", id)
	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
	return book, err
}
