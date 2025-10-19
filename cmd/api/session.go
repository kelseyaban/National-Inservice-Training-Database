package main

import (
    "fmt"
    "net/http"
    "errors"

    "github.com/kelseyaban/National-Inservice-Training-Database/internal/data"
    "github.com/kelseyaban/National-Inservice-Training-Database/internal/validator"
)

//------------------ CREATE ------------------
func (a *application) createSessionHandler(w http.ResponseWriter, r *http.Request) {
    var incomingData struct {
        CourseID      int64 `json:"course_id"`
        FormationID   int64 `json:"formation_id"`
        FacilitatorID int64 `json:"facilitator_id"`
    }

    err := a.readJSON(w, r, &incomingData)
    if err != nil {
        a.badRequestResponse(w, r, err)
        return
    }

    session := &data.Session{
        CourseID:      incomingData.CourseID,
        FormationID:   incomingData.FormationID,
        FacilitatorID: incomingData.FacilitatorID,
    }

    v := validator.New()
    data.ValidateSession(v, session)
    if !v.IsEmpty() {
        a.failedValidationResponse(w, r, v.Errors)
        return
    }

    err = a.sessionModel.Insert(session)
    if err != nil {
        a.serverErrorResponse(w, r, err)
        return
    }

    headers := make(http.Header)
    headers.Set("Location", fmt.Sprintf("/v1/session/%d", session.ID))

    data := envelope{
        "session": session,
    }

    err = a.writeJSON(w, http.StatusCreated, data, headers)
    if err != nil {
        a.serverErrorResponse(w, r, err)
    }
}

//------------------ DISPLAY ------------------
func (a *application) displaySessionHandler(w http.ResponseWriter, r *http.Request) {
    id, err := a.readIDParam(r)
    if err != nil {
        a.notFoundResponse(w, r)
        return
    }

    session, err := a.sessionModel.Get(id)
    if err != nil {
        switch {
        case errors.Is(err, data.ErrRecordNotFound):
            a.notFoundResponse(w, r)
        default:
            a.serverErrorResponse(w, r, err)
        }
        return
    }

    data := envelope{
        "session": session,
    }

    err = a.writeJSON(w, http.StatusOK, data, nil)
    if err != nil {
        a.serverErrorResponse(w, r, err)
    }
}

//------------------ UPDATE ------------------
func (a *application) updateSessionHandler(w http.ResponseWriter, r *http.Request) {
    id, err := a.readIDParam(r)
    if err != nil {
        a.notFoundResponse(w, r)
        return
    }

    session, err := a.sessionModel.Get(id)
    if err != nil {
        switch {
        case errors.Is(err, data.ErrRecordNotFound):
            a.notFoundResponse(w, r)
        default:
            a.serverErrorResponse(w, r, err)
        }
        return
    }

    var incomingData struct {
        CourseID      *int64 `json:"course_id"`
        FormationID   *int64 `json:"formation_id"`
        FacilitatorID *int64 `json:"facilitator_id"`
    }

    err = a.readJSON(w, r, &incomingData)
    if err != nil {
        a.badRequestResponse(w, r, err)
        return
    }

    if incomingData.CourseID != nil {
        session.CourseID = *incomingData.CourseID
    }
    if incomingData.FormationID != nil {
        session.FormationID = *incomingData.FormationID
    }
    if incomingData.FacilitatorID != nil {
        session.FacilitatorID = *incomingData.FacilitatorID
    }

    v := validator.New()
    data.ValidateSession(v, session)
    if !v.IsEmpty() {
        a.failedValidationResponse(w, r, v.Errors)
        return
    }

    err = a.sessionModel.Update(session)
    if err != nil {
        a.serverErrorResponse(w, r, err)
        return
    }

    data := envelope{
        "session": session,
    }

    err = a.writeJSON(w, http.StatusOK, data, nil)
    if err != nil {
        a.serverErrorResponse(w, r, err)
    }
}

//------------------ DELETE ------------------
func (a *application) deleteSessionHandler(w http.ResponseWriter, r *http.Request) {
    id, err := a.readIDParam(r)
    if err != nil {
        a.notFoundResponse(w, r)
        return
    }

    err = a.sessionModel.Delete(id)
    if err != nil {
        switch {
        case errors.Is(err, data.ErrRecordNotFound):
            a.notFoundResponse(w, r)
        default:
            a.serverErrorResponse(w, r, err)
        }
        return
    }

    data := envelope{
        "message": "session successfully deleted",
    }

    err = a.writeJSON(w, http.StatusOK, data, nil)
    if err != nil {
        a.serverErrorResponse(w, r, err)
    }
}

//------------------ LIST ------------------
func (a *application) listSessionHandler(w http.ResponseWriter, r *http.Request) {
    var queryParametersData struct {
        data.Filters
    }

    queryParameters := r.URL.Query()
    v := validator.New()

    queryParametersData.Filters.Page = a.getSingleIntegerParameter(queryParameters, "page", 1, v)
    queryParametersData.Filters.PageSize = a.getSingleIntegerParameter(queryParameters, "page_size", 10, v)
    queryParametersData.Filters.Sort = a.getSingleQueryParameter(queryParameters, "sort", "id")
    queryParametersData.Filters.SortSafeList = []string{"id", "-id"}

    data.ValidateFilters(v, queryParametersData.Filters)
    if !v.IsEmpty() {
        a.failedValidationResponse(w, r, v.Errors)
        return
    }

    sessions, metadata, err := a.sessionModel.GetAll(queryParametersData.Filters)
    if err != nil {
        a.serverErrorResponse(w, r, err)
        return
    }

    data := envelope{
        "session":   sessions,
        "@metadata": metadata,
    }

    err = a.writeJSON(w, http.StatusOK, data, nil)
    if err != nil {
        a.serverErrorResponse(w, r, err)
    }
}