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
			message.Text("    ")
			title := ""
			for i, element := range line {
				if i == 0 {
					continue
				}
				title += element
			}

			message.Title("    ", title)
		case "cmd":
			message.Text("    ")
			for i, element := range line {
				if i == 0 {
					continue
				}
				switch i {
				case 1:
					message.Text(color.Red)
					message.Text(element)
					message.Text(color.Reset)
				case 2:
					message.Text(color.Green)
					message.Text(element)
					message.Text(color.Reset)
				case 3:
					message.Text(color.Magenta)
					message.Text(element)
					message.Text(color.Reset)
				case 4:
					message.Text(color.Yellow)
					message.Text(element)
					message.Text(color.Reset)
				default:
					message.Text(element)
				}
			}
		default:
		}
	}
}
