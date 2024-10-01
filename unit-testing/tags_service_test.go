package unittesting

import (
	"errors"
	"log"
	"testing"

	"example.com/go-project/data/request"
	"example.com/go-project/model"
	"example.com/go-project/model/repository"
	"example.com/go-project/services"
	"github.com/go-playground/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDbForTagService() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
}

// Define the MockTagsRepository
type MockTagsRepository struct {
	mock.Mock
}

// Save implements repository.TagsRepository.
func (m *MockTagsRepository) Save(tags model.Tags) error {
	args := m.Called(tags)
	return args.Error(0)
}

// Update implements repository.TagsRepository.
func (m *MockTagsRepository) Update(tags model.Tags) error {
	args := m.Called(tags)
	return args.Error(0)
}

// Delete implements repository.TagsRepository.
func (m *MockTagsRepository) Delete(tagsId int) error {
	args := m.Called(tagsId)
	return args.Error(0)
}

// FindById implements repository.TagsRepository.
func (m *MockTagsRepository) FindById(tagsId int) (model.Tags, error) {
	args := m.Called(tagsId)
	return args.Get(0).(model.Tags), args.Error(1)
}

// FindAll implements repository.TagsRepository.
func (m *MockTagsRepository) FindAll(limit int, offset int) ([]model.Tags, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]model.Tags), args.Error(1)
}
func TestCreateTagService(t *testing.T) {
	// Set up the in-memory database (SQLite or similar)
	log.Print("\n\n\n Running Tags Service Test Cases.....\n\n\n")
	db, err := setupTestDbForTagService()
	assert.NoError(t, err, "Failed to set up test database")

	// Migrate your models here (make sure to define your models)
	err = db.AutoMigrate(&model.Tags{})
	assert.NoError(t, err, "Failed to migrate models")

	// Create mock repository and validator
	mockTagsRepository := repository.NewTagsRepositoryImpl(db)
	validate := validator.New()

	// Initialize the service
	tagsService := services.NewTagsServiceImpl(mockTagsRepository, validate)

	// Create a test request
	createTagsRequest := request.CreteTagsRequest{
		Name: "Test Tag",
	}

	log.Print("Requested tag to add: ", createTagsRequest.Name)

	// Call the service's Create method directly
	err = tagsService.Create(createTagsRequest)
	assert.NoError(t, err, "Service failed to create tag")

	// Fetch the tag from the database to verify it was created
	var createdTag model.Tags
	err = db.Where("name = ?", createTagsRequest.Name).First(&createdTag).Error
	assert.NoError(t, err, "Failed to fetch created tag from the database")

	// Assert the tag was created with the correct name
	assert.Equal(t, createTagsRequest.Name, createdTag.Name, "Tag name does not match")

	log.Print("Test case for adding tag passed successfully.")
}

func TestDeleteTagService(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockTagsRepository)

	// Initialize service
	tagsService := services.NewTagsServiceImpl(mockRepo, validator.New())

	// Set up expectations
	mockRepo.On("Delete", 1).Return(nil)

	// Test deleting a tag
	tagsService.Delete(1)

	// Assert that the repository Delete method was called
	mockRepo.AssertCalled(t, "Delete", 1)
	log.Print("Deleted Tag Test Case Passed.")
}

func TestDeleteTagService_NotFound(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockTagsRepository)

	// Initialize service
	tagsService := services.NewTagsServiceImpl(mockRepo, validator.New())

	// Set up expectations for non-existing tag
	mockRepo.On("Delete", 1).Return(errors.New("tag not found"))

	err := tagsService.Delete(1)

	// Assert the error is as expected
	assert.Error(t, err, "tag not found")
	log.Print("Tag Not Found Test Case Passed.")
}

func TestFindAllTagsService(t *testing.T) {
	mockRepo := new(MockTagsRepository)
	tagsService := services.NewTagsServiceImpl(mockRepo, validator.New())

	// Mock data
	mockTags := []model.Tags{
		{Id: 1, Name: "Tag 1"},
		{Id: 2, Name: "Tag 2"},
	}
	// Mock the repository method to return the mock data
	mockRepo.On("FindAll", 10, 0).Return(mockTags, nil)

	// Test fetching tags
	tags, err := tagsService.FindAll(10, 0)

	// Assert no error occurred
	assert.NoError(t, err)

	// Assert correct data is returned
	assert.Len(t, tags, 2)
	assert.Equal(t, "Tag 1", tags[0].Name)
	assert.Equal(t, "Tag 2", tags[1].Name)

	log.Print("Fetched All Tags Test Case Passed.")
}

func TestFindByIdTagsService(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockTagsRepository)

	// Initialize service
	tagsService := services.NewTagsServiceImpl(mockRepo, validator.New())
	mockTag := model.Tags{Id: 1, Name: "Tag 1"}
	// Set up expectations for the FindById method
	mockRepo.On("FindById", 1).Return(mockTag, nil)

	// Call the FindById method on the service
	tagResponse, err := tagsService.FindById(1)

	// Assert that there is no error
	assert.NoError(t, err, "Expected no error when finding tag")

	// Assert that the returned tag matches the expected values
	assert.Equal(t, 1, tagResponse.Id, "Tag ID does not match")
	assert.Equal(t, "Tag 1", tagResponse.Name, "Tag name does not match")

	log.Print("Get Tag By Id Test Case Passed.")
}

func TestFindByIdNotFoundTagsService(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockTagsRepository)

	// Initialize service
	tagsService := services.NewTagsServiceImpl(mockRepo, validator.New())

	// Set up expectations for non-existing tag
	mockRepo.On("FindById", 1).Return(model.Tags{}, errors.New("tag not found"))

	// Call FindById and expect an error
	tagResponse, err := tagsService.FindById(1)

	// Assert the error is as expected
	assert.Error(t, err)
	assert.EqualError(t, err, "tag not found")
	assert.Empty(t, tagResponse)

	log.Print("Tag Not Found Test Case Passed.")
}

func TestUpdateTagService_Success(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockTagsRepository)

	// Initialize service
	tagsService := services.NewTagsServiceImpl(mockRepo, validator.New())

	// Create mock tag data
	mockTag := model.Tags{Id: 1, Name: "Old Tag Name"}
	updateRequest := request.UpdateTagsRequest{Id: 1, Name: "New Tag Name"}

	// Set up expectations for FindById
	mockRepo.On("FindById", updateRequest.Id).Return(mockTag, nil)

	// Update the tag name before calling Update
	updatedTag := mockTag
	updatedTag.Name = updateRequest.Name

	// Set up expectations for Update with the updated tag name
	mockRepo.On("Update", updatedTag).Return(nil)

	// Test updating a tag
	err := tagsService.Update(updateRequest)

	// Assert that there are no errors
	assert.NoError(t, err)

	// Assert that the repository methods were called
	mockRepo.AssertCalled(t, "FindById", updateRequest.Id)
	mockRepo.AssertCalled(t, "Update", updatedTag)

	log.Print("Update Tag Success Test Case Passed.")
}

func TestUpdateTagService_TagNotFound(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockTagsRepository)

	// Initialize service
	tagsService := services.NewTagsServiceImpl(mockRepo, validator.New())

	// Create update request
	updateRequest := request.UpdateTagsRequest{Id: 1, Name: "New Tag Name"}

	// Set up expectations for non-existing tag
	mockRepo.On("FindById", updateRequest.Id).Return(model.Tags{}, errors.New("tag not found"))

	// Test updating a tag
	err := tagsService.Update(updateRequest)

	// Assert the error is as expected
	assert.Error(t, err, "tag not found")

	log.Print("Update Tag Not Found Test Case Passed.")
}

func TestUpdateTagService_ValidationError(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockTagsRepository)

	// Initialize service
	tagsService := services.NewTagsServiceImpl(mockRepo, validator.New())

	// Create update request with an empty name
	updateRequest := request.UpdateTagsRequest{Id: 1, Name: ""}

	// Test updating a tag
	err := tagsService.Update(updateRequest)

	// Assert validation error
	assert.Error(t, err, "tag name cannot be empty")

	log.Print("Validation Error Test Case Passed.")
}

func TestUpdateTagService_UpdateFails(t *testing.T) {
	// Create mock repository
	mockRepo := new(MockTagsRepository)

	// Initialize service
	tagsService := services.NewTagsServiceImpl(mockRepo, validator.New())

	// Create mock tag data
	mockTag := model.Tags{Id: 1, Name: "Old Tag Name"}
	updateRequest := request.UpdateTagsRequest{Id: 1, Name: "New Tag Name"}

	// Set up expectations for FindById
	mockRepo.On("FindById", updateRequest.Id).Return(mockTag, nil)

	// Update the tag name before calling Update
	updatedTag := mockTag
	updatedTag.Name = updateRequest.Name

	// Set up expectations for Update failure
	mockRepo.On("Update", updatedTag).Return(errors.New("failed to update tag"))

	// Test updating a tag
	err := tagsService.Update(updateRequest)

	// Assert the error is as expected
	assert.Error(t, err, "failed to update tag")

	// Assert that the repository methods were called
	mockRepo.AssertCalled(t, "FindById", updateRequest.Id)
	mockRepo.AssertCalled(t, "Update", updatedTag)

	log.Print("Update Fails Test Case Passed.")
}

