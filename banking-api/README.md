# 🏦 Go Echo Banking REST API

A financial Banking RESTful API built using the **Echo Web Framework** (`github.com/labstack/echo/v4`), **GORM ORM** (`gorm.io/gorm`), and **SQLite** (`gorm.io/driver/sqlite`).

---

## 🏛️ Architecture Overview

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

---

## 💡 Echo Framework Explanation

**Echo** is a minimalist, ultra-fast, and highly extensible Go web framework.
- **Context Abstraction (`echo.Context`)**: Simplifies JSON response serialization (`c.JSON`), request body binding (`c.Bind`), and route parameter retrieval (`c.Param`).
- **Middleware Ecosystem**: Ships with optimized built-in middleware for logging (`middleware.Logger()`), panic recovery (`middleware.Recover()`), and CORS (`middleware.CORSWithConfig()`).
- **Fast HTTP Router**: High performance matching with zero memory allocations on hit.

---

## ✨ Features

- **ACID Transactional Fund Transfer**: Inter-account fund transfers execute inside a GORM transaction with pessimistic row locking (`FOR UPDATE`) to prevent race conditions and double-spending.
- **Ledger Statements & Audit Trails**: Every financial movement creates an immutable `Transaction` record and writes security `AuditLog` events.
- **Pre-seeded Accounts**: Demo accounts automatically created:
  - Account 1 (`ACC1001`): Initial Balance `$5,000.00`
  - Account 2 (`ACC1002`): Initial Balance `$1,200.00`
- **Automated Unit Tests**: Tested with `go test -v ./...`.

---

## ⚙️ How to Run & Test

### 1. Run Server
```bash
cd banking-api
go run main.go
```
The server will start listening at `http://localhost:8082`.

### 2. Run Automated Unit Tests
```bash
cd banking-api
go test -v ./...
```

---

## 📡 API Endpoints & cURL Testing

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/api/v1/accounts` | Create a new bank account |
| `GET` | `/api/v1/accounts` | List all bank accounts |
| `GET` | `/api/v1/accounts/:id` | Fetch account balance & status |
| `POST` | `/api/v1/accounts/:id/deposit` | Credit funds to account |
| `POST` | `/api/v1/accounts/:id/withdraw` | Debit funds with balance validation |
| `POST` | `/api/v1/transfer` | Execute atomic fund transfer between accounts |
| `GET` | `/api/v1/accounts/:id/transactions` | Fetch account ledger transaction statement |

### cURL Examples

#### 1. Fetch All Bank Accounts
```bash
curl -i http://localhost:8082/api/v1/accounts
```

#### 2. Deposit Funds into Account 1 ($500.00)
```bash
curl -i -X POST http://localhost:8082/api/v1/accounts/1/deposit \
  -H "Content-Type: application/json" \
  -d '{"amount": 500.00, "description": "Salary Deposit"}'
```

#### 3. Withdraw Funds from Account 1 ($200.00)
```bash
curl -i -X POST http://localhost:8082/api/v1/accounts/1/withdraw \
  -H "Content-Type: application/json" \
  -d '{"amount": 200.00, "description": "ATM Cash"}'
```

#### 4. Execute Atomic Fund Transfer (Account 1 -> Account 2, Amount: $1000.00)
```bash
curl -i -X POST http://localhost:8082/api/v1/transfer \
  -H "Content-Type: application/json" \
  -d '{"fromAccountId": 1, "toAccountId": 2, "amount": 1000.00, "description": "Rent Payment"}'
```

#### 5. View Transaction Statement for Account 1
```bash
curl -i http://localhost:8082/api/v1/accounts/1/transactions
```
