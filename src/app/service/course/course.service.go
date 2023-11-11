package course

import (
	"context"
	"time"

	"github.com/bookpanda/mygraderlist-backend/src/app/model"
	"github.com/bookpanda/mygraderlist-backend/src/app/model/course"
	"github.com/bookpanda/mygraderlist-backend/src/config"
	constant "github.com/bookpanda/mygraderlist-backend/src/constant/problem"
	proto "github.com/bookpanda/mygraderlist-proto/MyGraderList/backend/course"
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
	FindAll(*[]*course.Course) error
	Create(*course.Course) error
	Update(string, *course.Course) error
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

func (s *Service) FindAll(_ context.Context, _ *proto.FindAllCourseRequest) (*proto.FindAllCourseResponse, error) {
	var courses []*course.Course
	err := s.cache.GetCache(constant.ProblemKey, &courses)
	if err != redis.Nil {

		if err != nil {
			log.Error().
				Err(err).
				Str("service", "course").
				Str("module", "find all").
				Msg("Error while get cache")

			return nil, status.Error(codes.Unavailable, "Service is down")
		}

		return &proto.FindAllCourseResponse{Courses: RawToDtoList(&courses)}, nil
	}

	err = s.repository.FindAll(&courses)
	if err != nil {

		log.Error().Err(err).
			Str("service", "course").
			Str("module", "find all").
			Msg("Error while querying all baans")

		return nil, status.Error(codes.Unavailable, "Internal error")
	}

	err = s.cache.SaveCache(constant.ProblemKey, &courses, s.conf.ProblemCacheTTL)
	if err != nil {

		log.Error().
			Err(err).
			Str("service", "problem").
			Str("module", "find all").
			Msg("Error while saving the cache")

		return nil, status.Error(codes.Unavailable, "Service is down")
	}

	return &proto.FindAllCourseResponse{Courses: RawToDtoList(&courses)}, nil
}

func (s *Service) Create(_ context.Context, req *proto.CreateCourseRequest) (res *proto.CreateCourseResponse, err error) {
	raw, _ := DtoToRaw(req.Course)

	err = s.repository.Create(raw)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create course")
	}

	return &proto.CreateCourseResponse{Course: RawToDto(raw)}, nil
}

func (s *Service) Update(_ context.Context, req *proto.UpdateCourseRequest) (res *proto.UpdateCourseResponse, err error) {

	raw := &course.Course{
		Course: req.Course,
		Name:   req.Name,
		Color:  req.Color,
	}

	err = s.repository.Update(req.Id, raw)
	if err != nil {
		return nil, status.Error(codes.NotFound, "course not found")
	}

	return &proto.UpdateCourseResponse{Course: RawToDto(raw)}, nil
}

func (s *Service) Delete(_ context.Context, req *proto.DeleteCourseRequest) (res *proto.DeleteCourseResponse, err error) {
	err = s.repository.Delete(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "something wrong when deleting course")
	}

	return &proto.DeleteCourseResponse{Success: true}, nil
}

func DtoToRaw(in *proto.Course) (result *course.Course, err error) {
	var id uuid.UUID
	if in.Id != "" {
		id, err = uuid.Parse(in.Id)
		if err != nil {
			return nil, err
		}
	}

	return &course.Course{
		Base: model.Base{
			ID:        id,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		},
		Course: in.Course,
		Name:   in.Name,
		Color:  in.Color,
	}, nil
}

func RawToDtoList(in *[]*course.Course) []*proto.Course {
	var result []*proto.Course
	for _, b := range *in {
		result = append(result, RawToDto(b))
	}

	return result
}

func RawToDto(in *course.Course) *proto.Course {
	return &proto.Course{
		Id:     in.ID.String(),
		Course: in.Course,
		Name:   in.Name,
		Color:  in.Color,
	}
}
