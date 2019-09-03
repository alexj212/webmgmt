package webmgmt

import "fmt"

type ServerMessage struct {
    Type string `json:"type"`
}

type TextMessage struct {
    ServerMessage
    Text  string `json:"text"`
    Color string `json:"color"`
}

func AppendText(text, color string) *TextMessage {
    something := &TextMessage{}
    something.Type = "text"
    something.Text = text
    something.Color = color
    return something
}

type RawTextMessage struct {
    ServerMessage
    Text     string   `json:"text"`
    Commands []string `json:"commands"`
}

func AppendRawText(text string, commands []string) *RawTextMessage {
    something := &RawTextMessage{}
    something.Type = "rawtext"
    something.Text = text
    something.Commands = commands
    return something
}

type Prompt struct {
    ServerMessage
    Prompt string `json:"prompt"`
}

func SetPrompt(prompt string) *Prompt {
    something := &Prompt{}
    something.Type = "prompt"
    something.Prompt = prompt
    return something
}

type HistoryMode struct {
    ServerMessage
    Val bool `json:"val"`
}

func SetHistoryMode(val bool) *HistoryMode {
    something := &HistoryMode{}
    something.Type = "history"
    something.Val = val
    return something
}

type Authenticated struct {
    ServerMessage
    Val bool `json:"val"`
}

func SetAuthenticated(val bool) *Authenticated {
    something := &Authenticated{}
    something.Type = "authenticated"
    something.Val = val
    return something
}

type Echo struct {
    ServerMessage
    Val bool `json:"val"`
}

func SetEchoOn(val bool) *Echo {
    something := &Echo{}
    something.Type = "echo"
    something.Val = val
    return something
}

type ClientMessage struct {
    Payload string `json:"payload"`
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
