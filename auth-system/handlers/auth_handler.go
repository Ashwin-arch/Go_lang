package handlers

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"time"

	"auth-system/models"
)

// AuthHandler handles authentication endpoints and serves UI templates
type AuthHandler struct {
	Store        *models.AuthStore
	TemplatePath string
}

// NewAuthHandler creates an AuthHandler instance
func NewAuthHandler(store *models.AuthStore, templatePath string) *AuthHandler {
	return &AuthHandler{
		Store:        store,
		TemplatePath: templatePath,
	}
}

// Response structs
type apiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	User    *models.User `json:"user,omitempty"`
}

type signupRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"fullName"`
}

type loginRequest struct {
	UsernameOrEmail string `json:"usernameOrEmail"`
	Password        string `json:"password"`
}

func writeJSON(w http.ResponseWriter, status int, resp apiResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

// ServeIndex serves the main index.html template
func (h *AuthHandler) ServeIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/index.html" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, filepath.Join(h.TemplatePath, "index.html"))
}

// Signup handles user registration requests
func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, apiResponse{Success: false, Message: "Method not allowed"})
		return
	}

	var req signupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Message: "Invalid JSON request body"})
		return
	}

	user, err := h.Store.RegisterUser(req.Username, req.Email, req.Password, req.FullName)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Message: err.Error()})
		return
	}

	// Auto login user after successful signup
	session, err := h.Store.CreateSession(user.ID)
	if err == nil {
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    session.Token,
			Path:     "/",
			Expires:  session.ExpiresAt,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})
	}

	writeJSON(w, http.StatusCreated, apiResponse{
		Success: true,
		Message: "Registration successful! Welcome aboard.",
		User:    &user,
	})
}

// Login handles user authentication and sets session cookie
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, apiResponse{Success: false, Message: "Method not allowed"})
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, apiResponse{Success: false, Message: "Invalid JSON request body"})
		return
	}

	user, err := h.Store.AuthenticateUser(req.UsernameOrEmail, req.Password)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, apiResponse{Success: false, Message: err.Error()})
		return
	}

	session, err := h.Store.CreateSession(user.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, apiResponse{Success: false, Message: "Failed to create session"})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.Token,
		Path:     "/",
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Message: "Welcome back! Login successful.",
		User:    &user,
	})
}

// Logout clears the user session cookie
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, apiResponse{Success: false, Message: "Method not allowed"})
		return
	}

	cookie, err := r.Cookie("session_token")
	if err == nil && cookie.Value != "" {
		h.Store.DeleteSession(cookie.Value)
	}

	// Clear Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		MaxAge:   -1,
	})

	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		Message: "Logged out successfully.",
	})
}

// Me checks active session and returns current user profile
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, apiResponse{Success: false, Message: "Method not allowed"})
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		writeJSON(w, http.StatusUnauthorized, apiResponse{Success: false, Message: "Unauthenticated"})
		return
	}

	user, ok := h.Store.GetSession(cookie.Value)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, apiResponse{Success: false, Message: "Session expired or invalid"})
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{
		Success: true,
		User:    &user,
	})
}
