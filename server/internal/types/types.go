package types

import "time"

type User struct {
	ID       int    `db:"id"`
	UserName string `db:"user_name"`

	Email             string    `db:"email"`
	EncryptedPassword string    `db:"encrypted_password"`
	CreatedAt         time.Time `db:"created_at"`
}

type CreateUserRequest struct {
	UserName string `json:"user_name" validate:"required"`

	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserSignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserSession struct {
	Id       int    `json:"id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
}
type ErrorResponse struct {
	Error string
}

type ShortenRequest struct {
	LongUrl string `json:"long_url" validate:"required,long_url"`
}

type Link struct {
	Id          int       `db:"id"`
	OriginalURL string    `db:"original_url"`
	ShortURL    string    `db:"short_url"`
	CreatedAt   time.Time `db:"created_at"`
	UserId      int       `db:"user_id"`
	IsEnabled   bool      `db:"is_enabled"`
}

type CreateShortURLResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"long_url"`
	LinkId      int    `json:"link_id"`
}

type Clicks struct {
	Id         int       `db:"id"`
	ShortCode  string    `db:"short_code"`
	Timestamp  time.Time `db:"time_stamp"`
	DeviceType string    `db:"device_type"`
	Location   string    `db:"location"`
}
