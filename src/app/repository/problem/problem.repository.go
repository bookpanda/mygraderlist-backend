package problem

import (
	"github.com/bookpanda/mygraderlist-backend/src/app/model/problem"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindAllProblem(result *[]*problem.Problem) error {
	return r.db.Model(&problem.Problem{}).Find(result).Error
}

func (r *Repository) Create(in *problem.Problem) error {
	return r.db.Create(&in).Error
}

func (r *Repository) Update(id string, result *problem.Problem) error {
	return r.db.Where(id, "id = ?", id).Updates(&result).First(&result, "id = ?", id).Error
}

func (r *Repository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&problem.Problem{}).Error
}
