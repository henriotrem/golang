package bob

import "regexp"

// My solution
const (
	upperCase = `[A-Z]`
	lowerCase = `[a-z]`
	questionMark = `\?\s*$`
	silence = `^\s*$`
)

func Hey(remark string) string {

	containUpperCase, _ := regexp.MatchString(upperCase, remark)
	containLowerCase, _ := regexp.MatchString(lowerCase, remark)
	containQuestionMark, _ := regexp.MatchString(questionMark, remark)
	containSilence, _ := regexp.MatchString(silence, remark)

	switch {
		case containUpperCase && !containLowerCase && containQuestionMark:
			return "Calm down, I know what I'm doing!"
		case containUpperCase && !containLowerCase:
			return "Whoa, chill out!"
		case containQuestionMark:
			return "Sure."
		case containSilence:
			return "Fine. Be that way!"
		default:
			return "Whatever."
	}
}


