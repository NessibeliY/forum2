package models

import (
	"errors"

	"forum/internal/validator"
)

type contextKey string

const (
	// Define a key for the session
	SessionKey contextKey = "session"
)

var (
	ErrUserExists             = errors.New("user already exists")
	ErrNotFound               = errors.New("error: Not found")
	ErrUserNotFound           = errors.New("user not found")
	ErrWrongPassword          = errors.New("wrong password")
	ErrSessionNotFound        = errors.New("session not found")
	ErrPostNotCreated         = errors.New("post not created")
	ErrUUIDCreate             = errors.New("id not generated")
	ErrNotCreated             = errors.New("error: not created")
	ErrIncrementLikeInPost    = errors.New("error: increment like in post")
	ErrDecrementLikeInPost    = errors.New("error: decrement like in post")
	ErrDecrementDisLikeInPost = errors.New("error: decrement dislike in post")
	ErrIncrementCommentInPost = errors.New("error: increment comment in post")
	ErrDecrementCommentInpost = errors.New("error: decrement comment in post")
	ErrDeletDisLikeInPost     = errors.New("error: delete dislike in post")
	ErrDeleteLikeInPost       = errors.New("error: delete like in post")
	ErrDeleteLikeInComment    = errors.New("error: delete like in comment")
	ErrDeletePost             = errors.New("error: delete post")
	ErrDeleteComment          = errors.New("error: delete comment")
)

type User struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UserRepository interface {
	AddUser(user User) error
	GetUserByEmail(email string) (*User, error)
	GetUserByUserID(id string) (*User, error)
}

type UserService interface {
	SignUpUser(signupRequest *SignupRequest) error
	Login(email, password string) (*User, error)
	GetUserByUserID(id string) (*User, error)
}

type SignupRequest struct {
	Username      string
	Email         string
	Password      string
	ErrorMessages ErrorMessage
	IsAuth        bool
}

type LoginRequest struct {
	User          AuthData
	ErrorMessages ErrorMessage
	IsAuth        bool
}

type AuthData struct {
	Email    string
	Password string
}

type ErrorMessage struct {
	Email       string
	Username    string
	Password    string
	Title       string
	Description string
	Tags        string
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(v.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePassword(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(v.MinChars(password, 8), "password", "password contain at least 8 characters")
	v.Check(v.ValidPassword(password), "password", "password must be contain 1 upper character,1 lower character and 1 digit")
}

func ValidateSignupRequest(v *validator.Validator, u *SignupRequest) {
	v.Check(u.Username != "", "username", "must be provided")
	v.Check(len(u.Username) >= 3, "username", "username must be at least 3 characters long")

	ValidateEmail(v, u.Email)
	ValidatePassword(v, u.Password)
}
