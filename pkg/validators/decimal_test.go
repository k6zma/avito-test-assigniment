package validators_test

import (
	"fmt"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"

	"github.com/k6zma/avito-test-assigniment/pkg/validators"
)

const decimalTestPrefix = "DecimalValidator"

type decimalGTStruct struct {
	Value decimal.Decimal `validate:"decimal_gt=10"`
}

type decimalGTEStruct struct {
	Value decimal.Decimal `validate:"decimal_gte=10"`
}

type decimalLTStruct struct {
	Value decimal.Decimal `validate:"decimal_lt=10"`
}

type decimalLTEStruct struct {
	Value decimal.Decimal `validate:"decimal_lte=10"`
}

type decimalTestCase struct {
	testName   string
	inputValue decimal.Decimal
	want       bool
}

type decimalGTFieldStruct struct {
	Value decimal.Decimal `validate:"decimal_gtfield=Other"`
	Other decimal.Decimal
}

type decimalGTEFieldStruct struct {
	Value decimal.Decimal `validate:"decimal_gtefield=Other"`
	Other decimal.Decimal
}

type decimalLTFieldStruct struct {
	Value decimal.Decimal `validate:"decimal_ltfield=Other"`
	Other decimal.Decimal
}

type decimalLTEFieldStruct struct {
	Value decimal.Decimal `validate:"decimal_ltefield=Other"`
	Other decimal.Decimal
}

type decimalFieldTestCase struct {
	testName string
	value    decimal.Decimal
	other    decimal.Decimal
	want     bool
}

func parseStringToDecimal(t *testing.T, s string) decimal.Decimal {
	t.Helper()

	val, err := decimal.NewFromString(s)
	if err != nil {
		t.Fatalf("failed to parse %q as decimal: %v", s, err)
	}

	return val
}

func TestDecimalGreaterThanValidator(t *testing.T) {
	tests := []decimalTestCase{
		{"less than threshold", parseStringToDecimal(t, "9.99"), false},
		{"equal to threshold", parseStringToDecimal(t, "10.00"), false},
		{"greater than threshold", parseStringToDecimal(t, "10.01"), true},
	}

	v := validator.New()
	if err := v.RegisterValidation("decimal_gt", validators.GreaterThanDecimal); err != nil {
		t.Fatalf("failed to register decimal_gt validator: %v", err)
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("[%s-GT]-%s-№%d", decimalTestPrefix, tt.testName, i+1),
			func(t *testing.T) {
				s := decimalGTStruct{
					Value: tt.inputValue,
				}

				err := v.Struct(s)

				got := err == nil
				if got != tt.want {
					t.Errorf(
						"decimal_gt validation failed for input value=%q: got %v, want %v (err: %v)",
						tt.inputValue,
						got,
						tt.want,
						err,
					)
				}
			},
		)
	}
}

func TestDecimalGreaterThanOrEqualValidator(t *testing.T) {
	tests := []decimalTestCase{
		{"less than threshold", parseStringToDecimal(t, "9.99"), false},
		{"equal to threshold", parseStringToDecimal(t, "10.00"), true},
		{"greater than threshold", parseStringToDecimal(t, "10.01"), true},
	}

	v := validator.New()
	if err := v.RegisterValidation("decimal_gte", validators.GreaterThanOrEqualDecimal); err != nil {
		t.Fatalf("failed to register decimal_gte validator: %v", err)
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("[%s-GTE]-%s-№%d", decimalTestPrefix, tt.testName, i+1),
			func(t *testing.T) {
				s := decimalGTEStruct{
					Value: tt.inputValue,
				}

				err := v.Struct(s)

				got := err == nil
				if got != tt.want {
					t.Errorf(
						"decimal_gte validation failed for input value=%q: got %v, want %v (err: %v)",
						tt.inputValue,
						got,
						tt.want,
						err,
					)
				}
			},
		)
	}
}

func TestDecimalLessThanValidator(t *testing.T) {
	tests := []decimalTestCase{
		{"less than threshold", parseStringToDecimal(t, "9.99"), true},
		{"equal to threshold", parseStringToDecimal(t, "10.00"), false},
		{"greater than threshold", parseStringToDecimal(t, "10.01"), false},
	}

	v := validator.New()
	if err := v.RegisterValidation("decimal_lt", validators.LessThanDecimal); err != nil {
		t.Fatalf("failed to register decimal_lt validator: %v", err)
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("[%s-LT]-%s-№%d", decimalTestPrefix, tt.testName, i+1),
			func(t *testing.T) {
				s := decimalLTStruct{
					Value: tt.inputValue,
				}

				err := v.Struct(s)

				got := err == nil
				if got != tt.want {
					t.Errorf(
						"decimal_lt validation failed for input value=%q: got %v, want %v (err: %v)",
						tt.inputValue,
						got,
						tt.want,
						err,
					)
				}
			},
		)
	}
}

func TestDecimalLessThanOrEqualValidator(t *testing.T) {
	tests := []decimalTestCase{
		{"less than threshold", parseStringToDecimal(t, "9.99"), true},
		{"equal to threshold", parseStringToDecimal(t, "10.00"), true},
		{"greater than threshold", parseStringToDecimal(t, "10.01"), false},
	}

	v := validator.New()
	if err := v.RegisterValidation("decimal_lte", validators.LessThanOrEqualDecimal); err != nil {
		t.Fatalf("failed to register decimal_lte validator: %v", err)
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("[%s-LTE]-%s-№%d", decimalTestPrefix, tt.testName, i+1),
			func(t *testing.T) {
				s := decimalLTEStruct{
					Value: tt.inputValue,
				}

				err := v.Struct(s)

				got := err == nil
				if got != tt.want {
					t.Errorf(
						"decimal_lte validation failed for input value=%q: got %v, want %v (err: %v)",
						tt.inputValue,
						got,
						tt.want,
						err,
					)
				}
			},
		)
	}
}

func TestDecimalGreaterThanFieldValidator(t *testing.T) {
	tests := []decimalFieldTestCase{
		{
			testName: "value_less_than_other",
			value:    parseStringToDecimal(t, "9.99"),
			other:    parseStringToDecimal(t, "10.00"),
			want:     false,
		},
		{
			testName: "value_equal_other",
			value:    parseStringToDecimal(t, "10.00"),
			other:    parseStringToDecimal(t, "10.00"),
			want:     false,
		},
		{
			testName: "value_greater_than_other",
			value:    parseStringToDecimal(t, "10.01"),
			other:    parseStringToDecimal(t, "10.00"),
			want:     true,
		},
	}

	v := validator.New()
	if err := v.RegisterValidation("decimal_gtfield", validators.GreaterThanDecimalField); err != nil {
		t.Fatalf("failed to register decimal_gtfield validator: %v", err)
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("[%s-GTField]-%s-№%d", decimalTestPrefix, tt.testName, i+1),
			func(t *testing.T) {
				s := decimalGTFieldStruct{
					Value: tt.value,
					Other: tt.other,
				}

				err := v.Struct(s)

				got := err == nil
				if got != tt.want {
					t.Errorf(
						"decimal_gtfield validation failed for value=%q, other=%q: got %v, want %v (err: %v)",
						tt.value,
						tt.other,
						got,
						tt.want,
						err,
					)
				}
			},
		)
	}
}

func TestDecimalGreaterThanOrEqualFieldValidator(t *testing.T) {
	tests := []decimalFieldTestCase{
		{
			testName: "value_less_than_other",
			value:    parseStringToDecimal(t, "9.99"),
			other:    parseStringToDecimal(t, "10.00"),
			want:     false,
		},
		{
			testName: "value_equal_other",
			value:    parseStringToDecimal(t, "10.00"),
			other:    parseStringToDecimal(t, "10.00"),
			want:     true,
		},
		{
			testName: "value_greater_than_other",
			value:    parseStringToDecimal(t, "10.01"),
			other:    parseStringToDecimal(t, "10.00"),
			want:     true,
		},
	}

	v := validator.New()
	if err := v.RegisterValidation("decimal_gtefield", validators.GreaterThanOrEqualDecimalField); err != nil {
		t.Fatalf("failed to register decimal_gtefield validator: %v", err)
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("[%s-GTEField]-%s-№%d", decimalTestPrefix, tt.testName, i+1),
			func(t *testing.T) {
				s := decimalGTEFieldStruct{
					Value: tt.value,
					Other: tt.other,
				}

				err := v.Struct(s)

				got := err == nil
				if got != tt.want {
					t.Errorf(
						"decimal_gtefield validation failed for value=%q, other=%q: got %v, want %v (err: %v)",
						tt.value,
						tt.other,
						got,
						tt.want,
						err,
					)
				}
			},
		)
	}
}

func TestDecimalLessThanFieldValidator(t *testing.T) {
	tests := []decimalFieldTestCase{
		{
			testName: "value_less_than_other",
			value:    parseStringToDecimal(t, "9.99"),
			other:    parseStringToDecimal(t, "10.00"),
			want:     true,
		},
		{
			testName: "value_equal_other",
			value:    parseStringToDecimal(t, "10.00"),
			other:    parseStringToDecimal(t, "10.00"),
			want:     false,
		},
		{
			testName: "value_greater_than_other",
			value:    parseStringToDecimal(t, "10.01"),
			other:    parseStringToDecimal(t, "10.00"),
			want:     false,
		},
	}

	v := validator.New()
	if err := v.RegisterValidation("decimal_ltfield", validators.LessThanDecimalField); err != nil {
		t.Fatalf("failed to register decimal_ltfield validator: %v", err)
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("[%s-LTField]-%s-№%d", decimalTestPrefix, tt.testName, i+1),
			func(t *testing.T) {
				s := decimalLTFieldStruct{
					Value: tt.value,
					Other: tt.other,
				}

				err := v.Struct(s)

				got := err == nil
				if got != tt.want {
					t.Errorf(
						"decimal_ltfield validation failed for value=%q, other=%q: got %v, want %v (err: %v)",
						tt.value,
						tt.other,
						got,
						tt.want,
						err,
					)
				}
			},
		)
	}
}

func TestDecimalLessThanOrEqualFieldValidator(t *testing.T) {
	tests := []decimalFieldTestCase{
		{
			testName: "value_less_than_other",
			value:    parseStringToDecimal(t, "9.99"),
			other:    parseStringToDecimal(t, "10.00"),
			want:     true,
		},
		{
			testName: "value_equal_other",
			value:    parseStringToDecimal(t, "10.00"),
			other:    parseStringToDecimal(t, "10.00"),
			want:     true,
		},
		{
			testName: "value_greater_than_other",
			value:    parseStringToDecimal(t, "10.01"),
			other:    parseStringToDecimal(t, "10.00"),
			want:     false,
		},
	}

	v := validator.New()
	if err := v.RegisterValidation("decimal_ltefield", validators.LessThanOrEqualDecimalField); err != nil {
		t.Fatalf("failed to register decimal_ltefield validator: %v", err)
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("[%s-LTEField]-%s-№%d", decimalTestPrefix, tt.testName, i+1),
			func(t *testing.T) {
				s := decimalLTEFieldStruct{
					Value: tt.value,
					Other: tt.other,
				}

				err := v.Struct(s)

				got := err == nil
				if got != tt.want {
					t.Errorf(
						"decimal_ltefield validation failed for value=%q, other=%q: got %v, want %v (err: %v)",
						tt.value,
						tt.other,
						got,
						tt.want,
						err,
					)
				}
			},
		)
	}
}
