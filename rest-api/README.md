# Go REST API Example

A lightweight, thread-safe RESTful API built in Go using zero external dependencies, leveraging Go's standard library `net/http` package and Go 1.22+ routing features.

---

## 📋 Table of Contents
1. [Prerequisites](#-prerequisites)
2. [How to Run & Build](#-how-to-run--build)
3. [API Endpoints & Testing](#-api-endpoints--testing)
4. [Detailed Code Explanations & Architecture](#-detailed-code-explanations--architecture)
   - [1. Data Model & Concurrency (`models/book.go`)](#1-data-model--concurrency-modelsbookgo)
   - [2. HTTP Request Handlers (`handlers/book_handler.go`)](#2-http-request-handlers-handlersbook_handlergo)
   - [3. Routing & Server Setup (`main.go`)](#3-routing--server-setup-maingo)

---

## ⚙️ Prerequisites

- **Go 1.22+** installed on your system.
  ```bash
  go version
  ```

---

## 🚀 How to Run & Build

### 1. Run Directly from Source
```bash
cd rest-api
go run main.go
```
The server will start listening at `http://localhost:8080`.

### 2. Build a Standalone Executable Binary
To compile a production binary executable:
```bash
cd rest-api
go build -o rest-api main.go
```
Run the compiled binary:
- **Linux / macOS**: `./rest-api`
- **Windows**: `rest-api.exe`

---

## 📡 API Endpoints & Testing

### Available Endpoints

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `http://localhost:8080/` | Welcome message & server status |
| `GET` | `http://localhost:8080/api/books` | Fetch all books |
| `GET` | `http://localhost:8080/api/books/{id}` | Fetch a single book by ID |
| `POST` | `http://localhost:8080/api/books` | Create a new book |
| `PUT` | `http://localhost:8080/api/books/{id}` | Update an existing book |
| `DELETE` | `http://localhost:8080/api/books/{id}` | Delete a book by ID |

### cURL Test Commands

#### 1. Fetch All Books
```bash
curl -i http://localhost:8080/api/books
```

#### 2. Fetch Book by ID
```bash
curl -i http://localhost:8080/api/books/1
```

#### 3. Create a New Book
```bash
curl -i -X POST http://localhost:8080/api/books \
  -H "Content-Type: application/json" \
  -d '{"id":"3", "title":"Designing Data-Intensive Applications", "author":"Martin Kleppmann", "price":45.00}'
```

#### 4. Update an Existing Book
```bash
curl -i -X PUT http://localhost:8080/api/books/3 \
  -H "Content-Type: application/json" \
  -d '{"title":"Designing Data-Intensive Applications (2nd Ed)", "author":"Martin Kleppmann", "price":49.99}'
```

#### 5. Delete a Book
```bash
curl -i -X DELETE http://localhost:8080/api/books/3
```

---

## 📂 Detailed Code Explanations & Architecture

The application separates data access logic, HTTP controllers, and server configuration into distinct modules.

```
rest-api/
├── go.mod                # Module definition file
├── main.go               # Entry point & HTTP router initialization
├── models/
│   └── book.go           # Struct model & thread-safe in-memory store
└── handlers/
    └── book_handler.go   # HTTP handlers for CRUD API operations
```

---

### 1. Data Model & Concurrency (`models/book.go`)

This file defines the data structures and in-memory thread-safe store for books.

#### A. The `Book` Struct
```go
type Book struct {
    ID     string  `json:"id"`
    Title  string  `json:"title"`
    Author string  `json:"author"`
    Price  float64 `json:"price"`
}
```
- **Struct fields**: Defines the attributes of a book (`ID`, `Title`, `Author`, `Price`).
- **Struct Tags (`json:"..."`)**: Directs the standard `encoding/json` package on how to map Go struct fields to JSON property names during encoding and decoding.

#### B. Thread-Safe Store with `sync.RWMutex`
```go
type BookStore struct {
    mu    sync.RWMutex
    books map[string]Book
}
```
- **Map Storage (`map[string]Book`)**: Stores books in memory using the book's unique `ID` as the map key.
- **Thread Safety (`sync.RWMutex`)**: Go web servers process each HTTP request in a separate goroutine. Plain Go maps are not safe for concurrent reads and writes.
  - `s.mu.RLock()` / `s.mu.RUnlock()`: Used for read-only operations (`GetAll()`, `GetByID()`). Multiple HTTP GET requests can read data concurrently without blocking each other.
  - `s.mu.Lock()` / `s.mu.Unlock()`: Used for write operations (`Create()`, `Update()`, `Delete()`). Grants exclusive lock access to guarantee thread-safe modifications.

#### C. Idiomatic Go Practices
- **`defer` for Lock Unlocking**: `defer s.mu.RUnlock()` guarantees that locks are released when the function exits, preventing deadlocks.
- **Slice Allocation**: `make([]Book, 0, len(s.books))` pre-allocates slice capacity for optimal memory allocation when retrieving all books.

---

### 2. HTTP Request Handlers (`handlers/book_handler.go`)

Contains HTTP controller handlers that parse incoming HTTP requests, delegate operations to `BookStore`, and format JSON responses.

#### A. Dependency Injection
```go
type BookHandler struct {
    Store *models.BookStore
}

func NewBookHandler(store *models.BookStore) *BookHandler {
    return &BookHandler{Store: store}
}
```
- Receives `*models.BookStore` pointer via constructor, providing handlers access to the data store.

#### B. Response Formatting Helpers
```go
func writeJSON(w http.ResponseWriter, status int, data any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}
```
- Sets `Content-Type: application/json` headers, assigns HTTP status codes (200 OK, 201 Created, 404 Not Found, 400 Bad Request), and streams JSON output.

#### C. Path Parameters & Request Body Parsing
- **Route Parameter Extraction**: `id := r.PathValue("id")` utilizes Go 1.22+ standard library path matching to extract `{id}` from request URLs.
- **JSON Body Decoding**: `json.NewDecoder(r.Body).Decode(&book)` parses client payload directly into Go structs with error validation.

---

### 3. Routing & Server Setup (`main.go`)

Initializes dependencies and configures the HTTP server router using standard library primitives.

#### A. Go 1.22+ Standard HTTP Multiplexer
```go
mux := http.NewServeMux()

mux.HandleFunc("GET /api/books", bookHandler.GetBooks)
mux.HandleFunc("GET /api/books/{id}", bookHandler.GetBookByID)
mux.HandleFunc("POST /api/books", bookHandler.CreateBook)
mux.HandleFunc("PUT /api/books/{id}", bookHandler.UpdateBook)
mux.HandleFunc("DELETE /api/books/{id}", bookHandler.DeleteBook)
```
- **Method-based Routing**: Patterns like `"GET /api/books"` and `"POST /api/books"` ensure endpoints execute only for matching HTTP verbs.
- **Path Wildcards**: `{id}` automatically extracts parameters without needing external routing libraries.

#### B. Server Binding
```go
if err := http.ListenAndServe(":8080", mux); err != nil {
    log.Fatalf("Server failed to start: %v", err)
}
```
- Listens on TCP port `8080` and serves incoming requests.
