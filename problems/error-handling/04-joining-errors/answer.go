//go:build ignore

package main

import (
	"errors"
	"strings"
)

func ValidateForm(name string, age int, email string) error {
	var errs []error
	if name == "" {
		errs = append(errs, errors.New("name is required"))
	}
	if age < 0 || age > 150 {
		errs = append(errs, errors.New("age out of range"))
	}
	if !strings.Contains(email, "@") {
		errs = append(errs, errors.New("invalid email"))
	}
	return errors.Join(errs...)
}

func main() {}
