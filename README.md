# PhnENV: Parse OS Environment into a Struct in Go

PhnENV allows users to parse operating system environment variables into a custom Go struct.
The env to struct mapping can be configured using struct tags, similarly to the standard library `encoding/json` package.

To include PhnENV in your project, run the following go get command:

```
go get github.com/phonaputer/phnenv
```

GoDocs for PhnENV can be found here:

```
https://pkg.go.dev/github.com/phonaputer/phnenv
```

## Example Usage

The following is an example of parsing a few different field types from the environment.
For a full list of supported types, see below.

#### The environment being parsed:

```
A_STRING=hello
SOME_STRINGS=how,are,you?
AN_INT=123
A_BOOL=true
A_FRACTION=13.37
OPTIONAL_INT=456
```

#### The Go code:

```
import "github.com/phonaputer/phnenv"

type EnvConfig struct {
   AString string `phnenv:"A_STRING"`
   SomeStrings []string `phnenv:"SOME_STRINGS"`
   Nested struct {
      AnInt int `phnenv:"AN_INT"`
   }
   ABool bool `phenenv:"A_BOOL"`
   AFraction float64 `phnenv:"A_FRACTION"`
   OptionalInt *int `phnenv:"OPTIONAL_INT"`
   NotInEnv *string `phnenv:"NOT_IN_ENV"`
}

func loadEnvIntoStruct(){ 
  var e EnvConfig
  
  _ := phnenv.Parse(&e)
}
```

#### The Result:

After executing the above code, the struct instance `e` will look like this:

```
{
   AString: "hello",
   SomeStrings: []string{"how", "are", "you?"},
   Nested: {
      AnInt: 123
   }
   ABool: true,
   AFraction: 13.37,
   OptionalInt: &456,
   NotInEnv: nil,
}
```

## Supported Field Types

* string
* bool
* int, int8, int16, int32/rune, int64
* uint, uint8, uint16, uint32, uint64
* float32, float64
* complex64, complex128

In addition, pointers to and slices of the above types are supported (including slices of pointers, pointers to slices, etc.).

## Unsupported Field Types

* Arrays
* Slices of structs
* Nested slices
* Interfaces
* JSON, XML, YAML, etc. (although JSON may be added in the future)

## How does parsing work?

For each `phnenv` tagged field in your struct, first, the parser checks to see if a variable exists in the environment with the specified name.
It does this using the standard library `os.LookupEnv` function.
If the variable does NOT exist, the struct field's value is not modified and no error is thrown.
If the variable does exist, it is parsed based on the type of the struct field and the result of parsing is placed into the field.
If the variable cannot be parsed as the field type, parsing stops and an error is returned.

## How Different Types are Parsed

### String

Strings are copied directly from the environment variable into the struct field.

### Bool

The value of the environment variable is considered to be boolean `true` iff the env var equals (ignoring case) the string `"true"`.

### Int

Integer types are parsed using the standard library `strconv.ParseInt` function.

The `base` and `bitsize` parameters of `strconv.ParseInt` can be specified in the struct tag using the `base:` and `bitsize:` options. For example:

```
type Example struct {
    Binary8Bit int `phnenv:"BINARY_8_BIT,base:2,bitsize:8"`
}
```

#### Special case: rune/uint32

`uint32`/`rune` is parsed by default using `strconv.ParseInt`.
However, if you would like to parse a unicode character instead of a string representation of a number, you can specify the option `rune` in your struct tag.

For example:

```
type Example struct {
   ARune rune `phnenv:"A_RUNE,rune"`
}
```

### Uint

Unsigned integer types are parsed using the standard library `strconv.ParseUint` function.
The `base` and `bitsize` parameters of `strconv.ParseUint` can be specified in the struct tag (see **Int** above).

### Float

Float types are parsed using the standard library `strconv.ParseFloat` function.
The `bitsize` parameter of `strconv.ParseFloat` can be specified in the struct tag (see **Int** above).

### Complex

Complex types are parsed using the standard library `strconv.ParseComplex` function.
The `bitsize` parameter of `strconv.ParseComplex` can be specified in the struct tag (see **Int** above).

### Slice

To parse into a slice field, first the environment variable is split using the standard library `strings.Split` function.
Then each element in the split string is handled individually following the same parsing rules as above.

The `sep` (separator) argument to `strings.Split` can be specified in the struct tag using the `sep:` option:

```
type Example struct {
    DoubleBarSeparated []string `phnenv:"DOUBLE_BAR_SEPARATED,sep:||"`
}
```

Note that if you have a slice of `rune` or of numeric types, the `rune`, `bitsize:`, and `base:` options can also be provided. For example:

```
type Example struct {
    BarSeparatedSliceOf32BitBinary []int `phnenv:"SLICE_OF_32_BIT_BINARY,sep:|,base:2,bitsize:32"`
}
```

### Pointers

Pointers are parsed using the same rules mentioned above.

## Included Dependendies

The only dependency included within `phnenv` is `github.com/stretchr/testify`, and it is only used in the unit test.
So you can rest easy in the knowledge that this library does not bring with it potentially malicious code.

## License

This project is licensed under the Apache 2.0 license.


*== 2021 phonaputer ==*
