package like

import (
	"github.com/bookpanda/mygraderlist-backend/src/app/model"
	"github.com/google/uuid"
)

type Like struct {
	model.Base
	ProblemId *uuid.UUID `json:"problem_id" gorm:"primaryKey"`
	UserId    *uuid.UUID `json:"user_id" gorm:"primaryKey"`
}
