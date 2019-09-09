package main

import (
    "bufio"
    "bytes"
    "fmt"
    "io"
    "log"
    "sort"
    "strings"
    "time"

    "github.com/alexj212/webmgmt"
    "github.com/gorilla/mux"
    "github.com/potakhov/loge"
)
type commandFunc func(client webmgmt.Client, args *webmgmt.CommandArgs, out io.Writer) error

var commandMap map[string]commandFunc
var commands []string


func Setup(router *mux.Router ) (mgmtApp *webmgmt.MgmtApp, err error) {
    config := &webmgmt.Config{StaticHtmlDir: "./web"}
    config.DefaultPrompt = "$"
    config.WebPath = "/admin/"

    config.UserAuthenticator = func(client webmgmt.Client, s string, s2 string) bool {
        return s == "alex" && s2 == "bambam"
    }

    commandMap = make(map[string]commandFunc)
    commandMap["ticker"] = toggleTicker
    commandMap["image"] = func(client webmgmt.Client, args *webmgmt.CommandArgs, out io.Writer) (err error) {
        client.Send(webmgmt.AppendRawText(webmgmt.Image(200, 200, "https://avatars1.githubusercontent.com/u/174203?s=200&v=4", "me"), nil))
        return
    }

    commandMap["raw"] = func(client webmgmt.Client, args *webmgmt.CommandArgs, out io.Writer) (err error) {
        client.Send(webmgmt.AppendRawText(webmgmt.Image(100, 100, "https://avatars1.githubusercontent.com/u/174203?s=200&v=4", "me"), commands))
        return
    }
    commandMap["commands"] = func(client webmgmt.Client, args *webmgmt.CommandArgs, out io.Writer) (err error) {
        client.Send(webmgmt.AppendRawText("Set Commands", commands))
        return
    }

    commandMap["link"] = func(client webmgmt.Client, args *webmgmt.CommandArgs, out io.Writer) (err error) {
        client.Send(webmgmt.AppendRawText(webmgmt.Link("http://www.slashdot.org", webmgmt.Color("orange", "slashdot")), nil))
        return
    }

    commandMap["prompt"] = func(client webmgmt.Client, args *webmgmt.CommandArgs, out io.Writer) (err error) {
        client.Send(webmgmt.SetPrompt(webmgmt.Color("red", client.Username()) + "@" + webmgmt.Color("green", "myserver") + ":&nbsp;"))
        return
    }
    commandMap["http"] = displayHttpInfo
    commandMap["history"] = displayHistory
    commandMap["user"] = displayUserInfo
    commandMap["lines"] = lines

    commandMap["cls"] = func(client webmgmt.Client, args *webmgmt.CommandArgs, out io.Writer) (err error) {
        client.Send(webmgmt.Cls())
        return
    }

    commandMap["help"] = func(client webmgmt.Client, args *webmgmt.CommandArgs, out io.Writer) (err error) {
        client.Send(webmgmt.AppendText(fmt.Sprintf("Available Commands"), "green"))
        client.Send(webmgmt.AppendText(fmt.Sprintf("------------------"), "green"))
        for i, k := range commands {
            client.Send(webmgmt.AppendText(fmt.Sprintf("[%d] %s", i, k), "yellow"))
        }
        return
    }

    // commandMap["echo"] = func(client webmgmt.Client, args *webmgmt.CommandArgs) {
    //     args.FlagSet.
    //
    //
    //     client.Send(webmgmt.AppendText(fmt.Sprintf("Available Commands"), "green"))
    // }

    commandMap["status"] = func(client webmgmt.Client, args *webmgmt.CommandArgs, out io.Writer) (err error) {
        client.Send(webmgmt.SetStatus("Hello World"))
        return
    }
    commandMap["hidestatus"] = func(client webmgmt.Client, args *webmgmt.CommandArgs, out io.Writer) (err error) {
        client.Send(webmgmt.SetStatus(""))
        return
    }
    commandMap["eval"] = func(client webmgmt.Client, args *webmgmt.CommandArgs, out io.Writer) (err error) {
        client.Send(webmgmt.Eval("alert ('alex');"))
        return
    }

    commandMap["table"] = func(client webmgmt.Client, args *webmgmt.CommandArgs, out io.Writer) (err error) {

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
        client.Send(webmgmt.AppendRawText(tableCode, nil))
        return
    }

    commandMap["canvas"] = func(client webmgmt.Client, args *webmgmt.CommandArgs, out io.Writer) (err error) {

        id := fmt.Sprintf("id_%v", time.Now().Unix())
        width := 300
        height := 300
        canvasCode := fmt.Sprintf("<canvas id=\"%s\" width=\"%d\" height=\"%d\"></canvas>", id, width, height)
        client.Send(webmgmt.AppendRawText(canvasCode, nil))

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
    }

    config.HandleCommand = func(client webmgmt.Client, cmdLine string) {
        // loge.Info("handleMessage  - authenticated user message.Payload: [" + cmd+"]")

        var b bytes.Buffer
        writer := bufio.NewWriter(&b)

        parsed, err := webmgmt.NewCommandArgs(cmdLine, writer)

        if err != nil {
            client.Send(webmgmt.AppendText(fmt.Sprintf("Error parsing command: %v", err), "red"))
            return
        } else {
            cmdFunc, ok := commandMap[parsed.CmdName]
            if !ok {
                client.Send(webmgmt.AppendText(fmt.Sprintf("echo: %v", parsed.CmdLine), "green"))
                return
            }

            err = cmdFunc(client, parsed, writer)
            writer.Flush()

            if err != nil {
                client.Send(webmgmt.AppendRawText(fmt.Sprintf("%s\n\n", err), nil))
            }

            output := b.String()
            if output != "" {
                client.Send(webmgmt.AppendRawText(output, nil))
            }

        }

    }

    commands = make([]string, 0, len(commandMap))
    for k := range commandMap {
        commands = append(commands, k)
    }

    sort.Slice(commands, func(i, j int) bool { return strings.ToLower(commands[i]) < strings.ToLower(commands[j]) })

    config.NotifyClientAuthenticated = func(client webmgmt.Client) {
        client.Send(webmgmt.SetPrompt("$ "))
        loge.Info("New user authenticated on system: %v", client.Username())
    }

    config.UnregisterUser = func(client webmgmt.Client) {
        loge.Info("user logged off system: %v", client.Username())
    }

    config.WelcomeUser = func(client webmgmt.Client) {
        client.Send(webmgmt.AppendText("Welcome to the machine", "red"))
        client.Send(webmgmt.SetPrompt("Enter Username: "))
        client.Send(webmgmt.SetAuthenticated(false))
        client.Send(webmgmt.SetHistoryMode(false))
        client.Send(webmgmt.SetEchoOn(true))

    }


    config.ClientInitializer = func(client webmgmt.Client) {
        client.Misc()["aa"] = 112
    }

    mgmtApp, err = webmgmt.NewMgmtApp("testapp", "1", config, router)
    return
}



func displayHistory(client webmgmt.Client, args *webmgmt.CommandArgs, out io.Writer) (err error) {
    if len(client.History()) > 0 {
        for i, cmd := range client.History() {
            client.Send(webmgmt.AppendText(fmt.Sprintf("History[%d]: %v", i, cmd), "green"))
        }
    } else {
        client.Send(webmgmt.AppendText("History is empty", "green"))
    }
    return
}

func displayUserInfo(client webmgmt.Client, args *webmgmt.CommandArgs, out io.Writer) (err error) {
    client.Send(webmgmt.AppendText(fmt.Sprintf("Username       : %v", client.Username()), "green"))
    client.Send(webmgmt.AppendText(fmt.Sprintf("IsAuthenticated: %v", client.IsAuthenticated()), "green"))
    return
}

func displayHttpInfo(client webmgmt.Client, args *webmgmt.CommandArgs, out io.Writer) (err error) {
    client.Send(webmgmt.AppendText(fmt.Sprintf("GetIPAdress              : %v", webmgmt.GetIPAddress(client.HttpReq())), "green"))
    client.Send(webmgmt.AppendText(fmt.Sprintf("client.HttpReq.Host      : %v", client.HttpReq().Host), "green"))
    client.Send(webmgmt.AppendText(fmt.Sprintf("client.HttpReq.Method    : %v", client.HttpReq().Method), "green"))
    client.Send(webmgmt.AppendText(fmt.Sprintf("client.HttpReq.RemoteAddr: %v", client.HttpReq().RemoteAddr), "green"))
    client.Send(webmgmt.AppendText(fmt.Sprintf("client.HttpReq.RequestURI: %v", client.HttpReq().RequestURI), "green"))
    client.Send(webmgmt.AppendText(fmt.Sprintf("client.HttpReq.Referer() : %v", client.HttpReq().Referer()), "green"))
    for i, cookie := range client.HttpReq().Cookies() {
        client.Send(webmgmt.AppendText(fmt.Sprintf("client.HttpReq.Cookies[%-2d]: %-25v / %v", i, cookie.Name, cookie.Value), "green"))
    }

    for name, values := range client.HttpReq().Header {
        // Loop over all values for the name.
        for _, value := range values {
            client.Send(webmgmt.AppendText(fmt.Sprintf("client.HttpReq.Header[%-25v]:  %v", name, value), "green"))
        }
    }
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

