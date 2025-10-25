package main

import (
	// "bytes"
	// "context"
	// "database/sql"
	// "encoding/json"
	// "log"
	// "log/slog"
	// "net/http"
	// "net/http/httptest"
	// "os"

	// // "strconv"
	// "testing"
	// "time"

	// _ "github.com/lib/pq"

	// // Adjust the path as necessary for your project
	// "github.com/kelseyaban/National-Inservice-Training-Database/internal/data"
	// "github.com/kelseyaban/National-Inservice-Training-Database/internal/mailer"
)

// // --- Struct Definitions (Matching your application's types) ---

// type config struct {
// 	port int
// 	env  string
// 	cors struct {
// 		trustedOrigins []string
// 	}
// 	limiter struct {
// 		rps     float64
// 		burst   int
// 		enabled bool
// 	}
// }

// // type application struct {
// // 	config          configuration
// // 	logger          *slog.Logger
// // 	userModel       data.UserModel
// // 	tokenModel      data.TokenModel
// // 	permissionModel data.PermissionModel
// // 	roleModel       data.RoleModel
// // 	mailer          mailer.Mailer
// // 	attendanceModel data.AttendanceModel
// // }

// // --- Global Variables and Test Setup ---

// var testApp *application

// func TestMain(m *testing.M) {
// 	testApp = setupTestApp()
// 	exitCode := m.Run()
// 	os.Exit(exitCode)
// }

// // setupTestApp connects to the test database and sets up the app dependencies.
// func setupTestApp() *application {
// 	dsn := os.Getenv("TEST_DB_DSN")
// 	if dsn == "" {
// 		// Replace with your actual DSN if necessary
// 		dsn = "postgres://nationalitdb:t1advweb@localhost/nationalitdb?sslmode=disable"
// 	}

// 	db, err := sql.Open("postgres", dsn)
// 	if err != nil {
// 		log.Fatalf("error connecting to test database: %v", err)
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	if err = db.PingContext(ctx); err != nil {
// 		log.Fatalf("error pinging test database: %v", err)
// 	}
// 	log.Println("Successfully connected to test database.")

// 	testConfig := configuration{
// 		port: 4000,
// 		env:  "test",
// 		cors: struct {
// 			trustedOrigins []string
// 		}{trustedOrigins: []string{"http://localhost:3000"}},
// 		limiter: struct {
// 			rps     float64
// 			burst   int
// 			enabled bool
// 		}{rps: 2.0, burst: 4, enabled: false},
// 	}

// 	testMailer := mailer.New("sandbox.smtp.mailtrap.io", 2525, "213718fe792166", "54177abe81857f", "Testing <test@example.com>")
// 	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

// 	app := &application{
// 		config:          testConfig,
// 		userModel:       data.UserModel{DB: db},
// 		mailer:          testMailer,
// 		tokenModel:      data.TokenModel{DB: db},
// 		logger:          logger,
// 		permissionModel: data.PermissionModel{DB: db},
// 		roleModel:       data.RoleModel{DB: db},
// 		attendanceModel: data.AttendanceModel{DB: db},
// 	}
// 	return app
// }

// // --- Test Helper Functions ---

// // clearTestAttendance deletes all records from the attendance table.
// func clearTestAttendance(db *sql.DB) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()
// 	// Using TRUNCATE with RESTART IDENTITY can often be faster and resets sequence IDs.
// 	// You might also need CASCADE if 'attendance' is referenced by other tables.
// 	query := `TRUNCATE attendance RESTART IDENTITY CASCADE`
// 	_, err := db.ExecContext(ctx, query)
// 	if err != nil {
// 		// Log the error but don't fail the test setup unless absolutely necessary
// 		log.Printf("Warning: Failed to clear attendance table: %v", err)
// 	}
// }

// // clearTestUserSession deletes all records from the user_session table.
// func clearTestUserSession(db *sql.DB) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()
// 	// TRUNCATE the parent table as well.
// 	query := `TRUNCATE user_session RESTART IDENTITY CASCADE`
// 	_, err := db.ExecContext(ctx, query)
// 	if err != nil {
// 		log.Fatalf("Fatal: Failed to clear user_session table: %v", err)
// 	}
// }

// // createTestUserSession inserts a minimal, valid user_session record and returns its ID.
// // This is critical to satisfy the foreign key constraint for the attendance table.
// func createTestUserSession(db *sql.DB) (int64, error) {
// 	// NOTE: If 'session_id' or 'trainee' in user_session are FKs, they must exist first!
// 	// For this example, we assume we can insert placeholder values.

// 	sessionIDPlaceholder := 3
// 	traineeIDPlaceholder := 3

// 	query := `
//         INSERT INTO attendance (user_session_id, attendance,date ) 
//         VALUES ($1, $2, $3) 
//         RETURNING id`

// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	var userSessionID int64
// 	args := []any{
// 		sessionIDPlaceholder,
// 		0.0,                  // credithours_completed (placeholder)
// 		"A",                  // grade (placeholder)
// 		"",                   // feedback (placeholder)
// 		1,                    // version (placeholder)
// 		traineeIDPlaceholder, // trainee (placeholder)
// 	}

// 	err := db.QueryRowContext(ctx, query, args...).Scan(&userSessionID)

// 	return userSessionID, err
// }

// // --- Tests for /v1/attendance (Creation Handler) ---

// // TestCreateAttendanceHandler_Success tests the successful creation of an attendance record (HTTP 201).
// func TestCreateAttendanceHandler_Success(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping integration test")
// 	}

// 	// 1. Clean both parent and child tables before the test
// 	clearTestAttendance(testApp.attendanceModel.DB)
// 	clearTestUserSession(testApp.attendanceModel.DB)

// 	// 2. Setup: Insert the required parent record first (to avoid FK violation)
// 	validUserSessionID, err := createTestUserSession(testApp.attendanceModel.DB)
// 	if err != nil {
// 		t.Fatalf("Failed to create required UserSession parent record: %v", err)
// 	}

// 	// 3. Prepare Request Body using the valid ID
// 	testTime := "2025-10-23"
// 	postBody := map[string]any{
// 		"user_session_id": validUserSessionID,
// 		"attendance":      true,
// 		"date":            testTime,
// 	}
// 	body, _ := json.Marshal(postBody)

// 	// 4. Create Request and Recorder
// 	req := httptest.NewRequest(http.MethodPost, "/v1/attendance", bytes.NewReader(body))
// 	req.Header.Set("Content-Type", "application/json")
// 	rr := httptest.NewRecorder()

// 	// 5. Execute Handler
// 	testApp.createAttendanceHandler(rr, req)

// 	// 6. Assert Response
// 	if rr.Code != http.StatusCreated {
// 		t.Fatalf("expected status %d, got %d. Body: %s", http.StatusCreated, rr.Code, rr.Body.String())
// 	}

// 	// Assert Location header
// 	locationHeader := rr.Header().Get("Location")
// 	if locationHeader == "" {
// 		t.Error("expected Location header to be set")
// 	}
// }

// // TestCreateAttendanceHandler_ValidationFailure tests failure due to missing required fields (HTTP 422).
// func TestCreateAttendanceHandler_ValidationFailure(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping integration test")
// 	}
// 	// Clean the table before the test (optional for validation, but good practice)
// 	clearTestAttendance(testApp.attendanceModel.DB)

// 	// Missing 'user_session_id' (set to 0, which is invalid by your ValidateAttendance function)
// 	postBody := map[string]any{
// 		"user_session_id": 0,
// 		"attendance":      true,
// 		"date":            "2025-10-23",
// 	}
// 	body, _ := json.Marshal(postBody)

// 	req := httptest.NewRequest(http.MethodPost, "/v1/attendance", bytes.NewReader(body))
// 	req.Header.Set("Content-Type", "application/json")
// 	rr := httptest.NewRecorder()

// 	testApp.createAttendanceHandler(rr, req)

// 	// Assert response is a validation error (422 Unprocessable Entity)
// 	if rr.Code != http.StatusUnprocessableEntity {
// 		t.Fatalf("expected status %d for validation error, got %d. Body: %s", http.StatusUnprocessableEntity, rr.Code, rr.Body.String())
// 	}

// 	// Check if the expected error message for the missing field is present
// 	if !bytes.Contains(rr.Body.Bytes(), []byte("user_session_id")) {
// 		t.Error("expected validation error for user_session_id not found in response body")
// 	}
// }
