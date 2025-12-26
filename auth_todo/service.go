package auth_todo

import (
	"context"
	"errors"
	"fmt"
	"sort"
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
	ListTodos(ctx context.Context, userID string, limit, offset int) (todos []Todo, total int, err error)
	CompleteTodo(ctx context.Context, userID, todoID string) error
}

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTodoNotFound       = errors.New("todo not found")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrEmptyEmail         = errors.New("email cannot be empty")
	ErrEmptyPassword      = errors.New("password cannot be empty")
	ErrEmptyText          = errors.New("todo text cannot be empty")
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
	mu          sync.RWMutex
	todosByUser map[string][]Todo
	todosById   map[string]Todo
	counter     int
}

func NewTodoService() TodoService {
	return &todoService{
		todosByUser: make(map[string][]Todo),
		todosById:   make(map[string]Todo),
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
	s.todosById[todoID] = todo
	s.todosByUser[userID] = append(s.todosByUser[userID], todo)

	return todoID, nil
}

func (s *todoService) ListTodos(ctx context.Context, userID string, limit, offset int) ([]Todo, int, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	s.mu.RLock()
	userTodos, exists := s.todosByUser[userID]
	if !exists || len(userTodos) == 0 {
		s.mu.RUnlock()
		return []Todo{}, 0, nil
	}
	allTodos := make([]Todo, len(userTodos))
	copy(allTodos, userTodos)
	s.mu.RUnlock()

	sort.Slice(allTodos, func(i, j int) bool {
		return allTodos[i].CreatedAt.After(allTodos[j].CreatedAt)
	})

	total := len(allTodos)

	if offset >= total {
		return []Todo{}, total, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}

	result := allTodos[offset:end]
	return result, total, nil
}

func (s *todoService) CompleteTodo(ctx context.Context, userID, todoID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo, exists := s.todosById[todoID]
	if !exists {
		return ErrTodoNotFound
	}

	if todo.UserID != userID {
		return ErrUnauthorized
	}

	todo.Completed = true
	s.todosById[todoID] = todo
	
	userTodos := s.todosByUser[userID]
	for i, t := range userTodos {
		if t.ID == todoID {
			userTodos[i].Completed = true
			s.todosByUser[userID] = userTodos
			break
		}
	}

	return nil
}
