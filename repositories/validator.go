package repositories

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

type ValidationErrors struct {
	errors map[string][]string
}

func (validationErrors *ValidationErrors) add(field string, message string) {
	validationErrors.errors[field] = append(validationErrors.errors[field], message)
}

func (validationErrors *ValidationErrors) All() map[string][]string {
	return validationErrors.errors
}

// Validator dengan gaya validatorjs:
//
//	validation := repositories.NewValidator(body, map[string]string{
//		"email"    : "required|email",
//		"password" : "required|min:6",
//	})
//	if validation.Fails() {
//		validation.Errors.All()
//	}
//
// Rule: required, email, min, max, numeric, string, in
type Validator struct {
	body   Body
	rules  map[string]string
	Errors *ValidationErrors
}

func NewValidator(body Body, rules map[string]string) *Validator {
	return &Validator{body: body, rules: rules, Errors: &ValidationErrors{errors: map[string][]string{}}}
}

func (validator *Validator) Fails() bool {
	return !validator.Passes()
}

func (validator *Validator) Passes() bool {
	validator.Errors = &ValidationErrors{errors: map[string][]string{}}

	for field, rule := range validator.rules {
		value := validator.body[field]
		items := strings.Split(rule, "|")
		label := strings.ReplaceAll(field, "_", " ")

		if isEmpty(value) {
			if slices.Contains(items, "required") {
				validator.Errors.add(field, fmt.Sprintf("The %s field is required.", label))
			}
			continue
		}

		_, isNumber := value.(float64)
		numeric := slices.Contains(items, "numeric") || isNumber
		text := validator.body.String(field)

		for _, item := range items {
			name, param, _ := strings.Cut(item, ":")

			switch name {
			case "email":
				if !emailRegex.MatchString(text) {
					validator.Errors.add(field, fmt.Sprintf("The %s format is invalid.", label))
				}
			case "numeric":
				if _, err := strconv.ParseFloat(text, 64); err != nil {
					validator.Errors.add(field, fmt.Sprintf("The %s must be a number.", label))
				}
			case "string":
				if _, ok := value.(string); !ok {
					validator.Errors.add(field, fmt.Sprintf("The %s must be a string.", label))
				}
			case "in":
				if !slices.Contains(strings.Split(param, ","), text) {
					validator.Errors.add(field, fmt.Sprintf("The selected %s is invalid.", label))
				}
			case "min":
				limit, _ := strconv.ParseFloat(param, 64)
				if number, err := strconv.ParseFloat(text, 64); numeric && err == nil {
					if number < limit {
						validator.Errors.add(field, fmt.Sprintf("The %s must be at least %s.", label, param))
					}
				} else if !numeric && float64(len([]rune(text))) < limit {
					validator.Errors.add(field, fmt.Sprintf("The %s must be at least %s characters.", label, param))
				}
			case "max":
				limit, _ := strconv.ParseFloat(param, 64)
				if number, err := strconv.ParseFloat(text, 64); numeric && err == nil {
					if number > limit {
						validator.Errors.add(field, fmt.Sprintf("The %s may not be greater than %s.", label, param))
					}
				} else if !numeric && float64(len([]rune(text))) > limit {
					validator.Errors.add(field, fmt.Sprintf("The %s may not be greater than %s characters.", label, param))
				}
			}
		}
	}

	return len(validator.Errors.errors) == 0
}
