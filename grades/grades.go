package grades

import (
	"fmt"
	"sync"
)

type Student struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Grades    []Grade
}

func (s *Student) Average() float32 {
	var result float32
	for _, grade := range s.Grades {
		result += grade.Score
	}
	return result / float32(len(s.Grades))
}

type Students []Student

var (
	students      Students
	studentsMutex sync.Mutex
)

func (ss Students) GetByID(id int) (*Student, error) {
	studentsMutex.Lock()
	defer studentsMutex.Unlock()
	for i := range ss {
		if ss[i].ID == id {
			return &ss[i], nil
		}
	}
	return nil, fmt.Errorf("student %v not found", id)
}

type GradeType string

const (
	GradeQuiz = GradeType("Quiz")
	GradeTest = GradeType("Test")
	GradeExam = GradeType("Exam")
)

type Grade struct {
	Title string    `json:"title"`
	Type  GradeType `json:"type"`
	Score float32   `json:"score"`
}