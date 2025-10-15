// Filename: cmd/api/routes.go

package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// routes specifies our routes
func (app *application) routes() http.Handler {
	// setup a new routes
	router := httprouter.New()

	// handle 404
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	// handle 405
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// setup routes
	// router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	// router.HandlerFunc(http.MethodPost, "/v1/quotes", app.requirePermission("quotes:write", app.createQuoteHandler))
	// router.HandlerFunc(http.MethodGet, "/v1/quotes/:id", app.requirePermission("quotes:read", app.displayQuoteHandler))
	// router.HandlerFunc(http.MethodPatch, "/v1/quotes/:id", app.requirePermission("quotes:write", app.updateQuoteHandler))
	// router.HandlerFunc(http.MethodDelete, "/v1/quotes/:id", app.requirePermission("quotes:write", app.deleteQuoteHandler))
	// router.HandlerFunc(http.MethodGet, "/v1/quotes", app.requirePermission("quotes:read", app.listQuotesHandler))
	// router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	// router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	// router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	// router.Handler(http.MethodGet, "/v1/observability/quotes/metrics", expvar.Handler())

	// return app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router))))
	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}
