# 🔒 Go Authentication System

A secure, modern Authentication System built in Go featuring **Bcrypt password hashing**, **HTTP-Only session cookie management**, thread-safe user storage, and a vibrant **Glassmorphic web interface** built with HTML, CSS, and Vanilla JavaScript.

---

## 📋 Table of Contents
1. [Features](#-features)
2. [How to Run & Build](#-how-to-run--build)
3. [Interactive Web UI](#-interactive-web-ui)
4. [API Endpoints & cURL Testing](#-api-endpoints--curl-testing)
5. [Detailed Code Explanations & Architecture](#-detailed-code-explanations--architecture)
   - [1. User Model & Bcrypt Hashing (`models/user.go`)](#1-user-model--bcrypt-hashing-modelsusergo)
   - [2. HTTP Handlers & Session Cookies (`handlers/auth_handler.go`)](#2-http-handlers--session-cookies-handlersauth_handlergo)
   - [3. Modern Glassmorphism Web Interface (`templates/index.html` & `static/`)](#3-modern-glassmorphism-web-interface-templatesindexhtml--static)
   - [4. Server Setup & Routing (`main.go`)](#4-server-setup--routing-maingo)

---

## ✨ Features

- **Bcrypt Password Security**: Passwords are securely hashed using `golang.org/x/crypto/bcrypt` before persistence.
- **HTTP-Only Session Cookies**: Prevents XSS cookie theft using `HttpOnly`, `SameSite=Lax`, and 24-hour expiration tokens.
- **Thread-Safe Memory Store**: Guarded by Go's `sync.RWMutex` for concurrent request handling.
- **Real-Time Password Strength Meter**: Dynamic visual indicator (Weak, Moderate, Strong) as users type.
- **Glassmorphic Web Interface**: Custom dark-mode theme built with CSS backdrop-filter blur, Google Fonts (Outfit), glowing animated background orbs, and micro-interactions.
- **Pre-configured Demo Account**: Ready to test out of the box (`demo` / `password123`).
- **Automated Test Suite**: Tested with `go test -v ./...`.

---

## 🚀 How to Run & Build

### 1. Run Directly from Source
```bash
cd auth-system
go run main.go
```
The server will start listening at `http://localhost:8080`.

### 2. Run Automated Unit Tests
```bash
cd auth-system
go test -v ./...
```

### 3. Build Executable Binary
```bash
cd auth-system
go build -o auth-system main.go
./auth-system
```

---

## 🎨 Interactive Web UI

Open your browser and navigate to:
👉 **`http://localhost:8080/`**

Features included in the Web UI:
- **Sign In / Sign Up View Toggling**
- **Quick Fill Button** for instant Demo login (`demo` / `password123`)
- **Password Visibility Toggle** (👁️ / 🙈)
- **Real-time Password Strength Meter**
- **Authenticated Dashboard** displaying User Profile details and Session Status
- **Logout Action** that invalidates server sessions and clears client cookies

---

## 📡 API Endpoints & cURL Testing

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `http://localhost:8080/api/signup` | Register a new user account |
| `POST` | `http://localhost:8080/api/login` | Authenticate user & receive HTTP-only session cookie |
| `POST` | `http://localhost:8080/api/logout` | Invalidate session & clear cookie |
| `GET` | `http://localhost:8080/api/me` | Fetch authenticated user profile |
| `GET` | `http://localhost:8080/` | Serve Web UI |

### cURL Test Examples

#### 1. Login with Demo Account
```bash
curl -i -c cookies.txt -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"usernameOrEmail":"demo", "password":"password123"}'
```

#### 2. Check Session Status (`/api/me`)
```bash
curl -i -b cookies.txt http://localhost:8080/api/me
```

#### 3. Register a New User
```bash
curl -i -c cookies.txt -X POST http://localhost:8080/api/signup \
  -H "Content-Type: application/json" \
  -d '{"username":"alice", "email":"alice@example.com", "password":"securepassword123", "fullName":"Alice Smith"}'
```

#### 4. Logout
```bash
curl -i -b cookies.txt -X POST http://localhost:8080/api/logout
```

---

## 📂 Detailed Code Explanations & Architecture

```
auth-system/
├── go.mod                # Module definition file
├── go.sum                # Dependencies checksum
├── main.go               # Server entry point & HTTP multiplexer
├── auth_test.go          # Integration & unit tests
├── models/
│   └── user.go           # User/Session models, bcrypt hashing & store
├── handlers/
│   └── auth_handler.go   # HTTP API controller handlers
├── templates/
│   └── index.html        # Glassmorphic HTML template
├── static/
│   ├── style.css         # Custom CSS design system
│   └── app.js            # Client-side SPA interaction script
└── README.md             # Project documentation
```

---

### 1. User Model & Bcrypt Hashing (`models/user.go`)

#### A. User & Session Structs
```go
type User struct {
    ID           string    `json:"id"`
    Username     string    `json:"username"`
    Email        string    `json:"email"`
    PasswordHash string    `json:"-"` // Omitted from JSON output for security
    FullName     string    `json:"fullName"`
    CreatedAt    time.Time `json:"createdAt"`
}

type Session struct {
    Token     string    `json:"token"`
    UserID    string    `json:"userId"`
    ExpiresAt time.Time `json:"expiresAt"`
}
```

#### B. Bcrypt Hashing & Verification
- **Registration (`RegisterUser`)**:
  ```go
  hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
  ```
  Hashes raw passwords with adaptive salt cost factor before storing in memory.
- **Authentication (`AuthenticateUser`)**:
  ```go
  err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
  ```
  Compares candidate password against stored hash safely resisting timing attacks.

#### C. Session Management
- **`CreateSession`**: Generates a 32-byte cryptographically secure random token (`crypto/rand`) converted to hex string, valid for 24 hours.
- **`GetSession`**: Verifies token existence and expiration (`time.Now().After(session.ExpiresAt)`).

---

### 2. HTTP Handlers & Session Cookies (`handlers/auth_handler.go`)

- **Setting Cookies (`http.SetCookie`)**:
  ```go
  http.SetCookie(w, &http.Cookie{
      Name:     "session_token",
      Value:    session.Token,
      Path:     "/",
      Expires:  session.ExpiresAt,
      HttpOnly: true,
      SameSite: http.SameSiteLaxMode,
  })
  ```
  `HttpOnly: true` prevents client-side JavaScript from accessing session cookies, protecting against XSS token stealing.
- **Clearing Cookies on Logout**: Overwrites cookie value to `""` with `MaxAge: -1` and expired timestamp.

---

### 3. Modern Glassmorphism Web Interface (`templates/index.html` & `static/`)

- **HTML Structure (`templates/index.html`)**: Semantic markup with tabbed navigation, password strength visualizer, and dashboard view.
- **CSS Styling (`static/style.css`)**: Implements CSS backdrop-filter glassmorphism, animated glowing background globes (`@keyframes float`), and clean dark mode typography.
- **JS Application Logic (`static/app.js`)**: Evaluates password complexity in real time, handles async `fetch` login/signup calls, updates DOM views dynamically, and displays toast alerts.

---

### 4. Server Setup & Routing (`main.go`)

```go
mux := http.NewServeMux()

fileServer := http.FileServer(http.Dir("static"))
mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

mux.HandleFunc("POST /api/signup", authHandler.Signup)
mux.HandleFunc("POST /api/login", authHandler.Login)
mux.HandleFunc("POST /api/logout", authHandler.Logout)
mux.HandleFunc("GET /api/me", authHandler.Me)
mux.HandleFunc("GET /", authHandler.ServeIndex)
```
- Uses Go 1.22+ standard HTTP method routing (`"POST /api/login"`) to restrict endpoints to matching verbs.
