package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

type Bookmark struct {
	gorm.Model
	Name string `json:"name"`
	Url  string `json:"url"`
}

func InitDatabase() error {
	var err error
	db, err = gorm.Open(sqlite.Open("bookmark.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&Bookmark{})
	if err != nil {
		return err
	}

	return nil
}

// GetAllBookmarks gets all bookmarks
func GetAllBookmarks() ([]Bookmark, error) {
	var bookmarks []Bookmark

	err := db.Find(&bookmarks).Error
	if err != nil {
		return bookmarks, err
	}

	return bookmarks, nil
}

// GetBookmark gets a bookmark
func GetBookmark(id int) (Bookmark, error) {
	var bookmark Bookmark

	err := db.First(&bookmark, id).Error
	if err != nil {
		return bookmark, err
	}

	return bookmark, nil
}

// GetBookmarkByName gets a bookmark by name
func GetBookmarkByName(name string) (Bookmark, error) {
	var bookmark Bookmark

	err := db.Where("name = ?", name).First(&bookmark).Error
	if err != nil {
		return bookmark, err
	}

	return bookmark, nil
}

// CreateBookmark creates a bookmark
func CreateBookmark(bookmark Bookmark) error {
	err := db.Create(&bookmark).Error
	if err != nil {
		return err
	}

	return nil
}

// UpdateBookmark updates a bookmark
func UpdateBookmark(bookmark Bookmark) error {
	err := db.Save(&bookmark).Error
	if err != nil {
		return err
	}

	return nil
}

// DeleteBookmark deletes a bookmark
func DeleteBookmark(bookmark Bookmark) error {
	err := db.Delete(&bookmark).Error
	if err != nil {
		return err
	}

	return nil
}
