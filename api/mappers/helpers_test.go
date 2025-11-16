package mappers_test

import (
	"testing"

	"github.com/k6zma/avito-test-assigniment/pkg/validators"
)

func initValidators(t *testing.T) {
	t.Helper()

	if err := validators.InitValidators(); err != nil {
		t.Fatalf("failed to init validators: %v", err)
	}
}

