function wsrepl(customArgs) {

    let defArgs = {
        wsEndpoint: 'ws://127.0.0.1/ws-repl',

        wsClosedCallback: wsClosed,
        wsOpenedCallback: wsOpened,
        wsMsgEnteredCallback: wsMsgEntered,
        OnMessageReceived: wsMsgReceived,
    };

    let arg = Object.assign({}, defArgs, customArgs);


    function replyToPage(str) {
    }

    function cls() {

        replyToPage('');
    }


    arg.sendCmd = function (val) {

        // console.log("sendCmd[" + val + ']');
        if (val === ":cls") {
            // console.log("clearing screen");
            cls();
            setStatus('cleared screen');
            return;
        }


        arg.wsSend(JSON.stringify(val));

        if (arg.wsMsgEnteredCallback) {
            arg.wsMsgEnteredCallback(val)
        }
    };

    arg.wsSend = function wsSend(val) {
        if (arg.ws != null) {
            //console.log('wsSend actual', val)
            arg.ws.send(val);
        } else {
            //console.log('wsSend arg.ws is null', val)
        }
    };


    function setStatus(val) {
    }

    function setTextEntryDisabled(val) {

    }

    function wsClosed() {
    }

    function wsOpened() {
    }

    function wsMsgEntered(msg) {
    }

    function wsMsgReceived(msg) {
    }


    arg.connect = function connect(options) {
        Object.assign(arg, options);

        arg.ws = new WebSocket(arg.wsEndpoint);
        arg.ws.addEventListener('open', function (event) {
            console.log('ws.open', event);
            setStatus('connected to ' + arg.wsEndpoint);
            setTextEntryDisabled(false);
            if (arg.wsOpenedCallback) {
                arg.wsOpenedCallback();
            }
        });
        arg.ws.addEventListener('message', function (event) {

            // console.log('addEventListener message>> ' + event.data);

            let message = JSON.parse(event.data);

            if (arg.OnMessageRecevied) {
                arg.OnMessageRecevied(message)
            }
            replyToPage(message.payload, true);
        });

        arg.ws.addEventListener('close', function (event) {
            console.log('ws.close', event);
            setStatus('disconnected from ' + arg.wsEndpoint);
            setTextEntryDisabled(true);

            if (arg.wsClosedCallback) {
                arg.wsClosedCallback();
            }
        });
    };


    arg.isConnected = function () {
        return arg.ws.readyState === arg.ws.OPEN;
    };

    arg.setStatus = setStatus;
    arg.setTextEntryDisabled = setTextEntryDisabled;
    arg.cls = cls;
    arg.replyToPage = replyToPage;

    return arg;
}
