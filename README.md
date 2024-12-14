# Starter pack for using GO in CLI

Starter package providing useful tools to setup a CLI application

## Help

You can easily setup a help system.

1. Create your help lines:

```go
help := starter.help
help.Content{
    Title: "My Application name",
    Lines: [][]string{
        {"title", "part 1 Title"},
        {"cmd", "command1", "option1", "option2", "option3"},
        {"cmd", "command2", "option1", "option2", "option3"},
        {"option", "-o/--option", "Explanation"},
    },
}
```

Lines should start with the type: {title, cmd}, next entries are the content.

## Config

You can load a config file from a YAML file:

1. Create a struct that reflects the data of your file

```go
type ConfigEntity struct {
	Test string `yaml:"test"`
}
```

2. Give the relative filepath (from home directory) and initiate the struct:

```go
configFilePath := ".oceanspy/starter/test.yaml"
configEntity := ConfigEntity{}

c, err := config.Construct(configFilePath, configEntity)
if err != nil {

}
```

3. Access the data:

You can load the whole configuration:

```go
all, err := c.Get()
if err != nil {
    //...
}
```

Or access a specific parameter:

```go
// Get with error
test, err := c.GetValue("test")
if err != nil {
    //...
}

// Get without error
test := c.GetValueOrEmpty("test")
```

## Message/Print

To make the printing more standard, you can use the message package.

```go
// Print a line break
Ln()

// Print Info message
Info(message string, args ...interface{})

// Print a success message
Success(messages ...string)

// Print a warning
Warning(messages ...string)

// Print an error message
Error(args ...interface{})

// Print a text message
Text(messages ...string)
TextWithoutLn(messages ...string)

// Print a title
Title(messages ...string)

// Print a fixed length text, and fill the text with the rune
TextFixedLength(length int, char rune, messages ...string)

// Ask for a string input to the user
Ask(message string) string

// Ask for a bool input to the user
AskBool(message string) bool
```

## Command

You can also deal with command.
This will allow you to deal with parameters and have options with or without values.

1. Setup:

```go
// Constructing the command & allowed parameters options
var optionsNames = map[string]string{
    "help": "h",
    "limit": "l",
}
var optionsWithValue = []string{ "limit" }
c, _ := command.Construct(os.Args, optionsNames, optionsWithValue)
```

2. Use the command parameters:

```go
// Transform short option name to long:
c.ShortToLongOptionName(arg string) (string, error)

// Get the number of commands:
c.Count() int

// Get command at position:
c.Get(at int) string

// Does the command has the option:
c.HasOption(name string) bool

// Get the option (string) value
c.GetOption(name string) (string, error)
```

## Autocomplete command

A simple autocomplete is available.

### Setup:

Create a variable that will hold the yaml data:

```go
	cmdAutocompleteYaml := `
command1:
command2:
command3:
  a:
    wuj:
    xaaua:
  b:
  c:
  d:
    y:
    z:
`
```

Each key is a command that will get autocomplete.

#### Having file completion in your program:

If you want to get file completion, just return `_file` and a normal autocomplete will be done.

#### How to

To have the autocompletion, just setup a first argument like "commands", to get the autocompletion:

```go
// Setup the autocomplete:
// You can use the command tools to get the parameters without the options:
cmdAutocomplete := commandautocomplete.Init(cmdAutocompleteYaml, cmd.Get())

switch cmd.GetAt(1) {
case "commands":
    autocomplete := cmdAutocomplete.Get()
    if len(autocomplete) > 0 {
        message.Text(autocomplete)
    }
// ...
default:
    //...
}
```

You can also get the last argument inputted:

```go
// Get last argument inputted:
message.Info("Last arg: ", cmdAutocomplete.GetLastArg())
```

### Oh-my-ZSH

You can use this file to put into `~/.oh-my-zsh/completions/`

In order to make it work, change `{{EXEC}}` to your executable.

```zsh
#compdef {{EXEC}}

# Define the completion function
_{{EXEC}}_completion() {
    # Get the current word being completed
    local current_word
    current_word="${words[CURRENT]}"

    # Check if the last word on the command line is quoted
    if [[ $current_word == *"'"* || $current_word == *'"'* ]]; then
        return
    fi

    # If this is the first argument, provide completion suggestions based on '{{EXEC}} commands'
    if [[ CURRENT -eq 2 ]]; then
        local completion_output
        completion_output=($({{EXEC}} commands | sed 's/\\ /--TEMP--/g' | tr ' ' '\n'))  # Replace \ with a temporary string, then split

        # Restore the temporary marker back to space before adding completions
        for word in "${completion_output[@]}"; do
            word="${word//--TEMP--/ }"  # Replace the temporary string back with space
            compadd -Q -- "$word"  # Add without quotes
        done
        return
    fi

    # Generate the appropriate completion command based on the number of arguments
    local completion_command="{{EXEC}} commands"
    for ((i=2; i<=CURRENT; i++)); do
        completion_command+=" ${words[i]}"
    done

    # Execute the completion command and capture the output
    local completion_output
    completion_output=($(eval $completion_command | sed 's/\\ /--TEMP--/g' | tr ' ' '\n'))  # Replace \ with a temporary string, then split

    # If '_file' is returned, use file completion
    if [[ " ${completion_output[*]} " == *" _file "* ]]; then
        _files
        return
    fi

    # Add the processed output as completions without quotes
    for word in "${completion_output[@]}"; do
        word="${word//--TEMP--/ }"  # Replace the temporary string back with space
        compadd -Q -- "$word"  # Add without quotes
    done
}

# Call the completion function
compdef _{{EXEC}}_completion {{EXEC}}
```
