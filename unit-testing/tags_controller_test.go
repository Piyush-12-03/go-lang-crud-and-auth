package unittesting

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"example.com/go-project/controller"
	"example.com/go-project/data/request"
	"example.com/go-project/data/response"
	"example.com/go-project/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// mockTagsService implements services.TagsService for testing purposes.
type mockTagsService struct {
	tags map[int]struct{}
}

// Create implements services.TagsService.
func (m *mockTagsService) Create(tagsRequest request.CreteTagsRequest) error {
	// Simulate successful creation by adding to the mock tags map
	m.tags[len(m.tags)+1] = struct{}{}
	return nil
}

// Update implements services.TagsService.
func (m *mockTagsService) Update(tagsRequest request.UpdateTagsRequest) error {
	// Find the tag by ID
	if _, exists := m.tags[tagsRequest.Id]; !exists {
		return errors.New("Tag not found")
	}
	// update by just ensuring the tag exists
	m.tags[tagsRequest.Id] = struct{}{}
	return nil
}

// Delete implements services.TagsService.
func (m *mockTagsService) Delete(tagId int) error {
	if _, exists := m.tags[tagId]; !exists {
		return errors.New("Tag not found")
	}
	delete(m.tags, tagId)
	return nil
}

// FindById implements services.TagsService.
func (m *mockTagsService) FindById(tagId int) (response.TagsResponse, error) {
	if _, exists := m.tags[tagId]; !exists {
		return response.TagsResponse{}, errors.New("Tag not found")
	}
	return response.TagsResponse{Id: tagId, Name: "Tag " + strconv.Itoa(tagId)}, nil
}

// FindAll implements services.TagsService.
func (m *mockTagsService) FindAll(limit, offset int) ([]response.TagsResponse, error) {
	var allTags []response.TagsResponse
	for id := range m.tags {
		allTags = append(allTags, response.TagsResponse{Id: id, Name: "Tag " + strconv.Itoa(id)})
	}
	return allTags, nil
}

func setupTestDB() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
}

func TestCreateTag(t *testing.T) {
	// Set up the in-memory database
	db, err := setupTestDB()
	assert.NoError(t, err) // Assert no error

	// Migrate your models here (make sure to define your models)
	err = db.AutoMigrate(&model.Tags{})
	assert.NoError(t, err)

	// Create mock service
	tagsService := &mockTagsService{tags: make(map[int]struct{})}

	// Initialize the controller
	tagsController := controller.NewTagsController(tagsService)

	// Set up the router
	router := gin.Default()
	router.POST("/tags", tagsController.Create)

	// Create a test request
	createTagsRequest := request.CreteTagsRequest{Name: "Test Tag"}
	jsonRequest, _ := json.Marshal(createTagsRequest)

	log.Print("Requested tag to add : ", createTagsRequest.Name)
	// Create a request recorder
	req, _ := http.NewRequest(http.MethodPost, "/tags", bytes.NewBuffer(jsonRequest))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	// Perform the request
	router.ServeHTTP(recorder, req)

	// Assert the response code
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Check the response body
	var responseBody response.Response
	err = json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	assert.NoError(t, err)

	// Assert that the response data matches the request
	assert.Equal(t, createTagsRequest.Name, responseBody.Data.(map[string]interface{})["name"])
	log.Print("Test case for adding tag passed successfully.")
}

// Test for Delete method
func TestDeleteTag(t *testing.T) {
	// Set up the router
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Set up mock service with existing tags
	mockService := &mockTagsService{
		tags: map[int]struct{}{
			1: {},
		},
	}

	// Initialize the controller
	tagsController := controller.NewTagsController(mockService)

	// Set up the route
	router.DELETE("/tags/:tagId", tagsController.Delete)

	// Test deleting an existing tag
	req, _ := http.NewRequest(http.MethodDelete, "/tags/1", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Assert the response code for successful deletion
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Check the response body
	var responseBody response.Response
	err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Equal(t, "Deleted tag with id 1", responseBody.Data)

	// Test deleting a non-existing tag
	req, _ = http.NewRequest(http.MethodDelete, "/tags/999", nil)
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Assert the response code for not found
	assert.Equal(t, http.StatusNotFound, recorder.Code)

	// Check the response body for not found
	err = json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Equal(t, "Tag with id 999 not found", responseBody.Data)
}

func TestFindTagById(t *testing.T) {
	// Set up the router
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Set up mock service with existing tags
	mockService := &mockTagsService{
		tags: map[int]struct{}{
			1: {}, // Mock tag with ID 1 exists
		},
	}

	// Initialize the controller
	tagsController := controller.NewTagsController(mockService)

	// Set up the route
	router.GET("/tags/:tagId", tagsController.FindById)

	// Test finding an existing tag
	req, _ := http.NewRequest(http.MethodGet, "/tags/1", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Assert the response code for found tag
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Check the response body for the found tag
	var responseBody response.Response
	err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Equal(t, 1, int(responseBody.Data.(map[string]interface{})["id"].(float64)))
	assert.Equal(t, "Tag 1", responseBody.Data.(map[string]interface{})["name"])

	// Test finding a non-existing tag
	req, _ = http.NewRequest(http.MethodGet, "/tags/999", nil) // Tag ID 999 does not exist
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Assert the response code for not found
	assert.Equal(t, http.StatusNotFound, recorder.Code)

	// Check the response body for not found
	err = json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Equal(t, "Tag with id 999 not found", responseBody.Data)
}

func TestUpdateTag(t *testing.T) {
	// Set up the router
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Set up mock service with a tag already in the mock data
	mockService := &mockTagsService{
		tags: map[int]struct{}{
			1: {}, // Mock tag with ID 1 exists
		},
	}

	// Initialize the controller
	tagsController := controller.NewTagsController(mockService)

	// Set up the route
	router.PUT("/tags/:tagId", tagsController.Update)

	// Create an update request for a tag with Id 1
	updateTagsRequest := request.UpdateTagsRequest{
		Id:   1,
		Name: "Updated Tag",
	}
	jsonRequest, _ := json.Marshal(updateTagsRequest)

	// Test updating the existing tag with ID 1
	req, _ := http.NewRequest(http.MethodPut, "/tags/1", bytes.NewBuffer(jsonRequest))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Assert the response code for successful update
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Check the response body
	var responseBody response.Response
	err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	assert.NoError(t, err)

	// Assert that the response data matches the update request
	data := responseBody.Data.(map[string]interface{})

	// Assert the id and name
	assert.Equal(t, float64(1), data["id"].(float64)) // Ensure 'id' is present in the response
	assert.Equal(t, "Updated Tag", data["name"])
}

func TestFindAllTags(t *testing.T) {
	// Set up the router
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Set up mock service
	mockService := &mockTagsService{
		tags: map[int]struct{}{
			1: {},
			2: {},
		},
	}

	// Initialize the controller
	tagsController := controller.NewTagsController(mockService)

	// Set up the route
	router.GET("/tags", tagsController.FindAll)

	// Test finding all tags with pagination
	req, _ := http.NewRequest(http.MethodGet, "/tags?page=1&pageSize=2", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	// Assert the response code
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Check the response body
	var responseBody response.PaginatedResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	assert.NoError(t, err)

	// Assert that the correct number of tags are returned
	assert.Equal(t, 2, len(responseBody.Data.([]interface{})))

	// Convert float64 to int for comparison
	firstTagId := int(responseBody.Data.([]interface{})[0].(map[string]interface{})["id"].(float64))
	firstTagName := responseBody.Data.([]interface{})[0].(map[string]interface{})["name"].(string)

	// Assert the values are as expected
	assert.Equal(t, 1, firstTagId)
	assert.Equal(t, "Tag 1", firstTagName)
}
