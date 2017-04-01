package stringutil

import (
	"fmt"
	"regexp"
)

type stringArray []string

var numberRegex = regexp.MustCompile(`[0-9]`)

func ToInterfaces(array []string) []interface{} {
	ret := make([]interface{}, 0, len(array))
	for i := 0; i < len(array); i++ {
		ret = append(ret, array[i])
	}
	return ret
}

// IsNumber returns true, if argument contains only numeric.
func IsNumber(s string) bool {
	return numberRegex.MatchString(s)
}

type Cutter struct {
	maxLength int
	tail      string
}

var defaultCutter = Cutter{
	maxLength: 10,
	tail:      "...",
}

func (c Cutter) Cut(s string, length int) string {
	if len(s) > c.maxLength {
		return fmt.Sprintf("%s%s", s[:length], c.tail)
	}
	return s
}

func SetCutConfig(maxLength int, tail string) {
	defaultCutter.maxLength = maxLength
	defaultCutter.tail = tail
}

func Cut(s string, length int) string {
	return defaultCutter.Cut(s, length)
}
