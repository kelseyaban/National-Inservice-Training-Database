// Filename: cmd/api/facilitator_rating_test.go
package main

import (
    // "bytes"
    // "context"
    // "encoding/json"
    // "fmt"
    // "net/http"
    // "net/http/httptest"
    // "testing"
    // "time"

    // "github.com/kelseyaban/National-Inservice-Training-Database/internal/data"
)

// // createTestUserHelper inserts a minimal, activated user for testing FK constraints.
// // It returns the new user's ID and registers a cleanup function to delete them.
// func createTestUserHelper(t *testing.T) int64 {
//     t.Helper()

//     // Create a user with unique credentials
//     email := fmt.Sprintf("rating-user-%d@test.com", time.Now().UnixNano())
//     username := fmt.Sprintf("rating-user-%d", time.Now().UnixNano())

//     user := &data.User{
//         RegulationNumber: "R-RATE-TEST",
//         Username:         username,
//         FName:            "Rating",
//         LName:            "TestUser",
//         Email:            email,
//         Gender:           "N/A",
//         Formation:        1,
//         Rank:             1,
//         Postings:         1,
//         Activated:        true, // Activate them so we don't need token logic
//     }

//     // This call will now work on the zero-initialized 'Password' field.
//     if err := user.Password.Set("SafeP@ss123"); err != nil {
//         t.Fatalf("Failed to set/hash password for test user: %v", err)
//     }

//     // Use the model's Insert method to save the user (with the new hash)
//     err := testApp.userModel.Insert(user)
//     if err != nil {
//         t.Fatalf("Failed to create test user for rating: %v", err)
//     }

//     // Register the main user cleanup function (from users_test.go)
//     t.Cleanup(func() {
//         tearDownTestData(t, []string{email})
//     })

//     return user.ID
// }

// // tearDownFacilitatorRatings cleans up ratings created for a specific user.
// func tearDownFacilitatorRatings(t *testing.T, userID int64) {
//     t.Helper()
//     if userID == 0 {
//         return
//     }

//     ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//     defer cancel()

//     // We can use any model's DB connection
//     query := "DELETE FROM facilitator_rating WHERE user_id = $1"
//     _, err := testApp.userModel.DB.ExecContext(ctx, query, userID)
//     if err != nil {
//         t.Errorf("facilitator rating cleanup failed: %v", err)
//     }
// }

// // TestAddFacilitatorRating tests the POST /v1/facilitator_rating endpoint
// func TestAddFacilitatorRating(t *testing.T) {
//     // 1. Create a prerequisite test user. All ratings will be linked to this user.
//     testUserID := createTestUserHelper(t)
//     if testUserID == 0 {
//         t.Fatal("Test user ID is zero, setup failed")
//     }

//     // 2. Schedule cleanup for any ratings created for this user.
//     // This runs before the user cleanup (t.Cleanup is LIFO).
//     t.Cleanup(func() {
//         tearDownFacilitatorRatings(t, testUserID)
//     })

//     testCases := []struct {
//         name           string
//         payload        map[string]interface{}
//         expectedStatus int
//     }{
//         {
//             name: "Success: Valid Rating",
//             payload: map[string]interface{}{
//                 "user_id": testUserID,
//                 "rating":  4,
//             },
//             expectedStatus: http.StatusCreated,
//         },
//         {
//             name: "Failure: Rating Too High",
//             payload: map[string]interface{}{
//                 "user_id": testUserID,
//                 "rating":  6, // Validation rule is <= 5
//             },
//             expectedStatus: http.StatusUnprocessableEntity,
//         },
//         {
//             name: "Failure: Rating Too Low",
//             payload: map[string]interface{}{
//                 "user_id": testUserID,
//                 "rating":  0, // Validation rule is >= 1
//             },
//             expectedStatus: http.StatusUnprocessableEntity,
//         },
//         {
//             name: "Failure: Missing UserID",
//             payload: map[string]interface{}{
//                 // "user_id" is missing
//                 "rating": 3,
//             },
//             expectedStatus: http.StatusUnprocessableEntity, // Caught by v.Check(fr.UserID > 0, ...)
//         },
//         {
//             name: "Failure: Missing Rating",
//             payload: map[string]interface{}{
//                 "user_id": testUserID,
//                 // "rating" is missing (will be 0)
//             },
//             expectedStatus: http.StatusUnprocessableEntity, // Caught by v.Check(fr.Rating >= 1, ...)
//         },
//         {
//             name: "Failure: Non-existent UserID (FK Constraint)",
//             payload: map[string]interface{}{
//                 "user_id": 999999999, // A valid int, but not a real user
//                 "rating":  5,
//             },
//             // This fails at the database (foreign key violation), which the
//             // handler catches and returns a 500 Server Error.
//             expectedStatus: http.StatusInternalServerError,
//         },
//     }

//     for _, tc := range testCases {
//         t.Run(tc.name, func(t *testing.T) {
//             body, err := json.Marshal(tc.payload)
//             if err != nil {
//                 t.Fatalf("Failed to marshal payload: %v", err)
//             }

//             // Note: The endpoint path is /v1/facilitator_rating
//             req := httptest.NewRequest(http.MethodPost, "/v1/facilitator_rating", bytes.NewReader(body))
//             req.Header.Set("Content-Type", "application/json")

//             // executeRequest is available from users_test.go
//             rr := executeRequest(req)

//             if rr.Code != tc.expectedStatus {
//                 t.Errorf("expected %d (%s); got %d (%s)\nBody: %s",
//                     tc.expectedStatus, http.StatusText(tc.expectedStatus),
//                     rr.Code, http.StatusText(rr.Code),
//                     rr.Body.String())
//             }

//             // For the successful case, check the Location header
//             if tc.expectedStatus == http.StatusCreated {
//                 location := rr.Header().Get("Location")
//                 if location == "" {
//                     t.Errorf("expected 'Location' header for 201 Created, but it was missing")
//                 }
//             }
//         })
//     }
// }

// func TestDisplayFacilitatorRatingHandler(t *testing.T) {
//     // 1. Create a test user
//     testUserID := createTestUserHelper(t)

//     // 2. Schedule cleanup for any ratings created
//     t.Cleanup(func() {
//         tearDownFacilitatorRatings(t, testUserID)
//     })

//     // 3. Create a test rating using the model directly
//     rating := &data.FacilitatorRating{
//         UserID: testUserID,
//         Rating: 5,
//     }
//     err := testApp.facilitatorRatingModel.Insert(rating)
//     if err != nil {
//         t.Fatalf("Failed to insert test rating: %v", err)
//     }

//     // 4. Define test cases
//     testCases := []struct {
//         name           string
//         ratingID       string // Use string to simulate URL param
//         expectedStatus int
//     }{
//         {
//             name:           "Success: Valid ID",
//             ratingID:       fmt.Sprintf("%d", rating.ID),
//             expectedStatus: http.StatusOK,
//         },
//         {
//             name:           "Failure: Non-existent ID",
//             ratingID:       "999999",
//             expectedStatus: http.StatusNotFound,
//         },
//         {
//             name:           "Failure: Invalid ID format",
//             ratingID:       "abc", // readIDParam helper should catch this
//             expectedStatus: http.StatusNotFound,
//         },
//     }

//     for _, tc := range testCases {
//         t.Run(tc.name, func(t *testing.T) {
//             url := fmt.Sprintf("/v1/facilitator_rating/%s", tc.ratingID)
//             req := httptest.NewRequest(http.MethodGet, url, nil)

//             rr := executeRequest(req)

//             if rr.Code != tc.expectedStatus {
//                 t.Errorf("expected %d (%s); got %d (%s)\nBody: %s",
//                     tc.expectedStatus, http.StatusText(tc.expectedStatus),
//                     rr.Code, http.StatusText(rr.Code),
//                     rr.Body.String())
//             }

//             // If we got a 200 OK, check the body content
//             if tc.expectedStatus == http.StatusOK {
//                 var jsonResponse struct {
//                     Rating data.FacilitatorRating `json:"facilitator_rating"`
//                 }
//                 err := json.NewDecoder(rr.Body).Decode(&jsonResponse)
//                 if err != nil {
//                     t.Fatalf("Failed to decode JSON response: %v", err)
//                 }

//                 if jsonResponse.Rating.ID != rating.ID {
//                     t.Errorf("expected rating ID %d; got %d", rating.ID, jsonResponse.Rating.ID)
//                 }
//                 if jsonResponse.Rating.Rating != 5 {
//                     t.Errorf("expected rating value 5; got %d", jsonResponse.Rating.Rating)
//                 }
//             }
//         })
//     }
// }