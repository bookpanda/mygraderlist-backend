package seed

import (
	"github.com/bookpanda/mygraderlist-backend/src/app/model/course"
)

var Courses = []course.Course{
	{
		CourseCode: "liked",
		Name:       "Liked Problems",
		Color:      "#6114c7",
	},
	{
		CourseCode: "2110211",
		Name:       "Data Structure",
		Color:      "#0285c7",
	},
	{
		CourseCode: "2110263",
		Name:       "Digital Logic Lab",
		Color:      "#854d0e",
	},
	{
		CourseCode: "2110327",
		Name:       "Algorithm Design",
		Color:      "#e2cb18",
	},
}
