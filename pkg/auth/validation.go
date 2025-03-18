package auth

import (
	"fmt"
	"strings"
)

const allowedDomain = "@ens.uqam.ca"

func ValidateEmail(email string) error {
	if !strings.HasSuffix(email, allowedDomain) {
		return fmt.Errorf("seuls les emails du domaine %s sont autorisés", allowedDomain)
	}
	return nil
}
