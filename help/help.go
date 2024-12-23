package help

import (
	"github.com/oceanspy/go-cli-starter/color"
	"github.com/oceanspy/go-cli-starter/message"
)

type Content struct {
	Title string
	Lines [][]string
}

func Print(help Content) {
	message.Info("Help Center", help.Title)

	for _, line := range help.Lines {
		switch line[0] {
		case "title":
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

				if len(element) == 0 {
					continue
				}

				switch i {
				case 1:
					message.TextWithoutLn(color.Red)
					message.TextWithoutLn(" ")
					message.TextWithoutLn(element)
					message.TextWithoutLn(color.Reset)
				case 2:
					message.TextWithoutLn(color.Green)
					message.TextWithoutLn(" ")
					message.TextWithoutLn(element)
					message.TextWithoutLn(color.Reset)
				case 3:
					message.TextWithoutLn(color.Magenta)
					message.TextWithoutLn(" ")
					message.TextWithoutLn(element)
					message.TextWithoutLn(color.Reset)
				case 4:
					message.TextWithoutLn(color.Yellow)
					message.TextWithoutLn(" ")
					message.TextWithoutLn(element)
					message.TextWithoutLn(color.Reset)
				default:
					message.TextWithoutLn(" ")
					message.TextWithoutLn(element)
				}
			}
			message.Ln()
		case "option":
			message.TextWithoutLn("        ")
			for i, element := range line {
				if i == 0 {
					continue
				}

				if len(element) == 0 {
					continue
				}

				switch i {
				case 1:
					message.TextWithoutLn(color.Blue)
					message.TextWithoutLn(" ")
					message.TextWithoutLn(element)
					message.TextWithoutLn(color.Reset)
				default:
					message.TextWithoutLn(" ")
					message.TextWithoutLn(element)
				}
			}
			message.Ln()
		default:
		}
	}
}
