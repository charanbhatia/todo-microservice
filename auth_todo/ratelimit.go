package auth_todo

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
	"golang.org/x/time/rate"
)

var ErrRateLimitExceeded = errors.New("rate limit exceeded")

func NewRateLimitMiddleware(limit rate.Limit, burst int) endpoint.Middleware {
	limiter := rate.NewLimiter(limit, burst)
	
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			if !limiter.Allow() {
				return nil, ErrRateLimitExceeded
			}
			return next(ctx, request)
		}
	}
}
