package course

import (
	"github.com/bookpanda/mygraderlist-backend/src/app/model/course"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindAllCourse(result *[]*course.Course) error {
	return r.db.Model(&course.Course{}).Find(result).Error
}

func (r *Repository) Create(in *course.Course) error {
	return r.db.Create(&in).Error
}

func (r *Repository) Update(id string, result *course.Course) error {
	return r.db.Where(id, "id = ?", id).Updates(&result).First(&result, "id = ?", id).Error
}

func (r *Repository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&course.Course{}).Error
}
