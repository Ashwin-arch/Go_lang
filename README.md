# Go_lang

A repository for Go language projects and applications.

---

## 📁 Projects Overview

1. **[🧮 Calculator (GUI)](#-1-calculator-gui)**: A desktop GUI calculator built using Go and [Fyne v2](https://fyne.io/).
2. **[🌐 REST API](#-2-rest-api)**: A lightweight, thread-safe RESTful API built in Go using standard library `net/http` and Go 1.22+ routing features.

---

## 🧮 1. Calculator (GUI)

### Features
- Standard arithmetic operations (`+`, `-`, `×`, `÷`, `%`)
- Sign toggle (`±`), clear (`C`), and backspace (`⌫`)
- Modern dark-themed user interface

### How to Run Calculator
```bash
cd calculator
go run .
```

---

## 🌐 2. REST API

A thread-safe RESTful API built in Go using **zero external dependencies**, leveraging Go's standard library `net/http` package and Go 1.22+ routing features.

### ⚙️ Prerequisites

- **Go 1.22+** installed on your system.
  ```bash
  go version
  ```

### 🚀 How to Run & Build

#### 1. Run Directly from Source
```bash
cd rest-api
go run main.go
```
The server will start listening at `http://localhost:8080`.

#### 2. Build a Standalone Executable Binary
```bash
cd rest-api
go build -o rest-api main.go
```
Run the compiled binary:
- **Linux / macOS**: `./rest-api`
- **Windows**: `rest-api.exe`

---

### 📡 API Endpoints & cURL Testing

#### Available Endpoints

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `http://localhost:8080/` | Welcome message & server status |
| `GET` | `http://localhost:8080/api/books` | Fetch all books |
| `GET` | `http://localhost:8080/api/books/{id}` | Fetch a single book by ID |
| `POST` | `http://localhost:8080/api/books` | Create a new book |
| `PUT` | `http://localhost:8080/api/books/{id}` | Update an existing book |
| `DELETE` | `http://localhost:8080/api/books/{id}` | Delete a book by ID |

#### cURL Test Commands

```bash
# 1. Fetch All Books
curl -i http://localhost:8080/api/books

# 2. Fetch Book by ID
curl -i http://localhost:8080/api/books/1

# 3. Create a New Book
curl -i -X POST http://localhost:8080/api/books \
  -H "Content-Type: application/json" \
  -d '{"id":"3", "title":"Designing Data-Intensive Applications", "author":"Martin Kleppmann", "price":45.00}'

# 4. Update an Existing Book
curl -i -X PUT http://localhost:8080/api/books/3 \
  -H "Content-Type: application/json" \
  -d '{"title":"Designing Data-Intensive Applications (2nd Ed)", "author":"Martin Kleppmann", "price":49.99}'

# 5. Delete a Book
curl -i -X DELETE http://localhost:8080/api/books/3
```

---

### 📂 Detailed Code Explanations & Architecture

The REST API application separates data access logic, HTTP controllers, and server configuration into distinct modules:

```
rest-api/
├── go.mod                # Module definition file
├── main.go               # Entry point & HTTP router initialization
├── models/
│   └── book.go           # Struct model & thread-safe in-memory store
└── handlers/
    └── book_handler.go   # HTTP handlers for CRUD API operations
```

#### 1. Data Model & Concurrency (`rest-api/models/book.go`)

This file defines the data structures and in-memory thread-safe store for books.

- **`Book` Struct**:
  ```go
  type Book struct {
      ID     string  `json:"id"`
      Title  string  `json:"title"`
      Author string  `json:"author"`
      Price  float64 `json:"price"`
  }
  ```
  - **Struct fields**: Defines the attributes of a book (`ID`, `Title`, `Author`, `Price`).
  - **Struct Tags (`json:"..."`)**: Directs `encoding/json` on how to serialize/deserialize Go struct fields to/from JSON keys.

- **Thread-Safe Store (`BookStore`)**:
  ```go
  type BookStore struct {
      mu    sync.RWMutex
      books map[string]Book
  }
  ```
  - **Map Storage (`map[string]Book`)**: Stores books in memory using the book's unique `ID` as the map key.
  - **Thread Safety (`sync.RWMutex`)**: Go web servers process each HTTP request in a separate goroutine. Plain Go maps are not thread-safe.
    - `RLock()` / `RUnlock()`: Used for read-only operations (`GetAll()`, `GetByID()`). Multiple HTTP GET requests can read data concurrently without blocking each other.
    - `Lock()` / `Unlock()`: Used for write operations (`Create()`, `Update()`, `Delete()`). Grants exclusive lock access to guarantee thread-safe modifications.
  - **`defer` for Lock Management**: Ensures lock release upon function return (`defer s.mu.RUnlock()`), preventing deadlocks.
  - **Memory Pre-allocation**: `make([]Book, 0, len(s.books))` pre-allocates slice capacity for optimal memory usage.

#### 2. HTTP Request Handlers (`rest-api/handlers/book_handler.go`)

Contains HTTP controller handlers that parse incoming HTTP requests, delegate operations to `BookStore`, and format JSON responses.

- **Dependency Injection**: `BookHandler` receives `*models.BookStore` pointer via `NewBookHandler` constructor.
- **Response Formatting**: `writeJSON` sets `Content-Type: application/json` headers, assigns HTTP status codes (200 OK, 201 Created, 404 Not Found, 400 Bad Request), and streams JSON output using `json.NewEncoder(w).Encode(data)`.
- **Route Parameters**: `id := r.PathValue("id")` utilizes Go 1.22+ standard library path matching to extract `{id}` from request URLs.
- **Request Body Parsing**: `json.NewDecoder(r.Body).Decode(&book)` parses incoming client JSON payload directly into Go structs.

#### 3. Routing & Server Setup (`rest-api/main.go`)

Initializes dependencies and configures the HTTP server router using standard library primitives.

- **Go 1.22+ Standard Routing**:
  ```go
  mux := http.NewServeMux()

  mux.HandleFunc("GET /api/books", bookHandler.GetBooks)
  mux.HandleFunc("GET /api/books/{id}", bookHandler.GetBookByID)
  mux.HandleFunc("POST /api/books", bookHandler.CreateBook)
  mux.HandleFunc("PUT /api/books/{id}", bookHandler.UpdateBook)
  mux.HandleFunc("DELETE /api/books/{id}", bookHandler.DeleteBook)
  ```
  - **Method Matching**: `"GET /api/books"` and `"POST /api/books"` restrict execution strictly to matching HTTP verbs.
  - **Path Wildcards**: `{id}` automatically extracts path parameters without external dependencies.
- **Server Binding**: Binds TCP port `:8080` using `http.ListenAndServe(":8080", mux)`.
