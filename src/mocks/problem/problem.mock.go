package problem

import (
	"github.com/bookpanda/mygraderlist-backend/src/app/model/problem"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) FindAll(in *[]*problem.Problem) error {
	args := r.Called(in)

	if args.Get(0) != nil {
		*in = *args.Get(0).(*[]*problem.Problem)
	}

	return args.Error(1)
}

func (r *RepositoryMock) Create(in *problem.Problem) error {
	args := r.Called(in)

	if args.Get(0) != nil {
		*in = *args.Get(0).(*problem.Problem)
	}

	return args.Error(1)
}

func (r *RepositoryMock) Update(id string, result *problem.Problem) error {
	args := r.Called(id, result)

	if args.Get(0) != nil {
		*result = *args.Get(0).(*problem.Problem)
	}

	return args.Error(1)
}

func (r *RepositoryMock) Delete(id string) error {
	args := r.Called(id)

	return args.Error(0)
}
