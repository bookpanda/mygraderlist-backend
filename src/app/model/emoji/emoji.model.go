package emoji

import (
	"github.com/bookpanda/mygraderlist-backend/src/app/model"
	"github.com/google/uuid"
)

type Emoji struct {
	model.Base
	Emoji     string     `json:"emoji" gorm:"type:tinytext"`
	ProblemId *uuid.UUID `json:"problem_id" gorm:"primaryKey"`
	UserId    *uuid.UUID `json:"user_id" gorm:"primaryKey"`
}
