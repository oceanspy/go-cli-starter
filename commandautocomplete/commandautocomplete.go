package commandautocomplete

import (
	"strings"

	"gopkg.in/yaml.v3"
)

const SYSTEM_AUTOCOMPLETE string = "_file"

type Commands struct {
	CommandsYamlStr string
	OsArgs          []string
	Commands        interface{}

	// root keeps the parsed YAML tree so suggestions follow the order in
	// which commands are written (a map would lose that order).
	root *yaml.Node

	// Program options, set with WithOptions
	optionsNames     map[string]string
	optionsWithValue []string
}

func Init(commandsYamlStr string, osArgs []string) *Commands {
	c := Commands{CommandsYamlStr: commandsYamlStr, OsArgs: osArgs}
	c.parseYaml()
	return &c
}

func (c *Commands) Get() string {
	lastArg := c.GetLastArg()

	// For each argument, we move down one level of the command tree
	level := c.root
	nothingFoundForLastArg := false
	for i, cmd := range c.OsArgs {
		// We skip the two first arguments
		// 0 = main program
		// 1 = arg to call autocomplete
		if i < 2 {
			continue
		}

		value, exists := lookup(level, cmd)
		nextLevel := resolve(value)
		if isFileShorthand(nextLevel) {
			nextLevel = fileLevel
		}
		if nextLevel == nil || nextLevel.Kind != yaml.MappingNode {
			if len(c.OsArgs) == i+1 {
				if exists {
					return ""
				}

				nothingFoundForLastArg = true
				break
			}

			return ""
		}

		level = nextLevel
	}

	nextLevel := keys(level)

	// if current word is the beginning of one or more words of the level, we suggest these one
	if nothingFoundForLastArg {
		var autocompleteLastWord []string
		for _, completeLastWord := range nextLevel {
			if completeLastWord == SYSTEM_AUTOCOMPLETE {
				return SYSTEM_AUTOCOMPLETE
			}

			if strings.HasPrefix(completeLastWord, lastArg) {
				autocompleteLastWord = append(autocompleteLastWord, completeLastWord)
			}
		}
		return strings.Join(autocompleteLastWord, " ")
	}

	return strings.Join(nextLevel, " ")
}

func (c *Commands) GetLastArg() string {
	if len(c.OsArgs) == 0 {
		return ""
	}

	return c.OsArgs[len(c.OsArgs)-1]
}

func (c *Commands) GetPenultimateArg() string {
	if len(c.OsArgs) < 2 {
		return ""
	}

	return c.OsArgs[len(c.OsArgs)-2]
}

func (c *Commands) parseYaml() {
	var root yaml.Node
	if err := yaml.Unmarshal([]byte(c.CommandsYamlStr), &root); err != nil {
		panic(err)
	}

	var parsedData interface{}
	if err := root.Decode(&parsedData); err != nil {
		panic(err)
	}

	c.root = resolve(&root)
	c.Commands = parsedData
}

// fileLevel is a level that only offers file completion, like "cmd: {_file: }".
var fileLevel = &yaml.Node{
	Kind: yaml.MappingNode,
	Content: []*yaml.Node{
		{Kind: yaml.ScalarNode, Value: SYSTEM_AUTOCOMPLETE},
		{Kind: yaml.ScalarNode, Tag: "!!null"},
	},
}

// isFileShorthand reports whether a value is "_file" written without a colon,
// as in "cmd: _file" or:
//
//	cmd:
//	  _file
func isFileShorthand(n *yaml.Node) bool {
	return n != nil && n.Kind == yaml.ScalarNode && n.Value == SYSTEM_AUTOCOMPLETE
}

// resolve unwraps document and alias nodes to reach the actual content node.
func resolve(n *yaml.Node) *yaml.Node {
	for n != nil {
		switch {
		case n.Kind == yaml.DocumentNode && len(n.Content) > 0:
			n = n.Content[0]
		case n.Kind == yaml.AliasNode:
			n = n.Alias
		default:
			return n
		}
	}

	return nil
}

// lookup returns the value node of key in a mapping node, and whether the key exists.
func lookup(mapping *yaml.Node, key string) (*yaml.Node, bool) {
	if mapping == nil || mapping.Kind != yaml.MappingNode {
		return nil, false
	}

	for i := 0; i+1 < len(mapping.Content); i += 2 {
		if mapping.Content[i].Value == key {
			return mapping.Content[i+1], true
		}
	}

	return nil, false
}

// keys returns the keys of a mapping node, in the order they are written.
func keys(mapping *yaml.Node) []string {
	if mapping == nil || mapping.Kind != yaml.MappingNode {
		return nil
	}

	var result []string
	for i := 0; i < len(mapping.Content); i += 2 {
		result = append(result, mapping.Content[i].Value)
	}

	return result
}
