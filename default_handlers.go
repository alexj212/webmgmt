package webmgmt

import (
	"bufio"
	"bytes"
	"fmt"
)

// ForegroundColor default foreground color definition
var ForegroundColor = "white"

// DefaultCommands main struct for all commands
var DefaultCommands = &Command{ExecLevel: ALL}

// HttpCommand command to view http session details
var HttpCommand = &Command{Use: "http", Exec: displayHttpInfo, Short: "Display http request information", ExecLevel: ALL}

// HistoryCommand command to view history of commands executed
var HistoryCommand = &Command{Use: "history", Exec: displayHistory, Short: "Show the history of commands executed", ExecLevel: ALL}

// UserCommand command to view user details
var UserCommand = &Command{Use: "user", Exec: displayUserInfo, Short: "Show user details about logged in user", ExecLevel: ALL}

// ClsCommand command to clear web screen
var ClsCommand = &Command{Use: "cls", Exec: func(client Client, args *CommandArgs) (err error) {
	client.Send(Cls())
	return
}, Short: "send cls event to terminal client", ExecLevel: ALL}

func init() {
	DefaultCommands.AddCommand(HttpCommand)
	DefaultCommands.AddCommand(HistoryCommand)
	DefaultCommands.AddCommand(UserCommand)
	DefaultCommands.AddCommand(ClsCommand)
	return
}

// HandleCommands handler function to execute commands
func HandleCommands(Commands *Command) (handler func(Client, string)) {

	handler = func(client Client, cmdLine string) {
		// loge.Info("handleMessage  - authenticated user message.Payload: [" + cmd+"]")

		var b bytes.Buffer
		writer := bufio.NewWriter(&b)

		parsed, err := NewCommandArgs(cmdLine, writer)

		if err != nil {
			client.Send(AppendText(fmt.Sprintf("Error parsing command: %v", err), "red"))
			return
		}
		Commands.Execute(client, parsed)
		writer.Flush()
		result := b.String()
		client.Send(AppendText(result, "white"))
	}
	return
}

func displayUserInfo(client Client, args *CommandArgs) (err error) {
	client.Send(AppendText(fmt.Sprintf("Username       : %v", client.Username()), "green"))
	client.Send(AppendText(fmt.Sprintf("IsAuthenticated: %v", client.IsAuthenticated()), "green"))
	return
}

func displayHttpInfo(client Client, args *CommandArgs) (err error) {
	client.Send(AppendText(fmt.Sprintf("GetIPAdress              : %v", GetIPAddress(client.HttpReq())), "green"))
	client.Send(AppendText(fmt.Sprintf("client.HttpReq.Host      : %v", client.HttpReq().Host), "green"))
	client.Send(AppendText(fmt.Sprintf("client.HttpReq.Method    : %v", client.HttpReq().Method), "green"))
	client.Send(AppendText(fmt.Sprintf("client.HttpReq.RemoteAddr: %v", client.HttpReq().RemoteAddr), "green"))
	client.Send(AppendText(fmt.Sprintf("client.HttpReq.RequestURI: %v", client.HttpReq().RequestURI), "green"))
	client.Send(AppendText(fmt.Sprintf("client.HttpReq.Referer() : %v", client.HttpReq().Referer()), "green"))
	for i, cookie := range client.HttpReq().Cookies() {
		client.Send(AppendText(fmt.Sprintf("client.HttpReq.Cookies[%-2d]: %-25v / %v", i, cookie.Name, cookie.Value), "green"))
	}

	for name, values := range client.HttpReq().Header {
		// Loop over all values for the name.
		for _, value := range values {
			client.Send(AppendText(fmt.Sprintf("client.HttpReq.Header[%-25v]:  %v", name, value), "green"))
		}
	}
	return
}

func displayHistory(client Client, args *CommandArgs) (err error) {
	if len(client.History()) > 0 {
		for i, cmd := range client.History() {
			client.Send(AppendText(fmt.Sprintf("History[%d]: %v", i, cmd), "green"))
		}
	} else {
		client.Send(AppendText("History is empty", "green"))
	}
	return
}
