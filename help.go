package main

import (
    "github.com/charmbracelet/bubbles/key"
)
type listKeyMap struct {
    User  key.Binding
    Pass  key.Binding
}

var listKeys = listKeyMap{
    User: key.NewBinding(
        key.WithKeys("alt+enter"),
        key.WithHelp("Alt+Enter","Copy Username"),
        ),
    Pass: key.NewBinding(
        key.WithKeys("enter"),
        key.WithHelp("Enter","Copy Password"),
        ),
}
