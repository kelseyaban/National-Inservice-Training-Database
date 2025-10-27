package main

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestCreateRoleHandler_BadJSON(t *testing.T) {
    req := httptest.NewRequest(http.MethodPost, "/v1/roles", bytes.NewBufferString("{bad json"))
    rr := httptest.NewRecorder()

    testApp.createRoleHandler(rr, req)

    if rr.Code != http.StatusBadRequest {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusBadRequest, rr.Code, rr.Body.String())
    }
}

func TestCreateRoleHandler_InvalidData(t *testing.T) {
    payload := `{"role":""}`
    req := httptest.NewRequest(http.MethodPost, "/v1/roles", bytes.NewBufferString(payload))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()

    testApp.createRoleHandler(rr, req)

    if rr.Code != http.StatusUnprocessableEntity {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusUnprocessableEntity, rr.Code, rr.Body.String())
    }
}

func TestDisplayRoleHandler_InvalidID(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/v1/roles/", nil)
    rr := httptest.NewRecorder()

    testApp.displayRoleHandler(rr, req)

    if rr.Code != http.StatusNotFound {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusNotFound, rr.Code, rr.Body.String())
    }
}

func TestUpdateRoleHandler_InvalidID(t *testing.T) {
    req := httptest.NewRequest(http.MethodPatch, "/v1/roles/", nil)
    rr := httptest.NewRecorder()

    testApp.updateRoleHandler(rr, req)

    if rr.Code != http.StatusNotFound {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusNotFound, rr.Code, rr.Body.String())
    }
}

func TestDeleteRoleHandler_InvalidID(t *testing.T) {
    req := httptest.NewRequest(http.MethodDelete, "/v1/roles/", nil)
    rr := httptest.NewRecorder()

    testApp.deleteRoleHandler(rr, req)

    if rr.Code != http.StatusNotFound {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusNotFound, rr.Code, rr.Body.String())
    }
}

func TestListRoleHandler_InvalidQueryParam(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/v1/roles?page=notint", nil)
    rr := httptest.NewRecorder()

    testApp.listRoleHandler(rr, req)

    if rr.Code != http.StatusUnprocessableEntity {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusUnprocessableEntity, rr.Code, rr.Body.String())
    }
}

func TestGetUserRolesHandler_InvalidID(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/v1/roles/user/", nil)
    rr := httptest.NewRecorder()

    testApp.getUserRolesHandler(rr, req)

    if rr.Code != http.StatusNotFound {
        t.Fatalf("expected status %d; got %d; body=%s", http.StatusNotFound, rr.Code, rr.Body.String())
    }
}