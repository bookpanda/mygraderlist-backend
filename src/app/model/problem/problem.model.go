package problem

import (
	"github.com/bookpanda/mygraderlist-backend/src/app/model"
	"github.com/bookpanda/mygraderlist-backend/src/app/model/course"
	"github.com/google/uuid"
)

type Problem struct {
	model.Base
	Group    string         `json:"group" gorm:"type:tinytext"`
	Code     string         `json:"code" gorm:"type:tinytext"`
	Name     string         `json:"name" gorm:"type:tinytext"`
	CourseID *uuid.UUID     `json:"course_id"`
	Course   *course.Course `json:"course" gorm:"foreignKey:CourseID,constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
