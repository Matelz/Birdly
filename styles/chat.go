package styles

import (
	"birdly/types"
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
	infoStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	timeStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	usernameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Bold(true)
)

func MessageFormatter(m types.MessageState) string {

	switch m.Type {
	// Normal message
	case types.NormalMessage:
		normalFormat := fmt.Sprintf("%s %s %s %s\n", timeStyle.Render(m.Timestamp), timeStyle.Render("-"), usernameStyle.Render("<<"+m.Username+">>"), m.Message)
		return normalFormat
		// Info message
	case types.InfoMessage:
		infoFormat := fmt.Sprintf("%s - %s %s\n", m.Timestamp, m.Username, m.Message)
		return infoStyle.Render(infoFormat)
		// Command message
	case types.CommandMessage:
		return m.Message
	default:
		return "error"

	}
}
