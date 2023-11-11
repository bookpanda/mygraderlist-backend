package rating

import (
	"github.com/bookpanda/mygraderlist-backend/src/app/model"
	"github.com/google/uuid"
)

type Rating struct {
	model.Base
	ProblemId  uuid.UUID `json:"problem_id" gorm:"primaryKey"`
	UserId     uuid.UUID `json:"user_id" gorm:"primaryKey"`
	Score      int       `json:"score" gorm:"type:tinyint"`
	Difficulty int       `json:"difficulty" gorm:"type:tinyint"`
}
