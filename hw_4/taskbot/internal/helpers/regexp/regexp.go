package regexp

import "regexp"

func CheckRegExp(regExp string, str string) bool {
	ok, err := regexp.MatchString(regExp, str)
	if err != nil {
		return false
	}
	return ok
}
