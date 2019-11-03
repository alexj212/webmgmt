package webmgmt

import (
	"fmt"
	"html"
)

// ServerMessage is the interface that all messages from the server to client will implement.
type ServerMessage interface {
	Get() interface{}
}

// ServerMessageBase is the base structure all server to client messages will use.
type ServerMessageBase struct {
	Type string `json:"type"`
}

// Get will return the ServerMessageBase
func (c *ServerMessageBase) Get() interface{} {
	return c
}

// TextMessaage is the struct for the server message that is sent to the client to tell the client to display text in the terminal window.
type TextMessage struct {
	ServerMessageBase
	Text  string `json:"text"`
	Color string `json:"color"`
}

func (c *TextMessage) Get() interface{} {
	return c
}

// RawTextMessage is the struct for the server message that is sent to the client to tell the client to display text as raw text in the terminal window.
type RawTextMessage struct {
	ServerMessageBase
	Text     string   `json:"text"`
}

func (c *RawTextMessage) Get() interface{} {
	return c
}



// Clickable is the struct for the server message that is sent to the client to tell the client to display clickable.
type Clickable struct {
	ServerMessageBase
	Commands []string `json:"commands"`
}

func (c *Clickable) Get() interface{} {
	return c
}


// Prompt is the struct for the server message that is sent to the client to tell the client what the prompt should be
type Prompt struct {
	ServerMessageBase
	Prompt string `json:"prompt"`
}

func (c *Prompt) Get() interface{} {
	return c
}

// HistoryMode is the struct for the server message that is sent to the client to tell the client to turn history saving on and off.
type HistoryMode struct {
	ServerMessageBase
	Val bool `json:"val"`
}

func (c *HistoryMode) Get() interface{} {
	return c
}

// Authenticated is the struct for the server message that is sent to the client to tell the client that is has been authenticated or not.
type Authenticated struct {
	ServerMessageBase
	Val bool `json:"val"`
}

func (c *Authenticated) Get() interface{} {
	return c
}

// Echo is the struct for the server message that is sent to the client to tell the client to turn echo on or off.
type Echo struct {
	ServerMessageBase
	Val bool `json:"val"`
}

func (c *Echo) Get() interface{} {
	return c
}

// Status is the struct for the server message that is sent to the client to tell the client to set the status bar to the text defined in the message.
type Status struct {
	ServerMessageBase
	Text string `json:"text"`
}

func (c *Status) Get() interface{} {
	return c
}

// AppendRawText will return a command packet that will append the text to the bottom of the output in the web terminal. This will format the text message in the color defined.
func AppendText(text, color string) ServerMessage {
	something := &TextMessage{}
	something.Type = "text"
	something.Text = html.EscapeString(text)
	something.Color = color
	return something
}

// AppendRawText will return a command packet that will append the raw text to the bottom of the output in the web terminal
func AppendRawText(text string) ServerMessage {
	something := &RawTextMessage{}
	something.Type = "rawtext"
	something.Text = text
	return something
}


// ClickableCommands will return a command packet that will append the raw text to the bottom of the output in the web terminal
func ClickableCommands( commands []string) ServerMessage {
	something := &Clickable{}
	something.Type = "clickable"
	something.Commands = commands
	return something
}


// SetPrompt will return a command packet that will set the prompt in the web terminal.
func SetPrompt(prompt string) ServerMessage {
	something := &Prompt{}
	something.Type = "prompt"
	something.Prompt = html.EscapeString(prompt)
	return something
}

// SetHistoryMode will return a command packet that will notify the web terminal that is should capture/not capture commands entered into the client's history/
func SetHistoryMode(val bool) ServerMessage {
	something := &HistoryMode{}
	something.Type = "history"
	something.Val = val
	return something
}

// SetAuthenticated will return a command packet that will notify the web terminal that the client has authenticated.
func SetAuthenticated(val bool) ServerMessage {
	something := &Authenticated{}
	something.Type = "authenticated"
	something.Val = val
	return something
}

// SetEchoOn will return a command packet that will set the echo state of text that is executed. Useful for disabling the display of the password entry in the auth process.
func SetEchoOn(val bool) ServerMessage {
	something := &Echo{}
	something.Type = "echo"
	something.Val = val
	return something
}

// SetStatus will return a command packet that will clear the current browser.
func Cls() ServerMessage {
	something := &ServerMessageBase{}
	something.Type = "cls"
	return something
}

// SetStatus will return a command packet that will set the status in the html window. This is yet to be implemented on the client.
func SetStatus(text string) ServerMessage {
	something := &Status{}
	something.Type = "status"
	something.Text = html.EscapeString(text)
	return something
}

// Eval will return a command packet that will be evaluated on the client browser. The text is the java script that will be evaluated.
func Eval(text string) ServerMessage {
	something := &Status{}
	something.Type = "eval"
	something.Text = text
	return something
}

// Link is used to format a string for a href tag
func Link(url, text string) string {
	return fmt.Sprintf("<a href=\"%s\" target=\"_blank\">%s</a>", url, text)
}

// Color is used to a color string and return a string containing the color for a span.
func Color(color, text string) string {
	if text == "" {
		return html.EscapeString(text)
	}

	return fmt.Sprintf("<span style=\"color:%s\">%s</span>", color, html.EscapeString(text))
}

// Image is used to take several predefined fields and return an html snipped to display an image.
func Image(width, height int, src, alt string) string {
	return fmt.Sprintf("<img width=\"%d\" height=\"%d\" src=\"%s\" alt=\"%s\"/>", width, height, src, alt)
}

// ClientMessage is the struct for the client message that is sent from the client to the server. nThe contents are passed to the
// HandleCommand             func(c Client, cmd string) defined in the Config.
type ClientMessage struct {
	Payload string `json:"payload"`
}
