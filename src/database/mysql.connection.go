package database

import (
	"fmt"
	"log"
	"os/user"
	"strconv"

	"github.com/bookpanda/mygraderlist-backend/src/app/model/course"
	"github.com/bookpanda/mygraderlist-backend/src/app/model/emoji"
	"github.com/bookpanda/mygraderlist-backend/src/app/model/like"
	"github.com/bookpanda/mygraderlist-backend/src/app/model/problem"
	"github.com/bookpanda/mygraderlist-backend/src/app/model/rating"
	"github.com/bookpanda/mygraderlist-backend/src/config"
	seed "github.com/bookpanda/mygraderlist-backend/src/database/seeds"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDatabase(conf *config.Database) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True", conf.User, conf.Password, conf.Host, strconv.Itoa(conf.Port), conf.Name)

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

	var count int64
	for _, b := range seed.Courses {
		if db.Model(&course.Course{}).Where("course_code = ?", b.CourseCode).Count(&count); count == 0 {
			err := db.Create(&b).Error
			if err != nil {
				return nil, err
			}
		} else {
			err := db.Where("course_code = ?", b.CourseCode).Updates(&b).Error
			if err != nil {
				return nil, err
			}
		}
	}
	log.Println("✔️Seed", "courses", "succeed")

	for _, b := range seed.Problems {
		if db.Model(&problem.Problem{}).Where("code = ?", b.Code).Count(&count); count == 0 {
			err := db.Create(&b).Error
			if err != nil {
				return nil, err
			}
		} else {
			err := db.Where("code = ?", b.Code).Updates(&b).Error
			if err != nil {
				return nil, err
			}
		}
	}
	log.Println("✔️Seed", "problems", "succeed")

	return
}
