package validators_test

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"

	"github.com/k6zma/avito-test-assigniment/pkg/validators"
)

const (
	initTestPrefix = "InitValidators"
)

type initStringsTestCase struct {
	testName   string
	inputValue string
	want       bool
}

type decimalInitStruct struct {
	GT  decimal.Decimal `validate:"decimal_gt=10"`
	GTE decimal.Decimal `validate:"decimal_gte=10"`
	LT  decimal.Decimal `validate:"decimal_lt=10"`
	LTE decimal.Decimal `validate:"decimal_lte=10"`
}

type initDecimalTestCase struct {
	testName string
	input    decimalInitStruct
	want     bool
}

func TestInitValidators_Idempotent(t *testing.T) {
	t.Run(fmt.Sprintf("[%s]-idempotentTest-firstCall", initTestPrefix), func(t *testing.T) {
		if err := validators.InitValidators(); err != nil {
			t.Errorf("InitValidators() first call error = %v, want nil", err)
		}
	})

	t.Run(fmt.Sprintf("[%s]-idempotentTest-secondCall", initTestPrefix), func(t *testing.T) {
		if err := validators.InitValidators(); err != nil {
			t.Errorf("InitValidators() second call error = %v, want nil", err)
		}
	})
}

func TestInitValidators_CapitalizedRegistered(t *testing.T) {
	if err := validators.InitValidators(); err != nil {
		t.Fatalf("InitValidators() error = %v, want nil", err)
	}

	tests := []initStringsTestCase{
		{"capitalized", "Test", true},
		{"not capitalized", "test", false},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("[%s]-%s-№%d", initTestPrefix, tt.testName, i+1), func(t *testing.T) {
			s := capitalizedValidationStruct{
				Value: tt.inputValue,
			}

			err := validators.Validator.Struct(s)

			got := err == nil
			if got != tt.want {
				t.Errorf(
					"capitalized validation failed for input value=%q: got %v, want %v (err: %v)",
					tt.inputValue, got, tt.want, err,
				)
			}
		})
	}
}

func TestInitValidators_DecimalValidatorsRegistered(t *testing.T) {
	if err := validators.InitValidators(); err != nil {
		t.Fatalf("InitValidators() error = %v, want nil", err)
	}

	tests := []initDecimalTestCase{
		{
			testName: "all_valid",
			input: decimalInitStruct{
				GT:  parseStringToDecimal(t, "10.01"),
				GTE: parseStringToDecimal(t, "10.00"),
				LT:  parseStringToDecimal(t, "9.99"),
				LTE: parseStringToDecimal(t, "10.00"),
			},
			want: true,
		},
		{
			testName: "invalid_gt",
			input: decimalInitStruct{
				GT:  parseStringToDecimal(t, "10.00"),
				GTE: parseStringToDecimal(t, "10.00"),
				LT:  parseStringToDecimal(t, "9.99"),
				LTE: parseStringToDecimal(t, "10.00"),
			},
			want: false,
		},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("[%s-decimal]-%s-№%d", initTestPrefix, tt.testName, i+1),
			func(t *testing.T) {
				err := validators.Validator.Struct(tt.input)

				got := err == nil
				if got != tt.want {
					t.Errorf(
						"decimal validators failed for test=%q: got %v, want %v (err: %v)",
						tt.testName, got, tt.want, err,
					)
				}
			},
		)
	}
}
