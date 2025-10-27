// Filename: cmd/api/facilitator_rating_test.go
package main

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestAddFacilitatorRating_BadJSON(t *testing.T) {
    req := httptest.NewRequest(http.MethodPost, "/v1/facilitator-rating", bytes.NewBufferString("{bad json"))
    rr := httptest.NewRecorder()

    testApp.addFacilitatorRating(rr, req)

    if rr.Code != http.StatusBadRequest {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusBadRequest, rr.Code, rr.Body.String())
    }
}

func TestAddFacilitatorRating_InvalidData(t *testing.T) {
    payload := `{"user_id": 0, "rating": 0}`
    req := httptest.NewRequest(http.MethodPost, "/v1/facilitator-rating", bytes.NewBufferString(payload))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()

    testApp.addFacilitatorRating(rr, req)

    if rr.Code != http.StatusUnprocessableEntity {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusUnprocessableEntity, rr.Code, rr.Body.String())
    }
}

func TestDisplayFacilitatorRatingHandler_InvalidID(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/v1/facilitator-rating/", nil)
    rr := httptest.NewRecorder()

    testApp.displayFacilitatorRatingHandler(rr, req)

    if rr.Code != http.StatusNotFound {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusNotFound, rr.Code, rr.Body.String())
    }
}

func TestListFacilitatorRatingHandler_InvalidQueryParam(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/v1/facilitator-rating?page=notint", nil)
    rr := httptest.NewRecorder()

    testApp.listFacilitatorRatingHandler(rr, req)

    if rr.Code != http.StatusUnprocessableEntity {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusUnprocessableEntity, rr.Code, rr.Body.String())
    }
}