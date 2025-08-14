package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

func Prompt(label string) string {
	fmt.Print(label)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func PromptInt(label string) (int, error) {
	var n int
	fmt.Print(label)
	_, err := fmt.Scan(&n)
	return n, err
}
