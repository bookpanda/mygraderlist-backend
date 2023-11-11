package course

import (
	"github.com/bookpanda/mygraderlist-backend/src/app/model"
)

type Course struct {
	model.Base
	CourseCode string `json:"course_code" gorm:"type:tinytext"`
	Name       string `json:"name" gorm:"type:mediumtext"`
	Color      string `json:"color" gorm:"type:tinytext"`
}
