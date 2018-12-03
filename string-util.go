package utils

import (
	"fmt"
	"regexp"
	"strings"
)

type StringUtil struct{}

const (
	// Numeric represents regular expression for numeric
	Numeric string = "^[0-9]+$"
)

var (
	regexNumeric = regexp.MustCompile(Numeric)
)

func PadLeft(sourceStr string, totalLen int, char string) string {
	if len(sourceStr) >= totalLen {
		return sourceStr
	} else {
		targetStrArr := []string{}
		for i := 0; i < totalLen-len(sourceStr); i++ {
			targetStrArr = append(targetStrArr, char)
		}
		targetStrArr = append(targetStrArr, sourceStr)
		return strings.Join(targetStrArr, "")
	}
}

func ToString(v interface{}) string {
	str, ok := v.(string)
	if !ok {
		str = fmt.Sprintf("%#v", v)
	}
	return str
}

func IsNumeric(str string) bool {
	return regexNumeric.MatchString(str)
}

// Should be replaced with strings.HasPrefix.
func StartWith(sourceStr string, targetStr string) bool {
	regexStr := strings.Join([]string{targetStr, "([\\s\\S]*?)"}, "")
	return regexp.MustCompile(regexStr).MatchString(sourceStr)
}

// Should be replaced with strings.HasSuffix.
func EndWith(sourceStr string, targetStr string) bool {
	regexStr := strings.Join([]string{"([\\s\\S]*?)", targetStr}, "")
	return regexp.MustCompile(regexStr).MatchString(sourceStr)
}
