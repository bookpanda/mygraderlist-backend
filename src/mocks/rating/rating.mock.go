package rating

import (
	"github.com/bookpanda/mygraderlist-backend/src/app/model/rating"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) FindAll(in *[]*rating.Rating) error {
	args := r.Called(in)

	if args.Get(0) != nil {
		*in = *args.Get(0).(*[]*rating.Rating)
	}

	return args.Error(1)
}

func (r *RepositoryMock) FindByUserId(userId string, result *[]*rating.Rating) error {
	args := r.Called(userId, result)

	if args.Get(0) != nil {
		*result = *args.Get(0).(*[]*rating.Rating)
	}

	return args.Error(1)
}

func (r *RepositoryMock) Create(in *rating.Rating) error {
	args := r.Called(in)

	if args.Get(0) != nil {
		*in = *args.Get(0).(*rating.Rating)
	}

	return args.Error(1)
}

func (r *RepositoryMock) Update(id string, result *rating.Rating) error {
	args := r.Called(id, result)

	if args.Get(0) != nil {
		*result = *args.Get(0).(*rating.Rating)
	}

	return args.Error(1)
}

func (r *RepositoryMock) Delete(id string) error {
	args := r.Called(id)

	return args.Error(0)
}
