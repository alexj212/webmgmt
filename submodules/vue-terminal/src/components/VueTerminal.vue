<template>
  <div class="termpage-window" ref="terminal">

    <div fluid fill-height class="termpage-header" ref="terminal_header"></div>

    <div class="termpage-body" ref="terminal_body">
      <div ref="terminal_output"></div>
      <p ref="terminal_input_block" class="termpage-block termpage-input-block">
        <span ref="terminal_prompt" class="termpage-prompt"></span>
        <input ref="terminal_input" type="text" class="termpage-input"/>
      </p>
    </div>

    <div class="termpage-status" ref="terminal_status"></div>
  </div>
</template>


<script>
  import * as webterm from '../common/webterm';

  export default {
    name: 'VueTerminal',
    data: function () {
      return {
        id: null,
        terminal: null
      };
    },

    props: {
      prompt: {
        type: String,
        default: "$"
      },
      appendHistory: {
        type: Boolean,
        default: true
      },
      echo: {
        type: Boolean,
        default: true
      },
      autoFocus: {
        type: Boolean,
        default: true
      },

      intro: {
        type: String,
        default: ""
      },
      initialCommand: {
        type: String,
        default: ""
      },
      status: {
        type: String,
        default: ""
      }
    },

    methods: {
      term1ProcessInput(input = "") {
        console.log("term1ProcessInput", input);
        this.$emit("command", input);
      },

      updatePrompt(prompt) {
        this.terminal.updatePrompt(prompt);
      },

      append(output) {
        this.terminal.append(output);
      },

      cls() {
        this.terminal.cls();
      },

      setStatus(status) {
        console.log("setStatus called " + status);
        this.terminal.setStatus(status);
      },
      setHeader(header) {
        console.log("setHeader called " + header);
        this.terminal.setHeader(header);
      },
      setStatusVisible(v) {
        console.log("setStatusVisible called " + v);
        this.terminal.setStatusVisible(v);
      },
      setHeaderVisible(v) {
        console.log("setHeaderVisible called " + v);
        this.terminal.setHeaderVisible(v);
      },


      setAppendHistory(val) {
        console.log("setAppendHistory called " + val);
        this.terminal.options.appendHistory = val;
      },

      setEcho(val) {
        console.log("setEcho called " + val);
        this.terminal.options.echo = val;
      },

    },

    mounted: function () {
      let options = {
        prompt: this.prompt,
        appendHistory: this.appendHistory,
        echo: this.echo,
        autoFocus: this.autoFocus
      };


      let dom = {
        $winElement: this.$refs.terminal,
        $headerElement: this.$refs.terminal_header,
        $bodyElement: this.$refs.terminal_body,
        $statusElement: this.$refs.terminal_status,
        $inputBlock: this.$refs.terminal_input_block,
        $input: this.$refs.terminal_input,
        $output: this.$refs.terminal_output,
        $prompt: this.$refs.terminal_prompt,
      };

      this.terminal = webterm.WebTerm.init(
        dom,
        this.term1ProcessInput,
        options
      );


      if (this.status !== "") {
        console.log("this.status set " + this.status);
        this.terminal.setStatus(this.status);
      }

      if (this.intro !== "") {
        console.log("this.intro set " + this.status);
        this.terminal.append(this.intro);
      }

      if (this.initialCommand !== "") {
        console.log("this.initialCommand set " + this.initialCommand);
        this.terminal.handleCommand(this.initialCommand);
      }
    }
  };
</script>

<style lang="css">

  .termpage-window {
    display: grid;

    grid-template-columns: 100%;
    grid-template-rows: 30px auto 30px;
    grid-template-areas: "termpage-header" "termpage-body" "termpage-status";
    border: 1px solid #818;
    height: calc(100vh - 200px);
  }

  .termpage-body {
    background-color: #000;
    grid-area: termpage-body;
    overflow-y: auto;
  }

  .termpage-header {
    grid-area: termpage-header;
    background-color: blue;
  }

  .termpage-status {
    grid-area: termpage-status;
    background-color: blue;
  }

  .termpage-window a {
    background-color: #888;
    text-decoration: none;
    cursor: pointer;
  }

  .termpage-window a:hover {
    background-color: #333;
  }

  .termpage-window * {
    font-family: "Courier New", Courier, monospace;
    font-size: 16px;
    color: #ddd;
  }

  .termpage-header * {
    font-family: "Courier New", Courier, monospace;
    font-size: 16px;
    color: #ddd;
  }

  .termpage-status * {
    font-family: "Courier New", Courier, monospace;
    font-size: 16px;
    color: #ddd;
  }

  .termpage-input {
    background-color: black;
  }

  .termpage-input {
    background-color: #222;
    color: #ddd;
    caret-color: white;
  }

  .termpage-input-block {
    display: flex;
  }

  .termpage-input {
    border-width: 0;
    outline: 0;
    flex: 1;
    padding: 0;
  }

  .termpage-block,
  .termpage-input {
    line-height: 20px;
  }

  .termpage-block {
    padding-left: 5px;
    padding-right: 5px;
  }

  .termpage-block {
    margin: 0;
    padding: 0;
  }

  pre.termpage-block {
    word-break: keep-all;
    white-space: pre-wrap !important;
  }

  .termpage-menu {
    background-color: #888;
  }

  .termpage-menu {
    display: flex;
    list-style: none;
    margin: 0;
    padding: 0;
  }

  .termpage-menu li:hover {
    background-color: #666;
    cursor: pointer;
  }
</style>
