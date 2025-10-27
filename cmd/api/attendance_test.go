package main

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestCreateAttendanceHandler_BadJSON(t *testing.T) {
    // Prepare a request with invalid JSON
    req := httptest.NewRequest(http.MethodPost, "/v1/attendance", bytes.NewBufferString("{bad json"))
    rr := httptest.NewRecorder()

    // Call the handler directly (bypass middleware)
    testApp.createAttendanceHandler(rr, req)

    if rr.Code != http.StatusBadRequest {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusBadRequest, rr.Code, rr.Body.String())
    }
}

func TestCreateAttendanceHandler_InvalidData(t *testing.T) {
    // Missing date (empty string) should trigger validation error
    payload := `{"user_session_id": 1, "attendance": true, "date": ""}`
    req := httptest.NewRequest(http.MethodPost, "/v1/attendance", bytes.NewBufferString(payload))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()

    testApp.createAttendanceHandler(rr, req)

    if rr.Code != http.StatusUnprocessableEntity {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusUnprocessableEntity, rr.Code, rr.Body.String())
    }
}

func TestDisplayIndividualAttendanceHandler_InvalidID(t *testing.T) {
    // No ID param in context should result in 404
    req := httptest.NewRequest(http.MethodGet, "/v1/attendance/", nil)
    rr := httptest.NewRecorder()

    testApp.displayIndividualAttendanceHandler(rr, req)

    if rr.Code != http.StatusNotFound {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusNotFound, rr.Code, rr.Body.String())
    }
}

func TestUpdateAttendanceHandler_InvalidID(t *testing.T) {
    // No ID param in context should result in 404
    req := httptest.NewRequest(http.MethodPatch, "/v1/attendance/", nil)
    rr := httptest.NewRecorder()

    testApp.updateAttendanceHandler(rr, req)

    if rr.Code != http.StatusNotFound {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusNotFound, rr.Code, rr.Body.String())
    }
}