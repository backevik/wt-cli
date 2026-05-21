package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var noColor = os.Getenv("NO_COLOR") != ""

func green(s string) string {
	if noColor {
		return s
	}
	return "\033[32m" + s + "\033[0m"
}

func red(s string) string {
	if noColor {
		return s
	}
	return "\033[31m" + s + "\033[0m"
}

func Success(format string, a ...any) {
	fmt.Printf(green("✓ ")+format+"\n", a...)
}

func Info(format string, a ...any) {
	fmt.Printf(format+"\n", a...)
}

func Error(format string, a ...any) {
	fmt.Fprintf(os.Stderr, red("error: ")+format+"\n", a...)
}

func Fatal(format string, a ...any) {
	Error(format, a...)
	os.Exit(1)
}

// Confirm prompts the user and returns true if they answered yes.
func Confirm(prompt string) (bool, error) {
	fmt.Printf("%s [y/N] ", prompt)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		reply := strings.TrimSpace(scanner.Text())
		return strings.EqualFold(reply, "y") || strings.EqualFold(reply, "yes"), nil
	}
	if err := scanner.Err(); err != nil {
		return false, err
	}
	return false, nil
}
