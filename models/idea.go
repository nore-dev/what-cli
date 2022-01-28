package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	ideaStyle        = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).Padding(1).Width(40).Align(lipgloss.Center)
	titleStyle       = lipgloss.NewStyle().Bold(true).Underline(true)
	descriptionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#aaaaaa")).Margin(1).Italic(true)
	likeStyle        = lipgloss.NewStyle().UnsetAlign().Align(lipgloss.Left).Background(lipgloss.Color("#ff0000")).
				Foreground(lipgloss.Color("#ffffff")).PaddingLeft(1).PaddingRight(1)
	tagStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00"))
)

type Tag struct {
	value string
}

type Idea struct {
	description string
	id          int
	likes       int
	title       string
	tags        []Tag
}

type IdeaModel struct {
	list  []Idea
	url   string
	index int
	page  int
}

func NewIdeaModel() IdeaModel {
	return IdeaModel{
		index: 0,
		page:  0,
		list: []Idea{
			{
				title:       "Title",
				description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
				tags:        []Tag{{value: "tag"}},
				likes:       31,
				id:          0,
			},
			{
				title:       "Title 2",
				description: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAaaa",
				tags:        []Tag{{value: "awesome"}},
				likes:       13,
				id:          0,
			},
		},
		url: "https://what-to-code.com/api/ideas?order=POPULAR&page=0",
	}
}

func (i IdeaModel) Init() tea.Cmd {
	return nil
}

func (i IdeaModel) Update(msg tea.Msg) (IdeaModel, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "w", "up", "a":
			if i.index > 0 {
				i.index -= 1
			}
		case "s", "down", "d":
			if i.index < len(i.list)-1 {
				i.index += 1
			}
		}
	}

	return i, nil
}

func (i IdeaModel) renderTags() string {
	s := ""
	for _, tag := range i.list[i.index].tags {
		s += "#" + tag.value
	}

	return tagStyle.Render(s)
}

func (i IdeaModel) View() string {
	s := titleStyle.Render(i.list[i.index].title) + "\n"

	s += descriptionStyle.Render(i.list[i.index].description) + "\n"
	s += likeStyle.Render(fmt.Sprintf("♥ %d", i.list[i.index].likes)) + "\n\n"
	s += tagStyle.Render(i.renderTags())

	return ideaStyle.Render(s)
}
