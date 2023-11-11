package course

import (
	"github.com/bookpanda/mygraderlist-backend/src/app/model/course"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) FindAll(in *[]*course.Course) error {
	args := r.Called(in)

	if args.Get(0) != nil {
		*in = *args.Get(0).(*[]*course.Course)
	}

	return args.Error(1)
}

func (r *RepositoryMock) Create(in *course.Course) error {
	args := r.Called(in)

	if args.Get(0) != nil {
		*in = *args.Get(0).(*course.Course)
	}

	return args.Error(1)
}

func (r *RepositoryMock) Update(id string, result *course.Course) error {
	args := r.Called(id, result)

	if args.Get(0) != nil {
		*result = *args.Get(0).(*course.Course)
	}

	return args.Error(1)
}

func (r *RepositoryMock) Delete(id string) error {
	args := r.Called(id)

	return args.Error(0)
}
