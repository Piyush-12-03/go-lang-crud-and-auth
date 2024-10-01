package services

import (
	"errors"

	"example.com/go-project/data/request"
	"example.com/go-project/data/response"
	"example.com/go-project/model"
	"example.com/go-project/model/repository"
	"github.com/go-playground/validator"
)

type TagsServiceImpl struct {
	TagsRepository repository.TagsRepository
	validate       *validator.Validate
}

func NewTagsServiceImpl(tagsRepository repository.TagsRepository, validate *validator.Validate) TagsService {
	return &TagsServiceImpl{
		TagsRepository: tagsRepository,
		validate:       validate,
	}
}

// Create implements TagsService.
func (t *TagsServiceImpl) Create(tags request.CreteTagsRequest) error {
	// Validate the incoming request
	err := t.validate.Struct(tags)
	if err != nil {
		return err // Return validation error if any
	}

	// Create the tag model
	tagModel := model.Tags{
		Name: tags.Name,
	}

	// Save the tag model to the database
	err = t.TagsRepository.Save(tagModel)
	if err != nil {
		return err // Return database save error if any
	}

	return nil // Return nil if creation is successful
}

// Delete implements TagsService.
func (t *TagsServiceImpl) Delete(tagId int) error {
	// Call the repository's Delete method and return its result
	return t.TagsRepository.Delete(tagId) // Return the error if any
}

// FindAll implements TagsService.
func (t *TagsServiceImpl) FindAll(limit int, offset int) ([]response.TagsResponse, error) {
	// Fetch paginated results from the repository
	result, err := t.TagsRepository.FindAll(limit, offset)
	if err != nil {
		return nil, err // Return nil and the error if fetching fails
	}

	var tags []response.TagsResponse

	for _, value := range result {
		var necheResponses []response.NecheResponse
		for _, neche := range value.Neches {
			necheResponse := response.NecheResponse{
				Id:     neche.Id,
				Name:   neche.NecheType,
				TagsID: neche.TagID,
			}
			necheResponses = append(necheResponses, necheResponse)
		}

		tag := response.TagsResponse{
			Id:     value.Id,
			Name:   value.Name,
			Neches: necheResponses,
		}
		tags = append(tags, tag) // Append each tag to the tags slice
	}

	return tags, nil // Return the list of tags and nil for error
}

// FindById implements TagsService.
func (t *TagsServiceImpl) FindById(tagsId int) (response.TagsResponse, error) {
	tagData, err := t.TagsRepository.FindById(tagsId)
	if err != nil {
		return response.TagsResponse{}, err // Return empty response and error if not found
	}
	tagResponse := response.TagsResponse{
		Id:   tagData.Id,
		Name: tagData.Name,
	}
	return tagResponse, nil // Return the tag response and nil for error
}

// Update implements TagsService.
func (t *TagsServiceImpl) Update(tags request.UpdateTagsRequest) error {
	// Validation for empty name
	if tags.Name == "" {
		return errors.New("tag name cannot be empty")
	}

	// Validation for non-positive ID
	if tags.Id <= 0 {
		return errors.New("tag ID must be positive")
	}

	// Attempt to find the existing tag by ID
	tagsData, err := t.TagsRepository.FindById(tags.Id)
	if err != nil {
		// Explicitly return the "tag not found" error from the repository
		return err
	}

	// If the tag was found but is empty, return an error
	if tagsData.Id == 0 { // This should check if the tag does not exist
		return errors.New("tag not found")
	}

	// Update the tag name
	tagsData.Name = tags.Name

	// Call the repository's Update method and handle the error
	err = t.TagsRepository.Update(tagsData)
	if err != nil {
		return err // Return the error if the update fails
	}

	return nil // Return nil if everything is successful
}

// // FindPaginated implements TagsService.
// func (t *TagsServiceImpl) FindPaginated(page int, pageSize int) []response.TagsResponse {
// 	offset := (page - 1) * pageSize
// 	result := t.TagsRepository.FindPaginated(pageSize, offset)
// 	var tags []response.TagsResponse
// 	for _, value := range result {
// 		tag := response.TagsResponse{
// 			Id:   value.Id,
// 			Name: value.Name,
// 		}
// 		tags = append(tags, tag)
// 	}
// 	return tags
// }

// // FindAllSorted implements TagsService.
// func (t *TagsServiceImpl) FindAllSorted(sortBy string, order string) []response.TagsResponse {
// 	result := t.TagsRepository.FindAllSorted(sortBy, order)
// 	var tags []response.TagsResponse
// 	for _, value := range result {
// 		tag := response.TagsResponse{
// 			Id:   value.Id,
// 			Name: value.Name,
// 		}
// 		tags = append(tags, tag)
// 	}
// 	return tags
// }

// // FindByCustomFilter implements TagsService.
// func (t *TagsServiceImpl) FindByCustomFilter(startsWith string) []response.TagsResponse {
// 	result := t.TagsRepository.FindByCustomFilter(startsWith)
// 	var tags []response.TagsResponse
// 	for _, value := range result {
// 		tag := response.TagsResponse{
// 			Id:   value.Id,
// 			Name: value.Name,
// 		}
// 		tags = append(tags, tag)
// 	}
// 	return tags
// }
