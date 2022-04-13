WebMgmt = {

  init: (customArgs) => {

    let defArgs = {
      Endpoint: 'ws://127.0.0.1/ws-repl',

      ClosedCallback: false,
      OpenedCallback: false,
      MessageSentCallback: false,
      MessageReceivedCallback: false,
    };

    let arg = Object.assign({}, defArgs, customArgs);

    arg.sendMessage = function (val) {
      if (arg.ws != null) {
        let pkt = {payload: val};
        console.log('sendMessage pkt', pkt)

        let payload = JSON.stringify(pkt);
        console.log('sendMessage payload', payload)
        arg.ws.send(payload);

      } else {
        console.log('sendMessage arg.ws is null', val)
      }

      if (arg.MessageSentCallback) {
        arg.MessageSentCallback(val)
      }
    };


    arg.connect = function connect() {

      arg.ws = new WebSocket(arg.Endpoint);
      arg.ws.addEventListener('open', function (event) {
        console.log('ws.open', event);

        if (arg.OpenedCallback) {
          arg.OpenedCallback();
        }
      });

      arg.ws.addEventListener('message', function (event) {
        console.log('addEventListener message>> ' + event.data);
        let message = JSON.parse(event.data);

        if (arg.MessageReceivedCallback) {
          arg.MessageReceivedCallback(message)
        }
      });

      arg.ws.addEventListener('close', function (event) {
        console.log('ws.close', event);

        if (arg.ClosedCallback) {
          arg.ClosedCallback();
        }exit

      });
    };


    arg.isConnected = function () {
      return arg.ws.readyState === arg.ws.OPEN;
    };

    return arg;
  }
}
