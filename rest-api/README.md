# Go REST API Example

A simple, thread-safe, RESTful API built in Go using the standard library `net/http` package.

## 🚀 Quick Start

### Run the API Server
```bash
cd rest-api
go run main.go
```
The server will start at `http://localhost:8080`.

---

## 📡 API Endpoints & cURL Examples

### 1. Get All Books
```bash
curl -X GET http://localhost:8080/api/books
```

### 2. Get Book by ID
```bash
curl -X GET http://localhost:8080/api/books/1
```

### 3. Create a New Book
```bash
curl -X POST http://localhost:8080/api/books \
  -H "Content-Type: application/json" \
  -d '{"id":"3", "title":"Designing Data-Intensive Applications", "author":"Martin Kleppmann", "price":44.99}'
```

### 4. Update an Existing Book
```bash
curl -X PUT http://localhost:8080/api/books/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"The Go Programming Language (2nd Edition)", "author":"Alan A. A. Donovan", "price":49.99}'
```

### 5. Delete a Book
```bash
curl -X DELETE http://localhost:8080/api/books/3
```
