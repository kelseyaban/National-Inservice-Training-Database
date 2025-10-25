// filename: cmd/api/course_posting_test.go
package main

import (
    // "bytes"
    // "context"
    // "encoding/json"
    // "errors"
    // "fmt"
    // "net/http"
    // "net/http/httptest"
    // "testing"
    // "time"

    // "github.com/kelseyaban/National-Inservice-Training-Database/internal/data"
)

// // --- Test Dependencies ---
// // Need change based on what is on the db
// const validPostingID int64 = 1
// const validRankID int64 = 1

// // tearDownCoursePostings truncates the course_posting table.
// func tearDownCoursePostings(t *testing.T) {
//     t.Helper()

//     ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//     defer cancel()

//     query := "TRUNCATE course_posting RESTART IDENTITY"
//     _, err := testApp.coursepostingModel.DB.ExecContext(ctx, query)
//     if err != nil {
//         t.Fatalf("Failed to truncate course_posting table: %v", err)
//     }
// }

// // createTestCoursePosting helper inserts a course_posting for testing GET/PUT/DELETE.
// // It relies on createTestCourse (from course_test.go) and the constant IDs.
// func createTestCoursePosting(t *testing.T, courseID int64) *data.CoursePosting {
//     t.Helper()

//     cp := &data.CoursePosting{
//         CourseID:    courseID,
//         PostingID:   validPostingID,
//         Mandatory:   true,
//         CreditHours: 10,
//         RankID:      validRankID,
//     }

//     err := testApp.coursepostingModel.Insert(cp)
//     if err != nil {
//         t.Fatalf("Failed to insert test course posting: %v", err)
//     }
//     return cp
// }

// // --- Tests ---

// // TestCreateCoursePostingHandler tests the POST /v1/course/postings endpoint
// func TestCreateCoursePostingHandler(t *testing.T) {
//     // Ensure tables are clean after this test function
//     t.Cleanup(func() {
//         tearDownCoursePostings(t)
//         // We must also clean up the courses we create
//         tearDownCourses(t)
//     })

//     // 1. Create the dependency: a valid Course
//     // We use the helper function from course_test.go
//     course := createTestCourse(t, "CP Test Course", "A course for CP testing")

//     testCases := []struct {
//         name           string
//         payload        map[string]interface{}
//         expectedStatus int
//     }{
//         {
//             name: "Success: Valid Course Posting",
//             payload: map[string]interface{}{
//                 "course_id":   course.ID,
//                 "posting_id":  validPostingID,
//                 "mandatory":   true,
//                 "credithours": 15,
//                 "rank_id":     validRankID,
//             },
//             expectedStatus: http.StatusCreated,
//         },
//         {
//             name: "Failure: Missing CourseID",
//             payload: map[string]interface{}{
//                 // "course_id" missing
//                 "posting_id":  validPostingID,
//                 "mandatory":   true,
//                 "credithours": 15,
//                 "rank_id":     validRankID,
//             },
//             expectedStatus: http.StatusUnprocessableEntity, // Caught by validation
//         },
//         {
//             name: "Failure: Missing PostingID",
//             payload: map[string]interface{}{
//                 "course_id": course.ID,
//                 // "posting_id" missing
//                 "mandatory":   true,
//                 "credithours": 15,
//                 "rank_id":     validRankID,
//             },
//             expectedStatus: http.StatusUnprocessableEntity, // Caught by validation
//         },
//         {
//             name: "Failure: Missing RankID",
//             payload: map[string]interface{}{
//                 "course_id":   course.ID,
//                 "posting_id":  validPostingID,
//                 "mandatory":   true,
//                 "credithours": 15,
//                 // "rank_id" missing
//             },
//             expectedStatus: http.StatusUnprocessableEntity, // Caught by validation
//         },
//         {
//             name: "Failure: Non-existent CourseID (FK Violation)",
//             payload: map[string]interface{}{
//                 "course_id":   999999, // Does not exist
//                 "posting_id":  validPostingID,
//                 "mandatory":   true,
//                 "credithours": 15,
//                 "rank_id":     validRankID,
//             },
//             expectedStatus: http.StatusInternalServerError, // DB foreign key constraint fail
//         },
//     }

//     for _, tc := range testCases {
//         t.Run(tc.name, func(t *testing.T) {
//             body, _ := json.Marshal(tc.payload)
//             // Note: Your handler seems to be at /v1/course/postings based on the create header
//             // If it's at /v1/course_postings, change this URL.
//             req := httptest.NewRequest(http.MethodPost, "/v1/course/postings", bytes.NewReader(body))
//             req.Header.Set("Content-Type", "application/json")

//             rr := executeRequest(req)

//             if rr.Code != tc.expectedStatus {
//                 t.Errorf("expected %d (%s); got %d (%s)\nBody: %s",
//                     tc.expectedStatus, http.StatusText(tc.expectedStatus),
//                     rr.Code, http.StatusText(rr.Code),
//                     rr.Body.String())
//             }
//         })
//     }
// }

// // TestDisplayCoursePostingHandler tests GET /v1/course/postings/:id
// func TestDisplayCoursePostingHandler(t *testing.T) {
//     t.Cleanup(func() {
//         tearDownCoursePostings(t)
//         tearDownCourses(t)
//     })

//     // 1. Create dependencies
//     course := createTestCourse(t, "Display Course", "Desc")
//     cp := createTestCoursePosting(t, course.ID)

//     // 2. Run tests
//     t.Run("Success: Valid ID", func(t *testing.T) {
//         url := fmt.Sprintf("/v1/course/postings/%d", cp.ID)
//         req := httptest.NewRequest(http.MethodGet, url, nil)
//         rr := executeRequest(req)

//         if rr.Code != http.StatusOK {
//             t.Fatalf("expected 200 OK; got %d", rr.Code)
//         }

//         var jsonResponse struct {
//             Posting data.CoursePosting `json:"course_posting"`
//         }
//         if err := json.NewDecoder(rr.Body).Decode(&jsonResponse); err != nil {
//             t.Fatalf("Failed to decode response: %v", err)
//         }

//         if jsonResponse.Posting.ID != cp.ID || jsonResponse.Posting.CreditHours != 10 {
//             t.Errorf("got unexpected data: %+v", jsonResponse.Posting)
//         }
//     })

//     t.Run("Failure: Non-existent ID", func(t *testing.T) {
//         url := "/v1/course/postings/999999"
//         req := httptest.NewRequest(http.MethodGet, url, nil)
//         rr := executeRequest(req)
//         if rr.Code != http.StatusNotFound {
//             t.Errorf("expected 404 Not Found; got %d", rr.Code)
//         }
//     })
// }

// // TestUpdateCoursePostingHandler tests PUT /v1/course/postings/:id
// func TestUpdateCoursePostingHandler(t *testing.T) {
//     t.Cleanup(func() {
//         tearDownCoursePostings(t)
//         tearDownCourses(t)
//     })

//     // 1. Create dependencies
//     course := createTestCourse(t, "Update Course", "Desc")
//     cp := createTestCoursePosting(t, course.ID) // Creates with credithours: 10

//     // 2. Run tests
//     t.Run("Success: Valid Update", func(t *testing.T) {
//         payload := map[string]interface{}{
//             "credithours": 25, // Change this
//         }
//         body, _ := json.Marshal(payload)
//         url := fmt.Sprintf("/v1/course/postings/%d", cp.ID)
//         req := httptest.NewRequest(http.MethodPut, url, bytes.NewReader(body))
//         req.Header.Set("Content-Type", "application/json")

//         rr := executeRequest(req)

//         if rr.Code != http.StatusOK {
//             t.Fatalf("expected 200 OK; got %d\nBody: %s", rr.Code, rr.Body.String())
//         }

//         var jsonResponse struct {
//             Posting data.CoursePosting `json:"course_posting"`
//         }
//         if err := json.NewDecoder(rr.Body).Decode(&jsonResponse); err != nil {
//             t.Fatalf("Failed to decode response: %v", err)
//         }

//         // Check that the field updated
//         if jsonResponse.Posting.CreditHours != 25 {
//             t.Errorf("expected credithours to be 25, got %d", jsonResponse.Posting.CreditHours)
//         }
//         // Check that other fields remained the same
//         if jsonResponse.Posting.CourseID != cp.CourseID {
//             t.Errorf("expected course_id to be %d, got %d", cp.CourseID, jsonResponse.Posting.CourseID)
//         }
//     })

//     t.Run("Failure: Invalid Update (validation fail)", func(t *testing.T) {
//         payload := map[string]interface{}{
//             "rank_id": 0, // Invalid (must be > 0)
//         }
//         body, _ := json.Marshal(payload)
//         url := fmt.Sprintf("/v1/course/postings/%d", cp.ID)
//         req := httptest.NewRequest(http.MethodPut, url, bytes.NewReader(body))
//         req.Header.Set("Content-Type", "application/json")

//         rr := executeRequest(req)

//         if rr.Code != http.StatusUnprocessableEntity {
//             t.Errorf("expected 422 Unprocessable Entity; got %d", rr.Code)
//         }
//     })
// }

// // TestDeleteCoursePostingHandler tests DELETE /v1/course/postings/:id
// func TestDeleteCoursePostingHandler(t *testing.T) {
//     t.Cleanup(func() {
//         tearDownCoursePostings(t)
//         tearDownCourses(t)
//     })

//     // 1. Create dependencies
//     course := createTestCourse(t, "Delete Course", "Desc")
//     cp := createTestCoursePosting(t, course.ID)

//     // 2. Run tests
//     t.Run("Success: Valid ID", func(t *testing.T) {
//         url := fmt.Sprintf("/v1/course/postings/%d", cp.ID)
//         req := httptest.NewRequest(http.MethodDelete, url, nil)
//         rr := executeRequest(req)

//         if rr.Code != http.StatusOK {
//             t.Fatalf("expected 200 OK; got %d", rr.Code)
//         }

//         // Verify it's gone
//         _, err := testApp.coursepostingModel.Get(cp.ID)
//         if !errors.Is(err, data.ErrRecordNotFound) {
//             t.Error("expected course posting to be deleted, but it still exists")
//         }
//     })

//     t.Run("Failure: Non-existent ID", func(t *testing.T) {
//         url := "/v1/course/postings/999999"
//         req := httptest.NewRequest(http.MethodDelete, url, nil)
//         rr := executeRequest(req)
//         if rr.Code != http.StatusNotFound {
//             t.Errorf("expected 404 Not Found; got %d", rr.Code)
//         }
//     })
// }

// // TestListCoursePostingsHandler tests GET /v1/course/postings
// func TestListCoursePostingsHandler(t *testing.T) {
//     t.Cleanup(func() {
//         tearDownCoursePostings(t)
//         tearDownCourses(t)
//     })

//     // 1. Create dependencies
//     course1 := createTestCourse(t, "C1", "Desc1")
//     course2 := createTestCourse(t, "C2", "Desc2")

//     // Create 3 postings
//     // CP1: Course 1, Rank 1, 10 credits
//     createTestCoursePosting(t, course1.ID)
//     // CP2: Course 2, Rank 1, 20 credits
//     cp2 := &data.CoursePosting{CourseID: course2.ID, PostingID: validPostingID, Mandatory: false, CreditHours: 20, RankID: validRankID}
//     testApp.coursepostingModel.Insert(cp2)
//     // CP3: Course 1, Rank 2 (assuming ID 2 exists), 10 credits
//     cp3 := &data.CoursePosting{CourseID: course1.ID, PostingID: validPostingID, Mandatory: true, CreditHours: 10, RankID: 2} // Assuming rank 2 exists
//     testApp.coursepostingModel.Insert(cp3)

//     type listResponse struct {
//         Postings []data.CoursePosting `json:"course_postings"`
//         Metadata data.Metadata        `json:"@metadata"`
//     }

//     // 2. Run tests
//     t.Run("List all (no filters)", func(t *testing.T) {
//         req := httptest.NewRequest(http.MethodGet, "/v1/course/postings", nil)
//         rr := executeRequest(req)
//         if rr.Code != http.StatusOK {
//             t.Fatalf("expected 200 OK; got %d", rr.Code)
//         }

//         var resp listResponse
//         if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
//             t.Fatalf("Failed to decode JSON: %v", err)
//         }

//         if resp.Metadata.TotalRecords != 3 {
//             t.Errorf("expected 3 total records; got %d", resp.Metadata.TotalRecords)
//         }
//     })

//     t.Run("Filter by course_id", func(t *testing.T) {
//         url := fmt.Sprintf("/v1/course/postings?course_id=%d", course1.ID)
//         req := httptest.NewRequest(http.MethodGet, url, nil)
//         rr := executeRequest(req)
//         if rr.Code != http.StatusOK {
//             t.Fatalf("expected 200 OK; got %d", rr.Code)
//         }

//         var resp listResponse
//         if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
//             t.Fatalf("Failed to decode JSON: %v", err)
//         }

//         if resp.Metadata.TotalRecords != 2 {
//             t.Errorf("expected 2 total records; got %d", resp.Metadata.TotalRecords)
//         }
//     })

//     t.Run("Filter by credithours", func(t *testing.T) {
//         req := httptest.NewRequest(http.MethodGet, "/v1/course/postings?credithours=10", nil)
//         rr := executeRequest(req)
//         if rr.Code != http.StatusOK {
//             t.Fatalf("expected 200 OK; got %d", rr.Code)
//         }

//         var resp listResponse
//         if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
//             t.Fatalf("Failed to decode JSON: %v", err)
//         }

//         if resp.Metadata.TotalRecords != 2 {
//             t.Errorf("expected 2 total records; got %d", resp.Metadata.TotalRecords)
//         }
//     })

//     t.Run("Filter by rank_id", func(t *testing.T) {
//         req := httptest.NewRequest(http.MethodGet, "/v1/course/postings?rank_id=2", nil)
//         rr := executeRequest(req)
//         if rr.Code != http.StatusOK {
//             t.Fatalf("expected 200 OK; got %d", rr.Code)
//         }

//         var resp listResponse
//         if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
//             t.Fatalf("Failed to decode JSON: %v", err)
//         }

//         if resp.Metadata.TotalRecords != 1 {
//             t.Errorf("expected 1 total record; got %d", resp.Metadata.TotalRecords)
//         }
//     })

//     t.Run("Failure: Invalid Filter Value", func(t *testing.T) {
//         req := httptest.NewRequest(http.MethodGet, "/v1/course/postings?course_id=abc", nil)
//         rr := executeRequest(req)
//         if rr.Code != http.StatusUnprocessableEntity {
//             t.Errorf("expected 422 Unprocessable Entity; got %d", rr.Code)
//         }
//     })
// }
