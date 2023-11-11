package problem

import (
	"github.com/bookpanda/mygraderlist-backend/src/app/model"
	"github.com/google/uuid"
)

type Problem struct {
	model.Base
	CourseId *uuid.UUID `json:"course_id" gorm:"primaryKey"`
	Group    string     `json:"group" gorm:"type:tinytext"`
	Code     string     `json:"code" gorm:"type:tinytext"`
	Name     string     `json:"name" gorm:"type:tinytext"`
}
