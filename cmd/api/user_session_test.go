package main

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestCreateUserSessionHandler_BadJSON(t *testing.T) {
    req := httptest.NewRequest(http.MethodPost, "/v1/usersessions", bytes.NewBufferString("{bad json"))
    rr := httptest.NewRecorder()

    testApp.createUserSessionHandler(rr, req)

    if rr.Code != http.StatusBadRequest {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusBadRequest, rr.Code, rr.Body.String())
    }
}

func TestCreateUserSessionHandler_InvalidData(t *testing.T) {
    payload := `{"trainee_id":0,"session_id":0,"credithours_completed":0,"grade":"","feedback":""}`
    req := httptest.NewRequest(http.MethodPost, "/v1/usersessions", bytes.NewBufferString(payload))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()

    testApp.createUserSessionHandler(rr, req)

    if rr.Code != http.StatusUnprocessableEntity {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusUnprocessableEntity, rr.Code, rr.Body.String())
    }
}

func TestGetUserSessionHandler_InvalidID(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/v1/usersessions/", nil)
    rr := httptest.NewRecorder()

    testApp.getUserSessionHandler(rr, req)

    if rr.Code != http.StatusNotFound {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusNotFound, rr.Code, rr.Body.String())
    }
}

func TestUpdateUserSessionHandler_InvalidID(t *testing.T) {
    req := httptest.NewRequest(http.MethodPatch, "/v1/usersessions/", nil)
    rr := httptest.NewRecorder()

    testApp.updateUserSessionHandler(rr, req)

    if rr.Code != http.StatusNotFound {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusNotFound, rr.Code, rr.Body.String())
    }
}

func TestDeleteUserSessionHandler_InvalidID(t *testing.T) {
    req := httptest.NewRequest(http.MethodDelete, "/v1/usersessions/", nil)
    rr := httptest.NewRecorder()

    testApp.deleteUserSessionHandler(rr, req)

    if rr.Code != http.StatusNotFound {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusNotFound, rr.Code, rr.Body.String())
    }
}