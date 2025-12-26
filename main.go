package main

import (
	"net/http"
	"os"
	"time"

	"todo-microservice/auth_todo"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	kitlog "github.com/go-kit/log"
	"github.com/gorilla/mux"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	logger := kitlog.NewLogfmtLogger(os.Stderr)
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)
	logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)

	fieldKeys := []string{"method"}

	authRequestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "auth_todo",
		Subsystem: "auth_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)

	authRequestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "auth_todo",
		Subsystem: "auth_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	todoRequestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "auth_todo",
		Subsystem: "todo_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)

	todoRequestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "auth_todo",
		Subsystem: "todo_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	var authSvc auth_todo.AuthService
	authSvc = auth_todo.NewAuthService()
	authSvc = auth_todo.NewLoggingAuthMiddleware(logger, authSvc)
	authSvc = auth_todo.NewInstrumentingAuthMiddleware(authRequestCount, authRequestLatency, authSvc)

	var todoSvc auth_todo.TodoService
	todoSvc = auth_todo.NewTodoService()
	todoSvc = auth_todo.NewCachedTodoService(30*time.Second, todoSvc)
	todoSvc = auth_todo.NewLoggingTodoMiddleware(logger, todoSvc)
	todoSvc = auth_todo.NewInstrumentingTodoMiddleware(todoRequestCount, todoRequestLatency, todoSvc)

	endpoints := auth_todo.MakeEndpoints(authSvc, todoSvc)

	r := mux.NewRouter()

	r.Handle("/signup", auth_todo.MakeSignupHandler(endpoints)).Methods("POST")
	r.Handle("/login", auth_todo.MakeLoginHandler(endpoints)).Methods("POST")
	r.Handle("/validate", auth_todo.MakeValidateTokenHandler(endpoints)).Methods("POST", "GET")
	r.Handle("/todos", auth_todo.MakeCreateTodoHandler(endpoints)).Methods("POST")
	r.Handle("/todos", auth_todo.MakeListTodosHandler(endpoints)).Methods("GET")
	r.Handle("/todos/{id}/complete", auth_todo.MakeCompleteTodoHandler(endpoints)).Methods("POST")

	r.Handle("/metrics", promhttp.Handler())

	logger.Log("msg", "HTTP server started", "addr", ":8080")
	http.ListenAndServe(":8080", r)
}
