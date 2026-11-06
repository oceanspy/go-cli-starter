package commandautocomplete

import (
	"reflect"
	"testing"
)

var completeYaml = `
invoice:
  list:
    to:
    from:
  send:
    _file
  generate:
    pdf:
      _file:
quote:
  send:
  sender:
my item:
copy: _file
`

var completeOptionsNames = map[string]string{
	"help":  "h",
	"limit": "l",
	"all":   "",
}

var completeOptionsWithValue = []string{"limit"}

func complete(words ...string) []string {
	args := append([]string{"prog", "commands", "--"}, words...)
	return Init(completeYaml, args).WithOptions(completeOptionsNames, completeOptionsWithValue).Complete()
}

func TestComplete(t *testing.T) {
	tests := []struct {
		name     string
		words    []string
		expected []string
	}{
		// Point 1: the last word is always the word being typed
		{name: "nothing typed", words: []string{""}, expected: []string{"invoice", "quote", "my item", "copy"}},
		{name: "partial word", words: []string{"inv"}, expected: []string{"invoice"}},
		{name: "complete word without space", words: []string{"invoice"}, expected: []string{"invoice"}},
		{name: "complete word with space", words: []string{"invoice", ""}, expected: []string{"list", "send", "generate"}},
		{name: "partial sub command", words: []string{"invoice", "l"}, expected: []string{"list"}},
		{name: "final command without space", words: []string{"invoice", "list", "to"}, expected: []string{"to"}},
		{name: "after final command", words: []string{"invoice", "list", "to", ""}, expected: nil},
		{name: "command that starts another one", words: []string{"quote", "send"}, expected: []string{"send", "sender"}},
		{name: "no match", words: []string{"x"}, expected: nil},
		{name: "unknown word in the middle", words: []string{"unknown", ""}, expected: nil},
		{name: "case sensitive", words: []string{"INV"}, expected: nil},

		// Point 3: suggestions can contain spaces
		{name: "suggestion with space", words: []string{"my"}, expected: []string{"my item"}},
		{name: "after a word with space", words: []string{"my item", ""}, expected: nil},

		// Files
		{name: "file level", words: []string{"invoice", "generate", "pdf", ""}, expected: []string{"_file"}},
		{name: "file level while typing", words: []string{"invoice", "generate", "pdf", "~/Doc"}, expected: []string{"_file"}},
		{name: "file shorthand", words: []string{"invoice", "send", "inv"}, expected: []string{"_file"}},
		{name: "file shorthand on same line", words: []string{"copy", ""}, expected: []string{"_file"}},
		{name: "after a file", words: []string{"copy", "a.txt", ""}, expected: nil},

		// Options
		{name: "dash", words: []string{"-"}, expected: []string{"--all", "--help", "-h", "--limit", "-l"}},
		{name: "double dash", words: []string{"--"}, expected: []string{"--all", "--help", "--limit"}},
		{name: "partial long option", words: []string{"--li"}, expected: []string{"--limit"}},
		{name: "complete long option", words: []string{"--limit"}, expected: []string{"--limit"}},
		{name: "options after commands", words: []string{"invoice", "list", "--h"}, expected: []string{"--help"}},
		{name: "used options are not suggested", words: []string{"--help", "-l", "5", "-"}, expected: []string{"--all"}},
		{name: "options are skipped", words: []string{"-h", "invoice", "--all", ""}, expected: []string{"list", "send", "generate"}},
		{name: "option value is skipped", words: []string{"--limit", "5", "invoice", ""}, expected: []string{"list", "send", "generate"}},
		{name: "short option value is skipped", words: []string{"invoice", "-l", "list", ""}, expected: []string{"list", "send", "generate"}},
		{name: "option value being typed", words: []string{"invoice", "--limit", ""}, expected: nil},
		{name: "unknown option is skipped", words: []string{"--nope", "invoice", ""}, expected: []string{"list", "send", "generate"}},
		{name: "after end of options", words: []string{"--", "-"}, expected: nil},
		{name: "commands after end of options", words: []string{"--", "invoice", ""}, expected: []string{"list", "send", "generate"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := complete(tt.words...); !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Complete(%q) = %q, expected %q", tt.words, result, tt.expected)
			}
		})
	}
}

func TestCompleteWithoutSeparator(t *testing.T) {
	// The program may pass the words without the "--" separator
	c := Init(completeYaml, []string{"prog", "commands", "invoice", ""})
	if result := c.Complete(); !reflect.DeepEqual(result, []string{"list", "send", "generate"}) {
		t.Errorf("Complete() = %q, expected [list send generate]", result)
	}
}

func TestCompleteNoWords(t *testing.T) {
	for _, args := range [][]string{nil, {"prog"}, {"prog", "commands"}, {"prog", "commands", "--"}} {
		c := Init(completeYaml, args)
		if result := c.Complete(); !reflect.DeepEqual(result, []string{"invoice", "quote", "my item", "copy"}) {
			t.Errorf("Complete() with args %q = %q, expected the top level", args, result)
		}
	}
}

func TestCompleteWithoutOptions(t *testing.T) {
	c := Init(completeYaml, []string{"prog", "commands", "--", "-"})
	if result := c.Complete(); result != nil {
		t.Errorf("Complete() = %q, expected no suggestion without options", result)
	}
}

func TestCompleteEmptyYaml(t *testing.T) {
	c := Init("", []string{"prog", "commands", "--", ""})
	if result := c.Complete(); result != nil {
		t.Errorf("Complete() = %q, expected nil", result)
	}
}
