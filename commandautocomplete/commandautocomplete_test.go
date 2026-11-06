package commandautocomplete

import (
	"testing"
)

var cmdAutocompleteYaml = `
psw:
encode:
decode:
backup:
  all:
  home:
  system:
  applications:
touchpad:
onscreen-keyboard:
clean-unused-kernels:
connect:
ssh-tunnel:
test:
  a:
    wuj:
    xaaua:
  b:
  c:
  d:
    y:
    z:
`

func TestGetLastArg(t *testing.T) {
	var osArgs = []string{"test", "a"}
	c := Init(cmdAutocompleteYaml, osArgs)

	result := c.GetLastArg()
	if result != "a" {
		t.Errorf("Last argument doesn't match -- got: %s, expected: %s", result, "a")
	}

	osArgs = []string{"test", "a", "wuj"}
	c = Init(cmdAutocompleteYaml, osArgs)

	result = c.GetLastArg()
	if result != "wuj" {
		t.Errorf("Last argument doesn't match -- got: %s, expected: %s", result, "wuj")
	}

	osArgs = []string{"test", "a", "wu"}
	c = Init(cmdAutocompleteYaml, osArgs)

	result = c.GetLastArg()
	if result != "wu" {
		t.Errorf("Last argument doesn't match -- got: %s, expected: %s", result, "wu")
	}

	osArgs = []string{"test", "-h", "wu"}
	c = Init(cmdAutocompleteYaml, osArgs)

	result = c.GetLastArg()
	if result != "wu" {
		t.Errorf("Last argument doesn't match -- got: %s, expected: %s", result, "wu")
	}

	osArgs = []string{"test", "-h"}
	c = Init(cmdAutocompleteYaml, osArgs)

	result = c.GetLastArg()
	if result != "-h" {
		t.Errorf("Last argument doesn't match -- got: %s, expected: %s", result, "-h")
	}
}

func TestGet(t *testing.T) {
	var osArgs = []string{"myprog", "autocomplete", "test", "a"}
	c := Init(cmdAutocompleteYaml, osArgs)

	result := c.Get()
	if result != "wuj xaaua" {
		t.Errorf("Last argument doesn't match -- got: %s, expected: %s", result, "wuj xaaua")
	}

	osArgs = []string{"myprog", "autocomplete", "test", "a", "wuj"}
	c = Init(cmdAutocompleteYaml, osArgs)

	result = c.Get()
	if result != "" {
		t.Errorf("Last argument doesn't match -- got: %s, expected: %s", result, "")
	}

	osArgs = []string{"myprog", "autocomplete", "test", "a", "wu"}
	c = Init(cmdAutocompleteYaml, osArgs)

	result = c.Get()
	if result != "wuj" {
		t.Errorf("Last argument doesn't match -- got: %s, expected: %s", result, "wuj")
	}

	osArgs = []string{"myprog", "autocomplete", "test"}
	c = Init(cmdAutocompleteYaml, osArgs)

	result = c.Get()
	if result != "a b c d" {
		t.Errorf("Last argument doesn't match -- got: %s, expected: %s", result, "a b c d")
	}

	osArgs = []string{"myprog", "autocomplete", "backup"}
	c = Init(cmdAutocompleteYaml, osArgs)

	result = c.Get()
	if result != "all home system applications" {
		t.Errorf("Last argument doesn't match -- got: %s, expected: %s", result, "all home system applications")
	}
}

func TestGetOrder(t *testing.T) {
	// Suggestions must follow the YAML order, every time.
	for i := 0; i < 20; i++ {
		c := Init(cmdAutocompleteYaml, []string{"myprog", "autocomplete"})

		expected := "psw encode decode backup touchpad onscreen-keyboard clean-unused-kernels connect ssh-tunnel test"
		if result := c.Get(); result != expected {
			t.Fatalf("Get() = %q, expected %q", result, expected)
		}
	}
}

func TestGetCases(t *testing.T) {
	yamlStr := `
copy:
  _file:
backup:
  home:
  system:
  hosts:
open:
  _file:
  recent:
base: &base
  one:
  two:
alias: *base
send:
  _file
read: _file
invoice:
  send:
    _file
other: value
`

	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{name: "top level", args: []string{"myprog", "autocomplete"}, expected: "copy backup open base alias send read invoice other"},
		{name: "prefix at top level", args: []string{"myprog", "autocomplete", "b"}, expected: "backup base"},
		{name: "prefix keeps yaml order", args: []string{"myprog", "autocomplete", "backup", "h"}, expected: "home hosts"},
		{name: "no match for prefix", args: []string{"myprog", "autocomplete", "backup", "x"}, expected: ""},
		{name: "unknown command in the middle", args: []string{"myprog", "autocomplete", "unknown", "h"}, expected: ""},
		{name: "leaf command", args: []string{"myprog", "autocomplete", "backup", "home"}, expected: ""},
		{name: "after a leaf command", args: []string{"myprog", "autocomplete", "backup", "home", "h"}, expected: ""},
		{name: "file completion level", args: []string{"myprog", "autocomplete", "copy"}, expected: "_file"},
		{name: "file completion while typing", args: []string{"myprog", "autocomplete", "copy", "some/pa"}, expected: "_file"},
		{name: "file completion wins over other entries", args: []string{"myprog", "autocomplete", "open", "rec"}, expected: "_file"},
		{name: "yaml alias", args: []string{"myprog", "autocomplete", "alias"}, expected: "one two"},
		{name: "file shorthand on next line", args: []string{"myprog", "autocomplete", "send"}, expected: "_file"},
		{name: "file shorthand on next line while typing", args: []string{"myprog", "autocomplete", "send", "inv"}, expected: "_file"},
		{name: "file shorthand on same line", args: []string{"myprog", "autocomplete", "read"}, expected: "_file"},
		{name: "file shorthand on same line while typing", args: []string{"myprog", "autocomplete", "read", "some/pa"}, expected: "_file"},
		{name: "file shorthand nested", args: []string{"myprog", "autocomplete", "invoice", "send", "inv"}, expected: "_file"},
		{name: "other scalar value is a leaf", args: []string{"myprog", "autocomplete", "other"}, expected: ""},
		{name: "after a file argument", args: []string{"myprog", "autocomplete", "send", "file1", "x"}, expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Init(yamlStr, tt.args)
			if result := c.Get(); result != tt.expected {
				t.Errorf("Get() with args %v = %q, expected %q", tt.args, result, tt.expected)
			}
		})
	}
}

func TestGetEmptyYaml(t *testing.T) {
	c := Init("", []string{"myprog", "autocomplete", "a"})
	if result := c.Get(); result != "" {
		t.Errorf("Get() = %q, expected empty", result)
	}
}

func TestInitInvalidYamlPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("Init should panic on invalid YAML")
		}
	}()

	Init("a: [unclosed", []string{"myprog", "autocomplete"})
}

func TestInitKeepsParsedCommands(t *testing.T) {
	c := Init(cmdAutocompleteYaml, nil)

	commands, ok := c.Commands.(map[string]interface{})
	if !ok {
		t.Fatalf("Commands should be a map, got %T", c.Commands)
	}
	if _, exists := commands["backup"]; !exists {
		t.Error("Commands should contain the backup key")
	}
}

func TestGetPenultimateArg(t *testing.T) {
	tests := []struct {
		args     []string
		expected string
	}{
		{args: nil, expected: ""},
		{args: []string{"a"}, expected: ""},
		{args: []string{"a", "b"}, expected: "a"},
		{args: []string{"a", "b", "c"}, expected: "b"},
	}

	for _, tt := range tests {
		c := Init(cmdAutocompleteYaml, tt.args)
		if result := c.GetPenultimateArg(); result != tt.expected {
			t.Errorf("GetPenultimateArg() with args %v = %q, expected %q", tt.args, result, tt.expected)
		}
	}
}

func TestGetLastArgEmpty(t *testing.T) {
	c := Init(cmdAutocompleteYaml, nil)
	if result := c.GetLastArg(); result != "" {
		t.Errorf("GetLastArg() = %q, expected empty", result)
	}
}
