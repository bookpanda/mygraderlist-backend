package emoji

import (
	"context"
	"testing"
	"time"

	"github.com/bookpanda/mygraderlist-backend/src/app/model"
	"github.com/bookpanda/mygraderlist-backend/src/app/model/emoji"
	mock "github.com/bookpanda/mygraderlist-backend/src/mocks/emoji"
	proto "github.com/bookpanda/mygraderlist-proto/MyGraderList/backend/emoji"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EmojiServiceTest struct {
	suite.Suite
	problemID          uuid.UUID
	userID             uuid.UUID
	Emojis             []*emoji.Emoji
	Emoji              *emoji.Emoji
	EmojiDto           *proto.Emoji
	CreateEmojiReqMock *proto.CreateEmojiRequest
}

func TestEmojiService(t *testing.T) {
	suite.Run(t, new(EmojiServiceTest))
}

func (t *EmojiServiceTest) SetupTest() {
	t.Emojis = make([]*emoji.Emoji, 0)
	t.problemID = uuid.New()
	t.userID = uuid.New()

	t.Emoji = &emoji.Emoji{
		Base: model.Base{
			ID:        uuid.New(),
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			// DeletedAt: gorm.DeletedAt{},
		},
		Emoji:     "👍",
		ProblemID: &t.problemID,
		UserID:    &t.userID,
	}

	Emoji2 := &emoji.Emoji{
		Base: model.Base{
			ID:        uuid.New(),
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			// DeletedAt: gorm.DeletedAt{},
		},
		Emoji:     "💀",
		ProblemID: &t.problemID,
		UserID:    &t.userID,
	}
	Emoji3 := &emoji.Emoji{
		Base: model.Base{
			ID:        uuid.New(),
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			// DeletedAt: gorm.DeletedAt{},
		},
		Emoji:     "💩",
		ProblemID: &t.problemID,
		UserID:    &t.userID,
	}
	t.Emojis = append(t.Emojis, t.Emoji, Emoji2, Emoji3)

	t.EmojiDto = &proto.Emoji{
		Id:        t.Emoji.ID.String(),
		Emoji:     t.Emoji.Emoji,
		ProblemId: t.Emoji.ProblemID.String(),
		UserId:    t.Emoji.UserID.String(),
	}

	t.CreateEmojiReqMock = &proto.CreateEmojiRequest{
		Emoji: &proto.Emoji{
			Emoji:     t.Emoji.Emoji,
			ProblemId: t.Emoji.ProblemID.String(),
			UserId:    t.Emoji.UserID.String(),
		},
	}
}

func (t *EmojiServiceTest) TestFindAllSuccess() {
	want := &proto.FindAllEmojiResponse{Emojis: createEmojiDto(t.Emojis)}

	var EmojisIn []*emoji.Emoji

	repo := mock.RepositoryMock{}
	repo.On("FindAll", &EmojisIn).Return(&t.Emojis, nil)

	srv := NewService(&repo)
	actual, err := srv.FindAll(context.Background(), &proto.FindAllEmojiRequest{})

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func createEmojiDto(in []*emoji.Emoji) []*proto.Emoji {
	var result []*proto.Emoji

	for _, b := range in {
		r := &proto.Emoji{
			Id:        b.ID.String(),
			Emoji:     b.Emoji,
			ProblemId: b.ProblemID.String(),
			UserId:    b.UserID.String(),
		}

		result = append(result, r)
	}

	return result
}

func (t *EmojiServiceTest) TestFindByUserIdSuccess() {
	want := &proto.FindByUserIdEmojiResponse{Emojis: createEmojiDto(t.Emojis)}

	repo := &mock.RepositoryMock{}

	var Emojis []*emoji.Emoji
	repo.On("FindByUserId", t.Emoji.UserID.String(), &Emojis).Return(&t.Emojis, nil)

	srv := NewService(repo)
	actual, err := srv.FindByUserId(context.Background(), &proto.FindByUserIdEmojiRequest{UserId: t.EmojiDto.UserId})

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *EmojiServiceTest) TestFindByUserIdNotFound() {
	repo := &mock.RepositoryMock{}

	var Emojis []*emoji.Emoji
	repo.On("FindByUserId", t.Emoji.UserID.String(), &Emojis).Return(nil, errors.New("Not found user"))

	srv := NewService(repo)
	actual, err := srv.FindByUserId(context.Background(), &proto.FindByUserIdEmojiRequest{UserId: t.EmojiDto.UserId})

	st, ok := status.FromError(err)
	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.NotFound, st.Code())
}

func (t *EmojiServiceTest) TestCreateSuccess() {
	want := &proto.CreateEmojiResponse{Emoji: t.EmojiDto}

	in := &emoji.Emoji{
		Emoji:     t.Emoji.Emoji,
		ProblemID: &t.problemID,
		UserID:    &t.userID,
	}

	repo := &mock.RepositoryMock{}
	repo.On("Create", in).Return(t.Emoji, nil)

	srv := NewService(repo)
	actual, err := srv.Create(context.Background(), t.CreateEmojiReqMock)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *EmojiServiceTest) TestCreateInternalErr() {
	in := &emoji.Emoji{
		Emoji:     t.Emoji.Emoji,
		ProblemID: t.Emoji.ProblemID,
		UserID:    t.Emoji.UserID,
	}

	repo := &mock.RepositoryMock{}
	repo.On("Create", in).Return(nil, errors.New("something wrong"))

	srv := NewService(repo)
	actual, err := srv.Create(context.Background(), t.CreateEmojiReqMock)

	st, ok := status.FromError(err)

	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.Internal, st.Code())
}

func (t *EmojiServiceTest) TestDeleteSuccess() {
	want := &proto.DeleteEmojiResponse{Success: true}

	repo := &mock.RepositoryMock{}
	repo.On("Delete", t.Emoji.ID.String()).Return(nil)

	srv := NewService(repo)
	actual, err := srv.Delete(context.Background(), &proto.DeleteEmojiRequest{Id: t.EmojiDto.Id})

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *EmojiServiceTest) TestDeleteNotFound() {
	repo := &mock.RepositoryMock{}
	repo.On("Delete", t.Emoji.ID.String()).Return(errors.New("Not found Emoji"))

	srv := NewService(repo)
	actual, err := srv.Delete(context.Background(), &proto.DeleteEmojiRequest{Id: t.EmojiDto.Id})

	st, ok := status.FromError(err)
	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.NotFound, st.Code())
}
