package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"crud-api-golang-postgres/schema"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func initConnection() *sql.DB {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	POSTGRES_USER := os.Getenv("POSTGRES_USER")
	POSTGRES_PASSWORD := os.Getenv("POSTGRES_PASSWORD")
	POSTGRES_PORT := os.Getenv("PORT")
	POSTGRES_DATABASE := os.Getenv("DATABASE")

	postgres_path := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s", POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_PORT, POSTGRES_DATABASE)
	fmt.Print(postgres_path)
	db, err := sql.Open("postgres", postgres_path)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	var book schema.Book
	err := json.NewDecoder(r.Body).Decode(&book)

	if err != nil {
		log.Printf("Unable to decode the request body.  %v", err)
	}
	insertID := insertBook(book)

	res := schema.Response{
		ID:      insertID,
		Message: "Book created successfully",
	}
	json.NewEncoder(w).Encode(res)
}

func GetBook(w http.ResponseWriter, r *http.Request) {
	variables := mux.Vars(r)
	id, err := strconv.Atoi(variables["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}
	book, err := getBook(int64(id))

	if err != nil {
		log.Fatalf("Unable to get book. %v", err)
	}
	json.NewEncoder(w).Encode(book)
}

func GetAllBook(w http.ResponseWriter, r *http.Request) {
	books, err := getAllBooks()

	if err != nil {
		log.Fatalf("Unable to get all books. %v", err)
	}
	json.NewEncoder(w).Encode(books)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	variables := mux.Vars(r)
	id, err := strconv.Atoi(variables["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}
	var book schema.Book
	err = json.NewDecoder(r.Body).Decode(&book)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}
	updatedRows := updateBook(int64(id), book)
	var msg string
	if updatedRows == 1 {
		msg = "Book Updated successfully"
	} else {
		msg = fmt.Sprintf("Book updated. Total rows/record affected %v", updatedRows)

	}
	res := schema.Response{
		ID:      int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}
	deletedRows := deleteBook(int64(id))
	var msg string
	if deletedRows == 1 {
		msg = "Book Deleted successfully"
	} else {
		msg = fmt.Sprintf("Book Deleted. Total rows/record affected %v", deletedRows)

	}
	res := schema.Response{
		ID:      int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)
}

func insertBook(book schema.Book) int64 {
	db := initConnection()
	defer db.Close()
	sqlStatement := `INSERT INTO books (name, price, company) VALUES ($1, $2, $3) RETURNING bookid`

	var id int64
	err := db.QueryRow(sqlStatement, book.Name, book.Author, book.Publisher).Scan(&id)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	return id
}

func getBook(id int64) (schema.Book, error) {
	db := initConnection()
	defer db.Close()

	var book schema.Book
	sqlStatement := `SELECT * FROM books WHERE bookid=$1`
	row := db.QueryRow(sqlStatement, id)

	err := row.Scan(&book.BookID, &book.Name, &book.Author, &book.Publisher)
	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return book, nil
	case nil:
		return book, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}
	return book, err
}

func getAllBooks() ([]schema.Book, error) {
	db := initConnection()
	defer db.Close()

	var books []schema.Book
	sqlStatement := `SELECT * FROM books`
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var book schema.Book
		err = rows.Scan(&book.BookID, &book.Name, &book.Author, &book.Publisher)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}
		books = append(books, book)
	}
	return books, err
}

func updateBook(id int64, book schema.Book) int64 {

	db := initConnection()
	defer db.Close()

	sqlStatement := `UPDATE books SET name=$2, price=$3, company=$4 WHERE bookid=$1`

	res, err := db.Exec(sqlStatement, id, book.Name, book.Author, book.Publisher)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}
	return rowsAffected
}

func deleteBook(id int64) int64 {

	db := initConnection()
	defer db.Close()

	sqlStatement := `DELETE FROM books WHERE bookid=$1`
	res, err := db.Exec(sqlStatement, id)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}
	return rowsAffected
}
