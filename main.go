package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
}

type model struct {
	tasks               []task
	cursor              int
	renderAddItemPrompt bool
	addItemPrompt       textinput.Model
}

type task struct {
	description string
	selected    bool
}

func initialModel() *model {
	input := textinput.New()
	return &model{
		renderAddItemPrompt: false,
		addItemPrompt:       input,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.renderAddItemPrompt {
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "esc":
				m.addItemPrompt.SetValue("")
				m.renderAddItemPrompt = false
				return m, nil
			case "enter":
				task := task{description: m.addItemPrompt.Value(), selected: false}
				m.tasks = append(m.tasks, task)
				m.addItemPrompt.SetValue("")
				m.renderAddItemPrompt = false
				return m, nil
			}
			m.addItemPrompt, cmd = m.addItemPrompt.Update(msg)
			return m, cmd
		} else {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "a":
				m.renderAddItemPrompt = true
				m.addItemPrompt.Focus()
			case "d":
				if len(m.tasks) > 0 {
					m.tasks = append(m.tasks[:m.cursor], m.tasks[m.cursor+1:]...)
				}
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.tasks)-1 {
					m.cursor++
				}
			case "enter", " ":
				m.tasks[m.cursor].selected = !m.tasks[m.cursor].selected
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "What should we buy at the market?\n\n"
	for i, task := range m.tasks {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if task.selected {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, task.description)
	}

	if m.renderAddItemPrompt {
		s += "\n------------------------------------------------------------------------------------------------------\n"
		s += m.addItemPrompt.View()
		s += "\n------------------------------------------------------------------------------------------------------\n"
	} else {
		s += "\n[a]dd, [d]elete, [q]uit\n"
	}

	return s
}
