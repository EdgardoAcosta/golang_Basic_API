package main

import (
	"books-list/controllers"
	"books-list/driver"
	"books-list/models"
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
	"log"
	"net/http"
)

var books []models.Book
var db *sql.DB

func init() {
	gotenv.Load()
}

func main() {
	// start connection
	db = driver.ConnectDB()

	router := mux.NewRouter()

	controller := controllers.Controller{}

	// Append test data
	//books = append(books,
	//	Book{ID: 1, Title:"Book 1", Author:"Author 1", Year:"2018/01/10"},
	//	Book{ID: 2, Title:"Book 2", Author:"Author 2", Year:"2018/02/04"})

	// Get all books
	router.HandleFunc("/books/{id}", controller.GetBook(db)).Methods("GET")
	router.HandleFunc("/books", controller.GetBooks(db)).Methods("GET")
	router.HandleFunc("/books", controller.AddBook(db)).Methods("POST")
	router.HandleFunc("/books", controller.UpdateBook(db)).Methods("PUT")
	router.HandleFunc("/books/{id}", controller.DeleteBook(db)).Methods("DELETE")

	// Start server
	log.Fatal(http.ListenAndServe(":8080", router))

}
