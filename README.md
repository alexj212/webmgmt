## WebMgmt  
This library provides an easy way to embed a command execution shell to an existing service application. The execution shell is access
via a web browser at a specific url path. Authentication can be implemented or bypassed to allow all to access. Configuration can
set UI features such as prompt can be configured. Command functions can be defined so that when the client sends a command with arguments 
the function is invoked with the parsed command line. 



## Command func 
This is example of using flags to set a field, if the command is invoked with the `-help` flag, then the help output is displayed.

```.env

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
```

Here the command is defined in the Map of Commands. Within the Command struct Help along with the ExecLevel can be set. 
```.go
	cmd = &webmgmt.Command{Exec: lines, ExecLevel:webmgmt.ALL, Help: "Displays N lines of text"}
	webmgmt.Commands["lines"] = cmd
``` 



## Details
this project will allow for a web admin service to be embedded into a service. This will allow clients to open a browser to the access port and login to the service. Commands can be developed to access the server.

## Web
html terminal assets borrowed from  https://github.com/tautvilas/termpage



## Building & Running
There is a sample web terminal that is embedded in a test application. The source resides in `./example`. To build the project
You can follow the steps below.
```bash

make binaries            - make bianries into bin dir
cd bin
./example

open browser to:  http://localhost:1099

You can login with username: alex password: bambam

Commands to try to enter

help                     - display the commands available.
user                     - display user info
http                     - display http headers/cookies
prompt                   - chanmge prompt
link                     - respond with a link that can be clicked on in terminal window
ticker                   - display a ticker that updates info to window periodically
image                    - display an image in the terminal window
raw                      - dsiplay an image along with a command tool bar of commands that can be clicked on
commands                 - display command tool bar of commands that can be clicked on
history                  - display history of commands executed


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

## Post Authentication
The NotifyClientAuthenticated func is invoked when a client is authenticated. This can be used for logging purposes.
```.go
    config.NotifyClientAuthenticated= func(client webmgmt.Client) {

        client.Send(webmgmt.SetPrompt("$ "))
        loge.Info("New user authenticated on system: %v", client.Username())
    }

```

## Post Authentication Failure
The NotifyClientAuthenticatedFailed func is invoked when a client fails authentication. It will be invoked after the client is disconnected. . This can be used for logging purposes.
```.go
    config.NotifyClientAuthenticatedFailed= func(client webmgmt.Client) {
        loge.Info("user auth failed on system: %v - %v", client.Username(), client.Ip())
    }
```



## Client command handling
Below is adding of commands that will be available to be executed. A webmgmt.Command struct is used to store a reference to the func
To be executed, The Help text, along with the ExecRights needed to execute the command. Several commands are added by default such as 
`help` or `cls`. The `help` command with look at the clients ExecRights to see if they have access to exec that command. 
The func that is defined for the command has access to the client, to send them ServerMessages.
```.go
	cmd = &webmgmt.Command{Exec: func(client webmgmt.Client, args *webmgmt.CommandArgs, out io.Writer) (err error) {
		client.Send(webmgmt.AppendRawText(webmgmt.Image(200, 200, "https://avatars1.githubusercontent.com/u/174203?s=200&v=4", "me"), nil))
		return
	}, Help: "Returns raw html to display image in terminal"}
	webmgmt.Commands["image"] = cmd

	cmd = &webmgmt.Command{Exec: func(client webmgmt.Client, args *webmgmt.CommandArgs, out io.Writer) (err error) {
		client.Send(webmgmt.AppendRawText(webmgmt.Link("http://www.slashdot.org", webmgmt.Color("orange", "slashdot")), nil))
		return
	}, Help: "Displays clickable link in terminal"}
	webmgmt.Commands["link"] = cmd

	cmd = &webmgmt.Command{Exec: func(client webmgmt.Client, args *webmgmt.CommandArgs, out io.Writer) (err error) {
		client.Send(webmgmt.SetPrompt(webmgmt.Color("red", client.Username()) + "@" + webmgmt.Color("green", "myserver") + ":&nbsp;"))
		return
	}, Help: "Updates the prompt to a multi colored prompt"}
	webmgmt.Commands["prompt"] = cmd
```


## Client Disconnect
The Unregister function is invoked when a client websocket is broken, such as when the client is shut down, or 
navigates away from the web terminal.
```.go
    config.UnregisterUser = func(client webmgmt.Client) {
        loge.Info("user logged off system: %v", client.Username())
    }
```

## Creation
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
	NotifyClientAuthenticatedFailed func(client Client)
	WelcomeUser                     func(client Client)
	UnregisterUser                  func(client Client)
	ClientInitializer               func(client Client)
}
```

    StaticHtmlDir       Used to define where assets are located. If the field is not set or not a directory, then the 
                        embedded assets are used.
                        
    DefaultPrompt       prompt that should be displayed in the terminal.
    WebPath             http routing path to access the web terminal
    UserAuthenticator   function to be invoked to authorize a user. The return value is the authenticated state
    HandleCommand       function to be invoked when a command is submitted to the server. By default a handler function
                        that validates the clients ExecLevel. If the client does not have access to a command then the
                        client will receive a response indicating they do not have rights to execute.  If a command is
                        not available then an error is sent. Otherwise the command is executed. Commands are case sensitive.
                        This can be set if the developer needed additional functionality such as logging etc.
    NotifyClientAuthenticated This is invoked when a client is authenticated. Can be used for logging this information etc.
    notifyClientAuthenticatedFailed This is invoked when a client has failed 3 password attempts for a username. The client 
                        is already disconnected. This can be used for logging etc.
     
                         


## Html Assets
The webmgmt uses several html and js resources that are delivered to the client. They are embedded into the webmgmt library with the 
use of embed that can mimic a filesystem, while the assets are encoded into a go file via the packr command. Users of the library 
can set a directory to be used instead of the embedded assets. There is a utility method  
`func webmgmt.SaveAssets(outputDir string) error` This will save all assets into a directory specified. Then a developer can customize 
the assets. In webmgmt.Config the field `StaticHtmlDir` if defined and it exists will be used to serve assets from. If that field is not set or 
does not exist the embedded assets will be used.


## Example App
In this snippet we use a flag passed to the app to write out assets so that they can be customized. 
```.go
    var saveTemplateDir string

	flag.StringVar(&saveTemplateDir, "save", "", "save assets to directory")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] \n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()


	if saveTemplateDir != "" {
		err = webmgmt.SaveAssets(saveTemplateDir)
		if err != nil {
			loge.Printf("Error writing assets: %v", err)
			os.Exit(-1)
		}
	}
```


## ServerMessages
Server sends json payloads via the websocket to the client terminal. The Client terminal will process the ServerMessages to
trigger begaviors within the client terminal. Such as display text, clear screen, render html, eval javascript code etc.

```.json
// TextMessaage is the json for the server message that is sent to the client to tell the client to display text in the terminal window.
{
	type="text",
	text="Hello World",
	color="red"
}


// RawTextMessage is the struct for the server message that is sent to the client to tell the client to display text as raw text in the terminal window.
{
	type="rawtext",
	text="<a href=....",
}

// Prompt is the struct for the server message that is sent to the client to tell the client what the prompt should be
{
	type="prompt",
	prompt=" > $ "
}


// HistoryMode is the struct for the server message that is sent to the client to tell the client to turn history saving on and off.
{
	type="history",
	val=true
}


// Authenticated is the struct for the server message that is sent to the client to tell the client that is has been authenticated or not.
{
	type="authenticated",
	val=true
}


// Echo is the struct for the server message that is sent to the client to tell the client to turn echo on or off.
{
	type="echo",
	val=true
}


// Status is the struct for the server message that is sent to the client to tell the client to set the status bar to the text defined in the message.
{
	type="status",
	text="Custom Status Text"
}


// cls is the struct for the server message that is sent to the client to tell the client to set the status bar to the text defined in the message.
{
	type="cls",
}

// Eval is the struct for the server message that is sent to the client to tell the client to call js.eval on the val set.
{
	type="eval",
	text="alert('hello world');"
}



```

The webmgmt has several utility functions that will create the various ServerMessages structs to send to the client.


## Acknowledgements
Original inspiration was taken from https://github.com/tautvilas/termpage I hacked upon the js code to add functionality needed
to create various ServerMessages to manipulate the web terminal from the server. The original assets are located under 
[./web]

## TODO
Implement a StatusBar within the web ui. The ServerMessage is defined just needs someone with js/css styling skills.
