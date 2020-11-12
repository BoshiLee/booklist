package repository

import (
	"bookList/model"
	"database/sql"
	"fmt"
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

func (r *BookRepository) CreateABook(db *sql.DB, book model.Book) (string, error) {
	var id int
	err := db.QueryRow("insert into books (title, author, year) values($1, $2, $3) RETURNING id;", &book.Title, &book.Author, &book.Year).Scan(&id)
	return fmt.Sprintf("Book Id %v has been created", id), err
}

func (r *BookRepository) UpdateBook(db *sql.DB, book *model.Book) (string, error) {
	_, err := db.Exec("update books set title=$1, author=$2, year=$3 where id=$4 RETURNING id;", book.Title, book.Author, book.Year, book.ID)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Book Id %v has been update", book.ID), nil
}

func (r *BookRepository) DeleteABook(db *sql.DB, id string) (string, error) {
	result, err := db.Exec("delete from books where id=$1", id)
	if err != nil {
		return "", err
	}
	if rows, err := result.RowsAffected(); err != nil || rows == 0 {
		return "Resource Not Found", err
	}
	return fmt.Sprintf("Book Id %v has been delete", id), err
}
