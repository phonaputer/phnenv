package phnenv

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Parse_MultipleFieldsInOneStruct_MapsAllFields(t *testing.T) {
	s := struct {
		A int        `phnenv:"A"`
		B uint       `phnenv:"B"`
		C rune       `phnenv:"C,rune"`
		D float64    `phnenv:"D"`
		E complex128 `phnenv:"E"`
		F bool       `phnenv:"F"`
		G struct {
			A string   `phnenv:"GA"`
			B int      `phnenv:"GB"`
			C []string `phnenv:"GC"`
		}
		H string `phnenv:"H"`
		I string `phnenv:"notfound"`
		J string
	}{}

	g := func(s string) (string, bool) {
		switch s {
		case "A":
			return "123", true
		case "B":
			return "234", true
		case "C":
			return "X", true
		case "D":
			return "12.7", true
		case "E":
			return "22.8", true
		case "F":
			return "true", true
		case "GA":
			return "what", true
		case "GB":
			return "678", true
		case "GC":
			return "a,b,c", true
		case "H":
			return "who", true
		}

		return "", false
	}

	err := parse(g, &s)

	assert.Nil(t, err)
	assert.Equal(t, 123, s.A)
	assert.Equal(t, uint(234), s.B)
	assert.Equal(t, 'X', s.C)
	assert.Equal(t, 12.7, s.D)
	assert.Equal(t, complex128(22.8), s.E)
	assert.Equal(t, true, s.F)
	assert.Equal(t, "what", s.G.A)
	assert.Equal(t, 678, s.G.B)
	assert.Equal(t, []string{"a", "b", "c"}, s.G.C)
	assert.Equal(t, "who", s.H)
	var zeroStr string
	assert.Equal(t, zeroStr, s.I)
	assert.Equal(t, zeroStr, s.J)
}

var test_parse_TagOpts_PositiveCases = []struct {
	Name     string
	Input    interface{}
	Conf     string
	Expected interface{}
}{
	{"int base",
		&struct {
			F int `phnenv:"TESTENV,base:2"`
		}{},
		"10",
		&struct {
			F int `phnenv:"TESTENV,base:2"`
		}{F: 2}},
	{"uint base",
		&struct {
			F uint `phnenv:"TESTENV,base:2"`
		}{},
		"10",
		&struct {
			F uint `phnenv:"TESTENV,base:2"`
		}{F: 2}},
	{"int bitsize",
		&struct {
			F int `phnenv:"TESTENV,bitsize:8"`
		}{},
		"125",
		&struct {
			F int `phnenv:"TESTENV,bitsize:8"`
		}{F: 125}},
	{"uint bitsize",
		&struct {
			F uint `phnenv:"TESTENV,bitsize:8"`
		}{},
		"125",
		&struct {
			F uint `phnenv:"TESTENV,bitsize:8"`
		}{F: 125}},
	{"float bitsize",
		&struct {
			F float64 `phnenv:"TESTENV,bitsize:8"`
		}{},
		"125",
		&struct {
			F float64 `phnenv:"TESTENV,bitsize:8"`
		}{F: 125}},
	{"complex bitsize",
		&struct {
			F complex128 `phnenv:"TESTENV,bitsize:8"`
		}{},
		"125",
		&struct {
			F complex128 `phnenv:"TESTENV,bitsize:8"`
		}{F: 125}},
	{"int base & bitsize",
		&struct {
			F int `phnenv:"TESTENV,bitsize:8,base:2"`
		}{},
		"10",
		&struct {
			F int `phnenv:"TESTENV,bitsize:8,base:2"`
		}{F: 2}},
	{"uint base & bitsize",
		&struct {
			F uint `phnenv:"TESTENV,base:2,bitsize:8"`
		}{},
		"10",
		&struct {
			F uint `phnenv:"TESTENV,base:2,bitsize:8"`
		}{F: 2}},
	{"int32 rune",
		&struct {
			F int32 `phnenv:"TESTENV,rune"`
		}{},
		"B",
		&struct {
			F int32 `phnenv:"TESTENV,rune"`
		}{F: 'B'}},
	{"slice sep",
		&struct {
			F []string `phnenv:"TESTENV,sep:||"`
		}{},
		"abc||123",
		&struct {
			F []string `phnenv:"TESTENV,sep:||"`
		}{F: []string{"abc", "123"}}},
	{"int8 slice sep, base & bitsize",
		&struct {
			F []int8 `phnenv:"TESTENV,sep:||,base:2,bitsize:8"`
		}{},
		"10||101",
		&struct {
			F []int8 `phnenv:"TESTENV,sep:||,base:2,bitsize:8"`
		}{F: []int8{2, 5}}},
	{"ignore irrelevant tags - float",
		&struct {
			F float64 `phnenv:"TESTENV,rune,base:2"`
		}{},
		"123",
		&struct {
			F float64 `phnenv:"TESTENV,rune,base:2"`
		}{F: 123}},
	{"ignore irrelevant tags - string",
		&struct {
			F string `phnenv:"TESTENV,rune,base:2,bitsize:32"`
		}{},
		"123",
		&struct {
			F string `phnenv:"TESTENV,rune,base:2,bitsize:32"`
		}{F: "123"}},
}

func Test_parse_TagOpts_PositiveCases(t *testing.T) {
	for _, c := range test_parse_TagOpts_PositiveCases {
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

var test_parse_TagOpts_NegativeCases = []struct {
	Name            string
	Input           interface{}
	Conf            string
	ExpectedErrPart string
}{
	{"empty tag",
		&struct {
			F int `phnenv:""`
		}{},
		"10",
		"must contain at minimum an environment variable"},
	{"tag missing conf key",
		&struct {
			F int `phnenv:",base:2"`
		}{},
		"10",
		"must contain at minimum an environment variable"},
	{"unknown option",
		&struct {
			F int `phnenv:"E,what"`
		}{},
		"10",
		"unsupported struct tag option"},
	{"duplicate rune",
		&struct {
			F int `phnenv:"E,rune,rune"`
		}{},
		"10",
		"rune option must only be provided once"},
	{"duplicate base",
		&struct {
			F int `phnenv:"E,base:1,base:1"`
		}{},
		"10",
		"base option must only be provided once"},
	{"duplicate bitsize",
		&struct {
			F int `phnenv:"E,bitsize:8,bitsize:8"`
		}{},
		"10",
		"bitsize option must only be provided once"},
	{"duplicate sep",
		&struct {
			F int `phnenv:"E,sep:8,sep:8"`
		}{},
		"10",
		"sep option must only be provided once"},
	{"empty bitsize",
		&struct {
			F int `phnenv:"E,bitsize:"`
		}{},
		"10",
		""},
	{"empty base",
		&struct {
			F int `phnenv:"E,base:"`
		}{},
		"10",
		""},
	{"invalid bitsize",
		&struct {
			F int `phnenv:"E,bitsize:abc"`
		}{},
		"10",
		""},
	{"invalid base",
		&struct {
			F int `phnenv:"E,base:abc"`
		}{},
		"10",
		""},
	{"empty sep",
		&struct {
			F int `phnenv:"E,sep:"`
		}{},
		"10",
		"separator must not be empty string"},
}

func Test_parse_TagOpts_NegativeCases(t *testing.T) {
	for _, c := range test_parse_TagOpts_NegativeCases {
		t.Run(c.Name, func(t *testing.T) {
			g := func(s string) (string, bool) {
				if s == "TESTENV" {
					return c.Conf, true
				}

				return "", false
			}

			err := parse(g, c.Input)

			if assert.NotNil(t, err) {
				assert.Contains(t, err.Error(), c.ExpectedErrPart)
			}
		})
	}
}

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
	Name  string
	Input interface{}
}{
	{"array",
		&struct {
			F [2]string `phnenv:"sdgas"`
		}{}},
	{"struct slice",
		&struct {
			F []struct{} `phnenv:"sdgas"`
		}{}},
	{"nested slice",
		&struct {
			F [][]string `phnenv:"sdgas"`
		}{}},
	{"map",
		&struct {
			F map[string]string `phnenv:"sdgas"`
		}{}},
	{"func",
		&struct {
			F func() `phnenv:"sdgas"`
		}{}},
	{"chan",
		&struct {
			F chan string `phnenv:"sdgas"`
		}{}},
	{"interface",
		&struct {
			F interface{} `phnenv:"sdgas"`
		}{}},
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
	Name     string
	Input    interface{}
	Conf     string
	Expected interface{}
}{
	{"string", &struct{ F string }{"yes"}, "hello", &struct{ F string }{F: "yes"}},
	{"bool", &struct{ F bool }{true}, "false", &struct{ F bool }{F: true}},
	{"int", &struct{ F int }{123}, "234", &struct{ F int }{F: 123}},
	{"int8", &struct{ F int8 }{123}, "234", &struct{ F int8 }{F: 123}},
	{"int16", &struct{ F int16 }{123}, "234", &struct{ F int16 }{F: 123}},
	{"int32", &struct{ F int32 }{123}, "234", &struct{ F int32 }{F: 123}},
	{"rune", &struct{ F rune }{'A'}, "B", &struct{ F rune }{F: 'A'}},
	{"int64", &struct{ F int64 }{123}, "234", &struct{ F int64 }{F: 123}},
	{"uint", &struct{ F uint }{123}, "234", &struct{ F uint }{F: 123}},
	{"uint8", &struct{ F uint8 }{123}, "234", &struct{ F uint8 }{F: 123}},
	{"uint16", &struct{ F uint16 }{123}, "234", &struct{ F uint16 }{F: 123}},
	{"uint32", &struct{ F uint32 }{123}, "234", &struct{ F uint32 }{F: 123}},
	{"uint64", &struct{ F uint64 }{123}, "234", &struct{ F uint64 }{F: 123}},
	{"float32", &struct{ F float32 }{123}, "234", &struct{ F float32 }{F: 123}},
	{"float64", &struct{ F float64 }{123}, "234", &struct{ F float64 }{F: 123}},
	{"complex64", &struct{ F complex64 }{123}, "234", &struct{ F complex64 }{F: 123}},
	{"complex128", &struct{ F complex128 }{123}, "234", &struct{ F complex128 }{F: 123}},
	{"pointer", &struct{ F *string }{&helloStr}, "asdg", &struct{ F *string }{F: &helloStr}},
	{"slice", &struct{ F []string }{[]string{"hello", "goodbye"}}, "nope,2", &struct{ F []string }{F: []string{"hello", "goodbye"}}},
	{"slice pointer", &struct{ F *[]string }{&[]string{"hello", "goodbye"}}, "nope,2", &struct{ F *[]string }{F: &[]string{"hello", "goodbye"}}},
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
	Name     string
	Input    interface{}
	Conf     string
	Expected interface{}
}{
	{"string",
		&struct {
			F string `phnenv:"wrong"`
		}{"yes"},
		"hello",
		&struct {
			F string `phnenv:"wrong"`
		}{F: "yes"}},
	{"bool",
		&struct {
			F bool `phnenv:"wrong"`
		}{true},
		"false",
		&struct {
			F bool `phnenv:"wrong"`
		}{F: true}},
	{"int",
		&struct {
			F int `phnenv:"wrong"`
		}{123},
		"234",
		&struct {
			F int `phnenv:"wrong"`
		}{F: 123}},
	{"int8",
		&struct {
			F int8 `phnenv:"wrong"`
		}{123},
		"234",
		&struct {
			F int8 `phnenv:"wrong"`
		}{F: 123}},
	{"int16",
		&struct {
			F int16 `phnenv:"wrong"`
		}{123},
		"234",
		&struct {
			F int16 `phnenv:"wrong"`
		}{F: 123}},
	{"int32",
		&struct {
			F int32 `phnenv:"wrong"`
		}{123},
		"234",
		&struct {
			F int32 `phnenv:"wrong"`
		}{F: 123}},
	{"rune",
		&struct {
			F rune `phnenv:"wrong,rune"`
		}{'A'},
		"B",
		&struct {
			F rune `phnenv:"wrong,rune"`
		}{F: 'A'}},
	{"int64",
		&struct {
			F int64 `phnenv:"wrong"`
		}{123},
		"234",
		&struct {
			F int64 `phnenv:"wrong"`
		}{F: 123}},
	{"uint",
		&struct {
			F uint `phnenv:"wrong"`
		}{123},
		"234",
		&struct {
			F uint `phnenv:"wrong"`
		}{F: 123}},
	{"uint8",
		&struct {
			F uint8 `phnenv:"wrong"`
		}{123},
		"234",
		&struct {
			F uint8 `phnenv:"wrong"`
		}{F: 123}},
	{"uint16",
		&struct {
			F uint16 `phnenv:"wrong"`
		}{123},
		"234",
		&struct {
			F uint16 `phnenv:"wrong"`
		}{F: 123}},
	{"uint32",
		&struct {
			F uint32 `phnenv:"wrong"`
		}{123},
		"234",
		&struct {
			F uint32 `phnenv:"wrong"`
		}{F: 123}},
	{"uint64",
		&struct {
			F uint64 `phnenv:"wrong"`
		}{123},
		"234",
		&struct {
			F uint64 `phnenv:"wrong"`
		}{F: 123}},
	{"float32",
		&struct {
			F float32 `phnenv:"wrong"`
		}{123},
		"234",
		&struct {
			F float32 `phnenv:"wrong"`
		}{F: 123}},
	{"float64",
		&struct {
			F float64 `phnenv:"wrong"`
		}{123},
		"234",
		&struct {
			F float64 `phnenv:"wrong"`
		}{F: 123}},
	{"complex64",
		&struct {
			F complex64 `phnenv:"wrong"`
		}{123},
		"234",
		&struct {
			F complex64 `phnenv:"wrong"`
		}{F: 123}},
	{"complex128",
		&struct {
			F complex128 `phnenv:"wrong"`
		}{123},
		"234",
		&struct {
			F complex128 `phnenv:"wrong"`
		}{F: 123}},
	{"pointer",
		&struct {
			F *string `phnenv:"wrong"`
		}{&helloStr},
		"asdg",
		&struct {
			F *string `phnenv:"wrong"`
		}{F: &helloStr}},
	{"slice",
		&struct {
			F []string `phnenv:"wrong"`
		}{[]string{"hello", "goodbye"}},
		"nope,2",
		&struct {
			F []string `phnenv:"wrong"`
		}{F: []string{"hello", "goodbye"}}},
	{"slice pointer",
		&struct {
			F *[]string `phnenv:"wrong"`
		}{&[]string{"hello", "goodbye"}},
		"nope,2",
		&struct {
			F *[]string `phnenv:"wrong"`
		}{F: &[]string{"hello", "goodbye"}}},
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
	Name     string
	Input    interface{}
	Conf     string
	Expected interface{}
}{
	{"string",
		&struct {
			F string `phnenv:"TESTENV"`
		}{},
		"hello",
		&struct {
			F string `phnenv:"TESTENV"`
		}{F: "hello"}},
	{"bool",
		&struct {
			F bool `phnenv:"TESTENV"`
		}{},
		"true",
		&struct {
			F bool `phnenv:"TESTENV"`
		}{F: true}},
	{"int",
		&struct {
			F int `phnenv:"TESTENV"`
		}{},
		"123",
		&struct {
			F int `phnenv:"TESTENV"`
		}{F: 123}},
	{"int8",
		&struct {
			F int8 `phnenv:"TESTENV"`
		}{},
		"123",
		&struct {
			F int8 `phnenv:"TESTENV"`
		}{F: 123}},
	{"int16",
		&struct {
			F int16 `phnenv:"TESTENV"`
		}{},
		"123",
		&struct {
			F int16 `phnenv:"TESTENV"`
		}{F: 123}},
	{"int32",
		&struct {
			F int32 `phnenv:"TESTENV"`
		}{},
		"123",
		&struct {
			F int32 `phnenv:"TESTENV"`
		}{F: 123}},
	{"rune",
		&struct {
			F rune `phnenv:"TESTENV,rune"`
		}{},
		"Z",
		&struct {
			F rune `phnenv:"TESTENV,rune"`
		}{F: 'Z'}},
	{"int64",
		&struct {
			F int64 `phnenv:"TESTENV"`
		}{},
		"123",
		&struct {
			F int64 `phnenv:"TESTENV"`
		}{F: 123}},
	{"uint",
		&struct {
			F uint `phnenv:"TESTENV"`
		}{},
		"123",
		&struct {
			F uint `phnenv:"TESTENV"`
		}{F: 123}},
	{"uint8",
		&struct {
			F uint8 `phnenv:"TESTENV"`
		}{},
		"123",
		&struct {
			F uint8 `phnenv:"TESTENV"`
		}{F: 123}},
	{"uint16",
		&struct {
			F uint16 `phnenv:"TESTENV"`
		}{},
		"123",
		&struct {
			F uint16 `phnenv:"TESTENV"`
		}{F: 123}},
	{"uint32",
		&struct {
			F uint32 `phnenv:"TESTENV"`
		}{},
		"123",
		&struct {
			F uint32 `phnenv:"TESTENV"`
		}{F: 123}},
	{"uint64",
		&struct {
			F uint64 `phnenv:"TESTENV"`
		}{},
		"123",
		&struct {
			F uint64 `phnenv:"TESTENV"`
		}{F: 123}},
	{"float32",
		&struct {
			F float32 `phnenv:"TESTENV"`
		}{},
		"123",
		&struct {
			F float32 `phnenv:"TESTENV"`
		}{F: 123}},
	{"float64",
		&struct {
			F float64 `phnenv:"TESTENV"`
		}{},
		"123",
		&struct {
			F float64 `phnenv:"TESTENV"`
		}{F: 123}},
	{"complex64",
		&struct {
			F complex64 `phnenv:"TESTENV"`
		}{},
		"123",
		&struct {
			F complex64 `phnenv:"TESTENV"`
		}{F: 123}},
	{"complex128",
		&struct {
			F complex128 `phnenv:"TESTENV"`
		}{},
		"123",
		&struct {
			F complex128 `phnenv:"TESTENV"`
		}{F: 123}},
	{"pointer",
		&struct {
			F *string `phnenv:"TESTENV"`
		}{},
		"hello",
		&struct {
			F *string `phnenv:"TESTENV"`
		}{F: &helloStr}},
	{"slice",
		&struct {
			F []string `phnenv:"TESTENV"`
		}{},
		"hello,goodbye",
		&struct {
			F []string `phnenv:"TESTENV"`
		}{F: []string{"hello", "goodbye"}}},
	{"slice pointer",
		&struct {
			F *[]string `phnenv:"TESTENV"`
		}{},
		"hello,goodbye",
		&struct {
			F *[]string `phnenv:"TESTENV"`
		}{F: &[]string{"hello", "goodbye"}}},
	{"nested struct",
		&struct {
			F struct {
				F2 string `phnenv:"TESTENV"`
			}
		}{},
		"hello",
		&struct {
			F struct {
				F2 string `phnenv:"TESTENV"`
			}
		}{F: struct {
			F2 string `phnenv:"TESTENV"`
		}{F2: "hello"}}},
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

var test_parse_EnvVarExistsButIsInvalid_ShouldReturnError = []struct {
	Name  string
	Input interface{}
	Conf  string
}{
	{"int",
		&struct {
			F int `phnenv:"TESTENV"`
		}{},
		"abc"},
	{"int8",
		&struct {
			F int8 `phnenv:"TESTENV"`
		}{},
		"abc"},
	{"int8 overflow",
		&struct {
			F int8 `phnenv:"TESTENV"`
		}{},
		"129"},
	{"int16",
		&struct {
			F int16 `phnenv:"TESTENV"`
		}{},
		"abc"},
	{"int16 overflow",
		&struct {
			F int16 `phnenv:"TESTENV"`
		}{},
		"99999"},
	{"int32",
		&struct {
			F int32 `phnenv:"TESTENV"`
		}{},
		"abc"},
	{"int32 rune",
		&struct {
			F rune `phnenv:"TESTENV,rune"`
		}{},
		"abc"},
	{"int32 overflow",
		&struct {
			F int32 `phnenv:"TESTENV"`
		}{},
		"999999999999"},
	{"int64",
		&struct {
			F int64 `phnenv:"TESTENV"`
		}{},
		"abc"},
	{"int64 overflow",
		&struct {
			F int64 `phnenv:"TESTENV"`
		}{},
		"9999999999999999999999999"},
	{"uint",
		&struct {
			F uint `phnenv:"TESTENV"`
		}{},
		"abc"},
	{"uint8",
		&struct {
			F uint8 `phnenv:"TESTENV"`
		}{},
		"abc"},
	{"uint8 overflow",
		&struct {
			F uint8 `phnenv:"TESTENV"`
		}{},
		"257"},
	{"uint16",
		&struct {
			F uint16 `phnenv:"TESTENV"`
		}{},
		"abc"},
	{"uint16 overflow",
		&struct {
			F uint16 `phnenv:"TESTENV"`
		}{},
		"999999"},
	{"uint32",
		&struct {
			F uint32 `phnenv:"TESTENV"`
		}{},
		"abc"},
	{"uint32 overflow",
		&struct {
			F uint32 `phnenv:"TESTENV"`
		}{},
		"999999999999"},
	{"uint64",
		&struct {
			F uint64 `phnenv:"TESTENV"`
		}{},
		"abc"},
	{"uint64 overflow",
		&struct {
			F uint64 `phnenv:"TESTENV"`
		}{},
		"99999999999999999999999"},
	{"float32",
		&struct {
			F float32 `phnenv:"TESTENV"`
		}{},
		"abc"},
	{"float32 overflow",
		&struct {
			F float32 `phnenv:"TESTENV"`
		}{},
		"3.40282346638528859811704183484516925440e+39"},
	{"float64",
		&struct {
			F float64 `phnenv:"TESTENV"`
		}{},
		"abc"},
	{"float64 overflow",
		&struct {
			F float64 `phnenv:"TESTENV"`
		}{},
		"1.797693134862315708145274237317043567981e+309"},
	{"complex64",
		&struct {
			F complex64 `phnenv:"TESTENV"`
		}{},
		"abc"},
	{"complex64 overflow",
		&struct {
			F complex64 `phnenv:"TESTENV"`
		}{},
		"3.40282346638528859811704183484516925440e+39"},
	{"complex128",
		&struct {
			F complex128 `phnenv:"TESTENV"`
		}{},
		"abc"},
	{"complex128 overflow",
		&struct {
			F complex128 `phnenv:"TESTENV"`
		}{},
		"1.797693134862315708145274237317043567981e+309"},
	{"pointer",
		&struct {
			F *int `phnenv:"TESTENV"`
		}{},
		"abc"},
	{"slice",
		&struct {
			F []int `phnenv:"TESTENV"`
		}{},
		"hello,goodbye"},
	{"slice pointer",
		&struct {
			F *[]int `phnenv:"TESTENV"`
		}{},
		"hello,goodbye"},
	{"nested struct",
		&struct {
			F struct {
				F2 int `phnenv:"TESTENV"`
			}
		}{},
		"hello"},
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
