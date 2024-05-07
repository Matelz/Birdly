package main

import (
	"example/network"
	"example/other"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var sub chan struct{} = make(chan struct{})

type model struct {
	sub 	chan struct{}
	width    int
	height   int
	viewport viewport.Model
	done	bool
	messages []string
	input	textinput.Model
}

type responseMsg struct{}

func waitForActivity(sub chan struct{}) tea.Cmd {
	return func() tea.Msg {
		return responseMsg(<-sub)
	}
}

func NewModel() model {
	go network.ConnectToServer(sub)

	ti := textinput.New()
	ti.Placeholder = "Type something..."
	ti.Focus()

	return model{
		sub: sub,
		viewport: viewport.Model{},
		input: ti,
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 18

		if !m.done {
			m.viewport.Init()
			m.viewport.Style.Border(lipgloss.BlockBorder())
			m.done = true
		}
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.input.Value() != "" {
				network.SendMessage(m.input.Value())
				m.input.Reset()
			}
		}
	case responseMsg:
		m.messages = other.Messages
		return m, waitForActivity(m.sub)
	}

	m.input, cmd = m.input.Update(msg)

	return m, cmd
}

func (m model) View() string {
	all := ""
	for _, msg := range m.messages {
		all += msg + "\n"
	}
	m.viewport.SetContent(all)
	m.viewport.GotoBottom()

	return fmt.Sprintf("%s\n%s",m.viewport.View(),m.input.View())
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		waitForActivity(m.sub),
	)
}

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "host" {
			go network.CreateServer()
		}
	}
	p := tea.NewProgram(NewModel(), tea.WithOutput(os.Stdout))
    if _, err := p.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}