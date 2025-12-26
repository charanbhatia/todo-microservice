package auth_todo

import (
	"context"
	"testing"
)

func TestAuthService(t *testing.T) {
	svc := NewAuthService()
	ctx := context.Background()

	// Test Signup
	userID, err := svc.Signup(ctx, "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Signup failed: %v", err)
	}
	if userID == "" {
		t.Fatal("Expected non-empty userID")
	}

	// Test duplicate signup
	_, err = svc.Signup(ctx, "test@example.com", "password123")
	if err != ErrUserExists {
		t.Fatalf("Expected ErrUserExists, got: %v", err)
	}

	// Test Login
	token, err := svc.Login(ctx, "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if token == "" {
		t.Fatal("Expected non-empty token")
	}

	// Test invalid login
	_, err = svc.Login(ctx, "test@example.com", "wrongpassword")
	if err != ErrInvalidCredentials {
		t.Fatalf("Expected ErrInvalidCredentials, got: %v", err)
	}

	// Test ValidateToken
	validUserID, err := svc.ValidateToken(ctx, token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}
	if validUserID != userID {
		t.Fatalf("Expected userID %s, got %s", userID, validUserID)
	}

	// Test invalid token
	_, err = svc.ValidateToken(ctx, "invalid_token")
	if err != ErrInvalidToken {
		t.Fatalf("Expected ErrInvalidToken, got: %v", err)
	}
}

func TestTodoService(t *testing.T) {
	svc := NewTodoService()
	ctx := context.Background()
	userID := "user_1"

	// Test CreateTodo
	todoID, err := svc.CreateTodo(ctx, userID, "Buy groceries")
	if err != nil {
		t.Fatalf("CreateTodo failed: %v", err)
	}
	if todoID == "" {
		t.Fatal("Expected non-empty todoID")
	}

	// Test ListTodos
	todos, total, err := svc.ListTodos(ctx, userID, 50, 0)
	if err != nil {
		t.Fatalf("ListTodos failed: %v", err)
	}
	if len(todos) != 1 {
		t.Fatalf("Expected 1 todo, got %d", len(todos))
	}
	if total != 1 {
		t.Fatalf("Expected total 1, got %d", total)
	}
	if todos[0].Text != "Buy groceries" {
		t.Fatalf("Expected 'Buy groceries', got '%s'", todos[0].Text)
	}

	// Test CompleteTodo
	err = svc.CompleteTodo(ctx, userID, todoID)
	if err != nil {
		t.Fatalf("CompleteTodo failed: %v", err)
	}

	// Verify completion
	todos, _, _ = svc.ListTodos(ctx, userID, 50, 0)
	if !todos[0].Completed {
		t.Fatal("Expected todo to be completed")
	}

	// Test unauthorized complete
	err = svc.CompleteTodo(ctx, "different_user", todoID)
	if err != ErrUnauthorized {
		t.Fatalf("Expected ErrUnauthorized, got: %v", err)
	}

	// Test complete non-existent todo
	err = svc.CompleteTodo(ctx, userID, "nonexistent")
	if err != ErrTodoNotFound {
		t.Fatalf("Expected ErrTodoNotFound, got: %v", err)
	}
}
