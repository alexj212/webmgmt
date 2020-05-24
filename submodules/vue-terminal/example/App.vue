<template>
  <VueTerminal ref="term"
               prompt="$"
               status="my custom status"
               intro="my custom intro1"
               @command="onCliCommand"></VueTerminal>

</template>

<script>
  import VueTerminal from '../src/components/VueTerminal';
  import * as webterm from '../src/common/webterm';
  import * as webmgmt from './webmgmt';


  export default {
    layout: 'dashboard',
    components: {
      VueTerminal
    },
    data: () => ({
      intro: 'Hello World',
      webmgt: false,
    }),
    computed: {},
    mounted() {
      let path = 'ws';

      let tcp = window.location.protocol === 'https:' ? 'wss://' : 'ws://';
      let ws_url = tcp + window.location.host + window.location.pathname + path;

      ws_url = "ws://localhost:1099/admin/ws";

      let options = {
        Endpoint: ws_url,
        OpenedCallback: () => {
          terminal.append(webterm.WebTerm.color('red', "Connection Opened\n\n"));
        },

        ClosedCallback: () => {
          terminal.append(webterm.WebTerm.color('red', "Connection Closed\n\n"));
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
              this.$refs.term.append(webterm.WebTerm.color(color, msg.text));
              break;
            case "rawtext":
              this.$refs.term.append(msg);
              break;
            case "clickable":
              this.$refs.term.append(msg);
              break;
            case "history":
              this.$refs.term.setAppendHistory(msg.val);
              break;
            case "echo":
              this.$refs.term.setEcho(msg.val);
              break;
            case "authenticated":
              this.$refs.term.setAppendHistory(msg.val);
              break;
            case "prompt":
              this.$refs.term.updatePrompt(msg.prompt);
              break;
            case "cls":
              this.$refs.term.cls();
              break;
            case "status":
              console.log('status: ' + msg.text);
              this.$refs.term.setStatus(msg.text);
              break;
            case "eval":
              console.log('eval: ' + msg.text);
              eval(msg.text);
              break;

          }
        },
      };


      this.webmgt = webmgmt.WebMgmt.init(options);

      let terminal = this.$refs.term;

      this.webmgt.connect();
      console.log('this webmgt', this.webmgt);

    },
    methods: {

      onCliCommand(data) {
        if (!this.webmgt.isConnected()) {
          console.log('onCliCommand not connected reconnecting data: ' + data);
          this.webmgt.connect();
          return;
        }

        console.log('onCliCommand ' + data);

        if (!data || data === '') {
          return;
        } else if (data === ':cls') {
          this.$refs.term.cls();

        } else if (data.startsWith(':status ')) {
          let status = data.substring(8);
          console.log('status: ' + status);

          if (status.startsWith("show")) {
            this.$refs.term.setStatusVisible(true);
          } else if (status.startsWith("hide")) {
            this.$refs.term.setStatusVisible(false);
          } else {
            this.$refs.term.setStatus(status);
          }

        } else if (data.startsWith(':header ')) {

          let header = data.substring(8);
          console.log('header: ' + header);

          if (header.startsWith("show")) {
            this.$refs.term.setHeaderVisible(true);
          } else if (header.startsWith("hide")) {
            this.$refs.term.setHeaderVisible(false);
          } else {
            this.$refs.term.setHeader(header);
          }

        } else if (data.startsWith(':prompt ')) {


          let prompt = data.substring(8);
          console.log('prompt: ' + prompt);
          this.$refs.term.updatePrompt(prompt);

        } else {
          console.log('onCliCommand ' + data);
          this.webmgt.sendMessage(data);
          console.log('onCliCommand done  this.webmgt ', this.webmgt);
        }
      }

    },
  };
</script>


<style>

</style>
