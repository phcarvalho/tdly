package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	Editing string = "editing"
)

type item struct {
	text        string    `json:"text"`
	createdAt   time.Time `json:"createdAt,omitempty"`
	completedAt time.Time `json:"completedAt,omitempty"`
}

type data struct {
	items []item `json:"items,omitempty"`
}

type model struct {
	keys       keyMap
	editKeys   editKeyMap
	help       help.Model
	inputStyle lipgloss.Style
	textInput  textinput.Model
	state      string
	items      []item
	cursor     int
	focused    bool
}

func main() {
	p := tea.NewProgram(initialModel())

	m, err := p.Run()
	if err != nil {
		fmt.Printf("An error has occured: %v", err)
		os.Exit(1)
	}

	if m, ok := m.(model); ok {
		saveToFile(m.items)
	}
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "To do item text"
	ti.CharLimit = 255
	ti.Cursor.Blink = true
	ti.Prompt = ""

	items := readFromFile()

	return model{
		keys:       keys,
		editKeys:   editKeys,
		help:       help.New(),
		inputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")),
		textInput:  ti,
		items:      items,
		cursor:     0,
		focused:    false,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Teardown() tea.Cmd {
	saveToFile(m.items)

	return nil
}

func readFromFile() []item {
	filename := "data"

	file, err := os.Open(filename)
	if err != nil {
		return []item{}
	}
	defer file.Close()

  scanner := bufio.NewScanner(file)
  scanner.Split(bufio.ScanLines)

	items := []item{}

  for scanner.Scan() {
    arr := strings.Split(scanner.Text(), ";")

    createdAt, err := time.Parse(time.RFC3339Nano, arr[1])
    if err != nil {
      createdAt = time.Now()
    }

    completedAt, err := time.Parse(time.RFC3339Nano, arr[2])
    if err != nil {
      completedAt = time.Time{}
    }

    items = append(items, item{
      text: arr[0],
      createdAt: createdAt,
      completedAt: completedAt,
    })
  }

	return items
}

func saveToFile(items []item) {
	filename := "data"
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

  s := ""

  for _, item := range(items) {
    completedAt := ""

    if !item.completedAt.IsZero() {
      completedAt = item.completedAt.Format(time.RFC3339Nano)
    }

    s += fmt.Sprintf("%s;%s;%s\n", 
      item.text, 
      item.createdAt.Format(time.RFC3339Nano), 
      completedAt,
    )
  }

	_, err = file.WriteString(s)
	if err != nil {
		fmt.Println("Error writing data to file:", err)
		return
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.state == Editing {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.editKeys.Cancel):
				m.textInput.SetValue("")
				m.textInput.Blur()
				m.state = ""

				if m.items[m.cursor].createdAt.IsZero() {
					m.items = m.items[:m.cursor]

					m.cursor--
				}
			case key.Matches(msg, m.editKeys.Enter):
				text := m.textInput.Value()

				m.textInput.SetValue("")
				m.textInput.Blur()
				m.state = ""

				m.items[m.cursor] = item{text: text, createdAt: time.Now()}
			}

      if msg.String() == ";" {
        return m, cmd
      }

			m.textInput, cmd = m.textInput.Update(msg)

			return m, cmd
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case key.Matches(msg, m.keys.MoveUp):
			if m.cursor <= 0 {
				break
			}

			item := m.items[m.cursor]

			m.items[m.cursor] = m.items[m.cursor-1]
			m.items[m.cursor-1] = item

			m.cursor--
		case key.Matches(msg, m.keys.MoveDown):
			if m.cursor == -1 || m.cursor >= len(m.items)-1 {
				break
			}

			item := m.items[m.cursor]

			m.items[m.cursor] = m.items[m.cursor+1]
			m.items[m.cursor+1] = item

			m.cursor++
		case key.Matches(msg, m.keys.Add):
			m.state = Editing
			m.textInput.Focus()

			m.items = append(m.items, item{})
			m.cursor = len(m.items) - 1
		case key.Matches(msg, m.keys.Edit):
			m.state = Editing
			m.textInput.Focus()
			m.textInput.SetValue(m.items[m.cursor].text)
		case key.Matches(msg, m.keys.Complete):
			if m.cursor == -1 {
				break
			}

			item := &m.items[m.cursor]

			if item.completedAt.IsZero() {
				item.completedAt = time.Now()
			} else {
				item.completedAt = time.Time{}
			}
		case key.Matches(msg, m.keys.Delete):
			if m.cursor == -1 {
				break
			}

			m.items = append(m.items[:m.cursor], m.items[m.cursor+1:]...)

			if m.cursor >= len(m.items) || m.cursor == 0 {
				m.cursor--
			}
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "\n"

	if len(m.items) == 0 {
		s = "\nPress \"a\" to start adding items\n"
	}

	for i, item := range m.items {
		cursor := " "

		if i == m.cursor {
			cursor = ">"
		}

		completed := "[ ]"

		if !item.completedAt.IsZero() {
			completed = "[x]"
		}

		itemText := item.text

		if m.cursor == i && m.state == Editing {
			itemText = m.textInput.View()

			if item.createdAt.IsZero() {
				completed = " + "
			}
		}

		lineText := fmt.Sprintf("%s %s %s", cursor, completed, itemText)

		if m.cursor == i {
			lineText = m.inputStyle.Render(lineText)
		}

		s += lineText + "\n"
	}

	help := m.help.View(m.keys)

	if m.state == Editing {
		help = m.help.View(m.editKeys)
	}

	s += "\n" + help

	return s
}
