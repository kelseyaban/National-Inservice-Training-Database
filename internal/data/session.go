package data

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "time"

    "github.com/kelseyaban/National-Inservice-Training-Database/internal/validator"
)

type Session struct {
    ID            int64     `json:"id"`
    CourseID      int64     `json:"course_id"`
    FormationID   int64     `json:"formation_id"`
    FacilitatorID int64     `json:"facilitator_id"`
    CreatedAt     time.Time `json:"created_at"`
}



func ValidateSession(v *validator.Validator, session *Session) {
    v.Check(session.CourseID > 0, "course_id", "must be provided")
    v.Check(session.FormationID > 0, "formation_id", "must be provided")
    v.Check(session.FacilitatorID > 0, "facilitator_id", "must be provided")
}


type SessionModel struct {
    DB *sql.DB
}


// Insert a new row in the role table
func (s SessionModel) Insert(session *Session) error {
    query := `
        INSERT INTO session (course_id, formation_id, facilitator_id)
        VALUES ($1, $2, $3)
        RETURNING id, created_at
    `
    args := []any{session.CourseID, session.FormationID, session.FacilitatorID}

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    return s.DB.QueryRowContext(ctx, query, args...).Scan(&session.ID, &session.CreatedAt)
}
func (s SessionModel) Get(id int64) (*Session, error) {
    if id < 1 {
        return nil, ErrRecordNotFound
    }

    query := `
        SELECT id, course_id, formation_id, facilitator_id, created_at
        FROM session
        WHERE id = $1
    `

    var session Session
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    err := s.DB.QueryRowContext(ctx, query, id).Scan(
        &session.ID,
        &session.CourseID,
        &session.FormationID,
        &session.FacilitatorID,
        &session.CreatedAt,
    )

    if err != nil {
        switch {
        case errors.Is(err, sql.ErrNoRows):
            return nil, ErrRecordNotFound
        default:
            return nil, err
        }
    }

    return &session, nil
}
func (s SessionModel) Update(session *Session) error {
    query := `
        UPDATE session
        SET course_id = $1, formation_id = $2, facilitator_id = $3
        WHERE id = $4
    `

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    _, err := s.DB.ExecContext(ctx, query, session.CourseID, session.FormationID, session.FacilitatorID, session.ID)
    return err
}
func (s SessionModel) Delete(id int64) error {
    if id < 1 {
        return ErrRecordNotFound
    }

    query := `
	DELETE FROM session WHERE id = $1
	`

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    result, err := s.DB.ExecContext(ctx, query, id)
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
func (s SessionModel) GetAll(filters Filters) ([]*Session, Metadata, error) {
    query := fmt.Sprintf(`
        SELECT COUNT(*) OVER(), id, course_id, formation_id, facilitator_id, created_at
        FROM session
        ORDER BY %s %s, id ASC
        LIMIT $1 OFFSET $2
    `, filters.sortColumn(), filters.sortDirection())

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    rows, err := s.DB.QueryContext(ctx, query, filters.limit(), filters.offset())
    if err != nil {
        return nil, Metadata{}, err
    }
    defer rows.Close()

    totalRecords := 0
    sessions := []*Session{}

    for rows.Next() {
        var session Session
        err := rows.Scan(
            &totalRecords,
            &session.ID,
            &session.CourseID,
            &session.FormationID,
            &session.FacilitatorID,
            &session.CreatedAt,
        )
        if err != nil {
            return nil, Metadata{}, err
        }
        sessions = append(sessions, &session)
    }

    if err = rows.Err(); err != nil {
        return nil, Metadata{}, err
    }

    metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
    return sessions, metadata, nil
}
