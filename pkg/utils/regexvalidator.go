package utils

import "regexp"

func IsValidWithRegex(expression, str string) bool {
	re := regexp.MustCompile(expression)
	return re.MatchString(str)
}
