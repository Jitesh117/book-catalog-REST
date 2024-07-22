package main

import (
	"fmt"
	"log"
	"net/http"

	"book-catalog/handlers"
	"book-catalog/storage"
)

func main() {
	store, err := storage.NewPostgresStorage(
		"localhost", "5432", "postgres", "mysecretpassword", "book_catalog")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer store.Close()

	if err := store.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	bookHandler := handlers.NewBookHandler(store)

	http.HandleFunc("/books", bookHandler.HandleBooks)
	http.HandleFunc("/books/", bookHandler.HandleBook)

	fmt.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
