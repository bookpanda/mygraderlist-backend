package course

import (
	"github.com/bookpanda/mygraderlist-backend/src/app/model"
)

type Course struct {
	model.Base
	Course string `json:"course" gorm:"type:tinytext"`
	Name   string `json:"name" gorm:"type:mediumtext"`
	Color  string `json:"color" gorm:"type:tinytext"`
}
