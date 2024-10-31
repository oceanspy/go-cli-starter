package command

import (
	"errors"
	"github.com/oceanspy/go-cli-starter/message"
	"os"
	"slices"
	"strings"
)

type Command struct {
	args     []string
	options  map[string]string
	commands []string
}

type Program struct {
	optionsNames     map[string]string
	optionsWithValue []string
}

func Construct(args []string, optionsNames map[string]string, optionsWithValue []string) (*Command, *Program) {
	p := Program{optionsNames: optionsNames, optionsWithValue: optionsWithValue}
	cmd := Command{args: args}
	p.extractCommandsAndOptions(&cmd)
	return &cmd, &p
}

func (p *Program) extractCommandsAndOptions(c *Command) {
	isOptionValue := false
	for i, arg := range c.args {
		if isOptionValue {
			isOptionValue = false
			continue
		}

		// If string starts with --
		if strings.HasPrefix(arg, "--") {
			optionName := strings.TrimPrefix(arg, "--")

			if _, ok := p.optionsNames[optionName]; !ok {
				message.Error("Couldn't parse the command")
				message.Info("Option", arg, " is not valid")
				os.Exit(1)
			}

			optionValue := ""
			if slices.Contains(p.optionsWithValue, optionName) {
				optionValue = c.args[i+1]
				isOptionValue = true
			}

			if c.options == nil {
				c.options = make(map[string]string)
			}
			c.options[optionName] = optionValue

			continue
		}

		// If string starts with -
		if strings.HasPrefix(arg, "-") {
			shortOptionName := strings.TrimPrefix(arg, "-")
			optionName, err := p.ShortToLongOptionName(shortOptionName)
			if err != nil {
				message.Error("Couldn't parse the command")
				message.Info("Option", arg, " is not valid")
				os.Exit(1)
			}

			if _, ok := p.optionsNames[optionName]; !ok {
				message.Error("Couldn't parse the command")
				message.Info("Option", arg, " is not valid")
				os.Exit(1)
			}

			optionValue := ""
			if slices.Contains(p.optionsWithValue, optionName) {
				optionValue = c.args[i+1]
				isOptionValue = true
			}

			if c.options == nil {
				c.options = make(map[string]string)
			}
			c.options[optionName] = optionValue

			continue
		}

		// Else
		c.commands = append(c.commands, arg)
	}
}

func (p *Program) ShortToLongOptionName(arg string) (string, error) {
	for i, v := range p.optionsNames {
		if v == arg {
			return i, nil
		}
	}
	return "", errors.New("Option not found")
}

func (c *Command) Count() int {
	return len(c.commands)
}

func (c *Command) Get(at int) (string, error) {
	if len(c.commands) > at {
		return c.commands[at], nil
	}

	return "", errors.New("No argument found at: " + string(at))
}

func (c *Command) HasOption(name string) bool {
	if _, exists := c.options[name]; exists {
		return true
	}

	return false
}

func (c *Command) GetOption(name string) (string, error) {
	if _, exists := c.options[name]; exists {
		return c.options[name], nil
	}

	return "", errors.New("Option not found")

}
