package auth_todo

import (
	"context"
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/log"
)

type loggingAuthMiddleware struct {
	logger log.Logger
	next   AuthService
}

func NewLoggingAuthMiddleware(logger log.Logger, svc AuthService) AuthService {
	return &loggingAuthMiddleware{
		logger: logger,
		next:   svc,
	}
}

func (mw *loggingAuthMiddleware) Signup(ctx context.Context, email, password string) (userID string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "Signup",
			"email", email,
			"user_id", userID,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.Signup(ctx, email, password)
}

func (mw *loggingAuthMiddleware) Login(ctx context.Context, email, password string) (token string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "Login",
			"email", email,
			"token_generated", token != "",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.Login(ctx, email, password)
}

func (mw *loggingAuthMiddleware) ValidateToken(ctx context.Context, token string) (userID string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "ValidateToken",
			"user_id", userID,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.ValidateToken(ctx, token)
}

type instrumentingAuthMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           AuthService
}

func NewInstrumentingAuthMiddleware(counter metrics.Counter, latency metrics.Histogram, svc AuthService) AuthService {
	return &instrumentingAuthMiddleware{
		requestCount:   counter,
		requestLatency: latency,
		next:           svc,
	}
}

func (mw *instrumentingAuthMiddleware) Signup(ctx context.Context, email, password string) (string, error) {
	defer func(begin time.Time) {
		mw.requestCount.With("method", "Signup").Add(1)
		mw.requestLatency.With("method", "Signup").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mw.next.Signup(ctx, email, password)
}

func (mw *instrumentingAuthMiddleware) Login(ctx context.Context, email, password string) (string, error) {
	defer func(begin time.Time) {
		mw.requestCount.With("method", "Login").Add(1)
		mw.requestLatency.With("method", "Login").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mw.next.Login(ctx, email, password)
}

func (mw *instrumentingAuthMiddleware) ValidateToken(ctx context.Context, token string) (string, error) {
	defer func(begin time.Time) {
		mw.requestCount.With("method", "ValidateToken").Add(1)
		mw.requestLatency.With("method", "ValidateToken").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mw.next.ValidateToken(ctx, token)
}

type loggingTodoMiddleware struct {
	logger log.Logger
	next   TodoService
}

func NewLoggingTodoMiddleware(logger log.Logger, svc TodoService) TodoService {
	return &loggingTodoMiddleware{
		logger: logger,
		next:   svc,
	}
}

func (mw *loggingTodoMiddleware) CreateTodo(ctx context.Context, userID, text string) (todoID string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "CreateTodo",
			"user_id", userID,
			"todo_id", todoID,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.CreateTodo(ctx, userID, text)
}

func (mw *loggingTodoMiddleware) ListTodos(ctx context.Context, userID string, limit, offset int) (todos []Todo, total int, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "ListTodos",
			"user_id", userID,
			"limit", limit,
			"offset", offset,
			"count", len(todos),
			"total", total,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.ListTodos(ctx, userID, limit, offset)
}

func (mw *loggingTodoMiddleware) CompleteTodo(ctx context.Context, userID, todoID string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "CompleteTodo",
			"user_id", userID,
			"todo_id", todoID,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.CompleteTodo(ctx, userID, todoID)
}

type instrumentingTodoMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           TodoService
}

func NewInstrumentingTodoMiddleware(counter metrics.Counter, latency metrics.Histogram, svc TodoService) TodoService {
	return &instrumentingTodoMiddleware{
		requestCount:   counter,
		requestLatency: latency,
		next:           svc,
	}
}

func (mw *instrumentingTodoMiddleware) CreateTodo(ctx context.Context, userID, text string) (string, error) {
	defer func(begin time.Time) {
		mw.requestCount.With("method", "CreateTodo").Add(1)
		mw.requestLatency.With("method", "CreateTodo").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mw.next.CreateTodo(ctx, userID, text)
}

func (mw *instrumentingTodoMiddleware) ListTodos(ctx context.Context, userID string, limit, offset int) ([]Todo, int, error) {
	defer func(begin time.Time) {
		mw.requestCount.With("method", "ListTodos").Add(1)
		mw.requestLatency.With("method", "ListTodos").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mw.next.ListTodos(ctx, userID, limit, offset)
}

func (mw *instrumentingTodoMiddleware) CompleteTodo(ctx context.Context, userID, todoID string) error {
	defer func(begin time.Time) {
		mw.requestCount.With("method", "CompleteTodo").Add(1)
		mw.requestLatency.With("method", "CompleteTodo").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return mw.next.CompleteTodo(ctx, userID, todoID)
}
