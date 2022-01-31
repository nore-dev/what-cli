package main

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nore-dev/what-cli/models"
)

type App struct {
	list        list.Model
	ideaModel   models.IdeaModel
	submitModel models.SubmitModel
	selected    bool
	isIdeaPage  bool
}

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

type Page string

func (p Page) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int  { return 1 }
func (d itemDelegate) Spacing() int { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Page)
	if !ok {
		return
	}

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprintf(w, fn(string(i)))
}

func NewApp() App {
	items := []list.Item{
		Page("Popular"),
		Page("Recent"),
		Page("Rising"),
		Page("Oldest"),
		Page("Random"),
		Page("Submit"),
	}

	app := App{
		list:        list.New(items, itemDelegate{}, 0, 0),
		selected:    false,
		isIdeaPage:  false,
		ideaModel:   models.NewIdeaModel(),
		submitModel: models.NewSubmitModel(),
	}

	app.list.Title = "What CLI"

	app.list.Styles.PaginationStyle = paginationStyle
	app.list.Styles.HelpStyle = helpStyle

	app.list.SetHeight(14)
	app.list.SetShowStatusBar(false)
	app.list.SetFilteringEnabled(false)
	app.list.SetShowFilter(false)

	return app
}

func (a App) Init() tea.Cmd {
	return nil
}

func (a App) getOrder() string {
	switch a.list.SelectedItem() {
	case Page("Popular"):
		return "POPULAR"
	case Page("Oldest"):
		return "OLDEST"
	case Page("Rising"):
		return "RISING"
	case Page("Recent"):
		return "RECENT"
	case Page("Random"):
		return "RANDOM"
	}

	return ""
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			if a.selected {
				a.selected = false
			} else {
				return a, tea.Quit
			}
		case "enter", " ":
			switch a.list.SelectedItem() {
			case Page("Submit"):
				a.isIdeaPage = false
				a.selected = true
			default:
				if !a.selected {
					a.ideaModel.Clear()
					a.ideaModel.Order = a.getOrder()
				}

				a.selected = true
				a.isIdeaPage = true
			}
		}
	}

	var cmd tea.Cmd
	if a.selected {
		if a.isIdeaPage {
			a.ideaModel, _ = a.ideaModel.Update(msg)
		} else {
			a.submitModel, cmd = a.submitModel.Update(msg)
		}
	} else {
		a.list, _ = a.list.Update(msg)
	}

	return a, cmd
}

func (a App) View() string {
	v := a.list.View()

	if a.selected {
		if a.isIdeaPage {
			v = a.ideaModel.View()
		} else {
			v = a.submitModel.View()
		}
	}

	return v
}
