# Go_lang

A repository for Go language projects and applications.

---

## 📁 Projects Overview

1. **[🧮 1. Calculator (GUI)](#-1-calculator-gui)**: A desktop GUI calculator built using Go and [Fyne v2](https://fyne.io/).
2. **[🌐 2. REST API](#-2-rest-api)**: A lightweight, thread-safe RESTful API built in Go using standard library `net/http` and Go 1.22+ routing features.
3. **[📊 3. GraphQL API](#-3-graphql-api)**: A thread-safe GraphQL API built in Go with nested relational schema, CRUD mutations, unit tests, and an embedded **GraphiQL IDE**.

---

## 🧮 1. Calculator (GUI)

A cross-platform desktop calculator application built using Go and Fyne v2.

### ✨ Features
- Standard arithmetic operations (`+`, `-`, `×`, `÷`, `%`)
- Sign toggle (`±`), clear (`C`), and backspace (`⌫`)
- Divide-by-zero error handling (`Error: Div by 0`)
- Modern dark-themed user interface with color-coded button importance levels

### ⚙️ Prerequisites & Running

```bash
cd calculator
go run .
```

To build a standalone desktop executable:
```bash
cd calculator
go build -o calculator main.go
```

### 📂 Detailed Code Explanations & Architecture (`calculator/main.go`)

The GUI Calculator logic is contained within `calculator/main.go`:

#### 1. State Management (`calculator` struct)
```go
type calculator struct {
    display *widget.Entry
    op      string
    val1    float64
    newNum  bool
}
```
- **`display`**: Reference to the Fyne input widget showing current values and results.
- **`op`**: Holds the pending arithmetic operator (`"+"`, `"-"`, `"×"`, `"÷"`, `"%"`).
- **`val1`**: Stores the first operand when an operation is triggered.
- **`newNum`**: Boolean flag indicating whether entering a digit should replace the display text or append to it.

#### 2. Operations & Error Handling
- **`input(digit)`**: Appends numbers or decimals to the display, preventing duplicate decimal points.
- **`setOp(op)`**: Saves the current display value into `val1`, assigns `op`, and sets `newNum = true`. Chained operations (e.g., `5 + 3 + 2`) compute intermediate results automatically.
- **`compute()`**: Parses operand 2, executes the pending operation via a `switch` statement, handles divide-by-zero (`val2 == 0`), and formats floating-point output cleanly with `strconv.FormatFloat`.
- **`toggleSign()`, `backspace()`, `clear()`**: Helper methods for sign inversion, character deletion, and calculator reset.

#### 3. Fyne UI Composition
- **Custom Button Builder (`makeButton`)**: Assigns Fyne importance levels (`HighImportance` for primary operators, `MediumImportance` for numbers, `LowImportance` for functions) to render distinct dark-mode button colors.
- **Grid Layout (`container.NewGridWithColumns(4, ...)`)**: Arranges 20 buttons in a 4-column responsive grid.
- **Border Layout (`container.NewBorder(...)`)**: Places the display entry at the top with padding and locks the button grid into the remaining space.

---

## 🌐 2. REST API

A thread-safe RESTful API built in Go using **zero external dependencies**, leveraging Go's standard library `net/http` package and Go 1.22+ routing features.

### ⚙️ Prerequisites & Running

- **Go 1.22+** installed on your system (`go version`).

#### 1. Run Directly from Source
```bash
cd rest-api
go run main.go
```
The server will start listening at `http://localhost:8080`.

#### 2. Build Executable Binary
```bash
cd rest-api
go build -o rest-api main.go
./rest-api
```

---

### 📡 API Endpoints & cURL Testing

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
- **`Book` Struct**: Defines fields (`ID`, `Title`, `Author`, `Price`) with `json:"..."` struct tags for JSON serialization.
- **Thread-Safe Store (`BookStore`)**:
  - Uses `sync.RWMutex` to guard access to `map[string]Book`.
  - `RLock()` / `RUnlock()`: Used for read-only operations (`GetAll()`, `GetByID()`) allowing concurrent read access without blocking.
  - `Lock()` / `Unlock()`: Used for write operations (`Create()`, `Update()`, `Delete()`) granting exclusive write access.
  - `defer` ensures mutex locks are always unlocked upon function exit.

#### 2. HTTP Request Handlers (`rest-api/handlers/book_handler.go`)
- **Dependency Injection**: `BookHandler` receives `*models.BookStore` pointer via constructor (`NewBookHandler`).
- **Response Formatting**: `writeJSON` sets `Content-Type: application/json` headers, assigns HTTP status codes, and encodes responses using `json.NewEncoder(w).Encode(data)`.
- **Path Matching**: `r.PathValue("id")` extracts URL path parameters using Go 1.22+ routing primitives.

#### 3. Routing & Server Setup (`rest-api/main.go`)
- Configures `http.NewServeMux()` with HTTP method matching (e.g. `"GET /api/books"`, `"POST /api/books"`).
- Binds TCP port `:8080` using `http.ListenAndServe(":8080", mux)`.

---

## 📊 3. GraphQL API

A lightweight, thread-safe GraphQL API built in Go using `github.com/graphql-go/graphql` with zero code-generation dependencies. Features full CRUD operations for Books and Authors, nested relationship resolvers, and an embedded interactive **GraphiQL IDE**.

### ✨ Features
- **Nested Relational Resolvers**: Fetch books along with author details (`Book.author`), or authors with their published books (`Author.books`).
- **Filtering & Arguments**: Filter books by `authorId` and/or `maxPrice`.
- **Full Mutations**: Create, update, and delete operations for both Books and Authors.
- **Embedded GraphiQL IDE**: Access `http://localhost:8080/` in browser for interactive query building.
- **Automated Test Suite**: Tested with `go test -v ./...`.

---

### ⚙️ Prerequisites & Running

```bash
cd graphql-api
go run main.go
```
The server will start listening at `http://localhost:8080`.

#### 1. Interactive GraphiQL IDE Browser Interface
Open your browser and navigate to:
👉 **`http://localhost:8080/`**

#### 2. Run Automated Unit Tests
```bash
cd graphql-api
go test -v ./...
```

#### 3. Build Executable Binary
```bash
cd graphql-api
go build -o graphql-api main.go
./graphql-api
```

---

### 📡 Example GraphQL cURL Commands

```bash
# 1. Fetch Books with Nested Author Data
curl -i -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query": "{ books { id title price author { name bio } } }"}'

# 2. Fetch Single Book by ID
curl -i -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query": "{ book(id: \"1\") { id title price author { name } } }"}'

# 3. Create a New Author (Mutation)
curl -i -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query": "mutation { createAuthor(name: \"Rob Pike\", bio: \"Co-creator of Go.\") { id name bio } }"}'

# 4. Create a New Book (Mutation)
curl -i -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query": "mutation { createBook(title: \"Go in Action\", authorId: \"1\", price: 34.99) { id title price author { name } } }"}'

# 5. Update a Book (Mutation)
curl -i -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query": "mutation { updateBook(id: \"1\", price: 42.50) { id title price } }"}'

# 6. Delete a Book (Mutation)
curl -i -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query": "mutation { deleteBook(id: \"1\") }"}'
```

---

### 📂 Detailed Code Explanations & Architecture

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

#### 1. Data Model & Concurrency (`graphql-api/models/store.go`)
- **`Author` & `Book` Structs**:
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
- **Thread-Safe Data Store (`Store`)**:
  - Encapsulates `authors map[string]Author` and `books map[string]Book` behind `sync.RWMutex`.
  - Provides concurrency-safe methods: `GetAllBooks()`, `GetBookByID()`, `GetBooksByAuthorID()`, `CreateBook()`, `UpdateBook()`, `DeleteBook()`, `GetAllAuthors()`, `GetAuthorByID()`, `CreateAuthor()`, `UpdateAuthor()`, `DeleteAuthor()`.

#### 2. GraphQL Schema & Resolvers (`graphql-api/schema/schema.go`)
- **GraphQL Object Types (`graphql.NewObject`)**:
  - Defines `AuthorType` and `BookType` with `graphql.FieldsThunk`.
  - **Nested Resolver for `Book.author`**: Looks up author in `store.GetAuthorByID(book.AuthorID)` when queried.
  - **Nested Resolver for `Author.books`**: Queries `store.GetBooksByAuthorID(author.ID)` to return all books written by the author.
- **Root Query (`Query`)**:
  - `books`: Supports optional `authorId` and `maxPrice` argument filters.
  - `book(id)`: Fetches a single book by ID.
  - `authors` & `author(id)`: Fetch all authors or single author.
- **Root Mutations (`Mutation`)**:
  - Provides mutations for creating, updating, and deleting books and authors.
- **Automated Tests (`schema/schema_test.go`)**:
  - Validates queries and mutations against `graphql.Do(...)` in memory without spawning a network server.

#### 3. HTTP Handler & GraphiQL IDE (`graphql-api/handlers/graphql.go`)
- **HTTP Routing**:
  - Handles incoming POST requests by parsing JSON payloads (`{"query": "...", "variables": {...}}`) into `graphql.Params` and executing via `graphql.Do(...)`.
  - Enables CORS (`Access-Control-Allow-Origin: *`).
  - Renders an embedded HTML/JS GraphiQL web interface when accessed via GET requests in a browser.

#### 4. Server Entry Point (`graphql-api/main.go`)
- Initializes `models.NewStore()`, constructs `schema.BuildSchema(store)`, wires up handlers on `http.NewServeMux()`, and binds port `:8080`.
