package phnenv

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_parseBitSize_notBitSizeTag_returnsFalse(t *testing.T) {
	_, ok, err := parseBitSize("what")

	assert.Nil(t, err)
	assert.False(t, ok)
}

func Test_parseSep_notSepTag_returnsFalse(t *testing.T) {
	_, ok, err := parseSep("what")

	assert.Nil(t, err)
	assert.False(t, ok)
}
