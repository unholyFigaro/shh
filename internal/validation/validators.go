package validation

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Validator func(value any) error

type Schema struct {
	Validators map[string]Validator
	Required   []string
}

func Validate(input map[string]any, s Schema) error {
	missingFields := []string{}
	for _, field := range s.Required {
		if _, ok := input[field]; !ok {
			missingFields = append(missingFields, field)
		}
	}
	if len(missingFields) > 0 {
		return fmt.Errorf("missing required fields: %v", missingFields)
	}
	validationErrors := []error{}
	for k, v := range input {
		if val, ok := s.Validators[k]; ok {
			err := val(v)
			if err != nil {
				validationErrors = append(validationErrors, fmt.Errorf("field %q: %w", k, err))
			}
		} else {
			validationErrors = append(validationErrors, fmt.Errorf("unknown field: %q", k))
		}
	}
	if len(validationErrors) > 0 {
		return errors.Join(validationErrors...)
	}
	return nil
}

func validateName(value any) error {
	s, ok := value.(string)
	if !ok || s == "" {
		return fmt.Errorf("name must be a non-empty string")
	}
	if strings.ContainsAny(s, " \t\n@:/") {
		return fmt.Errorf("name contains invalid characters")
	}
	return nil
}

func validateUser(value any) error {
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("user must be a string")
	}
	if strings.ContainsAny(s, " \t\n@:/") {
		return fmt.Errorf("user contains invalid characters")
	}
	return nil
}

func validatePort(v any) error {
	switch t := v.(type) {
	case int:
		if t == 0 {
			return nil
		}
		if t < 1 || t > 65535 {
			return fmt.Errorf("out of range: %d", t)
		}
		return nil
	case string:
		s := strings.TrimSpace(t)
		if s == "" {
			return nil
		}
		n, err := strconv.Atoi(s)
		if err != nil {
			return errors.New("must be integer")
		}
		return validatePort(n)
	default:
		return fmt.Errorf("invalid type: %T", t)
	}
}

func validateForce(value any) error {
	return nil
}

func validateHost(value any) error {
	return nil
}
