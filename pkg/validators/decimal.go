package validators

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

type decimalComparator func(a, b decimal.Decimal) bool

func decimalParam(fl validator.FieldLevel) (decimal.Decimal, bool) {
	param := fl.Param()
	if param == "" {
		return decimal.Decimal{}, false
	}

	d, err := decimal.NewFromString(param)
	if err != nil {
		return decimal.Decimal{}, false
	}

	return d, true
}

func decimalFromValue(v reflect.Value) (decimal.Decimal, bool) {
	if !v.IsValid() {
		return decimal.Decimal{}, false
	}

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return decimal.Decimal{}, false
		}

		v = v.Elem()
	}

	d, ok := v.Interface().(decimal.Decimal)
	if !ok {
		return decimal.Decimal{}, false
	}

	return d, true
}

func decimalField(fl validator.FieldLevel) (decimal.Decimal, bool) {
	return decimalFromValue(fl.Field())
}

func decimalOtherField(fl validator.FieldLevel) (decimal.Decimal, bool) {
	fieldName := fl.Param()
	if fieldName == "" {
		return decimal.Decimal{}, false
	}

	parent := fl.Parent()
	if !parent.IsValid() {
		return decimal.Decimal{}, false
	}

	if parent.Kind() == reflect.Ptr {
		if parent.IsNil() {
			return decimal.Decimal{}, false
		}

		parent = parent.Elem()
	}

	if parent.Kind() != reflect.Struct {
		return decimal.Decimal{}, false
	}

	otherField := parent.FieldByName(fieldName)
	if !otherField.IsValid() {
		return decimal.Decimal{}, false
	}

	return decimalFromValue(otherField)
}

func compareWithParam(fl validator.FieldLevel, cmp decimalComparator) bool {
	threshold, ok := decimalParam(fl)
	if !ok {
		return false
	}

	value, ok := decimalField(fl)
	if !ok {
		return false
	}

	return cmp(value, threshold)
}

func compareWithField(fl validator.FieldLevel, cmp decimalComparator) bool {
	value, ok := decimalField(fl)
	if !ok {
		return false
	}

	other, ok := decimalOtherField(fl)
	if !ok {
		return false
	}

	return cmp(value, other)
}

func GreaterThanDecimal(fl validator.FieldLevel) bool {
	return compareWithParam(fl, func(a, b decimal.Decimal) bool {
		return a.GreaterThan(b)
	})
}

func GreaterThanOrEqualDecimal(fl validator.FieldLevel) bool {
	return compareWithParam(fl, func(a, b decimal.Decimal) bool {
		return a.GreaterThanOrEqual(b)
	})
}

func LessThanDecimal(fl validator.FieldLevel) bool {
	return compareWithParam(fl, func(a, b decimal.Decimal) bool {
		return a.LessThan(b)
	})
}

func LessThanOrEqualDecimal(fl validator.FieldLevel) bool {
	return compareWithParam(fl, func(a, b decimal.Decimal) bool {
		return a.LessThanOrEqual(b)
	})
}

func GreaterThanDecimalField(fl validator.FieldLevel) bool {
	return compareWithField(fl, func(a, b decimal.Decimal) bool {
		return a.GreaterThan(b)
	})
}

func GreaterThanOrEqualDecimalField(fl validator.FieldLevel) bool {
	return compareWithField(fl, func(a, b decimal.Decimal) bool {
		return a.GreaterThanOrEqual(b)
	})
}

func LessThanDecimalField(fl validator.FieldLevel) bool {
	return compareWithField(fl, func(a, b decimal.Decimal) bool {
		return a.LessThan(b)
	})
}

func LessThanOrEqualDecimalField(fl validator.FieldLevel) bool {
	return compareWithField(fl, func(a, b decimal.Decimal) bool {
		return a.LessThanOrEqual(b)
	})
}
