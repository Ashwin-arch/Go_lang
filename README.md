# Go_lang

A repository for Go language projects and applications.

---

## 📁 Projects Overview

1. **[🧮 1. Calculator (GUI)](#-1-calculator-gui)**: A desktop GUI calculator built using Go and [Fyne v2](https://fyne.io/).
2. **[🌐 2. REST API](#-2-rest-api)**: A lightweight, thread-safe RESTful API built in Go using standard library `net/http` and Go 1.22+ routing features.
3. **[📊 3. GraphQL API](#-3-graphql-api)**: A thread-safe GraphQL API built in Go with nested relational schema, CRUD mutations, unit tests, and an embedded **GraphiQL IDE**.
4. **[🔒 4. Authentication System](#-4-authentication-system)**: A secure authentication system in Go featuring **Bcrypt password hashing**, **HTTP-Only session cookies**, and a modern **Glassmorphic web interface** (HTML/CSS/JS).
5. **[🛒 5. E-commerce REST API (Gin)](#-5-e-commerce-rest-api-gin)**: A full-featured e-commerce backend built with the **Gin Framework**, **GORM**, and **SQLite**, featuring atomic database transactions for checkout.
6. **[🏦 6. Banking REST API (Echo)](#-6-banking-rest-api-echo)**: A financial banking backend built with the **Echo Framework**, **GORM**, and **SQLite**, featuring ACID transactional transfers with pessimistic row locking (`FOR UPDATE`) and audit logging.

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

---

## 📊 3. GraphQL API

A lightweight, thread-safe GraphQL API built in Go using `github.com/graphql-go/graphql` with zero code-generation dependencies. Features full CRUD operations for Books and Authors, nested relationship resolvers, and an embedded interactive **GraphiQL IDE**.

---

## 🔒 4. Authentication System

A secure authentication system in Go featuring **Bcrypt password hashing**, **HTTP-Only session cookies**, and a modern **Glassmorphic web interface** (HTML/CSS/JS).

---

## 🛒 5. E-commerce REST API (Gin)

A production-grade E-commerce RESTful API built using the **Gin Web Framework** (`github.com/gin-gonic/gin`), **GORM ORM** (`gorm.io/gorm`), and **SQLite** (`gorm.io/driver/sqlite`).

### 🏛️ Architecture & Gin Framework

```
Client (cURL / Web / Postman)
   │
   ▼
Gin Web Framework (Radix Tree Router, Middlewares, Context Controllers)
   │
   ▼
GORM ORM Layer (Database Transactions, Model Relations, Auto-Migration)
   │
   ▼
SQLite Database (Relational Store for Products, Categories, Cart, Orders)
```

- **Gin Framework**: Utilizes Gin's Radix tree router for high-performance HTTP request handling, custom middleware groups (`AuthMiddleware`, `AdminMiddleware`), and request binding validation.
- **Atomic Checkout Transaction**: Checkout is executed inside a `db.Transaction(...)` block that validates product stock, deducts inventory, creates order items, and empties the shopping cart atomically.

### 🚀 How to Run & Test
```bash
cd ecommerce-api
go run main.go
```
Server runs at `http://localhost:8081`.

#### Run Unit Tests
```bash
cd ecommerce-api
go test -v ./...
```

### 📡 cURL Example
```bash
# Add product to cart
curl -i -X POST http://localhost:8081/api/v1/cart/items \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 1" \
  -d '{"productId": 1, "quantity": 2}'

# Execute Atomic Checkout
curl -i -X POST http://localhost:8081/api/v1/orders/checkout \
  -H "X-User-ID: 1"
```

---

## 🏦 6. Banking REST API (Echo)

A financial Banking RESTful API built using the **Echo Web Framework** (`github.com/labstack/echo/v4`), **GORM ORM** (`gorm.io/gorm`), and **SQLite** (`gorm.io/driver/sqlite`).

### 🏛️ Architecture & Echo Framework

```
Client (Mobile App / Web / ATM / cURL)
   │
   ▼
Echo Web Server (Context Binding, Echo Middlewares & Handlers)
   │
   ▼
GORM ORM Layer (ACID Transactions, Pessimistic Row Locking, Audit Ledger)
   │
   ▼
SQLite Database (Persistent Relational Store for Accounts & Financial Transactions)
```

- **Echo Framework**: Features Echo's context abstraction (`echo.Context`), built-in logger & recovery middleware, and clean JSON controller signatures.
- **ACID Transactional Fund Transfer**: Inter-account fund transfers execute inside a GORM transaction with pessimistic row locking (`FOR UPDATE`) to prevent race conditions and double-spending. Audit log records are created for every financial transaction.

### 🚀 How to Run & Test
```bash
cd banking-api
go run main.go
```
Server runs at `http://localhost:8082`.

#### Run Unit Tests
```bash
cd banking-api
go test -v ./...
```

### 📡 cURL Example
```bash
# Deposit Funds into Account 1 ($500.00)
curl -i -X POST http://localhost:8082/api/v1/accounts/1/deposit \
  -H "Content-Type: application/json" \
  -d '{"amount": 500.00, "description": "Salary Deposit"}'

# Execute Atomic Fund Transfer (Account 1 -> Account 2, Amount: $1000.00)
curl -i -X POST http://localhost:8082/api/v1/transfer \
  -H "Content-Type: application/json" \
  -d '{"fromAccountId": 1, "toAccountId": 2, "amount": 1000.00, "description": "Rent Payment"}'
```
