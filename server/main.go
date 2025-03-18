package main

import (
	"flag"
	"log"

	"github.com/felixlheureux/uqam-grade-notifier/pkg/auth"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/db"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/helper"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/store"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Config struct {
	DBConnString     string `json:"db_conn_string"`
	Port             string `json:"port"`
	JWTSecret        string `json:"jwt_secret"`
	BaseURL          string `json:"base_url"`
	GmailAppPassword string `json:"gmail_app_password"`
	MailFrom         string `json:"mail_from"`
}

func main() {
	config := &Config{}
	flag.StringVar(&config.DBConnString, "db", "", "Chaîne de connexion à la base de données PostgreSQL")
	flag.StringVar(&config.Port, "port", "8080", "Port du serveur")
	flag.Parse()

	// Charger la configuration
	helper.MustLoadConfig("config/dev.config.json", config)

	// Initialiser la connexion à la base de données
	database, err := db.New(config.DBConnString)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// Initialiser les stores
	tokenStore, err := store.NewTokenStore(database)
	if err != nil {
		log.Fatal(err)
	}

	courseStore, err := store.NewCourseStore(database)
	if err != nil {
		log.Fatal(err)
	}

	// Initialiser le gestionnaire d'authentification
	tokenManager := auth.NewTokenManager(config.JWTSecret)
	handler := auth.NewHandler(tokenStore, tokenManager, courseStore, config.BaseURL, config.GmailAppPassword, config.MailFrom)

	// Configurer le serveur Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes d'authentification
	e.POST("/auth/request-login", handler.RequestLogin)
	e.GET("/auth/login", handler.Login)

	// Routes des cours (protégées)
	courses := e.Group("/courses")
	courses.Use(handler.RequireAuth)
	courses.POST("/:semester", handler.SaveCourses)
	courses.GET("/:semester", handler.GetCourses)
	courses.DELETE("/:semester/:course", handler.DeleteCourse)
	courses.PUT("/:semester/:course/grade", handler.UpdateGrade)
	courses.DELETE("/:semester", handler.DeleteCourses)

	// Démarrer le serveur
	e.Logger.Fatal(e.Start(":" + config.Port))
}
