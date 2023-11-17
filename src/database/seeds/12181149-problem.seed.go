package seed

import "github.com/bookpanda/mygraderlist-backend/src/app/model/problem"

func (s Seed) ProblemSeed12181149() error {
	for _, b := range Problems {
		err := s.db.Model(&problem.Problem{}).Create(&b).Error

		if err != nil {
			return err
		}
	}
	return nil
}
