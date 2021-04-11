package utils

import "strings"

func ConvertToString(s string) string {
	ByteOrderMarkAsString := string('\uFEFF')
	str := strings.TrimPrefix(s, ByteOrderMarkAsString)
	return str
}

func ConvertToLowerCaseString(s string) string {
	str := strings.ToLower(ConvertToString(s))
	return str
}
