package user

import (
	"errors"
	"fmt"

	"forum/internal/models"
	"forum/internal/service/helpers"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct { //nolint:revive
	UserRepo models.UserRepository
}

func NewUserService(userRepo models.UserRepository) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

func (u *UserService) SignUpUser(signupRequest *models.SignupRequest) error {
	existingUser, _ := u.UserRepo.GetUserByEmail(signupRequest.Email)

	if existingUser != nil && (existingUser.Email == signupRequest.Email || existingUser.Username == signupRequest.Username) {
		return models.ErrUserExists
	}

	id, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("create uuid: %w", err)
	}

	hashPass, err := helpers.HashPassword(signupRequest.Password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	userModel := models.User{
		ID:       id.String(),
		Username: signupRequest.Username,
		Email:    signupRequest.Email,
		Password: hashPass,
	}

	if err := u.UserRepo.AddUser(userModel); err != nil {
		return fmt.Errorf("add user: %w", err)
	}

	return nil
}

func (u *UserService) Login(email, password string) (*models.User, error) {
	user, err := u.UserRepo.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return nil, models.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, models.ErrWrongPassword
	}

	return user, nil
}

func (u *UserService) GetUserByUserID(id string) (*models.User, error) {
	return u.UserRepo.GetUserByUserID(id) //nolint:wrapcheck
}
