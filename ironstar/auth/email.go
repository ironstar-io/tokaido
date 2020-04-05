package auth

import (
	"errors"
	"regexp"
	"strings"
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

func GetUserEmail(args []string) (string, error) {
	if len(args) == 0 {
		return args[0], nil
	}

	return args[0], nil
}
