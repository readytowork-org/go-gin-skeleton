package user

import (
	"boilerplate-api/config/constants"
	"boilerplate-api/config/errors"
	"fmt"
	"regexp"

	validator "github.com/go-playground/validator/v10"
)

// Validator structure
type Validator struct {
	Validate *validator.Validate
}

// Register Custom Validators
func UserValidator() Validator {
	v := validator.New()
	_ = v.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		if fl.Field().String() != "" {
			match, _ := regexp.MatchString("^[- +()]*[0-9][- +()0-9]*$", fl.Field().String())
			return match
		}
		return true
	})
	_ = v.RegisterValidation("gender", func(fl validator.FieldLevel) bool {
		if fl.Field().String() != "" {
			var val_type constants.Gender
			if err := val_type.IsValidVal(fl.Field().String()); err != nil {
				return false
			}
		}
		return true
	})
	_ = v.RegisterValidation("email", func(fl validator.FieldLevel) bool {
		if fl.Field().String() != "" {
			match, _ := regexp.MatchString(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`, fl.Field().String())
			return match
		}
		return true
	})
	return Validator{
		Validate: v,
	}
}

func (cv Validator) generateValidationMessage(field string, rule string) (message string) {
	switch rule {
	case "required":
		return fmt.Sprintf("Field '%s' is '%s'.", field, rule)
	case "phone":
		return fmt.Sprintf("Field '%s' is not valid.", field)
	case "gender":
		return fmt.Sprintf("Field '%s' is not valid.", field)
	case "email":
		return fmt.Sprintf("Field '%s' is not valid.", field)
	default:
		return fmt.Sprintf("Field '%s' is not valid.", field)
	}
}

func (cv Validator) GenerateValidationResponse(err error) []errors.ErrorContext {
	var validations []errors.ErrorContext
	for _, value := range err.(validator.ValidationErrors) {
		field, rule := value.Field(), value.Tag()
		validation := errors.ErrorContext{Field: field, Message: cv.generateValidationMessage(field, rule)}
		validations = append(validations, validation)
	}
	return validations
}
