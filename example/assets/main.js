let terminalOptions = {
    prompt: "$",
    appendHistory: true,
    echo: true,
    autoFocus: true,
};


let dom = {
    $winElement: document.getElementById("terminal"),
    $headerElement: document.getElementById("terminal_header"),
    $bodyElement: document.getElementById("terminal_body"),
    $statusElement: document.getElementById("terminal_status"),
    $inputBlock: document.getElementById("terminal_input_block"),
    $input: document.getElementById("terminal_input"),
    $output: document.getElementById("terminal_output"),
    $prompt: document.getElementById("terminal_prompt"),
};

let terminal = WebTerm.init(
    dom,
    onCliCommand,
    terminalOptions
);





let path = 'ws';
let tcp = window.location.protocol === 'https:' ? 'wss://' : 'ws://';
let ws_url = tcp + window.location.host + window.location.pathname + path;
ws_url = "ws://localhost:1099/admin/ws";


let webmgmtOptions = {
    Endpoint: ws_url,
    OpenedCallback: () => {
        terminal.append(WebTerm.color('red', "Connection Opened\n\n"));
    },

    ClosedCallback: () => {
        terminal.append(WebTerm.color('red', "Connection Closed\n\n"));
        terminal.updatePrompt('');
    },

    MessageSentCallback: (msg) => {
        console.log('MessageSentCallback> ', msg);
    },

    MessageReceivedCallback: (msg) => {
        console.log('MessageReceivedCallback> ', msg);

        switch (msg.type) {
            case "text":
                let color = msg.color;
                if (!color) {
                    color = "white";
                }
                terminal.append(WebTerm.color(color, msg.text));
                break;
            case "rawtext":
                terminal.append(msg);
                break;
            case "clickable":
                terminal.append(msg);
                break;
            case "history":
                terminal.setAppendHistory(msg.val);
                break;
            case "echo":
                terminal.setEcho(msg.val);
                break;
            case "authenticated":
                terminal.setAppendHistory(msg.val);
                break;
            case "prompt":
                terminal.updatePrompt(msg.prompt);
                break;
            case "cls":
                terminal.cls();
                break;
            case "status":
                console.log('status: ' + msg.text);
                terminal.setStatus(msg.text);
                break;
            case "eval":
                console.log('eval: ' + msg.text);
                eval(msg.text);
                break;

        }
    },
};





function onCliCommand(data)
{
    if (!webmgt.isConnected()) {
        console.log('onCliCommand not connected reconnecting data: ' + data);
        webmgt.connect();
        return;
    }

    console.log('onCliCommand ' + data);

    if (!data || data === '') {
        return;
    } else if (data === ':cls') {
        terminal.cls();

    } else if (data.startsWith(':status ')) {
        let status = data.substring(8);
        console.log('status: ' + status);

        if (status.startsWith("show")) {
            terminal.setStatusVisible(true);
        } else if (status.startsWith("hide")) {
            terminal.setStatusVisible(false);
        } else {
            terminal.setStatus(status);
        }

    } else if (data.startsWith(':header ')) {

        let header = data.substring(8);
        console.log('header: ' + header);

        if (header.startsWith("show")) {
            terminal.setHeaderVisible(true);
        } else if (header.startsWith("hide")) {
            terminal.setHeaderVisible(false);
        } else {
            terminal.setHeader(header);
        }

    } else if (data.startsWith(':prompt ')) {


        let prompt = data.substring(8);
        console.log('prompt: ' + prompt);
        terminal.updatePrompt(prompt);

    } else {
        console.log('onCliCommand ' + data);
        webmgt.sendMessage(data);
        console.log('onCliCommand done  this.webmgt ', webmgt);
    }
}


let webmgt = WebMgmt.init(webmgmtOptions);
webmgt.connect();
console.log('this webmgt', webmgt);


