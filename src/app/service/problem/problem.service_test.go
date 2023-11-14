package problem

import (
	"context"
	"testing"
	"time"

	"github.com/bookpanda/mygraderlist-backend/src/app/model"
	"github.com/bookpanda/mygraderlist-backend/src/app/model/problem"
	"github.com/bookpanda/mygraderlist-backend/src/config"
	constant "github.com/bookpanda/mygraderlist-backend/src/constant/problem"
	cacheMock "github.com/bookpanda/mygraderlist-backend/src/mocks/cache"
	mock "github.com/bookpanda/mygraderlist-backend/src/mocks/problem"
	proto "github.com/bookpanda/mygraderlist-proto/MyGraderList/backend/problem"
	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type ProblemServiceTest struct {
	suite.Suite
	Problems             []*problem.Problem
	Problem              *problem.Problem
	UpdateProblem        *problem.Problem
	ProblemDto           *proto.Problem
	CreateProblemReqMock *proto.CreateProblemRequest
	UpdateProblemReqMock *proto.UpdateProblemRequest
	conf                 config.App
}

func TestProblemService(t *testing.T) {
	suite.Run(t, new(ProblemServiceTest))
}

func (t *ProblemServiceTest) SetupTest() {
	t.Problems = make([]*problem.Problem, 0)

	t.Problem = &problem.Problem{
		Base: model.Base{
			ID:        uuid.New(),
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		},
		Group:      faker.Name(),
		Code:       faker.Name(),
		Name:       faker.Name(),
		CourseCode: faker.Name(),
		Order:      1,
	}

	Problem2 := &problem.Problem{
		Base: model.Base{
			ID:        uuid.New(),
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		},
		Group:      faker.Name(),
		Code:       faker.Name(),
		Name:       faker.Name(),
		CourseCode: faker.Name(),
		Order:      2,
	}
	Problem3 := &problem.Problem{
		Base: model.Base{
			ID:        uuid.New(),
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		},
		Group:      faker.Name(),
		Code:       faker.Name(),
		Name:       faker.Name(),
		CourseCode: faker.Name(),
		Order:      3,
	}
	t.Problems = append(t.Problems, t.Problem, Problem2, Problem3)

	t.ProblemDto = &proto.Problem{
		Id:         t.Problem.ID.String(),
		Group:      t.Problem.Group,
		Code:       t.Problem.Code,
		Name:       t.Problem.Name,
		CourseCode: t.Problem.CourseCode,
		Order:      int32(t.Problem.Order),
	}

	t.CreateProblemReqMock = &proto.CreateProblemRequest{
		Problem: &proto.Problem{
			Group:      t.Problem.Group,
			Code:       t.Problem.Code,
			Name:       t.Problem.Name,
			CourseCode: t.Problem.CourseCode,
			Order:      int32(t.Problem.Order),
		},
	}

	t.UpdateProblemReqMock = &proto.UpdateProblemRequest{
		Id:         t.Problem.ID.String(),
		Group:      t.Problem.Group,
		Code:       t.Problem.Code,
		Name:       t.Problem.Name,
		CourseCode: t.Problem.CourseCode,
		Order:      int32(t.Problem.Order),
	}

	t.UpdateProblem = &problem.Problem{
		Group:      t.Problem.Group,
		Code:       t.Problem.Code,
		Name:       t.Problem.Name,
		CourseCode: t.Problem.CourseCode,
		Order:      t.Problem.Order,
	}

	t.conf = config.App{
		Port:            3001,
		Debug:           false,
		ProblemCacheTTL: 90,
	}
}

func (t *ProblemServiceTest) TestFindAllSuccess() {
	want := &proto.FindAllProblemResponse{Problems: createProblemDto(t.Problems)}

	var problemsIn []*problem.Problem

	repo := &mock.RepositoryMock{}
	repo.On("FindAll", &problemsIn).Return(&t.Problems, nil)

	c := &cacheMock.RepositoryMock{
		V: map[string]interface{}{},
	}
	c.On("SaveCache", constant.ProblemKey, &t.Problems, t.conf.ProblemCacheTTL).Return(nil)
	c.On("GetCache", constant.ProblemKey, &problemsIn).Return(&t.Problems, nil)

	srv := NewService(repo, c, t.conf)
	actual, err := srv.FindAll(context.Background(), &proto.FindAllProblemRequest{})

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func createProblemDto(in []*problem.Problem) []*proto.Problem {
	var result []*proto.Problem

	for _, b := range in {
		r := &proto.Problem{
			Id:         b.ID.String(),
			Group:      b.Group,
			Code:       b.Code,
			Name:       b.Name,
			CourseCode: b.CourseCode,
			Order:      int32(b.Order),
		}

		result = append(result, r)
	}

	return result
}

func (t *ProblemServiceTest) TestCreateSuccess() {
	want := &proto.CreateProblemResponse{Problem: t.ProblemDto}

	var problemsIn []*problem.Problem

	in := &problem.Problem{
		Group:      t.Problem.Group,
		Code:       t.Problem.Code,
		Name:       t.Problem.Name,
		CourseCode: t.Problem.CourseCode,
		Order:      t.Problem.Order,
	}

	repo := &mock.RepositoryMock{}
	repo.On("Create", in).Return(t.Problem, nil)

	c := &cacheMock.RepositoryMock{
		V: map[string]interface{}{},
	}
	c.On("SaveCache", constant.ProblemKey, &t.Problems, t.conf.ProblemCacheTTL).Return(nil)
	c.On("GetCache", constant.ProblemKey, &problemsIn).Return(&t.Problems, nil)

	srv := NewService(repo, c, t.conf)
	actual, err := srv.Create(context.Background(), t.CreateProblemReqMock)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *ProblemServiceTest) TestCreateInternalErr() {
	in := &problem.Problem{
		Group:      t.Problem.Group,
		Code:       t.Problem.Code,
		Name:       t.Problem.Name,
		CourseCode: t.Problem.CourseCode,
		Order:      t.Problem.Order,
	}
	var problemsIn []*problem.Problem

	repo := &mock.RepositoryMock{}
	repo.On("Create", in).Return(nil, errors.New("something wrong"))

	c := &cacheMock.RepositoryMock{
		V: map[string]interface{}{},
	}
	c.On("SaveCache", constant.ProblemKey, &t.Problems, t.conf.ProblemCacheTTL).Return(nil)
	c.On("GetCache", constant.ProblemKey, &problemsIn).Return(&t.Problems, nil)

	srv := NewService(repo, c, t.conf)
	actual, err := srv.Create(context.Background(), t.CreateProblemReqMock)

	st, ok := status.FromError(err)

	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.Internal, st.Code())
}

func (t *ProblemServiceTest) TestUpdateSuccess() {
	want := &proto.UpdateProblemResponse{Problem: t.ProblemDto}

	var problemsIn []*problem.Problem

	repo := &mock.RepositoryMock{}
	repo.On("Update", t.Problem.ID.String(), t.UpdateProblem).Return(t.Problem, nil)

	c := &cacheMock.RepositoryMock{
		V: map[string]interface{}{},
	}
	c.On("SaveCache", constant.ProblemKey, &t.Problems, t.conf.ProblemCacheTTL).Return(nil)
	c.On("GetCache", constant.ProblemKey, &problemsIn).Return(&t.Problems, nil)

	srv := NewService(repo, c, t.conf)
	actual, err := srv.Update(context.Background(), t.UpdateProblemReqMock)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *ProblemServiceTest) TestUpdateNotFound() {
	var problemsIn []*problem.Problem

	repo := &mock.RepositoryMock{}
	repo.On("Update", t.Problem.ID.String(), t.UpdateProblem).Return(nil, errors.New("Not found Problem"))

	c := &cacheMock.RepositoryMock{
		V: map[string]interface{}{},
	}
	c.On("SaveCache", constant.ProblemKey, &t.Problems, t.conf.ProblemCacheTTL).Return(nil)
	c.On("GetCache", constant.ProblemKey, &problemsIn).Return(&t.Problems, nil)

	srv := NewService(repo, c, t.conf)
	actual, err := srv.Update(context.Background(), t.UpdateProblemReqMock)

	st, ok := status.FromError(err)

	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.NotFound, st.Code())
}

func (t *ProblemServiceTest) TestDeleteSuccess() {
	want := &proto.DeleteProblemResponse{Success: true}

	var problemsIn []*problem.Problem

	repo := &mock.RepositoryMock{}
	repo.On("Delete", t.Problem.ID.String()).Return(nil)

	c := &cacheMock.RepositoryMock{
		V: map[string]interface{}{},
	}
	c.On("SaveCache", constant.ProblemKey, &t.Problems, t.conf.ProblemCacheTTL).Return(nil)
	c.On("GetCache", constant.ProblemKey, &problemsIn).Return(&t.Problems, nil)

	srv := NewService(repo, c, t.conf)
	actual, err := srv.Delete(context.Background(), &proto.DeleteProblemRequest{Id: t.ProblemDto.Id})

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *ProblemServiceTest) TestDeleteNotFound() {
	repo := &mock.RepositoryMock{}
	repo.On("Delete", t.Problem.ID.String()).Return(errors.New("Not found Problem"))

	var problemsIn []*problem.Problem

	c := &cacheMock.RepositoryMock{
		V: map[string]interface{}{},
	}
	c.On("SaveCache", constant.ProblemKey, &t.Problems, t.conf.ProblemCacheTTL).Return(nil)
	c.On("GetCache", constant.ProblemKey, &problemsIn).Return(&t.Problems, nil)

	srv := NewService(repo, c, t.conf)
	actual, err := srv.Delete(context.Background(), &proto.DeleteProblemRequest{Id: t.ProblemDto.Id})

	st, ok := status.FromError(err)
	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.NotFound, st.Code())
}
