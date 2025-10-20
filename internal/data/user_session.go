package data

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "time"

    // "github.com/lib/pq"
    "github.com/kelseyaban/National-Inservice-Training-Database/internal/validator"
)

type UserSession struct {
    ID                   int64     `json:"id"`
    TraineeID            int64     `json:"trainee_id"`
    SessionID            int64     `json:"session_id"`
    CreditHoursCompleted int64     `json:"credithours_completed"`
    Grade                string    `json:"grade"`
    Feedback             string    `json:"feedback"`
    Version              int       `json:"-"`
    CreatedAt            time.Time `json:"created_at"`
}

// ------------------- VALIDATION -------------------

func ValidateUserSession(v *validator.Validator, us *UserSession) {
    v.Check(us.SessionID > 0, "session_id", "must be provided and greater than zero")
    v.Check(us.CreditHoursCompleted >= 0, "credithours_completed", "must be provided")
    v.Check(us.Grade != "", "grade", "must be provided")
    v.Check(len(us.Grade) <= 25, "grade", "must not be more than 25 bytes long")
    v.Check(us.Feedback != "", "feedback", "must be provided")
    v.Check(len(us.Feedback) <= 255, "feedback", "must not be more than 255 bytes long")
}

// ------------------- MODEL STRUCT -------------------

type UserSessionModel struct {
    DB *sql.DB
}

// ------------------- ADD -------------------

func (m UserSessionModel) AddUserSession(us *UserSession) error {
    query := `
        INSERT INTO user_session (trainee_id, session_id, credithours_completed, grade, feedback)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at, version
    `
    args := []any{us.TraineeID, us.SessionID, us.CreditHoursCompleted, us.Grade, us.Feedback}

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    return m.DB.QueryRowContext(ctx, query, args...).Scan(&us.ID, &us.CreatedAt, &us.Version)
}

// ------------------- GET BY ID -------------------

func (m UserSessionModel) GetUserSession(id int64) (*UserSession, error) {
    query := `
        SELECT id, trainee_id, session_id, credithours_completed, grade, feedback, created_at, version
        FROM user_session
        WHERE id = $1
    `
    var us UserSession

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    err := m.DB.QueryRowContext(ctx, query, id).Scan(
        &us.ID,
        &us.TraineeID,
        &us.SessionID,
        &us.CreditHoursCompleted,
        &us.Grade,
        &us.Feedback,
        &us.CreatedAt,
        &us.Version,
    )

    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("no user session found with id %d", id)
        }
        return nil, err
    }

    return &us, nil
}

// ------------------- UPDATE -------------------

func (m UserSessionModel) UpdateUserSession(us *UserSession) error {
    query := `
        UPDATE user_session
        SET credithours_completed = $1, grade = $2, feedback = $3, version = version + 1
        WHERE id = $4
        RETURNING version
    `
    args := []any{us.CreditHoursCompleted, us.Grade, us.Feedback, us.ID}

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    err := m.DB.QueryRowContext(ctx, query, args...).Scan(&us.Version)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return fmt.Errorf("no user session found to update with id %d", us.ID)
        }
        return err
    }

    return nil
}

// ------------------- DELETE -------------------

func (m UserSessionModel) DeleteUserSession(id int64) error {
    query := `DELETE FROM user_session WHERE id = $1`

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    result, err := m.DB.ExecContext(ctx, query, id)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return fmt.Errorf("no user session found to delete with id %d", id)
    }

    return nil
}

// ------------------- GET ALL -------------------

func (m UserSessionModel) GetAllUserSessions() ([]*UserSession, error) {
    query := `
        SELECT id, trainee_id, session_id, credithours_completed, grade, feedback, created_at, version
        FROM user_session
        ORDER BY created_at DESC
    `

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    rows, err := m.DB.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var sessions []*UserSession

    for rows.Next() {
        var us UserSession
        err := rows.Scan(
            &us.ID,
            &us.TraineeID,
            &us.SessionID,
            &us.CreditHoursCompleted,
            &us.Grade,
            &us.Feedback,
            &us.CreatedAt,
            &us.Version,
        )
        if err != nil {
            return nil, err
        }
        sessions = append(sessions, &us)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return sessions, nil
}