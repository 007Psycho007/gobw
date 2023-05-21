package main

import (
    "fmt"
    "os"
    "github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ListSelectedEntry struct {
    item list.Item
    prop property
}

func SelectListSelectedEntry(id list.Item,prop property) tea.Cmd {
    return func() tea.Msg {
        return ListSelectedEntry{id,prop}
    }
}

type UIList struct {
	list list.Model
    selected int
}

func NewUIList(h,v int) UIList{
    d := list.NewDefaultDelegate()
    d.Styles.SelectedTitle = d.Styles.SelectedTitle.Foreground(listSelectedStyle).BorderLeftForeground(listSelectedStyle)
    d.Styles.SelectedDesc = d.Styles.SelectedTitle.Copy()
    width, height := docStyle.GetFrameSize()
    l := list.New(nil, d,h-width,v-height)
    l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
            listKeys.User,
            listKeys.Pass,
		}
    }
    l.Styles.Title = titleStyle

    return UIList{list: l}
}

func (m UIList) Init() tea.Cmd {
	return nil
}

func (m *UIList) GetEntries() {
    bwi,err := bwm.GetList()
    if err != nil {
        fmt.Println("Error: " + err.Error())
        os.Exit(1)
    }
    m.list.Title = " Bitwarden Vault of: " + bwm.UserMail + " "
    m.list.SetItems(bwi)
}
func (m UIList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
    case LoadingDone:
        m.GetEntries()        
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
        if msg.String() == "enter" {
			return m, SelectListSelectedEntry(m.list.SelectedItem(),copyPassword)
        }
        if msg.String() == "alt+enter" {
			return m, SelectListSelectedEntry(m.list.SelectedItem(),copyUsername)
        }
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
        return m , tea.ClearScreen
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m UIList) View() string {
	return docStyle.Render(m.list.View())
}

