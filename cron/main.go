package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/felixlheureux/uqam-grade-notifier/pkg/alert"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/auth"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/db"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/grade"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/helper"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/store"
)

type Config struct {
	DBConnString     string `json:"db_conn_string"`
	GmailAppPassword string `json:"gmail_app_password"`
	MailFrom         string `json:"mail_from"`
	EncryptionKey    string `json:"encryption_key"`
	LogPath          string `json:"log_path"`
}

func setupLogging(logPath string) (*os.File, error) {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		return nil, fmt.Errorf("error creating logs directory: %v", err)
	}

	// Open log file in append mode
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("error opening log file: %v", err)
	}

	// Configure logger to write to file
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	return logFile, nil
}

func main() {
	configPath := flag.String("config", "cron/config.json", "Path to configuration file")
	flag.Parse()

	var cfg Config
	helper.MustLoadConfig(*configPath, &cfg)

	// Configure logging
	logFile, err := setupLogging(cfg.LogPath)
	if err != nil {
		fmt.Printf("Error configuring logs: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()

	log.Printf("Starting grade check")

	// Initialize database connection
	database, err := db.New(cfg.DBConnString)
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		os.Exit(1)
	}
	defer database.Close()

	// Initialize stores
	userStore, err := store.NewUserStore(database, []byte(cfg.EncryptionKey))
	if err != nil {
		log.Printf("Error initializing user store: %v", err)
		os.Exit(1)
	}

	courseStore, err := store.NewCourseStore(database)
	if err != nil {
		log.Printf("Error initializing course store: %v", err)
		os.Exit(1)
	}

	// Get all users
	emails, err := userStore.GetAllUsers()
	if err != nil {
		log.Printf("Error retrieving users: %v", err)
		os.Exit(1)
	}

	if len(emails) == 0 {
		log.Printf("No users found")
		os.Exit(0)
	}

	log.Printf("Checking grades for %d users", len(emails))

	// For each user
	for _, email := range emails {
		log.Printf("Processing user %s", email)

		// Get UQAM credentials
		uqamUsername, uqamPassword, err := userStore.GetUser(email)
		if err != nil {
			log.Printf("Error for user %s: %v", email, err)
			continue
		}

		// Generate UQAM token
		token := auth.MustGenerateToken(uqamUsername, uqamPassword)

		// Get user's courses
		courses, err := courseStore.GetCourses(email, "20251") // TODO: Handle different semesters
		if err != nil {
			log.Printf("Error for user %s: %v", email, err)
			continue
		}

		if len(courses) == 0 {
			log.Printf("No courses found for user %s", email)
			continue
		}

		log.Printf("Checking %d courses for user %s", len(courses), email)

		// For each course
		for courseCode := range courses {
			log.Printf("Checking course %s for user %s", courseCode, email)

			// Get grade from UQAM
			newGrade, err := grade.FetchGrade(token, "20251", courseCode)
			if err != nil {
				log.Printf("Error for course %s of user %s: %v", courseCode, email, err)
				continue
			}

			// Update grade if it changed
			change, err := courseStore.UpdateGrade(email, "20251", courseCode, newGrade)
			if err != nil {
				log.Printf("Error updating grade for course %s of user %s: %v", courseCode, email, err)
				continue
			}

			// If grade changed, send notification
			if change != nil {
				subject := "Grade Update"
				body := fmt.Sprintf("Grade for course %s has been changed from %s to %s.", courseCode, change.OldGrade, change.NewGrade)

				if err := alert.SendEmail(cfg.GmailAppPassword, cfg.MailFrom, email, subject, body); err != nil {
					log.Printf("Error sending notification for user %s: %v", email, err)
					continue
				}

				log.Printf("Notification sent to %s for course %s", email, courseCode)
			}
		}

		// Wait a bit between each user to avoid overwhelming the UQAM server
		time.Sleep(2 * time.Second)
	}

	log.Printf("Grade check completed")
	os.Exit(0)
}
