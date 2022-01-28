package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

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
	Value string `json:"value,omitempty"`
}

type Idea struct {
	Description string `json:"description,omitempty"`
	Id          int    `json:"id,omitempty"`
	Likes       int    `json:"likes,omitempty"`
	Title       string `json:"title,omitempty"`
	Tags        []Tag  `json:"tags,omitempty"`
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
		list:  []Idea{},
		url:   "https://what-to-code.com/api/ideas?order=POPULAR&page=0",
	}
}

func (i IdeaModel) Init() tea.Cmd {
	return nil
}

func (i IdeaModel) Update(msg tea.Msg) (IdeaModel, tea.Cmd) {

	// First Time
	if len(i.list) == 0 {
		i.getMoreIdeas()
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "w", "up", "a":
			if i.index > 0 {
				i.index -= 1
			}
		case "s", "down", "d":
			if i.index < len(i.list) {
				i.index += 1
			}

			if i.index == len(i.list) {
				i.getMoreIdeas()
			}

		}
	}

	return i, nil
}

func (i IdeaModel) renderTags() string {
	s := ""
	for _, tag := range i.list[i.index].Tags {
		s += "#" + tag.Value
	}

	return tagStyle.Render(s)
}

func (i *IdeaModel) getMoreIdeas() {
	var newIdeas []Idea

	res, _ := http.Get(i.url)

	body, _ := ioutil.ReadAll(res.Body)
	err := json.Unmarshal(body, &newIdeas)

	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	i.list = append(i.list, newIdeas...)
}

func (i IdeaModel) View() string {
	s := titleStyle.Render(i.list[i.index].Title) + "\n"

	s += descriptionStyle.Render(i.list[i.index].Description) + "\n"
	s += likeStyle.Render(fmt.Sprintf("♥ %d", i.list[i.index].Likes)) + "\n\n"
	s += tagStyle.Render(i.renderTags())

	return ideaStyle.Render(s)
}
