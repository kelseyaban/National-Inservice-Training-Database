// filename: cmd/api/course_posting_test.go
package main

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestCreateCoursePostingHandler_BadJSON(t *testing.T) {
    req := httptest.NewRequest(http.MethodPost, "/v1/course/posting", bytes.NewBufferString("{bad json"))
    rr := httptest.NewRecorder()

    testApp.createCoursePostingHandler(rr, req)

    if rr.Code != http.StatusBadRequest {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusBadRequest, rr.Code, rr.Body.String())
    }
}

func TestCreateCoursePostingHandler_InvalidData(t *testing.T) {
    // missing required numeric fields (course_id, posting_id, rank_id)
    payload := `{"course_id": 0, "posting_id": 0, "mandatory": true, "credithours": -1, "rank_id": 0}`
    req := httptest.NewRequest(http.MethodPost, "/v1/course/posting", bytes.NewBufferString(payload))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()

    testApp.createCoursePostingHandler(rr, req)

    if rr.Code != http.StatusUnprocessableEntity {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusUnprocessableEntity, rr.Code, rr.Body.String())
    }
}

func TestDisplayCoursePostingHandler_InvalidID(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/v1/course/posting/", nil)
    rr := httptest.NewRecorder()

    testApp.displayCoursePostingHandler(rr, req)

    if rr.Code != http.StatusNotFound {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusNotFound, rr.Code, rr.Body.String())
    }
}

func TestUpdateCoursePostingHandler_InvalidID(t *testing.T) {
    req := httptest.NewRequest(http.MethodPatch, "/v1/course/posting/", nil)
    rr := httptest.NewRecorder()

    testApp.updateCoursePostingHandler(rr, req)

    if rr.Code != http.StatusNotFound {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusNotFound, rr.Code, rr.Body.String())
    }
}

func TestDeleteCoursePostingHandler_InvalidID(t *testing.T) {
    req := httptest.NewRequest(http.MethodDelete, "/v1/course/posting/", nil)
    rr := httptest.NewRecorder()

    testApp.deleteCoursePostingHandler(rr, req)

    if rr.Code != http.StatusNotFound {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusNotFound, rr.Code, rr.Body.String())
    }
}

func TestListCoursePostingsHandler_InvalidQueryParam(t *testing.T) {
    // invalid page (should be integer)
    req := httptest.NewRequest(http.MethodGet, "/v1/course/posting?page=notint", nil)
    rr := httptest.NewRecorder()

    testApp.listCoursePostingsHandler(rr, req)

    if rr.Code != http.StatusUnprocessableEntity {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusUnprocessableEntity, rr.Code, rr.Body.String())
    }
}