package main

import (
	//   "encoding/json"
	"fmt"
	"net/http"
	"errors"
	"github.com/kelseyaban/National-Inservice-Training-Database/internal/data"
	"github.com/kelseyaban/National-Inservice-Training-Database/internal/validator"
)

func (a *application) addFacilitatorRating(w http.ResponseWriter, r *http.Request) { 
    var incomingData struct {
        UserID int   `json:"user_id"`
        Rating int   `json:"rating"`
    }
    
    // Decode JSON body
    err := a.readJSON(w, r, &incomingData)
    if err != nil {
        a.badRequestResponse(w, r, err)
        return
    }

    fr := &data.FacilitatorRating{
        UserID: int64(incomingData.UserID),
        Rating: incomingData.Rating,
    }

    // Initialize a validator instance
    v := validator.New()
    data.ValidateFacilitatorRating(v, fr)
    if !v.IsEmpty() {
        a.failedValidationResponse(w, r, v.Errors)
        return
    }

    // Insert into database
    err = a.facilitatorRatingModel.Insert(fr)
    if err != nil {
        a.serverErrorResponse(w, r, err)
        return
    }
 
    // Set Location header
    headers := make(http.Header)
    headers.Set("Location", fmt.Sprintf("/v1/facilitator_rating/%d", fr.ID))

    // Send JSON response with 201 status
    data := envelope{
        "facilitator_rating": fr,
    }
    err = a.writeJSON(w, http.StatusCreated, data, headers)
    if err != nil {
        a.serverErrorResponse(w, r, err)
    }
}
func (a *application) displayFacilitatorRatingHandler(w http.ResponseWriter, r *http.Request) {
    id, err := a.readIDParam(r)
    if err != nil {
        a.notFoundResponse(w, r)
        return 
    }
 
    fr, err := a.facilitatorRatingModel.Get(id)
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
        "facilitator_rating": fr,
    }
    err = a.writeJSON(w, http.StatusOK, data, nil)
    if err != nil {
        a.serverErrorResponse(w, r, err)
    }
}
func (a *application) listFacilitatorRatingHandler(w http.ResponseWriter, r *http.Request) {
    var queryData struct {
        UserID int
        data.Filters
    }

    q := r.URL.Query()
    queryData.UserID = a.getSingleIntegerParameter(q, "user_id", 0, validator.New()) // 0 = all users

    // Setup pagination & sorting
    v := validator.New()
    queryData.Filters.Page = a.getSingleIntegerParameter(q, "page", 1, v)
    queryData.Filters.PageSize = a.getSingleIntegerParameter(q, "page_size", 10, v)
    queryData.Filters.Sort = a.getSingleQueryParameter(q, "sort", "id")
    queryData.Filters.SortSafeList = []string{"id", "user_id", "rating", "-id", "-user_id", "-rating"}

    data.ValidateFilters(v, queryData.Filters)
    if !v.IsEmpty() {
        a.failedValidationResponse(w, r, v.Errors)
        return
    }

    ratings, metadata, err := a.facilitatorRatingModel.GetAll(int64(queryData.UserID), queryData.Filters)
    if err != nil {
        a.serverErrorResponse(w, r, err)
        return
    }

    data := envelope{
        "facilitator_rating": ratings,
        "@metadata": metadata,
    }
    err = a.writeJSON(w, http.StatusOK, data, nil)
    if err != nil {
        a.serverErrorResponse(w, r, err)
    }
}
