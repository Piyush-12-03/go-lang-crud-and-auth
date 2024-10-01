package services

import (
	"example.com/go-project/data/request"
	"example.com/go-project/data/response"
)

type TagsService interface {
	Create(tags request.CreteTagsRequest) error
	Update(tags request.UpdateTagsRequest) error
	Delete(tagId int) error
	FindById(tagsId int) (response.TagsResponse, error)
	FindAll(limit int, offset int) ([]response.TagsResponse, error)
	// FindPaginated(page int, pageSize int) []response.TagsResponse
	// FindAllSorted(sortBy string, order string) []response.TagsResponse
	// FindByCustomFilter(startsWith string) []response.TagsResponse
}
