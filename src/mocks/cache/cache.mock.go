package cache

import (
	dto "github.com/bookpanda/mygraderlist-backend/src/app/model/problem"
	proto "github.com/bookpanda/mygraderlist-proto/MyGraderList/backend/problem"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
	V map[string]interface{}
}

func (t *RepositoryMock) SaveCache(key string, v interface{}, ttl int) error {
	args := t.Called(key, v, ttl)

	t.V[key] = v

	return args.Error(0)
}

func (t *RepositoryMock) GetCache(key string, v interface{}) error {
	args := t.Called(key, v)

	if args.Get(0) != nil {
		switch args.Get(0).(type) {
		case *[]*dto.Problem:
			*v.(*[]*dto.Problem) = *args.Get(0).(*[]*dto.Problem)
		case *dto.Problem:
			*v.(*dto.Problem) = *args.Get(0).(*dto.Problem)
		case *[]*proto.Problem:
			*v.(*[]*proto.Problem) = *args.Get(0).(*[]*proto.Problem)
		}
	}

	return args.Error(1)
}

func (t *RepositoryMock) RemoveCache(key string) (err error) {
	delete(t.V, key)
	return err
}
