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
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
}


type UserSignInRequest struct {
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserSession struct {
	Id int
	FirstName string
	LastName  string
	Email     string
}
type ErrorResponse struct {
	Error string
}


type ShortenRequest struct {
	LongUrl string
}
