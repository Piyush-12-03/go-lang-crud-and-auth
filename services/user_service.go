package services

import (
	"errors"

	"example.com/go-project/config"
	"example.com/go-project/model"
	"example.com/go-project/model/repository"
)

type UsersService struct {
	usersRepo *repository.UsersRepository
}

func NewUsersService(repo *repository.UsersRepository) *UsersService {
	return &UsersService{usersRepo: repo}
}

func (service *UsersService) Create(user model.Users) error {
	// Check if the email already exists
	existingUser, err := service.usersRepo.FindByEmail(user.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("email already exists")
	}

	// Hash the user's password before saving
	hashedPassword, err := config.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	// Save the user with the hashed password
	return service.usersRepo.Save(user)
}

func (service *UsersService) Authenticate(email string, password string) (*model.Users, error) {
	// Retrieve the user by email
	user, err := service.usersRepo.FindByEmail(email)
	if err != nil || user == nil {
		return nil, errors.New("invalid credentials")
	}

	// Compare the provided password with the stored hashed password
	if !config.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *UsersService) FindUserByEmail(email string) (*model.Users, error) {
	return s.usersRepo.FindByEmail(email)
}
