package repository

import "example.com/go-project/model"

type TagsRepository interface {
	Save(tags model.Tags) error
	Update(tags model.Tags) error
	Delete(tagsId int) error
	FindById(tagsId int) (tags model.Tags, err error)
	FindAll(int, int) ([]model.Tags, error)
	// FindPaginated(limit int, offset int) []model.Tags
	// FindAllSorted(sortBy string, order string) []model.Tags
	// FindByCustomFilter(startsWith string) []model.Tags
}
