package models

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents an authenticated user entity
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Excluded from JSON output
	FullName     string    `json:"fullName"`
	CreatedAt    time.Time `json:"createdAt"`
}

// Session represents an active user login session
type Session struct {
	Token     string    `json:"token"`
	UserID    string    `json:"userId"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// AuthStore manages thread-safe user accounts and login sessions
type AuthStore struct {
	mu           sync.RWMutex
	users        map[string]User    // ID -> User
	usernames    map[string]string  // Lowercase Username -> User ID
	emails       map[string]string  // Lowercase Email -> User ID
	sessions     map[string]Session // Session Token -> Session
	nextUserID   int
}

// NewAuthStore initializes a thread-safe AuthStore with pre-populated demo user
func NewAuthStore() *AuthStore {
	store := &AuthStore{
		users:        make(map[string]User),
		usernames:    make(map[string]string),
		emails:       make(map[string]string),
		sessions:     make(map[string]Session),
		nextUserID:   1,
	}

	// Register a default demo user for instant testing: (demo / password123)
	_, _ = store.RegisterUser("demo", "demo@example.com", "password123", "Demo Account")

	return store
}

// RegisterUser creates a new user account with bcrypt hashed password
func (s *AuthStore) RegisterUser(username, email, password, fullName string) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	usernameClean := strings.TrimSpace(strings.ToLower(username))
	emailClean := strings.TrimSpace(strings.ToLower(email))

	if usernameClean == "" || len(usernameClean) < 3 {
		return User{}, errors.New("username must be at least 3 characters long")
	}
	if emailClean == "" || !strings.Contains(emailClean, "@") {
		return User{}, errors.New("a valid email address is required")
	}
	if len(password) < 6 {
		return User{}, errors.New("password must be at least 6 characters long")
	}

	if _, exists := s.usernames[usernameClean]; exists {
		return User{}, errors.New("username is already taken")
	}
	if _, exists := s.emails[emailClean]; exists {
		return User{}, errors.New("email is already registered")
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	userID := fmt.Sprintf("usr_%d", s.nextUserID)
	s.nextUserID++

	if fullName == "" {
		fullName = username
	}

	user := User{
		ID:           userID,
		Username:     username,
		Email:        emailClean,
		PasswordHash: string(hashedBytes),
		FullName:     fullName,
		CreatedAt:    time.Now(),
	}

	s.users[userID] = user
	s.usernames[usernameClean] = userID
	s.emails[emailClean] = userID

	return user, nil
}

// AuthenticateUser validates username/email and password using bcrypt
func (s *AuthStore) AuthenticateUser(usernameOrEmail, password string) (User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	identifier := strings.TrimSpace(strings.ToLower(usernameOrEmail))
	userID, exists := s.usernames[identifier]
	if !exists {
		userID, exists = s.emails[identifier]
	}

	if !exists {
		return User{}, errors.New("invalid credentials")
	}

	user := s.users[userID]

	// Verify bcrypt password hash
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return User{}, errors.New("invalid credentials")
	}

	return user, nil
}

// CreateSession generates a secure session token valid for 24 hours
func (s *AuthStore) CreateSession(userID string) (Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return Session{}, errors.New("failed to generate secure session token")
	}

	token := hex.EncodeToString(tokenBytes)
	session := Session{
		Token:     token,
		UserID:    userID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	s.sessions[token] = session
	return session, nil
}

// GetSession validates a session token and returns the associated user profile
func (s *AuthStore) GetSession(token string) (User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[token]
	if !exists {
		return User{}, false
	}

	if time.Now().After(session.ExpiresAt) {
		return User{}, false
	}

	user, exists := s.users[session.UserID]
	return user, exists
}

// DeleteSession invalidates a session token (Logout)
func (s *AuthStore) DeleteSession(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, token)
}

// GetUserByID fetches a user by unique ID
func (s *AuthStore) GetUserByID(userID string) (User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, exists := s.users[userID]
	return user, exists
}
