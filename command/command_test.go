package command

import (
	"errors"
	"os"
	"os/exec"
	"reflect"
	"strings"
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

func TestCount(t *testing.T) {
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

	result := c.Count()
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

	result2 := c2.Count()
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

	result3 := c3.Count()
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

func TestGetAt(t *testing.T) {
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

	result, _ := c.GetAt(0)
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

	result2, _ := c2.GetAt(1)
	if result2 != "command2" {
		t.Errorf("Command 1 is: %s ; expected %s", result2, "command2")
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

	result3, _ := c3.GetAt(1)
	if result3 != "command2" {
		t.Errorf("Command 1 is: %s ; expected %s", result3, "command2")
	}
}

func newTestCommand(args []string) (*Command, *Program) {
	var optionsNames = map[string]string{
		"help":  "h",
		"test":  "t",
		"chose": "c",
	}
	var optionsWithValue = []string{
		"chose",
	}

	return Construct(args, optionsNames, optionsWithValue)
}

func TestGetAtOutOfRange(t *testing.T) {
	c, _ := newTestCommand([]string{"command1", "command2"})

	for _, at := range []int{2, 10} {
		result, err := c.GetAt(at)
		if err == nil {
			t.Errorf("GetAt(%d) should return an error, got %q", at, result)
		}
		if result != "" {
			t.Errorf("GetAt(%d) = %q, expected empty", at, result)
		}
	}
}

func TestGetAfter(t *testing.T) {
	c, _ := newTestCommand([]string{"command1", "--chose", "value", "command2", "-t", "command3"})

	tests := []struct {
		nb       int
		expected []string
		wantErr  bool
	}{
		{nb: 0, expected: []string{"command1", "command2", "command3"}},
		{nb: 1, expected: []string{"command2", "command3"}},
		{nb: 2, expected: []string{"command3"}},
		{nb: 3, wantErr: true},
		{nb: 10, wantErr: true},
	}

	for _, tt := range tests {
		result, err := c.GetAfter(tt.nb)
		if (err != nil) != tt.wantErr {
			t.Errorf("GetAfter(%d) error = %v, wantErr %v", tt.nb, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf("GetAfter(%d) = %#v, expected %#v", tt.nb, result, tt.expected)
		}
	}
}

func TestGetAndGetRaw(t *testing.T) {
	args := []string{"command1", "-c", "value", "--test", "command2"}
	c, _ := newTestCommand(args)

	if result := c.Get(); !reflect.DeepEqual(result, []string{"command1", "command2"}) {
		t.Errorf("Get() = %#v, expected [command1 command2]", result)
	}

	if result := c.GetRaw(); !reflect.DeepEqual(result, args) {
		t.Errorf("GetRaw() = %#v, expected %#v", result, args)
	}
}

func TestNoParameters(t *testing.T) {
	c, _ := newTestCommand([]string{"--test"})

	if c.Count() != 0 {
		t.Errorf("Count() = %d, expected 0", c.Count())
	}
	if _, err := c.GetAt(0); err == nil {
		t.Error("GetAt(0) should return an error when there are no parameters")
	}
	if _, err := c.GetAfter(0); err == nil {
		t.Error("GetAfter(0) should return an error when there are no parameters")
	}
}

func TestGetOption(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		option   string
		expected string
		wantErr  bool
	}{
		{name: "long option with value", args: []string{"cmd", "--chose", "value"}, option: "chose", expected: "value"},
		{name: "short option with value", args: []string{"cmd", "-c", "value"}, option: "chose", expected: "value"},
		{name: "value looking like a parameter", args: []string{"--chose", "cmd", "other"}, option: "chose", expected: "cmd"},
		{name: "option without value", args: []string{"cmd", "--test"}, option: "test", expected: ""},
		{name: "short option without value", args: []string{"cmd", "-t"}, option: "test", expected: ""},
		{name: "option not given", args: []string{"cmd"}, option: "test", wantErr: true},
		{name: "short name is not an option name", args: []string{"cmd", "-t"}, option: "t", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := newTestCommand(tt.args)

			result, err := c.GetOption(tt.option)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetOption(%q) error = %v, wantErr %v", tt.option, err, tt.wantErr)
			}
			if result != tt.expected {
				t.Errorf("GetOption(%q) = %q, expected %q", tt.option, result, tt.expected)
			}
		})
	}
}

// TestInvalidCommandExits checks that a command that can't be parsed stops
// the program. os.Exit can't be caught, so the parsing runs in a child test process.
func TestInvalidCommandExits(t *testing.T) {
	if args := os.Getenv("GO_CLI_STARTER_INVALID_COMMAND"); args != "" {
		newTestCommand(strings.Fields(args))
		return
	}

	tests := []struct {
		args     string
		expected string
	}{
		{args: "cmd --unknown", expected: "is not valid"},
		{args: "cmd -u", expected: "is not valid"},
		{args: "cmd --chose", expected: "requires a value"},
		{args: "cmd -c", expected: "requires a value"},
	}

	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			cmd := exec.Command(os.Args[0], "-test.run=^TestInvalidCommandExits$")
			cmd.Env = append(os.Environ(), "GO_CLI_STARTER_INVALID_COMMAND="+tt.args)
			out, err := cmd.CombinedOutput()

			var exitErr *exec.ExitError
			if !errors.As(err, &exitErr) || exitErr.ExitCode() != 1 {
				t.Fatalf("expected exit code 1, got error %v\noutput: %s", err, out)
			}
			if !strings.Contains(string(out), "Couldn't parse the command") || !strings.Contains(string(out), tt.expected) {
				t.Errorf("expected a parse error message containing %q, got: %s", tt.expected, out)
			}
		})
	}
}

func TestEndOfOptions(t *testing.T) {
	c, _ := newTestCommand([]string{"prog", "commands", "--", "invoice", "--chose", "-t", "--unknown", "-", "--", ""})

	expected := []string{"prog", "commands", "invoice", "--chose", "-t", "--unknown", "-", "--", ""}
	if result := c.Get(); !reflect.DeepEqual(result, expected) {
		t.Errorf("Get() = %#v, expected %#v", result, expected)
	}
	if c.HasOption("chose") || c.HasOption("test") {
		t.Error("options after -- should not be parsed")
	}

	c, _ = newTestCommand([]string{"prog", "-t", "--", "-c"})
	if !c.HasOption("test") {
		t.Error("options before -- should still be parsed")
	}
	if result := c.Get(); !reflect.DeepEqual(result, []string{"prog", "-c"}) {
		t.Errorf("Get() = %#v, expected [prog -c]", result)
	}
}
