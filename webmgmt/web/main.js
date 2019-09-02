

let path = '/ws';

let tcp = window.location.protocol === 'https:' ? 'wss://' : 'ws://';
let ws_url = tcp + window.location.host + path;

console.log("tcp: " + tcp);
console.log("host: " + window.location.host);
console.log("ws_url: " + ws_url);

let options = {
  wsEndpoint: ws_url,
};


let repl = wsrepl(options);






const homeText = Termpage.replace(`
...My name is Tautvilas, welcome to my termpage!...
...Type HELP for the list of available commands....\n\n`,
    {
        HELP: Termpage.color('orange'),
        Tautvilas: Termpage.link('http://www.tautvilas.lt'),
    });

function processInput(input = '') {


    input = input.toLowerCase().trim();


    const commands = ['home', 'help', 'image'];
    if (input === "help") {
        return {text: "Available commands are `home`, `help` and `image`", commands: commands};
    } else if (input === 'home') {
        return {text: homeText, commands: commands};
    } else if (input === 'image') {
        return {
            text: '<img width="200" height="200" src="https://i.imgur.com/RDsb26sb.jpg" alt="me"/>',
            commands: commands
        };
    } else {
        return {text: 'Command not found\n', commands: commands};
    }
}

let term1 = Termpage.init(
    document.getElementById('window1'),
    term1ProcessInput,
    {
        initialCommand: window.location.hash.substr(1) || 'home'
    }
);

let term2 = Termpage.init(
    document.getElementById('window2'),
    (input) => {
        if (input === 'home') {
            return 'This terminal demonstrates async commands that have 50% chance of failing';
        }
        let resolveP;
        const promise = new Promise((resolve, reject) => {
            setTimeout(() => {
                if (Math.random() > 0.5) {
                    resolve({
                        text: 'async request was successfull'
                    });
                } else {
                    reject()
                }
            }, 500);
        });
        return promise;
    },
    {
        initialCommand: 'home',
        prompt: Termpage.color('green', 'type_anything:') + ':',
        autoFocus: false
    }
);


console.log('term1', term1);
console.log('term2', term2);





function term1ProcessInput(input = '') {
    let val = {payload: input};
    if (repl != null) {
        console.log('repl wsSend', val);
        repl.sendCmd(val);
    } else {
      console.log('repl is null');
    }
}






let connectOptions = {
  wsOpenedCallback: () => {
    console.log('opened');
  },
  wsClosedCallback: () => {
    console.log('closed');
  },
  wsMsgEnteredCallback: (msg) => {
    console.log('msgEntered> ', msg);
    term1.updatePrompt(msg);
    term2.updatePrompt(msg);
  },
  OnMessageRecevied: (msg) => {
    console.log('wsMsgReceived> ', msg);
    term1.updatePrompt(msg.prompt);
    term2.updatePrompt(msg.prompt);
  },
};

repl.connect(connectOptions);
