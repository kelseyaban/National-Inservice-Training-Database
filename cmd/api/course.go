// Filename: cmd/api/course.go
package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/kelseyaban/National-Inservice-Training-Database/internal/data"
	"github.com/kelseyaban/National-Inservice-Training-Database/internal/validator"
)

func (app *application) createCourseHandler(w http.ResponseWriter, r *http.Request) {
	var incomingData struct {
		Course_Name string `json:"course"`
		Description string `json:"description"`
	}

	err := app.readJSON(w, r, &incomingData)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	course := &data.Course{
		Course_Name: incomingData.Course_Name,
		Description: incomingData.Description,
	}

	// Validate the course data
	v := validator.New()
	data.ValidateCourse(v, course)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Insert the course into the database
	err = app.courseModel.Insert(course)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/courses/%d", course.ID))

	// Send a JSOn response with 201 (new resource created) status code
	data := envelope{
		"course": course,
	}

	err = app.writeJSON(w, http.StatusCreated, data, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// Displays quotes
func (app *application) displayCourseHandler(w http.ResponseWriter, r *http.Request) {
	// get the id from the URL
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	course, err := app.courseModel.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// send the course as JSON response
	data := envelope{
		"course": course,
	}

	err = app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// Edit course
func (app *application) updateCourseHandler(w http.ResponseWriter, r *http.Request) {
	// get the id from the URL
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	course, err := app.courseModel.Get(id)
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
		Course_Name string `json:"course"`
		Description string `json:"description"`
	}

	err = app.readJSON(w, r, &incomingData)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// check to see which fields need to be updated
	if incomingData.Description != "" {
		course.Description = incomingData.Description
	}
	if incomingData.Course_Name != "" {
		course.Course_Name = incomingData.Course_Name
	}

	// validate the updated course data
	v := validator.New()
	data.ValidateCourse(v, course)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// update the course in the database
	err = app.courseModel.Update(course)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	data := envelope{
		"course": course,
	}

	// send the updated course as JSON response
	err = app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// delete course
func (app *application) deleteCourseHandler(w http.ResponseWriter, r *http.Request) {
	// get the id from the URL
	id, err := app.readIDParam(r)

	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.courseModel.Delete(id)
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
	data := envelope{"message": "course successfully deleted"}

	err = app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// List all Courses
func (app *application) listCoursesHandler(w http.ResponseWriter, r *http.Request) {
	var queryParametersData struct {
		Course_Name string
		Description string
		data.Filters
	}

	// get the query parameters from the URL
	queryParameters := r.URL.Query()

	// Load the query parameters into our struct
	queryParametersData.Course_Name = app.getSingleQueryParameter(queryParameters, "course", "")
	queryParametersData.Description = app.getSingleQueryParameter(queryParameters, "description", "")

	// validation
	v := validator.New()
	queryParametersData.Filters.Page = app.getSingleIntegerParameter(queryParameters, "page", 1, v)
	queryParametersData.Filters.PageSize = app.getSingleIntegerParameter(queryParameters, "page_size", 10, v)
	queryParametersData.Sort = app.getSingleQueryParameter(queryParameters, "sort", "id")
	queryParametersData.Filters.SortSafeList = []string{"id", "course", "-id", "-course"}

	// Check if the filters are valid
	data.ValidateFilters(v, queryParametersData.Filters)
	if !v.IsEmpty() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// get the list of courses from the database
	courses, metadata, err := app.courseModel.GetAll(queryParametersData.Course_Name, queryParametersData.Description, queryParametersData.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// send the courses as JSON response
	data := envelope{
		"courses":   courses,
		"@metadata": metadata,
	}
	err = app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
