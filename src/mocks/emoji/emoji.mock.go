package emoji

import (
	"github.com/bookpanda/mygraderlist-backend/src/app/model/emoji"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) FindAll(in *[]*emoji.Emoji) error {
	args := r.Called(in)

	if args.Get(0) != nil {
		*in = *args.Get(0).(*[]*emoji.Emoji)
	}

	return args.Error(1)
}

func (r *RepositoryMock) FindByUserId(userId string, result *[]*emoji.Emoji) error {
	args := r.Called(userId, result)

	if args.Get(0) != nil {
		*result = *args.Get(0).(*[]*emoji.Emoji)
	}

	return args.Error(1)
}

func (r *RepositoryMock) Create(in *emoji.Emoji) error {
	args := r.Called(in)

	if args.Get(0) != nil {
		*in = *args.Get(0).(*emoji.Emoji)
	}

	return args.Error(1)
}

func (r *RepositoryMock) Delete(id string) error {
	args := r.Called(id)

	return args.Error(0)
}
