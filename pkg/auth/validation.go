package auth

import (
	"fmt"
	"strings"
)

const allowedDomain = "@ens.uqam.ca"

func ValidateEmail(email string) error {
	if !strings.HasSuffix(email, allowedDomain) {
		return fmt.Errorf("only emails from domain %s are allowed", allowedDomain)
	}
	return nil
}
