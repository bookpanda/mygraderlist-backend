package course

import (
	"github.com/bookpanda/mygraderlist-backend/src/app/model"
)

type Course struct {
	model.Base
	CourseCode string `json:"course_code" gorm:"size:191"`
	Name       string `json:"name" gorm:"type:mediumtext"`
	Color      string `json:"color" gorm:"type:tinytext"`
	// Problems   []*problem.Problem `gorm:"foreignKey:CourseCode;constraint:OnUpdate:CASCADE;OnDelete:SET NULL;"`
}
