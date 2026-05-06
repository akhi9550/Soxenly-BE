package cleanvalidator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func GetErrorMessage(err error) string {
	if errs, ok := err.(validator.ValidationErrors); ok {
		var errMsgs []string
		for _, e := range errs {
			field := e.Field()
			tag := e.Tag()

			switch tag {
			case "required":
				errMsgs = append(errMsgs, fmt.Sprintf("%s is required", field))
			case "email":
				errMsgs = append(errMsgs, fmt.Sprintf("%s must be a valid email", field))
			case "min":
				errMsgs = append(errMsgs, fmt.Sprintf("%s must be at least %s characters long", field, e.Param()))
			default:
				errMsgs = append(errMsgs, fmt.Sprintf("%s failed on %s", field, tag))
			}
		}
		return strings.Join(errMsgs, ", ")
	}
	return err.Error()
}
