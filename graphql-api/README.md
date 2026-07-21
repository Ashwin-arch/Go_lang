# 📊 Go GraphQL API

A lightweight, thread-safe, high-performance GraphQL API built in Go using `github.com/graphql-go/graphql` with zero code-generation dependencies. Features full CRUD operations for Books and Authors, nested relationship resolvers, unit tests, and an embedded interactive **GraphiQL IDE**.

---

## 📋 Table of Contents
1. [Features](#-features)
2. [How to Run & Build](#-how-to-run--build)
3. [Interactive GraphiQL UI](#-interactive-graphiql-ui)
4. [GraphQL Queries & cURL Testing](#-graphql-queries--curl-testing)
5. [Detailed Code Explanations & Architecture](#-detailed-code-explanations--architecture)
   - [1. Data Model & Concurrency (`models/store.go`)](#1-data-model--concurrency-modelsstorego)
   - [2. GraphQL Schema & Resolvers (`schema/schema.go`)](#2-graphql-schema--resolvers-schemaschemago)
   - [3. Automated Unit Tests (`schema/schema_test.go`)](#3-automated-unit-tests-schemaschema_testgo)
   - [4. HTTP Handler & GraphiQL IDE (`handlers/graphql.go`)](#4-http-handler--graphiql-ide-handlersgraphqlgo)
   - [5. Server Setup & Routing (`main.go`)](#5-server-setup--routing-maingo)

---

## ✨ Features

- **Nested Relational Resolvers**: Fetch books along with author details (`Book.author`), or authors with their published books (`Author.books`) in a single GraphQL query.
- **Thread-Safe In-Memory Data Store**: Concurrency safe via Go's `sync.RWMutex`.
- **Interactive GraphiQL IDE**: Access `http://localhost:8080/` in any browser for an interactive query builder & schema documentation browser.
- **Filtering & Arguments**: Filter books by `authorId` and/or `maxPrice`.
- **Full Mutations**: Create, update, and delete operations for both Books and Authors.
- **CORS Support**: Cross-Origin Resource Sharing enabled for seamless integration with frontend applications.

---

## 🚀 How to Run & Build

### 1. Run Directly from Source
```bash
cd graphql-api
go run main.go
```
The server will start listening at `http://localhost:8080`.

### 2. Run Automated Unit Tests
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

## 📡 GraphQL Queries & cURL Testing

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

## 📂 Detailed Code Explanations & Architecture

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

---

### 1. Data Model & Concurrency (`models/store.go`)

This file defines the domain models and thread-safe data store logic.

#### A. Domain Structs
```go
type Author struct {
    ID   string `json:"id"`
    Name string `json:"name"`
    Bio  string `json:"bio"`
}

type Book struct {
    ID       string  `json:"id"`
    Title    string  `json:"title"`
    AuthorID string  `json:"authorId"`
    Price    float64 `json:"price"`
}
```

#### B. Thread-Safe Store (`Store`)
```go
type Store struct {
    mu           sync.RWMutex
    authors      map[string]Author
    books        map[string]Book
    nextAuthorID int
    nextBookID   int
}
```
- **Concurrency Management**: Uses `sync.RWMutex` so multiple goroutines can safely read data simultaneously using `RLock()`, while write operations (`Create`, `Update`, `Delete`) acquire exclusive `Lock()`.
- **Pre-populated Seed Data**: Initialized via `NewStore()` with sample books and authors (e.g. Alan A. A. Donovan, Katherine Cox-Buday, Martin Kleppmann).

---

### 2. GraphQL Schema & Resolvers (`schema/schema.go`)

Defines the GraphQL schema using `github.com/graphql-go/graphql` programmatic builder without code generation.

#### A. Type Thunks & Relational Resolvers
```go
authorType = graphql.NewObject(graphql.ObjectConfig{
    Name: "Author",
    Fields: graphql.FieldsThunk(func() graphql.Fields {
        return graphql.Fields{
            "books": &graphql.Field{
                Type: graphql.NewList(bookType),
                Resolve: func(p graphql.ResolveParams) (interface{}, error) {
                    author := p.Source.(models.Author)
                    return store.GetBooksByAuthorID(author.ID), nil
                },
            },
        }
    }),
})
```
- **`graphql.FieldsThunk`**: Solves circular type dependencies (e.g., `Author` referencing `Book` and `Book` referencing `Author`) by lazily evaluating field definitions.
- **Nested Resolution**: When `book.author` is requested, the resolver extracts `p.Source.(models.Book)` and delegates to `store.GetAuthorByID(book.AuthorID)`.

#### B. Root Query & Arguments
- `books`: Supports optional `authorId` and `maxPrice` argument filtering.
- `book`: Requires non-null `id` argument (`graphql.NewNonNull(graphql.ID)`).

#### C. Root Mutation
- Declares input arguments for `createBook`, `updateBook`, `deleteBook`, `createAuthor`, `updateAuthor`, `deleteAuthor` and maps incoming argument maps (`p.Args`) to store calls.

---

### 3. Automated Unit Tests (`schema/schema_test.go`)

Executes unit testing directly against the GraphQL schema using `graphql.Do(...)` without starting an HTTP server:
```go
func executeQuery(query string, gqlSchema graphql.Schema) *graphql.Result {
    return graphql.Do(graphql.Params{
        Schema:        gqlSchema,
        RequestString: query,
    })
}
```
- Tests nested relational queries, single item lookups, and mutation state changes.

---

### 4. HTTP Handler & GraphiQL IDE (`handlers/graphql.go`)

- **CORS Setup**: Attaches CORS response headers (`Access-Control-Allow-Origin: *`).
- **Browser Detection**: Detects incoming `GET` browser requests or `Accept: text/html` and serves an embedded single-page GraphiQL HTML user interface.
- **Payload Parsing**: Unmarshals incoming HTTP POST request JSON bodies (`query`, `operationName`, `variables`) into `graphql.Params` and encodes results back to JSON.

---

### 5. Server Setup & Routing (`main.go`)

```go
mux := http.NewServeMux()
mux.Handle("/graphql", gqlHandler)
mux.Handle("/", gqlHandler)
http.ListenAndServe(":" + port, mux)
```
- Binds port `:8080` and routes root requests `/` to GraphiQL IDE and `/graphql` to API query execution.
