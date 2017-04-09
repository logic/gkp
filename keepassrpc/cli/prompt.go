package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Prompt is a simple implementation of keepassrpc.Passworder.
func Prompt() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the code provided by KeePass: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}
