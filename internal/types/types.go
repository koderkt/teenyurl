package types

import "time"

type User struct {
	ID                int       `db:"id"`
	FirstName         string    `db:"first_name"`
	LastName          string    `db:"last_name"`
	Email             string    `db:"email"`
	EncryptedPassword string    `db:"encrypted_password"`
	CreatedAt         time.Time `db:"created_at"`
}

type CreateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type ErrorResponse struct {
	Error string
}
