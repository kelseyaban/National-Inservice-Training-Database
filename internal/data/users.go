// Filename: internal/data/users.go
package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"github.com/kelseyaban/National-Inservice-Training-Database/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

// Duplicate error message
var ErrDuplicateEmail = errors.New("duplicate email")
var AnonymousUser = &User{}

type User struct {
	ID               int64     `json:"id"`
	RegulationNumber string    `json:"regulation_number"`
	Username         string    `json:"username"`
	FName            string    `json:"fname"`
	LName            string    `json:"lname"`
	Email            string    `json:"email"`
	Gender           string    `json:"gender"`
	Formation        int       `json:"formation"`
	Rank             int       `json:"rank"`
	Postings         int       `json:"postings"`
	Password         password  `json:"-"`
	Activated        bool      `json:"activated"`
	Version          int       `json:"-"`
	CreatedAt        time.Time `json:"created_at"`
}

type password struct {
	plaintext *string
	hash      []byte
}

type UserModel struct {
	DB *sql.DB
}

// Set the hashing of the password
func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plaintextPassword
	p.hash = hash
	return nil
}

// Compare the client-provided plaintext password with saved-hashed version
func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))

	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

// Validate the email address
func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

// Check that a valid password is provided
func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

// Validate a user
func ValidateUser(v *validator.Validator, user User) {
	v.Check(user.RegulationNumber != "", "regulation_number", "must be provided")
	v.Check(len(user.RegulationNumber) <= 100, "regulation_number", "must not be more than 100 bytes long")

	v.Check(user.Username != "", "username", "must be provided")
	v.Check(len(user.Username) <= 200, "username", "must not be more than 200 bytes long")

	v.Check(user.FName != "", "fname", "must be provided")
	v.Check(len(user.FName) <= 200, "fname", "must not be more than 200 bytes long")

	v.Check(user.LName != "", "lname", "must be provided")
	v.Check(len(user.LName) <= 200, "lname", "must not be more than 200 bytes long")

	v.Check(user.Formation >= 0, "formation", "must be a valid formation id")
	v.Check(user.Rank >= 0, "rank", "must be a valid rank id")
	v.Check(user.Postings >= 0, "postings", "must be a valid posting id")

	// Optional: simple gender presence check (adjust allowed values as needed)
	v.Check(user.Gender != "", "gender", "must be provided")

	// validate email for user
	ValidateEmail(v, user.Email)
	// validate the plain text password
	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}
	// Check if we messed up in our codebase
	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}

// Inserty a new user into the db
func (u UserModel) Insert(user *User) error {
	query := `
			INSERT INTO users (regulation_number, username, fname, lname, email, password_hash, activated, gender, formation_id, rank_id, posting_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			RETURNING id, created_at, version
			`

	args := []any{user.RegulationNumber, user.Username, user.FName, user.LName, user.Email, user.Password.hash, user.Activated, user.Gender, user.Formation, user.Rank, user.Postings}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// If the email address is already used, error message will be sent
	err := u.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violated unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

// Get a user from the db based on their email provided
func (u UserModel) GetByEmail(email string) (*User, error) {
	query := `
			SELECT id, regulation_number, username, fname, lname, email, password_hash,
		       activated, gender, formation_id, rank_id, posting_id, version, created_at
			FROM users
			WHERE email = $1
			`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.RegulationNumber,
		&user.Username,
		&user.FName,
		&user.LName,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Gender,
		&user.Formation,
		&user.Rank,
		&user.Postings,
		&user.Version,
		&user.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

// Update a user. The version number determins id the query will me ran
// if it doesn't match the previous edit, query will fail and user will need to try again later
func (u UserModel) Update(user *User) error {
	query := `
		UPDATE users
		SET regulation_number = $1, username = $2, fname = $3, lname = $4, email = $5, password_hash = $6, activated = $7, gender = $8, formation_id = $9, rank_id = $10, posting_id = $11, version = version + 1 
		WHERE id = $12 AND version = $13
		RETURNING version
		`
	args := []any{
		&user.RegulationNumber,
		&user.Username,
		&user.FName,
		&user.LName,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Gender,
		&user.Formation,
		&user.Rank,
		&user.Postings,
		&user.ID,
		&user.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)

	// Check for errors during update
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violated unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

// Verify token to user. We need to hash the passed in token
func (u UserModel) GetForToken(tokenScope, tokenPlaintext string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	// Join
	query := `
        SELECT users.id, users.created_at, users.username,
               users.email, users.password_hash, users.activated, users.version
        FROM users
        INNER JOIN tokens
        ON users.id = tokens.user_id
        WHERE tokens.hash = $1
        AND tokens.scope = $2 
        AND tokens.expiry > $3
	`
	args := []any{tokenHash[:], tokenScope, time.Now()}
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := u.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	// Return the matching user.
	return &user, nil
}

// Let's check if the current user is anonymous
func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

// Activatng the user
func (u UserModel) Activate(user *User) error {
	query := `
        UPDATE users
        SET activated = true, version = version + 1
        WHERE id = $1 AND version = $2
        RETURNING version
    `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return u.DB.QueryRowContext(ctx, query, user.ID, user.Version).Scan(&user.Version)
}

// Get all users from the database
func (u UserModel) GetAll(id int64, regNumber, username, fname, lname, email, gender string, formation, rank, postings int, filters Filters) ([]*User, Metadata, error) {
	// Build query using these parameters
	query := fmt.Sprintf(`
        SELECT COUNT(*) OVER(), id, regulation_number, username, fname, lname, email, 
               gender, formation_id, rank_id, posting_id
        FROM users
        WHERE (to_tsvector('simple', username) @@ plainto_tsquery('simple', $1) OR $1 = '')
        ORDER BY %s %s, id ASC
        LIMIT $2 OFFSET $3`,
		filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use username for $1, then page size and offset
	rows, err := u.DB.QueryContext(ctx, query, username, filters.PageSize, (filters.Page-1)*filters.PageSize)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	users := []*User{}

	for rows.Next() {
		var user User
		err := rows.Scan(
			&totalRecords,
			&user.ID,
			&user.RegulationNumber,
			&user.Username,
			&user.FName,
			&user.LName,
			&user.Email,
			&user.Gender,
			&user.Formation,
			&user.Rank,
			&user.Postings,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return users, metadata, nil
}

// Update User Information
func (u UserModel) UpdateUser(user *User) error {
	query := `
		UPDATE users
		SET regulation_number = $1, username = $2, fname = $3, lname = $4, email = $5, gender = $6, formation_id = $7, rank_id = $8, posting_id = $9, version = version + 1 
		WHERE id = $10 AND version = $11
		RETURNING version
		`
	args := []any{
		&user.RegulationNumber,
		&user.Username,
		&user.FName,
		&user.LName,
		&user.Email,
		&user.Gender,
		&user.Formation,
		&user.Rank,
		&user.Postings,
		&user.ID,
		&user.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return u.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.RegulationNumber, &user.Username, &user.FName, &user.LName, &user.Email, &user.Gender, &user.Formation, &user.Rank, &user.Postings)

}

// Delete user
func (u UserModel) Delete(id int64) error {

	// Check
	if id < 1 {
		return ErrRecordNotFound
	}

	// the SQL query to be executed
	query := `
		DELETE FROM users
		WHERE id = $1
		`

	// Context with a 3-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// execute query against the database
	result, err := u.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	// Check how many rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

// Get by id
func (u UserModel) GetByID(id int64) (*User, error) {
	// Check if the id is valid
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, regulation_number, username, fname, lname, email, password_hash,
		       activated, gender, formation_id, rank_id, posting_id, version, created_at
		FROM users
		WHERE id = $1
		`

	var user User

	// Context with a 3-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.RegulationNumber,
		&user.Username,
		&user.FName,
		&user.LName,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Gender,
		&user.Formation,
		&user.Rank,
		&user.Postings,
		&user.Version,
		&user.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

// UpdatePassword updates the user's password hash in the database.
func (u UserModel) UpdatePassword(id int64, newPassword string) error {
	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update query
	query := `
		UPDATE users
		SET password_hash = $1, version = version + 1
		WHERE id = $2
		RETURNING id
	`

	// Execute the update
	var returnedID int64
	err = u.DB.QueryRow(query, hashedPassword, id).Scan(&returnedID)
	if err != nil {
		return err
	}

	return nil
}
