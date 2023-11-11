package like

import (
	"github.com/bookpanda/mygraderlist-backend/src/app/model/like"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) FindByUserId(userId string, result *[]*like.Like) error {
	args := r.Called(userId, result)

	if args.Get(0) != nil {
		*result = *args.Get(0).(*[]*like.Like)
	}

	return args.Error(1)
}

func (r *RepositoryMock) Create(in *like.Like) error {
	args := r.Called(in)

	if args.Get(0) != nil {
		*in = *args.Get(0).(*like.Like)
	}

	return args.Error(1)
}

func (r *RepositoryMock) Delete(id string) error {
	args := r.Called(id)

	return args.Error(0)
}
