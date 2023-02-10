package common

import (
	"regexp"
	"strings"
)

func GetKey(input string) string {
	return strings.Split(input, ":")[0]
}

func GetValue(input string) string {
	return strings.Split(input, ":")[1]
}

func GetIdentifierFromDesc(input string) string {
	r, _ := regexp.Compile("host_name=\"(.+)\"")
	return strings.Split(r.FindString(input), "=")[1]
}
