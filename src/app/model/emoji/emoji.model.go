package emoji

import (
	"github.com/bookpanda/mygraderlist-backend/src/app/model"
	"github.com/bookpanda/mygraderlist-backend/src/app/model/problem"
	"github.com/bookpanda/mygraderlist-backend/src/app/model/user"
	"github.com/google/uuid"
)

type Emoji struct {
	model.Base
	Emoji     string           `json:"emoji" gorm:"type:tinytext"`
	ProblemID *uuid.UUID       `json:"problem_id" gorm:"index:idx_name,unique"`
	Problem   *problem.Problem `json:"problem" gorm:"foreignKey:ProblemID;constraint:OnUpdate:CASCADE;OnDelete:SET NULL;"`
	UserID    *uuid.UUID       `json:"user_id" gorm:"index:idx_name,unique"`
	User      *user.User       `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE;OnDelete:SET NULL;"`
}
