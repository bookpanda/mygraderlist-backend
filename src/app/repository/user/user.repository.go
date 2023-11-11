package user

import (
	"github.com/bookpanda/mygraderlist-backend/src/app/model/user"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindOne(id string, result *user.User) error {
	return r.db.First(&result, "id = ?", id).Error
}

func (r *Repository) FindByEmail(email string, result *user.User) error {
	return r.db.Model(&user.User{}).First(&result, "email = ?", email).Error
}

func (r *Repository) Create(in *user.User) error {
	return r.db.Create(&in).Error
}

func (r *Repository) Update(id string, result *user.User) error {
	return r.db.Where(id, "id = ?", id).Updates(&result).First(&result, "id = ?", id).Error
}

func (r *Repository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&user.User{}).Error
}
