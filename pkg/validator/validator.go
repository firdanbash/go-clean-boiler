package validator

import (
	"github.com/firdanbash/go-clean-boiler/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct validates a struct and returns validation errors
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// BindAndValidate binds request body and validates it
func BindAndValidate(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return false
	}

	if err := ValidateStruct(obj); err != nil {
		validationErrors := FormatValidationErrors(err)
		response.BadRequest(c, "Validation failed", validationErrors)
		return false
	}

	return true
}

// FormatValidationErrors formats validator errors into a map
func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errors[e.Field()] = formatErrorMessage(e)
		}
	}

	return errors
}

func formatErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Minimum length is " + e.Param()
	case "max":
		return "Maximum length is " + e.Param()
	case "eqfield":
		return "Must match " + e.Param()
	default:
		return "Invalid value"
	}
}
