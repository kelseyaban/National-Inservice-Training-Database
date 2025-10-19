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
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	// Users
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	// Roles
	router.HandlerFunc(http.MethodPost, "/v1/roles", app.createRoleHandler)
	router.HandlerFunc(http.MethodGet, "/v1/roles/:id", app.displayRoleHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/roles/:id", app.updateRoleHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/roles/:id", app.deleteRoleHandler)
	router.HandlerFunc(http.MethodGet, "/v1/roles", app.listRoleHandler)

	//User Roles
	router.HandlerFunc(http.MethodPost, "/v1/users/assign-role", app.assignRoleHandler)
	router.HandlerFunc(http.MethodGet, "/v1/users/user_roles/:id", app.getUserRolesHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/users/update-role/:id", app.updateUserRoleHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/users/delete-role/:id", app.deleteUserRoleHandler)
	router.HandlerFunc(http.MethodGet, "/v1/users/user_roles", app.listUsersWithRolesHandler)

	//Facilitator Rating
	router.HandlerFunc(http.MethodPost, "/v1/facilitator-rating", app.addFacilitatorRating)
	router.HandlerFunc(http.MethodGet, "/v1/facilitator-rating/:id", app.displayFacilitatorRatingHandler)
	router.HandlerFunc(http.MethodGet, "/v1/facilitator-rating", app.listFacilitatorRatingHandler)


	// Courses
	// Permission structure: router.HandlerFunc(http.MethodGet, "/v1/quotes/:id", app.requirePermission("quotes:read", app.displayQuoteHandler))
	router.HandlerFunc(http.MethodPost, "/v1/courses", app.createCourseHandler)
	router.HandlerFunc(http.MethodGet, "/v1/courses/:id", app.displayCourseHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/courses/:id", app.updateCourseHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/courses/:id", app.deleteCourseHandler)
	router.HandlerFunc(http.MethodGet, "/v1/courses", app.listCoursesHandler)



	// return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(router))))
}
