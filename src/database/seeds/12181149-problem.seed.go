package seed

func (s Seed) ProblemSeed12181149() error {
	for _, b := range problems {
		err := s.db.Create(&b).Error

		if err != nil {
			return err
		}
	}
	return nil
}
