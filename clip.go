package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
)

type tickMsg time.Time

func tick() tea.Msg {
	time.Sleep(time.Second)
	return tickMsg{}
}

type property int

const (
    copyUsername property = iota
    copyPassword 
)


type UIClip struct {
    timer int 
    object string
    prop property
}

func NewUIClip() tea.Model{
    return UIClip{
        timer: 10,
    }
}

func (c UIClip) Init() tea.Cmd {
    return nil
}

func (c UIClip) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case ListSelectedEntry:
        item,ok := msg.item.(BWListItem)
        if !ok {
            panic("Could not get BWListItem")
        }
        switch msg.prop {
        case copyPassword:
            data,err := bwm.GetPassword(item.ID)       
            if err != nil {
                panic("Error getting Password")
            }
            clipboard.WriteAll(data)
            c.prop = copyPassword
            data = ""
        case copyUsername:
            data := item.UserName
            clipboard.WriteAll(data)
            c.prop = copyUsername
            data = ""
        }
        c.object = item.ObjectName
        return c, tick
    case tickMsg:
        c.timer--
        if c.timer <= 0 {
            clipboard.WriteAll("")
            c.timer = 10
			return c, SelectLoadingDone()
		}
		return c, tick
    case tea.KeyMsg:
		if msg.String() == "q" {
            clipboard.WriteAll("")
            return c, SelectLoadingDone()
        }
    }
    var cmd tea.Cmd
    return c, cmd
}
func (c UIClip) View() string {
	var b strings.Builder
    var p string
    b.WriteString(titleStyle.Render(" Bitwarden TUI "))
    b.WriteString("\n\n")
    b.WriteString("Type: ")
    switch c.prop{
    case copyPassword:
        p = "Password"
    case copyUsername:
        p = "Username"
    }

    b.WriteString(focusedStyle.Render(p))
    b.WriteString("\n")
    b.WriteString("Object: ")
    b.WriteString(focusedStyle.Render(c.object))
    b.WriteString("\n\nClipboard will be cleared in " + focusedStyle.Render(fmt.Sprint(c.timer)) + " seconds.")
    b.WriteString("\nPress " + focusedStyle.Render("q")+ " to delete Clipboard now and return to Vault.")
    return docStyle.Render(b.String())
}

