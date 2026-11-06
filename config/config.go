package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Filepath   string
	YamlFile   []byte
	Data       interface{}
	DataLoaded bool

	// Entity is the struct (or pointer to struct) the YAML is loaded into.
	// When nil, the YAML is loaded into a map.
	Entity interface{}
}

// Construct reads the YAML file at filepath, relative to the home directory.
//
// entity describes the data of the file. With a pointer to a struct, that
// struct is filled when the data is loaded. With a struct value, a new struct
// of the same type is filled. With nil, the data is loaded into a map.
func Construct(path string, entity interface{}) (Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, fmt.Errorf("Error fetching home directory: %w", err)
	}

	configFilePath := filepath.Join(homeDir, path)

	yamlFile, err := os.ReadFile(configFilePath)
	if err != nil {
		return Config{}, errors.New("Error opening config file: " + configFilePath)
	}

	return Config{Filepath: configFilePath, YamlFile: yamlFile, DataLoaded: false, Entity: entity}, nil
}

func (c *Config) Get() (interface{}, error) {
	if c.DataLoaded {
		return c.Data, nil
	}

	var target interface{}
	switch {
	case c.Entity == nil:
		target = &c.Data
	case reflect.TypeOf(c.Entity).Kind() == reflect.Ptr && !reflect.ValueOf(c.Entity).IsNil():
		target = c.Entity
	default:
		t := reflect.TypeOf(c.Entity)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		target = reflect.New(t).Interface()
	}

	err := yaml.Unmarshal(c.YamlFile, target)
	if err != nil {
		return c.Data, errors.New("Error parsing config file. Is it a valid yaml file ?")
	}

	if c.Entity != nil {
		c.Data = target
	}
	c.DataLoaded = true

	return c.Data, nil
}

// lookup finds the value of name in the loaded data.
// For a map, name is the key. For a struct, name is the YAML key or the field name.
func (c *Config) lookup(name string) (reflect.Value, bool, error) {
	data, err := c.Get()
	if err != nil {
		return reflect.Value{}, false, err
	}

	v := reflect.ValueOf(data)
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return reflect.Value{}, false, nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Map:
		if v.Type().Key().Kind() != reflect.String {
			return reflect.Value{}, false, nil
		}
		value := v.MapIndex(reflect.ValueOf(name).Convert(v.Type().Key()))
		return value, value.IsValid(), nil
	case reflect.Struct:
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if !f.IsExported() {
				continue
			}
			// Without a yaml tag, the YAML key is the lowercased field name
			key := strings.Split(f.Tag.Get("yaml"), ",")[0]
			if key == "" {
				key = strings.ToLower(f.Name)
			}
			if key == name || f.Name == name {
				return v.Field(i), true, nil
			}
		}
	}

	return reflect.Value{}, false, nil
}

// unwrap follows interfaces and pointers; it returns an invalid value for nil.
func unwrap(v reflect.Value) reflect.Value {
	for v.IsValid() && (v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr) {
		if v.IsNil() {
			return reflect.Value{}
		}
		v = v.Elem()
	}

	return v
}

func format(v reflect.Value) string {
	v = unwrap(v)
	if !v.IsValid() {
		return ""
	}

	return fmt.Sprintf("%v", v.Interface())
}

func (c *Config) GetValue(name string) (string, error) {
	value, found, err := c.lookup(name)
	if err != nil {
		return "", err
	}
	if !found {
		return "", errors.New("Field was not found")
	}

	return format(value), nil
}

func (c *Config) GetValueOrEmpty(name string) string {
	v, _ := c.GetValue(name)
	return v
}

func (c *Config) GetFloat(name string) float64 {
	strValue := c.GetValueOrEmpty(name)
	if strValue == "" {
		return 0
	}

	value, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		return 0
	}

	return value
}

func (c *Config) GetInt(name string) int {
	strValue := c.GetValueOrEmpty(name)
	if strValue == "" {
		return 0
	}

	value, err := strconv.Atoi(strValue)
	if err != nil {
		return 0
	}

	return value
}

// GetValues returns the values of a list, or a single value as a one-element list.
// A missing key in a map, or an empty value, returns an empty list.
// A missing field in a struct returns an error.
func (c *Config) GetValues(name string) ([]string, error) {
	value, found, err := c.lookup(name)
	if err != nil {
		return nil, err
	}

	if !found {
		data := unwrap(reflect.ValueOf(c.Data))
		if data.IsValid() && data.Kind() == reflect.Struct {
			return nil, errors.New("Field was not found")
		}
		return []string{}, nil
	}

	value = unwrap(value)
	if !value.IsValid() {
		return []string{}, nil
	}

	if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
		result := make([]string, 0, value.Len())
		for i := 0; i < value.Len(); i++ {
			result = append(result, format(value.Index(i)))
		}
		return result, nil
	}

	return []string{format(value)}, nil
}

func (c *Config) GetValuesOrEmpty(name string) []string {
	v, _ := c.GetValues(name)
	return v
}
