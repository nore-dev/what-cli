package models

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	submitStyle     = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).Width(35).Height(10).Padding(0, 1, 1)
	backgroundStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).Padding(0, 1)
)

type SubmitModel struct {
	fields []textinput.Model
	index  int
}

func NewSubmitModel() SubmitModel {
	title := textinput.New()
	title.Placeholder = "Title"
	title.CharLimit = 100
	title.Focus()

	description := textinput.New()
	description.Placeholder = "Description"
	description.CharLimit = 1000

	tags := textinput.New()
	tags.Placeholder = "Tags"

	return SubmitModel{
		fields: []textinput.Model{
			title,
			description,
			tags,
		},
		index: 0,
	}
}

func (s SubmitModel) Init() tea.Cmd {
	return textinput.Blink
}

func (s *SubmitModel) Clear() {
	for i, _ := range s.fields {
		s.fields[i].SetValue("")
	}

	s.index = 0
}

func (s SubmitModel) Update(msg tea.Msg) (SubmitModel, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if s.index > 0 {
				s.index -= 1
				s.fields[s.index].Focus()
			}
		case "down":
			if s.index < len(s.fields)-1 {
				s.index += 1
				s.fields[s.index].Focus()
			}
		}
	}

	s.fields[s.index], _ = s.fields[s.index].Update(msg)

	return s, nil
}

func (s SubmitModel) View() string {
	return backgroundStyle.Render("SUBMIT") + "\n" +
		submitStyle.Render(s.fields[s.index].View()) +
		helpStyle.Render("↑ previous field • ↓/j/s next field • enter submit • q/esc go back")
}
