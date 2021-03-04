package phnenv

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var test_strToBool_ShouldReturnTrue = []struct {
	Name  string
	Input string
}{
	{"true lower", "true"},
	{"true upper", "TRUE"},
	{"true mixed", "TrUe"},
}

func Test_strToBool_ShouldReturnTrue(t *testing.T) {
	for _, c := range test_strToBool_ShouldReturnTrue {
		t.Run(c.Name, func(t *testing.T) {
			res := strToBool(c.Input)

			assert.True(t, res)
		})
	}
}

var test_strToBool_ShouldReturnFalse = []struct {
	Name  string
	Input string
}{
	{"false lower", "false"},
	{"false upper", "FALSE"},
	{"random string that's not 'true'", "asdgasd"},
	{"t", "t"},
}

func Test_strToBool_ShouldReturnFalse(t *testing.T) {
	for _, c := range test_strToBool_ShouldReturnFalse {
		t.Run(c.Name, func(t *testing.T) {
			res := strToBool(c.Input)

			assert.False(t, res)
		})
	}
}
