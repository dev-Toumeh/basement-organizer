package validate

import (
	"basement/main/internal/logg"
	"errors"
	"fmt"
	"html"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/gofrs/uuid/v5"
)

// StringField
type StringField struct {
	Input               string
	Value               string
	DefaultMaxLength    int64
	DefaultMinLength    int64
	DefaultRegexPattern string
}

func NewStringField(input string) StringField {
	trimmed := strings.TrimSpace(input)
	sanitized := html.EscapeString(trimmed)

	return StringField{
		Input:               input,
		Value:               sanitized,
		DefaultMaxLength:    255,
		DefaultMinLength:    1,
		DefaultRegexPattern: `^[\p{L}\p{N} _.-]+$`,
	}
}

func (s StringField) String() string   { return s.Value }
func (s StringField) IsEmpty() bool    { return s.Value == "" }
func (s StringField) MaxLength() error { return s.MaxLengthCustom(s.DefaultMaxLength) }
func (s StringField) MaxLengthCustom(limit int64) error {
	if int64(len(s.Value)) > limit {
		return fmt.Errorf("string exceeds maximum length of %d", limit)
	}
	return nil
}
func (s StringField) MinLength() error { return s.MinLengthCustom(s.DefaultMinLength) }
func (s StringField) MinLengthCustom(limit int64) error {
	if int64(len(s.Value)) < limit {
		return fmt.Errorf("string shorter than minimum length of %d", limit)
	}
	return nil
}
func (s StringField) MatchesRegex() error { return s.MatchesRegexCustom(s.DefaultRegexPattern) }
func (s StringField) MatchesRegexCustom(pattern string) error {
	matched, err := regexp.MatchString(pattern, s.Value)
	if err != nil {
		return fmt.Errorf("regex error: %w", err)
	}
	if !matched {
		return errors.New("string does not match required pattern")
	}
	return nil
}
func (s StringField) IsEmailFormat() error {
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	return s.MatchesRegexCustom(emailRegex)
}

func (s StringField) ValidatePictureFormat() error {
	allowed := map[string]bool{"image/png": true, "image/jpeg": true, "image/jpg": true}
	if !allowed[strings.TrimSpace(s.Input)] {
		return fmt.Errorf("unsupported image format: %s", s.Input)
	}
	return nil
}

// IntField
type IntField struct {
	value           int64
	input           string
	DefaultMaxValue int64
	DefaultMinValue int64
	Err             error
}

func NewIntField(input string) IntField {
	val, err := strconv.ParseInt(strings.TrimSpace(input), 10, 64)
	return IntField{
		input:           input,
		value:           val,
		DefaultMaxValue: math.MaxInt64,
		DefaultMinValue: 0,
		Err:             err,
	}
}

func (n IntField) Int() int64 { return n.value }
func (n IntField) IsZeroOrPositive() error {
	if n.value < 0 {
		return errors.New("integer value is negative")
	}
	return n.Err
}
func (n IntField) IsPositive() error {
	if n.value <= 0 {
		return errors.New("integer value is not positive")
	}
	return n.Err
}
func (n IntField) MinValueCustom(min int64) error {
	if n.value < min {
		return fmt.Errorf("integer value less than minimum (%d)", min)
	}
	return n.Err
}
func (n IntField) MaxValueCustom(max int64) error {
	if n.value > max {
		return fmt.Errorf("integer value greater than maximum (%d)", max)
	}
	return n.Err
}
func (n IntField) MinValue() error { return n.MinValueCustom(n.DefaultMinValue) }
func (n IntField) MaxValue() error { return n.MaxValueCustom(n.DefaultMaxValue) }
func (s IntField) IsEmpty() bool   { return s.input == "" }

type FloatField struct {
	value           float64
	input           string
	DefaultMaxValue float64
	DefaultMinValue float64
	Err             error
}

// NewFloatField
func NewFloatField(input string) FloatField {
	logg.Debugf("input is: %s", input)
	val, err := strconv.ParseFloat(strings.TrimSpace(input), 64)
	return FloatField{
		input:           input,
		value:           val,
		DefaultMaxValue: math.MaxFloat64,
		DefaultMinValue: 0.0,
		Err:             err,
	}
}

func (f FloatField) Float64() float64 { return f.value }
func (f FloatField) IsZeroOrPositive() error {
	if f.value < 0 {
		return errors.New("float value is negative")
	}
	return nil
}
func (f FloatField) IsPositive() error {
	if f.value <= 0 {
		return errors.New("float value is not positive")
	}
	return nil
}
func (n FloatField) MinValue() error { return n.MinValueCustom(n.DefaultMinValue) }
func (n FloatField) MaxValue() error { return n.MaxValueCustom(n.DefaultMaxValue) }
func (n FloatField) MinValueCustom(min float64) error {
	if n.value < min {
		return fmt.Errorf("float value less than minimum (%f)", min)
	}
	return nil
}
func (n FloatField) MaxValueCustom(max float64) error {
	if n.value > max {
		return fmt.Errorf("float value greater than maximum (%f)", max)
	}
	return nil
}

func (f FloatField) IsEmpty() bool { return f.input == "" }

type UUIDField struct {
	Value uuid.UUID
	Input string
	Err   error
}

// NewUUIDField
func NewUUIDField(input string) UUIDField {
	input = strings.TrimSpace(input)
	value, err := uuid.FromString(input)
	return UUIDField{Value: value, Input: input, Err: err}
}

func (u UUIDField) UUID() uuid.UUID { return u.Value }
func (u UUIDField) IsValid() error  { return u.Err }
func (u UUIDField) IsNil() bool     { return u.Value == uuid.Nil }
func (s UUIDField) IsEmpty() bool   { return s.Input == "" }
