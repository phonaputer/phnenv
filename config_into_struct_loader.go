package phnenv

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
)

const (
	phnEnvStructTag = "phnenv"

	errWrapFmt   = "phnenv: %w"
	fieldWrapFmt = `field "%s": %w`
)

var (
	errMustBeStructPtr = errors.New("input must be a pointer to a struct")
	errNumericOverflow = errors.New("environment value overflows numeric type")
	errCantSet         = errors.New("can't set field")
	errUnsupportedType = errors.New("unsupported field type")
)

// Function for getting a string value for a string key from a config source
// The bool result should be true if the input key exists in the source.
// If it does not exist the bool result will be false.
type confGetter func(string) (string, bool)

// Parse reads OS environment variables and fills the struct in the value pointed to by v.
// If v is nil or not a pointer to a struct, Parse returns an error.
//
// Parse examines the tags on the fields of the struct pointed to by v in order to
// determine which environment variables should be read for which field, and how
// the string environment variables should be parsed into the type of each field.
//
// The struct field types supported by Parse are:
//   string
//   bool
//   int, int8, int16, int32/rune, int64
//   uint, uint8, uint16, uint32, uint64
//   float32, float64
//   complex64, complex128
// In addition, pointers to and slices of the above types are supported.
// Nested structs are supported.
//
// Types which are not supported are:
//   arrays
//   slices of structs
//   interfaces
//   json, xml, yaml, etc. (although JSON may be added in the future)
//
// The struct tags used by phnenv must have the key "phnenv" and include at least
// an evironment variable name:
//
//   s := struct {
//       FieldName string `phnenv:"ENVIRONMENT_VAR"`
//   }{}
//
//   err := phnenv.Parse(&s)
//
// In the above example, the ENVIRONMENT_VAR variable will be parsed as a string
// and the result set into s.FieldName.
//
// Additional options for how to parse the environment variable can be added
// to a field's tag if necessary.
//
// The parsing options supported by this version of phnenv are:
//
//   rune
//   bitsize:
//   base:
//   sep:
//
// The `rune` parsing option can be applied to fields of type int32
// (the standard Go rune type is an alias for int32, so you can also use the type rune).
// This option tells the parser to parse the environment variable as a single character/rune
// instead of as a base 10 number (the default behavior for int32).
// For example, the following code will print "字" instead of returning a string to integer
// parsing error.
//
//   // In the OS: `ENV_VAR=字`
//
//   s := struct {
//      Field rune `phnenv:"ENV_VAR,rune"`
//   }{}
//
//   phnenv.Parse(&s)
//   fmt.Println(s.Field)
//
// The `bitsize:` option specifies the bit width of the numeric type that the result must fit into.
// This option can be applied to fields of int, uint, float, and complex.
// This is the same as the bitsize parameter to the standard library strconv.ParseInt
// (or parse float, complex, uint) function. So please check the documentation in strconv
// for more information. Example usage where the parsing result must fit into 8 bits
// (Note that a similar check could be achieved by making Field's type uint8):
//
//   s := struct {
//      Field int64 `phnenv:"ENV_VAR,bitsize:8"`
//   }{}
//
// The `base:` option specifies the mathematical base that the string will be parsed as.
// This option can be applied to fields of int and uint.
// It functions the same as the `base` parameter to strconv.ParseInt and strconv.ParseUint.
// Please check the documentation of those functions in strconv for more information.
// Example usage where the string environment variable ENV_VAR should be parsed to int as
// base 2 (binary):
//
//   s struct {
//      Field int `phnenv:"ENV_VAR,base:2"`
//   }{}
//
// The `sep:` option specifies the string used to split a single environment variable into a
// list of strings when parsing a slice type. This works the same as the sep argument to the
// standard library strings.Split function. Please check the strings docs for more info.
// The default separator is ",".
// The following is an example where the environment variable will be split on "||".
// The resulting value in s.Field will be ["abc", "123"]:
//
//   // In the OS: `ENV_VAR=abc||123`
//
//   s struct {
//      Field []string `phnenv:"ENV_VAR,sep:||"`
//   }{}
//
//   phnenv.Parse(&s)
//
// Brief overview of how parsing works for each type:
//
//   string: copied directly from the environment variable
//   int: parsed using strconv.ParseInt
//   uint: parsed using strconv.ParseUint
//   float: parsed using strconv.ParseFloat
//   complex: parsed using strconv.ParseComplex
//   bool: if the environment variable's string equals (ignoring case) "true" then the bool
//      will be true
//   slices: the environment variable's string will be split with strings.Split using a configurable
//      separator. Then, each index will be parsed individually as the slice element type.
//
// Errors will be returned by Parse in the following cases:
//    1. Parsing one or more field fails for any reason.
//    2. One or more struct tags is malformed or invalid.
//    3. The input v is not a pointer to a struct.
//    4. A phnenv struct tag was placed on a struct field of an unsupported type.
func Parse(v interface{}) error {
	err := parse(os.LookupEnv, v)
	if err != nil {
		return fmt.Errorf(errWrapFmt, err)
	}

	return nil
}

func parse(c confGetter, v interface{}) error {
	sv, err := validateInput(v)
	if err != nil {
		return err
	}

	return iterateStruct(c, sv)
}

func validateInput(v interface{}) (reflect.Value, error) {
	var res reflect.Value

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return res, errMustBeStructPtr
	}

	res = rv.Elem()
	if res.Kind() != reflect.Struct {
		return res, errMustBeStructPtr
	}

	return res, nil
}

func iterateStruct(c confGetter, sv reflect.Value) error {
	for i := 0; i < sv.NumField(); i++ {
		err := loadConfAndSetField(c, sv.Type().Field(i), sv.Field(i))
		if err != nil {
			return fmt.Errorf(fieldWrapFmt, sv.Type().Field(i).Name, err)
		}
	}

	return nil
}

func loadConfAndSetField(c confGetter, sf reflect.StructField, fv reflect.Value) error {
	if fv.Kind() == reflect.Struct {
		return iterateStruct(c, fv)
	}

	conf, to, ok, err := parseStructTagAndLoadConf(c, sf)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	err = setField(conf, to, fv)
	if err != nil {
		return err
	}

	return nil
}

func parseStructTagAndLoadConf(c confGetter, sf reflect.StructField) (string, tagOpts, bool, error) {
	tagStr, ok := sf.Tag.Lookup(phnEnvStructTag)
	if !ok {
		return "", tagOpts{}, false, nil // if there's no phnenv tag this is not an error, but we should skip this field
	}

	key, opts, err := parseTag(tagStr)
	if err != nil {
		return "", opts, false, err
	}

	conf, ok := c(key)

	return conf, opts, ok, nil
}

func setField(conf string, to tagOpts, fieldVal reflect.Value) error {
	if !fieldVal.CanSet() {
		return errCantSet
	}

	switch fieldVal.Kind() {
	case reflect.Bool:
		setBasicBool(conf, fieldVal)
	case reflect.String:
		setBasicStr(conf, fieldVal)
	case reflect.Int32:
		return setBasicInt32(conf, to, fieldVal)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int64:
		return setBasicInt(conf, to, fieldVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return setBasicUint(conf, to, fieldVal)
	case reflect.Float32, reflect.Float64:
		return setBasicFloat(conf, to, fieldVal)
	case reflect.Complex64, reflect.Complex128:
		return setBasicComplex(conf, to, fieldVal)
	case reflect.Ptr:
		return setPtr(conf, to, fieldVal)
	case reflect.Slice:
		return setSlice(conf, to, fieldVal)
	default:
		return errUnsupportedType
	}

	return nil
}

func setBasicStr(conf string, fieldVal reflect.Value) {
	fieldVal.SetString(conf)
}

func setBasicBool(conf string, fieldVal reflect.Value) {
	fieldVal.SetBool(strToBool(conf))
}

func setBasicInt(conf string, to tagOpts, fieldVal reflect.Value) error {
	v, err := strToInt(conf, to.NumBitSize, to.NumBase)
	if err != nil {
		return err
	}

	if fieldVal.OverflowInt(v) {
		return errNumericOverflow
	}

	fieldVal.SetInt(v)

	return nil
}

func setBasicInt32(conf string, to tagOpts, fieldVal reflect.Value) error {
	if !to.IsRune {
		return setBasicInt(conf, to, fieldVal)
	}

	v, err := strToIntRune(conf)
	if err != nil {
		return err
	}

	if fieldVal.OverflowInt(v) {
		return errNumericOverflow
	}

	fieldVal.SetInt(v)

	return nil
}

func setBasicUint(conf string, to tagOpts, fieldVal reflect.Value) error {
	v, err := strToUint(conf, to.NumBitSize, to.NumBase)
	if err != nil {
		return err
	}

	if fieldVal.OverflowUint(v) {
		return errNumericOverflow
	}

	fieldVal.SetUint(v)

	return nil
}

func setBasicFloat(conf string, to tagOpts, fieldVal reflect.Value) error {
	v, err := strToFloat(conf, to.NumBitSize)
	if err != nil {
		return err
	}

	if fieldVal.OverflowFloat(v) {
		return errNumericOverflow
	}

	fieldVal.SetFloat(v)

	return nil
}

func setBasicComplex(conf string, to tagOpts, fieldVal reflect.Value) error {
	v, err := strToComplex(conf, to.NumBitSize)
	if err != nil {
		return err
	}

	if fieldVal.OverflowComplex(v) {
		return errNumericOverflow
	}

	fieldVal.SetComplex(v)

	return nil
}

func setSlice(conf string, to tagOpts, fv reflect.Value) error {
	var splt []string
	if len(conf) > 0 {
		splt = strings.Split(conf, to.SliceSep)
	}

	res := reflect.MakeSlice(fv.Type(), len(splt), len(splt))

	for i := 0; i < len(splt); i++ {
		err := setField(splt[i], to, res.Index(i))
		if err != nil {
			return err
		}
	}

	fv.Set(res)

	return nil
}

func setPtr(conf string, to tagOpts, fieldVal reflect.Value) error {
	newPtr := reflect.New(fieldVal.Type().Elem())

	err := setField(conf, to, reflect.Indirect(newPtr))
	if err != nil {
		return err
	}

	fieldVal.Set(newPtr)

	return nil
}
