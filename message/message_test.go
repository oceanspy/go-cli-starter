package message

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/oceanspy/go-cli-starter/color"
)

// captureStdout runs fn and returns everything it printed to stdout.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	original := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = original }()

	done := make(chan string)
	go func() {
		out, _ := io.ReadAll(r)
		done <- string(out)
	}()

	fn()
	w.Close()

	return <-done
}

// withStdin replaces stdin with the given input while fn runs.
func withStdin(t *testing.T, input string, fn func()) {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.WriteString(input); err != nil {
		t.Fatal(err)
	}
	w.Close()

	original := os.Stdin
	os.Stdin = r
	defer func() {
		os.Stdin = original
		r.Close()
	}()

	fn()
}

func TestLn(t *testing.T) {
	if out := captureStdout(t, Ln); out != "\n" {
		t.Errorf("Ln() printed %q, expected a single newline", out)
	}
}

func TestSuccess(t *testing.T) {
	out := captureStdout(t, func() { Success("all ", "good") })
	expected := "✅ " + color.Green + "all good" + color.Reset + "\n"
	if out != expected {
		t.Errorf("Success printed %q, expected %q", out, expected)
	}
}

func TestWarning(t *testing.T) {
	out := captureStdout(t, func() { Warning("be ", "careful") })
	expected := "⚠️ " + color.Yellow + "be careful" + color.Reset + "\n"
	if out != expected {
		t.Errorf("Warning printed %q, expected %q", out, expected)
	}
}

func TestError(t *testing.T) {
	tests := []struct {
		name     string
		args     []interface{}
		expected string
	}{
		{name: "no args", args: nil, expected: "❌\n"},
		{name: "message only", args: []interface{}{"failed"}, expected: "❌ " + color.Red + "failed" + color.Reset + "\n"},
		{name: "message and error", args: []interface{}{"failed", errors.New("boom")}, expected: "❌ " + color.Red + "failed" + color.Reset + ":boom\n"},
		{name: "message and nil error", args: []interface{}{"failed", nil}, expected: "❌ " + color.Red + "failed" + color.Reset + "\n"},
		{name: "message and typed nil error", args: []interface{}{"failed", error(nil)}, expected: "❌ " + color.Red + "failed" + color.Reset + "\n"},
		{name: "error as message", args: []interface{}{errors.New("boom")}, expected: "❌ " + color.Red + "boom" + color.Reset + "\n"},
		{name: "non-string extra arg", args: []interface{}{"failed", 42}, expected: "❌ " + color.Red + "failed" + color.Reset + ":42\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := captureStdout(t, func() { Error(tt.args...) })
			if out != tt.expected {
				t.Errorf("Error printed %q, expected %q", out, tt.expected)
			}
		})
	}
}

func TestInfo(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		args     []interface{}
		expected string
	}{
		{name: "message only", message: "hello", expected: "ℹ️  " + color.Blue + "hello" + color.Reset + "\n"},
		{name: "with args", message: "hello", args: []interface{}{"a", 2}, expected: "ℹ️  " + color.Blue + "hello" + color.Reset + color.Gray + ": a, 2\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := captureStdout(t, func() { Info(tt.message, tt.args...) })
			if out != tt.expected {
				t.Errorf("Info printed %q, expected %q", out, tt.expected)
			}
		})
	}
}

func TestText(t *testing.T) {
	if out := captureStdout(t, func() { Text("a", "b", "c") }); out != "abc\n" {
		t.Errorf("Text printed %q, expected %q", out, "abc\n")
	}
}

func TestTextWithoutLn(t *testing.T) {
	if out := captureStdout(t, func() { TextWithoutLn("a", "b", "c") }); out != "abc" {
		t.Errorf("TextWithoutLn printed %q, expected %q", out, "abc")
	}
}

func TestTitle(t *testing.T) {
	out := captureStdout(t, func() { Title("My ", "title") })
	if !strings.HasPrefix(out, "\n"+color.Yellow+"My title"+color.Reset+"\n") {
		t.Errorf("Title printed %q, expected it to start with the colored title", out)
	}
	if !strings.Contains(out, "──────") {
		t.Errorf("Title printed %q, expected an underline", out)
	}
}

func TestFixedTextLength(t *testing.T) {
	out := captureStdout(t, func() { FixedTextLength(6, '.', "ab", "c") })
	if out != "abc..." {
		t.Errorf("FixedTextLength printed %q, expected %q", out, "abc...")
	}
}

func TestAdjustTextLength(t *testing.T) {
	tests := []struct {
		text     string
		length   int
		char     rune
		expected string
	}{
		{text: "abc", length: 6, char: '.', expected: "abc..."},
		{text: "abcdef", length: 3, char: '.', expected: "abc"},
		{text: "abc", length: 3, char: '.', expected: "abc"},
		{text: "", length: 2, char: '-', expected: "--"},
		{text: "abc", length: 0, char: '.', expected: ""},
		{text: "abc", length: -1, char: '.', expected: ""},
		{text: "héllo", length: 2, char: '.', expected: "hé"},
		{text: "héllo", length: 7, char: '.', expected: "héllo.."},
		{text: "✅ ok", length: 3, char: ' ', expected: "✅ o"},
		{text: "ab", length: 4, char: '─', expected: "ab──"},
	}

	for _, tt := range tests {
		if result := adjustTextLength(tt.text, tt.length, tt.char); result != tt.expected {
			t.Errorf("adjustTextLength(%q, %d, %q) = %q, expected %q", tt.text, tt.length, tt.char, result, tt.expected)
		}
	}
}

func TestAsk(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "Hello World\n", expected: "hello world"},
		{input: "  spaced  \n", expected: "spaced"},
		{input: "\n", expected: ""},
		{input: "", expected: ""}, // EOF
	}

	for _, tt := range tests {
		var result string
		out := captureStdout(t, func() {
			withStdin(t, tt.input, func() { result = Ask("Question? ") })
		})

		if result != tt.expected {
			t.Errorf("Ask with input %q = %q, expected %q", tt.input, result, tt.expected)
		}
		if !strings.Contains(out, "Question? ") {
			t.Errorf("Ask printed %q, expected the question", out)
		}
	}
}

func TestAskRaw(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "Hello World\n", expected: "Hello World"},
		{input: "  Keep Case  \n", expected: "Keep Case"},
		{input: "", expected: ""}, // EOF
	}

	for _, tt := range tests {
		var result string
		captureStdout(t, func() {
			withStdin(t, tt.input, func() { result = AskRaw("Question? ") })
		})

		if result != tt.expected {
			t.Errorf("AskRaw with input %q = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestAskBool(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{input: "y\n", expected: true},
		{input: "Y\n", expected: true},
		{input: "yes\n", expected: true},
		{input: " YES \n", expected: true},
		{input: "n\n", expected: false},
		{input: "no\n", expected: false},
		{input: "maybe\n", expected: false},
		{input: "\n", expected: false},
		{input: "", expected: false}, // EOF
	}

	for _, tt := range tests {
		var result bool
		out := captureStdout(t, func() {
			withStdin(t, tt.input, func() { result = AskBool("Continue?") })
		})

		if result != tt.expected {
			t.Errorf("AskBool with input %q = %v, expected %v", tt.input, result, tt.expected)
		}
		if !strings.Contains(out, "Continue? (y/n) ") {
			t.Errorf("AskBool printed %q, expected the question with (y/n)", out)
		}
	}
}

func TestAskInt(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{input: "42\n", expected: 42},
		{input: " -7 \n", expected: -7},
		{input: "\n", expected: 0},
		{input: "abc\n", expected: 0},
		{input: "3.5\n", expected: 0},
		{input: "", expected: 0}, // EOF
	}

	for _, tt := range tests {
		var result int
		captureStdout(t, func() {
			withStdin(t, tt.input, func() { result = AskInt("Number? ") })
		})

		if result != tt.expected {
			t.Errorf("AskInt with input %q = %d, expected %d", tt.input, result, tt.expected)
		}
	}
}

func TestAskFloat(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{input: "3.5\n", expected: 3.5},
		{input: "42\n", expected: 42},
		{input: " -0.25 \n", expected: -0.25},
		{input: "\n", expected: 0},
		{input: "abc\n", expected: 0},
		{input: "", expected: 0}, // EOF
	}

	for _, tt := range tests {
		var result float64
		captureStdout(t, func() {
			withStdin(t, tt.input, func() { result = AskFloat("Number? ") })
		})

		if result != tt.expected {
			t.Errorf("AskFloat with input %q = %v, expected %v", tt.input, result, tt.expected)
		}
	}
}
