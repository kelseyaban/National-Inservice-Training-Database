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

	router.HandlerFunc(http.MethodPatch, "/v1/users/update/:id", app.requirePermission("users:write", app.requireActivatedUser(app.updateUserHandler)),)
	router.HandlerFunc(http.MethodGet, "/v1/users/details", app.requirePermission("users:read", app.requireActivatedUser(app.listUsersHandler)),)
	router.HandlerFunc(http.MethodDelete, "/v1/users/delete/:id", app.requirePermission("users:write", app.requireActivatedUser(app.deleteUserHandler)),)
	router.HandlerFunc(http.MethodPatch, "/v1/users/update-password/:id", app.requirePermission("users:write", app.requireActivatedUser(app.updatePasswordHandler)),)

	// Roles
	router.HandlerFunc(http.MethodPost, "/v1/roles", app.requirePermission("role:write", app.requireActivatedUser(app.createRoleHandler)),)
	router.HandlerFunc(http.MethodGet, "/v1/roles/:id", app.requirePermission("role:read", app.requireActivatedUser(app.displayRoleHandler)),)
	router.HandlerFunc(http.MethodPatch, "/v1/roles/:id", app.requirePermission("role:write", app.requireActivatedUser(app.updateRoleHandler)),)
	router.HandlerFunc(http.MethodDelete, "/v1/roles/:id", app.requirePermission("role:write", app.requireActivatedUser(app.deleteRoleHandler)),)
	router.HandlerFunc(http.MethodGet, "/v1/roles", app.requirePermission("role:read", app.requireActivatedUser(app.listRoleHandler)),)

	//User Roles
	router.HandlerFunc(http.MethodPost, "/v1/users/assign-role", app.requirePermission("role:write", app.requireActivatedUser(app.assignRoleHandler)),)
	router.HandlerFunc(http.MethodGet, "/v1/users/user_roles/:id", app.requirePermission("role:read", app.requireActivatedUser(app.getUserRolesHandler)),)
	router.HandlerFunc(http.MethodPatch, "/v1/users/update-role/:id", app.requirePermission("role:write", app.requireActivatedUser(app.updateUserRoleHandler)),)
	router.HandlerFunc(http.MethodDelete, "/v1/users/delete-role/:id", app.requirePermission("role:write", app.requireActivatedUser(app.deleteUserRoleHandler)),)
	router.HandlerFunc(http.MethodGet, "/v1/users/user_roles", app.requirePermission("role:read", app.requireActivatedUser(app.listUsersWithRolesHandler)),)

	//Facilitator Rating
	router.HandlerFunc(http.MethodPost, "/v1/facilitator-rating", app.requirePermission("facilitator_rating:write", app.requireActivatedUser(app.addFacilitatorRating)),)
	router.HandlerFunc(http.MethodGet, "/v1/facilitator-rating/:id", app.requirePermission("facilitator_rating:read", app.requireActivatedUser(app.displayFacilitatorRatingHandler)),)
	router.HandlerFunc(http.MethodGet, "/v1/facilitator-rating", app.requirePermission("facilitator_rating:read", app.requireActivatedUser(app.listFacilitatorRatingHandler)),)

	// Courses
	router.HandlerFunc(http.MethodPost, "/v1/courses", app.requirePermission("course:write", app.requireActivatedUser(app.createCourseHandler)),)
	router.HandlerFunc(http.MethodGet, "/v1/courses/:id", app.requirePermission("course:read", app.requireActivatedUser(app.displayCourseHandler)),)
	router.HandlerFunc(http.MethodPatch, "/v1/courses/:id", app.requirePermission("course:write", app.requireActivatedUser(app.updateCourseHandler)),)
	router.HandlerFunc(http.MethodDelete, "/v1/courses/:id", app.requirePermission("course:write", app.requireActivatedUser(app.deleteCourseHandler)),)
	router.HandlerFunc(http.MethodGet, "/v1/courses", app.requirePermission("course:read", app.requireActivatedUser(app.listCoursesHandler)),)

	// Course Postings
	router.HandlerFunc(http.MethodPost, "/v1/course/posting", app.requirePermission("course_posting:write", app.requireActivatedUser(app.createCoursePostingHandler)),)
	router.HandlerFunc(http.MethodGet, "/v1/course/posting/:id", app.requirePermission("course_posting:read", app.requireActivatedUser(app.displayCoursePostingHandler)),)
	router.HandlerFunc(http.MethodPatch, "/v1/course/posting/:id", app.requirePermission("course_posting:write", app.requireActivatedUser(app.updateCoursePostingHandler)),)
	router.HandlerFunc(http.MethodDelete, "/v1/course/posting/:id", app.requirePermission("course_posting:write", app.requireActivatedUser(app.deleteCoursePostingHandler)),)
	router.HandlerFunc(http.MethodGet, "/v1/course/posting", app.requirePermission("course_posting:read", app.requireActivatedUser(app.listCoursePostingsHandler)),)

	// Sessions
	router.HandlerFunc(http.MethodPost, "/v1/session", app.requirePermission("session:write", app.requireActivatedUser(app.createSessionHandler)),)
	router.HandlerFunc(http.MethodGet, "/v1/session/:id", app.requirePermission("session:read", app.requireActivatedUser(app.displaySessionHandler)),)
	router.HandlerFunc(http.MethodPatch, "/v1/session/:id", app.requirePermission("session:write", app.requireActivatedUser(app.updateSessionHandler)),)
	router.HandlerFunc(http.MethodDelete, "/v1/session/:id", app.requirePermission("session:write", app.requireActivatedUser(app.deleteSessionHandler)),)
	router.HandlerFunc(http.MethodGet, "/v1/session", app.requirePermission("session:read", app.requireActivatedUser(app.listSessionHandler)),)

	//User Session
	router.HandlerFunc(http.MethodPost, "/v1/user_session", app.requirePermission("user_session:write", app.requireActivatedUser(app.createUserSessionHandler)),)
	router.HandlerFunc(http.MethodGet, "/v1/user_session/:id", app.requirePermission("user_session:read", app.requireActivatedUser(app.getUserSessionHandler)),)
	router.HandlerFunc(http.MethodPatch, "/v1/user_session/:id", app.requirePermission("user_session:write", app.requireActivatedUser(app.updateUserSessionHandler)),)
	router.HandlerFunc(http.MethodDelete, "/v1/user_session/:id", app.requirePermission("user_session:write", app.requireActivatedUser(app.deleteUserSessionHandler)),)
	router.HandlerFunc(http.MethodGet, "/v1/user_session", app.requirePermission("user_session:read", app.requireActivatedUser(app.listUserSessionHandler)),)

	// Attendance
	router.HandlerFunc(http.MethodPost, "/v1/attendance", app.requirePermission("attendance:write", app.requireActivatedUser(app.createAttendanceHandler)),)
	router.HandlerFunc(http.MethodGet, "/v1/attendance/:id", app.requirePermission("user_session:read", app.requireActivatedUser(app.displayIndividualAttendanceHandler)),)
	router.HandlerFunc(http.MethodPatch, "/v1/attendance/:id", app.requirePermission("user_session:write", app.requireActivatedUser(app.updateAttendanceHandler)),)

	router.Handler(http.MethodGet, "/v1/observability/course/metrics", expvar.Handler())

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
	//return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(router))))
}
