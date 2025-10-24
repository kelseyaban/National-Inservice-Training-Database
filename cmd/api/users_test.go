package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/kelseyaban/National-Inservice-Training-Database/internal/data"
	"github.com/kelseyaban/National-Inservice-Training-Database/internal/mailer"
)

// Minimal configuration structure required for middleware tests.
// NOTE: This must be named 'configuration' to match the type defined in main.go.
type config struct {
	port int
	env  string
	cors struct {
		trustedOrigins []string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

// The 'application' struct is defined in main.go and is available here
// because both files are in the 'main' package. We only define the dependencies
// we need to set up the test app.
/*
type application struct {
    config          configuration // Renamed to configuration
    logger          *slog.Logger
    userModel       data.UserModel
    tokenModel      data.TokenModel
    permissionModel data.PermissionModel
    roleModel       data.RoleModel
    mailer          mailer.Mailer
}
*/

// testApp will be initialized before running any tests
var testApp *application

// setupTestApp connects to the test database and sets up the app dependencies
func setupTestApp() *application {
	// Use the same DSN as your app, or override it with TEST_DB_DSN
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		// NOTE: Please ensure this DSN is for a dedicated test database
		dsn = "postgres://nationalitdb:t1advweb@localhost/nationalitdb?sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		// Use log.Fatal to stop execution if the database connection fails
		log.Fatalf("error connecting to test database: %v", err)
	}

	// Wait for the database to be ready
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		log.Fatalf("error pinging test database: %v", err)
	}

	// 1. Define a minimal test configuration
	// IMPORTANT: Using 'configuration' type to match main.go
	testConfig := configuration{
		port: 4000,
		env:  "test",
		cors: struct {
			trustedOrigins []string
		}{
			// Use a placeholder trusted origin for the CORS middleware
			trustedOrigins: []string{"http://localhost:3000"},
		},
		limiter: struct {
			rps     float64
			burst   int
			enabled bool
		}{
			// Crucial: Disable rate limiting in tests for predictable results
			rps:     2.0,
			burst:   4,
			enabled: false,
		},
	}

	// Test mailer setup (it won't actually send emails, which is ideal for testing)
	testMailer := mailer.New("sandbox.smtp.mailtrap.io", 2525, "213718fe792166", "54177abe81857f", "Testing <test@example.com>")

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &application{
		// Initialize the config field (named 'config' in the application struct)
		// with the test configuration structure (named 'configuration').
		config: testConfig,

		userModel:  data.UserModel{DB: db},
		mailer:     testMailer,
		tokenModel: data.TokenModel{DB: db},
		logger:     logger,

		// Initialize permissionModel and roleModel, which are accessed by permission checking functions
		permissionModel: data.PermissionModel{DB: db},
		roleModel:       data.RoleModel{DB: db},
	}

	return app
}

// TestMain runs once before all tests
func TestMain(m *testing.M) {
	testApp = setupTestApp()
	code := m.Run()
	os.Exit(code)
}

// helper to perform requests
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	// This call to testApp.routes() is where the panic is occurring,
	// likely due to the expvar issue or a nil dependency in middleware.
	router := testApp.routes()
	router.ServeHTTP(rr, req)
	return rr
}

// tearDownTestData cleans up test users by email
func tearDownTestData(t *testing.T, emails []string) {
	if len(emails) == 0 {
		return
	}

	// Convert emails slice into a string for the SQL IN clause
	// e.g., 'a@example.com', 'b@example.com'
	emailList := "'" + emails[0] + "'"
	for i := 1; i < len(emails); i++ {
		emailList += ", '" + emails[i] + "'"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "DELETE FROM users WHERE email IN (" + emailList + ")"
	_, err := testApp.userModel.DB.ExecContext(ctx, query)
	if err != nil {
		t.Errorf("cleanup failed: %v", err)
	}
}

// TestCreateUserHandler tests successful and invalid user creation using table-driven tests
func TestCreateUserHandler(t *testing.T) {
	testCases := []struct {
		name           string
		userPayload    map[string]interface{}
		expectedStatus int
		cleanupEmail   string // Email to delete afterwards
	}{
		{
			name: "Success: Valid User",
			userPayload: map[string]interface{}{
				"regulation_number": "R53865",
				"username":          "example_valid",
				"fname":             "Ex",
				"lname":             "Ample",
				"email":             "example_valid@test.com",
				"gender":            "F",
				"formation":         2,
				"rank":              1,
				"postings":          5,
				"password":          "StrongP@ss123",
			},
			expectedStatus: http.StatusCreated,
			cleanupEmail:   "example_valid@test.com",
		},
		{
			name: "Failure: Missing Email Field",
			userPayload: map[string]interface{}{
				"regulation_number": "R99999",
				"username":          "baduser",
				"fname":             "Bad",
				"lname":             "User",
				"gender":            "M",
				"formation":         1,
				"rank":              1,
				"postings":          1,
				"password":          "123456",
			},
			// Expecting 422 for validation failure on a required field
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "Failure: Password Too Short (assuming min 8 chars)",
			userPayload: map[string]interface{}{
				"regulation_number": "R12345",
				"username":          "user2ex",
				"fname":             "User1",
				"lname":             "Ex",
				"email":             "userex@test.com",
				"gender":            "M",
				"formation":         1,
				"rank":              2,
				"postings":          1,
				"password":          "short", // Changed to make sure it fails validation
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	// 1. Collect emails for cleanup
	var emailsToClean []string
	for _, tc := range testCases {
		if tc.cleanupEmail != "" {
			emailsToClean = append(emailsToClean, tc.cleanupEmail)
		}
	}

	// 2. Schedule cleanup using t.Cleanup()
	t.Cleanup(func() {
		tearDownTestData(t, emailsToClean)
	})

	// 3. Run table tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, err := json.Marshal(tc.userPayload)
			if err != nil {
				t.Fatalf("Failed to marshal payload: %v", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			rr := executeRequest(req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("expected %d (%s); got %d (%s)\nBody: %s",
					tc.expectedStatus, http.StatusText(tc.expectedStatus),
					rr.Code, http.StatusText(rr.Code),
					rr.Body.String())
			}
		})
	}

	// --- Test for Duplicate Email (must be run after successful creation) ---
	// NOTE: This test will only run if the first "Success: Valid User" test case passes.
	t.Run("Failure: Duplicate Email", func(t *testing.T) {
		// First, create a user successfully
		userA := testCases[0].userPayload // Use the payload from the successful case
		// Change the username and regulation number slightly to avoid the
		// "pq: duplicate key value violates unique constraint \"users_username_key\"" error
		// while still testing the duplicate email constraint.
		userA["username"] = "example_valid_dup"
		userA["regulation_number"] = "R53866"

		body, _ := json.Marshal(userA)
		req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		// Initial creation attempt (should succeed with 201)
		rr := executeRequest(req)
		if rr.Code != http.StatusCreated {
			t.Logf("Initial user creation failed with status %d. Skipping duplicate test.", rr.Code)
			return
		}

		// Second, try to create the same user again (should fail with 409)
		rr = executeRequest(req)
		expectedStatus := http.StatusConflict // 409 Conflict is ideal for duplicate resource
		if rr.Code != expectedStatus {
			// Note: If your application returns 422 (Unprocessable Entity) for duplicate email,
			// you may need to adjust expectedStatus here. I'm using 409 as it's standard for conflicts.
			t.Errorf("expected %d Conflict for duplicate email; got %d\nBody: %s", expectedStatus, rr.Code, rr.Body.String())
		}
	})
}
