package rating

import (
	"context"
	"testing"
	"time"

	"github.com/bookpanda/mygraderlist-backend/src/app/model"
	"github.com/bookpanda/mygraderlist-backend/src/app/model/rating"
	mock "github.com/bookpanda/mygraderlist-backend/src/mocks/rating"
	proto "github.com/bookpanda/mygraderlist-proto/MyGraderList/backend/rating"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type RatingServiceTest struct {
	suite.Suite
	problemID           uuid.UUID
	userID              uuid.UUID
	Ratings             []*rating.Rating
	Rating              *rating.Rating
	UpdateRating        *rating.Rating
	RatingDto           *proto.Rating
	CreateRatingReqMock *proto.CreateRatingRequest
	UpdateRatingReqMock *proto.UpdateRatingRequest
}

func TestRatingService(t *testing.T) {
	suite.Run(t, new(RatingServiceTest))
}

func (t *RatingServiceTest) SetupTest() {
	t.Ratings = make([]*rating.Rating, 0)
	t.problemID = uuid.New()
	t.userID = uuid.New()

	t.Rating = &rating.Rating{
		Base: model.Base{
			ID:        uuid.New(),
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		},
		Score:      1,
		Difficulty: 2,
		ProblemID:  &t.problemID,
		UserID:     &t.userID,
	}

	Rating2 := &rating.Rating{
		Base: model.Base{
			ID:        uuid.New(),
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		},
		Score:      1,
		Difficulty: 2,
		ProblemID:  &t.problemID,
		UserID:     &t.userID,
	}
	Rating3 := &rating.Rating{
		Base: model.Base{
			ID:        uuid.New(),
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		},
		Score:      1,
		Difficulty: 2,
		ProblemID:  &t.problemID,
		UserID:     &t.userID,
	}
	t.Ratings = append(t.Ratings, t.Rating, Rating2, Rating3)

	t.RatingDto = &proto.Rating{
		Id:         t.Rating.ID.String(),
		Score:      int32(t.Rating.Score),
		Difficulty: int32(t.Rating.Difficulty),
		ProblemId:  t.Rating.ProblemID.String(),
		UserId:     t.Rating.UserID.String(),
	}

	t.CreateRatingReqMock = &proto.CreateRatingRequest{
		Rating: &proto.Rating{
			Score:      int32(t.Rating.Score),
			Difficulty: int32(t.Rating.Difficulty),
			ProblemId:  t.Rating.ProblemID.String(),
			UserId:     t.Rating.UserID.String(),
		},
	}

	t.UpdateRatingReqMock = &proto.UpdateRatingRequest{
		Id:         t.Rating.ID.String(),
		Score:      int32(t.Rating.Score),
		Difficulty: int32(t.Rating.Difficulty),
	}

	t.UpdateRating = &rating.Rating{
		Score:      1,
		Difficulty: 2,
	}
}

func (t *RatingServiceTest) TestFindAllSuccess() {
	want := &proto.FindAllRatingResponse{Ratings: createRatingDto(t.Ratings)}

	var RatingsIn []*rating.Rating

	repo := mock.RepositoryMock{}
	repo.On("FindAll", &RatingsIn).Return(&t.Ratings, nil)

	srv := NewService(&repo)
	actual, err := srv.FindAll(context.Background(), &proto.FindAllRatingRequest{})

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func createRatingDto(in []*rating.Rating) []*proto.Rating {
	var result []*proto.Rating

	for _, b := range in {
		r := &proto.Rating{
			Id:         b.ID.String(),
			Score:      int32(b.Score),
			Difficulty: int32(b.Difficulty),
			ProblemId:  b.ProblemID.String(),
			UserId:     b.UserID.String(),
		}

		result = append(result, r)
	}

	return result
}

func (t *RatingServiceTest) TestFindByUserIdSuccess() {
	want := &proto.FindByUserIdRatingResponse{Ratings: createRatingDto(t.Ratings)}

	repo := &mock.RepositoryMock{}

	var ratings []*rating.Rating
	repo.On("FindByUserId", t.Rating.UserID.String(), &ratings).Return(&t.Ratings, nil)

	srv := NewService(repo)
	actual, err := srv.FindByUserId(context.Background(), &proto.FindByUserIdRatingRequest{UserId: t.RatingDto.UserId})

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *RatingServiceTest) TestFindByUserIdNotFound() {
	repo := &mock.RepositoryMock{}

	var ratings []*rating.Rating
	repo.On("FindByUserId", t.Rating.UserID.String(), &ratings).Return(nil, errors.New("Not found user"))

	srv := NewService(repo)
	actual, err := srv.FindByUserId(context.Background(), &proto.FindByUserIdRatingRequest{UserId: t.RatingDto.UserId})

	st, ok := status.FromError(err)
	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.NotFound, st.Code())
}

func (t *RatingServiceTest) TestCreateSuccess() {
	want := &proto.CreateRatingResponse{Rating: t.RatingDto}

	in := &rating.Rating{
		Score:      1,
		Difficulty: 2,
		ProblemID:  &t.problemID,
		UserID:     &t.userID,
	}

	repo := &mock.RepositoryMock{}
	repo.On("Create", in).Return(t.Rating, nil)

	srv := NewService(repo)
	actual, err := srv.Create(context.Background(), t.CreateRatingReqMock)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *RatingServiceTest) TestCreateInternalErr() {
	in := &rating.Rating{
		Score:      t.Rating.Score,
		Difficulty: t.Rating.Difficulty,
		ProblemID:  t.Rating.ProblemID,
		UserID:     t.Rating.UserID,
	}

	repo := &mock.RepositoryMock{}
	repo.On("Create", in).Return(nil, errors.New("something wrong"))

	srv := NewService(repo)
	actual, err := srv.Create(context.Background(), t.CreateRatingReqMock)

	st, ok := status.FromError(err)

	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.Internal, st.Code())
}

func (t *RatingServiceTest) TestUpdateSuccess() {
	want := &proto.UpdateRatingResponse{Rating: t.RatingDto}

	repo := &mock.RepositoryMock{}
	repo.On("Update", t.Rating.ID.String(), t.UpdateRating).Return(t.Rating, nil)

	srv := NewService(repo)
	actual, err := srv.Update(context.Background(), t.UpdateRatingReqMock)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *RatingServiceTest) TestUpdateNotFound() {
	repo := &mock.RepositoryMock{}
	repo.On("Update", t.Rating.ID.String(), t.UpdateRating).Return(nil, errors.New("Not found Rating"))

	srv := NewService(repo)
	actual, err := srv.Update(context.Background(), t.UpdateRatingReqMock)

	st, ok := status.FromError(err)

	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.NotFound, st.Code())
}

func (t *RatingServiceTest) TestDeleteSuccess() {
	want := &proto.DeleteRatingResponse{Success: true}

	repo := &mock.RepositoryMock{}
	repo.On("Delete", t.Rating.ID.String()).Return(nil)

	srv := NewService(repo)
	actual, err := srv.Delete(context.Background(), &proto.DeleteRatingRequest{Id: t.RatingDto.Id})

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *RatingServiceTest) TestDeleteNotFound() {
	repo := &mock.RepositoryMock{}
	repo.On("Delete", t.Rating.ID.String()).Return(errors.New("Not found Rating"))

	srv := NewService(repo)
	actual, err := srv.Delete(context.Background(), &proto.DeleteRatingRequest{Id: t.RatingDto.Id})

	st, ok := status.FromError(err)
	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.NotFound, st.Code())
}
