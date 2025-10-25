// filename: cmd/api/course_test.go
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

// // tearDownCourses truncates the course table to ensure a clean state.
// func tearDownCourses(t *testing.T) {
//     t.Helper()

//     ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//     defer cancel()

//     // TRUNCATE is faster than DELETE and resets any auto-incrementing counters.
//     query := "TRUNCATE course RESTART IDENTITY"
//     _, err := testApp.courseModel.DB.ExecContext(ctx, query)
//     if err != nil {
//         t.Fatalf("Failed to truncate course table: %v", err)
//     }
// }

// // helper to create a test course for use in GET/UPDATE/DELETE tests
// func createTestCourse(t *testing.T, name, description string) *data.Course {
//     t.Helper()

//     course := &data.Course{
//         Course_Name: name,
//         Description: description,
//     }

//     err := testApp.courseModel.Insert(course)
//     if err != nil {
//         t.Fatalf("Failed to insert test course: %v", err)
//     }
//     return course
// }

// // TestCreateCourseHandler tests the POST /v1/courses endpoint
// func TestCreateCourseHandler(t *testing.T) {
//     // Clean up the table after this test function runs
//     t.Cleanup(func() {
//         tearDownCourses(t)
//     })

//     testCases := []struct {
//         name           string
//         payload        map[string]interface{}
//         expectedStatus int
//     }{
//         {
//             name: "Success: Valid Course",
//             payload: map[string]interface{}{
//                 "course":      "Go Fundamentals",
//                 "description": "A great course about Go.",
//             },
//             expectedStatus: http.StatusCreated,
//         },
//         {
//             name: "Failure: Missing Course Name",
//             payload: map[string]interface{}{
//                 "description": "A description without a name.",
//             },
//             expectedStatus: http.StatusUnprocessableEntity,
//         },
//         {
//             name: "Failure: Missing Description",
//             payload: map[string]interface{}{
//                 "course": "A course without a description.",
//             },
//             expectedStatus: http.StatusUnprocessableEntity,
//         },
//         {
//             name: "Failure: Course Name Too Long",
//             payload: map[string]interface{}{
//                 "course":      "This Course Name Is Definitely Way Too Long To Be Valid", // > 25 bytes
//                 "description": "Valid description.",
//             },
//             expectedStatus: http.StatusUnprocessableEntity,
//         },
//         {
//             name: "Failure: Description Too Long",
//             payload: map[string]interface{}{
//                 "course":      "Valid Course",
//                 "description": "This description is very, very, very, very, very, very, very, very, very, very, very, very, very long and should fail validation.", // > 100 bytes
//             },
//             expectedStatus: http.StatusUnprocessableEntity,
//         },
//     }

//     for _, tc := range testCases {
//         t.Run(tc.name, func(t *testing.T) {
//             body, _ := json.Marshal(tc.payload)
//             req := httptest.NewRequest(http.MethodPost, "/v1/courses", bytes.NewReader(body))
//             req.Header.Set("Content-Type", "application/json")

//             rr := executeRequest(req)

//             if rr.Code != tc.expectedStatus {
//                 t.Errorf("expected %d (%s); got %d (%s)\nBody: %s",
//                     tc.expectedStatus, http.StatusText(tc.expectedStatus),
//                     rr.Code, http.StatusText(rr.Code),
//                     rr.Body.String())
//             }

//             if rr.Code == http.StatusCreated {
//                 location := rr.Header().Get("Location")
//                 if location == "" {
//                     t.Error("expected 'Location' header to be set on 201 Created")
//                 }
//             }
//         })
//     }
// }

// // TestDisplayCourseHandler tests the GET /v1/courses/:id endpoint
// func TestDisplayCourseHandler(t *testing.T) {
//     t.Cleanup(func() {
//         tearDownCourses(t)
//     })

//     // 1. Create a known course to fetch
//     course := createTestCourse(t, "Test Course", "Test Description")

//     // 2. Run tests
//     t.Run("Success: Valid ID", func(t *testing.T) {
//         url := fmt.Sprintf("/v1/courses/%d", course.ID)
//         req := httptest.NewRequest(http.MethodGet, url, nil)
//         rr := executeRequest(req)

//         if rr.Code != http.StatusOK {
//             t.Fatalf("expected 200 OK; got %d", rr.Code)
//         }

//         var jsonResponse struct {
//             Course data.Course `json:"course"`
//         }
//         if err := json.NewDecoder(rr.Body).Decode(&jsonResponse); err != nil {
//             t.Fatalf("Failed to decode response: %v", err)
//         }

//         if jsonResponse.Course.ID != course.ID || jsonResponse.Course.Course_Name != "Test Course" {
//             t.Errorf("got unexpected course data: %+v", jsonResponse.Course)
//         }
//     })

//     t.Run("Failure: Non-existent ID", func(t *testing.T) {
//         url := "/v1/courses/999999"
//         req := httptest.NewRequest(http.MethodGet, url, nil)
//         rr := executeRequest(req)

//         if rr.Code != http.StatusNotFound {
//             t.Errorf("expected 404 Not Found; got %d", rr.Code)
//         }
//     })

//     t.Run("Failure: Invalid ID", func(t *testing.T) {
//         url := "/v1/courses/abc"
//         req := httptest.NewRequest(http.MethodGet, url, nil)
//         rr := executeRequest(req)

//         if rr.Code != http.StatusNotFound {
//             t.Errorf("expected 404 Not Found; got %d", rr.Code)
//         }
//     })
// }

// // TestUpdateCourseHandler tests the PUT /v1/courses/:id endpoint
// func TestUpdateCourseHandler(t *testing.T) {
//     t.Cleanup(func() {
//         tearDownCourses(t)
//     })

//     // 1. Create a known course to update
//     course := createTestCourse(t, "Original Name", "Original Description")

//     // 2. Run tests
//     t.Run("Success: Valid Partial Update", func(t *testing.T) {
//         payload := map[string]interface{}{
//             "description": "Updated Description Only",
//         }
//         body, _ := json.Marshal(payload)
//         url := fmt.Sprintf("/v1/courses/%d", course.ID)
//         req := httptest.NewRequest(http.MethodPut, url, bytes.NewReader(body))
//         req.Header.Set("Content-Type", "application/json")

//         rr := executeRequest(req)

//         if rr.Code != http.StatusOK {
//             t.Fatalf("expected 200 OK; got %d\nBody: %s", rr.Code, rr.Body.String())
//         }

//         var jsonResponse struct {
//             Course data.Course `json:"course"`
//         }
//         if err := json.NewDecoder(rr.Body).Decode(&jsonResponse); err != nil {
//             t.Fatalf("Failed to decode response: %v", err)
//         }

//         // Name should be original, description should be new
//         if jsonResponse.Course.Course_Name != "Original Name" {
//             t.Errorf("expected name to be 'Original Name', got '%s'", jsonResponse.Course.Course_Name)
//         }
//         if jsonResponse.Course.Description != "Updated Description Only" {
//             t.Errorf("expected description to be 'Updated Description Only', got '%s'", jsonResponse.Course.Description)
//         }
//     })

//     t.Run("Failure: Invalid Update (validation fail)", func(t *testing.T) {
//         payload := map[string]interface{}{
//             "course": "", // Empty name is invalid
//         }
//         body, _ := json.Marshal(payload)
//         url := fmt.Sprintf("/v1/courses/%d", course.ID)
//         req := httptest.NewRequest(http.MethodPut, url, bytes.NewReader(body))
//         req.Header.Set("Content-Type", "application/json")

//         rr := executeRequest(req)

//         if rr.Code != http.StatusUnprocessableEntity {
//             t.Errorf("expected 422 Unprocessable Entity; got %d", rr.Code)
//         }
//     })

//     t.Run("Failure: Update Non-existent ID", func(t *testing.T) {
//         payload := map[string]interface{}{
//             "course": "Does not matter",
//         }
//         body, _ := json.Marshal(payload)
//         url := "/v1/courses/999999"
//         req := httptest.NewRequest(http.MethodPut, url, bytes.NewReader(body))
//         req.Header.Set("Content-Type", "application/json")

//         rr := executeRequest(req)

//         if rr.Code != http.StatusNotFound {
//             t.Errorf("expected 404 Not Found; got %d", rr.Code)
//         }
//     })
// }

// // TestDeleteCourseHandler tests the DELETE /v1/courses/:id endpoint
// func TestDeleteCourseHandler(t *testing.T) {
//     t.Cleanup(func() {
//         tearDownCourses(t)
//     })

//     // 1. Create a known course to delete
//     course := createTestCourse(t, "To Be Deleted", "Delete me")

//     // 2. Run tests
//     t.Run("Success: Valid ID", func(t *testing.T) {
//         url := fmt.Sprintf("/v1/courses/%d", course.ID)
//         req := httptest.NewRequest(http.MethodDelete, url, nil)
//         rr := executeRequest(req)

//         if rr.Code != http.StatusOK {
//             t.Fatalf("expected 200 OK; got %d", rr.Code)
//         }

//         // 3. Verify it's gone
//         _, err := testApp.courseModel.Get(course.ID)
//         if err == nil || err != data.ErrRecordNotFound {
//             t.Error("expected course to be deleted, but it still exists")
//         }
//     })

//     t.Run("Failure: Non-existent ID", func(t *testing.T) {
//         url := "/v1/courses/999999"
//         req := httptest.NewRequest(http.MethodDelete, url, nil)
//         rr := executeRequest(req)

//         if rr.Code != http.StatusNotFound {
//             t.Errorf("expected 404 Not Found; got %d", rr.Code)
//         }
//     })
// }

// // TestListCoursesHandler tests the GET /v1/courses endpoint
// func TestListCoursesHandler(t *testing.T) {
//     t.Cleanup(func() {
//         tearDownCourses(t)
//     })

//     // 1. Create multiple courses
//     createTestCourse(t, "Go Programming", "Learn Go")
//     createTestCourse(t, "SQL Basics", "Learn SQL")
//     createTestCourse(t, "Go Web Development", "Build web apps")

//     // 2. Define response struct
//     type listResponse struct {
//         Courses  []data.Course `json:"courses"`
//         Metadata data.Metadata `json:"@metadata"`
//     }

//     // 3. Run tests
//     t.Run("List all (no filters)", func(t *testing.T) {
//         req := httptest.NewRequest(http.MethodGet, "/v1/courses", nil)
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
//         if len(resp.Courses) != 3 {
//             t.Errorf("expected 3 courses in list; got %d", len(resp.Courses))
//         }
//     })

//     t.Run("Filter by course name", func(t *testing.T) {
//         // This should match "Go Programming" and "Go Web Development"
//         req := httptest.NewRequest(http.MethodGet, "/v1/courses?course=Go", nil)
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
//         if len(resp.Courses) != 2 {
//             t.Errorf("expected 2 courses in list; got %d", len(resp.Courses))
//         }
//     })

//     t.Run("Pagination", func(t *testing.T) {
//         req := httptest.NewRequest(http.MethodGet, "/v1/courses?page=2&page_size=2", nil)
//         rr := executeRequest(req)
//         if rr.Code != http.StatusOK {
//             t.Fatalf("expected 200 OK; got %d", rr.Code)
//         }

//         var resp listResponse
//         if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
//             t.Fatalf("Failed to decode JSON: %v", err)
//         }

//         if resp.Metadata.TotalRecords != 3 { // Total is still 3
//             t.Errorf("expected 3 total records; got %d", resp.Metadata.TotalRecords)
//         }
//         if len(resp.Courses) != 1 { // Only 1 record on page 2
//             t.Errorf("expected 1 course on page 2; got %d", len(resp.Courses))
//         }
//     })
// }
