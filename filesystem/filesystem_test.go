package filesystem

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestExpandHome(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home) // used by os.UserHomeDir on Windows

	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{name: "tilde only", value: "~", expected: home},
		{name: "tilde slash", value: "~/", expected: home},
		{name: "file in home", value: "~/config.yaml", expected: filepath.Join(home, "config.yaml")},
		{name: "nested path in home", value: "~/.config/app/config.yaml", expected: filepath.Join(home, ".config", "app", "config.yaml")},
		{name: "path is cleaned", value: "~/a/../b//c", expected: filepath.Join(home, "b", "c")},
		{name: "absolute path unchanged", value: "/etc/app/config.yaml", expected: "/etc/app/config.yaml"},
		{name: "relative path unchanged", value: "config/app.yaml", expected: "config/app.yaml"},
		{name: "empty path unchanged", value: "", expected: ""},
		{name: "other user's home not expanded", value: "~bob/file", expected: "~bob/file"},
		{name: "tilde not at start not expanded", value: "dir/~/file", expected: "dir/~/file"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExpandHome(tt.value)
			if err != nil {
				t.Fatalf("ExpandHome(%q) returned unexpected error: %v", tt.value, err)
			}

			if result != tt.expected {
				t.Errorf("ExpandHome(%q) = %q, expected %q", tt.value, result, tt.expected)
			}
		})
	}
}

func TestExpandHomeWithoutHomeDir(t *testing.T) {
	if runtime.GOOS == "windows" || runtime.GOOS == "plan9" {
		t.Skip("os.UserHomeDir uses a different environment variable on this platform")
	}
	t.Setenv("HOME", "")

	if _, err := ExpandHome("~/config.yaml"); err == nil {
		t.Error("expected an error when the home directory cannot be determined")
	}

	// Paths that need no expansion must not depend on the home directory.
	result, err := ExpandHome("/etc/config.yaml")
	if err != nil {
		t.Fatalf("unexpected error for a path without tilde: %v", err)
	}
	if result != "/etc/config.yaml" {
		t.Errorf("ExpandHome(%q) = %q, expected it unchanged", "/etc/config.yaml", result)
	}
}
