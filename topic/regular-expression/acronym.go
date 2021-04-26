package acronym

import "regexp"
import "strings"

func Abbreviate(s string) string {
	re := regexp.MustCompile(`[^a-zA-Z]*([a-zA-Z])[a-zA-Z']*`)
	output := re.ReplaceAllString(s, "$1")
	return strings.ToUpper(output)
}