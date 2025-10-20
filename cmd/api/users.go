// Filename: cmd/api/users.go
package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/kelseyaban/National-Inservice-Training-Database/internal/data"
	"github.com/kelseyaban/National-Inservice-Training-Database/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
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

func (a *application) listUsersHandler(w http.ResponseWriter, r *http.Request) {
	var queryParametersData struct {
		ID               int64
		RegulationNumber string
		Username         string
		FName            string
		LName            string
		Email            string
		Gender           string
		Formation        int
		Rank             int
		Postings         int
		data.Filters
	}
	// Read the query parameters into the struct
	queryParameters := r.URL.Query()

	queryParametersData.ID = int64(a.getSingleIntegerParameter(queryParameters, "id", 0, nil))
	queryParametersData.RegulationNumber = a.getSingleQueryParameter(queryParameters, "regulation_number", "")
	queryParametersData.Username = a.getSingleQueryParameter(queryParameters, "username", "")
	queryParametersData.FName = a.getSingleQueryParameter(queryParameters, "fname", "")
	queryParametersData.LName = a.getSingleQueryParameter(queryParameters, "lname", "")
	queryParametersData.Email = a.getSingleQueryParameter(queryParameters, "email", "")
	queryParametersData.Gender = a.getSingleQueryParameter(queryParameters, "gender", "")
	queryParametersData.Formation = a.getSingleIntegerParameter(queryParameters, "formation", 0, nil)
	queryParametersData.Rank = a.getSingleIntegerParameter(queryParameters, "rank", 0, nil)
	queryParametersData.Postings = a.getSingleIntegerParameter(queryParameters, "postings", 0, nil)

	v := validator.New()
	// Add pagination and sorting
	queryParametersData.Filters.Page = a.getSingleIntegerParameter(queryParameters, "page", 1, v)
	queryParametersData.Filters.PageSize = a.getSingleIntegerParameter(queryParameters, "page_size", 20, v)
	queryParametersData.Filters.Sort = a.getSingleQueryParameter(queryParameters, "sort", "id")
	queryParametersData.Filters.SortSafeList = []string{"id", "regulation_number", "username", "fname", "lname", "email", "-formation", "rank", "postings", "-id", "-regulation_number", "-username", "-fname", "-lname", "-email", "-formation", "-rank", "-postings"}

	// Check if the filters are valid
	data.ValidateFilters(v, queryParametersData.Filters)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Get the list of users
	users, metadata, err := a.userModel.GetAll(queryParametersData.ID, queryParametersData.RegulationNumber, queryParametersData.Username, queryParametersData.FName, queryParametersData.LName, queryParametersData.Email, queryParametersData.Gender, queryParametersData.Formation, queryParametersData.Rank, queryParametersData.Postings, queryParametersData.Filters)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// Send the users as JSON response
	data := envelope{
		"users":     users,
		"@metadata": metadata,
	}

	err = a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

// Update an existing user
func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Get the id from the URL
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	user, err := app.userModel.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Read the incoming JSON
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
	}

	err = app.readJSON(w, r, &incomingData)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Update fields only if provided
	if incomingData.RegulationNumber != "" {
		user.RegulationNumber = incomingData.RegulationNumber
	}
	if incomingData.Username != "" {
		user.Username = incomingData.Username
	}
	if incomingData.FName != "" {
		user.FName = incomingData.FName
	}
	if incomingData.LName != "" {
		user.LName = incomingData.LName
	}
	if incomingData.Email != "" {
		user.Email = incomingData.Email
	}
	if incomingData.Gender != "" {
		user.Gender = incomingData.Gender
	}
	if incomingData.Formation != 0 {
		user.Formation = incomingData.Formation
	}
	if incomingData.Rank != 0 {
		user.Rank = incomingData.Rank
	}
	if incomingData.Postings != 0 {
		user.Postings = incomingData.Postings
	}

	// Validate the updated user
	v := validator.New()
	data.ValidateUser(v, *user)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Update the user in the database
	err = app.userModel.Update(user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Send a JSON response with the updated user
	data := envelope{
		"user": user,
	}

	err = app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// Delete user
func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	// Get the id from the URL
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.userModel.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// display the quote
	data := envelope{"message": "user successfully deleted"}

	err = app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) updatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	// Parse user ID from the URL
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Read the new password from JSON
	var input struct {
		NewPassword string `json:"new_password"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.NewPassword == "" {
		app.failedValidationResponse(w, r, map[string]string{"new_password": "must be provided"})
		return
	}

	// Call model method
	err = app.userModel.UpdatePassword(id, input.NewPassword)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Respond success
	env := envelope{"message": "password updated successfully"}
	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
