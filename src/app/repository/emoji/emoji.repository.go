package course

import (
	"github.com/bookpanda/mygraderlist-backend/src/app/model/course"
	"github.com/bookpanda/mygraderlist-backend/src/app/model/emoji"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindAll(result *[]*emoji.Emoji) error {
	return r.db.Model(&course.Course{}).Find(result).Error
}

func (r *Repository) FindByUserId(userId string, result *[]*emoji.Emoji) error {
	return r.db.Model(&emoji.Emoji{}).Find(result, "user_id = ?", userId).Error
}

func (r *Repository) Create(in *emoji.Emoji) error {
	return r.db.Create(&in).Error
}

func (r *Repository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&emoji.Emoji{}).Error
}
