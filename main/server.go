package main

import (
	"flag"
	"log"

	"github.com/felixlheureux/uqam-grade-notifier/pkg/auth"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/helper"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/store"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	serverConfigPath = flag.String("config", "main/config/dev.config.json", "Path to the configuration file")
	serverPort       = flag.String("port", "8080", "Port to listen on")
)

type serverConfig struct {
	GmailAppPassword string `json:"gmail_app_password"`
	MailFrom         string `json:"mail_from"`
	StorePath        string `json:"store_path"`
	BaseURL          string `json:"base_url"`
	JWTSecret        string `json:"jwt_secret"`
}

func startServer() {
	flag.Parse()

	var cfg serverConfig
	helper.MustLoadConfig(*serverConfigPath, &cfg)

	tokenStore, err := store.NewTokenStore(cfg.StorePath)
	if err != nil {
		log.Fatal(err)
	}

	courseStore, err := store.NewCourseStore("main/data/courses")
	if err != nil {
		log.Fatal(err)
	}

	tokenManager := auth.NewTokenManager(cfg.JWTSecret)
	authHandler := auth.NewHandler(tokenStore, tokenManager, courseStore, cfg.GmailAppPassword, cfg.MailFrom, cfg.BaseURL)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes d'authentification
	e.POST("/auth/request-login", authHandler.RequestLogin)
	e.GET("/auth/login", authHandler.Login)

	// Routes des cours
	e.POST("/courses", authHandler.SaveCourses)
	e.GET("/courses", authHandler.GetCourses)
	e.PUT("/courses/:semester/:course", authHandler.UpdateGrade)
	e.DELETE("/courses/:semester/:course", authHandler.DeleteCourse)

	// Démarrage du serveur
	e.Logger.Fatal(e.Start(":" + *serverPort))
}
