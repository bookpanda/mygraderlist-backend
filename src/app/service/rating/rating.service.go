package rating

import (
	"context"
	"time"

	"github.com/bookpanda/mygraderlist-backend/src/app/model"
	"github.com/bookpanda/mygraderlist-backend/src/app/model/rating"
	"github.com/bookpanda/mygraderlist-backend/src/config"
	proto "github.com/bookpanda/mygraderlist-proto/MyGraderList/backend/rating"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type Service struct {
	repository IRepository
	conf       config.App
}

type IRepository interface {
	FindAll(*[]*rating.Rating) error
	FindByUserId(string, *[]*rating.Rating) error
	Create(*rating.Rating) error
	Update(string, *rating.Rating) error
	Delete(string) error
}

func NewService(repository IRepository, conf config.App) *Service {
	return &Service{
		repository: repository,
		conf:       conf,
	}
}

func (s *Service) FindAll(_ context.Context, _ *proto.FindAllRatingRequest) (*proto.FindAllRatingResponse, error) {
	var ratings []*rating.Rating

	err := s.repository.FindAll(&ratings)
	if err != nil {

		log.Error().Err(err).
			Str("service", "ratings").
			Str("module", "find all").
			Msg("Error while querying all ratings")

		return nil, status.Error(codes.Unavailable, "Internal error")
	}

	return &proto.FindAllRatingResponse{Ratings: RawToDtoList(&ratings)}, nil
}

func (s *Service) FindByUserId(_ context.Context, req *proto.FindByUserIdRatingRequest) (res *proto.FindByUserIdRatingResponse, err error) {
	var ratings []*rating.Rating

	err = s.repository.FindByUserId(req.UserId, &ratings)
	if err != nil {

		log.Error().Err(err).
			Str("service", "rating").
			Str("module", "find by userId").
			Str("useId", req.UserId).
			Msg("Not found")

		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &proto.FindByUserIdRatingResponse{Ratings: RawToDtoList(&ratings)}, nil
}

func (s *Service) Create(_ context.Context, req *proto.CreateRatingRequest) (res *proto.CreateRatingResponse, err error) {
	raw, _ := DtoToRaw(req.Rating)

	err = s.repository.Create(raw)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create rating")
	}

	return &proto.CreateRatingResponse{Rating: RawToDto(raw)}, nil
}

func (s *Service) Update(_ context.Context, req *proto.UpdateRatingRequest) (res *proto.UpdateRatingResponse, err error) {

	raw := &rating.Rating{
		Score:      int(req.Score),
		Difficulty: int(req.Difficulty),
	}

	err = s.repository.Update(req.Id, raw)
	if err != nil {
		return nil, status.Error(codes.NotFound, "rating not found")
	}

	return &proto.UpdateRatingResponse{Rating: RawToDto(raw)}, nil
}

func (s *Service) Delete(_ context.Context, req *proto.DeleteRatingRequest) (res *proto.DeleteRatingResponse, err error) {
	err = s.repository.Delete(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "something wrong when deleting rating")
	}

	return &proto.DeleteRatingResponse{Success: true}, nil
}

func DtoToRaw(in *proto.Rating) (result *rating.Rating, err error) {
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

	return &rating.Rating{
		Base: model.Base{
			ID:        id,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		},
		ProblemId:  &problemId,
		UserId:     &userId,
		Score:      int(in.Score),
		Difficulty: int(in.Difficulty),
	}, nil
}

func RawToDtoList(in *[]*rating.Rating) []*proto.Rating {
	var result []*proto.Rating
	for _, b := range *in {
		result = append(result, RawToDto(b))
	}

	return result
}

func RawToDto(in *rating.Rating) *proto.Rating {
	return &proto.Rating{
		Id:         in.ID.String(),
		ProblemId:  in.ProblemId.String(),
		UserId:     in.UserId.String(),
		Score:      int32(in.Score),
		Difficulty: int32(in.Difficulty),
	}
}
