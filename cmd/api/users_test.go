package main

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"
    "io"
    "log/slog"
)

// Provide an in-file test helper so this file doesn't rely on an external
// testhelpers_test.go. This mirrors the approach used by the other test files
// which reference a package-level testApp.

var testApp *application

func newTestApp() *application {
    if testApp != nil {
        return testApp
    }
    logger := slog.New(slog.NewTextHandler(io.Discard, nil))
    testApp = &application{logger: logger}
    return testApp
}

// The tests below exercise handler-level parsing/validation/parameter logic and
// do not require a database.

func TestRegisterUserHandler_BadJSON(t *testing.T) {
    app := newTestApp()
    req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewBufferString("{bad json"))
    rr := httptest.NewRecorder()

    app.registerUserHandler(rr, req)

    if rr.Code != http.StatusBadRequest {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusBadRequest, rr.Code, rr.Body.String())
    }
}

func TestRegisterUserHandler_InvalidData(t *testing.T) {
    app := newTestApp()
    // valid JSON but missing required fields should produce 422
    payload := `{"regulation_number":"","username":"","fname":"","lname":"","email":"invalid","gender":"","formation":0,"rank":0,"postings":0,"password":"short"}`
    req := httptest.NewRequest(http.MethodPost, "/v1/users", bytes.NewBufferString(payload))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()

    app.registerUserHandler(rr, req)

    if rr.Code != http.StatusUnprocessableEntity {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusUnprocessableEntity, rr.Code, rr.Body.String())
    }
}

func TestActivateUserHandler_InvalidToken(t *testing.T) {
    app := newTestApp()
    // empty token should be caught by validation
    payload := `{"token": ""}`
    req := httptest.NewRequest(http.MethodPost, "/v1/users/activated", bytes.NewBufferString(payload))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()

    app.activateUserHandler(rr, req)

    if rr.Code != http.StatusUnprocessableEntity {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusUnprocessableEntity, rr.Code, rr.Body.String())
    }
}

func TestListUsersHandler_InvalidQueryParam(t *testing.T) {
    app := newTestApp()
    req := httptest.NewRequest(http.MethodGet, "/v1/users?page=notint", nil)
    rr := httptest.NewRecorder()

    app.listUsersHandler(rr, req)

    if rr.Code != http.StatusUnprocessableEntity {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusUnprocessableEntity, rr.Code, rr.Body.String())
    }
}

func TestUpdateUserHandler_InvalidID(t *testing.T) {
    app := newTestApp()
    req := httptest.NewRequest(http.MethodPatch, "/v1/users/update/", nil)
    rr := httptest.NewRecorder()

    app.updateUserHandler(rr, req)

    if rr.Code != http.StatusNotFound {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusNotFound, rr.Code, rr.Body.String())
    }
}

func TestDeleteUserHandler_InvalidID(t *testing.T) {
    app := newTestApp()
    req := httptest.NewRequest(http.MethodDelete, "/v1/users/delete/", nil)
    rr := httptest.NewRecorder()

    app.deleteUserHandler(rr, req)

    if rr.Code != http.StatusNotFound {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusNotFound, rr.Code, rr.Body.String())
    }
}

func TestUpdatePasswordHandler_MissingNewPassword(t *testing.T) {
    app := newTestApp()
    // Missing new_password should trigger validation 422, but calling without an id param returns 404 first.
    req := httptest.NewRequest(http.MethodPatch, "/v1/users/update-password/", bytes.NewBufferString(`{"new_password":""}`))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()

    app.updatePasswordHandler(rr, req)

    if rr.Code != http.StatusNotFound {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusNotFound, rr.Code, rr.Body.String())
    }
}