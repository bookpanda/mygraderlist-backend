package seed

func (s Seed) CourseSeed12180464() error {
	for _, b := range courses {
		err := s.db.Save(&b).Error

		if err != nil {
			return err
		}
	}
	return nil
}
