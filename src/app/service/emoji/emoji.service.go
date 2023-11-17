package emoji

import (
	"context"
	"time"

	"github.com/bookpanda/mygraderlist-backend/src/app/model"
	"github.com/bookpanda/mygraderlist-backend/src/app/model/emoji"
	proto "github.com/bookpanda/mygraderlist-proto/MyGraderList/backend/emoji"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	repository IRepository
}

type IRepository interface {
	FindAll(*[]*emoji.Emoji) error
	FindByUserId(string, *[]*emoji.Emoji) error
	Create(*emoji.Emoji) error
	Delete(string) error
}

func NewService(repository IRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) FindAll(_ context.Context, _ *proto.FindAllEmojiRequest) (*proto.FindAllEmojiResponse, error) {
	var emojis []*emoji.Emoji

	err := s.repository.FindAll(&emojis)
	if err != nil {

		log.Error().Err(err).
			Str("service", "emojis").
			Str("module", "find all").
			Msg("Error while querying all emojis")

		return nil, status.Error(codes.Unavailable, "Internal error")
	}

	return &proto.FindAllEmojiResponse{Emojis: RawToDtoList(&emojis)}, nil
}

func (s *Service) FindByUserId(_ context.Context, req *proto.FindByUserIdEmojiRequest) (res *proto.FindByUserIdEmojiResponse, err error) {
	var emojis []*emoji.Emoji

	err = s.repository.FindByUserId(req.UserId, &emojis)
	if err != nil {

		log.Error().Err(err).
			Str("service", "emoji").
			Str("module", "find by userId").
			Str("useId", req.UserId).
			Msg("Not found")

		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &proto.FindByUserIdEmojiResponse{Emojis: RawToDtoList(&emojis)}, nil
}

func (s *Service) Create(_ context.Context, req *proto.CreateEmojiRequest) (res *proto.CreateEmojiResponse, err error) {
	raw, _ := DtoToRaw(req.Emoji)

	err = s.repository.Create(raw)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create emoji")
	}

	return &proto.CreateEmojiResponse{Emoji: RawToDto(raw)}, nil
}

func (s *Service) Delete(_ context.Context, req *proto.DeleteEmojiRequest) (res *proto.DeleteEmojiResponse, err error) {
	err = s.repository.Delete(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "something wrong when deleting emoji")
	}

	return &proto.DeleteEmojiResponse{Success: true}, nil
}

func DtoToRaw(in *proto.Emoji) (result *emoji.Emoji, err error) {
	var id uuid.UUID
	if in.Id != "" {
		id, err = uuid.Parse(in.Id)
		if err != nil {
			return nil, err
		}
	}

	problemId, err := uuid.Parse(in.ProblemId)
	if err != nil {
		return nil, err
	}
	userId, err := uuid.Parse(in.UserId)
	if err != nil {
		return nil, err
	}

	return &emoji.Emoji{
		Base: model.Base{
			ID:        id,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			// DeletedAt: gorm.DeletedAt{},
		},
		Emoji:     in.Emoji,
		ProblemID: &problemId,
		UserID:    &userId,
	}, nil
}

func RawToDtoList(in *[]*emoji.Emoji) []*proto.Emoji {
	var result []*proto.Emoji
	for _, b := range *in {
		result = append(result, RawToDto(b))
	}

	return result
}

func RawToDto(in *emoji.Emoji) *proto.Emoji {
	return &proto.Emoji{
		Id:        in.ID.String(),
		Emoji:     in.Emoji,
		ProblemId: in.ProblemID.String(),
		UserId:    in.UserID.String(),
	}
}
