package repository

import "example.com/go-project/model"

type NecheRepository interface {
	Save(neche model.Neche) error
	FindAll() ([]model.Neche, error)
	FindById(id int) (*model.Neche, error)
	Delete(id int) error
}
