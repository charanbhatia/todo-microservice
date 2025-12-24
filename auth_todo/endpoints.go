package auth_todo

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type signupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type signupResponse struct {
	UserID string `json:"user_id,omitempty"`
	Err    string `json:"error,omitempty"`
}

func makeSignupEndpoint(svc AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(signupRequest)
		userID, err := svc.Signup(ctx, req.Email, req.Password)
		if err != nil {
			return signupResponse{Err: err.Error()}, nil
		}
		return signupResponse{UserID: userID}, nil
	}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token,omitempty"`
	Err   string `json:"error,omitempty"`
}

func makeLoginEndpoint(svc AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(loginRequest)
		token, err := svc.Login(ctx, req.Email, req.Password)
		if err != nil {
			return loginResponse{Err: err.Error()}, nil
		}
		return loginResponse{Token: token}, nil
	}
}

type validateTokenRequest struct {
	Token string `json:"token"`
}

type validateTokenResponse struct {
	UserID string `json:"user_id,omitempty"`
	Err    string `json:"error,omitempty"`
}

func makeValidateTokenEndpoint(svc AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(validateTokenRequest)
		userID, err := svc.ValidateToken(ctx, req.Token)
		if err != nil {
			return validateTokenResponse{Err: err.Error()}, nil
		}
		return validateTokenResponse{UserID: userID}, nil
	}
}

type createTodoRequest struct {
	UserID string `json:"user_id"`
	Text   string `json:"text"`
}

type createTodoResponse struct {
	TodoID string `json:"todo_id,omitempty"`
	Err    string `json:"error,omitempty"`
}

func makeCreateTodoEndpoint(svc TodoService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createTodoRequest)
		todoID, err := svc.CreateTodo(ctx, req.UserID, req.Text)
		if err != nil {
			return createTodoResponse{Err: err.Error()}, nil
		}
		return createTodoResponse{TodoID: todoID}, nil
	}
}

type listTodosRequest struct {
	UserID string `json:"user_id"`
}

type listTodosResponse struct {
	Todos []Todo `json:"todos,omitempty"`
	Err   string `json:"error,omitempty"`
}

func makeListTodosEndpoint(svc TodoService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listTodosRequest)
		todos, err := svc.ListTodos(ctx, req.UserID)
		if err != nil {
			return listTodosResponse{Err: err.Error()}, nil
		}
		return listTodosResponse{Todos: todos}, nil
	}
}

type completeTodoRequest struct {
	UserID string `json:"user_id"`
	TodoID string `json:"todo_id"`
}

type completeTodoResponse struct {
	Err string `json:"error,omitempty"`
}

func makeCompleteTodoEndpoint(svc TodoService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(completeTodoRequest)
		err := svc.CompleteTodo(ctx, req.UserID, req.TodoID)
		if err != nil {
			return completeTodoResponse{Err: err.Error()}, nil
		}
		return completeTodoResponse{}, nil
	}
}

type Endpoints struct {
	SignupEndpoint        endpoint.Endpoint
	LoginEndpoint         endpoint.Endpoint
	ValidateTokenEndpoint endpoint.Endpoint
	CreateTodoEndpoint    endpoint.Endpoint
	ListTodosEndpoint     endpoint.Endpoint
	CompleteTodoEndpoint  endpoint.Endpoint
}

func MakeEndpoints(authSvc AuthService, todoSvc TodoService) Endpoints {
	return Endpoints{
		SignupEndpoint:        makeSignupEndpoint(authSvc),
		LoginEndpoint:         makeLoginEndpoint(authSvc),
		ValidateTokenEndpoint: makeValidateTokenEndpoint(authSvc),
		CreateTodoEndpoint:    makeCreateTodoEndpoint(todoSvc),
		ListTodosEndpoint:     makeListTodosEndpoint(todoSvc),
		CompleteTodoEndpoint:  makeCompleteTodoEndpoint(todoSvc),
	}
}
