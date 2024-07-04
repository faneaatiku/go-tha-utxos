package services

import (
	"regexp"
	"strconv"
)

func RemoveLineBreaks(text string) string {
	re := regexp.MustCompile(`\x{000D}\x{000A}|[\x{000A}\x{000B}\x{000C}\x{000D}\x{0085}\x{2028}\x{2029}]`)
	return re.ReplaceAllString(text, ``)
}

func FloatToString(input float64) string {
	return strconv.FormatFloat(input, 'f', -1, 64)
}
