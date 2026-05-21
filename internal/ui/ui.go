package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var noColor = os.Getenv("NO_COLOR") != ""

func color(code, s string) string {
	if noColor {
		return s
	}
	return code + s + "\033[0m"
}

func green(s string) string  { return color("\033[1;32m", s) }
func red(s string) string    { return color("\033[1;31m", s) }
func cyan(s string) string   { return color("\033[1;36m", s) }
func yellow(s string) string { return color("\033[1;33m", s) }
func dim(s string) string    { return color("\033[2m", s) }
func bold(s string) string   { return color("\033[1m", s) }

func Success(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%s %s\n", green(">>"), msg)
}

func Info(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%s %s\n", cyan("::"), msg)
}

func Step(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%s %s\n", dim("  ->"), msg)
}

func Error(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Fprintf(os.Stderr, "%s %s\n", red("!!"), msg)
}

func Fatal(format string, a ...any) {
	Error(format, a...)
	os.Exit(1)
}

func Warn(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%s %s\n", yellow("??"), msg)
}

func Banner(s string) {
	fmt.Printf("\n%s\n", bold(s))
}

func KeyValue(key, value string) {
	fmt.Printf("  %s %s\n", dim(key+":"), value)
}

// Confirm prompts the user and returns true if they answered yes.
func Confirm(prompt string) (bool, error) {
	fmt.Printf("%s %s %s ", yellow("??"), prompt, dim("[y/N]"))
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
