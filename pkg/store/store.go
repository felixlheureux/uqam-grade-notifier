package store

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type store struct {
	path string
}

func New(path string) (*store, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			ioutil.WriteFile(path, []byte("{}"), 0644)
		} else {
			return nil, err
		}
	}
	return &store{path}, nil
}

type gradeStore map[string]map[string]string

func (store *store) loadGrades() (gradeStore, error) {
	grades := make(gradeStore)
	file, err := ioutil.ReadFile(store.path)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(file, &grades)
	return grades, nil
}

func (store *store) SaveGrade(semester, course, grade string) error {
	grades, err := store.loadGrades()
	if err != nil {
		return err
	}
	if grades[semester] == nil {
		grades[semester] = make(map[string]string)
	}
	grades[semester][course] = grade
	data, err := json.MarshalIndent(grades, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(store.path, data, 0644)
}

func (store *store) GetGrade(semester, course string) (string, error) {
	grades, err := store.loadGrades()
	if err != nil {
		return "", err
	}
	if semesterGrades, exists := grades[semester]; exists {
		if grade, found := semesterGrades[course]; found {
			return grade, nil
		}
	}
	return "0.00", nil
}
