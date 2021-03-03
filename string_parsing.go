package phnenv

import (
	"strconv"
	"strings"
)

func strToInt(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func strToBool(s string) bool {
	return strings.ToLower(s) == "true"
}

func strToFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func strToUint(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

func strToComplex(s string) (complex128, error) {
	return strconv.ParseComplex(s, 128)
}