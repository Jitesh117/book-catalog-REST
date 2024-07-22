package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"book-catalog/models"
	"book-catalog/storage"
)

type BookHandler struct {
	storage *storage.PostgresStorage
}

func NewBookHandler(storage *storage.PostgresStorage) *BookHandler {
	return &BookHandler{storage: storage}
}

func (h *BookHandler) HandleBooks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getBooks(w, r)
	case http.MethodPost:
		h.createBook(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *BookHandler) HandleBook(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/books/")
	switch r.Method {
	case http.MethodGet:
		h.getBook(w, r, id)
	case http.MethodPut:
		h.updateBook(w, r, id)
	case http.MethodDelete:
		h.deleteBook(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *BookHandler) getBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.storage.GetAll()
	if err != nil {
		http.Error(w, "Failed to retrieve books: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func (h *BookHandler) getBook(w http.ResponseWriter, r *http.Request, id string) {
	book, err := h.storage.Get(id)
	if err != nil {
		if err.Error() == "book not found" {
			http.Error(w, "Book not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve book: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) createBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
		return
	}
	err = h.storage.Create(book)
	if err != nil {
		http.Error(w, "Failed to create book: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) updateBook(w http.ResponseWriter, r *http.Request, id string) {
	var book models.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
		return
	}
	book.ID = id
	err = h.storage.Update(book)
	if err != nil {
		if err.Error() == "book not found" {
			http.Error(w, "Book not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to update book: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) deleteBook(w http.ResponseWriter, r *http.Request, id string) {
	err := h.storage.Delete(id)
	if err != nil {
		if err.Error() == "book not found" {
			http.Error(w, "Book not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to delete book: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
