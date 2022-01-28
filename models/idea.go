package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	ideaStyle        = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).Padding(1).Width(40).Align(lipgloss.Center)
	titleStyle       = lipgloss.NewStyle().Bold(true).Underline(true)
	descriptionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#aaa")).Margin(1).Italic(true)
	likeStyle        = lipgloss.NewStyle().UnsetAlign().Align(lipgloss.Left).Background(lipgloss.Color("#CD5C5C")).
				Foreground(lipgloss.Color("#fff")).PaddingLeft(1).PaddingRight(1)
	tagStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#0f0"))
	orderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)
	helpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).MarginTop(2)
	modelStyle = lipgloss.NewStyle().MarginLeft(2)
)

const API_URL = "https://what-to-code.com/api"

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
	index int
	page  int
	Order string
}

func NewIdeaModel() IdeaModel {
	return IdeaModel{
		index: 0,
		page:  0,
		list:  []Idea{},
		Order: "POPULAR",
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
		case "w", "up", "a", "k":
			if i.index > 0 {
				i.index -= 1
			}
		case "s", "down", "d", "j":
			if i.index < len(i.list) {
				i.index += 1
			}

			if i.index == len(i.list) {
				i.getMoreIdeas()
			}

		case " ", "enter":
			i.Like()

		}
	}

	return i, nil
}

func (i IdeaModel) renderTags() string {
	s := ""
	for _, tag := range i.currentIdea().Tags {
		s += "#" + tag.Value
	}

	return s
}

func (i *IdeaModel) getMoreIdeas() {
	var newIdeas []Idea

	res, _ := http.Get(i.getListUrl())

	body, _ := ioutil.ReadAll(res.Body)
	err := json.Unmarshal(body, &newIdeas)

	if err != nil {
		fmt.Print(err)
	}

	defer res.Body.Close()

	i.list = append(i.list, newIdeas...)
	i.page += 1
}

func (i IdeaModel) getListUrl() string {
	return fmt.Sprintf("%s/ideas?order=%s&page=%d", API_URL, i.Order, i.page)
}

func text_default(text string, def string) string {
	if len(text) == 0 {
		return "There is no " + def
	}
	return text
}

func (i *IdeaModel) Clear() {
	i.list = nil
}

func (i IdeaModel) Like() {
	http.Post(fmt.Sprintf("%s/ideas/%d/like", API_URL, i.currentIdea().Id), "", nil)

	*i.currentIdea() = i.getIdea(i.currentIdea().Id)
}

func (i IdeaModel) currentIdea() *Idea {
	return &i.list[i.index]
}

func (i IdeaModel) getIdea(id int) Idea {
	var idea Idea

	res, _ := http.Get(fmt.Sprintf("%s/ideas/%d", API_URL, i.currentIdea().Id))

	body, _ := ioutil.ReadAll(res.Body)
	err := json.Unmarshal(body, &idea)

	if err != nil {
		fmt.Print(err)
	}

	defer res.Body.Close()

	return idea
}

func (i IdeaModel) View() string {
	s := titleStyle.Render(i.currentIdea().Title) + "\n"

	s += descriptionStyle.Render(text_default(i.currentIdea().Description, "description")) + "\n"
	s += likeStyle.Render(fmt.Sprintf("♥ %d", i.currentIdea().Likes)) + "\n\n"
	s += tagStyle.Render(text_default(i.renderTags(), "tag"))

	return modelStyle.Render(orderStyle.Render(i.Order) + "\n" + ideaStyle.Render(s) +
		helpStyle.Render("↑/k/w previous idea • ↓/j/s next idea • space like/dislike • q/esc quit"))
}
