// Filename: cmd/api/tokens.go
package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/kelseyaban/National-Inservice-Training-Database/internal/data"
	"github.com/kelseyaban/National-Inservice-Training-Database/internal/validator"
)

func (a *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Get the body from the request and store in a temporary struct
	// The client will give us their email and password. We will will give them
	// a Bearer token
	var incomingData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	// Validate the email and password provided by the client.
	v := validator.New()

	data.ValidateEmail(v, incomingData.Email)
	data.ValidatePasswordPlaintext(v, incomingData.Password)

	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Is there an associated user for the provided email?
	user, err := a.userModel.GetByEmail(incomingData.Email)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.invalidCredentialsResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}
	// The user is found. Does their password match?
	match, err := user.Password.Matches(incomingData.Password)

	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// Wrong password
	// We will define invalidCredentialsResponse() later
	if !match {
		a.invalidCredentialsResponse(w, r)
		return
	}
	token, err := a.tokenModel.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	data := envelope{
		"authentication_token": token,
	}

	// Return the bearer token
	err = a.writeJSON(w, http.StatusCreated, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
