package colorfmt

import (
	"fmt"
)

const yellow = "\033[33m"
const reset = "\033[0m"
const italic = "\033[3m"

// Error formats a CLI-style error message, fully colored.
func Error(msg string, args ...any) error {
	text := fmt.Sprintf(msg, args...)
	// Prepend with your standard prefix
	return fmt.Errorf("%s%s%s", yellow+italic, "darkroom: <ERROR> "+text, reset)
}
