package repository

import (
	"example.com/go-project/model"
	"gorm.io/gorm"
)

type NecheRepositoryImpl struct {
	Db *gorm.DB
}

func NewNecheRepositoryImpl(Db *gorm.DB) NecheRepository {
	return &NecheRepositoryImpl{Db: Db}
}

// Save Neche
func (n *NecheRepositoryImpl) Save(neche model.Neche) error {
	result := n.Db.Create(&neche)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Find all Neches
func (n *NecheRepositoryImpl) FindAll() ([]model.Neche, error) {
	var neches []model.Neche
	result := n.Db.Find(&neches)
	if result.Error != nil {
		return nil, result.Error
	}
	return neches, nil
}

// Find Neche by ID
func (n *NecheRepositoryImpl) FindById(id int) (*model.Neche, error) {
	var neche model.Neche
	result := n.Db.First(&neche, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &neche, nil
}

// Delete Neche by ID
func (n *NecheRepositoryImpl) Delete(id int) error {
	result := n.Db.Delete(&model.Neche{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
