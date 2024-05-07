package main

import (
	"example/network"
	"example/styles"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

var sub chan struct{} = make(chan struct{})

type model struct {
	sub 	chan struct{}
	width    int
	height   int
	viewport viewport.Model
	done	bool
	messages []network.Message
	input	textinput.Model
	header	string
	footer	string
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
			m.viewport.Style = styles.ChatStyle
			m.input.Width = msg.Width - 5

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
		m.messages = network.Messages
		return m, waitForActivity(m.sub)
	}

	m.input, cmd = m.input.Update(msg)

	return m, cmd
}

func (m model) View() string {
	all := ""
	for _, msg := range m.messages {
		all += styles.MessageFormat(msg)
	}
	m.viewport.SetContent(all)
	m.viewport.GotoBottom()

	return fmt.Sprintf("%s\n%s\n%s\n%s", styles.Header, m.viewport.View(), styles.ChatStyle.Render(m.input.View()), styles.Footer)
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