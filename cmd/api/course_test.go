package main

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestCreateCourseHandler_BadJSON(t *testing.T) {
    req := httptest.NewRequest(http.MethodPost, "/v1/courses", bytes.NewBufferString("{bad json"))
    rr := httptest.NewRecorder()

    testApp.createCourseHandler(rr, req)

    if rr.Code != http.StatusBadRequest {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusBadRequest, rr.Code, rr.Body.String())
    }
}

func TestCreateCourseHandler_InvalidData(t *testing.T) {
    // empty course name should trigger validation error
    payload := `{"course":"", "description":""}`
    req := httptest.NewRequest(http.MethodPost, "/v1/courses", bytes.NewBufferString(payload))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()

    testApp.createCourseHandler(rr, req)

    if rr.Code != http.StatusUnprocessableEntity {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusUnprocessableEntity, rr.Code, rr.Body.String())
    }
}

func TestDisplayCourseHandler_InvalidID(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/v1/courses/", nil)
    rr := httptest.NewRecorder()

    testApp.displayCourseHandler(rr, req)

    if rr.Code != http.StatusNotFound {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusNotFound, rr.Code, rr.Body.String())
    }
}

func TestUpdateCourseHandler_InvalidID(t *testing.T) {
    req := httptest.NewRequest(http.MethodPatch, "/v1/courses/", nil)
    rr := httptest.NewRecorder()

    testApp.updateCourseHandler(rr, req)

    if rr.Code != http.StatusNotFound {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusNotFound, rr.Code, rr.Body.String())
    }
}

func TestDeleteCourseHandler_InvalidID(t *testing.T) {
    req := httptest.NewRequest(http.MethodDelete, "/v1/courses/", nil)
    rr := httptest.NewRecorder()

    testApp.deleteCourseHandler(rr, req)

    if rr.Code != http.StatusNotFound {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusNotFound, rr.Code, rr.Body.String())
    }
}

func TestListCoursesHandler_InvalidQueryParam(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/v1/courses?page=notint", nil)
    rr := httptest.NewRecorder()

    testApp.listCoursesHandler(rr, req)

    if rr.Code != http.StatusUnprocessableEntity {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusUnprocessableEntity, rr.Code, rr.Body.String())
    }
}
