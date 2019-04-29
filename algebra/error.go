package algebra

import "fmt"

type undefinedError struct {
	message string
}

func (e *undefinedError) Error() string {
	return fmt.Sprintf("The result is not defined: %s", e.message)
}
