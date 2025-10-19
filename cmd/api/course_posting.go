// Filename: cmd/api/course_posting.go
package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kelseyaban/National-Inservice-Training-Database/internal/data"
	"github.com/kelseyaban/National-Inservice-Training-Database/internal/validator"
)

func (app *application) createCoursePostingHandler(w http.ResponseWriter, r *http.Request) {
	var incomingData struct {
		CourseID    int64 `json:"course_id"`
		PostingID   int64 `json:"posting_id"`
		Mandatory   bool  `json:"mandatory"`
		CreditHours int64 `json:"credithours"`
		RankID      int64 `json:"rank_id"`
	}

	err := app.readJSON(w, r, &incomingData)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	coursePosting := &data.CoursePosting{
		CourseID:    incomingData.CourseID,
		PostingID:   incomingData.PostingID,
		Mandatory:   incomingData.Mandatory,
		CreditHours: incomingData.CreditHours,
		RankID:      incomingData.RankID,
	}

	// Validate the course posting data
	v := validator.New()
	data.ValidateCoursePosting(v, coursePosting)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Insert the course posting into the database
	err = app.coursepostingModel.Insert(coursePosting)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	header := make(http.Header)
	header.Set("Location", fmt.Sprintf("/v1/course/postings/%d", coursePosting.ID))

	// Send a JSOn response with 201 (new resource created) status code
	data := envelope{
		"course_posting": coursePosting,
	}
	err = app.writeJSON(w, http.StatusCreated, data, header)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// Displays a course posting
func (app *application) displayCoursePostingHandler(w http.ResponseWriter, r *http.Request) {
	// get the id from the URL
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	coursePosting, err := app.coursepostingModel.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Send the course posting data as JSON response
	data := envelope{
		"course_posting": coursePosting,
	}

	err = app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// Edit course posting
func (app *application) updateCoursePostingHandler(w http.ResponseWriter, r *http.Request) {
	// get the id from the URL
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Fetch the existing course posting
	coursePosting, err := app.coursepostingModel.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var incomingData struct {
		CourseID    int64 `json:"course_id"`
		PostingID   int64 `json:"posting_id"`
		Mandatory   bool  `json:"mandatory"`
		CreditHours int64 `json:"credithours"`
		RankID      int64 `json:"rank_id"`
	}

	err = app.readJSON(w, r, &incomingData)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Check to see which field needs to be updated
	if incomingData.CourseID != 0 {
		coursePosting.CourseID = incomingData.CourseID
	}
	if incomingData.PostingID != 0 {
		coursePosting.PostingID = incomingData.PostingID
	}
	if !incomingData.Mandatory {
		coursePosting.Mandatory = incomingData.Mandatory
	}
	if incomingData.CreditHours != 0 {
		coursePosting.CreditHours = incomingData.CreditHours
	}
	if incomingData.RankID != 0 {
		coursePosting.RankID = incomingData.RankID
	}

	// Validate the updated course posting data
	v := validator.New()
	data.ValidateCoursePosting(v, coursePosting)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Update the course posting in the database
	err = app.coursepostingModel.Update(coursePosting)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Send the updated course posting as JSON response
	data := envelope{
		"course_posting": coursePosting,
	}

	err = app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// Delete a course posting
func (app *application) deleteCoursePostingHandler(w http.ResponseWriter, r *http.Request) {
	// get the id from the URL
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Delete the course posting from the database
	err = app.coursepostingModel.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// display the quote
	data := envelope{"message": "course posting successfully deleted"}

	err = app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// List all course postings
func (app *application) listCoursePostingsHandler(w http.ResponseWriter, r *http.Request) {
	var queryParametersData struct {
		CourseID    int64
		PostingID   int64
		Mandatory   bool
		CreditHours int64
		RankID      int64
		data.Filters
	}

	// get the query parameters from the URL
	queryParameters := r.URL.Query()

	v := validator.New()
	// Load the query parameters into our struct
	courseIDStr := app.getSingleQueryParameter(queryParameters, "course_id", "")
	if courseIDStr != "" {
		courseID, err := strconv.ParseInt(courseIDStr, 10, 64)
		if err != nil {
			v.AddError("course_id", "must be a valid integer")
		} else {
			queryParametersData.CourseID = courseID
		}
	}

	postingIDStr := app.getSingleQueryParameter(queryParameters, "posting_id", "")
	if postingIDStr != "" {
		postingID, err := strconv.ParseInt(postingIDStr, 10, 64)
		if err != nil {
			v.AddError("posting_id", "must be a valid integer")
		} else {
			queryParametersData.PostingID = postingID
		}
	}

	mandatoryStr := app.getSingleQueryParameter(queryParameters, "mandatory", "")
	if mandatoryStr != "" {
		mandatory, err := strconv.ParseBool(mandatoryStr)
		if err != nil {
			v.AddError("mandatory", "must be a valid boolean")
		} else {
			queryParametersData.Mandatory = mandatory
		}
	}

	creditHoursStr := app.getSingleQueryParameter(queryParameters, "credithours", "")
	if creditHoursStr != "" {
		creditHours, err := strconv.ParseInt(creditHoursStr, 10, 64)
		if err != nil {
			v.AddError("credithours", "must be a valid integer")
		} else {
			queryParametersData.CreditHours = creditHours
		}
	}

	rankIDStr := app.getSingleQueryParameter(queryParameters, "rank_id", "")
	if rankIDStr != "" {
		rankID, err := strconv.ParseInt(rankIDStr, 10, 64)
		if err != nil {
			v.AddError("rank_id", "must be a valid integer")
		} else {
			queryParametersData.RankID = rankID
		}
	}

	queryParametersData.Filters.Page = app.getSingleIntegerParameter(queryParameters, "page", 1, v)
	queryParametersData.Filters.PageSize = app.getSingleIntegerParameter(queryParameters, "page_size", 10, v)
	queryParametersData.Filters.Sort = app.getSingleQueryParameter(queryParameters, "sort", "id")
	queryParametersData.Filters.SortSafeList = []string{"id", "course_id", "posting_id", "mandatory", "credithours", "rank_id", "-id", "-course_id", "-posting_id", "-mandatory", "-credithours", "-rank_id"}

	// Check if the filters are valid
	data.ValidateFilters(v, queryParametersData.Filters)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Get the list of course postings from the database
	coursePostings, metadata, err := app.coursepostingModel.GetAll(
		queryParametersData.CourseID,
		queryParametersData.PostingID,
		queryParametersData.Mandatory,
		queryParametersData.CreditHours,
		queryParametersData.RankID,
		queryParametersData.Filters,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// send the course postings as JSON response
	data := envelope{
		"course_postings": coursePostings,
		"@metadata":       metadata,
	}
	err = app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
