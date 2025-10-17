// Filename: cmd/api/users.go
package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/kelseyaban/National-Inservice-Training-Database/internal/data"
	"github.com/kelseyaban/National-Inservice-Training-Database/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter,
	r *http.Request) {
	// Get the passed in data from the request body and store in a temporary struct
	var incomingData struct {
		RegulationNumber string `json:"regulation_number"`
		Username         string `json:"username"`
		FName            string `json:"fname"`
		LName            string `json:"lname"`
		Email            string `json:"email"`
		Gender           string `json:"gender"`
		Formation        int    `json:"formation"`
		Rank             int    `json:"rank"`
		Postings         int    `json:"postings"`
		Password         string `json:"password"`
	}

	err := app.readJSON(w, r, &incomingData)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// we will add the password later after we have hashed it
	user := &data.User{
		RegulationNumber: incomingData.RegulationNumber,
		Username:         incomingData.Username,
		FName:            incomingData.FName,
		LName:            incomingData.LName,
		Email:            incomingData.Email,
		Gender:           incomingData.Gender,
		Formation:        incomingData.Formation,
		Rank:             incomingData.Rank,
		Postings:         incomingData.Postings,
		Activated:        false,
	}

	// hash the password and store it along with the cleartext version
	err = user.Password.Set(incomingData.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Perform validation for the User
	v := validator.New()

	data.ValidateUser(v, *user)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.userModel.Insert(user) // we will add userModel to main() later
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Add the read permission for new users
	err = app.permissionModel.AddForUser(user.ID, "session:read")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Generate a new activation token which expires in 3 days
	token, err := app.tokenModel.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	data := envelope{
		"user": user,
	}

	// Send the email as a  go routine
	app.background(func() {
		data := map[string]any{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}

		err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			app.logger.Error(err.Error())
		}
	})

	// Status code 201 resource created
	err = app.writeJSON(w, http.StatusCreated, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (a *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Get the body from the request and store in temporary struct
	var incomingData struct {
		TokenPlaintext string `json:"token"`
	}
	err := a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	// Validate the data
	v := validator.New()
	data.ValidateTokenPlaintext(v, incomingData.TokenPlaintext)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Let's check if the token provided belongs to the user
	user, err := a.userModel.GetForToken(data.ScopeActivation,
		incomingData.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			a.failedValidationResponse(w, r, v.Errors)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	err = a.userModel.Activate(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			a.editConflictResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	// Re-fetch the full user from the database so all fields are populated
	user, err = a.userModel.GetByEmail(user.Email)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// User has been activated so let's delete the activation token to
	// prevent reuse.
	err = a.tokenModel.DeleteAllForUser(data.ScopeActivation, user.ID)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// Send a response
	data := envelope{
		"user": user,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
