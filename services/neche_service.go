package services

import (
	"example.com/go-project/data/request"
	"example.com/go-project/model"
)

type NecheService interface {
	Create(neche request.CreateNecheRequest) error
	FindAll() ([]model.Neche, error)
	FindById(id int) (*model.Neche, error)
	Delete(id int) error
}