package unittesting

import (
	"log"
	"testing"

	"example.com/go-project/model"
	"example.com/go-project/model/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDbForTagRepository() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
}

func TestTagsRepositoryImpl_Save_Success(t *testing.T) {
	log.Print("\n\n\n Running Tags Repository Test Cases.....\n\n\n")
	db, err := setupTestDbForTagRepository()
	if err != nil {
		t.Fatalf("could not connect to db: %v", err)
	}

	// Migrate the schema
	db.AutoMigrate(&model.Tags{})

	// Create repository instance
	repo := repository.NewTagsRepositoryImpl(db)

	// Test data
	mockTag := model.Tags{Id: 1, Name: "Sample Tag"}

	// Call the repository method
	err = repo.Save(mockTag)

	// Assertions
	assert.NoError(t, err)

	// Verify that the tag was saved
	var savedTag model.Tags
	db.First(&savedTag, 1) // Assuming ID = 1
	assert.Equal(t, mockTag.Name, savedTag.Name)
}

func TestTagsRepositoryImpl_Save_Failure(t *testing.T) {
	db, err := setupTestDbForTagRepository()
	if err != nil {
		t.Fatalf("could not connect to db: %v", err)
	}

	// Create repository instance
	repo := repository.NewTagsRepositoryImpl(db)

	// Simulate an error (like a unique constraint violation)
	db.AutoMigrate(&model.Tags{}) // Ensure the schema is migrated

	mockTag := model.Tags{Id: 1, Name: "Sample Tag"}
	// First save the tag
	err = repo.Save(mockTag)
	assert.NoError(t, err)

	// Try saving the same tag again to trigger a unique constraint error
	err = repo.Save(mockTag)
	assert.Error(t, err)
}

func TestTagsRepositoryImpl_Delete_Success(t *testing.T) {
	db, err := setupTestDbForTagRepository()
	if err != nil {
		t.Fatalf("could not connect to db: %v", err)
	}

	// Migrate the schema
	db.AutoMigrate(&model.Tags{})

	// Create a tag to delete
	db.Create(&model.Tags{Id: 1, Name: "Sample Tag"})

	// Create repository instance
	repo := repository.NewTagsRepositoryImpl(db)

	err = repo.Delete(1)

	assert.NoError(t, err)

	// Check if the tag is deleted
	var count int64
	db.Model(&model.Tags{}).Where("id = ?", 1).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestTagsRepositoryImpl_Delete_NotFound(t *testing.T) {
	db, err := setupTestDbForTagRepository()
	if err != nil {
		t.Fatalf("could not connect to db: %v", err)
	}

	// Migrate the schema
	db.AutoMigrate(&model.Tags{})

	// Create repository instance
	repo := repository.NewTagsRepositoryImpl(db)

	// Attempt to delete a non-existent tag
	err = repo.Delete(99)  // Assuming 99 doesn't exist
	assert.NoError(t, err) // No error should be returned even if it doesn't exist
}

func TestTagsRepositoryImpl_FindAll_Success(t *testing.T) {
	db, err := setupTestDbForTagRepository()
	if err != nil {
		t.Fatalf("could not connect to db: %v", err)
	}

	// Migrate the schema for both Tags and Neches
	db.AutoMigrate(&model.Tags{}, &model.Neche{}) // Ensure both models are migrated

	// Create repository instance
	repo := repository.NewTagsRepositoryImpl(db)

	// Add some tags and associated Neches for testing
	db.Create(&model.Tags{Id: 1, Name: "Tag 1"})
	db.Create(&model.Tags{Id: 2, Name: "Tag 2"})
	db.Create(&model.Neche{Id: 1, TagID: 1})
	db.Create(&model.Neche{Id: 2, TagID: 2})

	// Call the repository method
	tags, err := repo.FindAll(10, 0)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, tags, 2) // Expecting 2 tags
}

func TestTagsRepositoryImpl_FindAll_NoTags(t *testing.T) {
	db, err := setupTestDbForTagRepository()
	if err != nil {
		t.Fatalf("could not connect to db: %v", err)
	}

	// Migrate the schema
	db.AutoMigrate(&model.Tags{})

	// Create repository instance
	repo := repository.NewTagsRepositoryImpl(db)

	// Call the repository method with no tags
	tags, err := repo.FindAll(10, 0)

	// Assertions
	assert.NoError(t, err)
	assert.Empty(t, tags) // Expecting no tags
}

func TestTagsRepositoryImpl_FindById_Success(t *testing.T) {
	db, err := setupTestDbForTagRepository()
	if err != nil {
		t.Fatalf("could not connect to db: %v", err)
	}

	// Migrate the schema for both Tags and Neches
	db.AutoMigrate(&model.Tags{}, &model.Neche{})

	// Create a tag to find
	db.Create(&model.Tags{Id: 1, Name: "Sample Tag"})

	// Create repository instance
	repo := repository.NewTagsRepositoryImpl(db)

	tag, err := repo.FindById(1)

	assert.NoError(t, err)
	assert.Equal(t, "Sample Tag", tag.Name)
}

func TestTagsRepositoryImpl_FindById_NotFound(t *testing.T) {
	db, err := setupTestDbForTagRepository()
	if err != nil {
		t.Fatalf("could not connect to db: %v", err)
	}

	// Migrate the schema
	db.AutoMigrate(&model.Tags{})

	// Create repository instance
	repo := repository.NewTagsRepositoryImpl(db)

	tag, err := repo.FindById(99) // Assuming 99 doesn't exist

	assert.Error(t, err)
	assert.Equal(t, "tag not found", err.Error())
	assert.Equal(t, model.Tags{}, tag)
}

func TestTagsRepositoryImpl_Update_Success(t *testing.T) {
	db, err := setupTestDbForTagRepository()
	if err != nil {
		t.Fatalf("could not connect to db: %v", err)
	}

	// Migrate the schema
	db.AutoMigrate(&model.Tags{})

	// Create a tag to update
	db.Create(&model.Tags{Id: 1, Name: "Sample Tag"})

	// Create repository instance
	repo := repository.NewTagsRepositoryImpl(db)

	// Update tag
	updatedTag := model.Tags{Id: 1, Name: "Updated Tag"}
	err = repo.Update(updatedTag)

	assert.NoError(t, err)

	// Verify that the tag was updated
	var savedTag model.Tags
	db.First(&savedTag, 1)
	assert.Equal(t, "Updated Tag", savedTag.Name)
}

func TestTagsRepositoryImpl_Update_NotFound(t *testing.T) {
	db, err := setupTestDbForTagRepository()
	if err != nil {
		t.Fatalf("could not connect to db: %v", err)
	}

	// Migrate the schema for both Tags and Neches
	db.AutoMigrate(&model.Tags{}, &model.Neche{})

	// Create repository instance
	repo := repository.NewTagsRepositoryImpl(db)

	// Attempt to update a non-existent tag
	err = repo.Update(model.Tags{Id: 99, Name: "Non-Existent Tag"})

	assert.Error(t, err)
	assert.Equal(t, "tag not found", err.Error())
}
