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

    // Users
    // router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
    // router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
    // router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

    // Courses
    // router.HandlerFunc(http.MethodPost, "/v1/course", app.requirePermission("course:write", app.createCourseHandler))
    // router.HandlerFunc(http.MethodGet, "/v1/courses/:id", app.requirePermission("course:read", app.displayCourseHandler))
    // router.HandlerFunc(http.MethodPatch, "/v1/courses/:id", app.requirePermission("course:write", app.updateCourseHandler))
    // router.HandlerFunc(http.MethodDelete, "/v1/courses/:id", app.requirePermission("course:write", app.deleteCourseHandler))
    // router.HandlerFunc(http.MethodGet, "/v1/courses", app.requirePermission("course:read", app.listCoursesHandler))

    // router.HandlerFunc(http.MethodPost, "/v1/session", app.requirePermission("session:write", app.createSessionHandler))
    // router.HandlerFunc(http.MethodGet, "/v1/sessions/:id", app.requirePermission("session:read", app.displaySessionHandler))
    // router.HandlerFunc(http.MethodPatch, "/v1/sessions/:id", app.requirePermission("session:write", app.updateSessionHandler))
    // router.HandlerFunc(http.MethodDelete, "/v1/sessions/:id", app.requirePermission("session:write", app.deleteSessionHandler))
    // router.HandlerFunc(http.MethodGet, "/v1/sessions", app.requirePermission("session:read", app.listSessionHandler))

    // Trainee
    // router.HandlerFunc(http.MethodPost, "/v1/trainee", app.requirePermission("trainee:write", app.createTraineeHandler))
    // router.HandlerFunc(http.MethodGet, "/v1/trainees/:id", app.requirePermission("trainee:read", app.displayTraineeHandler))
    // router.HandlerFunc(http.MethodPatch, "/v1/trainees/:id", app.requirePermission("trainee:write", app.updateTraineeHandler))
    // router.HandlerFunc(http.MethodDelete, "/v1/trainees/:id", app.requirePermission("trainee:write", app.deleteTraineeHandler))
    // router.HandlerFunc(http.MethodGet, "/v1/trainees", app.requirePermission("trainee:read", app.listTraineeHandler))

    // router.HandlerFunc(http.MethodPost, "/v1/enrollement", app.requirePermission("enrollement:write", app.createEnrollementHandler))
    // router.HandlerFunc(http.MethodGet, "/v1/enrollements/:id", app.requirePermission("enrollement:read", app.displayEnrollementHandler))
    // router.HandlerFunc(http.MethodPatch, "/v1/enrollements/:id", app.requirePermission("enrollement:write", app.updateEnrollementHandler))
    // router.HandlerFunc(http.MethodDelete, "/v1/enrollements/:id", app.requirePermission("enrollement:write", app.deleteEnrollementHandler))
    // router.HandlerFunc(http.MethodGet, "/v1/enrollements", app.requirePermission("enrollement:read", app.listEnrollementHandler))

    // router.HandlerFunc(http.MethodPost, "/v1/attendance", app.requirePermission("attendance:write", app.createAttendanceHandler))
    // router.HandlerFunc(http.MethodGet, "/v1/attendance/:id", app.requirePermission("attendance:read", app.displayAttendanceHandler))
    // router.HandlerFunc(http.MethodPatch, "/v1/attendance/:id", app.requirePermission("attendance:write", app.updateAttendanceHandler))
    // router.HandlerFunc(http.MethodGet, "/v1/session/attendance/:id", app.requirePermission("attendance:read", app.listEnrollementHandler))

    // Facilitator Rating
    // router.HandlerFunc(http.MethodPost, "/v1/frating", app.requirePermission("frating:write", app.createFacilitatorRatingHandler))
    // router.HandlerFunc(http.MethodGet, "/v1/fratings/:id", app.requirePermission("frating:read", app.displayFacilitatorRatingHandler))
    // router.HandlerFunc(http.MethodPatch, "/v1/fratings/:id", app.requirePermission("frating:write", app.updateFacilitatorRatingHandler))
    // router.HandlerFunc(http.MethodDelete, "/v1/fratings/:id", app.requirePermission("frating:write", app.deleteFacilitatorRatingHandler))
    // router.HandlerFunc(http.MethodGet, "/v1/fratings", app.requirePermission("frating:read", app.listFacilitatorRatingHandler))

    // Posting
    // router.HandlerFunc(http.MethodPost, "/v1/posting", app.requirePermission("posting:write", app.createPostingHandler))
    // router.HandlerFunc(http.MethodGet, "/v1/postings/:id", app.requirePermission("posting:read", app.displayPostingHandler))
    // router.HandlerFunc(http.MethodPatch, "/v1/postings/:id", app.requirePermission("posting:write", app.updatePostingHandler))
    // router.HandlerFunc(http.MethodDelete, "/v1/postings/:id", app.requirePermission("posting:write", app.deletePostingHandler))
    // router.HandlerFunc(http.MethodGet, "/v1/postings", app.requirePermission("posting:read", app.listPostingHandler))

    // Posting
    // router.HandlerFunc(http.MethodPost, "/v1/ranking", app.requirePermission("ranking:write", app.createRankingHandler))
    // router.HandlerFunc(http.MethodGet, "/v1/rankings/:id", app.requirePermission("ranking:read", app.displayRankingHandler))
    // router.HandlerFunc(http.MethodPatch, "/v1/rankings/:id", app.requirePermission("ranking:write", app.updateRankingHandler))
    // router.HandlerFunc(http.MethodDelete, "/v1/rankings/:id", app.requirePermission("ranking:write", app.deleteRankingHandler))
    // router.HandlerFunc(http.MethodGet, "/v1/rankings", app.requirePermission("ranking:read", app.listRankingHandler))

    // Checks
    // router.Handler(http.MethodGet, "/v1/observability/quotes/metrics", expvar.Handler())


	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}
