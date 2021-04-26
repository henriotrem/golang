package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	ID    primitive.ObjectID `json:"_id,omitempty"`
	Isbn  string             `json:"isbn,omitempty"`
	Title string             `json:"title,omitempty"`
}

type ResponseMessage struct {
	Message string `json:"message"`
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", r))
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var books []Book

	for i := 1; i <= 10; i++ {
		id, _ := primitive.ObjectIDFromHex(strconv.Itoa(i))

		books = append(books, Book{
			ID:    id,
			Isbn:  "isbn " + strconv.Itoa(i),
			Title: "livre " + strconv.Itoa(i),
		})
	}

	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)

	id, _ := primitive.ObjectIDFromHex(params["id"])

	book := Book{
		ID:    id,
		Isbn:  "isbn 1",
		Title: "livre 1",
	}

	json.NewEncoder(w).Encode(book)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := ResponseMessage{
		Message: "Inserted ID 32lhl542",
	}

	json.NewEncoder(w).Encode(response)
}
