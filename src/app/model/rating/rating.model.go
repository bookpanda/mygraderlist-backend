package rating

import (
	"github.com/bookpanda/mygraderlist-backend/src/app/model"
	"github.com/bookpanda/mygraderlist-backend/src/app/model/problem"
	"github.com/bookpanda/mygraderlist-backend/src/app/model/user"
	"github.com/google/uuid"
)

type Rating struct {
	model.Base
	Score      int              `json:"score" gorm:"type:tinyint"`
	Difficulty int              `json:"difficulty" gorm:"type:tinyint"`
	ProblemID  *uuid.UUID       `json:"problem_id" gorm:"index:idx_name,unique"`
	Problem    *problem.Problem `json:"problem" gorm:"foreignKey:ProblemID;constraint:OnUpdate:CASCADE;OnDelete:SET NULL;"`
	UserID     *uuid.UUID       `json:"user_id" gorm:"index:idx_name,unique"`
	User       *user.User       `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE;OnDelete:SET NULL;"`
}
