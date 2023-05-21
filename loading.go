package main

import (
	"errors"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type loginType int
const (
    login loginType = iota
    unlock
)

type errMsg error

type LoadingLoginFailed struct {
    lt loginType
}

func SelectLoadingFailed(lt loginType) tea.Cmd {
    return func() tea.Msg {
        return LoadingLoginFailed{lt}
    }
}

type LoadingDone struct {
}

func SelectLoadingDone() tea.Cmd {
    return func() tea.Msg {
        return LoadingDone{}
    }
}
type UILoading struct {
    un      string
    pw      string
    lt    loginType
}

func NewUILoading() UILoading {
	return UILoading{}
}

func (m UILoading) Init() tea.Cmd {
	return nil
}

func (m UILoading) Login() error {
    switch m.lt {
    case login:
        err := bwm.Login(m.un,m.pw)
        if err != nil {
            m.un,m.pw = "",""
            return errors.New("Login Failed")
        }
    case unlock:
        err := bwm.Unlock(m.pw)
        if err != nil {
            m.un,m.pw = "",""
            return errors.New("Unlock Failed")
        }
    }
    m.un,m.pw = "",""
    bwm.UpdateList()
    return nil
}

func (m UILoading) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
    case LoginSubmit:
        m.un = msg.un
        m.pw = msg.pw
        m.lt = msg.lt
        return m, tick
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		default:
			return m, nil
		}

	case errMsg:
		return m, tea.Quit

	default:
        err := m.Login()
        if err != nil {
            return m,SelectLoadingFailed(m.lt)
        }
		return m, SelectLoadingDone()
	}
}

func (m UILoading) View() string {
    var b strings.Builder
    b.WriteString(titleStyle.Render(" Bitwarden TUI "))
    b.WriteString("\n\n Logging in. Please wait\n\n")
	return docStyle.Render(b.String())
}
