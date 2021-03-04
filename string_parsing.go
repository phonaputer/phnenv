package phnenv

import (
	"errors"
	"strconv"
	"strings"
)

var (
	errRuneLength = errors.New("less/more than 1 rune found for rune type")
)

func strToInt(s string, bitsize *int, base *int) (int64, error) {
	bsz := 64
	if bitsize != nil {
		bsz = *bitsize
	}

	bse := 10
	if base != nil {
		bse = *base
	}

	return strconv.ParseInt(s, bse, bsz)
}

func strToIntRune(s string) (int64, error) {
	rns := []rune(s)

	if len(rns) != 1 {
		return 0, errRuneLength
	}

	return int64(rns[0]), nil
}

func strToBool(s string) bool {
	return strings.ToLower(s) == "true"
}

func strToFloat(s string, bitsize *int) (float64, error) {
	bs := 64
	if bitsize != nil {
		bs = *bitsize
	}

	return strconv.ParseFloat(s, bs)
}

func strToUint(s string, bitsize *int, base *int) (uint64, error) {
	bsz := 64
	if bitsize != nil {
		bsz = *bitsize
	}

	bse := 10
	if base != nil {
		bse = *base
	}

	return strconv.ParseUint(s, bse, bsz)
}

func strToComplex(s string, bitsize *int) (complex128, error) {
	bs := 128
	if bitsize != nil {
		bs = *bitsize
	}

	return strconv.ParseComplex(s, bs)
}
