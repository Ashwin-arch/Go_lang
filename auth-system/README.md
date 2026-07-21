# 🔒 Go Authentication System

A secure, modern Authentication System built in Go featuring **Bcrypt password hashing**, **HTTP-Only session cookie management**, thread-safe user storage, SVG vector icon UI, and a vibrant **Glassmorphic web interface** built with HTML, CSS, and Vanilla JavaScript.

---

## 📋 Table of Contents
1. [Features](#-features)
2. [How to Run & Build](#-how-to-run--build)
3. [Interactive Web UI](#-interactive-web-ui)
4. [API Endpoints & cURL Testing](#-api-endpoints--curl-testing)
5. [Comprehensive Code Explanations & Implementation Guide](#-comprehensive-code-explanations--implementation-guide)
   - [1. User Model & Bcrypt Security (`models/user.go`)](#1-user-model--bcrypt-security-modelsusergo)
   - [2. HTTP Controller Handlers (`handlers/auth_handler.go`)](#2-http-controller-handlers-handlersauth_handlergo)
   - [3. Server Setup & Route Multiplexing (`main.go`)](#3-server-setup--route-multiplexing-maingo)
   - [4. Glassmorphism Web Interface (`templates/index.html`)](#4-glassmorphism-web-interface-templatesindexhtml)
   - [5. CSS Design System & Keyframe Animations (`static/style.css`)](#5-css-design-system--keyframe-animations-staticstylecss)
   - [6. Client-side Application Logic & SVG Toggling (`static/app.js`)](#6-client-side-application-logic--svg-toggling-staticappjs)
   - [7. Automated Unit & Integration Tests (`auth_test.go`)](#7-automated-unit--integration-tests-auth_testgo)

---

## ✨ Features

- **Bcrypt Password Security**: Passwords are securely hashed using `golang.org/x/crypto/bcrypt` before persistence.
- **HTTP-Only Session Cookies**: Prevents XSS cookie theft using `HttpOnly`, `SameSite=Lax`, and 24-hour expiration tokens.
- **Thread-Safe Memory Store**: Guarded by Go's `sync.RWMutex` for concurrent request handling.
- **Real-Time Password Strength Meter**: Dynamic visual indicator (Weak, Moderate, Strong) as users type.
- **Resolution-Independent SVG Vector Icons**: Pure SVG vector icons used throughout inputs, buttons, badges, and alerts.
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
- **Password Visibility Toggle** with SVG Eye / Eye-Off icons
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

## 📂 Comprehensive Code Explanations & Implementation Guide

---

### 1. User Model & Bcrypt Security (`models/user.go`)

#### A. User & Session Struct Definitions
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
- **`PasswordHash` Tag (`json:"-"`)**: Ensures password hashes are never exposed in JSON API responses.
- **`Session`**: Links a unique 32-byte session token to a user ID with an expiration timestamp (`ExpiresAt`).

#### B. Bcrypt Hashing Implementation
```go
// RegisterUser hashes password before storing
hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
if err != nil {
    return User{}, fmt.Errorf("failed to hash password: %w", err)
}
user.PasswordHash = string(hashedBytes)
```
- Uses `bcrypt.GenerateFromPassword` with adaptive salt cost factor (`bcrypt.DefaultCost = 10`), rendering password hashes immune to rainbow table attacks.

```go
// AuthenticateUser verifies bcrypt hash
err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
if err != nil {
    return User{}, errors.New("invalid credentials")
}
```
- `bcrypt.CompareHashAndPassword` securely checks candidate passwords in constant time, preventing timing side-channel attacks.

#### C. Thread-Safe Store (`AuthStore`)
```go
type AuthStore struct {
    mu         sync.RWMutex
    users      map[string]User    // ID -> User
    usernames  map[string]string  // Username -> User ID
    emails     map[string]string  // Email -> User ID
    sessions   map[string]Session // Session Token -> Session
    nextUserID int
}
```
- **Concurrency Safety**: Guards maps with `sync.RWMutex`.
- **Read Operations (`GetSession`, `AuthenticateUser`)**: Acquire `s.mu.RLock()` and `defer s.mu.RUnlock()`.
- **Write Operations (`RegisterUser`, `CreateSession`, `DeleteSession`)**: Acquire `s.mu.Lock()` and `defer s.mu.Unlock()`.

---

### 2. HTTP Controller Handlers (`handlers/auth_handler.go`)

#### A. JSON Response Formatting Helper
```go
func writeJSON(w http.ResponseWriter, status int, resp apiResponse) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(resp)
}
```
- Configures response headers, HTTP status code (200, 201, 400, 401), and streams JSON payload.

#### B. Setting HTTP-Only Session Cookies
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
- **`HttpOnly: true`**: Crucial security flag that blocks client-side scripts (`document.cookie`) from reading the cookie, mitigating XSS token theft.
- **`SameSite: Lax`**: Restricts cross-site cookie transmission, defending against CSRF attacks.

#### C. Logout Session Invalidation
```go
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("session_token")
    if err == nil && cookie.Value != "" {
        h.Store.DeleteSession(cookie.Value)
    }

    http.SetCookie(w, &http.Cookie{
        Name:     "session_token",
        Value:    "",
        Path:     "/",
        Expires:  time.Unix(0, 0),
        HttpOnly: true,
        MaxAge:   -1,
    })
    writeJSON(w, http.StatusOK, apiResponse{Success: true, Message: "Logged out successfully."})
}
```
- Deletes session token from server-side store and expires client cookie immediately (`MaxAge: -1`).

---

### 3. Server Setup & Route Multiplexing (`main.go`)

```go
mux := http.NewServeMux()

// Static assets handler with explicit verb matching
fileServer := http.FileServer(http.Dir("static"))
mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

// API Endpoints
mux.HandleFunc("POST /api/signup", authHandler.Signup)
mux.HandleFunc("POST /api/login", authHandler.Login)
mux.HandleFunc("POST /api/logout", authHandler.Logout)
mux.HandleFunc("GET /api/me", authHandler.Me)

// HTML UI Page
mux.HandleFunc("GET /", authHandler.ServeIndex)
```
- **Go 1.22+ Standard Routing**: Explicit method prefixes (`"POST /api/login"`, `"GET /static/"`, `"GET /"`) guarantee strict verb matching without route pattern conflicts.

---

### 4. Glassmorphic Web Interface (`templates/index.html`)

#### A. Header Brand & SVG Vector Logo
```html
<div class="logo-icon">
  <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
    <rect x="3" y="11" width="18" height="11" rx="2.5" ry="2.5"></rect>
    <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
  </svg>
</div>
```
- Inline SVG vectors scale crisp and sharp on retina screens without image pixelation.

#### B. Input Field Wrapper with SVG Icons
```html
<div class="input-wrapper">
  <span class="input-icon">
    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
      <circle cx="12" cy="7" r="4"></circle>
    </svg>
  </span>
  <input type="text" id="login-username" placeholder="demo or demo@example.com" required>
</div>
```

---

### 5. CSS Design System & Keyframe Animations (`static/style.css`)

#### A. CSS Custom Properties
```css
:root {
  --bg-dark: #0b0f19;
  --card-bg: rgba(23, 32, 54, 0.65);
  --card-border: rgba(255, 255, 255, 0.12);
  --primary: #6366f1;
  --primary-hover: #4f46e5;
  --primary-glow: rgba(99, 102, 241, 0.4);
  --accent: #ec4899;
}
```

#### B. Glassmorphism & Animated Glowing Orbs
```css
.glass-card {
  background: var(--card-bg);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid var(--card-border);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
}

@keyframes float {
  0% { transform: translate(0, 0) scale(1); }
  50% { transform: translate(60px, 40px) scale(1.1); }
  100% { transform: translate(-40px, 60px) scale(0.95); }
}
```

---

### 6. Client-side Application Logic & SVG Toggling (`static/app.js`)

#### A. Dynamic Password Visibility Toggle
```javascript
function togglePasswordVisibility(inputId, btn) {
  const input = document.getElementById(inputId);
  if (input.type === 'password') {
    input.type = 'text';
    btn.innerHTML = `<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"></path><line x1="1" y1="1" x2="23" y2="23"></line></svg>`;
  } else {
    input.type = 'password';
    btn.innerHTML = `<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path><circle cx="12" cy="12" r="3"></circle></svg>`;
  }
}
```

#### B. Real-time Password Strength Meter Algorithm
```javascript
function checkPasswordStrength(password) {
  let score = 0;
  if (password.length >= 6) score++;
  if (password.length >= 10) score++;
  if (/[0-9]/.test(password)) score++;
  if (/[^a-zA-Z0-9]/.test(password)) score++;

  if (score <= 1) {
    bar.style.width = '25%'; bar.style.backgroundColor = '#ef4444';
  } else if (score === 2 || score === 3) {
    bar.style.width = '65%'; bar.style.backgroundColor = '#f59e0b';
  } else {
    bar.style.width = '100%'; bar.style.backgroundColor = '#10b981';
  }
}
```

---

### 7. Automated Unit & Integration Tests (`auth_test.go`)

```go
func TestAuthSystem(t *testing.T) {
    store := models.NewAuthStore()
    authHandler := handlers.NewAuthHandler(store, "templates")

    mux := http.NewServeMux()
    mux.HandleFunc("POST /api/signup", authHandler.Signup)
    mux.HandleFunc("POST /api/login", authHandler.Login)
    mux.HandleFunc("GET /api/me", authHandler.Me)

    // Test Login & Session Cookie
    body, _ := json.Marshal(map[string]string{"usernameOrEmail": "demo", "password": "password123"})
    req := httptest.NewRequest("POST", "/api/login", bytes.NewBuffer(body))
    w := httptest.NewRecorder()

    mux.ServeHTTP(w, req)
    if w.Code != http.StatusOK {
        t.Fatalf("Expected 200, got %d", w.Code)
    }
}
```
- Tests user signup, authentication, password verification, cookie issuance, and session expiration via `httptest.NewRecorder()`.
