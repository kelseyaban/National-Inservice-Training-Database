// Filename: internal/data/users.go
package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
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
	FName            string    `json:"fname"`
	Username         string    `json:"username"`
	LName            string    `json:"lname"`
	Email            string    `json:"email"`
	Gender           string    `json:"gender"`
	Role             string    `json:"role"`
	Formation        string    `json:"formation"`
	Rank             string    `json:"rank"`
	Postings         string    `json:"postings"`
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
	v.Check(user.Username != "", "username", "must be provided")
	v.Check(len(user.Username) <= 200, "username", "must not be more than 200 bytes long")

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
			INSERT INTO users (username, email, password_hash, activated)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at, version
			`

	args := []any{user.Username, user.Email, user.Password.hash, user.Activated}

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

// Get a user from the db based on their emial provided
func (u UserModel) GetByEmail(email string) (*User, error) {
	query := `
			SELECT id, created_at, username, email, password_hash, activated, version 
			FROM users
			WHERE email = $1
			`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, email).Scan(
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
	return &user, nil
}

// Update a user. The version number determins id the query will me ran
// if it doesn't match the previous edit, query will fail and user will need to try again later
func (u UserModel) Update(user *User) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, password_hash = $3, activated = $4, version = version + 1 
		WHERE id = $5 AND version = $6
		RETURNING version
		`
	args := []any{
		user.Username,
		user.Email,
		user.Password.hash,
		user.Activated,
		user.ID,
		user.Version,
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
