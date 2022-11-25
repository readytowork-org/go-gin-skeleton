package validators

import (
	"boilerplate-api/errors"
	"fmt"

	validator "github.com/go-playground/validator/v10"
)

// UserValidator structure
type CategoryValidator struct {
	Validate *validator.Validate
}

// Register Custom Validators
func NewCategoryValidator() CategoryValidator {
	v := validator.New()
	return CategoryValidator{
		Validate: v,
	}
}

func (cv CategoryValidator) generateValidationMessage(field string, rule string) (message string) {
	switch rule {
	case "required":
		return fmt.Sprintf("Field '%s' is '%s'.", field, rule)
	default:
		return fmt.Sprintf("Field '%s' is not valid.", field)
	}
}

func (cv CategoryValidator) GenerateValidationResponse(err error) []errors.ErrorContext {
	var validations []errors.ErrorContext
	for _, value := range err.(validator.ValidationErrors) {
		field, rule := value.Field(), value.Tag()
		validation := errors.ErrorContext{Field: field, Message: cv.generateValidationMessage(field, rule)}
		validations = append(validations, validation)
	}
	return validations
}
