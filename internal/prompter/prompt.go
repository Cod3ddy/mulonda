package prompter

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Confirm asks for interactive consent and defaults to "no".
func Confirm(message string) bool {
	fmt.Printf("%s [y/N]: ", message)

	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return false
	}

	v := strings.ToLower(strings.TrimSpace(input))
	return v == "y" || v == "yes"
}
