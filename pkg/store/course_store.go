package store

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/db"
)

type CourseStore struct {
	db *db.DB
}

type CourseData map[string]map[string]string

type GradeChange struct {
	Course   string
	OldGrade string
	NewGrade string
	Semester string
}

func NewCourseStore(db *db.DB) (*CourseStore, error) {
	// Créer les tables si elles n'existent pas
	if err := createTables(db); err != nil {
		return nil, err
	}
	return &CourseStore{db: db}, nil
}

func createTables(db *db.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS courses (
			id VARCHAR(50) PRIMARY KEY,
			email VARCHAR(255) NOT NULL,
			semester VARCHAR(10) NOT NULL,
			course_code VARCHAR(20) NOT NULL,
			grade VARCHAR(10) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(email, semester, course_code)
		);
	`

	_, err := db.Exec(query)
	return err
}

func (s *CourseStore) SaveCourses(email, semester string, courses map[string]string) ([]GradeChange, error) {
	var changes []GradeChange

	// Démarrer une transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("erreur lors du début de la transaction: %w", err)
	}
	defer tx.Rollback()

	// Pour chaque cours
	for course, grade := range courses {
		// Vérifier si le cours existe déjà
		var oldGrade string
		var courseID string
		err := squirrel.Select("id", "grade").
			From("courses").
			Where(squirrel.Eq{
				"email":       email,
				"semester":    semester,
				"course_code": course,
			}).
			RunWith(tx).
			QueryRow().
			Scan(&courseID, &oldGrade)

		if err == sql.ErrNoRows {
			// Nouveau cours
			courseID = GenerateID(CoursePrefix)
			_, err = squirrel.Insert("courses").
				Columns("id", "email", "semester", "course_code", "grade").
				Values(courseID, email, semester, course, grade).
				RunWith(tx).
				Exec()
			if err != nil {
				return nil, fmt.Errorf("erreur lors de l'insertion du cours: %w", err)
			}
		} else if err != nil {
			return nil, fmt.Errorf("erreur lors de la vérification du cours: %w", err)
		} else if oldGrade != grade {
			// Mise à jour de la note
			_, err = squirrel.Update("courses").
				Set("grade", grade).
				Set("updated_at", squirrel.Expr("CURRENT_TIMESTAMP")).
				Where(squirrel.Eq{"id": courseID}).
				RunWith(tx).
				Exec()
			if err != nil {
				return nil, fmt.Errorf("erreur lors de la mise à jour de la note: %w", err)
			}

			changes = append(changes, GradeChange{
				Course:   course,
				OldGrade: oldGrade,
				NewGrade: grade,
				Semester: semester,
			})
		}
	}

	// Valider la transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("erreur lors de la validation de la transaction: %w", err)
	}

	return changes, nil
}

func (s *CourseStore) GetCourses(email, semester string) (map[string]string, error) {
	rows, err := squirrel.Select("course_code", "grade").
		From("courses").
		Where(squirrel.Eq{
			"email":    email,
			"semester": semester,
		}).
		RunWith(s.db).
		Query()
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des cours: %w", err)
	}
	defer rows.Close()

	courses := make(map[string]string)
	for rows.Next() {
		var course, grade string
		if err := rows.Scan(&course, &grade); err != nil {
			return nil, fmt.Errorf("erreur lors de la lecture des cours: %w", err)
		}
		courses[course] = grade
	}

	return courses, nil
}

func (s *CourseStore) DeleteCourse(email, semester, course string) error {
	var courseID string
	err := squirrel.Select("id").
		From("courses").
		Where(squirrel.Eq{
			"email":       email,
			"semester":    semester,
			"course_code": course,
		}).
		RunWith(s.db).
		QueryRow().
		Scan(&courseID)

	if err == sql.ErrNoRows {
		return fmt.Errorf("cours non trouvé")
	}
	if err != nil {
		return fmt.Errorf("erreur lors de la récupération de l'ID du cours: %w", err)
	}

	if err := ValidateID(courseID, CoursePrefix); err != nil {
		return fmt.Errorf("ID de cours invalide: %w", err)
	}

	result, err := squirrel.Delete("courses").
		Where(squirrel.Eq{"id": courseID}).
		RunWith(s.db).
		Exec()
	if err != nil {
		return fmt.Errorf("erreur lors de la suppression du cours: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erreur lors de la vérification de la suppression: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("cours non trouvé")
	}

	return nil
}

func (s *CourseStore) UpdateGrade(email, semester, course, grade string) (*GradeChange, error) {
	var oldGrade string
	var courseID string
	err := squirrel.Select("id", "grade").
		From("courses").
		Where(squirrel.Eq{
			"email":       email,
			"semester":    semester,
			"course_code": course,
		}).
		RunWith(s.db).
		QueryRow().
		Scan(&courseID, &oldGrade)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("cours non trouvé")
	}
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération de la note: %w", err)
	}

	if err := ValidateID(courseID, CoursePrefix); err != nil {
		return nil, fmt.Errorf("ID de cours invalide: %w", err)
	}

	if oldGrade == grade {
		return nil, nil
	}

	_, err = squirrel.Update("courses").
		Set("grade", grade).
		Set("updated_at", squirrel.Expr("CURRENT_TIMESTAMP")).
		Where(squirrel.Eq{"id": courseID}).
		RunWith(s.db).
		Exec()
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la mise à jour de la note: %w", err)
	}

	return &GradeChange{
		Course:   course,
		OldGrade: oldGrade,
		NewGrade: grade,
		Semester: semester,
	}, nil
}

func (s *CourseStore) DeleteCourses(email, semester string) error {
	_, err := squirrel.Delete("courses").
		Where(squirrel.Eq{
			"email":    email,
			"semester": semester,
		}).
		RunWith(s.db).
		Exec()
	if err != nil {
		return fmt.Errorf("erreur lors de la suppression des cours: %w", err)
	}

	return nil
}
