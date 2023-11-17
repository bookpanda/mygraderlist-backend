package user

import (
	"context"
	"time"

	"github.com/bookpanda/mygraderlist-backend/src/app/model"
	"github.com/bookpanda/mygraderlist-backend/src/app/model/user"
	proto "github.com/bookpanda/mygraderlist-proto/MyGraderList/backend/user"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	repository IRepository
}

type IRepository interface {
	FindOne(string, *user.User) error
	FindByEmail(string, *user.User) error
	Create(*user.User) error
	Update(string, *user.User) error
	Delete(string) error
}

func NewService(repository IRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) FindOne(_ context.Context, req *proto.FindOneUserRequest) (res *proto.FindOneUserResponse, err error) {
	raw := user.User{}

	err = s.repository.FindOne(req.Id, &raw)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &proto.FindOneUserResponse{User: RawToDto(&raw)}, nil
}

func (s *Service) FindByEmail(_ context.Context, req *proto.FindByEmailUserRequest) (res *proto.FindByEmailUserResponse, err error) {
	raw := user.User{}

	err = s.repository.FindByEmail(req.Email, &raw)
	if err != nil {

		log.Error().Err(err).
			Str("service", "user").
			Str("module", "find one").
			Str("email", req.Email).
			Msg("Not found")

		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &proto.FindByEmailUserResponse{User: RawToDto(&raw)}, nil
}

func (s *Service) Create(_ context.Context, req *proto.CreateUserRequest) (res *proto.CreateUserResponse, err error) {
	raw, _ := DtoToRaw(req.User)

	err = s.repository.Create(raw)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	return &proto.CreateUserResponse{User: RawToDto(raw)}, nil
}

func (s *Service) Update(_ context.Context, req *proto.UpdateUserRequest) (res *proto.UpdateUserResponse, err error) {
	raw := &user.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	err = s.repository.Update(req.Id, raw)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &proto.UpdateUserResponse{User: RawToDto(raw)}, nil
}

func (s *Service) Delete(_ context.Context, req *proto.DeleteUserRequest) (res *proto.DeleteUserResponse, err error) {
	err = s.repository.Delete(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "something wrong when deleting user")
	}

	return &proto.DeleteUserResponse{Success: true}, nil
}

func DtoToRaw(in *proto.User) (result *user.User, err error) {
	var id uuid.UUID
	if in.Id != "" {
		id, err = uuid.Parse(in.Id)
		if err != nil {
			return nil, err
		}
	}

	return &user.User{
		Base: model.Base{
			ID:        id,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			// DeletedAt: gorm.DeletedAt{},
		},
		Username: in.Username,
		Email:    in.Email,
		Password: in.Password,
	}, nil
}

func RawToDtoList(in *[]*user.User) []*proto.User {
	var result []*proto.User
	for _, b := range *in {
		result = append(result, RawToDto(b))
	}

	return result
}

func RawToDto(in *user.User) *proto.User {
	return &proto.User{
		Id:       in.ID.String(),
		Username: in.Username,
		Email:    in.Email,
		Password: in.Password,
	}
}
