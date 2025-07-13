package validation

import (
	"fmt"
	"reflect"
)

type BaseModel interface {
	Validate() error
}

type validator struct {
	checks       []bool
	descriptions []string
}

func NewValidator(model any) validator {
	return validator{
		checks:       make([]bool, 0, reflect.ValueOf(model).NumField()),
		descriptions: make([]string, 0, reflect.ValueOf(model).NumField()),
	}
}

func (v validator) Add(check bool, desc string) validator {
	v.checks = append(v.checks, check)
	v.descriptions = append(v.descriptions, desc)
	return v
}

func (v validator) Validate() error {
	for ix, check := range v.checks {
		desc := v.descriptions[ix]
		if !check {
			return fmt.Errorf("validation error: %s", desc)
		}
	}
	return nil
}
