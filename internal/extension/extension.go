package extension

import (
	"strings"
)

type (
	// RegisterError is custom implementation of error interface
	RegisterError []error
)

// String return all error message like string with ";" delimiter
func (re RegisterError) String() string {
	var slice []string
	for _, err := range re {
		slice = append(slice, err.Error())
	}
	return strings.Join(slice, `; `)
}

// Error is implementation of error interface
func (re RegisterError) Error() string {
	return re.String()
}
