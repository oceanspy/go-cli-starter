package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

const testYaml = `
name: starter
float: 3.5
int: 42
notanumber: abc
empty:
list:
  - a
  - b
`

func newTestConfig(yamlStr string) *Config {
	return &Config{Filepath: "test.yaml", YamlFile: []byte(yamlStr)}
}

type testEntity struct {
	Name  string
	Tags  []string
	Empty []string
}

// newStructConfig returns a config whose data is already loaded as a struct.
func newStructConfig() *Config {
	return &Config{
		Data:       &testEntity{Name: "starter", Tags: []string{"a", "b"}},
		DataLoaded: true,
	}
}

type fileEntity struct {
	Name    string   `yaml:"name"`
	Limit   int      `yaml:"limit"`
	Ratio   float64  `yaml:"ratio"`
	Enabled bool     `yaml:"enabled"`
	Tags    []string `yaml:"tags"`
	Nothing *string  `yaml:"nothing"`
	NoTag   string
	private string
}

const fileYaml = `
name: starter
limit: 10
ratio: 0.5
enabled: true
tags:
  - a
  - b
notag: untagged
`

// writeConfigFile creates a config file in a temporary home directory and
// returns its path relative to that home.
func writeConfigFile(t *testing.T, content string) string {
	t.Helper()

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home) // used by os.UserHomeDir on Windows

	relPath := filepath.Join(".myapp", "config.yaml")
	if err := os.MkdirAll(filepath.Join(home, ".myapp"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(home, relPath), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	return relPath
}

func TestConstructMissingFile(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("USERPROFILE", os.Getenv("HOME"))

	_, err := Construct(".myapp/this-file-does-not-exist.yaml", nil)
	if err == nil {
		t.Error("expected an error for a missing config file")
	}
}

func TestConstructReadsFileFromHome(t *testing.T) {
	relPath := writeConfigFile(t, fileYaml)

	c, err := Construct(relPath, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Filepath != filepath.Join(os.Getenv("HOME"), relPath) {
		t.Errorf("Filepath = %q, expected it inside the home directory", c.Filepath)
	}
	if result := c.GetValueOrEmpty("name"); result != "starter" {
		t.Errorf("GetValueOrEmpty(name) = %q, expected starter", result)
	}
}

func TestConstructWithEntityPointer(t *testing.T) {
	relPath := writeConfigFile(t, fileYaml)

	entity := fileEntity{}
	c, err := Construct(relPath, &entity)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := c.Get()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// The struct given by the caller is filled.
	if entity.Name != "starter" || entity.Limit != 10 || !reflect.DeepEqual(entity.Tags, []string{"a", "b"}) {
		t.Errorf("entity was not filled: %#v", entity)
	}
	if data != &entity {
		t.Errorf("Get() should return the given entity pointer, got %#v", data)
	}
}

func TestConstructWithEntityValue(t *testing.T) {
	relPath := writeConfigFile(t, fileYaml)

	c, err := Construct(relPath, fileEntity{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := c.Get()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entity, ok := data.(*fileEntity)
	if !ok {
		t.Fatalf("Get() should return a *fileEntity, got %T", data)
	}
	if entity.Name != "starter" || !entity.Enabled {
		t.Errorf("entity was not filled: %#v", entity)
	}
}

func TestEntityValues(t *testing.T) {
	var nilEntity *fileEntity // a nil pointer is treated like a struct value
	for name, entity := range map[string]interface{}{"value": fileEntity{}, "pointer": &fileEntity{}, "nil pointer": nilEntity} {
		t.Run(name, func(t *testing.T) {
			c := &Config{YamlFile: []byte(fileYaml), Entity: entity}

			values := map[string]string{
				"name":    "starter",
				"limit":   "10",
				"ratio":   "0.5",
				"enabled": "true",
				"nothing": "",
				"notag":   "untagged",
				"NoTag":   "untagged",
				"Name":    "starter", // field name also works
			}
			for key, expected := range values {
				result, err := c.GetValue(key)
				if err != nil {
					t.Errorf("GetValue(%q) returned an error: %v", key, err)
				}
				if result != expected {
					t.Errorf("GetValue(%q) = %q, expected %q", key, result, expected)
				}
			}

			if _, err := c.GetValue("private"); err == nil {
				t.Error("unexported fields should not be found")
			}
			if c.GetInt("limit") != 10 {
				t.Errorf("GetInt(limit) = %d, expected 10", c.GetInt("limit"))
			}
			if c.GetFloat("ratio") != 0.5 {
				t.Errorf("GetFloat(ratio) = %v, expected 0.5", c.GetFloat("ratio"))
			}
			if result, _ := c.GetValues("tags"); !reflect.DeepEqual(result, []string{"a", "b"}) {
				t.Errorf("GetValues(tags) = %#v, expected [a b]", result)
			}
			if result, err := c.GetValues("nothing"); err != nil || !reflect.DeepEqual(result, []string{}) {
				t.Errorf("GetValues(nothing) = %#v, %v ; expected empty list", result, err)
			}
		})
	}
}

func TestEntityInvalidYaml(t *testing.T) {
	c := &Config{YamlFile: []byte("limit: not-a-number"), Entity: &fileEntity{}}

	if _, err := c.Get(); err == nil {
		t.Error("expected an error when the YAML doesn't match the entity")
	}
}

func TestGet(t *testing.T) {
	c := newTestConfig(testYaml)

	data, err := c.Get()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !c.DataLoaded {
		t.Error("DataLoaded should be true after Get")
	}

	m, ok := data.(map[string]interface{})
	if !ok {
		t.Fatalf("expected a map, got %T", data)
	}
	if m["name"] != "starter" {
		t.Errorf("name = %v, expected starter", m["name"])
	}

	// Once loaded, the YAML must not be parsed again.
	c.YamlFile = []byte("name: changed")
	data, _ = c.Get()
	if data.(map[string]interface{})["name"] != "starter" {
		t.Error("Get should return cached data once loaded")
	}
}

func TestGetInvalidYaml(t *testing.T) {
	c := newTestConfig("name: [unclosed")

	if _, err := c.Get(); err == nil {
		t.Error("expected an error for invalid YAML")
	}
	if c.DataLoaded {
		t.Error("DataLoaded should stay false when parsing fails")
	}
}

func TestGetValue(t *testing.T) {
	c := newTestConfig(testYaml)

	tests := []struct {
		name     string
		expected string
		wantErr  bool
	}{
		{name: "name", expected: "starter"},
		{name: "float", expected: "3.5"},
		{name: "int", expected: "42"},
		{name: "empty", expected: ""},
		{name: "missing", expected: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := c.GetValue(tt.name)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetValue(%q) error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if result != tt.expected {
				t.Errorf("GetValue(%q) = %q, expected %q", tt.name, result, tt.expected)
			}
		})
	}
}

func TestGetValueStruct(t *testing.T) {
	c := newStructConfig()

	result, err := c.GetValue("Name")
	if err != nil || result != "starter" {
		t.Errorf("GetValue(Name) = %q, %v ; expected starter, nil", result, err)
	}

	if _, err := c.GetValue("Missing"); err == nil {
		t.Error("expected an error for a missing struct field")
	}
}

func TestGetValueInvalidYaml(t *testing.T) {
	c := newTestConfig("name: [unclosed")

	if _, err := c.GetValue("name"); err == nil {
		t.Error("expected an error for invalid YAML")
	}
	if result := c.GetValueOrEmpty("name"); result != "" {
		t.Errorf("GetValueOrEmpty = %q, expected empty", result)
	}
}

func TestGetValueOrEmpty(t *testing.T) {
	c := newTestConfig(testYaml)

	if result := c.GetValueOrEmpty("name"); result != "starter" {
		t.Errorf("GetValueOrEmpty(name) = %q, expected starter", result)
	}
	if result := c.GetValueOrEmpty("missing"); result != "" {
		t.Errorf("GetValueOrEmpty(missing) = %q, expected empty", result)
	}
}

func TestGetFloat(t *testing.T) {
	c := newTestConfig(testYaml)

	tests := map[string]float64{
		"float":      3.5,
		"int":        42,
		"notanumber": 0,
		"missing":    0,
	}

	for name, expected := range tests {
		if result := c.GetFloat(name); result != expected {
			t.Errorf("GetFloat(%q) = %v, expected %v", name, result, expected)
		}
	}
}

func TestGetInt(t *testing.T) {
	c := newTestConfig(testYaml)

	tests := map[string]int{
		"int":        42,
		"float":      0, // "3.5" is not an int
		"notanumber": 0,
		"missing":    0,
	}

	for name, expected := range tests {
		if result := c.GetInt(name); result != expected {
			t.Errorf("GetInt(%q) = %v, expected %v", name, result, expected)
		}
	}
}

func TestGetValues(t *testing.T) {
	c := newTestConfig(testYaml)

	tests := []struct {
		name     string
		expected []string
	}{
		{name: "name", expected: []string{"starter"}},
		{name: "empty", expected: []string{}},
		{name: "missing", expected: []string{}},
		{name: "list", expected: []string{"a", "b"}},
		{name: "int", expected: []string{"42"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := c.GetValues(tt.name)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GetValues(%q) = %#v, expected %#v", tt.name, result, tt.expected)
			}
		})
	}
}

func TestGetValuesStruct(t *testing.T) {
	c := newStructConfig()

	tests := []struct {
		name     string
		expected []string
		wantErr  bool
	}{
		{name: "Name", expected: []string{"starter"}},
		{name: "Tags", expected: []string{"a", "b"}},
		{name: "Empty", expected: []string{}},
		{name: "Missing", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := c.GetValues(tt.name)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetValues(%q) error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GetValues(%q) = %#v, expected %#v", tt.name, result, tt.expected)
			}
		})
	}
}

func TestGetValuesOrEmpty(t *testing.T) {
	c := newTestConfig("name: [unclosed")

	if result := c.GetValuesOrEmpty("name"); result != nil {
		t.Errorf("GetValuesOrEmpty = %#v, expected nil on error", result)
	}

	c = newTestConfig(testYaml)
	if result := c.GetValuesOrEmpty("name"); !reflect.DeepEqual(result, []string{"starter"}) {
		t.Errorf("GetValuesOrEmpty(name) = %#v, expected [starter]", result)
	}
}
