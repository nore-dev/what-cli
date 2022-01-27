package main

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type App struct {
	list list.Model
}

var (
	// titleStyle        = lipgloss.NewStyle().MarginLeft(2)
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
		Page("Liked"),
	}

	app := App{list: list.New(items, itemDelegate{}, 0, 0)}
	app.list.Title = "What CLI"

	// app.list.Styles.Title = titleStyle
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

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return a, tea.Quit
		}
	}

	var cmd tea.Cmd
	a.list, cmd = a.list.Update(msg)

	return a, cmd
}

func (a App) View() string {
	return a.list.View()
}
