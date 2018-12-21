package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/subosito/gotenv"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   string `json:"year"`
}

var books []Book
var db *sql.DB

func init() {
	gotenv.Load()
}

func logFatal(err error) {
	if err != nil {
		// Stop execution if error
		log.Fatal(err)
	}
}

func main() {

	pgUrl, err := pq.ParseURL(os.Getenv("ELEPHANTSQL_URL"))
	logFatal(err)

	db, err = sql.Open("postgres", pgUrl)
	logFatal(err)

	err = db.Ping()
	logFatal(err)

	router := mux.NewRouter()

	// Append test data
	//books = append(books,
	//	Book{ID: 1, Title:"Book 1", Author:"Author 1", Year:"2018/01/10"},
	//	Book{ID: 2, Title:"Book 2", Author:"Author 2", Year:"2018/02/04"})

	// Get all books
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books", addBook).Methods("POST")
	router.HandleFunc("/books", updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", removeBook).Methods("DELETE")

	// Start server
	log.Fatal(http.ListenAndServe(":8080", router))

}

// Get single record by ID
func getBook(w http.ResponseWriter, r *http.Request) {
	log.Println("Get single book")
	var book Book

	// Read value from url
	params := mux.Vars(r)
	// pars str to int
	id, _ := strconv.Atoi(params["id"])
	//log.Println(reflect.TypeOf(id))

	rows := db.QueryRow("SELECT * FROM books WHERE id=$1", id)

	// map result to a Book
	err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
	// Stop if error
	logFatal(err)

	json.NewEncoder(w).Encode(book)

}

// Get all books from table
func getBooks(w http.ResponseWriter, r *http.Request) {
	log.Println("Get all book")
	var book Book

	books = []Book{}

	rows, err := db.Query("SELECT * FROM books")
	logFatal(err)

	defer rows.Close()

	for rows.Next() {
		// map each row to a instance of Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		logFatal(err)
		books = append(books, book)
	}

	json.NewEncoder(w).Encode(books)
}

// Add new record to DB
func addBook(w http.ResponseWriter, r *http.Request) {
	log.Println("Add a book")

	var book Book
	var bookID int

	// Decode POST body
	_ = json.NewDecoder(r.Body).Decode(&book)
	err := db.QueryRow("INSERT INTO books (title, author, year) VALUES($1, $2, $3) RETURNING id",
		book.Title, book.Author, book.Year).Scan(&bookID)
	logFatal(err)

	json.NewEncoder(w).Encode(bookID)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	log.Println("update a book")

	var book Book

	json.NewDecoder(r.Body).Decode(&book)

	result, err := db.Exec("UPDATE books SET title=$1, author=$2, year=$3 WHERE id=$4 RETURNING id",
		book.Title, book.Author, book.Year,  book.ID)

	// Check number of rows affected
	rowsUpdated, err := result.RowsAffected()
	logFatal(err)

	json.NewEncoder(w).Encode(rowsUpdated)

}

func removeBook(w http.ResponseWriter, r *http.Request) {
	log.Println("remove book")

	// Read param from url
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	result, err :=db.Exec("DELETE FROM books WHERE id = $1", id)
	logFatal(err)

	rowsDeleted, err := result.RowsAffected()
	logFatal(err)

	json.NewEncoder(w).Encode(rowsDeleted)
}
