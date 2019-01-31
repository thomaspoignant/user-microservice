package api

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	validator "gopkg.in/go-playground/validator.v8"
)

// register all the validator
func validatorsRegistration() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("emailvalidator", emailValidator)
	}
}

// valid that the email has a @ in it
func emailValidator(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	if email, ok := field.Interface().(string); ok {
		return strings.Contains(email, "@")
	}
	return false
}
