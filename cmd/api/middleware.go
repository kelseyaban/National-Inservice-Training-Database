// Filename: cmd/api/middleware.go
package main

import (
	// "errors"
	"expvar"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kelseyaban/National-Inservice-Training-Database/internal/data"
	"github.com/kelseyaban/National-Inservice-Training-Database/internal/validator"
	"golang.org/x/time/rate"
)

func (a *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// defer will be called when the stack unwinds
		defer func() {
			// recover() checks for panics
			err := recover()
			if err != nil {
				w.Header().Set("Connection", "close")
				a.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// add CORS headers to the response
func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Vary", "Origin")
		// The request method can vary so don't rely on cache
		w.Header().Add("Vary", "Access-Control-Request-Method")
		// Check if the request origin is in the trusted list
		origin := r.Header.Get("Origin")

		if origin != "" {
			for i := range app.config.cors.trustedOrigins {
				if origin == app.config.cors.trustedOrigins[i] {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					// Check if its a preflight CORS request
					if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-method") != "" {
						w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
						w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
						w.WriteHeader(http.StatusOK)
						return
					}
					break
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time // remove map enteries that are stale
	}

	var mu sync.Mutex                      // use to synchronize the map
	var clients = make(map[string]*client) // the actual map

	// A goroutine to remove stale entries from the map

	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock() // begin cleanup
			// delete any entry not seen in three minutes
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock() // finish clean up
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if app.config.limiter.enabled {
			// get the IP address
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			mu.Lock() // exclusive access to the map
			// check if ip address already in map, if not add it
			_, found := clients[ip]
			if !found {
				clients[ip] = &client{limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst)}
			}

			// Update the last seem for the client
			clients[ip].lastSeen = time.Now()

			// Check the rate limit status
			if !clients[ip].limiter.Allow() {
				mu.Unlock() // no longer need exclusive access to the map
				app.rateLimitExceededResponse(w, r)
				return
			}

			mu.Unlock() // others are free to get exclusive access to the map
		}
		next.ServeHTTP(w, r)
	})

}

func (a *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Tells the servers not to cache the response when
		// the Authorization header changes.
		w.Header().Add("Vary", "Authorization")

		// Get the Authorization header from the request. It should have the
		// Bearer token
		authorizationHeader := r.Header.Get("Authorization")

		// // If there is no Authorization header then we have an Anonymous user
		// if authorizationHeader == "" {
		// 	r = a.contextSetUser(r, data.AnonymousUser)
		// 	next.ServeHTTP(w, r)
		// 	return
		// }

		// Bearer token present so parse it. The Bearer token is in the form
		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			a.invalidAuthenticationTokenResponse(w, r)
			return
		}

		// Get the actual token
		token := headerParts[1]
		// Validate
		v := validator.New()

		data.ValidateTokenPlaintext(v, token)
		if !v.IsEmpty() {
			a.invalidAuthenticationTokenResponse(w, r)
			return
		}

		// Get the user info associated with this authentication token
		// user, err := a.userModel.GetForToken(data.ScopeAuthentication, token)
		// if err != nil {
		// 	switch {
		// 	case errors.Is(err, data.ErrRecordNotFound):
		// 		a.invalidAuthenticationTokenResponse(w, r)
		// 	default:
		// 		a.serverErrorResponse(w, r, err)
		// 	}
		// 	return
		// }

		// Add the retrieved user info to the context
		// r = a.contextSetUser(r, user)

		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}

// Note: 401 is Unauthorized  and 403 is Forbidden (

func (a *application) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	a.errorResponseJSON(w, r, http.StatusUnauthorized, message)
}

func (a *application) inactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account must be activated to access this resource"
	a.errorResponseJSON(w, r, http.StatusForbidden, message)
}

// This middleware checks if the user is authenticated (not anonymous)
// func (a *application) requireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		user := a.contextGetUser(r)

// 		if user.IsAnonymous() {
// 			a.authenticationRequiredResponse(w, r)
// 			return
// 		}
// 		next.ServeHTTP(w, r)
// 	})
// }

// This middleware checks if the user is activated
// It call the authentication middleware to help it do its job
// func (a *application) requireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
// 	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		user := a.contextGetUser(r)

// 		if !user.Activated {
// 			a.inactiveAccountResponse(w, r)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// 	// Only check if the user is activated if they are actually authenticated.
// 	return a.requireAuthenticatedUser(fn)
// }

// Checks if the user has the right permissions
// We send the permission that is expected as an argument
// func (a *application) requirePermission(permissionCode string, next http.HandlerFunc) http.HandlerFunc {

// 	fn := func(w http.ResponseWriter, r *http.Request) {
// 		user := a.contextGetUser(r)
// 		// get all the permissions associated with the user
// 		permissions, err := a.permissionModel.GetAllForUser(user.ID)
// 		if err != nil {
// 			a.serverErrorResponse(w, r, err)
// 			return
// 		}
// 		if !permissions.Include(permissionCode) {
// 			a.notPermittedResponse(w, r)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	}

// 	return a.requireActivatedUser(fn)

// }

// Run for every request received
func (a *application) metrics(next http.Handler) http.Handler {
	// Setup our variable to track the metrics
	var (
		totalResponsesSentByStatus      = expvar.NewMap("total_responses_sent_by_status")
		totalRequestsReceived           = expvar.NewInt("total_requests_received")
		totalResponsesSent              = expvar.NewInt("total_responses_sent")
		totalProcessingTimeMicroseconds = expvar.NewInt("total_processing_time_μs")
	)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// start is when we receive the request and start processing it
		start := time.Now()
		// update our request received counter
		totalRequestsReceived.Add(1)

		// create a custom responseWriter
		mw := newMetricsResponseWriter(w)
		// we send our custom responseWriter down the middleware chain
		next.ServeHTTP(mw, r)

		// When we return back to our middleware; we will increment
		// the responses sent counter
		totalResponsesSent.Add(1)

		// extract the status code for use in our metrics since we have returned
		// from the middleware chain.
		totalResponsesSentByStatus.Add(strconv.Itoa(mw.statusCode), 1)
		// calculate the processing time for this request.
		duration := time.Since(start).Microseconds()
		totalProcessingTimeMicroseconds.Add(duration)
	})
}

// A custom response writer to capture the status code
type metricsResponseWriter struct {
	wrapped       http.ResponseWriter // the original http.ResponseWriter
	statusCode    int                 // this will contain the status code we need
	headerWritten bool                // has the response headers already been written?
}

// Create an new instance of our custom http.ResponseWriter
// We will set the status code to 200 by default
func newMetricsResponseWriter(w http.ResponseWriter) *metricsResponseWriter {
	return &metricsResponseWriter{
		wrapped:    w,
		statusCode: http.StatusOK,
	}
}

// Call the original http.ResponseWriter's Header()
// method when our custom http.ResponseWriter's Header() method is called
func (mw *metricsResponseWriter) Header() http.Header {
	return mw.wrapped.Header()
}

// Write the status code that is provided
func (mw *metricsResponseWriter) WriteHeader(statusCode int) {
	mw.wrapped.WriteHeader(statusCode)
	// After the call to WriteHeader() returns, we record
	// the status code for use in our metrics
	if !mw.headerWritten {
		mw.statusCode = statusCode
		mw.headerWritten = true
	}
}

// The write() method simply calls the original http.ResponseWriter's
// Write() method which write the data to the connection
func (mw *metricsResponseWriter) Write(b []byte) (int, error) {
	mw.headerWritten = true
	return mw.wrapped.Write(b)
}

// We need a function to get the original http.ResponseWriter
func (mw *metricsResponseWriter) Unwrap() http.ResponseWriter {
	return mw.wrapped
}
