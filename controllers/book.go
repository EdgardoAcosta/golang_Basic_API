package controllers

import (
	"books-list/models"
	"books-list/repository/book"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Controller struct{}

var books []models.Book

func logFatal(err error) {
	if err != nil {
		// Stop execution if error
		log.Fatal(err)
	}
}

func (c Controller) GetBooks(db *sql.DB) http.HandlerFunc {
	// Get all books from table
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Get all book")
		var book models.Book

		books = []models.Book{}
		bookRepo := bookRepository.BookRepository{}

		books = bookRepo.GetBooks(db, book, books)

		json.NewEncoder(w).Encode(books)
	}
}

func (c Controller) GetBook(db *sql.DB) http.HandlerFunc {
	// Get single record by ID
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Get single book")
		var book models.Book

		// Read value from url
		params := mux.Vars(r)
		// pars str to int
		id, _ := strconv.Atoi(params["id"])
		//log.Println(reflect.TypeOf(id))
		bookRepo := bookRepository.BookRepository{}

		book = bookRepo.GetBook(db, book, id)

		json.NewEncoder(w).Encode(book)

	}
}

func (c Controller) AddBook(db *sql.DB) http.HandlerFunc {
	// Add new record to DB
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Add a book")

		var book models.Book
		var bookID int

		// Decode POST body
		_ = json.NewDecoder(r.Body).Decode(&book)

		bookRepo := bookRepository.BookRepository{}

		bookID = bookRepo.AddBook(db, book)

		json.NewEncoder(w).Encode(bookID)
	}
}

func (c Controller) UpdateBook(db *sql.DB) http.HandlerFunc {
	// Update single book
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("update a book")

		var book models.Book

		json.NewDecoder(r.Body).Decode(&book)

		bookRepo := bookRepository.BookRepository{}

		rowsUpdated := bookRepo.UpdateBook(db, book)

		json.NewEncoder(w).Encode(rowsUpdated)

	}
}

func (c Controller) DeleteBook(db *sql.DB) http.HandlerFunc {
	// Remove book by id
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("remove book")

		// Read param from url
		params := mux.Vars(r)

		id, _ := strconv.Atoi(params["id"])
		bookRepo := bookRepository.BookRepository{}

		rowsDeleted := bookRepo.DeleteBook(db, id)

		json.NewEncoder(w).Encode(rowsDeleted)
	}

}
