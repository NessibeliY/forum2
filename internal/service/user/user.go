package user

import (
	"forum/internal/models"
	"forum/internal/service/helpers"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo models.UserRepository
}

func NewUserService(userRepo models.UserRepository) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

func (u *UserService) SignUpUser(signupRequest *models.SignupRequest) error {
	existingUser, err := u.UserRepo.GetUserByEmail(signupRequest.Email)
	if err != nil {
		return err
	}
	if existingUser != nil && (existingUser.Email == signupRequest.Email || existingUser.UserName == signupRequest.UserName) {
		return models.ErrUserExists
	}

	id, err := uuid.NewV4()
	if err != nil {
		return err
	}

	hashPass, err := helpers.HashPassword(signupRequest.Password)
	if err != nil {
		return err
	}

	userModel := models.User{
		ID:       id.String(),
		UserName: signupRequest.UserName,
		Email:    signupRequest.Email,
		Password: hashPass,
	}

	if err := u.UserRepo.AddUser(userModel); err != nil {
		return err
	}

	return nil
}

func (u *UserService) Login(email, password string) (*models.User, error) {
	user, err := u.UserRepo.GetUserByEmail(email)
	if err != nil {
		if err == models.ErrUserNotFound {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, models.ErrWrongPassword
	}

	return user, nil
}

func (u *UserService) GetUserByUserID(id string) (*models.User, error) {
	user, err := u.UserRepo.GetUserByUserID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
