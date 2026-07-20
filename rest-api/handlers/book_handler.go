package handlers

import (
	"encoding/json"
	"net/http"
	"rest-api/models"
)

// BookHandler holds reference to the data store and handles HTTP requests
type BookHandler struct {
	Store *models.BookStore
}

// NewBookHandler creates a new BookHandler instance
func NewBookHandler(store *models.BookStore) *BookHandler {
	return &BookHandler{Store: store}
}

// writeJSON is a helper function to send JSON responses with appropriate headers and HTTP status codes
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError is a helper function to send formatted JSON error responses
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// GetBooks handles GET /api/books - returns all books
func (h *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	books := h.Store.GetAll()
	writeJSON(w, http.StatusOK, books)
}

// GetBookByID handles GET /api/books/{id} - returns a single book by ID
func (h *BookHandler) GetBookByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	book, found := h.Store.GetByID(id)
	if !found {
		writeError(w, http.StatusNotFound, "Book not found")
		return
	}

	writeJSON(w, http.StatusOK, book)
}

// CreateBook handles POST /api/books - creates a new book
func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var newBook models.Book

	if err := json.NewDecoder(r.Body).Decode(&newBook); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if newBook.ID == "" || newBook.Title == "" || newBook.Author == "" {
		writeError(w, http.StatusBadRequest, "Fields 'id', 'title', and 'author' are required")
		return
	}

	createdBook, err := h.Store.Create(newBook)
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, createdBook)
}

// UpdateBook handles PUT /api/books/{id} - updates an existing book
func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var updatedBook models.Book
	if err := json.NewDecoder(r.Body).Decode(&updatedBook); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	book, err := h.Store.Update(id, updatedBook)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, book)
}

// DeleteBook handles DELETE /api/books/{id} - deletes a book by ID
func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.Store.Delete(id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Book deleted successfully"})
}
