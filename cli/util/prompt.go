package util

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// Prompt prompts the user for input
func Prompt(question string) (string, error) {
	return PromptWithDefault(question, "")
}

// PromptWithDefault prompts the user for input with a default value
func PromptWithDefault(question, defaultValue string) (string, error) {
	if defaultValue != "" {
		fmt.Printf("%s %s[%s]%s: ", question, ColorGray, defaultValue, ColorReset)
	} else {
		fmt.Printf("%s: ", question)
	}
	
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	
	input = strings.TrimSpace(input)
	if input == "" && defaultValue != "" {
		return defaultValue, nil
	}
	
	return input, nil
}

// Confirm prompts the user for yes/no confirmation
func Confirm(question string) (bool, error) {
	return ConfirmWithDefault(question, false)
}

// ConfirmWithDefault prompts the user for yes/no with a default
func ConfirmWithDefault(question string, defaultValue bool) (bool, error) {
	defaultStr := "y/N"
	if defaultValue {
		defaultStr = "Y/n"
	}
	
	fmt.Printf("%s %s[%s]%s: ", question, ColorGray, defaultStr, ColorReset)
	
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	
	input = strings.TrimSpace(strings.ToLower(input))
	
	if input == "" {
		return defaultValue, nil
	}
	
	return input == "y" || input == "yes", nil
}

// Select prompts the user to select from a list of options
func Select(question string, options []string) (int, error) {
	return SelectWithDefault(question, options, 0)
}

// SelectWithDefault prompts the user to select from a list with a default
func SelectWithDefault(question string, options []string, defaultIndex int) (int, error) {
	fmt.Printf("%s\n", Bold(question))
	for i, opt := range options {
		marker := " "
		if i == defaultIndex {
			marker = ">"
		}
		fmt.Printf(" %s %d) %s\n", marker, i+1, opt)
	}
	
	fmt.Printf("\nSelect option %s[%d]%s: ", ColorGray, defaultIndex+1, ColorReset)
	
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}
	
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultIndex, nil
	}
	
	var selected int
	_, err = fmt.Sscanf(input, "%d", &selected)
	if err != nil || selected < 1 || selected > len(options) {
		return 0, fmt.Errorf("invalid selection")
	}
	
	return selected - 1, nil
}

// PromptReader is an interface for reading user input (for testing)
type PromptReader interface {
	ReadString(delim byte) (string, error)
}

// SetPromptReader sets a custom prompt reader (for testing)
var promptReader PromptReader = bufio.NewReader(os.Stdin)

// ReadLine reads a line from the prompt reader
func ReadLine() (string, error) {
	line, err := promptReader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

