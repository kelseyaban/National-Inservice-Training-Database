package main

import (
    "fmt"
    "net/http"
    "errors"

    "github.com/kelseyaban/National-Inservice-Training-Database/internal/data"
    "github.com/kelseyaban/National-Inservice-Training-Database/internal/validator"
)

// ---------------- CREATE ----------------
func (a *application) createUserSessionHandler(w http.ResponseWriter, r *http.Request) {
    var input struct {
        TraineeID            int64  `json:"trainee_id"`
        SessionID            int64  `json:"session_id"`
        CreditHoursCompleted int64  `json:"credithours_completed"`
        Grade                string `json:"grade"`
        Feedback             string `json:"feedback"`
    }

    err := a.readJSON(w, r, &input)
    if err != nil {
        a.badRequestResponse(w, r, err)
        return
    }

    us := &data.UserSession{
        TraineeID:               input.TraineeID,
        SessionID:            input.SessionID,
        CreditHoursCompleted: input.CreditHoursCompleted,
        Grade:                input.Grade,
        Feedback:             input.Feedback,
    }

    // Validate input
    v := validator.New()
    data.ValidateUserSession(v, us)
    if !v.IsEmpty() {
        a.failedValidationResponse(w, r, v.Errors)
        return
    }

    // Insert record into DB
    err = a.userSessionModel.AddUserSession(us)
    if err != nil {
        a.serverErrorResponse(w, r, err)
        return
    }

    headers := make(http.Header)
    headers.Set("Location", fmt.Sprintf("/v1/usersessions/%d", us.ID))

    data := envelope{
        "user_session": us,
    }

    err = a.writeJSON(w, http.StatusCreated, data, headers)
    if err != nil {
        a.serverErrorResponse(w, r, err)
    }
}

// ---------------- GET ----------------
func (a *application) getUserSessionHandler(w http.ResponseWriter, r *http.Request) {
    id, err := a.readIDParam(r)
    if err != nil {
        a.notFoundResponse(w, r)
        return
    }

    us, err := a.userSessionModel.GetUserSession(id)
    if err != nil {
        if errors.Is(err, data.ErrRecordNotFound) {
            a.notFoundResponse(w, r)
        } else {
            a.serverErrorResponse(w, r, err)
        }
        return
    }

    data := envelope{
        "user_session": us,
    }

    err = a.writeJSON(w, http.StatusOK, data, nil)
    if err != nil {
        a.serverErrorResponse(w, r, err)
    }
}

// ---------------- UPDATE ----------------
func (a *application) updateUserSessionHandler(w http.ResponseWriter, r *http.Request) {
    id, err := a.readIDParam(r)
    if err != nil {
        a.notFoundResponse(w, r)
        return
    }

    us, err := a.userSessionModel.GetUserSession(id)
    if err != nil {
        if errors.Is(err, data.ErrRecordNotFound) {
            a.notFoundResponse(w, r)
        } else {
            a.serverErrorResponse(w, r, err)
        }
        return
    }

    var input struct {
        CreditHoursCompleted *int64  `json:"credithours_completed"`
        Grade                *string `json:"grade"`
        Feedback             *string `json:"feedback"`
    }

    err = a.readJSON(w, r, &input)
    if err != nil {
        a.badRequestResponse(w, r, err)
        return
    }

    // Apply updates if fields are provided
    if input.CreditHoursCompleted != nil {
        us.CreditHoursCompleted = *input.CreditHoursCompleted
    }
    if input.Grade != nil {
        us.Grade = *input.Grade
    }
    if input.Feedback != nil {
        us.Feedback = *input.Feedback
    }

    v := validator.New()
    data.ValidateUserSession(v, us)
    if !v.IsEmpty() {
        a.failedValidationResponse(w, r, v.Errors)
        return
    }

    err = a.userSessionModel.UpdateUserSession(us)
    if err != nil {
        a.serverErrorResponse(w, r, err)
        return
    }

    data := envelope{
        "user_session": us,
    }

    err = a.writeJSON(w, http.StatusOK, data, nil)
    if err != nil {
        a.serverErrorResponse(w, r, err)
    }
}

// ---------------- DELETE ----------------
func (a *application) deleteUserSessionHandler(w http.ResponseWriter, r *http.Request) {
    id, err := a.readIDParam(r)
    if err != nil {
        a.notFoundResponse(w, r)
        return
    }

    err = a.userSessionModel.DeleteUserSession(id)
    if err != nil {
        if errors.Is(err, data.ErrRecordNotFound) {
            a.notFoundResponse(w, r)
        } else {
            a.serverErrorResponse(w, r, err)
        }
        return
    }

    data := envelope{
        "message": fmt.Sprintf("user session %d successfully deleted", id),
    }

    err = a.writeJSON(w, http.StatusOK, data, nil)
    if err != nil {
        a.serverErrorResponse(w, r, err)
    }
}

// ---------------- LIST ----------------
func (a *application) listUserSessionHandler(w http.ResponseWriter, r *http.Request) {
    sessions, err := a.userSessionModel.GetAllUserSessions()
    if err != nil {
        a.serverErrorResponse(w, r, err)
        return
    }

    data := envelope{
        "user_session": sessions,
    }

    err = a.writeJSON(w, http.StatusOK, data, nil)
    if err != nil {
        a.serverErrorResponse(w, r, err)
    }
}