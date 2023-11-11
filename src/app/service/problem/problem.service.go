package problem

import (
	"context"
	"time"

	"github.com/bookpanda/mygraderlist-backend/src/app/model"
	"github.com/bookpanda/mygraderlist-backend/src/app/model/problem"
	"github.com/bookpanda/mygraderlist-backend/src/config"
	constant "github.com/bookpanda/mygraderlist-backend/src/constant/problem"
	proto "github.com/bookpanda/mygraderlist-proto/MyGraderList/backend/problem"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type Service struct {
	repository IRepository
	cache      ICacheRepository
	conf       config.App
}

type IRepository interface {
	FindAll(*[]*problem.Problem) error
	Create(*problem.Problem) error
	Update(string, *problem.Problem) error
	Delete(string) error
}

type ICacheRepository interface {
	SaveCache(string, interface{}, int) error
	GetCache(string, interface{}) error
}

func NewService(repository IRepository, cache ICacheRepository, conf config.App) *Service {
	return &Service{
		repository: repository,
		cache:      cache,
		conf:       conf,
	}
}

func (s *Service) FindAll(_ context.Context, _ *proto.FindAllProblemRequest) (*proto.FindAllProblemResponse, error) {
	var problems []*problem.Problem
	err := s.cache.GetCache(constant.ProblemKey, &problems)
	if err != redis.Nil {

		if err != nil {
			log.Error().
				Err(err).
				Str("service", "problem").
				Str("module", "find all").
				Msg("Error while get cache")

			return nil, status.Error(codes.Unavailable, "Service is down")
		}

		return &proto.FindAllProblemResponse{Problems: RawToDtoList(&problems)}, nil
	}

	err = s.repository.FindAll(&problems)
	if err != nil {

		log.Error().Err(err).
			Str("service", "problem").
			Str("module", "find all").
			Msg("Error while querying all problems")

		return nil, status.Error(codes.Unavailable, "Internal error")
	}

	err = s.cache.SaveCache(constant.ProblemKey, &problems, s.conf.ProblemCacheTTL)
	if err != nil {

		log.Error().
			Err(err).
			Str("service", "problem").
			Str("module", "find all").
			Msg("Error while saving the cache")

		return nil, status.Error(codes.Unavailable, "Service is down")
	}

	return &proto.FindAllProblemResponse{Problems: RawToDtoList(&problems)}, nil
}

func (s *Service) Create(_ context.Context, req *proto.CreateProblemRequest) (res *proto.CreateProblemResponse, err error) {
	raw, _ := DtoToRaw(req.Problem)

	err = s.repository.Create(raw)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create problem")
	}

	return &proto.CreateProblemResponse{Problem: RawToDto(raw)}, nil
}

func (s *Service) Update(_ context.Context, req *proto.UpdateProblemRequest) (res *proto.UpdateProblemResponse, err error) {
	raw := &problem.Problem{
		CourseCode: req.CourseCode,
		Group:      req.Group,
		Code:       req.Code,
		Name:       req.Name,
	}

	err = s.repository.Update(req.Id, raw)
	if err != nil {
		return nil, status.Error(codes.NotFound, "problem not found")
	}

	return &proto.UpdateProblemResponse{Problem: RawToDto(raw)}, nil
}

func (s *Service) Delete(_ context.Context, req *proto.DeleteProblemRequest) (res *proto.DeleteProblemResponse, err error) {
	err = s.repository.Delete(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "something wrong when deleting problem")
	}

	return &proto.DeleteProblemResponse{Success: true}, nil
}

func DtoToRaw(in *proto.Problem) (result *problem.Problem, err error) {
	var id uuid.UUID
	if in.Id != "" {
		id, err = uuid.Parse(in.Id)
		if err != nil {
			return nil, err
		}
	}

	return &problem.Problem{
		Base: model.Base{
			ID:        id,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		},
		CourseCode: in.CourseCode,
		Group:      in.Group,
		Code:       in.Code,
		Name:       in.Name,
	}, nil
}

func RawToDtoList(in *[]*problem.Problem) []*proto.Problem {
	var result []*proto.Problem
	for _, b := range *in {
		result = append(result, RawToDto(b))
	}

	return result
}

func RawToDto(in *problem.Problem) *proto.Problem {
	return &proto.Problem{
		Id:         in.ID.String(),
		CourseCode: in.CourseCode,
		Group:      in.Group,
		Code:       in.Code,
		Name:       in.Name,
	}
}
