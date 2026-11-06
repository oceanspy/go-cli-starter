package help

import (
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

func TestPrint(t *testing.T) {
	out := captureStdout(t, func() {
		Print(Content{
			Title: "My Application",
			Lines: [][]string{
				{"title", "Part ", "1"},
				{"cmd", "command1", "option1", "option2", "option3", "extra"},
				{"cmd", "command2", "", "option2"},
				{"option", "-o/--option", "Explanation"},
				{"unknown", "should not be printed"},
				{},
			},
		})
	})

	expected := []string{
		"Help Center",
		"My Application",
		"Part 1",
		"    " + color.Red + " command1" + color.Reset + color.Green + " option1" + color.Reset + color.Magenta + " option2" + color.Reset + color.Yellow + " option3" + color.Reset + " extra\n",
		"    " + color.Red + " command2" + color.Reset + color.Magenta + " option2" + color.Reset + "\n",
		"        " + color.Blue + " -o/--option" + color.Reset + " Explanation\n",
	}

	for _, e := range expected {
		if !strings.Contains(out, e) {
			t.Errorf("Print output doesn't contain %q\noutput: %q", e, out)
		}
	}

	if strings.Contains(out, "should not be printed") {
		t.Errorf("Print output contains a line with an unknown type\noutput: %q", out)
	}
}

func TestPrintNoLines(t *testing.T) {
	out := captureStdout(t, func() { Print(Content{Title: "Empty"}) })

	if !strings.Contains(out, "Empty") {
		t.Errorf("Print output doesn't contain the title\noutput: %q", out)
	}
}
