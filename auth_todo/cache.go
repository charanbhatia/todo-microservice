package auth_todo

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type cacheEntry struct {
	todos     []Todo
	total     int
	timestamp time.Time
}

type cachedTodoService struct {
	mu    sync.RWMutex
	cache map[string]cacheEntry
	ttl   time.Duration
	next  TodoService
}

func NewCachedTodoService(ttl time.Duration, svc TodoService) TodoService {
	cached := &cachedTodoService{
		cache: make(map[string]cacheEntry),
		ttl:   ttl,
		next:  svc,
	}

	go cached.cleanupExpired()

	return cached
}

func (s *cachedTodoService) cleanupExpired() {
	ticker := time.NewTicker(s.ttl)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for key, entry := range s.cache {
			if now.Sub(entry.timestamp) > s.ttl {
				delete(s.cache, key)
			}
		}
		s.mu.Unlock()
	}
}

func (s *cachedTodoService) CreateTodo(ctx context.Context, userID, text string) (string, error) {
	todoID, err := s.next.CreateTodo(ctx, userID, text)
	if err != nil {
		return "", err
	}

	s.mu.Lock()
	for key := range s.cache {
		if len(key) > len(userID) && key[:len(userID)] == userID {
			delete(s.cache, key)
		}
	}
	s.mu.Unlock()

	return todoID, nil
}

func (s *cachedTodoService) ListTodos(ctx context.Context, userID string, limit, offset int) ([]Todo, int, error) {
	cacheKey := fmt.Sprintf("%s:%d:%d", userID, limit, offset)

	s.mu.RLock()
	entry, exists := s.cache[cacheKey]
	s.mu.RUnlock()

	if exists && time.Since(entry.timestamp) < s.ttl {
		return entry.todos, entry.total, nil
	}

	todos, total, err := s.next.ListTodos(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	s.mu.Lock()
	s.cache[cacheKey] = cacheEntry{
		todos:     todos,
		total:     total,
		timestamp: time.Now(),
	}
	s.mu.Unlock()

	return todos, total, nil
}

func (s *cachedTodoService) CompleteTodo(ctx context.Context, userID, todoID string) error {
	err := s.next.CompleteTodo(ctx, userID, todoID)
	if err != nil {
		return err
	}

	s.mu.Lock()
	for key := range s.cache {
		if len(key) > len(userID) && key[:len(userID)] == userID {
			delete(s.cache, key)
		}
	}
	s.mu.Unlock()

	return nil
}
