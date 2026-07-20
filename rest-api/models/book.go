package models

import (
	"errors"
	"sync"
)

// Book represents a book item in our library inventory
type Book struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
}

// BookStore provides thread-safe in-memory CRUD operations for Books
type BookStore struct {
	mu    sync.RWMutex
	books map[string]Book
}

// NewBookStore initializes a new BookStore with sample data
func NewBookStore() *BookStore {
	store := &BookStore{
		books: make(map[string]Book),
	}
	// Initial seed data
	store.books["1"] = Book{ID: "1", Title: "The Go Programming Language", Author: "Alan A. A. Donovan", Price: 39.99}
	store.books["2"] = Book{ID: "2", Title: "Concurrency in Go", Author: "Katherine Cox-Buday", Price: 29.99}
	return store
}

// GetAll returns all books from the store
func (s *BookStore) GetAll() []Book {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]Book, 0, len(s.books))
	for _, book := range s.books {
		result = append(result, book)
	}
	return result
}

// GetByID fetches a single book by its unique ID
func (s *BookStore) GetByID(id string) (Book, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	book, exists := s.books[id]
	return book, exists
}

// Create adds a new book to the store
func (s *BookStore) Create(book Book) (Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.books[book.ID]; exists {
		return Book{}, errors.New("book with this ID already exists")
	}

	s.books[book.ID] = book
	return book, nil
}

// Update modifies an existing book in the store
func (s *BookStore) Update(id string, updatedBook Book) (Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.books[id]; !exists {
		return Book{}, errors.New("book not found")
	}

	updatedBook.ID = id
	s.books[id] = updatedBook
	return updatedBook, nil
}

// Delete removes a book by ID from the store
func (s *BookStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.books[id]; !exists {
		return errors.New("book not found")
	}

	delete(s.books, id)
	return nil
}
