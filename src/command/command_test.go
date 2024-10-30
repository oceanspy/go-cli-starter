package command

import (
	"testing"
)

type TestElement struct {
	value    string
	expected string
}

func TestShortToLongOptionName(t *testing.T) {
	var optionsNames = map[string]string{
		"help": "h",
		"test": "t",
	}
	var optionsWithValue []string
	var args = []string{
		"command1",
		"command2",
	}

	testsElements := []TestElement{
		{value: "h", expected: "help"},
		{value: "t", expected: "test"},
		{value: "w", expected: ""},
	}

	_, p := Construct(args, optionsNames, optionsWithValue)

	for _, e := range testsElements {
		result, err := p.ShortToLongOptionName(e.value)
		if err != nil && e.expected != "" {
			t.Errorf("Result %s is throwing an error while an actual value is expected: %s", result, e.expected)
			continue
		}

		if result != e.expected {
			t.Errorf("Result %s is not what is expected: %s", result, e.expected)
		}
	}

}

func TestLen(t *testing.T) {
	var optionsNames = map[string]string{
		"help": "h",
		"test": "t",
	}
	var optionsWithValue []string
	var args = []string{
		"command1",
		"command2",
	}

	c, _ := Construct(args, optionsNames, optionsWithValue)

	result := c.Len()
	if result != 2 {
		t.Errorf("Command length is: %d ; expected %d", result, 2)
	}

	var optionsNames2 = map[string]string{
		"help":  "h",
		"test":  "t",
		"chose": "c",
	}
	var optionsWithValue2 = []string{
		"chose",
	}
	var args2 = []string{
		"command1",
		"command2",
		"--chose",
		"testvalue",
	}

	c2, _ := Construct(args2, optionsNames2, optionsWithValue2)

	result2 := c2.Len()
	if result2 != 2 {
		t.Errorf("Command length is: %d ; expected %d", result2, 2)
	}

	var optionsNames3 = map[string]string{
		"help":  "h",
		"test":  "t",
		"chose": "c",
	}
	var optionsWithValue3 = []string{
		"chose",
	}
	var args3 = []string{
		"command1",
		"command2",
		"-c",
		"testvalue",
		"--test",
		"-t",
	}

	c3, _ := Construct(args3, optionsNames3, optionsWithValue3)

	result3 := c3.Len()
	if result3 != 2 {
		t.Errorf("Command length is: %d ; expected %d", result3, 2)
	}
}

func TestHasOption(t *testing.T) {
	var optionsNames = map[string]string{
		"help": "h",
		"test": "t",
	}
	var optionsWithValue []string
	var args = []string{
		"command1",
		"command2",
	}

	c, _ := Construct(args, optionsNames, optionsWithValue)

	result := c.HasOption("test")
	if result != false {
		t.Errorf("Command doesn't have test option but it returns true")
	}

	var optionsNames2 = map[string]string{
		"help":  "h",
		"test":  "t",
		"chose": "c",
	}
	var optionsWithValue2 = []string{
		"chose",
	}
	var args2 = []string{
		"command1",
		"--chose",
		"testvalue",
		"command2",
	}

	c2, _ := Construct(args2, optionsNames2, optionsWithValue2)

	result2 := c2.HasOption("test")
	if result2 != false {
		t.Errorf("Command doesn't have test option but it returns true")
	}

	result2b := c2.HasOption("chose")
	if result2b != true {
		t.Errorf("Command have chose option but it returns false")
	}

	var optionsNames3 = map[string]string{
		"help":  "h",
		"test":  "t",
		"chose": "c",
	}
	var optionsWithValue3 = []string{
		"chose",
	}
	var args3 = []string{
		"-c",
		"testvalue",
		"--test",
		"-h",
		"command1",
		"command2",
	}

	c3, _ := Construct(args3, optionsNames3, optionsWithValue3)

	result3 := c3.HasOption("chose")
	if result3 != true {
		t.Errorf("Command have chose option but it returns false")
	}

	result3b := c3.HasOption("test")
	if result3b != true {
		t.Errorf("Command have chose option but it returns false")
	}

	result3c := c3.HasOption("help")
	if result3c != true {
		t.Errorf("Command have chose option but it returns false")
	}

	result3d := c3.HasOption("h")
	if result3d != false {
		t.Errorf("Command doesn't have 'h' option but it returns true. Only long name options are accepted")
	}
}

func TestGet(t *testing.T) {
	var optionsNames = map[string]string{
		"help": "h",
		"test": "t",
	}
	var optionsWithValue []string
	var args = []string{
		"command1",
		"command2",
	}

	c, _ := Construct(args, optionsNames, optionsWithValue)

	result := c.Get(0)
	if result != "command1" {
		t.Errorf("Command 0 is: %s ; expected %s", result, "command1")
	}

	var optionsNames2 = map[string]string{
		"help":  "h",
		"test":  "t",
		"chose": "c",
	}
	var optionsWithValue2 = []string{
		"chose",
	}
	var args2 = []string{
		"command1",
		"--chose",
		"testvalue",
		"command2",
	}

	c2, _ := Construct(args2, optionsNames2, optionsWithValue2)

	result2 := c2.Get(1)
	if result2 != "command2" {
		t.Errorf("Command 1 is: %s ; expected %s", result, "command2")
	}

	var optionsNames3 = map[string]string{
		"help":  "h",
		"test":  "t",
		"chose": "c",
	}
	var optionsWithValue3 = []string{
		"chose",
	}
	var args3 = []string{
		"-c",
		"testvalue",
		"--test",
		"-t",
		"command1",
		"command2",
	}

	c3, _ := Construct(args3, optionsNames3, optionsWithValue3)

	result3 := c3.Get(1)
	if result3 != "command2" {
		t.Errorf("Command 1 is: %s ; expected %s", result, "command2")
	}
}
