package main


import (
    "golang.org/x/term"
    "github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
    "fmt"
)
var (
    listSelectedStyle   = lipgloss.Color("4")
    titleStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Background(lipgloss.Color("4"))
    docStyle            = lipgloss.NewStyle().Margin(1, 2)
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type sessionState int 

const (
    viewLogin sessionState = iota
    viewUnlock
    viewLoading
    viewList
    viewClip
)
type MainModel struct {
    state sessionState
    ModelLogin tea.Model
    ModelUnlock tea.Model
    ModelLoading tea.Model
    ModelList tea.Model
    ModelClip tea.Model
    bwm BWManager
}

func NewMainModel() MainModel{
    var initialState sessionState
    h,v,_ := term.GetSize(0)
    if bwm.Status == "unauthenticated" {
        initialState = viewLogin
    } else {
        initialState = viewUnlock
    }
    return MainModel{
        state: initialState, 
        ModelLogin: NewUILogin(),
        ModelUnlock: NewUIUnlock(),
        ModelLoading: NewUILoading(),
        ModelList: NewUIList(h,v),
        ModelClip: NewUIClip(),
    }
}

func (m MainModel) Init() tea.Cmd {
    return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model,tea.Cmd) {
    var cmd tea.Cmd
	var cmds []tea.Cmd
    switch msg := msg.(type) {
    case LoadingLoginFailed:
        if msg.lt == login {
            m.state = viewLogin
        } else if msg.lt == unlock {
            m.state = viewUnlock
        }
	case LoginSubmit:
        m.state = viewLoading
    case ListSelectedEntry:
        m.state = viewClip
    case LoadingDone:
        m.state = viewList
    }
    switch m.state {
	case viewList:
        v,h,_ := term.GetSize(0)
        size := tea.WindowSizeMsg{Height:h,Width:v}
        m.ModelList.Update(size)
		newList, newCmd := m.ModelList.Update(msg)
		list, ok := newList.(UIList)
		if !ok {
			panic("could not perform assertion on List model")
		}
		m.ModelList = list
		cmd = newCmd
    case viewLoading:
        newLoading,newCmd := m.ModelLoading.Update(msg)
        loading, ok := newLoading.(UILoading)
		if !ok {
			panic("could not perform assertion on Loading model")
		}
        m.ModelLoading = loading
        cmd = newCmd
    case viewLogin:
        newLogin,newCmd := m.ModelLogin.Update(msg)
        login, ok := newLogin.(UILogin)
		if !ok {
			panic("could not perform assertion on Login model")
		}
        m.ModelLogin = login
        cmd = newCmd
    case viewUnlock:
        newUnlock,newCmd := m.ModelUnlock.Update(msg)
        unlock, ok := newUnlock.(UIUnlock)
		if !ok {
			panic("could not perform assertion on Unlock model")
		}
        m.ModelUnlock = unlock
        cmd = newCmd
    case viewClip:
        newClip,newCmd := m.ModelClip.Update(msg)
        clip, ok := newClip.(UIClip)
		if !ok {
			panic("could not perform assertion on Clip model")
		}
        m.ModelClip = clip
        cmd = newCmd
    }
    cmds = append(cmds, cmd)
    return m, tea.Batch(cmds...)
}

func (m MainModel) View() string {
	switch m.state {
	case viewLogin:
		return m.ModelLogin.View()
    case viewList:
        return m.ModelList.View()
    case viewLoading:
        return m.ModelLoading.View()
    case viewClip:
        return m.ModelClip.View()
    case viewUnlock:
        return m.ModelUnlock.View()
    default:
        return m.ModelLogin.View()
    }
}
