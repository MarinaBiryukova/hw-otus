package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

type stringValidators struct {
	len    *int64
	regexp *regexp.Regexp
	in     []string
}

type intValidators struct {
	min *int64
	max *int64
	in  []int64
}

const (
	lenPrefix    = "len:"
	regexpPrefix = "regexp:"
	inPrefix     = "in:"
	minPrefix    = "min:"
	maxPrefix    = "max:"
)

var (
	errNotStruct     = errors.New("not a struct")
	errInvalidLen    = errors.New("invalid length of string")
	errNoMatchRegexp = errors.New("string does not match regexp")
	errStringNotIn   = errors.New("string is not present in set")
	errViolatedMin   = errors.New("value is less than min")
	errViolatedMax   = errors.New("value is greater than max")
	errIntegerNotIn  = errors.New("integer is not present in set")
)

func (v ValidationErrors) Error() string {
	var res string
	for _, err := range v {
		res += err.Field + ": " + err.Err.Error() + "\n"
	}
	return res
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	tVal := val.Type()

	if val.Kind() != reflect.Struct {
		return errNotStruct
	}

	result := make(ValidationErrors, 0)
	for i := 0; i < val.NumField(); i++ {
		tField := tVal.Field(i)
		vField := val.Field(i)

		tag := tField.Tag.Get("validate")
		if tag == "" {
			continue
		}

		switch vField.Kind() { //nolint:exhaustive
		case reflect.String:
			err := checkString(vField.String(), tag)
			if err != nil {
				result = append(result, ValidationError{Field: tField.Name, Err: err})
			}
		case reflect.Int:
			err := checkInt(vField.Int(), tag)
			if err != nil {
				result = append(result, ValidationError{Field: tField.Name, Err: err})
			}
		case reflect.Slice:
			err := checkSlice(&vField, tField.Name, tag)
			if err != nil {
				result = append(result, *err)
			}
		default:
			continue
		}
	}

	if len(result) != 0 {
		return result
	}

	return nil
}

func checkString(s, tag string) error {
	v, err := parseStringValidators(tag)
	if err != nil {
		return err
	}

	if v.len != nil && int64(len(s)) != *v.len {
		return errInvalidLen
	}

	if v.regexp != nil && !v.regexp.MatchString(s) {
		return errNoMatchRegexp
	}

	if len(v.in) > 0 {
		var found bool
		for _, in := range v.in {
			if s == in {
				found = true
			}
		}
		if !found {
			return errStringNotIn
		}
	}

	return nil
}

func checkInt(i int64, tag string) error {
	v, err := parseIntValidators(tag)
	if err != nil {
		return err
	}

	if v.min != nil && i < *v.min {
		return errViolatedMin
	}

	if v.max != nil && i > *v.max {
		return errViolatedMax
	}

	if len(v.in) > 0 {
		var found bool
		for _, in := range v.in {
			if i == in {
				found = true
			}
		}
		if !found {
			return errIntegerNotIn
		}
	}

	return nil
}

func checkSlice(slice *reflect.Value, fieldName, tag string) *ValidationError {
	for i := 0; i < slice.Len(); i++ {
		sliceField := slice.Index(i)
		switch sliceField.Kind() { //nolint:exhaustive
		case reflect.String:
			err := checkString(sliceField.String(), tag)
			if err != nil {
				return &ValidationError{Field: fieldName, Err: err}
			}
		case reflect.Int:
			err := checkInt(sliceField.Int(), tag)
			if err != nil {
				return &ValidationError{Field: fieldName, Err: err}
			}
		default:
			continue
		}
	}

	return nil
}

func parseStringValidators(tag string) (*stringValidators, error) {
	var res stringValidators
	parts := strings.Split(tag, "|")

	for _, part := range parts {
		switch part[:strings.Index(part, ":")+1] {
		case lenPrefix:
			if res.len != nil {
				return nil, errors.New("duplicate len validator")
			}

			vLen, err := strconv.ParseInt(strings.TrimPrefix(part, lenPrefix), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("parse len: %w", err)
			}

			if vLen < 0 {
				return nil, errors.New("negative len")
			}

			res.len = &vLen
		case regexpPrefix:
			if res.regexp != nil {
				return nil, errors.New("duplicate regexp validator")
			}

			reg, err := regexp.Compile(strings.TrimPrefix(part, regexpPrefix))
			if err != nil {
				return nil, fmt.Errorf("compile regexp: %w", err)
			}

			res.regexp = reg
		case inPrefix:
			if len(res.in) > 0 {
				return nil, errors.New("duplicate in validator")
			}

			res.in = strings.Split(strings.TrimPrefix(part, inPrefix), ",")
		default:
			return nil, fmt.Errorf("unknown validator: %s", part)
		}
	}

	return &res, nil
}

func parseIntValidators(tag string) (*intValidators, error) {
	var res intValidators
	parts := strings.Split(tag, "|")

	for _, part := range parts {
		switch part[:strings.Index(part, ":")+1] {
		case minPrefix:
			if res.min != nil {
				return nil, errors.New("duplicate min validator")
			}

			vMin, err := strconv.ParseInt(strings.TrimPrefix(part, minPrefix), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("parse min: %w", err)
			}

			res.min = &vMin
		case maxPrefix:
			if res.max != nil {
				return nil, errors.New("duplicate max validator")
			}

			vMax, err := strconv.ParseInt(strings.TrimPrefix(part, maxPrefix), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("parse max: %w", err)
			}

			res.max = &vMax
		case inPrefix:
			if len(res.in) > 0 {
				return nil, errors.New("duplicate in validator")
			}

			vals := strings.Split(strings.TrimPrefix(part, inPrefix), ",")
			res.in = make([]int64, 0, len(vals))

			for _, strVal := range vals {
				val, err := strconv.ParseInt(strVal, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("parse in value: %w", err)
				}

				res.in = append(res.in, val)
			}
		default:
			return nil, fmt.Errorf("unknown validator: %s", part)
		}
	}

	return &res, nil
}
