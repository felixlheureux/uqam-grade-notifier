package store

import (
	"encoding/json"
	"io/ioutil"
)

const file = "grades.json"

type GradeStore map[string]map[string]string

func loadGrades() GradeStore {
	grades := make(GradeStore)
	file, err := ioutil.ReadFile(file)
	if err != nil {
		return grades
	}
	json.Unmarshal(file, &grades)
	return grades
}

func SaveGrade(semester, course, grade string) error {
	grades := loadGrades()
	if grades[semester] == nil {
		grades[semester] = make(map[string]string)
	}
	grades[semester][course] = grade
	data, err := json.MarshalIndent(grades, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, data, 0644)
}

func GetGrade(semester, course string) string {
	grades := loadGrades()
	if semesterGrades, exists := grades[semester]; exists {
		if grade, found := semesterGrades[course]; found {
			return grade
		}
	}
	return "0.00"
}
