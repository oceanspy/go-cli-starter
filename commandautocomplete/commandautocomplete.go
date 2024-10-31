package commandautocomplete

import (
	"errors"
	"gopkg.in/yaml.v3"
	"strings"
)

type Commands struct {
	CommandsYamlStr string
	OsArgs          []string
	Commands        interface{}
}

func Init(commandsYamlStr string, osArgs []string) *Commands {
	c := Commands{CommandsYamlStr: commandsYamlStr, OsArgs: osArgs}
	c.parseYaml()
	return &c
}

func (c *Commands) Get() string {
	nextCmds := ""
	lastArg := c.GetLastArg()

	// For each string, we get the current level and nexts
	var allNextLevels = c.Commands.(map[string]interface{})
	var nothingFoundForLastArg = false
	for i, cmd := range c.OsArgs {
		// We skip the two first arguments
		// 0 = main program
		// 1 = arg to call autocomplete
		if i < 2 {
			continue
		}

		nextLevels, err := c.getLevel(allNextLevels, cmd)
		if err != nil {
			if len(c.OsArgs) == i+1 {
				if _, exists := allNextLevels[cmd]; exists {
					return ""
				}

				nothingFoundForLastArg = true
				break
			}

			return ""
		}

		allNextLevels = nextLevels
	}

	// Since map don't keep the order => we transform the map to a slice of struct key/value
	var nextLevel []string
	for _, kv := range c.mapToKeyValue(allNextLevels) {
		nextLevel = append(nextLevel, kv.key)
	}

	// if current word is the beginning of one or more words of the level, we suggest these one
	if nothingFoundForLastArg {
		var autocompleteLastWord []string
		for _, completeLastWord := range nextLevel {
			if strings.HasPrefix(completeLastWord, lastArg) {
				autocompleteLastWord = append(autocompleteLastWord, completeLastWord)
			}
		}
		nextCmds = strings.Join(autocompleteLastWord, " ")
		return nextCmds
	}

	nextCmds = strings.Join(nextLevel, " ")

	return nextCmds
}

func (c *Commands) getLevel(allNextLevels map[string]interface{}, cmd string) (map[string]interface{}, error) {
	nextLevels, ok := allNextLevels[cmd].(map[string]interface{})
	if !ok {
		return nextLevels, errors.New("Couldn't find any match -- is Yaml correct ?")
	}

	return nextLevels, nil
}

func (c *Commands) GetLastArg() string {
	if len(c.OsArgs) == 0 {
		return ""
	}

	return c.OsArgs[len(c.OsArgs)-1]
}

func (c *Commands) parseYaml() {
	var parsedData interface{}

	err := yaml.Unmarshal([]byte(c.CommandsYamlStr), &parsedData)
	if err != nil {
		panic(err)
	}

	c.Commands = parsedData
}

type kv struct {
	key   string
	value interface{}
}

func (c *Commands) mapToKeyValue(m map[string]interface{}) []kv {
	var kvs []kv
	for k, v := range m {
		kvs = append(kvs, kv{key: k, value: v})
	}
	return kvs
}
