package database

import (
	"fmt"
	"os/user"
	"strconv"

	"github.com/bookpanda/mygraderlist-backend/src/app/model/course"
	"github.com/bookpanda/mygraderlist-backend/src/app/model/emoji"
	"github.com/bookpanda/mygraderlist-backend/src/app/model/like"
	"github.com/bookpanda/mygraderlist-backend/src/app/model/problem"
	"github.com/bookpanda/mygraderlist-backend/src/app/model/rating"
	"github.com/bookpanda/mygraderlist-backend/src/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDatabase(conf *config.Database) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", conf.User, conf.Password, conf.Host, strconv.Itoa(conf.Port), conf.Name)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&user.User{}, &course.Course{})
	err = db.AutoMigrate(&problem.Problem{})
	err = db.AutoMigrate(&rating.Rating{}, &emoji.Emoji{}, &like.Like{})
	if err != nil {
		return nil, err
	}

	return
}
