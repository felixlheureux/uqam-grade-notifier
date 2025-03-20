package main

import (
	"flag"
	"log"

	"github.com/felixlheureux/uqam-grade-notifier/pkg/auth"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/db"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/helper"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/service"
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
	flag.StringVar(&config.DBConnString, "db", "", "PostgreSQL database connection string")
	flag.StringVar(&config.Port, "port", "8080", "Server port")
	flag.Parse()

	// Load configuration
	helper.MustLoadConfig("config/dev.config.json", config)

	// Initialize database connection
	database, err := db.New(config.DBConnString)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// Initialize services
	authService, err := service.NewAuthService(database, config.JWTSecret)
	if err != nil {
		log.Fatal(err)
	}

	courseService, err := service.NewCourseService(database)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize authentication handler
	handler := auth.NewHandler(authService, courseService, config.BaseURL, config.GmailAppPassword, config.MailFrom)

	// Configure Echo server
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Authentication routes
	e.POST("/auth/request-login", handler.RequestLogin)
	e.GET("/auth/login", handler.Login)

	// Course routes (protected)
	courses := e.Group("/courses")
	courses.Use(handler.RequireAuth)

	// GET /courses/{semester}
	courses.GET("/:semester", handler.GetCourses)

	// POST /courses/{semester}
	courses.POST("/:semester", handler.SaveCourses)

	// DELETE /courses/{semester}
	courses.DELETE("/:semester", handler.DeleteCourses)

	// DELETE /courses/{semester}/{course}
	courses.DELETE("/:semester/:course", handler.DeleteCourse)

	// PUT /courses/{semester}/{course}/grade
	courses.PUT("/:semester/:course/grade", handler.UpdateGrade)

	// Start server
	e.Logger.Fatal(e.Start(":" + config.Port))
}
