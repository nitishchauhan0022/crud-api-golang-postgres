package main

import (
	"fmt"
	"crud-api-golang-postgres/handler"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	route := mux.NewRouter()
	route.HandleFunc("/api/book/{id}", handler.GetBook).Methods("GET", "OPTIONS")
	route.HandleFunc("/api/book", handler.GetAllBook).Methods("GET", "OPTIONS")
	route.HandleFunc("/api/newbook", handler.CreateBook).Methods("POST", "OPTIONS")
	route.HandleFunc("/api/book/{id}", handler.UpdateBook).Methods("PUT", "OPTIONS")
	route.HandleFunc("/api/deletebook/{id}", handler.DeleteBook).Methods("DELETE", "OPTIONS")

	fmt.Println("Starting server on the port 8080....")
	log.Fatal(http.ListenAndServe(":8080", route))
}
