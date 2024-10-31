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
		t.Errorf("Last argument doesn't match -- got: %s, expected: %s", result, "a b c d")
	}
}
