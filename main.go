package main

import (
	"github.com/felixlheureux/uqam-grade-notifier/pkg/alert"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/auth"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/grade"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/helper"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/store"
	"log"
)

type config struct {
	Username         string   `json:"username"`
	Password         string   `json:"password"`
	GmailAppPassword string   `json:"gmail_app_password"`
	MailTo           string   `json:"mail_to"`
	MailFrom         string   `json:"mail_from"`
	Semester         string   `json:"semester"`
	Courses          []string `json:"courses"`
}

func main() {
	config := config{}
	helper.MustLoadConfig(&config)
	token := auth.MustGenerateToken(config.Username, config.Password)

	for _, course := range config.Courses {
		newGrade, err := grade.FetchGrade(token, config.Semester, course)
		if err != nil {
			log.Fatal(err)
		}

		oldGrade := store.GetGrade(config.Semester, course)

		if newGrade != oldGrade {
			if err := store.SaveGrade(config.Semester, course, newGrade); err != nil {
				log.Fatal(err)
			}

			if err = alert.SendEmail(config.MailTo, config.MailFrom, config.GmailAppPassword, course, newGrade); err != nil {
				log.Fatal(err)
			}
		}
		err = alert.SendEmail(config.MailTo, config.MailFrom, config.GmailAppPassword, course, newGrade)
		if err != nil {
			log.Fatal(err)
		}
	}
}
