// Filename: cmd/api/attendance.go
package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/kelseyaban/National-Inservice-Training-Database/internal/data"
	"github.com/kelseyaban/National-Inservice-Training-Database/internal/validator"
)

func (app *application) createAttendanceHandler(w http.ResponseWriter, r *http.Request) {
	var incomingData struct {
		UserSessionID    int64  `json:"user_session_id"`
		AttendanceStatus bool   `json:"attendance"`
		Date             string `json:"date"`
	}

	err := app.readJSON(w, r, &incomingData)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	attendance := &data.Attendance{
		UserSessionID:    incomingData.UserSessionID,
		AttendanceStatus: incomingData.AttendanceStatus,
		Date:             parseDate(incomingData.Date),
	}

	// Validate the attendance data
	v := validator.New()
	data.ValidateAttendance(v, attendance)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Insert the attendance into the database
	err = app.attendanceModel.Insert(attendance)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/attendance/%d", attendance.ID))
	// Send a JSOn response with 201 (new resource created) status code
	data := envelope{
		"attendance": attendance,
	}

	err = app.writeJSON(w, http.StatusCreated, data, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// Displays attendance record of a specific user
func (app *application) displayIndividualAttendanceHandler(w http.ResponseWriter, r *http.Request) {
	// get the id from the URL
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	attendance, err := app.attendanceModel.GetIdividualAttendance(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// send the attendance record as JSON response
	data := envelope{
		"attendance": attendance,
	}

	err = app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// Updates attendance record of a specific user
func (app *application) updateAttendanceHandler(w http.ResponseWriter, r *http.Request) {
	// Get the id from the URL
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	attendance, err := app.attendanceModel.GetIdividualAttendance(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Read the incoming JSON
	var incomingData struct {
		AttendanceStatus *bool  `json:"attendance"`
		Date             string `json:"date"`
	}

	err = app.readJSON(w, r, &incomingData)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Update fields only if provided
	if incomingData.AttendanceStatus != nil {
		attendance.AttendanceStatus = *incomingData.AttendanceStatus
	}

	if incomingData.Date != "" {
		attendance.Date = parseDate(incomingData.Date)
	}

	// Validate the updated attendance data
	v := validator.New()
	data.ValidateAttendance(v, attendance)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Update the record in the DB
	err = app.attendanceModel.Update(attendance)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Send a JSON response with the updated record
	data := envelope{
		"attendance": attendance,
	}

	err = app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
