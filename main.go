package main

import (
	"fmt"
	"github.com/felixlheureux/uqam-grade-notifier/pkg/auth"
)

type config struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {

	token := auth.MustGenerateToken("SERF08069508", "14795")
	fmt.Println()
}
