package main

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type editKeyMap struct {
	Enter  key.Binding
	Cancel key.Binding
}

func (k editKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Enter, k.Cancel}
}

func (k editKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Enter, k.Cancel}}
}

var editKeys = editKeyMap{
	Enter: key.NewBinding(
		key.WithKeys(tea.KeyEnter.String()),
		key.WithHelp("enter", "confirm"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "cancel"),
	),
}

type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	MoveUp   key.Binding
	MoveDown key.Binding
	Add      key.Binding
	Delete   key.Binding
	Edit     key.Binding
	Complete key.Binding
	Help     key.Binding
	Quit     key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	empty := key.NewBinding(
		key.WithKeys(""),
		key.WithHelp("", ""),
	)
	return [][]key.Binding{
		{k.Up, k.Down, k.MoveUp, k.MoveDown},
		{k.Add, k.Edit, k.Delete, k.Complete},
		{empty, empty, k.Help, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "cursor up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "cursor down"),
	),
	MoveUp: key.NewBinding(
		key.WithKeys("shift+up", "K"),
		key.WithHelp("shift+↑/K", "move up"),
	),
	MoveDown: key.NewBinding(
		key.WithKeys("shift+down", "J"),
		key.WithHelp("shift+↓/J", "move down"),
	),
	Add: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	Delete: key.NewBinding(
		key.WithKeys("D"),
		key.WithHelp("D", "delete"),
	),
	Complete: key.NewBinding(
		key.WithKeys(tea.KeySpace.String()),
		key.WithHelp("space", "complete"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
