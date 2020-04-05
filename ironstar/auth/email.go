package auth

import (
	"errors"
	"regexp"
	"strings"

	"github.com/ironstar-io/tokaido/ironstar/utils"
)

// ValidateEmail - Ensure the supplied email address is valid
func ValidateEmail(email string) error {
	// Validator will handle the requirement
	if email == "" {
		return errors.New("Empty email suppied. Code: aG3rdr")
	}

	userRegexp := regexp.MustCompile("^[a-zA-Z0-9!#$%&'*+/=?^_`{|}~.-]+$")
	hostRegexp := regexp.MustCompile("^[^\\s]+\\.[^\\s]+$")
	userDotRegexp := regexp.MustCompile("(^[.]{1})|([.]{1}$)|([.]{2,})")
	if len(email) < 6 || len(email) > 254 {
		return errors.New("Invalid email suppied '" + email + "'. Code: sCkqzE")
	}
	at := strings.LastIndex(email, "@")
	if at <= 0 || at > len(email)-3 {
		return errors.New("Invalid email suppied '" + email + "'. Code: nRsEkq")
	}
	user := email[:at]
	host := email[at+1:]
	if len(user) > 64 {
		return errors.New("Invalid email suppied '" + email + "'. Code: uPQ3ov")
	}
	if userDotRegexp.MatchString(user) || !userRegexp.MatchString(user) || !hostRegexp.MatchString(host) {
		return errors.New("Invalid email suppied '" + email + "'. Code: t6ud7X")
	}

	return nil
}

func GetCLIEmail(args []string) (string, error) {
	var email string
	if len(args) == 0 {
		input, err := utils.StdinPrompt("Email: ")
		if err != nil {
			return "", err
		}
		email = input
	} else {
		email = args[0]
	}

	err := ValidateEmail(email)
	if err != nil {
		return "", err
	}

	return email, nil
}
