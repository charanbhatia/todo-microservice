package auth_todo

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Todo struct {
	ID        string
	UserID    string
	Text      string
	Completed bool
	CreatedAt time.Time
}

type AuthService interface {
	Signup(ctx context.Context, email, password string) (userID string, err error)
	Login(ctx context.Context, email, password string) (token string, err error)
	ValidateToken(ctx context.Context, token string) (userID string, err error)
}

type TodoService interface {
	CreateTodo(ctx context.Context, userID, text string) (todoID string, err error)
	ListTodos(ctx context.Context, userID string) ([]Todo, error)
	CompleteTodo(ctx context.Context, userID, todoID string) error
}

var (
	ErrUserExists       = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken     = errors.New("invalid token")
	ErrTodoNotFound     = errors.New("todo not found")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrEmptyEmail       = errors.New("email cannot be empty")
	ErrEmptyPassword    = errors.New("password cannot be empty")
	ErrEmptyText        = errors.New("todo text cannot be empty")
)

type user struct {
	ID       string
	Email    string
	Password string
}

type authService struct {
	mu      sync.RWMutex
	users   map[string]user
	tokens  map[string]string
	counter int
}

func NewAuthService() AuthService {
	return &authService{
		users:  make(map[string]user),
		tokens: make(map[string]string),
	}
}

func (s *authService) Signup(ctx context.Context, email, password string) (string, error) {
	if email == "" {
		return "", ErrEmptyEmail
	}
	if password == "" {
		return "", ErrEmptyPassword
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[email]; exists {
		return "", ErrUserExists
	}

	s.counter++
	userID := fmt.Sprintf("user_%d", s.counter)
	s.users[email] = user{
		ID:       userID,
		Email:    email,
		Password: password,
	}

	return userID, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	if email == "" {
		return "", ErrEmptyEmail
	}
	if password == "" {
		return "", ErrEmptyPassword
	}

	s.mu.RLock()
	u, exists := s.users[email]
	s.mu.RUnlock()

	if !exists || u.Password != password {
		return "", ErrInvalidCredentials
	}

	token := fmt.Sprintf("token_%s_%d", u.ID, time.Now().Unix())
	
	s.mu.Lock()
	s.tokens[token] = u.ID
	s.mu.Unlock()

	return token, nil
}

func (s *authService) ValidateToken(ctx context.Context, token string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userID, exists := s.tokens[token]
	if !exists {
		return "", ErrInvalidToken
	}

	return userID, nil
}

type todoService struct {
	mu      sync.RWMutex
	todos   map[string]Todo
	counter int
}

func NewTodoService() TodoService {
	return &todoService{
		todos: make(map[string]Todo),
	}
}

func (s *todoService) CreateTodo(ctx context.Context, userID, text string) (string, error) {
	if text == "" {
		return "", ErrEmptyText
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.counter++
	todoID := fmt.Sprintf("todo_%d", s.counter)
	todo := Todo{
		ID:        todoID,
		UserID:    userID,
		Text:      text,
		Completed: false,
		CreatedAt: time.Now(),
	}
	s.todos[todoID] = todo

	return todoID, nil
}

func (s *todoService) ListTodos(ctx context.Context, userID string) ([]Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []Todo
	for _, todo := range s.todos {
		if todo.UserID == userID {
			result = append(result, todo)
		}
	}

	return result, nil
}

func (s *todoService) CompleteTodo(ctx context.Context, userID, todoID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo, exists := s.todos[todoID]
	if !exists {
		return ErrTodoNotFound
	}

	if todo.UserID != userID {
		return ErrUnauthorized
	}

	todo.Completed = true
	s.todos[todoID] = todo

	return nil
}
