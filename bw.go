package main

import (
    "errors"
	"fmt"
	"os/exec"
    "encoding/json"
    "github.com/charmbracelet/bubbles/list"
)

type BWListLogin struct {
    UserName string `json:"username"`
    Password string `json:"password"`
}

type BWItems struct {
    ID string `json:"id"`
    ObjectName string `json:"name"`
    Login BWListLogin `json:"login"`
}

func (bwi *BWItems) GetData() BWListItem {
    var data BWListItem
    data.ID = bwi.ID
    data.ObjectName = bwi.ObjectName
    data.UserName = bwi.Login.UserName
    return data
}

type BWListItem struct {
    ID string 
    ObjectName string
    UserName string
}
func (bwl BWListItem) Title() string       { return bwl.ObjectName }
func (bwl BWListItem) Description() string { return bwl.UserName }
func (bwl BWListItem) FilterValue() string { return bwl.ObjectName }

type BWManager struct {
    ServerUrl string `json:"serverUrl"`
    LastSync string `json:"lastSync"`
    UserMail string `json:"userEmail"`
    UserId string `json:"userId"`
    Status string `json:"status"`
    items []BWItems
    token string
}

func NewBWManager() BWManager {
    var bwm BWManager
    bwm.UpdateStatus()
    return bwm
}

func (bwm *BWManager) Login(un string, pw string) error {
    if bwm.Status != "unauthenticated" {
        fmt.Println("Already Logged in")
        return nil
    }
    out,err := exec.Command("bw", "login", un, pw, "--raw").Output()
    if err != nil {
        return errors.New(err.Error())
    }
    bwm.UpdateStatus()
    (*bwm).token = string(out)
    return nil
} 

func (bwm *BWManager) Unlock(pw string) error {
    if bwm.Status == "unauthenticated" {
        return errors.New("Not Logged in")
    }
    out,err := exec.Command("bw", "unlock", pw, "--raw").Output()
    if err != nil {
        return errors.New(err.Error())
    }
    bwm.UpdateStatus()
    (*bwm).token = string(out)
    return nil
} 

func (bwm *BWManager) Logout() error {
    if bwm.Status == "unauthenticated" {
        return errors.New("Not Logged in")
    }
    _,err := exec.Command("bw", "logout").Output()
    if err != nil {
        return errors.New(err.Error())
    }
    (*bwm).token = ""
    bwm.UpdateStatus()
    return nil
}

func (bwm *BWManager) UpdateStatus() error{
    out, err := exec.Command("bw", "status").Output()
    if err != nil {
        return errors.New(err.Error())
    }
    json.Unmarshal(out,&bwm)
    return nil
}

func (bwm *BWManager) UpdateList() error{
    if bwm.Status == "unauthenticated"  {
        fmt.Println("Not Logged In")
        return errors.New("Not Logged In")
    }
    out,err := exec.Command("bw", "list", "items", "--session", bwm.token).Output()
    if err != nil {
        fmt.Println("Error: " + err.Error())
        return errors.New(err.Error())
    }
    json.Unmarshal(out,&bwm.items)
    
    return nil
}

func (bwm *BWManager) GetList() ([]list.Item,error){
    if bwm.Status == "unauthenticated" {
        return nil,errors.New("Not Logged In")
    }
    data := []list.Item{}
    for _,v := range(bwm.items) {
        data = append(data,v.GetData())
    }
    return data,nil
}

func (bwm *BWManager) GetPassword(id string) (string,error){
    for _,v := range(bwm.items) {
        if v.ID == id {
            return v.Login.Password,nil
        }
    }
    return "",errors.New("No entry matching ID found.")
}
