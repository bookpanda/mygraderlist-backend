package course

import (
	"context"
	"time"

	"github.com/bookpanda/mygraderlist-backend/src/app/model"
	"github.com/bookpanda/mygraderlist-backend/src/app/model/course"
	proto "github.com/bookpanda/mygraderlist-proto/MyGraderList/backend/course"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	repository IRepository
}

type IRepository interface {
	FindAll(*[]*course.Course) error
	Create(*course.Course) error
	Update(string, *course.Course) error
	Delete(string) error
}

func NewService(repository IRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) FindAll(_ context.Context, _ *proto.FindAllCourseRequest) (*proto.FindAllCourseResponse, error) {
	var courses []*course.Course

	err := s.repository.FindAll(&courses)
	if err != nil {

		log.Error().Err(err).
			Str("service", "course").
			Str("module", "find all").
			Msg("Error while querying all courses")

		return nil, status.Error(codes.Unavailable, "Internal error")
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
		CourseCode: req.CourseCode,
		Name:       req.Name,
		Color:      req.Color,
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
			// DeletedAt: gorm.DeletedAt{},
		},
		CourseCode: in.CourseCode,
		Name:       in.Name,
		Color:      in.Color,
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
		Id:         in.ID.String(),
		CourseCode: in.CourseCode,
		Name:       in.Name,
		Color:      in.Color,
	}
}
