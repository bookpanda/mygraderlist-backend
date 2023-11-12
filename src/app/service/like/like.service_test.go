package like

import (
	"context"
	"testing"
	"time"

	"github.com/bookpanda/mygraderlist-backend/src/app/model"
	"github.com/bookpanda/mygraderlist-backend/src/app/model/like"
	mock "github.com/bookpanda/mygraderlist-backend/src/mocks/like"
	proto "github.com/bookpanda/mygraderlist-proto/MyGraderList/backend/like"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type LikeServiceTest struct {
	suite.Suite
	problemID         uuid.UUID
	userID            uuid.UUID
	Likes             []*like.Like
	Like              *like.Like
	LikeDto           *proto.Like
	CreateLikeReqMock *proto.CreateLikeRequest
}

func TestLikeService(t *testing.T) {
	suite.Run(t, new(LikeServiceTest))
}

func (t *LikeServiceTest) SetupTest() {
	t.Likes = make([]*like.Like, 0)
	t.problemID = uuid.New()
	t.userID = uuid.New()

	t.Like = &like.Like{
		Base: model.Base{
			ID:        uuid.New(),
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		},
		ProblemID: &t.problemID,
		UserID:    &t.userID,
	}

	Like2 := &like.Like{
		Base: model.Base{
			ID:        uuid.New(),
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		},
		ProblemID: &t.problemID,
		UserID:    &t.userID,
	}
	Like3 := &like.Like{
		Base: model.Base{
			ID:        uuid.New(),
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		},
		ProblemID: &t.problemID,
		UserID:    &t.userID,
	}
	t.Likes = append(t.Likes, t.Like, Like2, Like3)

	t.LikeDto = &proto.Like{
		Id:        t.Like.ID.String(),
		ProblemId: t.Like.ProblemID.String(),
		UserId:    t.Like.UserID.String(),
	}

	t.CreateLikeReqMock = &proto.CreateLikeRequest{
		Like: &proto.Like{
			ProblemId: t.Like.ProblemID.String(),
			UserId:    t.Like.UserID.String(),
		},
	}
}

func createLikeDto(in []*like.Like) []*proto.Like {
	var result []*proto.Like

	for _, b := range in {
		r := &proto.Like{
			Id:        b.ID.String(),
			ProblemId: b.ProblemID.String(),
			UserId:    b.UserID.String(),
		}

		result = append(result, r)
	}

	return result
}

func (t *LikeServiceTest) TestFindByUserIdSuccess() {
	want := &proto.FindByUserIdLikeResponse{Likes: createLikeDto(t.Likes)}

	repo := &mock.RepositoryMock{}

	var Likes []*like.Like
	repo.On("FindByUserId", t.Like.UserID.String(), &Likes).Return(&t.Likes, nil)

	srv := NewService(repo)
	actual, err := srv.FindByUserId(context.Background(), &proto.FindByUserIdLikeRequest{UserId: t.LikeDto.UserId})

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *LikeServiceTest) TestFindByUserIdNotFound() {
	repo := &mock.RepositoryMock{}

	var Likes []*like.Like
	repo.On("FindByUserId", t.Like.UserID.String(), &Likes).Return(nil, errors.New("Not found user"))

	srv := NewService(repo)
	actual, err := srv.FindByUserId(context.Background(), &proto.FindByUserIdLikeRequest{UserId: t.LikeDto.UserId})

	st, ok := status.FromError(err)
	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.NotFound, st.Code())
}

func (t *LikeServiceTest) TestCreateSuccess() {
	want := &proto.CreateLikeResponse{Like: t.LikeDto}

	in := &like.Like{
		ProblemID: &t.problemID,
		UserID:    &t.userID,
	}

	repo := &mock.RepositoryMock{}
	repo.On("Create", in).Return(t.Like, nil)

	srv := NewService(repo)
	actual, err := srv.Create(context.Background(), t.CreateLikeReqMock)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *LikeServiceTest) TestCreateInternalErr() {
	in := &like.Like{
		ProblemID: t.Like.ProblemID,
		UserID:    t.Like.UserID,
	}

	repo := &mock.RepositoryMock{}
	repo.On("Create", in).Return(nil, errors.New("something wrong"))

	srv := NewService(repo)
	actual, err := srv.Create(context.Background(), t.CreateLikeReqMock)

	st, ok := status.FromError(err)

	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.Internal, st.Code())
}

func (t *LikeServiceTest) TestDeleteSuccess() {
	want := &proto.DeleteLikeResponse{Success: true}

	repo := &mock.RepositoryMock{}
	repo.On("Delete", t.Like.ID.String()).Return(nil)

	srv := NewService(repo)
	actual, err := srv.Delete(context.Background(), &proto.DeleteLikeRequest{Id: t.LikeDto.Id})

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *LikeServiceTest) TestDeleteNotFound() {
	repo := &mock.RepositoryMock{}
	repo.On("Delete", t.Like.ID.String()).Return(errors.New("Not found Like"))

	srv := NewService(repo)
	actual, err := srv.Delete(context.Background(), &proto.DeleteLikeRequest{Id: t.LikeDto.Id})

	st, ok := status.FromError(err)
	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.NotFound, st.Code())
}
