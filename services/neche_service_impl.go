package services

import (
	"fmt"
	"example.com/go-project/data/request"
	"example.com/go-project/model"
	"example.com/go-project/model/repository"
	"github.com/go-playground/validator"
)

type NecheServiceImpl struct {
	NecheRepository repository.NecheRepository
	validate        *validator.Validate
	TagsRepository  repository.TagsRepository
}

func NewNecheServiceImpl(necheRepository repository.NecheRepository, validate *validator.Validate, tagsRepository repository.TagsRepository) NecheService {
	return &NecheServiceImpl{
		NecheRepository: necheRepository,
		validate:        validate,
		TagsRepository:  tagsRepository,
	}
}

// Create Neche
func (n *NecheServiceImpl) Create(necheReq request.CreateNecheRequest) error {
	err := n.validate.Struct(necheReq)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	tag, err := n.TagsRepository.FindById(necheReq.TagID)
	if err != nil {
		return fmt.Errorf("tag with ID %d not found", necheReq.TagID)
	}

	neche := model.Neche{
		TagID:     tag.Id,
		NecheType: necheReq.Name,
	}

	return n.NecheRepository.Save(neche)
}

// Find all Neches
func (n *NecheServiceImpl) FindAll() ([]model.Neche, error) {
	return n.NecheRepository.FindAll()
}

// Find Neche by ID
func (n *NecheServiceImpl) FindById(id int) (*model.Neche, error) {
	return n.NecheRepository.FindById(id)
}

// Delete Neche by ID
func (n *NecheServiceImpl) Delete(id int) error {
	return n.NecheRepository.Delete(id)
}
