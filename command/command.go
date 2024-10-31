package command

import (
	"errors"
	"github.com/oceanspy/go-cli-starter/message"
	"os"
	"slices"
	"strings"
)

type Command struct {
	Raw        []string
	Options    map[string]string
	Parameters []string
}

type Program struct {
	Args             []string
	OptionsNames     map[string]string
	OptionsWithValue []string
}

func Construct(args []string, optionsNames map[string]string, optionsWithValue []string) (*Command, *Program) {
	p := Program{Args: args, OptionsNames: optionsNames, OptionsWithValue: optionsWithValue}
	cmd := Command{Raw: args}
	p.extractCommandsAndOptions(&cmd)
	return &cmd, &p
}

func (p *Program) ShortToLongOptionName(arg string) (string, error) {
	for i, v := range p.OptionsNames {
		if v == arg {
			return i, nil
		}
	}
	return "", errors.New("Option not found")
}

func (c *Command) Count() int {
	return len(c.Parameters)
}

func (c *Command) GetAt(at int) (string, error) {
	if len(c.Parameters) > at {
		return c.Parameters[at], nil
	}

	return "", errors.New("No argument found at: " + string(at))
}

func (c *Command) Get() []string {
	return c.Parameters
}

func (c *Command) GetRaw() []string {
	return c.Raw
}

func (c *Command) HasOption(name string) bool {
	if _, exists := c.Options[name]; exists {
		return true
	}

	return false
}

func (c *Command) GetOption(name string) (string, error) {
	if _, exists := c.Options[name]; exists {
		return c.Options[name], nil
	}

	return "", errors.New("Option not found")

}

func (p *Program) extractCommandsAndOptions(c *Command) {
	isOptionValue := false
	for i, arg := range c.Raw {
		if isOptionValue {
			isOptionValue = false
			continue
		}

		// If string starts with --
		if strings.HasPrefix(arg, "--") {
			optionName := strings.TrimPrefix(arg, "--")

			if _, ok := p.OptionsNames[optionName]; !ok {
				message.Error("Couldn't parse the command")
				message.Info("Option", arg, " is not valid")
				os.Exit(1)
			}

			optionValue := ""
			if slices.Contains(p.OptionsWithValue, optionName) {
				optionValue = c.Raw[i+1]
				isOptionValue = true
			}

			if c.Options == nil {
				c.Options = make(map[string]string)
			}
			c.Options[optionName] = optionValue

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

			if _, ok := p.OptionsNames[optionName]; !ok {
				message.Error("Couldn't parse the command")
				message.Info("Option", arg, " is not valid")
				os.Exit(1)
			}

			optionValue := ""
			if slices.Contains(p.OptionsWithValue, optionName) {
				optionValue = c.Raw[i+1]
				isOptionValue = true
			}

			if c.Options == nil {
				c.Options = make(map[string]string)
			}
			c.Options[optionName] = optionValue

			continue
		}

		// Else
		c.Parameters = append(c.Parameters, arg)
	}
}
