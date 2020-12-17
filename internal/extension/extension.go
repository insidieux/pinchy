package extension

import (
	"strings"
)

type (
	RegisterError []error
)

func (re RegisterError) String() string {
	var slice []string
	for _, err := range re {
		slice = append(slice, err.Error())
	}
	return strings.Join(slice, `; `)
}

func (re RegisterError) Error() string {
	return re.String()
}
