# Go_lang

A repository for Go language projects and applications.

---

## 📁 Projects Overview

1. **[🧮 1. Calculator (GUI)](#-1-calculator-gui)**: A desktop GUI calculator built using Go and [Fyne v2](https://fyne.io/).
2. **[🌐 2. REST API](#-2-rest-api)**: A lightweight, thread-safe RESTful API built in Go using standard library `net/http` and Go 1.22+ routing features.
3. **[📊 3. GraphQL API](#-3-graphql-api)**: A thread-safe GraphQL API built in Go with nested relational schema, CRUD mutations, unit tests, and an embedded **GraphiQL IDE**.
4. **[🔒 4. Authentication System](#-4-authentication-system)**: A secure authentication system in Go featuring **Bcrypt password hashing**, **HTTP-Only session cookies**, and a modern **Glassmorphic web interface** (HTML/CSS/JS).

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
./calculator
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
- **`Author` & `Book` Structs**: Encapsulates entity attributes.
- **Thread-Safe Data Store (`Store`)**: Encapsulates maps behind `sync.RWMutex` to allow concurrent safe operations across goroutines.

#### 2. GraphQL Schema & Resolvers (`graphql-api/schema/schema.go`)
- **GraphQL Object Types & Thunks**: Uses `graphql.FieldsThunk` to resolve circular references between `Author` and `Book`.
- **Relational Resolvers**: `Book.author` resolves author by ID; `Author.books` resolves all books by author ID.
- **Queries & Mutations**: Configures `Query` and `Mutation` root objects with arguments and resolver handlers.

#### 3. HTTP Handler & GraphiQL IDE (`graphql-api/handlers/graphql.go`)
- **HTTP Routing**: Parses POST JSON query bodies into `graphql.Params` and executes via `graphql.Do(...)`.
- **GraphiQL Interface**: Detects GET requests in browsers to serve an embedded GraphiQL interface.

#### 4. Server Entry Point (`graphql-api/main.go`)
- Initializes data store, builds schema, registers HTTP handlers on `http.NewServeMux()`, and listens on port `:8080`.

---

## 🔒 4. Authentication System

A secure authentication system in Go featuring **Bcrypt password hashing**, **HTTP-Only session cookies**, and a modern **Glassmorphic web interface** (HTML/CSS/JS).

### ✨ Features
- **Bcrypt Password Security**: Hashes passwords with `golang.org/x/crypto/bcrypt`.
- **HTTP-Only Session Cookies**: Mitigates XSS token theft using `HttpOnly` flags and server-side session tokens.
- **Real-time Password Strength Meter**: Displays visual feedback (Weak, Moderate, Strong) as users type.
- **Interactive Web Interface**: Custom Glassmorphism UI with HTML5, CSS backdrop filters, glowing orb animations, and Vanilla JS.
- **Demo Credentials Ready**: Quick fill button for instant login (`demo` / `password123`).

---

### ⚙️ Prerequisites & Running

```bash
cd auth-system
go run main.go
```
The server will start listening at `http://localhost:8080`.

#### 1. Interactive Web Application
Open your browser and navigate to:
👉 **`http://localhost:8080/`**

#### 2. Run Automated Unit Tests
```bash
cd auth-system
go test -v ./...
```

#### 3. Build Executable Binary
```bash
cd auth-system
go build -o auth-system main.go
./auth-system
```

---

### 📡 API Endpoints & cURL Testing

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `http://localhost:8080/api/signup` | Register a new user account |
| `POST` | `http://localhost:8080/api/login` | Authenticate user & set session cookie |
| `POST` | `http://localhost:8080/api/logout` | Invalidate session & clear cookie |
| `GET` | `http://localhost:8080/api/me` | Fetch authenticated user profile |
| `GET` | `http://localhost:8080/` | Serve Web UI |

#### cURL Test Commands

```bash
# 1. Login with Demo Account & Save Cookie
curl -i -c cookies.txt -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"usernameOrEmail":"demo", "password":"password123"}'

# 2. Check Session Status (/api/me)
curl -i -b cookies.txt http://localhost:8080/api/me

# 3. Register a New User
curl -i -c cookies.txt -X POST http://localhost:8080/api/signup \
  -H "Content-Type: application/json" \
  -d '{"username":"alice", "email":"alice@example.com", "password":"securepassword123", "fullName":"Alice Smith"}'

# 4. Logout
curl -i -b cookies.txt -X POST http://localhost:8080/api/logout
```

---

### 📂 Detailed Code Explanations & Architecture

```
auth-system/
├── go.mod                # Module definition file
├── go.sum                # Dependencies checksum
├── main.go               # Application entry point & HTTP multiplexer
├── auth_test.go          # Unit & integration tests
├── models/
│   └── user.go           # User/Session models, bcrypt hashing & store
├── handlers/
│   └── auth_handler.go   # HTTP handlers for auth endpoints
├── templates/
│   └── index.html        # Glassmorphic HTML template
├── static/
│   ├── style.css         # Custom CSS design system
│   └── app.js            # Client-side SPA interaction script
└── README.md             # Project documentation
```

#### 1. Password Hashing & Sessions (`auth-system/models/user.go`)
- **Bcrypt Password Security**: Hashes raw passwords upon signup using `bcrypt.GenerateFromPassword(..., bcrypt.DefaultCost)` and verifies credentials via `bcrypt.CompareHashAndPassword(...)`.
- **Session Tokens**: Generates 32-byte cryptographically secure tokens (`crypto/rand`) valid for 24 hours.

#### 2. HTTP Controller Handlers (`auth-system/handlers/auth_handler.go`)
- **HTTP-Only Cookies**: Sets `HttpOnly: true`, `SameSite: Lax`, and expiration on session cookies (`session_token`) upon login/signup to shield tokens from client script access.
- **Session Invalidation**: Clears cookies and deletes session tokens from the server store during logout.

#### 3. Web Interface UI (`templates/index.html` & `static/`)
- **`templates/index.html`**: Form markup for login, signup, password strength bar, and user profile dashboard.
- **`static/style.css`**: Backdrop-filter glassmorphism, responsive cards, glowing orb keyframe animations, and status badges.
- **`static/app.js`**: Dynamic tab switching, realtime password score calculation, AJAX API request handling, and toast notifications.
