package seed

import "github.com/bookpanda/mygraderlist-backend/src/app/model/course"

func (s Seed) CourseSeed12180464() error {
	for _, b := range courses {
		err := s.db.Model(&course.Course{}).Create(&b).Error

		if err != nil {
			return err
		}
	}
	return nil
}
