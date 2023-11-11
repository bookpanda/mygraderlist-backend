package user

import (
	"github.com/bookpanda/mygraderlist-backend/src/app/model/user"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) FindOne(id string, result *user.User) error {
	args := r.Called(id, result)

	if args.Get(0) != nil {
		*result = *args.Get(0).(*user.User)
	}

	return args.Error(1)
}

func (r *RepositoryMock) FindByEmail(email string, result *user.User) error {
	args := r.Called(email, result)

	if args.Get(0) != nil {
		*result = *args.Get(0).(*user.User)
	}

	return args.Error(1)
}

func (r *RepositoryMock) Create(in *user.User) error {
	args := r.Called(in)

	if args.Get(0) != nil {
		*in = *args.Get(0).(*user.User)
	}

	return args.Error(1)
}

func (r *RepositoryMock) Update(id string, result *user.User) error {
	args := r.Called(id, result)

	if args.Get(0) != nil {
		*result = *args.Get(0).(*user.User)
	}

	return args.Error(1)
}

func (r *RepositoryMock) Delete(id string) error {
	args := r.Called(id)

	return args.Error(0)
}
