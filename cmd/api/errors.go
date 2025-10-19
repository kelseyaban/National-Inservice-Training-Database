// Filename: cmd/api/errors.go
package main

import (
	"fmt"
	"net/http"
)

// log an error message
func (a *application) logError(r *http.Request, err error) {
	method := r.Method
	uri := r.URL.RequestURI()
	a.logger.Error(err.Error(), "method", method, "uri", uri)
}

// send an error response in JSON
func (a *application) errorResponseJSON(w http.ResponseWriter, r *http.Request, status int, message any) {
	errorData := envelope{"error": message}
	err := a.writeJSON(w, status, errorData, nil)
	if err != nil {
		a.logError(r, err)
		w.WriteHeader(500)
	}
}

// send an error response if our server messes up
func (a *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	// 1st log the error message
	a.logError(r, err)
	// prepare message to response to send to the client
	message := "the server encountered a problem and could not process your request"
	a.errorResponseJSON(w, r, http.StatusInternalServerError, message)
}

// send an error response if our client messes up with a 404
func (a *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	// Only log server errors, not client errors
	// Prepare a response to send to the client
	message := "the requested resource could not be found"
	a.errorResponseJSON(w, r, http.StatusNotFound, message)
}

// send an error response if our client messes up a 405
func (a *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	// we don't log, since its a client error
	// Prepare a formatted response to send to the client
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)

	a.errorResponseJSON(w, r, http.StatusMethodNotAllowed, message)
}

// Sending an error response if client messes up with 400(bad request)
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponseJSON(w, r, http.StatusBadRequest, err.Error())
}

// How to responds to validation errors in HTTP requests
func (a *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	a.errorResponseJSON(w, r, http.StatusUnprocessableEntity, errors)
}

// Send and error response if rate limit exceeded(429 - too many requests)
func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	app.errorResponseJSON(w, r, http.StatusTooManyRequests, message)
}

// send an error response if we have an edit conflict status 409
func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	app.errorResponseJSON(w, r, http.StatusConflict, message)
}

// Return a 401 status code
func (a *application) invalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	a.errorResponseJSON(w, r, http.StatusUnauthorized, message)
}

// We set the WWW-Authenticate header to give a hint to the user as
// to what they need to provide. Don't want to leave them guessing
func (a *application) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("WWW-Authenticate", "Bearer")

	message := "invalid or missing authentication token"
	a.errorResponseJSON(w, r, http.StatusUnauthorized, message)

}

// 403 Forbidden status if bad permission
func (a *application) notPermittedResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account doesn't have the necessary permissions to access this resource"

	a.errorResponseJSON(w, r, http.StatusForbidden, message)
}

// duplicateRoleResponse returns a 409 Conflict if a user already has that role.
func (a *application) duplicateRoleResponse(w http.ResponseWriter, r *http.Request, roleName string) {
    message := fmt.Sprintf("User has already been assigned the '%s' role", roleName)
    a.errorResponseJSON(w, r, http.StatusConflict, message)
}