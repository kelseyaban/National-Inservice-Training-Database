// Filename: cmd/api/routes.go

package main

import (
	"expvar"
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

	router.HandlerFunc(http.MethodPatch, "/v1/users/update/:id", app.requirePermission("users:write", app.updateUserHandler))
	router.HandlerFunc(http.MethodGet, "/v1/users/details", app.requirePermission("users:read", app.listUsersHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/users/delete/:id", app.requirePermission("users:write", app.deleteUserHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/users/update-password/:id", app.requirePermission("users:write", app.updatePasswordHandler))

	// Roles
	router.HandlerFunc(http.MethodPost, "/v1/roles", app.requirePermission("role:write", app.createRoleHandler))
	router.HandlerFunc(http.MethodGet, "/v1/roles/:id", app.requirePermission("role:read", app.displayRoleHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/roles/:id", app.requirePermission("role:write", app.updateRoleHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/roles/:id", app.requirePermission("role:write", app.deleteRoleHandler))
	router.HandlerFunc(http.MethodGet, "/v1/roles", app.requirePermission("role:read", app.listRoleHandler))

	//User Roles
	router.HandlerFunc(http.MethodPost, "/v1/users/assign-role", app.requirePermission("role:write", app.assignRoleHandler))
	router.HandlerFunc(http.MethodGet, "/v1/users/user_roles/:id", app.requirePermission("role:read", app.getUserRolesHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/users/update-role/:id", app.requirePermission("role:write", app.updateUserRoleHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/users/delete-role/:id", app.requirePermission("role:write", app.deleteUserRoleHandler))
	router.HandlerFunc(http.MethodGet, "/v1/users/user_roles", app.requirePermission("role:read", app.listUsersWithRolesHandler))

	//Facilitator Rating
	router.HandlerFunc(http.MethodPost, "/v1/facilitator-rating", app.requirePermission("facilitator_rating:write", app.addFacilitatorRating))
	router.HandlerFunc(http.MethodGet, "/v1/facilitator-rating/:id", app.requirePermission("facilitator_rating:read", app.displayFacilitatorRatingHandler))
	router.HandlerFunc(http.MethodGet, "/v1/facilitator-rating", app.requirePermission("facilitator_rating:read", app.listFacilitatorRatingHandler))

	// Courses
	router.HandlerFunc(http.MethodPost, "/v1/courses", app.requirePermission("course:write", app.createCourseHandler))
	router.HandlerFunc(http.MethodGet, "/v1/courses/:id", app.requirePermission("course:read", app.displayCourseHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/courses/:id", app.requirePermission("course:write", app.updateCourseHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/courses/:id", app.requirePermission("course:write", app.deleteCourseHandler))
	router.HandlerFunc(http.MethodGet, "/v1/courses", app.requirePermission("course:read", app.listCoursesHandler))

	// Course Postings
	router.HandlerFunc(http.MethodPost, "/v1/course/posting", app.requirePermission("course_posting:write", app.createCoursePostingHandler))
	router.HandlerFunc(http.MethodGet, "/v1/course/posting/:id", app.requirePermission("course_posting:read", app.displayCoursePostingHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/course/posting/:id", app.requirePermission("course_posting:write", app.updateCoursePostingHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/course/posting/:id", app.requirePermission("course_posting:write", app.deleteCoursePostingHandler))
	router.HandlerFunc(http.MethodGet, "/v1/course/posting", app.requirePermission("course_posting:read", app.listCoursePostingsHandler))

	// Sessions
	router.HandlerFunc(http.MethodPost, "/v1/session", app.requirePermission("session:write", app.createSessionHandler))
	router.HandlerFunc(http.MethodGet, "/v1/session/:id", app.requirePermission("session:read", app.displaySessionHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/session/:id", app.requirePermission("session:write", app.updateSessionHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/session/:id", app.requirePermission("session:write", app.deleteSessionHandler))
	router.HandlerFunc(http.MethodGet, "/v1/session", app.requirePermission("session:read", app.listSessionHandler))

	//User Session
	router.HandlerFunc(http.MethodPost, "/v1/user_session", app.requirePermission("user_session:write", app.createUserSessionHandler))
	router.HandlerFunc(http.MethodGet, "/v1/user_session/:id", app.requirePermission("user_session:read", app.getUserSessionHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/user_session/:id", app.requirePermission("user_session:write", app.updateUserSessionHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/user_session/:id", app.requirePermission("user_session:write", app.deleteUserSessionHandler))
	router.HandlerFunc(http.MethodGet, "/v1/user_session", app.requirePermission("user_session:read", app.listUserSessionHandler))

	// Attendance
	router.HandlerFunc(http.MethodPost, "/v1/attendance", app.requirePermission("attendance:write", app.createAttendanceHandler))
	router.HandlerFunc(http.MethodGet, "/v1/attendance/:id", app.requirePermission("user_session:read", app.displayIndividualAttendanceHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/attendance/:id", app.requirePermission("user_session:write", app.updateAttendanceHandler))

	router.Handler(http.MethodGet, "/v1/observability/quotes/metrics", expvar.Handler())

	// return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(router))))
}
