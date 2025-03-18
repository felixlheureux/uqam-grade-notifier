package store

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type CourseStore struct {
	baseDir string
}

type CourseData map[string]map[string]string

type GradeChange struct {
	Course   string
	OldGrade string
	NewGrade string
	Semester string
}

func NewCourseStore(baseDir string) (*CourseStore, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create course store directory: %w", err)
	}
	return &CourseStore{baseDir: baseDir}, nil
}

func (s *CourseStore) getFilePath(email string) string {
	return filepath.Join(s.baseDir, fmt.Sprintf("%s.json", email))
}

func (s *CourseStore) load(email string) (CourseData, error) {
	filePath := s.getFilePath(email)
	data := make(CourseData)

	// Si le fichier n'existe pas, retourner une map vide
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return data, nil
	}

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read course file: %w", err)
	}

	if err := json.Unmarshal(file, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal course data: %w", err)
	}

	return data, nil
}

func (s *CourseStore) save(email string, data CourseData) error {
	filePath := s.getFilePath(email)
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal course data: %w", err)
	}

	if err := ioutil.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write course file: %w", err)
	}

	return nil
}

func (s *CourseStore) SaveCourses(email, semester string, courses map[string]string) ([]GradeChange, error) {
	data, err := s.load(email)
	if err != nil {
		return nil, err
	}

	if data[semester] == nil {
		data[semester] = make(map[string]string)
	}

	var changes []GradeChange
	for course, grade := range courses {
		if oldGrade, exists := data[semester][course]; exists && oldGrade != grade {
			changes = append(changes, GradeChange{
				Course:   course,
				OldGrade: oldGrade,
				NewGrade: grade,
				Semester: semester,
			})
		}
		data[semester][course] = grade
	}

	if err := s.save(email, data); err != nil {
		return nil, err
	}

	return changes, nil
}

func (s *CourseStore) GetCourses(email, semester string) (map[string]string, error) {
	data, err := s.load(email)
	if err != nil {
		return nil, err
	}

	if courses, exists := data[semester]; exists {
		return courses, nil
	}

	return make(map[string]string), nil
}

func (s *CourseStore) DeleteCourse(email, semester, course string) error {
	data, err := s.load(email)
	if err != nil {
		return err
	}

	if _, exists := data[semester]; exists {
		delete(data[semester], course)
		if len(data[semester]) == 0 {
			delete(data, semester)
		}
		return s.save(email, data)
	}

	return nil
}

func (s *CourseStore) UpdateGrade(email, semester, course, grade string) (*GradeChange, error) {
	data, err := s.load(email)
	if err != nil {
		return nil, err
	}

	if _, exists := data[semester]; !exists {
		return nil, fmt.Errorf("semester %s does not exist", semester)
	}

	oldGrade, exists := data[semester][course]
	if !exists {
		return nil, fmt.Errorf("course %s does not exist in semester %s", course, semester)
	}

	if oldGrade == grade {
		return nil, nil
	}

	change := &GradeChange{
		Course:   course,
		OldGrade: oldGrade,
		NewGrade: grade,
		Semester: semester,
	}

	data[semester][course] = grade
	if err := s.save(email, data); err != nil {
		return nil, err
	}

	return change, nil
}

func (s *CourseStore) DeleteCourses(email, semester string) error {
	data, err := s.load(email)
	if err != nil {
		return err
	}

	delete(data, semester)
	return s.save(email, data)
}
