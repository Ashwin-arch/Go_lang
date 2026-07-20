# Go_lang

A repository for Go language projects and applications.

---

## 📁 Projects

1. **[🧮 Calculator (GUI)](./calculator)**: A desktop GUI calculator built using Go and [Fyne v2](https://fyne.io/).
2. **[🌐 REST API](./rest-api)**: A lightweight, thread-safe RESTful API built in Go using standard library `net/http`.

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

### Features
- Full CRUD operations for a Book inventory (`GET`, `POST`, `PUT`, `DELETE`)
- Thread-safe in-memory data storage using `sync.RWMutex`
- Standard library routing (`net/http`) without third-party dependencies

### How to Run REST API
```bash
cd rest-api
go run main.go
```
The server will run on `http://localhost:8080`.

#### Test Endpoint Example
```bash
curl http://localhost:8080/api/books
```
