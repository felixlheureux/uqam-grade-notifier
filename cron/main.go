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
	// Créer le dossier logs s'il n'existe pas
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		return nil, fmt.Errorf("erreur lors de la création du dossier logs: %v", err)
	}

	// Ouvrir le fichier de log en mode append
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de l'ouverture du fichier de log: %v", err)
	}

	// Configurer le logger pour écrire dans le fichier
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	return logFile, nil
}

func main() {
	configPath := flag.String("config", "cron/config.json", "Chemin du fichier de configuration")
	flag.Parse()

	var cfg Config
	helper.MustLoadConfig(*configPath, &cfg)

	// Configurer le logging
	logFile, err := setupLogging(cfg.LogPath)
	if err != nil {
		fmt.Printf("Erreur lors de la configuration des logs: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()

	log.Printf("Démarrage de la vérification des notes")

	// Initialiser la connexion à la base de données
	database, err := db.New(cfg.DBConnString)
	if err != nil {
		log.Printf("Erreur lors de la connexion à la base de données: %v", err)
		os.Exit(1)
	}
	defer database.Close()

	// Initialiser les stores
	userStore, err := store.NewUserStore(database, []byte(cfg.EncryptionKey))
	if err != nil {
		log.Printf("Erreur lors de l'initialisation du store utilisateur: %v", err)
		os.Exit(1)
	}

	courseStore, err := store.NewCourseStore(database)
	if err != nil {
		log.Printf("Erreur lors de l'initialisation du store des cours: %v", err)
		os.Exit(1)
	}

	// Récupérer tous les utilisateurs
	emails, err := userStore.GetAllUsers()
	if err != nil {
		log.Printf("Erreur lors de la récupération des utilisateurs: %v", err)
		os.Exit(1)
	}

	if len(emails) == 0 {
		log.Printf("Aucun utilisateur trouvé")
		os.Exit(0)
	}

	log.Printf("Vérification des notes pour %d utilisateurs", len(emails))

	// Pour chaque utilisateur
	for _, email := range emails {
		log.Printf("Traitement de l'utilisateur %s", email)

		// Récupérer les identifiants UQAM
		uqamUsername, uqamPassword, err := userStore.GetUser(email)
		if err != nil {
			log.Printf("Erreur pour l'utilisateur %s: %v", email, err)
			continue
		}

		// Générer le token UQAM
		token := auth.MustGenerateToken(uqamUsername, uqamPassword)

		// Récupérer les cours de l'utilisateur
		courses, err := courseStore.GetCourses(email, "20251") // TODO: Gérer les différents semestres
		if err != nil {
			log.Printf("Erreur pour l'utilisateur %s: %v", email, err)
			continue
		}

		if len(courses) == 0 {
			log.Printf("Aucun cours trouvé pour l'utilisateur %s", email)
			continue
		}

		log.Printf("Vérification de %d cours pour l'utilisateur %s", len(courses), email)

		// Pour chaque cours
		for courseCode := range courses {
			log.Printf("Vérification du cours %s pour l'utilisateur %s", courseCode, email)

			// Récupérer la note depuis UQAM
			newGrade, err := grade.FetchGrade(token, "20251", courseCode)
			if err != nil {
				log.Printf("Erreur pour le cours %s de l'utilisateur %s: %v", courseCode, email, err)
				continue
			}

			// Mettre à jour la note si elle a changé
			change, err := courseStore.UpdateGrade(email, "20251", courseCode, newGrade)
			if err != nil {
				log.Printf("Erreur lors de la mise à jour de la note pour le cours %s de l'utilisateur %s: %v", courseCode, email, err)
				continue
			}

			// Si la note a changé, envoyer une notification
			if change != nil {
				subject := "Modification de note"
				body := fmt.Sprintf("La note du cours %s a été modifiée de %s à %s.", courseCode, change.OldGrade, change.NewGrade)

				if err := alert.SendEmail(cfg.GmailAppPassword, cfg.MailFrom, email, subject, body); err != nil {
					log.Printf("Erreur lors de l'envoi de la notification pour l'utilisateur %s: %v", email, err)
					continue
				}

				log.Printf("Notification envoyée à %s pour le cours %s", email, courseCode)
			}
		}

		// Attendre un peu entre chaque utilisateur pour ne pas surcharger le serveur UQAM
		time.Sleep(2 * time.Second)
	}

	log.Printf("Vérification des notes terminée")
	os.Exit(0)
}
