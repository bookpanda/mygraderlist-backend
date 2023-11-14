package problem

import (
	"github.com/bookpanda/mygraderlist-backend/src/app/model"
)

type Problem struct {
	model.Base
	Group      string `json:"group" gorm:"type:tinytext"`
	Code       string `json:"code" gorm:"type:tinytext"`
	Name       string `json:"name" gorm:"type:tinytext"`
	CourseCode string `json:"course_code" gorm:"size:191"`
	Order      int    `json:"order" gorm:"type:tinyint"`
	// Course     *course.Course `json:"course" gorm:"foreignKey:CourseCode;references:CourseCode;constraint:OnUpdate:CASCADE;OnDelete:SET NULL;"`
}
