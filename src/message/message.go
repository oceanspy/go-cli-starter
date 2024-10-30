package message

import (
	"bufio"
	"fmt"
	"github.com/oceanspy/go-cli-starter/src/color"
	"os"
	"strings"
)

func Ln() {
	fmt.Println()
}

func Success(messages ...string) {
	message := ""
	for _, m := range messages {
		message += m
	}

	fmt.Print("✅ ", color.Get("green"), message, color.Reset)
	fmt.Println()
}

func Warning(messages ...string) {
	message := ""
	for _, m := range messages {
		message += m
	}

	fmt.Print("⚠️ ", color.Get("yellow"), message, color.Reset)
	fmt.Println()
}

func Error(args ...interface{}) {
	if len(args) == 0 {
		fmt.Println("❌")
		return
	}

	message := args[0].(string)

	if len(args) == 1 {
		fmt.Print("❌ ", color.Get("red"), message, color.Reset)
		fmt.Println()
		return
	}

	err := args[1].(error)
	fmt.Print("❌ ", color.Get("red"), message, color.Reset)
	if err != nil {
		fmt.Print(":", err)
		fmt.Println()
	}
}

func Info(message string, args ...interface{}) {
	fmt.Print("ℹ️  ", color.Get("blue"), message, color.Reset)
	if len(args) > 0 {
		fmt.Print(color.Get("gray"), ": ")
		for i, arg := range args {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Print(arg)
		}
	}
	fmt.Println()
}

func Text(messages ...string) {
	message := ""
	for _, m := range messages {
		message += m
	}

	fmt.Print(message)
	fmt.Println()
}

func FixedTextLength(length int, char rune, messages ...string) {
	message := ""
	for _, m := range messages {
		message += m
	}

	fmt.Print(adjustTextLength(message, length, char))
}

func adjustTextLength(text string, length int, char rune) string {
	textLen := len(text)

	if textLen > length {
		// If the text is too long, shorten it
		return text[:length]
	} else if textLen < length {
		// If the text is too short, pad it with the given char
		padding := strings.Repeat(string(char), length-textLen)
		return text + padding
	}
	// If the text is already the correct length, return it as is
	return text
}

func Title(messages ...string) {
	message := ""
	for _, m := range messages {
		message += m
	}

	fmt.Println()
	fmt.Print(color.Get("yellow"), message, color.Reset)
	fmt.Println()
	fmt.Println(color.Get("yellow"), "──────────────────────────────", color.Reset)
}

func Ask(message string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(color.Get("yellow"), message, color.Reset)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return ""
	}

	return strings.TrimSpace(strings.ToLower(response))
}

func AskBool(message string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(color.Get("yellow"), message, " (y/n) ", color.Reset)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return false
	}

	response = strings.TrimSpace(strings.ToLower(response))

	if response == "y" || response == "yes" {
		return true
	}

	return false
}
