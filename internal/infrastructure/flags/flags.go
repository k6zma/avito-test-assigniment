package flags

import (
	"flag"
	"fmt"

	"github.com/k6zma/avito-test-assigniment/pkg/validators"
)

const (
	environmentFlagName         = "environment"
	environmentFlagDefaultValue = "dev"
	environmentFlagDescription  = "The application environment (dev, prod)"
)

var Environment = flag.String(
	environmentFlagName,
	environmentFlagDefaultValue,
	environmentFlagDescription,
)

type HubFlags struct {
	Environment string `validate:"required,oneof=dev prod"`
}

func GetHubFlags() (*HubFlags, error) {
	flag.Parse()

	result := &HubFlags{
		Environment: *Environment,
	}

	if err := validators.Validator.Struct(result); err != nil {
		return nil, fmt.Errorf("failed validate flags: %w", err)
	}

	return result, nil
}

func (f *HubFlags) String() string {
	return fmt.Sprintf("Environment: %s", f.Environment)
}
