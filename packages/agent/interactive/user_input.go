package interactive

import (
	"deckronomicon/packages/agent/actionparser"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func (a *Agent) Prompt(message string) {
	fmt.Printf("%s %s", message, a.prompt)
}

// ReadNumberMany prompts the user to enter many numbers, and returns a slice of
// numbers.
func (a *Agent) ReadNumberMany(max int) ([]int, error) {
	tokens := a.ReadInputTokens()
	var numbers []int
	for _, part := range tokens {
		number, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("failed to parse number: %w", err)
		}
		if number < 0 {
			return nil, fmt.Errorf("number less than zero")
		}
		if max != -1 && number > max {
			return nil, fmt.Errorf("number greater than max")
		}
		numbers = append(numbers, number)
	}
	return numbers, nil
}

// ReadInputTokens reads input from the interactive console and returns a
// slice of tokens.
func (a *Agent) ReadInputTokens() []string {
	if !a.scanner.Scan() {
		if err := a.scanner.Err(); err != nil {
			fmt.Println("Input error or interrupted:", err)
		} else {
			fmt.Println("EOF received. Exiting.")
		}
		os.Exit(0)
	}
	parts := strings.Fields(a.scanner.Text())
	var tokens []string
	for _, part := range parts {
		for _, s := range strings.Split(part, ",") {
			tokens = append(tokens, strings.TrimSpace(s))
		}
	}
	return tokens
}

// ReadInput from the interactive console. Trim whitespace and return only the first token.
// func (a *Agent) ReadInput() (action, target string) {
func (a *Agent) ReadInput() string {
	if !a.scanner.Scan() {
		if err := a.scanner.Err(); err != nil {
			fmt.Println("Input error or interrupted:", err)
		} else {
			// TODO: can we move the signal handling to main?
			fmt.Println("EOF received. Exiting.")
		}
		os.Exit(0)
	}
	input := strings.TrimSpace(a.scanner.Text())
	return input
	/*
		switch len(parts) {
		case 0:
			return "", ""
		case 1:
			return parts[0], ""
		default:
			return parts[0], strings.Join(parts[1:], " ")
		}
	*/
}

// ReadInputConfirm reads input from the interactive console and returns a
// boolean.
func (a *Agent) ReadInputConfirm() (bool, error) {
	input := a.ReadInput()
	accept := strings.ToLower(input)
	if accept == "y" || accept == "yes" {
		return true, nil
	}
	if accept == "n" || accept == "no" {
		return false, nil
	}
	return false, fmt.Errorf("invalid input: %s", accept)
}

func (a Agent) ReadInputNumber(max int) (int, error) {
	input := a.ReadInput()
	number, err := strconv.Atoi(input)
	if err != nil {
		return -1, fmt.Errorf("failed to parse number: %w", err)
	}
	if number < 0 {
		return -1, fmt.Errorf("number less than zero")
	}
	if max == -1 {
		return number, nil
	}
	// TODO: 0 based arrays... make this less confusing
	if number > max {
		return -1, fmt.Errorf("number greater than max")
	}
	return number, nil
}

func PrintCommands(cheatsEnabled bool) {
	var commands []string
	var cheats []string
	for name, command := range actionparser.Commands {
		if command.Cheat {
			cheats = append(cheats, name)
			continue
		}
		commands = append(commands, name)
	}
	sort.Strings(commands)
	sort.Strings(cheats)
	fmt.Println("Available commands:", strings.Join(commands, ", "))
	if cheatsEnabled {
		fmt.Println("Available cheats:", strings.Join(cheats, ", "))
	}
}

func PrintHelp(cheatsEnabled bool) {
	var commands []string
	var cheats []string
	var aliases []string
	for name, command := range actionparser.Commands {
		if command.Cheat {
			cheats = append(cheats, fmt.Sprintf("%s :: %s", name, command.Description))
			continue
		}
		commands = append(commands, fmt.Sprintf("%s :: %s", name, command.Description))
	}
	sort.Strings(aliases)
	sort.Strings(commands)
	sort.Strings(cheats)
	fmt.Println("Available actions:")
	for _, command := range commands {
		fmt.Println(command)
	}
	if cheatsEnabled {
		fmt.Println("Available cheats:")
		for _, cheat := range cheats {
			fmt.Println(cheat)
		}
	}
}
