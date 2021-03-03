package phnenv

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_parse_NonStructPointerInput_ErrReturned(t *testing.T) {
	err := Parse(struct{}{})
	assert.True(t, errors.Is(err, errMustBeStructPtr))

	str := "agad"
	err = Parse(&str)
	assert.True(t, errors.Is(err, errMustBeStructPtr))

	err = Parse(nil)
	assert.True(t, errors.Is(err, errMustBeStructPtr))
}

var test_parse_UnsupportedFieldType_ReturnsError = []struct {
	Name string
	Input interface{}
} {
	{"array", &struct{F [2]string `phnenv:"sdgas"`}{}},
	{"struct slice", &struct{F []struct{} `phnenv:"sdgas"`}{}},
	{"map", &struct{F map[string]string `phnenv:"sdgas"`}{}},
	{"func", &struct{F func() `phnenv:"sdgas"`}{}},
	{"chan", &struct{F chan string `phnenv:"sdgas"`}{}},
	{"interface", &struct{F interface{} `phnenv:"sdgas"`}{}},
}
func Test_parse_UnsupportedFieldType_ReturnsError(t *testing.T) {
	for _, c := range test_parse_UnsupportedFieldType_ReturnsError {
		t.Run(c.Name, func(t *testing.T) {
			g := func(s string) (string, bool) {
				return "dgs", true
			}

			err := parse(g, c.Input)

			assert.True(t, errors.Is(err, errUnsupportedType))
		})
	}
}

var test_parse_NoPhnEnvTagInStruct_ShouldKeepOriginalValue = []struct {
	Name string
	Input interface{}
	Conf string
	Expected interface{}
} {
	{"string", &struct{F string}{"yes"}, "hello", &struct{F string}{F: "yes"}},
	{"bool", &struct{F bool}{true}, "false", &struct{F bool}{F: true}},
	{"int", &struct{F int}{123}, "234", &struct{F int}{F: 123}},
	{"int8", &struct{F int8}{123}, "234", &struct{F int8}{F: 123}},
	{"int16", &struct{F int16}{123}, "234", &struct{F int16}{F: 123}},
	{"int32", &struct{F int32}{123}, "234", &struct{F int32}{F: 123}},
	{"int64", &struct{F int64}{123}, "234", &struct{F int64}{F: 123}},
	{"uint", &struct{F uint}{123}, "234", &struct{F uint}{F: 123}},
	{"uint8", &struct{F uint8}{123}, "234", &struct{F uint8}{F: 123}},
	{"uint16", &struct{F uint16}{123}, "234", &struct{F uint16}{F: 123}},
	{"uint32", &struct{F uint32}{123}, "234", &struct{F uint32}{F: 123}},
	{"uint64", &struct{F uint64}{123}, "234", &struct{F uint64}{F: 123}},
	{"float32", &struct{F float32}{123}, "234", &struct{F float32}{F: 123}},
	{"float64", &struct{F float64}{123}, "234", &struct{F float64}{F: 123}},
	{"complex64", &struct{F complex64}{123}, "234", &struct{F complex64}{F: 123}},
	{"complex128", &struct{F complex128}{123}, "234", &struct{F complex128}{F: 123}},
	{"pointer", &struct{F *string}{&helloStr}, "asdg", &struct{F *string}{F: &helloStr}},
	{"slice", &struct{F []string}{[]string{"hello", "goodbye"}}, "nope,2", &struct{F []string}{F: []string{"hello", "goodbye"}}},
	{"slice pointer", &struct{F *[]string}{&[]string{"hello", "goodbye"}}, "nope,2", &struct{F *[]string}{F: &[]string{"hello", "goodbye"}}},
}
func Test_parse_NoPhnEnvTagInStruct_ShouldKeepOriginalValue(t *testing.T) {
	for _, c := range test_parse_NoPhnEnvTagInStruct_ShouldKeepOriginalValue {
		t.Run(c.Name, func(t *testing.T) {
			g := func(s string) (string, bool) {
				return c.Conf, true
			}

			err := parse(g, c.Input)

			if assert.Nil(t, err) {
				assert.Equal(t, c.Input, c.Expected)
			}
		})
	}
}

var test_parse_NoEnvVarFound_ShouldKeepOriginalValue = []struct {
	Name string
	Input interface{}
	Conf string
	Expected interface{}
} {
	{"string", &struct{F string `phnenv:"wrong"`}{"yes"}, "hello", &struct{F string `phnenv:"wrong"`}{F: "yes"}},
	{"bool", &struct{F bool `phnenv:"wrong"`}{true}, "false", &struct{F bool `phnenv:"wrong"`}{F: true}},
	{"int", &struct{F int `phnenv:"wrong"`}{123}, "234", &struct{F int `phnenv:"wrong"`}{F: 123}},
	{"int8", &struct{F int8 `phnenv:"wrong"`}{123}, "234", &struct{F int8 `phnenv:"wrong"`}{F: 123}},
	{"int16", &struct{F int16 `phnenv:"wrong"`}{123}, "234", &struct{F int16 `phnenv:"wrong"`}{F: 123}},
	{"int32", &struct{F int32 `phnenv:"wrong"`}{123}, "234", &struct{F int32 `phnenv:"wrong"`}{F: 123}},
	{"int64", &struct{F int64 `phnenv:"wrong"`}{123}, "234", &struct{F int64 `phnenv:"wrong"`}{F: 123}},
	{"uint", &struct{F uint `phnenv:"wrong"`}{123}, "234", &struct{F uint `phnenv:"wrong"`}{F: 123}},
	{"uint8", &struct{F uint8 `phnenv:"wrong"`}{123}, "234", &struct{F uint8 `phnenv:"wrong"`}{F: 123}},
	{"uint16", &struct{F uint16 `phnenv:"wrong"`}{123}, "234", &struct{F uint16 `phnenv:"wrong"`}{F: 123}},
	{"uint32", &struct{F uint32 `phnenv:"wrong"`}{123}, "234", &struct{F uint32 `phnenv:"wrong"`}{F: 123}},
	{"uint64", &struct{F uint64 `phnenv:"wrong"`}{123}, "234", &struct{F uint64 `phnenv:"wrong"`}{F: 123}},
	{"float32", &struct{F float32 `phnenv:"wrong"`}{123}, "234", &struct{F float32 `phnenv:"wrong"`}{F: 123}},
	{"float64", &struct{F float64 `phnenv:"wrong"`}{123}, "234", &struct{F float64 `phnenv:"wrong"`}{F: 123}},
	{"complex64", &struct{F complex64 `phnenv:"wrong"`}{123}, "234", &struct{F complex64 `phnenv:"wrong"`}{F: 123}},
	{"complex128", &struct{F complex128 `phnenv:"wrong"`}{123}, "234", &struct{F complex128 `phnenv:"wrong"`}{F: 123}},
	{"pointer", &struct{F *string `phnenv:"wrong"`}{&helloStr}, "asdg", &struct{F *string `phnenv:"wrong"`}{F: &helloStr}},
	{"slice", &struct{F []string `phnenv:"wrong"`}{[]string{"hello", "goodbye"}}, "nope,2", &struct{F []string `phnenv:"wrong"`}{F: []string{"hello", "goodbye"}}},
	{"slice pointer", &struct{F *[]string `phnenv:"wrong"`}{&[]string{"hello", "goodbye"}}, "nope,2", &struct{F *[]string `phnenv:"wrong"`}{F: &[]string{"hello", "goodbye"}}},
}
func Test_parse_NoEnvVarFound_ShouldKeepOriginalValue(t *testing.T) {
	for _, c := range test_parse_NoPhnEnvTagInStruct_ShouldKeepOriginalValue {
		t.Run(c.Name, func(t *testing.T) {
			g := func(s string) (string, bool) {
				if s == "TESTENV" {
					return c.Conf, true
				}

				return "", false
			}

			err := parse(g, c.Input)

			if assert.Nil(t, err) {
				assert.Equal(t, c.Input, c.Expected)
			}
		})
	}
}

var helloStr = "hello"
var test_parse_EnvVarExistsAndIsValid_ShouldSetInStructField = []struct {
	Name string
	Input interface{}
	Conf string
	Expected interface{}
} {
	{"string", &struct{F string `phnenv:"TESTENV"`}{}, "hello", &struct{F string `phnenv:"TESTENV"`}{F: "hello"}},
	{"bool", &struct{F bool `phnenv:"TESTENV"`}{}, "true", &struct{F bool `phnenv:"TESTENV"`}{F: true}},
	{"int", &struct{F int `phnenv:"TESTENV"`}{}, "123", &struct{F int `phnenv:"TESTENV"`}{F: 123}},
	{"int8", &struct{F int8 `phnenv:"TESTENV"`}{}, "123", &struct{F int8 `phnenv:"TESTENV"`}{F: 123}},
	{"int16", &struct{F int16 `phnenv:"TESTENV"`}{}, "123", &struct{F int16 `phnenv:"TESTENV"`}{F: 123}},
	{"int32", &struct{F int32 `phnenv:"TESTENV"`}{}, "123", &struct{F int32 `phnenv:"TESTENV"`}{F: 123}},
	{"int64", &struct{F int64 `phnenv:"TESTENV"`}{}, "123", &struct{F int64 `phnenv:"TESTENV"`}{F: 123}},
	{"uint", &struct{F uint `phnenv:"TESTENV"`}{}, "123", &struct{F uint `phnenv:"TESTENV"`}{F: 123}},
	{"uint8", &struct{F uint8 `phnenv:"TESTENV"`}{}, "123", &struct{F uint8 `phnenv:"TESTENV"`}{F: 123}},
	{"uint16", &struct{F uint16 `phnenv:"TESTENV"`}{}, "123", &struct{F uint16 `phnenv:"TESTENV"`}{F: 123}},
	{"uint32", &struct{F uint32 `phnenv:"TESTENV"`}{}, "123", &struct{F uint32 `phnenv:"TESTENV"`}{F: 123}},
	{"uint64", &struct{F uint64 `phnenv:"TESTENV"`}{}, "123", &struct{F uint64 `phnenv:"TESTENV"`}{F: 123}},
	{"float32", &struct{F float32 `phnenv:"TESTENV"`}{}, "123", &struct{F float32 `phnenv:"TESTENV"`}{F: 123}},
	{"float64", &struct{F float64 `phnenv:"TESTENV"`}{}, "123", &struct{F float64 `phnenv:"TESTENV"`}{F: 123}},
	{"complex64", &struct{F complex64 `phnenv:"TESTENV"`}{}, "123", &struct{F complex64 `phnenv:"TESTENV"`}{F: 123}},
	{"complex128", &struct{F complex128 `phnenv:"TESTENV"`}{}, "123", &struct{F complex128 `phnenv:"TESTENV"`}{F: 123}},
	{"pointer", &struct{F *string `phnenv:"TESTENV"`}{}, "hello", &struct{F *string `phnenv:"TESTENV"`}{F: &helloStr}},
	{"slice", &struct{F []string `phnenv:"TESTENV"`}{}, "hello,goodbye", &struct{F []string `phnenv:"TESTENV"`}{F: []string{"hello", "goodbye"}}},
	{"slice pointer", &struct{F *[]string `phnenv:"TESTENV"`}{}, "hello,goodbye", &struct{F *[]string `phnenv:"TESTENV"`}{F: &[]string{"hello", "goodbye"}}},
	{"nested struct", &struct{F struct{F2 string `phnenv:"TESTENV"`}}{}, "hello", &struct{F struct{F2 string `phnenv:"TESTENV"`}}{F: struct{F2 string  `phnenv:"TESTENV"`}{F2: "hello"}}},
}
func Test_parse_EnvVarExistsAndIsValid_ShouldSetInStructField(t *testing.T) {
	for _, c := range test_parse_EnvVarExistsAndIsValid_ShouldSetInStructField {
		t.Run(c.Name, func(t *testing.T) {
			g := func(s string) (string, bool) {
				if s == "TESTENV" {
					return c.Conf, true
				}

				return "", false
			}

			err := parse(g, c.Input)

			if assert.Nil(t, err) {
				assert.Equal(t, c.Input, c.Expected)
			}
		})
	}
}

var test_parse_EnvVarExistsButIsInvalid_ShouldReturnError = []struct{
	Name string
	Input interface{}
	Conf string
} {
	{"int", &struct{F int `phnenv:"TESTENV"`}{}, "abc"},
	{"int8", &struct{F int8 `phnenv:"TESTENV"`}{}, "abc"},
	{"int8 overflow", &struct{F int8 `phnenv:"TESTENV"`}{}, "129"},
	{"int16", &struct{F int16 `phnenv:"TESTENV"`}{}, "abc"},
	{"int16 overflow", &struct{F int16 `phnenv:"TESTENV"`}{}, "99999"},
	{"int32", &struct{F int32 `phnenv:"TESTENV"`}{}, "abc"},
	{"int32 overflow", &struct{F int32 `phnenv:"TESTENV"`}{}, "999999999999"},
	{"int64", &struct{F int64 `phnenv:"TESTENV"`}{}, "abc"},
	{"int64 overflow", &struct{F int64 `phnenv:"TESTENV"`}{}, "9999999999999999999999999"},
	{"uint", &struct{F uint `phnenv:"TESTENV"`}{}, "abc"},
	{"uint8", &struct{F uint8 `phnenv:"TESTENV"`}{}, "abc"},
	{"uint8 overflow", &struct{F uint8 `phnenv:"TESTENV"`}{}, "257"},
	{"uint16", &struct{F uint16 `phnenv:"TESTENV"`}{}, "abc"},
	{"uint16 overflow", &struct{F uint16 `phnenv:"TESTENV"`}{}, "999999"},
	{"uint32", &struct{F uint32 `phnenv:"TESTENV"`}{}, "abc"},
	{"uint32 overflow", &struct{F uint32 `phnenv:"TESTENV"`}{}, "999999999999"},
	{"uint64", &struct{F uint64 `phnenv:"TESTENV"`}{}, "abc"},
	{"uint64 overflow", &struct{F uint64 `phnenv:"TESTENV"`}{}, "99999999999999999999999"},
	{"float32", &struct{F float32 `phnenv:"TESTENV"`}{}, "abc"},
	{"float32 overflow", &struct{F float32 `phnenv:"TESTENV"`}{}, "3.40282346638528859811704183484516925440e+39"},
	{"float64", &struct{F float64 `phnenv:"TESTENV"`}{}, "abc"},
	{"float64 overflow", &struct{F float64 `phnenv:"TESTENV"`}{}, "1.797693134862315708145274237317043567981e+309"},
	{"complex64", &struct{F complex64 `phnenv:"TESTENV"`}{}, "abc"},
	{"complex64 overflow", &struct{F complex64 `phnenv:"TESTENV"`}{}, "3.40282346638528859811704183484516925440e+39"},
	{"complex128", &struct{F complex128 `phnenv:"TESTENV"`}{}, "abc"},
	{"complex128 overflow", &struct{F complex128 `phnenv:"TESTENV"`}{}, "1.797693134862315708145274237317043567981e+309"},
	{"pointer", &struct{F *int `phnenv:"TESTENV"`}{}, "abc"},
	{"slice", &struct{F []int `phnenv:"TESTENV"`}{}, "hello,goodbye"},
	{"slice pointer", &struct{F *[]int `phnenv:"TESTENV"`}{}, "hello,goodbye"},
	{"nested struct", &struct{F struct{F2 int `phnenv:"TESTENV"`}}{}, "hello"},

}
func Test_parse_EnvVarExistsButIsInvalid_ShouldReturnError(t *testing.T) {
	for _, c := range test_parse_EnvVarExistsButIsInvalid_ShouldReturnError {
		t.Run(c.Name, func(t *testing.T) {
			g := func(s string) (string, bool) {
				if s == "TESTENV" {
					return c.Conf, true
				}

				return "", false
			}

			err := parse(g, c.Input)

			assert.NotNil(t, err)
		})
	}
}