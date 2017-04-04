package stringutil

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
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

func NotEmpty(array []string) bool {
	return len(array) != 0
}

func Empty(array []string) bool {
	return len(array) != 0
}

// ToString converts interface array to string array.
func ToStrings(array []interface{}) ([]string, error) {
	ret := make([]string, 0, len(array))
	for i := 0; i < len(array); i++ {
		v, ok := array[i].(string)
		if !ok {
			return nil, errors.New(fmt.Sprintf("Can not convert to []string. %#v", array[i]))
		}

		ret = append(ret, v)
	}
	return ret, nil
}

func StringToStringArray(s string, separator string) []string {
	var ret []string
	if s == "" {
		ret = make([]string, 0)
	} else {
		if strings.Contains(s, separator) {
			fmt.Println(s)
			ret = strings.Split(s, separator)
			fmt.Println(ret)

		} else {
			ret = make([]string, 1)
			ret[0] = s
		}
	}

	return ret
}

func Map(array []string, fn func(s string) string) []string {
	var ret []string = make([]string, len(array))
	for i, v := range array {
		ret[i] = fn(v)
	}
	return ret
}

func Contains(array []string, s string) bool {
	for _, v := range array {
		if v == s {
			return true
		}
	}
	return false
}

// IsNumber returns true, if argument contains only numeric.
func IsNumber(s string) bool {
	return numberRegex.MatchString(s)
}

// String Cutter
type Cutter struct {
	maxLength int
	tail      string
}

var defaultCutter = Cutter{
	maxLength: 10,
	tail:      "...",
}

func (c Cutter) Cut(s string) string {
	if len(s) > c.maxLength {
		return fmt.Sprintf("%s%s", s[:c.maxLength], c.tail)
	}
	return s
}

func SetCutConfig(maxLength int, tail string) {
	defaultCutter.maxLength = maxLength
	defaultCutter.tail = tail
}

func Cut(s string) string {
	return defaultCutter.Cut(s)
}
