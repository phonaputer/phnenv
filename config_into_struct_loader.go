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

	listSeparator = ","

	errWrapFmt = "phnenv: %w"
	fieldWrapFmt = `field "%s": %w`
)

var (
	errMustBeStructPtr = errors.New("input must be a pointer to a struct")
	errNumericOverflow = errors.New("environment value overflows numeric type")
	errCantSet = errors.New("can't set field")
	errUnsupportedType = errors.New("unsupported field type")
)

// Function for getting a string value for a string key from a config source
// The bool result should be true if the input key exists in the source.
// If it does not exist the bool result will be false.
type confGetter func(string)(string, bool)

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

	conf, ok := loadConfStr(c, sf)
	if !ok {
		return nil
	}

	err := setField(conf, fv)
	if err != nil {
		return err
	}

	return nil
}

func loadConfStr(c confGetter, sf reflect.StructField) (string, bool) {
	confKey, ok := sf.Tag.Lookup(phnEnvStructTag)
	if !ok {
		return "", false
	}

	return c(confKey)
}

func setField(conf string, fieldVal reflect.Value) error {
	if !fieldVal.CanSet() {
		return errCantSet
	}

	switch fieldVal.Kind() {
	case reflect.Bool:
		setBasicBool(conf, fieldVal)
	case reflect.String:
		setBasicStr(conf, fieldVal)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return setBasicInt(conf, fieldVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return setBasicUint(conf, fieldVal)
	case reflect.Float32, reflect.Float64:
		return setBasicFloat(conf, fieldVal)
	case reflect.Complex64, reflect.Complex128:
		return setBasicComplex(conf, fieldVal)
	case reflect.Ptr:
		return setPtr(conf, fieldVal)
	case reflect.Slice:
		return setSlice(conf, fieldVal)
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

func setBasicInt(conf string, fieldVal reflect.Value) error {
	v, err := strToInt(conf)
	if err != nil {
		return err
	}

	if fieldVal.OverflowInt(v) {
		return errNumericOverflow
	}

	fieldVal.SetInt(v)

	return nil
}

func setBasicUint(conf string, fieldVal reflect.Value) error {
	v, err := strToUint(conf)
	if err != nil {
		return err
	}

	if fieldVal.OverflowUint(v) {
		return errNumericOverflow
	}

	fieldVal.SetUint(v)

	return nil
}

func setBasicFloat(conf string, fieldVal reflect.Value) error {
	v, err := strToFloat(conf)
	if err != nil {
		return err
	}

	if fieldVal.OverflowFloat(v) {
		return errNumericOverflow
	}

	fieldVal.SetFloat(v)

	return nil
}

func setBasicComplex(conf string, fieldVal reflect.Value) error {
	v, err := strToComplex(conf)
	if err != nil {
		return err
	}

	if fieldVal.OverflowComplex(v) {
		return errNumericOverflow
	}

	fieldVal.SetComplex(v)

	return nil
}

func setSlice(conf string, fv reflect.Value) error {
	var splt []string
	if len(conf) > 0 {
		splt = strings.Split(conf, listSeparator)
	}

	res := reflect.MakeSlice(fv.Type(), len(splt), len(splt))

	for i := 0; i < len(splt); i++ {
		err := setField(splt[i], res.Index(i))
		if err != nil {
			return err
		}
	}

	fv.Set(res)

	return nil
}

func setPtr(conf string, fieldVal reflect.Value) error {
	newPtr := reflect.New(fieldVal.Type().Elem())

	err := setField(conf, reflect.Indirect(newPtr))
	if err != nil {
		return err
	}

	fieldVal.Set(newPtr)

	return nil
}