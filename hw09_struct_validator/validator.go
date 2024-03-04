package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var sb strings.Builder
	for i, e := range v {
		sb.WriteString(e.Field)
		sb.WriteString(": ")
		sb.WriteString(e.Err.Error())
		if i < len(v)-1 {
			sb.WriteString("; ")
		}
	}
	return sb.String()
}

func Validate(v interface{}) error { //nolint:gocognit
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return errors.New("validation type must be a struct")
	}

	var validationErrors ValidationErrors

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldValue := val.Field(i)

		validateTag := field.Tag.Get("validate")
		if validateTag == "" {
			continue
		}

		// Проверяем, является ли первая буква имени поля заглавной.
		if field.PkgPath != "" {
			continue
		}

		defer func() {
			if r := recover(); r != nil {
				fmt.Println("painic: ", r)
				// Если паника, преобразуем ее в ошибку и добавляем в список validationErrors и это не сработает :).
				validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: fmt.Errorf("panic: %v", r)})
			}
		}()

		switch fieldValue.Kind() { //nolint:exhaustive
		case reflect.String:
			if err := ValidateString(fieldValue.String(), validateTag); err != nil {
				validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: err})
			}
		case reflect.Int:
			if err := ValidateInt(fieldValue.Int(), validateTag); err != nil {
				validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: err})
			}
		case reflect.Slice:
			// В случае если вдруг поле приватно, то проверка что бы избежать паники вида.
			// panic: reflect.Value.Interface: cannot return value obtained from unexported field or method
			if !fieldValue.CanInterface() {
				validationErrors = append(validationErrors, ValidationError{
					Field: field.Name,
					Err:   errors.New("cannot return value obtained from unexported field or method"),
				})
				continue
			}
			if fieldValue.Type().Elem().Kind() == reflect.String {
				if err := ValidateStringSlice(fieldValue.Interface().([]string), validateTag); err != nil {
					validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: err})
				}
			}
			if fieldValue.Type().Elem().Kind() == reflect.Int {
				if err := ValidateIntSlice(fieldValue.Interface().([]int), validateTag); err != nil {
					validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: err})
				}
			}
		default:
			continue
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}

func ValidateString(fieldValue string, validateTag string) error {
	rules := strings.Split(validateTag, "|")
	for _, rule := range rules {
		switch {
		case strings.HasPrefix(rule, "len:"):
			expectedLength, err := strconv.Atoi(strings.TrimPrefix(rule, "len:"))
			if err != nil {
				return fmt.Errorf("failed to parse len: %w", err)
			}
			if len(fieldValue) != expectedLength {
				return fmt.Errorf("validate: must be %d characters long", expectedLength)
			}
		case strings.HasPrefix(rule, "regexp:"):
			regexpPattern := strings.TrimPrefix(rule, "regexp:")
			match, err := regexp.MatchString(regexpPattern, fieldValue)
			if err != nil {
				return fmt.Errorf("failed to parse regexp: %w", err)
			}
			if !match {
				return fmt.Errorf("validate: must match regexp %s", regexpPattern)
			}
		case strings.HasPrefix(rule, "in:"):
			validValues := strings.Split(strings.TrimPrefix(rule, "in:"), ",")
			found := false
			for _, validValue := range validValues {
				if fieldValue == validValue {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("validate: must be one of %s", validValues)
			}
		}
	}
	return nil
}

func ValidateStringSlice(fieldValue []string, validateTag string) error {
	for _, value := range fieldValue {
		err := ValidateString(value, validateTag)
		if err != nil {
			return fmt.Errorf("validate: %w", err)
		}
	}

	return nil
}

func ValidateInt(fieldValue int64, validateTag string) error {
	rules := strings.Split(validateTag, "|")
	for _, rule := range rules {
		switch {
		case strings.HasPrefix(rule, "min:"):
			minValue, err := strconv.ParseInt(rule[4:], 10, 64)
			if err != nil {
				return fmt.Errorf("validate: invalid min value %s, %w", rule[4:], err)
			}
			if fieldValue < minValue {
				return fmt.Errorf("validate: must be at least %d", minValue)
			}
		case strings.HasPrefix(rule, "max:"):
			maxValue, err := strconv.ParseInt(rule[4:], 10, 64)
			if err != nil {
				return fmt.Errorf("validate: invalid max value %s, %w", rule[4:], err)
			}
			if fieldValue > maxValue {
				return fmt.Errorf("validate: must be less than %d", maxValue)
			}
		case strings.HasPrefix(rule, "in:"):
			validValues := strings.Split(rule[3:], ",")
			found := false
			for _, validValue := range validValues {
				value, err := strconv.ParseInt(validValue, 10, 64)
				if err != nil {
					return fmt.Errorf("validate: invalid value %s, %w", validValue, err)
				}
				if fieldValue == value {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("validate: must be one of %s", validValues)
			}
		}
	}
	return nil
}

func ValidateIntSlice(fieldValue []int, validateTag string) error {
	for _, value := range fieldValue {
		err := ValidateInt(int64(value), validateTag)
		if err != nil {
			return fmt.Errorf("validate: %w", err)
		}
	}

	return nil
}
