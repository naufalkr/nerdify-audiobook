package common

import "regexp"

func IsValidEmail(email string) bool {
	// Define the regex pattern for validating email
	regexPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regexPattern)

	// Match the email with the regex pattern
	return re.MatchString(email)
}
