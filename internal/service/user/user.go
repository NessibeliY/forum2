package user

import (
	"errors"

	"forum/internal/models"
	"forum/internal/service/helpers"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	Insert(user models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserUserID(id string) (*models.User, error)
}

type UserService struct {
	UserRepo UserRepo
}

func NewUserService(userRepo UserRepo) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

type IUserService interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	Login(email, password string) (*models.User, error)
	GetUserUserID(id string) (*models.User, error)
}

// create new user
func (u *UserService) CreateUser(user *models.User) error {
	// check user is not nil
	if user == nil {
		return errors.New("user cannot be nil")
	}

	// generate unique ID
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}

	// hash password
	hashPass, err := helpers.HashPassword(user.Password)
	if err != nil {
		return err
	}

	// new model user
	userModel := models.User{
		ID:       id.String(),
		UserName: user.UserName,
		Email:    user.Email,
		Password: hashPass,
	}

	// create new user in db
	if err := u.UserRepo.Insert(userModel); err != nil {
		return err
	}

	return nil
}

// get user by email
func (u *UserService) GetUserByEmail(email string) (*models.User, error) {
	user, err := u.UserRepo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// login user
func (u *UserService) Login(email, password string) (*models.User, error) {
	user, err := u.GetUserByEmail(email)
	if err != nil {
		return nil, models.ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, models.ErrWrongPassword
	}

	return user, nil
}

// get user by ID
func (u *UserService) GetUserUserID(id string) (*models.User, error) {
	user, err := u.UserRepo.GetUserUserID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
