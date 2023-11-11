package user

import (
	"github.com/bookpanda/mygraderlist-backend/src/app/model"
)

type User struct {
	model.Base
	Username string `json:"username" gorm:"type:tinytext"`
	Email    string `json:"email" gorm:"type:tinytext"`
	Password string `json:"password" gorm:"type:tinytext"`
}
