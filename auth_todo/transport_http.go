package auth_todo

import (
	"context"
	"encoding/json"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func decodeSignupRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req signupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeValidateTokenRequest(_ context.Context, r *http.Request) (interface{}, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		var req validateTokenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return nil, err
		}
		return req, nil
	}
	return validateTokenRequest{Token: token}, nil
}

func decodeCreateTodoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req createTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeListTodosRequest(_ context.Context, r *http.Request) (interface{}, error) {
	userID := r.URL.Query().Get("user_id")
	return listTodosRequest{UserID: userID}, nil
}

func decodeCompleteTodoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	todoID := vars["id"]
	
	var req struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	
	return completeTodoRequest{
		UserID: req.UserID,
		TodoID: todoID,
	}, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func MakeSignupHandler(endpoints Endpoints) http.Handler {
	return httptransport.NewServer(
		endpoints.SignupEndpoint,
		decodeSignupRequest,
		encodeResponse,
	)
}

func MakeLoginHandler(endpoints Endpoints) http.Handler {
	return httptransport.NewServer(
		endpoints.LoginEndpoint,
		decodeLoginRequest,
		encodeResponse,
	)
}

func MakeValidateTokenHandler(endpoints Endpoints) http.Handler {
	return httptransport.NewServer(
		endpoints.ValidateTokenEndpoint,
		decodeValidateTokenRequest,
		encodeResponse,
	)
}

func MakeCreateTodoHandler(endpoints Endpoints) http.Handler {
	return httptransport.NewServer(
		endpoints.CreateTodoEndpoint,
		decodeCreateTodoRequest,
		encodeResponse,
	)
}

func MakeListTodosHandler(endpoints Endpoints) http.Handler {
	return httptransport.NewServer(
		endpoints.ListTodosEndpoint,
		decodeListTodosRequest,
		encodeResponse,
	)
}

func MakeCompleteTodoHandler(endpoints Endpoints) http.Handler {
	return httptransport.NewServer(
		endpoints.CompleteTodoEndpoint,
		decodeCompleteTodoRequest,
		encodeResponse,
	)
}
