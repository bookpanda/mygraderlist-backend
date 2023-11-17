package user

import (
	"context"
	"testing"
	"time"

	"github.com/bookpanda/mygraderlist-backend/src/app/model"
	"github.com/bookpanda/mygraderlist-backend/src/app/model/user"
	mock "github.com/bookpanda/mygraderlist-backend/src/mocks/user"
	proto "github.com/bookpanda/mygraderlist-proto/MyGraderList/backend/user"
	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServiceTest struct {
	suite.Suite
	User              *user.User
	UpdateUser        *user.User
	UserDto           *proto.User
	CreateUserReqMock *proto.CreateUserRequest
	UpdateUserReqMock *proto.UpdateUserRequest
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(UserServiceTest))
}

func (t *UserServiceTest) SetupTest() {
	t.User = &user.User{
		Base: model.Base{
			ID:        uuid.New(),
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			// DeletedAt: gorm.DeletedAt{},
		},
		Username: faker.Username(),
		Email:    faker.Email(),
		Password: faker.Password(),
	}

	t.UserDto = &proto.User{
		Id:       t.User.ID.String(),
		Username: t.User.Username,
		Email:    t.User.Email,
		Password: t.User.Password,
	}

	t.CreateUserReqMock = &proto.CreateUserRequest{
		User: &proto.User{
			Username: t.User.Username,
			Email:    t.User.Email,
			Password: t.User.Password,
		},
	}

	t.UpdateUserReqMock = &proto.UpdateUserRequest{
		Id:       t.User.ID.String(),
		Username: t.User.Username,
		Email:    t.User.Email,
		Password: t.User.Password,
	}

	t.UpdateUser = &user.User{
		Username: t.User.Username,
		Email:    t.User.Email,
		Password: t.User.Password,
	}
}

func (t *UserServiceTest) TestFindOneSuccess() {
	want := &proto.FindOneUserResponse{User: t.UserDto}

	repo := &mock.RepositoryMock{}
	repo.On("FindOne", t.User.ID.String(), &user.User{}).Return(t.User, nil)

	srv := NewService(repo)
	actual, err := srv.FindOne(context.Background(), &proto.FindOneUserRequest{Id: t.User.ID.String()})

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *UserServiceTest) TestFindOneNotFound() {
	repo := &mock.RepositoryMock{}
	repo.On("FindOne", t.User.ID.String(), &user.User{}).Return(nil, errors.New("Not found user"))

	srv := NewService(repo)
	actual, err := srv.FindOne(context.Background(), &proto.FindOneUserRequest{Id: t.User.ID.String()})

	st, ok := status.FromError(err)

	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.NotFound, st.Code())
}

func (t *UserServiceTest) TestFindByEmailSuccess() {
	want := &proto.FindByEmailUserResponse{User: t.UserDto}

	repo := &mock.RepositoryMock{}

	repo.On("FindByEmail", t.User.Email, &user.User{}).Return(t.User, nil)

	srv := NewService(repo)
	actual, err := srv.FindByEmail(context.Background(), &proto.FindByEmailUserRequest{Email: t.UserDto.Email})

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *UserServiceTest) TestFindByEmailNotFound() {
	repo := &mock.RepositoryMock{}

	repo.On("FindByEmail", t.User.Email, &user.User{}).Return(nil, errors.New("Not found user"))

	srv := NewService(repo)
	actual, err := srv.FindByEmail(context.Background(), &proto.FindByEmailUserRequest{Email: t.UserDto.Email})

	st, ok := status.FromError(err)
	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.NotFound, st.Code())
}

func (t *UserServiceTest) TestCreateSuccess() {
	want := &proto.CreateUserResponse{User: t.UserDto}

	in := &user.User{
		Username: t.User.Username,
		Email:    t.User.Email,
		Password: t.User.Password,
	}

	repo := &mock.RepositoryMock{}
	repo.On("Create", in).Return(t.User, nil)

	srv := NewService(repo)
	actual, err := srv.Create(context.Background(), t.CreateUserReqMock)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *UserServiceTest) TestCreateInternalErr() {
	in := &user.User{
		Username: t.User.Username,
		Email:    t.User.Email,
		Password: t.User.Password,
	}

	repo := &mock.RepositoryMock{}
	repo.On("Create", in).Return(nil, errors.New("something wrong"))

	srv := NewService(repo)
	actual, err := srv.Create(context.Background(), t.CreateUserReqMock)

	st, ok := status.FromError(err)

	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.Internal, st.Code())
}

func (t *UserServiceTest) TestUpdateSuccess() {
	want := &proto.UpdateUserResponse{User: t.UserDto}

	repo := &mock.RepositoryMock{}
	repo.On("Update", t.User.ID.String(), t.UpdateUser).Return(t.User, nil)

	srv := NewService(repo)
	actual, err := srv.Update(context.Background(), t.UpdateUserReqMock)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *UserServiceTest) TestUpdateNotFound() {
	repo := &mock.RepositoryMock{}
	repo.On("Update", t.User.ID.String(), t.UpdateUser).Return(nil, errors.New("Not found user"))

	srv := NewService(repo)
	actual, err := srv.Update(context.Background(), t.UpdateUserReqMock)

	st, ok := status.FromError(err)

	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.NotFound, st.Code())
}

func (t *UserServiceTest) TestDeleteSuccess() {
	want := &proto.DeleteUserResponse{Success: true}

	repo := &mock.RepositoryMock{}
	repo.On("Delete", t.User.ID.String()).Return(nil)

	srv := NewService(repo)
	actual, err := srv.Delete(context.Background(), &proto.DeleteUserRequest{Id: t.UserDto.Id})

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *UserServiceTest) TestDeleteNotFound() {
	repo := &mock.RepositoryMock{}
	repo.On("Delete", t.User.ID.String()).Return(errors.New("Not found user"))

	srv := NewService(repo)
	actual, err := srv.Delete(context.Background(), &proto.DeleteUserRequest{Id: t.UserDto.Id})

	st, ok := status.FromError(err)
	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.NotFound, st.Code())
}
