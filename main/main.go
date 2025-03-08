package main

import (
	"github.com/felixlheureux/uqam-grade-notifier/pkg/alert"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/auth"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/grade"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/helper"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/store"
	"log"
)

const config_path = "main/config.json"

type config struct {
	Username         string   `json:"username"`
	Password         string   `json:"password"`
	GmailAppPassword string   `json:"gmail_app_password"`
	MailTo           string   `json:"mail_to"`
	MailFrom         string   `json:"mail_from"`
	Semester         string   `json:"semester"`
	Courses          []string `json:"courses"`
	StorePath        string   `json:"store_path"`
}

func main() {
	config := config{}
	helper.MustLoadConfig(config_path, &config)
	token := auth.MustGenerateToken(config.Username, config.Password)
	_store, err := store.New(config.StorePath)
	if err != nil {
		log.Fatal(err)
	}
	for _, course := range config.Courses {
		newGrade, err := grade.FetchGrade(token, config.Semester, course)
		if err != nil {
			log.Fatal(err)
		}
		oldGrade, err := _store.GetGrade(config.Semester, course)
		if err != nil {
			log.Fatal(err)
		}
		if newGrade != oldGrade {
			if err := _store.SaveGrade(config.Semester, course, newGrade); err != nil {
				log.Fatal(err)
			}
			if err = alert.SendEmail(config.GmailAppPassword, config.MailFrom, config.MailTo, course, newGrade); err != nil {
				log.Fatal(err)
			}
		}
	}
}
