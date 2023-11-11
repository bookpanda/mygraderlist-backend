package rating

import (
	"os/user"

	"github.com/bookpanda/mygraderlist-backend/src/app/model"
	"github.com/bookpanda/mygraderlist-backend/src/app/model/problem"
	"github.com/google/uuid"
)

type Rating struct {
	model.Base
	Score      int              `json:"score" gorm:"type:tinyint"`
	Difficulty int              `json:"difficulty" gorm:"type:tinyint"`
	ProblemId  *uuid.UUID       `json:"problem_id"`
	Problem    *problem.Problem `json:"problem" gorm:"foreignKey:ProblemID,constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UserId     *uuid.UUID       `json:"user_id"`
	User       *user.User       `json:"user" gorm:"foreignKey:UserID,constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
