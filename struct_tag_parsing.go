package phnenv

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	tagRune       = "rune"
	tagNumBase    = "base:"
	tagNumBitSize = "bitsize:"
	tagSeparator  = ","

	errTagBaseWrapFmt    = "base option: %w"
	errTagBitSizeWrapFmt = "base option: %w"
)

var (
	errTagMissingData      = errors.New("phnenv struct tags must contain at minimum an environment variable name")
	errTagDuplicateRune    = errors.New("struct tag rune option must only be provided once")
	errTagDuplicateBitSize = errors.New("struct tag bitsize option must only be provided once")
	errTagDuplicateBase    = errors.New("struct tag base option must only be provided once")
	errTagUnsupported      = errors.New("unsupported struct tag option provided")
)

type tagOpts struct {
	NumBase    *int
	NumBitSize *int
	IsRune     bool
}

// parseTag parses a phnenv struct tag to get:
// 1. the config key to retrieve to populate this struct field (the string result of parseTag)
// 2. options for parsing the config (the tagOpts struct result of parseTag)
func parseTag(t string) (string, tagOpts, error) {
	opts := tagOpts{IsRune: false}

	key, strOpts, err := validateTag(t)
	if err != nil {
		return "", opts, err
	}

	for _, opt := range strOpts {
		o, err := setOpt(opts, opt)
		if err != nil {
			return "", opts, err
		}

		opts = o
	}

	return key, opts, nil
}

// validateTag checks that no unknown tag options are provided. And that no tag options are provided more than once.
func validateTag(t string) (string, []string, error) {
	if len(t) < 1 {
		return "", nil, errTagMissingData
	}

	splitT := strings.Split(t, tagSeparator)

	if len(splitT[0]) < 1 {
		return "", nil, errTagMissingData
	}

	splitTWithoutKey := splitT[1:]

	foundBase := false
	foundRune := false
	foundBitSize := false
	for _, item := range splitTWithoutKey {
		if item == tagRune {
			if foundRune == true {
				return "", nil, errTagDuplicateRune
			}
			foundRune = true
		} else if len(item) >= len(tagNumBase) && item[:len(tagNumBase)] == tagNumBase {
			if foundBase == true {
				return "", nil, errTagDuplicateBase
			}
			foundBase = true
		} else if len(item) >= len(tagNumBase) && item[:len(tagNumBitSize)] == tagNumBitSize {
			if foundBitSize == true {
				return "", nil, errTagDuplicateBitSize
			}
			foundBitSize = true
		} else {
			return "", nil, errTagUnsupported
		}
	}

	return splitT[0], splitTWithoutKey, nil
}

func setOpt(to tagOpts, opt string) (tagOpts, error) {
	if isRune(opt) {
		to.IsRune = true
		return to, nil
	}

	base, ok, err := parseBase(opt)
	if err != nil {
		return to, fmt.Errorf(errTagBaseWrapFmt, err)
	}
	if ok {
		to.NumBase = &base
		return to, nil
	}

	bitSize, ok, err := parseBitSize(opt)
	if err != nil {
		return to, fmt.Errorf(errTagBitSizeWrapFmt, err)
	}
	if ok {
		to.NumBitSize = &bitSize
		return to, nil
	}

	return to, nil
}

func parseBase(s string) (int, bool, error) {
	if len(s) < len(tagNumBase) || s[:len(tagNumBase)] != tagNumBase {
		return 0, false, nil
	}

	numStr := s[len(tagNumBase):]

	base, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, false, err
	}

	return base, true, nil
}

func parseBitSize(s string) (int, bool, error) {
	if len(s) < len(tagNumBitSize) || s[:len(tagNumBitSize)] != tagNumBitSize {
		return 0, false, nil
	}

	numStr := s[len(tagNumBitSize):]

	bitSize, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, false, err
	}

	return bitSize, true, nil
}

func isRune(s string) bool {
	return s == tagRune
}
