package repository

import (
	"errors"

	"example.com/go-project/model"
	"gorm.io/gorm"
)

type TagsRepositoryImpl struct {
	Db *gorm.DB
}

func NewTagsRepositoryImpl(Db *gorm.DB) TagsRepository {
	return &TagsRepositoryImpl{Db: Db}
}

func (t *TagsRepositoryImpl) Delete(tagId int) error {
	var tags model.Tags
	result := t.Db.Where("id = ?", tagId).Delete(&tags)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (t *TagsRepositoryImpl) FindAll(limit int, offset int) ([]model.Tags, error) {
	var tags []model.Tags
	// Preload Neches and apply pagination using Limit and Offset
	result := t.Db.Preload("Neches").Limit(limit).Offset(offset).Find(&tags)

	if result.Error != nil {
		return nil, result.Error // Return nil and the error if fetching fails
	}

	return tags, nil // Return the fetched tags and nil for error
}

// FindById method to fetch a Tag along with its associated Neches
func (r *TagsRepositoryImpl) FindById(tagsId int) (model.Tags, error) {
	var tag model.Tags
	// Use Preload to eagerly load Neches and First to get a single record by ID
	result := r.Db.Preload("Neches").First(&tag, tagsId)

	// Check if no rows were found
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return tag, errors.New("tag not found")
	}

	return tag, result.Error
}

func (t *TagsRepositoryImpl) Save(tags model.Tags) error {
	// Try to create the tag in the database
	result := t.Db.Create(&tags)

	// Check for any errors and return them
	if result.Error != nil {
		return result.Error
	}

	// If successful, return nil
	return nil
}

func (t *TagsRepositoryImpl) Update(tags model.Tags) error {
	// Check if the tag exists
	var existingTag model.Tags
	if err := t.Db.First(&existingTag, tags.Id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("tag not found")
		}
		return err
	}

	result := t.Db.Model(&existingTag).Updates(tags)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// // Repository method for paginated retrieval
// func (r *TagsRepositoryImpl) FindPaginated(limit int, offset int) []model.Tags {
// 	var tags []model.Tags
// 	// Query to fetch paginated tags
// 	r.Db.Limit(limit).Offset(offset).Find(&tags)
// 	return tags
// }

// // Repository method for sorted retrieval
// func (r *TagsRepositoryImpl) FindAllSorted(sortBy string, order string) []model.Tags {
// 	var tags []model.Tags
// 	// Dynamic query to sort by a specific field
// 	r.Db.Order(fmt.Sprintf("%s %s", sortBy, order)).Find(&tags)
// 	return tags
// }

// func (r *TagsRepositoryImpl) FindByCustomFilter(startsWith string) []model.Tags {
// 	var tags []model.Tags

// 	err := r.Db.Raw("SELECT * FROM tags WHERE name LIKE ?", startsWith+"%").Scan(&tags).Error
// 	if err != nil {
// 		return []model.Tags{}
// 	}

// 	return tags
// }
