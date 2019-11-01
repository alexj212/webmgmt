let path = 'ws';

let tcp = window.location.protocol === 'https:' ? 'wss://' : 'ws://';
let ws_url = tcp + window.location.host + window.location.pathname + path;

let options = {
    wsEndpoint: ws_url,
};


let repl = wsrepl(options);

let term1 = Termpage.init(
    document.getElementById('window1'),
    term1ProcessInput,
    {}
);


function term1ProcessInput(input = '') {
    let val = {payload: input};
    if (repl.isConnected()) {
        repl.sendCmd(val);
    } else {
        repl.connect();
    }
}


let connectOptions = {
    wsOpenedCallback: () => {
        term1.append(Termpage.color('red', "Connection Opened\n\n"));
    },

    wsClosedCallback: () => {
        term1.append(Termpage.color('red', "Connection Closed\n\n"));
        term1.updatePrompt('');
    },

    wsMsgEnteredCallback: (msg) => {
        //console.log('msgEntered> ', msg);
    },

    OnMessageRecevied: (msg) => {
        //console.log('wsMsgReceived> ', msg);
        switch (msg.type) {
            case "text":
                let color = msg.color;
                if (!color) {
                    color = "white";
                }
                term1.append(Termpage.color(color, msg.text));
                break;
            case "rawtext":
                term1.append(msg);
                break;
            case "clickable":
                term1.append(msg);
                break;
            case "history":
                term1.options.appendHistory = msg.val;
                break;
            case "echo":
                term1.options.echo = msg.val;
                break;
            case "authenticated":
                term1.options.appendHistory = msg.val;
                break;
            case "prompt":
                term1.updatePrompt(msg.prompt);
                break;
            case "cls":
                term1.cls();
                break;
            case "status":
                console.log('status: ' + msg.text);
                term1.setStatus(msg.text);
                break;
            case "eval":
                console.log('eval: ' + msg.text);
                eval(msg.text);
                break;

        }

    },
};

repl.connect(connectOptions);
