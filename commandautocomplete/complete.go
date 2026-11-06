package commandautocomplete

import (
	"slices"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// WithOptions sets the options of the program, so Complete can suggest them
// and skip them (and their values) while reading the command.
// The parameters are the same as for command.Construct.
func (c *Commands) WithOptions(optionsNames map[string]string, optionsWithValue []string) *Commands {
	c.optionsNames = optionsNames
	c.optionsWithValue = optionsWithValue
	return c
}

// Complete returns the suggestions for the word being typed.
//
// Unlike Get, the last argument is always the word being typed, even when it
// is empty: "prog commands invoice" suggests "invoice" (still typing it) and
// "prog commands invoice <empty word>" suggests what comes after "invoice".
//
// Suggestions are returned in the YAML order. When the word is a file path,
// the only suggestion is SYSTEM_AUTOCOMPLETE ("_file").
// When the word starts with "-", the program options are suggested.
func (c *Commands) Complete() []string {
	// We skip the two first arguments
	// 0 = main program
	// 1 = arg to call autocomplete
	var args []string
	if len(c.OsArgs) > 2 {
		args = c.OsArgs[2:]
	}
	// "--" separates the autocomplete call from the words to complete
	if len(args) > 0 && args[0] == "--" {
		args = args[1:]
	}

	current := ""
	if len(args) > 0 {
		current = args[len(args)-1]
		args = args[:len(args)-1]
	}

	level := c.root
	usedOptions := map[string]bool{}
	endOfOptions := false
	for i := 0; i < len(args); i++ {
		arg := args[i]

		if !endOfOptions && arg == "--" {
			endOfOptions = true
			continue
		}

		if !endOfOptions && isOption(arg) {
			name := c.longOptionName(arg)
			usedOptions[name] = true
			if slices.Contains(c.optionsWithValue, name) {
				if i+1 == len(args) {
					// The word being typed is the value of this option
					return nil
				}
				i++
			}
			continue
		}

		value, exists := lookup(level, arg)
		next := resolve(value)
		if isFileShorthand(next) {
			next = fileLevel
		}

		switch {
		case exists && next != nil && next.Kind == yaml.MappingNode:
			level = next
		default:
			// A final command, a file path or an unknown word: nothing comes after it
			level = nil
		}
	}

	if !endOfOptions && strings.HasPrefix(current, "-") {
		return c.completeOptions(current, usedOptions)
	}

	if _, hasFile := lookup(level, SYSTEM_AUTOCOMPLETE); hasFile {
		return []string{SYSTEM_AUTOCOMPLETE}
	}

	var suggestions []string
	for _, key := range keys(level) {
		if strings.HasPrefix(key, current) {
			suggestions = append(suggestions, key)
		}
	}

	return suggestions
}

func isOption(arg string) bool {
	return len(arg) > 1 && strings.HasPrefix(arg, "-")
}

// longOptionName returns the long name of an option ("--limit" or "-l" gives "limit").
func (c *Commands) longOptionName(arg string) string {
	if strings.HasPrefix(arg, "--") {
		return strings.TrimPrefix(arg, "--")
	}

	short := strings.TrimPrefix(arg, "-")
	for long, s := range c.optionsNames {
		if s == short {
			return long
		}
	}

	return short
}

// completeOptions suggests the options (long, then short name) not used yet,
// sorted by long name.
func (c *Commands) completeOptions(current string, usedOptions map[string]bool) []string {
	names := make([]string, 0, len(c.optionsNames))
	for long := range c.optionsNames {
		if !usedOptions[long] {
			names = append(names, long)
		}
	}
	sort.Strings(names)

	var suggestions []string
	for _, long := range names {
		candidates := []string{"--" + long}
		if short := c.optionsNames[long]; short != "" {
			candidates = append(candidates, "-"+short)
		}

		for _, candidate := range candidates {
			if strings.HasPrefix(candidate, current) {
				suggestions = append(suggestions, candidate)
			}
		}
	}

	return suggestions
}
