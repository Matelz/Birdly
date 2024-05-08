package styles

import (
	"birdly/network"
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var Title = `
    ____  _          ____     
   / __ )(_)________/ / /_  __
  / __  / / ___/ __  / / / / /
 / /_/ / / /  / /_/ / / /_/ / 
/_____/_/_/   \__,_/_/\__, /  
                     /____/   
`

var (
	TitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Bold(true)

	ChatStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())

	HeaderStyle = lipgloss.NewStyle().Faint(true).Italic(true)

	UsernameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Bold(true)

	BroadcastStyle = lipgloss.NewStyle().Faint(true).Italic(true)
)

var Header = fmt.Sprintf("%s\n%s", TitleStyle.Render(Title), HeaderStyle.Render("v1.0.0\nA p2p tui chat client."))

var Footer = HeaderStyle.Render("https://github.com/Matelz/Birdly")

func MessageFormat(message network.Message) string {
	var res string
	var user = network.Users[message.UserID]

	switch message.MessageType {
	case 1:
		res = fmt.Sprintf("%s %s", UsernameStyle.Render("<<"+user.Name+">>"), message.Message)
	case 2, 3:
		res = BroadcastStyle.Render(fmt.Sprintf("%s %s", user.Name, message.Message))
	}

	res += "\n"

	return res
}