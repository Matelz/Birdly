package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var sentences = []string{
	"The quick brown fox jumps over the lazy dog.",
	"How much wood would a woodchuck chuck if a woodchuck could chuck wood?",
	"She sells sea shells by the sea shore.",
	"Peter Piper picked a peck of pickled peppers.",
	"The cat in the hat sat on the mat.",
	"A big black bug bit a big black dog on his big black nose.",
	"Betty Botter bought some butter but she said the butter’s bitter.",
	"Six sticky skeletons stacked skulls on a steep slope.",
	"Unique New York, Unique New York, Unique New York.",
	"Rubber baby buggy bumpers.",
}

type level struct {
	level  	int
	phrase 	string
	points 	int
	width 	int
	height 	int
	answerField textinput.Model
	startTime time.Time
}

var style = lipgloss.NewStyle().
			Bold(true).
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(lipgloss.Color("32")).
			Width(80)

var phraseStyle = lipgloss.NewStyle().Bold(true)

func initialLevel() level {
	ti := textinput.New()
	ti.Placeholder = "Type phrase"
	ti.Focus()
	ti.CharLimit = 450

	return level{
		level:  1,
		phrase: sentences[0],
		points: 0,
		answerField: ti,
	}
}

func (l level) checkAnswer(answer string) level{
	if l.phrase == answer{
		l.level += 1
		l.phrase = sentences[rand.Intn(len(sentences))]
		l.answerField.Reset()
		
		elapsed := int(math.Ceil(time.Since(l.startTime).Seconds()))

		l.points += int(math.Ceil(100 - 8 * float64(elapsed)))
		l.startTime = time.Time{}

		style.BorderForeground(lipgloss.Color("10"))
	} else {
		style.BorderForeground(lipgloss.Color("9"))
	}

	return l
}

func (l level) Init() tea.Cmd{
	return nil
}

func (l level) Update(msg tea.Msg) (tea.Model, tea.Cmd){
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		l.width = msg.Width
		l.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return l, tea.Quit
		
		case "enter":
			newLevel := l.checkAnswer(l.answerField.Value())
			return newLevel, cmd
		}
		
		if l.startTime.IsZero(){
			l.startTime = time.Now()
		}
	}
	
	l.answerField, cmd = l.answerField.Update((msg))
	return l, cmd
}

func (l level) View() string{
	return lipgloss.Place(l.width, l.height, lipgloss.Center, lipgloss.Center, lipgloss.JoinVertical(lipgloss.Center, lipgloss.JoinHorizontal(lipgloss.Position(0), fmt.Sprintf("Level %d", l.level), fmt.Sprintf(" - %d points", l.points)), phraseStyle.Render(l.phrase), style.Render(l.answerField.View()), "(press ctrl+c to exit)"))
}

func main() {
	p := tea.NewProgram(initialLevel())
    if _, err := p.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}