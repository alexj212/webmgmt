package main

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/alexj212/webmgmt"
	"github.com/gorilla/mux"
	"github.com/potakhov/loge"
)

func Setup(router *mux.Router) (mgmtApp *webmgmt.MgmtApp, err error) {

	//## Initialization
	// 1. Create a Config struct and set the template path to ./web
	// 2. Set the DefaultPrompt
	// 3. Set the Webpath that will be used to access the terminal via a browser

	config := &webmgmt.Config{StaticHtmlDir: "./web"}
	config.DefaultPrompt = "$"
	config.WebPath = "/admin/"

	// ## ClientInitialization
	// The Client initialization func is invoked when a client connects to the system. The handler func can access and modify the
	// client state. It has access the Misc() which is a Map available to save data for the client session.

	config.ClientInitializer = func(client webmgmt.Client) {
		client.Misc()["aa"] = 112
	}

	//    ## WelcomeUser
	//    The WelcomeUser func is invoked when the client connects. The Server has the ability to send ServerMessages to the
	//    client terminal. In the example below we send
	//    1. A welcome banner
	//    2. Set the Prompt
	//    3. The the authenticated state to the client
	//    4. Toggle the history mode for text sent from client to server to off.
	//    5. Toggle the echo text state for the client to true.
	config.WelcomeUser = func(client webmgmt.Client) {
		client.Send(webmgmt.AppendText("Welcome to the machine", "red"))
		client.Send(webmgmt.SetPrompt("Enter Username: "))
		client.Send(webmgmt.SetAuthenticated(false))
		client.Send(webmgmt.SetHistoryMode(false))
		client.Send(webmgmt.SetEchoOn(true))

	}

	//    ## Authentication
	//    1. Set the User Auth function, This function will have access to the Client interface, where you can access the IP, http Request etc.
	//    The submitted username and password will also be passed to validate the session. Function returns the state of authentication
	config.UserAuthenticator = func(client webmgmt.Client, s string, s2 string) bool {
		return s == "alex" && s2 == "bambam"
	}

	// #Post Authentication
	// The NotifyClientAuthenticated func is invoked when a client is authenticated. This can be used for logging purposes.
	config.NotifyClientAuthenticated = func(client webmgmt.Client) {
		client.SetExecLevel(webmgmt.ADMIN)
		client.Send(webmgmt.SetPrompt("$ "))
		loge.Info("New user authenticated on system: %v", client.Username())
	}

	// #Post Authentication Failure
	// The NotifyClientAuthenticatedFailed func is invoked when a client fails authentication. It will be invoked after the client is disconnected. . This can be used for logging purposes.

	config.NotifyClientAuthenticatedFailed = func(client webmgmt.Client) {

		loge.Info("user auth failed on system: %v - %v", client.Username(), client.Ip())
	}

	cmd := &webmgmt.Command{Exec: toggleTicker, Help: "Start/Stop ticker that periodically sends updates to client"}
	webmgmt.Commands["ticker"] = cmd

	cmd = &webmgmt.Command{Exec: func(client webmgmt.Client, args *webmgmt.CommandArgs, out io.Writer) (err error) {
		client.Send(webmgmt.AppendRawText(webmgmt.Image(200, 200, "https://avatars1.githubusercontent.com/u/174203?s=200&v=4", "me")))
		return
	}, Help: "Returns raw html to display image in terminal"}
	webmgmt.Commands["image"] = cmd

	cmd = &webmgmt.Command{Exec: func(client webmgmt.Client, args *webmgmt.CommandArgs, out io.Writer) (err error) {
		client.Send(webmgmt.AppendRawText(webmgmt.Link("http://www.slashdot.org", webmgmt.Color("orange", "slashdot"))))
		return
	}, Help: "Displays clickable link in terminal"}
	webmgmt.Commands["link"] = cmd

	cmd = &webmgmt.Command{Exec: func(client webmgmt.Client, args *webmgmt.CommandArgs, out io.Writer) (err error) {
		client.Send(webmgmt.SetPrompt(webmgmt.Color("red", client.Username()) + "@" + webmgmt.Color("green", "myserver") + ":&nbsp;"))
		return
	}, Help: "Updates the prompt to a multi colored prompt"}
	webmgmt.Commands["prompt"] = cmd

	cmd = &webmgmt.Command{Exec: lines, ExecLevel:webmgmt.ALL, Help: "Displays N lines of text"}
	webmgmt.Commands["lines"] = cmd

	cmd = &webmgmt.Command{Exec: func(client webmgmt.Client, args *webmgmt.CommandArgs, out io.Writer) (err error) {
		client.Send(webmgmt.SetStatus("Hello World"))
		return
	},
		Help: "Sets the status to Hello World"}
	webmgmt.Commands["status"] = cmd

	cmd = &webmgmt.Command{Exec: func(client webmgmt.Client, args *webmgmt.CommandArgs, out io.Writer) (err error) {
		client.Send(webmgmt.SetStatus(""))
		return
	},
		Help: "Clears the status"}
	webmgmt.Commands["hidestatus"] = cmd

	cmd = &webmgmt.Command{Exec: func(client webmgmt.Client, args *webmgmt.CommandArgs, out io.Writer) (err error) {
		client.Send(webmgmt.Eval("alert ('alex');"))
		return
	}, Help: "Evals sends js to the client to be evaluated"}
	webmgmt.Commands["eval"] = cmd



	cmd = &webmgmt.Command{Exec: func(client webmgmt.Client, args *webmgmt.CommandArgs, out io.Writer) (err error) {
		commands :=  []string{"help", "cls", "lines", "link", "image",}
		client.Send(webmgmt.ClickableCommands(commands))
		client.Send(webmgmt.Eval("alert ('alex');"))
		return
	}, Help: "Sends some clickable commands to be displayed"}
	webmgmt.Commands["clickable"] = cmd

	cmd = &webmgmt.Command{Exec: func(client webmgmt.Client, args *webmgmt.CommandArgs, out io.Writer) (err error) {

		var tableCode = `
<table style="width:100%">
  <tr>
    <th>Name</th>
    <th colspan="2">Telephone</th>
  </tr>
  <tr>
    <td>Bill Gates</td>
    <td>55577854</td>
    <td>55577855</td>
  </tr>
</table>
`
		client.Send(webmgmt.AppendRawText(tableCode))
		return
	},
		Help: "Example table returned to the client"}
	webmgmt.Commands["table"] = cmd

	cmd = &webmgmt.Command{Exec: func(client webmgmt.Client, args *webmgmt.CommandArgs, out io.Writer) (err error) {

		id := fmt.Sprintf("id_%v", time.Now().Unix())
		width := 300
		height := 300
		canvasCode := fmt.Sprintf("<canvas id=\"%s\" width=\"%d\" height=\"%d\"></canvas>", id, width, height)
		client.Send(webmgmt.AppendRawText(canvasCode))

		jsCodeTemplate := `
var canvas = document.getElementById('%s');
console.log('got canvas');
var ctx = canvas.getContext('2d');
ctx.fillStyle = "#FF0000";
ctx.fillRect(0, 0, 300, 300);

console.log('got canvas ctx');
ctx.beginPath();
ctx.arc(100, 100, 50, 1.5 * Math.PI, 0.5 * Math.PI, false);
ctx.lineWidth = 10;
ctx.stroke();
ctx.closePath();
console.log('got canvas closePath');
ctx.stroke();
console.log('got canvas stroke');
console.log('done');
`
		jsCode := fmt.Sprintf(jsCodeTemplate, id)
		client.Send(webmgmt.Eval(jsCode))
		return
	},
		Help: "Returns a canvas with js to update it"}
	webmgmt.Commands["canvas"] = cmd

	config.HandleCommand = webmgmt.HandleCommands()

	config.UnregisterUser = func(client webmgmt.Client) {
		loge.Info("user logged off system: %v", client.Username())
	}

	mgmtApp, err = webmgmt.NewMgmtApp("testapp", "1", config, router)
	return
}

func lines(client webmgmt.Client, args *webmgmt.CommandArgs, out io.Writer) (err error) {
	cnt := args.FlagSet.Int("cnt", 5, "number of lines to print")
	err = args.Parse()
	if err != nil {
		return
	}

	client.Send(webmgmt.AppendText(fmt.Sprintf("lines invoke"), "green"))
	log.Printf("lines invoked")
	for i := 0; i < *cnt; i++ {
		client.Send(webmgmt.AppendText(fmt.Sprintf("line[%d]", i), "green"))
	}
	return
}

func toggleTicker(client webmgmt.Client, args *webmgmt.CommandArgs, out io.Writer) (err error) {

	tickerDone, ok := client.Misc()["ticker_done"]
	if ok {
		delete(client.Misc(), "ticker_done")

		var done chan bool
		done = tickerDone.(chan bool)
		done <- true
	} else {
		done := make(chan bool)
		client.Misc()["ticker_done"] = done

		go func() {
			ticker := time.NewTicker(5 * time.Second)

			for {
				if !client.IsConnected() {
					break
				}

				select {

				case <-done:
					loge.Info("Ticker Done")
					return

				case t := <-ticker.C:
					loge.Info("Ticker ticked")
					if client.IsAuthenticated() {
						msg := fmt.Sprintf("Tick at: %v", t)
						client.Send(webmgmt.AppendText(msg, "blue"))
					}
				}
			}
			loge.Info("Ticker Stopped")
			ticker.Stop()
		}()

	}
	return
}
