package help

import (
	"github.com/oceanspy/go-cli-starter/src/color"
	"github.com/oceanspy/go-cli-starter/src/message"
)

type Content struct {
	Title string
	Lines [][]string
}

func Print(help Content) {
	message.Info("Usage:", help.Title, "<command>")
	message.Text(color.Gray, "────────────────────────────────────", color.Reset)

	for _, line := range help.Lines {
		switch line[0] {
		case "title":
			message.TextWithoutLn("    ")
			title := ""
			for i, element := range line {
				if i == 0 {
					continue
				}
				title += element
			}

			message.Title("    ", title)
		case "cmd":
			message.TextWithoutLn("    ")
			for i, element := range line {
				if i == 0 {
					continue
				}
				switch i {
				case 1:
					message.TextWithoutLn(color.Red)
					message.TextWithoutLn(element)
					message.TextWithoutLn(color.Reset)
				case 2:
					message.TextWithoutLn(color.Green)
					message.TextWithoutLn(element)
					message.TextWithoutLn(color.Reset)
				case 3:
					message.TextWithoutLn(color.Magenta)
					message.TextWithoutLn(element)
					message.TextWithoutLn(color.Reset)
				case 4:
					message.TextWithoutLn(color.Yellow)
					message.TextWithoutLn(element)
					message.TextWithoutLn(color.Reset)
				default:
					message.TextWithoutLn(element)
				}
			}
		default:
		}
		message.Ln()
	}
}
