// Filename: internal/data/attendance.go
package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/kelseyaban/National-Inservice-Training-Database/internal/validator"
)

type Attendance struct {
	ID               int64     `json:"id"`
	UserSessionID    int64     `json:"user_session_id"`
	AttendanceStatus bool      `json:"attendance"`
	Date             time.Time `json:"date"`
	CreatedAt        time.Time `json:"-"`
}

// Performs the checks
func ValidateAttendance(v *validator.Validator, attendance *Attendance) {
	// check if the UserSessionID field is provided
	v.Check(attendance.UserSessionID != 0, "user_session_id", "must be provided")
	// check if the AttendanceStatus field is provided
	v.Check(attendance.AttendanceStatus == true || attendance.AttendanceStatus == false, "attendance", "must be provided")
	// check if the Date field is provided
	v.Check(!attendance.Date.IsZero(), "date", "must be provided")
}

type AttendanceModel struct {
	DB *sql.DB
}

// Insert new attendance record into the database
func (a AttendanceModel) Insert(attendance *Attendance) error {
	query := `
		INSERT INTO attendance (user_session_id, attendance, date)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	// values to replace $1, $2, and $3
	args := []any{attendance.UserSessionID, attendance.AttendanceStatus, attendance.Date}

	// Context with a 3-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// execute query against the database
	return a.DB.QueryRowContext(ctx, query, args...).Scan(&attendance.ID, &attendance.CreatedAt)
}

// Get a specific attendance record of user from the database
func (a AttendanceModel) GetIdividualAttendance(id int64) (*Attendance, error) {
	// Check if the id is valid
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, user_session_id, attendance, date, created_at
		FROM attendance
		WHERE id = $1`

	var attendance Attendance

	// Context with a 3-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// execute query against the database
	err := a.DB.QueryRowContext(ctx, query, id).Scan(
		&attendance.ID,
		&attendance.UserSessionID,
		&attendance.AttendanceStatus,
		&attendance.Date,
		&attendance.CreatedAt,
	)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &attendance, nil
}

// Get all attendance records from the database
func (a AttendanceModel) GetAll() ([]*Attendance, error) {
	query := `
		SELECT id, user_session_id, attendance, date, created_at
		FROM attendance
		ORDER BY created_at DESC`

	// Context with a 3-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// execute query against the database
	rows, err := a.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	attendances := []*Attendance{}

	for rows.Next() {
		var attendance Attendance

		err := rows.Scan(
			&attendance.ID,
			&attendance.UserSessionID,
			&attendance.AttendanceStatus,
			&attendance.Date,
			&attendance.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		attendances = append(attendances, &attendance)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return attendances, nil
}

// Update a specific attendance record in the database
func (a AttendanceModel) Update(attendance *Attendance) error {
	query := `
		UPDATE attendance
		SET attendance = $1, date = $2
		WHERE id = $3
		RETURNING id, user_session_id,  attendance, date, created_at`

	// values to replace $1, $2, and $3
	args := []any{attendance.AttendanceStatus, attendance.Date, attendance.ID}

	// Context with a 3-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// execute query against the database
	return a.DB.QueryRowContext(ctx, query, args...).Scan(
		&attendance.ID,
		&attendance.UserSessionID,
		&attendance.AttendanceStatus,
		&attendance.Date,
		&attendance.CreatedAt,
	)
}
