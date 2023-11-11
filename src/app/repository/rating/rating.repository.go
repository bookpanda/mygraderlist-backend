package rating

import (
	"github.com/bookpanda/mygraderlist-backend/src/app/model/rating"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindAll(result *[]*rating.Rating) error {
	return r.db.Model(&rating.Rating{}).Find(result).Error
}

func (r *Repository) FindByUserId(userId string, result *[]*rating.Rating) error {
	return r.db.Model(&rating.Rating{}).Find(result, "user_id = ?", userId).Error
}

func (r *Repository) Create(in *rating.Rating) error {
	return r.db.Create(&in).Error
}

func (r *Repository) Update(id string, result *rating.Rating) error {
	return r.db.Where(id, "id = ?", id).Updates(&result).First(&result, "id = ?", id).Error
}

func (r *Repository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&rating.Rating{}).Error
}
