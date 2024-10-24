package types

import tea "github.com/charmbracelet/bubbletea"

// type model struct {
// 	// Set a bool so we know when the window was initialized
// 	Ready   bool
// 	State   RoomState
// 	Chatbox viewport.Model
// 	Input   textinput.Model
// }

const (
	// Normal message
	NormalMessage = 1
	// Info message
	InfoMessage = 2
	// Command message
	CommandMessage = 3
)

type Argument struct {
	Name     string
	Required bool
	// Default value of the argument
	Value any
}

type Command struct {
	Command string
	// Hint to be shown to the user when he makes a mistake eg: not enough arguments
	Description string
	Args        []Argument
	// Function to be called when the command is called
	Function func(m tea.Model, args ...interface{}) interface{}
}

type MessageState struct {
	Timestamp string
	Username  string
	Message   string
	Type      int
}

type RoomState struct {
	Messages []MessageState
	Clients  []string
}
