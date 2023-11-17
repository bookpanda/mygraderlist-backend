package course

import (
	"context"
	"testing"
	"time"

	"github.com/bookpanda/mygraderlist-backend/src/app/model"
	"github.com/bookpanda/mygraderlist-backend/src/app/model/course"
	mock "github.com/bookpanda/mygraderlist-backend/src/mocks/course"
	proto "github.com/bookpanda/mygraderlist-proto/MyGraderList/backend/course"
	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CourseServiceTest struct {
	suite.Suite
	Courses             []*course.Course
	Course              *course.Course
	UpdateCourse        *course.Course
	CourseDto           *proto.Course
	CreateCourseReqMock *proto.CreateCourseRequest
	UpdateCourseReqMock *proto.UpdateCourseRequest
}

func TestCourseService(t *testing.T) {
	suite.Run(t, new(CourseServiceTest))
}

func (t *CourseServiceTest) SetupTest() {
	t.Courses = make([]*course.Course, 0)

	t.Course = &course.Course{
		Base: model.Base{
			ID:        uuid.New(),
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			// DeletedAt: gorm.DeletedAt{},
		},
		CourseCode: faker.Name(),
		Name:       faker.Name(),
		Color:      faker.Word(),
	}

	course2 := &course.Course{
		Base: model.Base{
			ID:        uuid.New(),
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			// DeletedAt: gorm.DeletedAt{},
		},
		CourseCode: faker.Name(),
		Name:       faker.Name(),
		Color:      faker.Word(),
	}
	course3 := &course.Course{
		Base: model.Base{
			ID:        uuid.New(),
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			// DeletedAt: gorm.DeletedAt{},
		},
		CourseCode: faker.Name(),
		Name:       faker.Name(),
		Color:      faker.Word(),
	}
	t.Courses = append(t.Courses, t.Course, course2, course3)

	t.CourseDto = &proto.Course{
		Id:         t.Course.ID.String(),
		CourseCode: t.Course.CourseCode,
		Name:       t.Course.Name,
		Color:      t.Course.Color,
	}

	t.CreateCourseReqMock = &proto.CreateCourseRequest{
		Course: &proto.Course{
			CourseCode: t.Course.CourseCode,
			Name:       t.Course.Name,
			Color:      t.Course.Color,
		},
	}

	t.UpdateCourseReqMock = &proto.UpdateCourseRequest{
		Id:         t.Course.ID.String(),
		CourseCode: t.Course.CourseCode,
		Name:       t.Course.Name,
		Color:      t.Course.Color,
	}

	t.UpdateCourse = &course.Course{
		CourseCode: t.Course.CourseCode,
		Name:       t.Course.Name,
		Color:      t.Course.Color,
	}
}

func (t *CourseServiceTest) TestFindAllSuccess() {
	want := &proto.FindAllCourseResponse{Courses: createCourseDto(t.Courses)}

	var coursesIn []*course.Course

	repo := mock.RepositoryMock{}
	repo.On("FindAll", &coursesIn).Return(&t.Courses, nil)

	srv := NewService(&repo)
	actual, err := srv.FindAll(context.Background(), &proto.FindAllCourseRequest{})

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func createCourseDto(in []*course.Course) []*proto.Course {
	var result []*proto.Course

	for _, b := range in {
		r := &proto.Course{
			Id:         b.ID.String(),
			CourseCode: b.CourseCode,
			Name:       b.Name,
			Color:      b.Color,
		}

		result = append(result, r)
	}

	return result
}

func (t *CourseServiceTest) TestCreateSuccess() {
	want := &proto.CreateCourseResponse{Course: t.CourseDto}

	in := &course.Course{
		CourseCode: t.Course.CourseCode,
		Name:       t.Course.Name,
		Color:      t.Course.Color,
	}

	repo := &mock.RepositoryMock{}
	repo.On("Create", in).Return(t.Course, nil)

	srv := NewService(repo)
	actual, err := srv.Create(context.Background(), t.CreateCourseReqMock)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *CourseServiceTest) TestCreateInternalErr() {
	in := &course.Course{
		CourseCode: t.Course.CourseCode,
		Name:       t.Course.Name,
		Color:      t.Course.Color,
	}

	repo := &mock.RepositoryMock{}
	repo.On("Create", in).Return(nil, errors.New("something wrong"))

	srv := NewService(repo)
	actual, err := srv.Create(context.Background(), t.CreateCourseReqMock)

	st, ok := status.FromError(err)

	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.Internal, st.Code())
}

func (t *CourseServiceTest) TestUpdateSuccess() {
	want := &proto.UpdateCourseResponse{Course: t.CourseDto}

	repo := &mock.RepositoryMock{}
	repo.On("Update", t.Course.ID.String(), t.UpdateCourse).Return(t.Course, nil)

	srv := NewService(repo)
	actual, err := srv.Update(context.Background(), t.UpdateCourseReqMock)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *CourseServiceTest) TestUpdateNotFound() {
	repo := &mock.RepositoryMock{}
	repo.On("Update", t.Course.ID.String(), t.UpdateCourse).Return(nil, errors.New("Not found Course"))

	srv := NewService(repo)
	actual, err := srv.Update(context.Background(), t.UpdateCourseReqMock)

	st, ok := status.FromError(err)

	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.NotFound, st.Code())
}

func (t *CourseServiceTest) TestDeleteSuccess() {
	want := &proto.DeleteCourseResponse{Success: true}

	repo := &mock.RepositoryMock{}
	repo.On("Delete", t.Course.ID.String()).Return(nil)

	srv := NewService(repo)
	actual, err := srv.Delete(context.Background(), &proto.DeleteCourseRequest{Id: t.CourseDto.Id})

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), want, actual)
}

func (t *CourseServiceTest) TestDeleteNotFound() {
	repo := &mock.RepositoryMock{}
	repo.On("Delete", t.Course.ID.String()).Return(errors.New("Not found Course"))

	srv := NewService(repo)
	actual, err := srv.Delete(context.Background(), &proto.DeleteCourseRequest{Id: t.CourseDto.Id})

	st, ok := status.FromError(err)
	assert.True(t.T(), ok)
	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), codes.NotFound, st.Code())
}
