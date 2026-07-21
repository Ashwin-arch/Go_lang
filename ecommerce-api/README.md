# 🛒 Go Gin E-commerce REST API

A production-grade E-commerce RESTful API built using the **Gin Web Framework** (`github.com/gin-gonic/gin`), **GORM ORM** (`gorm.io/gorm`), and **SQLite** (`gorm.io/driver/sqlite`).

---

## 🏛️ Architecture Overview

```
Client (cURL / Web / Postman)
   │
   ▼
Gin Web Framework (Routers, Middleware, Request Binding & JSON Controllers)
   │
   ▼
GORM ORM Layer (Database Transactions, Model Relations, Auto-Migration)
   │
   ▼
SQLite Database (Persistent relational store for Products, Orders & Users)
```

---

## 💡 Gin Framework Explanation

**Gin** is one of the fastest and most popular web frameworks in Go.
- **Radix Tree Routing**: Gin uses a custom Radix tree router for extremely fast HTTP request routing with minimal memory allocations.
- **Middleware Chain**: Handlers can be chained with middleware using `r.Use()` or group-specific `.Use()`.
- **Context Management (`*gin.Context`)**: Encapsulates HTTP request parsing (`ShouldBindJSON`), parameter extraction (`Param`, `Query`), and JSON formatting (`c.JSON`).

---

## ✨ Features

- **Relational Data Modeling**: GORM models for Users, Categories, Products, CartItems, Orders, and OrderItems.
- **Atomic Checkout Transaction**: Checkout is executed inside a `db.Transaction(...)` block that validates stock, deducts product inventory, creates order items, and empties the shopping cart atomically.
- **Gin Middleware**: Includes CORS handler, Authentication middleware (`AuthMiddleware`), and Admin Authorization (`AdminMiddleware`).
- **Auto-Migration & Pre-seeded Data**: Automatically initializes database tables and seeds demo products, categories, customer account (`customer@example.com`), and admin account (`admin@example.com`).

---

## ⚙️ How to Run & Test

### 1. Run Server
```bash
cd ecommerce-api
go run main.go
```
The server will start listening at `http://localhost:8081`.

### 2. Run Automated Unit Tests
```bash
cd ecommerce-api
go test -v ./...
```

---

## 📡 API Endpoints & cURL Testing

| Method | Endpoint | Access | Description |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/auth/register` | Public | Register new user account |
| `POST` | `/api/v1/auth/login` | Public | Authenticate user & get ID token |
| `GET` | `/api/v1/products` | Public | Fetch products (optional `?categoryId=` or `?q=`) |
| `GET` | `/api/v1/products/:id` | Public | Fetch single product by ID |
| `POST` | `/api/v1/cart/items` | Customer | Add item to cart (`X-User-ID: 1`) |
| `GET` | `/api/v1/cart` | Customer | View items in shopping cart (`X-User-ID: 1`) |
| `POST` | `/api/v1/orders/checkout` | Customer | Atomic checkout & place order (`X-User-ID: 1`) |
| `GET` | `/api/v1/orders` | Customer | View placed order history (`X-User-ID: 1`) |
| `POST` | `/api/v1/admin/products` | Admin | Create product (`X-User-ID: 2`) |

### cURL Examples

#### 1. Fetch All Products
```bash
curl -i http://localhost:8081/api/v1/products
```

#### 2. Add Product to Shopping Cart
```bash
curl -i -X POST http://localhost:8081/api/v1/cart/items \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 1" \
  -d '{"productId": 1, "quantity": 2}'
```

#### 3. View Shopping Cart
```bash
curl -i -H "X-User-ID: 1" http://localhost:8081/api/v1/cart
```

#### 4. Execute Atomic Checkout
```bash
curl -i -X POST http://localhost:8081/api/v1/orders/checkout \
  -H "X-User-ID: 1"
```
