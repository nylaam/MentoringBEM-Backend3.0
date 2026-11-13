package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(s any) error {
	if err := validate.Struct(s); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var errMsgs []string
			for _, fe := range validationErrors {
				errMsgs = append(errMsgs, fmt.Sprintf("field '%s' gagal validasi '%s'", fe.Field(), fe.Tag()))
			}
			return fmt.Errorf("%s", strings.Join(errMsgs, ", "))
		}
		return err
	}
	return nil
}
