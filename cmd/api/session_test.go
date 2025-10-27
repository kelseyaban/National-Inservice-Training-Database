package main

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestCreateSessionHandler_BadJSON(t *testing.T) {
    req := httptest.NewRequest(http.MethodPost, "/v1/session", bytes.NewBufferString("{bad json"))
    rr := httptest.NewRecorder()

    testApp.createSessionHandler(rr, req)

    if rr.Code != http.StatusBadRequest {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusBadRequest, rr.Code, rr.Body.String())
    }
}

func TestCreateSessionHandler_InvalidData(t *testing.T) {
    payload := `{"course_id":0,"formation_id":0,"facilitator_id":0}`
    req := httptest.NewRequest(http.MethodPost, "/v1/session", bytes.NewBufferString(payload))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()

    testApp.createSessionHandler(rr, req)

    if rr.Code != http.StatusUnprocessableEntity {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusUnprocessableEntity, rr.Code, rr.Body.String())
    }
}

func TestDisplaySessionHandler_InvalidID(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/v1/session/", nil)
    rr := httptest.NewRecorder()

    testApp.displaySessionHandler(rr, req)

    if rr.Code != http.StatusNotFound {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusNotFound, rr.Code, rr.Body.String())
    }
}

func TestUpdateSessionHandler_InvalidID(t *testing.T) {
    req := httptest.NewRequest(http.MethodPatch, "/v1/session/", nil)
    rr := httptest.NewRecorder()

    testApp.updateSessionHandler(rr, req)

    if rr.Code != http.StatusNotFound {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusNotFound, rr.Code, rr.Body.String())
    }
}

func TestDeleteSessionHandler_InvalidID(t *testing.T) {
    req := httptest.NewRequest(http.MethodDelete, "/v1/session/", nil)
    rr := httptest.NewRecorder()

    testApp.deleteSessionHandler(rr, req)

    if rr.Code != http.StatusNotFound {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusNotFound, rr.Code, rr.Body.String())
    }
}

func TestListSessionHandler_InvalidQueryParam(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/v1/session?page=notint", nil)
    rr := httptest.NewRecorder()

    testApp.listSessionHandler(rr, req)

    if rr.Code != http.StatusUnprocessableEntity {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusUnprocessableEntity, rr.Code, rr.Body.String())
    }
}