package models

import (
	"errors"
	"fmt"
	"sync"
)

// Author represents a book author entity
type Author struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Bio  string `json:"bio"`
}

// Book represents a book item in the library store
type Book struct {
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	AuthorID string  `json:"authorId"`
	Price    float64 `json:"price"`
}

// Store provides thread-safe in-memory CRUD operations for Authors and Books
type Store struct {
	mu            sync.RWMutex
	authors       map[string]Author
	books         map[string]Book
	nextAuthorID  int
	nextBookID    int
}

// NewStore initializes a Store pre-populated with sample domain data
func NewStore() *Store {
	s := &Store{
		authors:      make(map[string]Author),
		books:        make(map[string]Book),
		nextAuthorID: 1,
		nextBookID:   1,
	}

	// Seed sample authors
	a1 := s.createAuthorInternal("Alan A. A. Donovan", "Co-author of The Go Programming Language book and Go core team member.")
	a2 := s.createAuthorInternal("Katherine Cox-Buday", "Computer scientist and author specializing in Go concurrency patterns.")
	a3 := s.createAuthorInternal("Martin Kleppmann", "Researcher and author of Designing Data-Intensive Applications.")

	// Seed sample books linked to authors
	s.createBookInternal("The Go Programming Language", a1.ID, 39.99)
	s.createBookInternal("Concurrency in Go", a2.ID, 29.99)
	s.createBookInternal("Designing Data-Intensive Applications", a3.ID, 45.00)

	return s
}

func (s *Store) createAuthorInternal(name, bio string) Author {
	id := fmt.Sprintf("%d", s.nextAuthorID)
	s.nextAuthorID++
	author := Author{ID: id, Name: name, Bio: bio}
	s.authors[id] = author
	return author
}

func (s *Store) createBookInternal(title, authorID string, price float64) Book {
	id := fmt.Sprintf("%d", s.nextBookID)
	s.nextBookID++
	book := Book{ID: id, Title: title, AuthorID: authorID, Price: price}
	s.books[id] = book
	return book
}

// --- Author Methods ---

func (s *Store) GetAuthorByID(id string) (Author, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	author, exists := s.authors[id]
	return author, exists
}

func (s *Store) GetAllAuthors() []Author {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Author, 0, len(s.authors))
	for _, author := range s.authors {
		result = append(result, author)
	}
	return result
}

func (s *Store) CreateAuthor(name, bio string) (Author, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if name == "" {
		return Author{}, errors.New("author name cannot be empty")
	}

	author := s.createAuthorInternal(name, bio)
	return author, nil
}

func (s *Store) UpdateAuthor(id string, name *string, bio *string) (Author, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	author, exists := s.authors[id]
	if !exists {
		return Author{}, errors.New("author not found")
	}

	if name != nil {
		author.Name = *name
	}
	if bio != nil {
		author.Bio = *bio
	}

	s.authors[id] = author
	return author, nil
}

func (s *Store) DeleteAuthor(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.authors[id]; !exists {
		return errors.New("author not found")
	}

	delete(s.authors, id)
	// Optionally disassociate books or leave them
	return nil
}

// --- Book Methods ---

func (s *Store) GetBookByID(id string) (Book, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	book, exists := s.books[id]
	return book, exists
}

func (s *Store) GetAllBooks(authorID string, maxPrice float64) []Book {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]Book, 0, len(s.books))
	for _, book := range s.books {
		if authorID != "" && book.AuthorID != authorID {
			continue
		}
		if maxPrice > 0 && book.Price > maxPrice {
			continue
		}
		result = append(result, book)
	}
	return result
}

func (s *Store) GetBooksByAuthorID(authorID string) []Book {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]Book, 0)
	for _, book := range s.books {
		if book.AuthorID == authorID {
			result = append(result, book)
		}
	}
	return result
}

func (s *Store) CreateBook(title, authorID string, price float64) (Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if title == "" {
		return Book{}, errors.New("book title cannot be empty")
	}
	if _, exists := s.authors[authorID]; !exists {
		return Book{}, fmt.Errorf("author with ID '%s' does not exist", authorID)
	}

	book := s.createBookInternal(title, authorID, price)
	return book, nil
}

func (s *Store) UpdateBook(id string, title *string, authorID *string, price *float64) (Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	book, exists := s.books[id]
	if !exists {
		return Book{}, errors.New("book not found")
	}

	if title != nil {
		book.Title = *title
	}
	if authorID != nil {
		if _, exists := s.authors[*authorID]; !exists {
			return Book{}, fmt.Errorf("author with ID '%s' does not exist", *authorID)
		}
		book.AuthorID = *authorID
	}
	if price != nil {
		book.Price = *price
	}

	s.books[id] = book
	return book, nil
}

func (s *Store) DeleteBook(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.books[id]; !exists {
		return errors.New("book not found")
	}

	delete(s.books, id)
	return nil
}
