package webmgmt

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
)

type CommandFunc func(client Client, args *CommandArgs, out io.Writer) error

type ExecLevel int32

const ALL = ExecLevel(0)
const USER = ExecLevel(1)
const ADMIN = ExecLevel(2)

type Command struct {
	Exec      CommandFunc
	Help      string
	ExecLevel ExecLevel
}

var Commands map[string]*Command
var ForegroundColor string = "white"

func init() {
	Commands = make(map[string]*Command)
	var cmd *Command

	cmd = &Command{Exec: displayHttpInfo, Help: "Display http request information", ExecLevel: ALL}
	Commands["http"] = cmd

	cmd = &Command{Exec: displayHistory, Help: "Show the history of commands executed", ExecLevel: ALL}
	Commands["history"] = cmd

	cmd = &Command{Exec: displayUserInfo, Help: "Show user details about logged in user", ExecLevel: ALL}
	Commands["user"] = cmd

	cmd = &Command{Exec: func(client Client, args *CommandArgs, out io.Writer) (err error) {
		client.Send(Cls())
		return
	}, Help: "send cls event to terminal client", ExecLevel: ALL}
	Commands["cls"] = cmd

	cmd = &Command{Exec: func(client Client, args *CommandArgs, out io.Writer) (err error) {

		commands := make([]string, 0, len(Commands))
		for k := range Commands {
			commands = append(commands, k)
		}

		sort.Slice(commands, func(i, j int) bool { return strings.ToLower(commands[i]) < strings.ToLower(commands[j]) })

		client.Send(AppendText(fmt.Sprintf("Available Commands"), "green"))
		client.Send(AppendText(fmt.Sprintf("------------------"), "green"))
		for i, k := range commands {
			cmd := Commands[k]
			if client.ExecLevel() >= cmd.ExecLevel {
				client.Send(AppendText(fmt.Sprintf("[%d] %s  - %s", i, k, cmd.Help), "yellow"))
			}
		}
		return
	},
		Help: "Display Help", ExecLevel: ALL}
	Commands["help"] = cmd

	return
}

func HandleCommands() (handler func(Client, string)) {

	handler = func(client Client, cmdLine string) {
		// loge.Info("handleMessage  - authenticated user message.Payload: [" + cmd+"]")

		var b bytes.Buffer
		writer := bufio.NewWriter(&b)

		parsed, err := NewCommandArgs(cmdLine, writer)

		if err != nil {
			client.Send(AppendText(fmt.Sprintf("Error parsing command: %v", err), "red"))
			return
		} else {
			cmd, ok := Commands[parsed.CmdName]
			if !ok {
				client.Send(AppendText(fmt.Sprintf("Command not found: %v", parsed.CmdLine), "red"))
				return
			}

			if client.ExecLevel() >= cmd.ExecLevel {

				err = cmd.Exec(client, parsed, writer)
				_ = writer.Flush()

				if err != nil {
					client.Send(AppendText(fmt.Sprintf("%s\n\n", err), ForegroundColor))
				}

				output := b.String()
				if output != "" {
					client.Send(AppendText(output, ForegroundColor))
				}
			} else {
				client.Send(AppendText(fmt.Sprintf("You do not have permission to execute: %v", parsed.CmdLine), "red"))
				return
			}
		}
	}
	return
}

func displayUserInfo(client Client, args *CommandArgs, out io.Writer) (err error) {
	client.Send(AppendText(fmt.Sprintf("Username       : %v", client.Username()), "green"))
	client.Send(AppendText(fmt.Sprintf("IsAuthenticated: %v", client.IsAuthenticated()), "green"))
	return
}

func displayHttpInfo(client Client, args *CommandArgs, out io.Writer) (err error) {
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

func displayHistory(client Client, args *CommandArgs, out io.Writer) (err error) {
	if len(client.History()) > 0 {
		for i, cmd := range client.History() {
			client.Send(AppendText(fmt.Sprintf("History[%d]: %v", i, cmd), "green"))
		}
	} else {
		client.Send(AppendText("History is empty", "green"))
	}
	return
}
