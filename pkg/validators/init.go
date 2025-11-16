package validators

import (
	"errors"
	"fmt"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	once      sync.Once
	Validator *validator.Validate
)

func InitValidators() error {
	var initErr error

	once.Do(func() {
		Validator = validator.New()

		err := registerStringValidators(Validator)
		if err != nil {
			initErr = errors.Join(
				err,
				fmt.Errorf("failed to register string validators: %w", err),
			)
		}

		err = registerDecimalValidators(Validator)
		if err != nil {
			initErr = errors.Join(
				err,
				fmt.Errorf("failed to register decimal validators: %w", err),
			)
		}
	})

	return initErr
}

func registerStringValidators(v *validator.Validate) error {
	err := v.RegisterValidation("capitalized", CapitalizedValidator)
	if err != nil {
		return fmt.Errorf("failed register string capitalized validator: %w", err)
	}

	return nil
}

func registerDecimalValidators(v *validator.Validate) error {
	if err := v.RegisterValidation("decimal_gt", GreaterThanDecimal); err != nil {
		return fmt.Errorf("failed register decimal greater than validator: %w", err)
	}

	if err := v.RegisterValidation("decimal_gte", GreaterThanOrEqualDecimal); err != nil {
		return fmt.Errorf("failed register decimal greater than or equal validator: %w", err)
	}

	if err := v.RegisterValidation("decimal_lt", LessThanDecimal); err != nil {
		return fmt.Errorf("failed register decimal less than validator: %w", err)
	}

	if err := v.RegisterValidation("decimal_lte", LessThanOrEqualDecimal); err != nil {
		return fmt.Errorf("failed register decimal less than or equal validator: %w", err)
	}

	if err := v.RegisterValidation("decimal_gtfield", GreaterThanDecimalField); err != nil {
		return fmt.Errorf("failed register decimal greater than field validator: %w", err)
	}

	if err := v.RegisterValidation("decimal_gtefield", GreaterThanOrEqualDecimalField); err != nil {
		return fmt.Errorf("failed register decimal greater than or equal field validator: %w", err)
	}

	if err := v.RegisterValidation("decimal_ltfield", LessThanDecimalField); err != nil {
		return fmt.Errorf("failed register decimal less than field validator: %w", err)
	}

	if err := v.RegisterValidation("decimal_ltefield", LessThanOrEqualDecimalField); err != nil {
		return fmt.Errorf("failed register decimal less than or equal field validator: %w", err)
	}

	return nil
}
