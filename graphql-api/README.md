# 🚀 Go GraphQL API

A lightweight, thread-safe, high-performance GraphQL API built in Go using `github.com/graphql-go/graphql` with zero code-generation dependencies. Features full CRUD operations for Books and Authors, nested relationship resolvers, and an embedded interactive **GraphiQL IDE**.

---

## 📁 Features

- **Nested Relational Resolvers**: Fetch books along with author details (`Book.author`), or authors with their published books (`Author.books`) in a single GraphQL query.
- **Thread-Safe In-Memory Data Store**: Concurrency safe via Go's `sync.RWMutex`.
- **Interactive GraphiQL IDE**: Access `http://localhost:8080/` in any browser for an interactive query builder & schema documentation browser.
- **Filtering & Arguments**: Filter books by `authorId` and/or `maxPrice`.
- **Full Mutations**: Create, update, and delete operations for both Books and Authors.
- **CORS Support**: Cross-Origin Resource Sharing enabled for seamless integration with frontend applications.

---

## ⚙️ How to Run & Build

### 1. Run Directly from Source
```bash
cd graphql-api
go run main.go
```
The server will start listening at `http://localhost:8080`.

### 2. Run Tests
```bash
cd graphql-api
go test -v ./...
```

### 3. Build Executable Binary
```bash
cd graphql-api
go build -o graphql-api main.go
./graphql-api
```

---

## 🎨 Interactive GraphiQL UI

Open your browser and navigate to:
👉 **`http://localhost:8080/`**

You can write queries, view autocomplete suggestions, and inspect the interactive GraphQL Schema documentation directly in your browser.

---

## 📡 GraphQL Queries & cURL Test Examples

The GraphQL HTTP endpoint accepts `POST` requests at `http://localhost:8080/graphql` with a JSON payload `{"query": "..."}`.

### 1. Fetch All Books with Author Details
```bash
curl -i -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query": "{ books { id title price author { id name bio } } }"}'
```

### 2. Fetch a Single Book by ID
```bash
curl -i -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query": "{ book(id: \"1\") { id title price author { name } } }"}'
```

### 3. Fetch Authors with List of Written Books (Nested Query)
```bash
curl -i -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query": "{ authors { id name bio books { id title price } } }"}'
```

### 4. Create a New Author (Mutation)
```bash
curl -i -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query": "mutation { createAuthor(name: \"Rob Pike\", bio: \"Co-creator of Go.\") { id name bio } }"}'
```

### 5. Create a New Book (Mutation)
```bash
curl -i -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query": "mutation { createBook(title: \"Go in Action\", authorId: \"1\", price: 34.99) { id title price author { name } } }"}'
```

### 6. Update a Book (Mutation)
```bash
curl -i -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query": "mutation { updateBook(id: \"1\", price: 42.50) { id title price } }"}'
```

### 7. Delete a Book (Mutation)
```bash
curl -i -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query": "mutation { deleteBook(id: \"1\") }"}'
```

---

## 📂 Architecture Overview

```
graphql-api/
├── go.mod                # Module definition file
├── go.sum                # Dependencies checksum
├── main.go               # Application entry point & HTTP server
├── models/
│   └── store.go          # Data structures & thread-safe store
├── schema/
│   ├── schema.go         # GraphQL schema, types, queries & mutations
│   └── schema_test.go    # Unit tests for GraphQL endpoints
├── handlers/
│   └── graphql.go        # HTTP handler & GraphiQL IDE web interface
└── README.md             # Project documentation
```
