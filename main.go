package main

import (
	"birdly/network"
	"birdly/styles"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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
		m.input.Width = msg.Width - 5

		if !m.done {
			m.viewport.Init()
			m.viewport.Style = styles.ChatStyle
			m.viewport.MouseWheelEnabled = true

			m.done = true
		}
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			trimmed := strings.TrimSpace(m.input.Value())
			if trimmed != "" {
				network.SendMessage(trimmed)
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
	m.viewport.LineDown(1)

	return fmt.Sprintf("%s\n%s\n%s\n%s", styles.Header, m.viewport.View(), styles.ChatStyle.Render(m.input.View()), styles.Footer)
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		waitForActivity(m.sub),
	)
}

func main() {
	// Parse flags
	flg := flag.NewFlagSet("Birdly", flag.ExitOnError)
	ip := flg.String("host", "127.0.0.1", "ip used to connect to the server")
	port := flg.String("port", "8080", "port used to connect/host the server")
	name := flg.String("name", "anon", "username used in chat")

	args := os.Args

	if len(args) < 2{
		log.Fatal("Please include a command, eg:. 'connect' | 'host'")
	}

	switch args[1]{
		case "host":
			flg.Parse(args[2:])
			go network.CreateServer(*port)
			time.Sleep(5 * time.Second)
			go network.ConnectToServer("localhost", *port, *name, sub)
		case "connect":
			flg.Parse(args[2:])
			go network.ConnectToServer(*ip, *port, *name, sub)
	}

	p := tea.NewProgram(NewModel(), tea.WithAltScreen())
    if _, err := p.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}