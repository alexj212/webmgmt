## WebMgmt  

#### Details
this project will allow for a web admin service to be embedded into a service. This will allow clients to open a browser to the access port and login to the service. Commands can be developed to access the server.

#### Web
html terminal borrowed from  https://github.com/tautvilas/termpage



#### Building & Running
```bash

make binaries            - make bianries into bin dir
cd bin
./example

open browser to:  http://localhost:1099

You can login with username: alex password: bambam
user                     - display user info
http                     - display http headers/cookies
prompt                   - chanmge prompt
link                     - respond with a link that can be clicked on in terminal window
ticker                   - display a ticker that updates info to window periodically
image                    - display an image in the terminal window
raw                      - dsiplay an image along with a command tool bar of commands that can be clicked on
commands                 - display command tool bar of commands that can be clicked on
history                  - display history of commands executed

any unknown command will be echoed back
```



## Initialization
1. Create a Config struct and set the template path to ./web
2. Set the DefaultPrompt
3. Set the Webpath that will be used to access the terminal via a browser
```.go
   config := &webmgmt.Config{StaticHtmlDir: "./web"}
    config.DefaultPrompt = "$"
    config.WebPath = "/admin/"
```


## ClientInitialization
The Client initialization func is invoked when a client connects to the system. The handler func can access and modify the
client state. It has access the Misc() which is a Map available to save data for the client session.
```.go

    config.ClientInitializer = func(client webmgmt.Client) {
        client.Misc()["aa"] = 112
    }

```


## WelcomeUser
The WelcomeUser func is invoked when the client connects. The Server has the ability to send ServerMessages to the
client terminal. In the example below we send 
1. A welcome banner
2. Set the Prompt
3. The the authenticated state to the client
4. Toggle the history mode for text sent from client to server to off. 
5. Toggle the echo text state for the client to true.
```.go

    config.WelcomeUser = func(client webmgmt.Client) {
        client.Send(webmgmt.AppendText("Welcome to the machine", "red"))
        client.Send(webmgmt.SetPrompt("Enter Username: "))
        client.Send(webmgmt.SetAuthenticated(false))
        client.Send(webmgmt.SetHistoryMode(false))
        client.Send(webmgmt.SetEchoOn(true))

    }

```


## Authentication
1. Set the User Auth function, This function will have access to the Client interface, where you can access the IP, http Request etc.
The submitted username and password will also be passed to validate the session. Function returns the state of authentication
```.go

    config.UserAuthenticator = func(client webmgmt.Client, username string, password string) bool {
        return username == "alex" && password == "bambam"
    }
```

#Post Authentication
The NotifyClientAuthenticated func is invoked when a client is authenticated. This can be used for logging purposes.
```.go
    config.NotifyClientAuthenticated= func(client webmgmt.Client) {

        client.Send(webmgmt.SetPrompt("$ "))
        loge.Info("New user authenticated on system: %v", client.Username())
    }

```

#Post Authentication Failure
The NotifyClientAuthenticatedFailed func is invoked when a client fails authentication. It will be invoked after the client is disconnected. . This can be used for logging purposes.
```.go
    config.NotifyClientAuthenticatedFailed= func(client webmgmt.Client) {

        loge.Info("Client Failed Authentication: %v", client.Username())
    }
```



# Client command handling
Below is an example of command handler. A Map is created and functions are stored in the map as a value for the command (key).
The Command function, will have references to the Client Terminal, The CommandArgs which the the parsed text entered, and an io.Writer
Via the Client reference, ServerMessages can be sent to the client. 
1. link command - sends text back to the client to be displayed in the terminal. In this example, its raw html containing an a href with color styling.
2. prompt command - send a Prompt ServerMessage to the client, to make the terminal change its prompt
3. help command - returns a list of the commands defined in the map.
4. config.HandleCommand  func - is invoked when the client sends text to the server. The text is parsed, the command function is looked up in the map,
and executed. Response is sent back to the client. If the command is not found, error message is sent back to the client.
```.go

    commandMap = make(map[string]commandFunc)

    commandMap["link"] = func(client webmgmt.Client, args *webmgmt.CommandArgs, out io.Writer) (err error) {
        client.Send(webmgmt.AppendRawText(webmgmt.Link("http://www.slashdot.org", webmgmt.Color("orange", "slashdot")), nil))
        return
    }

    commandMap["prompt"] = func(client webmgmt.Client, args *webmgmt.CommandArgs, out io.Writer) (err error) {
        client.Send(webmgmt.SetPrompt(webmgmt.Color("red", client.Username()) + "@" + webmgmt.Color("green", "myserver") + ":&nbsp;"))
        return
    }
  
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
```


## Client Disconnect - the Unregister function is invoked when a client websocket is broken, such as when the client is shut down, or 
navigates away from the web terminal.
```.go
    config.UnregisterUser = func(client webmgmt.Client) {
        loge.Info("user logged off system: %v", client.Username())

```

## Creation of the WebMgmt Terminal.
The NewMgmtApp func will create a new web terminal with the Config supplied and attach it to the http.Router passed in. The Name and instanceId
are used to set the "X-Server-Name" and "X-Server-Id" headers on the http server. This can be useful to identifying servers behind a load balancer.
An error will be returned if initialization fails. 
  
```.go

    mgmtApp, err = webmgmt.NewMgmtApp("testapp", "1", config, router)
```


## Config
The following fields are used to initialize the WebMgmt command server. 
```.go

// Config struct  is used to configure a WebMgmt admin handler.
type Config struct {
	StaticHtmlDir                   string
	DefaultPrompt                   string
	WebPath                         string
	UserAuthenticator               func(client Client, username string, password string) bool
	HandleCommand                   func(c Client, cmd string)
	NotifyClientAuthenticated       func(client Client)
	notifyClientAuthenticatedFailed func(client Client)
	WelcomeUser                     func(client Client)
	UnregisterUser                  func(client Client)
	ClientInitializer               func(client Client)
}

```
