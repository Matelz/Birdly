package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// User can pass variables to the exe eg:. (birdly.exe <host|serve> <optional: -p port -h host -u username -r room_id>)

// If none are given, start the exe

// On the first screen the user will be prompted to insert a username (optional)

// Then the user will be prompted with a list of current rooms in the local network

// If the user wants, it can insert a room id to connect to

// If the room has a password, it will be prompted now

// Then the user will connect to the chat room

// The user can send messages in the chat, that will be broadcasted to all connected users

// It can also send a private message to someone by using the command (:pm <to> <message>)

type command struct {
	command string
	description string
	function func()
}

type messageState struct {
	timestamp string
	username string
	message string
}

type roomState struct {
	messages []messageState
	clients  []string
}

type model struct {
	// Set a bool so we know when the window was initialized 
	ready bool
	state roomState
	chatbox viewport.Model
	input textinput.Model
}

var commands = []command {
	{
		"pm",
		"Send private message to user (:pm <to> <message>)",
		func(){},
	},
}

func initialModel() model {
	chatbox := viewport.New(10, 10) 
	chatbox.Style = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())

	input := textinput.New()
	input.Placeholder = "Type message"
	input.Focus()

	return model{
		false,
		// Initialize room state with no messages and no clients
		roomState{
			[]messageState{},
			[]string{},
		},
		chatbox,
		input,
	}
}

func handleCommand(message string) (command string, found bool) {
	// Command format (:<command> <arg1> <arg...>)
	r := regexp.MustCompile(`:[a-zA-Z]{1,}`)
	if cmd := r.FindString(message); cmd != "" {
		for _, c := range commands {
			cmd = strings.ReplaceAll(cmd, ":", "")
			if c.command == cmd {
				return c.description, true
			}
		}
	}

	return "Command not found", false
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			if !m.ready {
				m.chatbox.Width = msg.Width
				m.chatbox.Height = msg.Height - 3

				m.input.Width = msg.Width - 5
				m.ready = true
			}
		case tea.KeyMsg:
			switch msg.String() {
				case "ctrl+c":
					return m, tea.Quit
				case "enter":
					if strings.HasPrefix(m.input.Value(), ":") {
						cmd, fnd := handleCommand(m.input.Value())
						if fnd {
							m.state.messages = append(m.state.messages, messageState{time.Now().Format(time.Kitchen), "anon", fmt.Sprintf("%s found", cmd)})
						}
					}
					m.state.messages = append(m.state.messages, messageState{time.Now().Format(time.Kitchen), "anon", m.input.Value()})
					m.input.Reset()
		}
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() string {

	timeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	usernameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("203")).Bold(true)

	var messageString string

	for _, m := range m.state.messages{
		messageString += fmt.Sprintf("%s %s %s %s\n", timeStyle.Render(m.timestamp), timeStyle.Render("-"), usernameStyle.Render("<<" + m.username + ">>"), m.message)
	}

	m.chatbox.SetContent(messageString)

	return fmt.Sprintf("%s\n%s", m.chatbox.View(), lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Render(m.input.View()))
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
	
	defer p.Quit()
}
