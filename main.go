package main

import (
	"birdly/styles"
	"birdly/types"
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

type model struct {
	// Set a bool so we know when the window was initialized
	Ready   bool
	State   types.RoomState
	Chatbox viewport.Model
	Input   textinput.Model
}

var commands = []types.Command{
	{
		Command:     "pm",
		Description: "Send private message to user (:pm <to> <message>)",
		Args: []types.Argument{
			{
				Name:     "to",
				Required: true,
				Value:    "",
			},
		},
		Function: func(m tea.Model, args ...interface{}) interface{} {
			return "Private message sent\n"
		},
	},
}

func initialModel() model {
	chatbox := viewport.New(10, 10)
	chatbox.Style = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())

	input := textinput.New()
	input.Placeholder = "Type message"
	input.Focus()

	return model{
		Ready: false,
		// Initialize room state with no messages and no clients
		State: types.RoomState{
			Messages: []types.MessageState{},
			Clients:  []string{},
		},
		Chatbox: chatbox,
		Input:   input,
	}
}

func (m model) displayMessage(message string, _type int) tea.Model {
	m.State.Messages = append(m.State.Messages, types.MessageState{Timestamp: time.Now().Format(time.Kitchen), Username: "anon", Message: message, Type: _type})
	m.Input.Reset()

	return m
}

func (m model) handleCommand() tea.Model {
	// Command format (:<command> <arg1> <arg...>)
	r := regexp.MustCompile(`:[a-zA-Z]{1,}`)
	message := m.Input.Value()
	var returnModel tea.Model

	if cmd := r.FindString(message); cmd != "" {
		for _, c := range commands {
			cmd = strings.ReplaceAll(cmd, ":", "")
			if c.Command == cmd {
				// Retrieve args from command
				args := strings.Split(message, " ")

				if len(args)-2 < len(c.Args) {
					returnModel = m.displayMessage(fmt.Sprintf("Not enough arguments\n%s\n", c.Description), types.CommandMessage)
					return returnModel
				}

				args = args[1 : len(args)-1]

				returnModel = m.displayMessage(fmt.Sprintf("%s", c.Function(m)), types.CommandMessage)
				return returnModel
			}
		}
	}

	returnModel = m.displayMessage("Command not found", 2)
	return returnModel
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.Ready {
			m.Chatbox.Width = msg.Width
			m.Chatbox.Height = msg.Height - 3

			m.Input.Width = msg.Width - 5
			m.Ready = true
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.Input.Value() == "" {
				return m, cmd
			}

			if string(m.Input.Value()[0]) == ":" {
				return m.handleCommand(), cmd
			}
			return m.displayMessage(m.Input.Value(), types.NormalMessage), cmd
		}
	}

	m.Input, cmd = m.Input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var messageString string

	for _, m := range m.State.Messages {
		messageString += styles.MessageFormatter(m)
	}

	// Set the content of the chatbox and an empty line so the last message is not cut off
	m.Chatbox.SetContent(messageString + "\n")
	m.Chatbox.GotoBottom()

	return fmt.Sprintf("%s\n%s", m.Chatbox.View(), lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Render(m.Input.View()))
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	defer p.Quit()
}
