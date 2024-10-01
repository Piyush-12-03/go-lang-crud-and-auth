package repository

import (
	"example.com/go-project/model"
	"gorm.io/gorm"
)

type UsersRepository struct {
	Db *gorm.DB
}

func NewUsersRepository(db *gorm.DB) *UsersRepository {
	return &UsersRepository{Db: db}
}

func (repo *UsersRepository) Save(user model.Users) error {
	result := repo.Db.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// FindByEmail finds a user by their email address
func (repo *UsersRepository) FindByEmail(email string) (*model.Users, error) {
	var user model.Users
	result := repo.Db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // No user found, return nil without error
		}
		return nil, result.Error // Return any other database error
	}
	return &user, nil // User found, return the user
}

