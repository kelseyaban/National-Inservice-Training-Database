// Filename: cmd/api/context.go
package main

import (
	"context"
	"net/http"

	"github.com/kelseyaban/National-Inservice-Training-Database/internal/data"
)

// Define a custom type for the context key to avoid potential collisions
type contextKey string

const userContextKey = contextKey("user")

// Update the request context with the user information
// We return the request context with user-info added
func (a *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	// WithValue() expects the original context along with the new
	// key:value pair you want to update it with
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

// Retrieve the user info when we expect it to be present (registered users)
// We can panic here because it means something went unexpectedly wrong.
// (*data.User) converts the value from a generic type (any) to a User type
func (a *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}
