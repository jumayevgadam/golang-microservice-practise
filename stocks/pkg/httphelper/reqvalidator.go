package httphelper

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func RequestValidate(r *http.Request, data interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return fmt.Errorf("failed to decode request body: %w", err)
	}

	if err := ValidateRequest(data); err != nil {
		return fmt.Errorf("failed to validate request body: %w", err)
	}

	return nil
}

// ValidateRequest need for both services, grpc and http.
func ValidateRequest(data interface{}) error {
	return validate.Struct(data)
}
