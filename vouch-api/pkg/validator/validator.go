package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator wraps go-playground/validator with friendly error formatting.
type Validator struct {
	v *validator.Validate
}

// New constructs a Validator.
func New() *Validator {
	return &Validator{v: validator.New(validator.WithRequiredStructEnabled())}
}

// Struct validates a struct and returns a single human-readable error listing
// every failed field, or nil if valid.
func (val *Validator) Struct(s interface{}) error {
	if err := val.v.Struct(s); err != nil {
		var verrs validator.ValidationErrors
		if ok := asValidationErrors(err, &verrs); !ok {
			return err
		}
		msgs := make([]string, 0, len(verrs))
		for _, fe := range verrs {
			msgs = append(msgs, fmt.Sprintf("%s failed %q", strings.ToLower(fe.Field()), fe.Tag()))
		}
		return fmt.Errorf("validation: %s", strings.Join(msgs, "; "))
	}
	return nil
}

func asValidationErrors(err error, target *validator.ValidationErrors) bool {
	if v, ok := err.(validator.ValidationErrors); ok {
		*target = v
		return true
	}
	return false
}
