package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"os/user"
	"reflect"
)

type Config struct {
	Filepath   string
	YamlFile   []byte
	Data       interface{}
	DataLoaded bool
}

func Construct(filepath string, entity interface{}) (Config, error) {
	// Get current user
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Error fetching user:", err)
		os.Exit(1)
	}

	// Get the home directory
	homeDir := usr.HomeDir
	configFilePath := homeDir + "/" + filepath

	yamlFile, err := os.ReadFile(configFilePath)
	if err != nil {
		return Config{}, errors.New("Error opening config file: " + configFilePath)
	}

	return Config{Filepath: configFilePath, YamlFile: yamlFile, DataLoaded: false}, nil
}

func (c *Config) Get() (interface{}, error) {
	if c.DataLoaded {
		return c.Data, nil
	}

	err := yaml.Unmarshal(c.YamlFile, &c.Data)
	if err != nil {
		return c.Data, errors.New("Error parsing config file. Is it a valid yaml file ?")
	}

	c.DataLoaded = true

	return c.Data, nil
}

func (c *Config) GetValue(name string) (string, error) {
	config, err := c.Get()
	if err != nil {
		return "", err
	}

	v := reflect.ValueOf(config)

	// Handle map type
	if v.Kind() == reflect.Map {
		key := reflect.ValueOf(name)
		value := v.MapIndex(key)
		if !value.IsValid() {
			return "", errors.New("Field was not found")
		}
		return fmt.Sprintf("%v", value.Interface()), nil
	}

	// Ensure that we have a struct and not a pointer to a struct
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Check if the field exists
	field := v.FieldByName(name)
	if !field.IsValid() {
		return "", errors.New("Field was not found")
	}

	return field.String(), nil
}

func (c *Config) GetValueOrEmpty(name string) string {
	v, _ := c.GetValue(name)
	return v
}

func (c *Config) GetValues(name string) ([]string, error) {
	config, err := c.Get()
	if err != nil {
		return nil, err
	}

	v := reflect.ValueOf(config)

	// Handle map type
	if v.Kind() == reflect.Map {
		key := reflect.ValueOf(name)
		value := v.MapIndex(key)
		if !value.IsValid() {
			return nil, errors.New("Field was not found")
		}
		return []string{fmt.Sprintf("%v", value.Interface())}, nil
	}

	// Ensure that we have a struct and not a pointer to a struct
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Check if the field exists
	field := v.FieldByName(name)
	if !field.IsValid() {
		return nil, errors.New("Field was not found")
	}

	// Handle slice type
	if field.Kind() == reflect.Slice {
		if field.Len() == 0 {
			return []string{}, nil
		}

		var result []string
		for i := 0; i < field.Len(); i++ {
			result = append(result, fmt.Sprintf("%v", field.Index(i).Interface()))
		}

		return result, nil
	}

	return []string{field.String()}, nil
}

func (c *Config) GetValuesOrEmpty(name string) []string {
	v, _ := c.GetValues(name)
	return v
}
