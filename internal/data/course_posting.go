// Filename: internal/data/course_posting.go
package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/kelseyaban/National-Inservice-Training-Database/internal/validator"
)

type CoursePosting struct {
	ID          int64     `json:"id"`
	CourseID    int64     `json:"course_id"`
	PostingID   int64     `json:"posting_id"`
	Mandatory   bool      `json:"mandatory"`
	CreditHours int64     `json:"credithours"`
	RankID      int64     `json:"rank_id"`
	CreatedAt   time.Time `json:"-"`
}

func ValidateCoursePosting(v *validator.Validator, coursePosting *CoursePosting) {
	// check if the CourseID field is valid
	v.Check(coursePosting.CourseID > 0, "course_id", "must be provided and greater than zero")
	// check if the PostingID field is valid
	v.Check(coursePosting.PostingID > 0, "posting_id", "must be provided and greater than zero")
	// check if the RankID field is valid
	v.Check(coursePosting.RankID > 0, "rank_id", "must be provided and greater than zero")
	// check if the CreditHours field is valid
	v.Check(coursePosting.CreditHours >= 0, "credithours", "must be provided and non-negative")
}

type CoursePostingModel struct {
	DB *sql.DB
}

// Insert new course posting into the database
func (c CoursePostingModel) Insert(courseposting *CoursePosting) error {
	// Insert into course_postings table
	query := `
		INSERT INTO course_posting (course_id, posting_id, mandatory, credithours, rank_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
		`
	args := []any{courseposting.CourseID, courseposting.PostingID, courseposting.Mandatory, courseposting.CreditHours, courseposting.RankID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.DB.QueryRowContext(ctx, query, args...).Scan(&courseposting.ID, &courseposting.CreatedAt)

}

// Get a single course_posting by ID
func (c CoursePostingModel) Get(id int64) (*CoursePosting, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, course_id, posting_id, mandatory, credithours, rank_id, created_at
		FROM course_posting
		WHERE id = $1`

	var courseposting CoursePosting

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := c.DB.QueryRowContext(ctx, query, id).Scan(
		&courseposting.ID,
		&courseposting.CourseID,
		&courseposting.PostingID,
		&courseposting.Mandatory,
		&courseposting.CreditHours,
		&courseposting.RankID,
		&courseposting.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &courseposting, nil
}

// Update
func (c CoursePostingModel) Update(courseposting *CoursePosting) error {
	query := `
		UPDATE course_posting
		SET course_id = $1, posting_id = $2, mandatory = $3, credithours = $4, rank_id = $5
		WHERE id = $6
		RETURNING id, course_id, posting_id, mandatory, credithours, rank_id, created_at`

	args := []any{
		courseposting.CourseID,
		courseposting.PostingID,
		courseposting.Mandatory,
		courseposting.CreditHours,
		courseposting.RankID,
		courseposting.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.DB.QueryRowContext(ctx, query, args...).Scan(
		&courseposting.ID,
		&courseposting.CourseID,
		&courseposting.PostingID,
		&courseposting.Mandatory,
		&courseposting.CreditHours,
		&courseposting.RankID,
		&courseposting.CreatedAt,
	)
}

// Delete
func (c CoursePostingModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM course_posting
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := c.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

// Get all course postiings
func (c CoursePostingModel) GetAll(courseID, postingID int64, mandatory bool, credithours int64, rankID int64, filters Filters) ([]*CoursePosting, Metadata, error) {
	// Query to get all course postings
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER() AS total_records,
	       cp.id,
	       cp.course_id,
	       cp.posting_id,
	       cp.mandatory,
	       cp.credithours,
	       cp.rank_id,
	       cp.created_at
	FROM course_posting cp
	WHERE ($1 = 0 OR cp.course_id = $1)
	  AND ($2 = 0 OR cp.posting_id = $2)
	  AND ($3::boolean IS NULL OR cp.mandatory = $3)
	  AND ($4 = 0 OR cp.credithours = $4)
	  AND ($5 = 0 OR cp.rank_id = $5)
	ORDER BY %s %s
	LIMIT $6
	OFFSET $7;`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var mandatoryParam sql.NullBool
	if mandatory { // If you only want true/false, you can use this logic
		mandatoryParam = sql.NullBool{Bool: true, Valid: true}
	} else {
		mandatoryParam = sql.NullBool{Valid: false}
	}

	rows, err := c.DB.QueryContext(ctx, query,
		courseID,
		postingID,
		mandatoryParam,
		credithours,
		rankID,
		filters.limit(),
		filters.offset(),
	)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	coursePostings := []*CoursePosting{}

	for rows.Next() {
		var courseposting CoursePosting
		err := rows.Scan(
			&totalRecords,
			&courseposting.ID,
			&courseposting.CourseID,
			&courseposting.PostingID,
			&courseposting.Mandatory,
			&courseposting.CreditHours,
			&courseposting.RankID,
			&courseposting.CreatedAt,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		coursePostings = append(coursePostings, &courseposting)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return coursePostings, metadata, nil
}
