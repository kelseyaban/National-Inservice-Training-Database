package data

import (
	"context"
	"database/sql"
	"time"
	"errors"
	"fmt"
	// "github.com/lib/pq"
	"github.com/kelseyaban/National-Inservice-Training-Database/internal/validator"

)

type FacilitatorRating struct {
    ID        int64     `json:"id" ` 
    UserID    int64     `json:"user_id"`
    Rating    int       `json:"rating"`
    CreatedAt time.Time `json:"-"`
}

// Validator function for facilitator ratings
func ValidateFacilitatorRating(v *validator.Validator, fr *FacilitatorRating) {
    v.Check(fr.UserID > 0, "user_id", "must be provided")
    v.Check(fr.Rating >= 1 && fr.Rating <= 5, "rating", "must be between 1 and 5")
}

// Setup model
type FacilitatorRatingModel struct {
	DB *sql.DB
}


func (f FacilitatorRatingModel) Insert(fr *FacilitatorRating) error {
    query := `
        INSERT INTO facilitator_rating (user_id, rating)
        VALUES ($1, $2)
        RETURNING id, created_at
    `
    args := []any{fr.UserID, fr.Rating}

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    return f.DB.QueryRowContext(ctx, query, args...).Scan(&fr.ID, &fr.CreatedAt)
}


func (f FacilitatorRatingModel) Get(id int64) (*FacilitatorRating, error) {
    if id < 1 {
        return nil, ErrRecordNotFound
    }

    query := `
        SELECT id, user_id, rating, created_at
        FROM facilitator_rating
        WHERE id = $1
    `

    var fr FacilitatorRating

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    err := f.DB.QueryRowContext(ctx, query, id).Scan(&fr.ID, &fr.UserID, &fr.Rating, &fr.CreatedAt)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrRecordNotFound
        }
        return nil, err
    }

    return &fr, nil
}

func (f FacilitatorRatingModel) GetAll(userID int64, filters Filters) ([]*FacilitatorRating, Metadata, error) {
    query := fmt.Sprintf(`
        SELECT COUNT(*) OVER(), id, user_id, rating, created_at
        FROM facilitator_rating
        WHERE ($1 = 0 OR user_id = $1)
        ORDER BY %s %s, id ASC
        LIMIT $2 OFFSET $3
    `, filters.sortColumn(), filters.sortDirection())

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    rows, err := f.DB.QueryContext(ctx, query, userID, filters.limit(), filters.offset())
    if err != nil {
        return nil, Metadata{}, err
    }
    defer rows.Close()

    totalRecords := 0
    ratings := []*FacilitatorRating{}

    for rows.Next() {
        var fr FacilitatorRating
        err := rows.Scan(&totalRecords, &fr.ID, &fr.UserID, &fr.Rating, &fr.CreatedAt)
        if err != nil {
            return nil, Metadata{}, err
        }
        ratings = append(ratings, &fr)
    }

    if err = rows.Err(); err != nil {
        return nil, Metadata{}, err
    }

    metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
    return ratings, metadata, nil
}
