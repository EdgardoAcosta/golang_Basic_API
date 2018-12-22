package bookRepository

import (
	"books-list/models"
	"database/sql"
	"log"
)

type BookRepository struct{}

func logFatal(err error) {
	if err != nil {
		// Stop execution if error
		log.Fatal(err)
	}
}

func (b BookRepository) GetBooks(db *sql.DB, book models.Book, books []models.Book) []models.Book {
	rows, err := db.Query("SELECT * FROM books")
	logFatal(err)

	defer rows.Close()

	for rows.Next() {
		// map each row to a instance of Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		logFatal(err)
		books = append(books, book)
	}

	return books

}

func (b BookRepository) GetBook(db *sql.DB, book models.Book, id int) models.Book {

	rows := db.QueryRow("SELECT * FROM books WHERE id=$1", id)

	// map result to a Book
	err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
	// Stop if error
	logFatal(err)

	return book
}

func (b BookRepository) AddBook(db *sql.DB, book models.Book) int {
	var bookID int

	err := db.QueryRow("INSERT INTO books (title, author, year) VALUES($1, $2, $3) RETURNING id",
		book.Title, book.Author, book.Year).Scan(&bookID)
	logFatal(err)
	return bookID

}

func (b BookRepository) UpdateBook(db *sql.DB, book models.Book) int64 {

	result, err := db.Exec("UPDATE books SET title=$1, author=$2, year=$3 WHERE id=$4 RETURNING id",
		book.Title, book.Author, book.Year, book.ID)

	// Check number of rows affected
	rowsUpdated, err := result.RowsAffected()
	logFatal(err)

	return rowsUpdated

}

func (b BookRepository) DeleteBook(db *sql.DB, id int) int64 {

	result, err := db.Exec("DELETE FROM books WHERE id = $1", id)
	logFatal(err)

	rowsDeleted, err := result.RowsAffected()
	logFatal(err)

	return rowsDeleted

}
