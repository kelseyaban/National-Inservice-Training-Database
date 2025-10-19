// Filename: internal/data/course.go
package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/kelseyaban/National-Inservice-Training-Database/internal/validator"
)

type Course struct {
	ID          int64     `json:"id"`
	Course_Name string    `json:"course"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"-"` // database timestamp
}

// Performs the validation checks
func ValidateCourse(v *validator.Validator, course *Course) {
	// check if the Course name field is empty
	v.Check(course.Course_Name != "", "course", "must be provided")
	// check if the Description field us empty
	v.Check(course.Description != "", "description", "must be provided")
	// chekc if the content in the field is empty
	v.Check(len(course.Description) <= 100, "descrption", "must not be more than 100 bytes long")
	// check if the Author field is empty
	v.Check(len(course.Course_Name) <= 25, "course", "must not be more than 25 bytes long")
}

type CourseModel struct {
	DB *sql.DB
}

// Insert new course into the database
func (c CourseModel) Insert(course *Course) error {
	query := `
		INSERT INTO course (course, description)
		VALUES ($1, $2)
		RETURNING id, created_at`

	// values to replace $1 and $2
	args := []any{course.Course_Name, course.Description}

	// Context with a 3-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// execute query against the database
	return c.DB.QueryRowContext(ctx, query, args...).Scan(&course.ID, &course.CreatedAt)
}

// Get a specific course from the database
func (c CourseModel) Get(id int64) (*Course, error) {
	// Check if the id is valid
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	// the SQL query to be executed
	query := `
		SELECT id, course, description, created_at
		FROM course
		WHERE id = $1`

	// course variable to hold the data returned by the query
	var course Course

	// Context with a 3-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// execute the query against the database
	err := c.DB.QueryRowContext(ctx, query, id).Scan(
		&course.ID,
		&course.Course_Name,
		&course.Description,
		&course.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &course, nil
}

// Update a specific course in the database
func (c CourseModel) Update(course *Course) error {
	// the SQL query to be executed
	query := `
		UPDATE course
		SET course = $1, description = $2
		WHERE id = $3
		RETURNING id, course, description, created_at
		`

	// values to replace $1, $2, and $3
	args := []any{course.Course_Name, course.Description, course.ID}

	// Context with a 3-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.DB.QueryRowContext(ctx, query, args...).Scan(
		&course.ID,
		&course.Course_Name,
		&course.Description,
		&course.CreatedAt,
	)
}

// Delete a specific course from the database
func (c CourseModel) Delete(id int64) error {

	// Check
	if id < 1 {
		return ErrRecordNotFound
	}

	// the SQL query to be executed
	query := `
		DELETE FROM course
		WHERE id = $1
		`
	// Context with a 3-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// execute the query against the database
	result, err := c.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	// check how many rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// if no rows were affected, return a record not found error
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

// Get all courses from the database
func (c CourseModel) GetAll(course string, description string, filters Filters) ([]*Course, Metadata, error) {
	// the SQL query to be executed
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, course, description, created_at
		FROM course
		WHERE (to_tsvector('simple', course) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (to_tsvector('simple', description) @@ plainto_tsquery('simple', $2) OR $2 = '')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := c.DB.QueryContext(ctx, query, course, description, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	courses := []*Course{}

	for rows.Next() {
		var course Course
		err := rows.Scan(
			&totalRecords,
			&course.ID,
			&course.Course_Name,
			&course.Description,
			&course.CreatedAt,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		courses = append(courses, &course)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return courses, metadata, nil
}
