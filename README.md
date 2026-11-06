# Starter pack for using GO in CLI

Starter package providing useful tools to setup a CLI application

## Help

You can easily setup a help system.

1. Create your help lines and print them:

```go
help.Print(help.Content{
    Title: "My Application name",
    Lines: [][]string{
        {"title", "part 1 Title"},
        {"cmd", "command1", "option1", "option2", "option3"},
        {"cmd", "command2", "option1", "option2", "option3"},
        {"option", "-o/--option", "Explanation"},
    },
})
```

Lines should start with the type: {title, cmd, option}, next entries are the content.
Lines with another type are ignored, and empty entries are skipped.

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

c, err := config.Construct(configFilePath, &configEntity)
if err != nil {

}
```

With a pointer, `configEntity` is filled when the data is loaded.
You can also pass `ConfigEntity{}` (a new struct is filled and returned by `c.Get()`),
or `nil` to load the data into a `map[string]interface{}`.

3. Access the data:

You can load the whole configuration:

```go
all, err := c.Get()
if err != nil {
    //...
}
```

Or access a specific parameter, by its YAML key (or its struct field name):

```go
// Get with error
test, err := c.GetValue("test")
if err != nil {
    //...
}

// Get without error
test := c.GetValueOrEmpty("test")

// Get a number (0 if missing or not a number)
limit := c.GetInt("limit")
ratio := c.GetFloat("ratio")

// Get the items of a list (or a single value as a one-item list)
values, err := c.GetValues("test")
values := c.GetValuesOrEmpty("test")
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
FixedTextLength(length int, char rune, messages ...string)

// Ask for a string input to the user (trimmed and lowercased)
Ask(message string) string

// Ask for a string input to the user (trimmed, case kept)
AskRaw(message string) string

// Ask for a bool input to the user ("y" or "yes" is true)
AskBool(message string) bool

// Ask for a number input to the user (0 if empty or invalid)
AskInt(message string) int
AskFloat(message string) float64
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
c, p := command.Construct(os.Args, optionsNames, optionsWithValue)
```

An unknown option prints an error and exits the program with code 1.

2. Use the command parameters:

```go
// Transform short option name to long:
p.ShortToLongOptionName(arg string) (string, error)

// Get the number of commands:
c.Count() int

// Get command at position:
c.GetAt(at int) (string, error)

// Get commands from position to the end:
c.GetAfter(at int) ([]string, error)

// Get all commands (without the options) or the raw arguments:
c.Get() []string
c.GetRaw() []string

// Does the command has the option (long name only):
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
Suggestions are returned in the same order as in the YAML.

#### Having file completion in your program:

If you want to get file completion, add a `_file` key under the command: `Complete()` returns `_file`
and your shell completion script can switch to normal file completion (e.g. `_files` in zsh).

`_file` can be written as a key (with `:`) or as a value (without `:`), these are all the same:

```yaml
send:
  _file:
send:
  _file
send: _file
```

#### How to

Set up a first argument like "commands" that prints the suggestions, one per line.
Give the options of your program with `WithOptions` so they are suggested too:

```go
var optionsNames = map[string]string{
    "help":  "h",
    "limit": "l",
}
var optionsWithValue = []string{"limit"}
cmd, _ := command.Construct(os.Args, optionsNames, optionsWithValue)

arg, _ := cmd.GetAt(1)
switch arg {
case "commands":
    suggestions := commandautocomplete.Init(cmdAutocompleteYaml, cmd.Get()).
        WithOptions(optionsNames, optionsWithValue).
        Complete()
    if len(suggestions) > 0 {
        message.Text(strings.Join(suggestions, "\n"))
    }
// ...
default:
    //...
}
```

`Complete()` reads the words after `commands --`. The last word is always the word being typed,
even when it is empty:

| Shell calls                              | Suggestions                              |
|------------------------------------------|------------------------------------------|
| `myprog commands -- com`                 | `command1 command2 command3`              |
| `myprog commands -- command3`            | `command3` (still typing it)             |
| `myprog commands -- command3 ""`         | `a b c d` (what comes after `command3`)  |
| `myprog commands -- command3 --li`       | `--limit`                                |
| `myprog commands -- --limit 5 command3 ""` | `a b c d` (options and values are skipped) |

`--` ends the options in `command.Construct`, so the words being completed are never read as
options of your program.

### ZSH (and Oh-my-ZSH)

Put this file in a directory of your `$fpath` (for Oh-my-ZSH: `~/.oh-my-zsh/completions/`),
named `_{{EXEC}}`, and change `{{EXEC}}` to your executable:

```zsh
#compdef {{EXEC}}

_{{EXEC}}() {
    local -a suggestions

    # Send the words before the cursor (without their quotes) and the word being
    # typed, even when it is empty. "--" stops them being read as program options.
    suggestions=(${(f)"$({{EXEC}} commands -- "${(@Q)words[2,CURRENT]}" 2>/dev/null)"})

    # If '_file' is returned, use file completion
    if [[ ${suggestions[1]} == _file ]]; then
        _files
        return
    fi

    (( ${#suggestions} )) && compadd -- "${suggestions[@]}"
}

_{{EXEC}} "$@"
```

Words with spaces, quotes, `~` or `*` are handled by zsh: suggestions are inserted with the
right quoting, and nothing typed on the command line is run during completion.

### Legacy: `Get()`

`Get()` is the previous autocomplete. It still works, but can't tell "still typing a word" from
"word done", so `myprog command3<TAB>` doesn't complete. It returns the suggestions separated by spaces:

```go
cmdAutocomplete := commandautocomplete.Init(cmdAutocompleteYaml, cmd.Get())
autocomplete := cmdAutocomplete.Get()
if len(autocomplete) > 0 {
    message.Text(autocomplete)
}

// Get last argument inputted:
message.Info("Last arg: ", cmdAutocomplete.GetLastArg())
```

The matching zsh script:

```zsh
#compdef {{EXEC}}

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

_describe {{EXEC}}
```

## Filesystem

`ExpandHome` replaces a leading `~` with the home directory of the current user:

```go
path, err := filesystem.ExpandHome("~/.config/myapp/config.yaml")
if err != nil {
    //...
}
```

Paths that don't start with `~` or `~/` are returned unchanged.

## Testing

```sh
make test          # run all tests
make test-verbose  # run all tests with verbose output
make cover         # run all tests with a coverage summary
make vet           # run go vet
```
