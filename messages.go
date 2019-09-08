package webmgmt

import "fmt"

type ServerMessage interface {
    Get() interface{}
}


type ServerMessageBase struct {
    Type string `json:"type"`
}

func (c *ServerMessageBase) Get() interface{} {
    return c
}

type TextMessage struct {
    ServerMessageBase
    Text  string `json:"text"`
    Color string `json:"color"`
}

func (c *TextMessage) Get() interface{} {
    return c
}


type RawTextMessage struct {
    ServerMessageBase
    Text     string   `json:"text"`
    Commands []string `json:"commands"`
}

func (c *RawTextMessage) Get() interface{} {
    return c
}

type Prompt struct {
    ServerMessageBase
    Prompt string `json:"prompt"`
}
func (c *Prompt) Get() interface{} {
    return c
}

type HistoryMode struct {
    ServerMessageBase
    Val bool `json:"val"`
}
func (c *HistoryMode) Get() interface{} {
    return c
}


type Authenticated struct {
    ServerMessageBase
    Val bool `json:"val"`
}
func (c *Authenticated) Get() interface{} {
    return c
}


type Echo struct {
    ServerMessageBase
    Val bool `json:"val"`
}
func (c *Echo) Get() interface{} {
    return c
}


type ClientMessage struct {
    Payload string `json:"payload"`
}


type Status struct {
    ServerMessageBase
    Text string `json:"text"`
}
func (c *Status) Get() interface{} {
    return c
}







func AppendText(text, color string) ServerMessage {
    something := &TextMessage{}
    something.Type = "text"
    something.Text = text
    something.Color = color
    return something
}


func AppendRawText(text string, commands []string) ServerMessage {
    something := &RawTextMessage{}
    something.Type = "rawtext"
    something.Text = text
    something.Commands = commands
    return something
}

func SetPrompt(prompt string) ServerMessage {
    something := &Prompt{}
    something.Type = "prompt"
    something.Prompt = prompt
    return something
}




func SetHistoryMode(val bool) ServerMessage {
    something := &HistoryMode{}
    something.Type = "history"
    something.Val = val
    return something
}


func SetAuthenticated(val bool) ServerMessage {
    something := &Authenticated{}
    something.Type = "authenticated"
    something.Val = val
    return something
}


func SetEchoOn(val bool) ServerMessage {
    something := &Echo{}
    something.Type = "echo"
    something.Val = val
    return something
}


func Cls() ServerMessage {
    something := &ServerMessageBase{}
    something.Type = "cls"
    return something
}

func SetStatus(text string) ServerMessage {
    something := &Status{}
    something.Type = "status"
    something.Text = text
    return something
}

func Eval(text string) ServerMessage {
    something := &Status{}
    something.Type = "eval"
    something.Text = text
    return something
}


func Link(url, text string) string {
    return fmt.Sprintf("<a href=\"%s\" target=\"_blank\">%s</a>", url, text)
}

func Color(color, text string) string {
    if text == "" {
        return text
    }
    return fmt.Sprintf("<span style=\"color:%s\">%s</span>", color, text)
}

func Image(width, height int, src, alt string) string {
    return fmt.Sprintf("<img width=\"%d\" height=\"%d\" src=\"%s\" alt=\"%s\"/>", width, height, src, alt)
}
