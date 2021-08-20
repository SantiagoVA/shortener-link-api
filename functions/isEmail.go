package functions

import "strings"

func IsEmail(email string) bool {
	if len(email) > 254 || len(email) == 0 {
		return false
	}
	atIndex := strings.Index(email, "@")
	if atIndex == -1 {
		return false
	}
	if strings.Index(email[atIndex+1:], ".") == -1 {
		return false
	}
	return true
}
