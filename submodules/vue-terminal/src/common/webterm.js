/*
Original inspiration taken from https://github.com/tautvilas/termpage

Added ability to:
enable/disable adding to history
enable/disable display input when submitted.
set prompt
return reference to the terminal in the init function.

 */

export const WebTerm = {

  defaultOptions: {
    prompt: '$',
    autoFocus: true,
    appendHistory: true,
    echo: true,
  },

  _appendInput: (input, options, dom) => {
    //console.log('_appendInput.1 input: ', input);
    if (dom.$winElement.lastChild && dom.$winElement.lastChild.tagName === 'UL') {
      dom.$winElement.lastChild.remove();
    }
    //console.log('_appendInput.2');
    const prmpt = options.prompt + '&nbsp';
    //console.log('_appendInput.3');
    const pre = document.createElement("pre");
    //console.log('_appendInput.4');

    if (options.echo) {
      const encodedInput = input.replace(/[\u00A0-\u9999<>\&]/gim, function (i) {
        //console.log('_appendInput.5');
        return '&#' + i.charCodeAt(0) + ';';
      });
      //console.log('_appendInput.6');

      pre.innerHTML = prmpt + encodedInput + '\n';
    } else {
      pre.innerHTML = prmpt + '\n';
    }
    //console.log('_appendInput.7');
    pre.className = 'termpage-block';
    //console.log('_appendInput.8');
    dom.$output.appendChild(pre);
    //console.log('_appendInput.9');
  },

  _appendOutput: (output, options, dom) => {
    //console.log('_appendOutput.1 output: ', output);
    let outputText = "undefined";
    let commands = [];
    if (typeof (output) === "string") {
      outputText = output;
      //console.log('_appendOutput.2 '+outputText);
    } else if (typeof (output) === "object") {
      outputText = output.text;
      commands = output.commands || [];
      //console.log('_appendOutput.3 outputText: ' + outputText);
      //console.log('_appendOutput.3 commands: ' + commands);
    } else {
      //console.log('_appendOutput.3 not appending'+outputText);
      return;
    }
    //console.log('_appendOutput.4 ' + outputText);
    const pre = document.createElement("pre");
    pre.innerHTML = outputText;
    //console.log('_appendOutput.5');
    pre.className = 'termpage-block';
    //console.log('_appendOutput.6');
    dom.$output.appendChild(pre);
    //console.log('_appendOutput.7');
    if (commands.length) {
      //console.log('_appendOutput.8');
      const $commands = document.createElement('ul');
      $commands.className = 'termpage-menu termpage-block';
      //console.log('_appendOutput.9');
      commands.forEach(command => {
        //console.log('_appendOutput.10');
        const $command = document.createElement('li');
        //console.log('_appendOutput.11');
        $command.innerHTML = command + '&nbsp;';
        $commands.appendChild($command);
        $command.addEventListener('click', () => {
          //console.log('_appendOutput.12');
          WebTerm._appendInput($command.innerText, options, dom);
          const out = options.processInput($command.innerText);
          WebTerm._processInput(out, options, dom);
          //console.log('_appendOutput.13');
        });
        //console.log('_appendOutput.14');
      });
      //console.log('_appendOutput.15');
      dom.$winElement.appendChild($commands);
      //console.log('_appendOutput.16');
    }
    //console.log('_appendOutput.17');
    dom.$winElement.scrollTo(0, dom.$winElement.scrollHeight);
    dom.$bodyElement.scrollTo(0, dom.$bodyElement.scrollHeight);
    //console.log('_appendOutput.18');
    dom.$input.focus();
    //console.log('_appendOutput.19');
  },

  _processInput: (output, options, dom) => {
    //console.log('_processInput.1  output: ', output);
    if (output && output.then) {
      //console.log('_processInput.2 output: ' + output);
      const pre = document.createElement("pre");
      //console.log('_processInput.3');
      pre.innerHTML = '.';
      pre.className = 'termpage-loader termpage-block';
      dom.$output.appendChild(pre);
      dom.$inputBlock.setAttribute('style', 'display:none');
      dom.$winElement.scrollTo(0, dom.$winElement.scrollHeight);
      dom.$bodyElement.scrollTo(0, dom.$bodyElement.scrollHeight);

      //console.log('_processInput.4');
      output.then((out) => {
        //console.log('_processInput.5');
        pre.remove();
        dom.$inputBlock.setAttribute('style', 'display:flex');
        WebTerm._appendOutput(out, options, dom);
        //console.log('_processInput.6');
      });

      output.catch(() => {
        //console.log('_processInput.7');
        pre.remove();
        dom.$inputBlock.setAttribute('style', 'display:flex');
        WebTerm._appendOutput(WebTerm.color('red', 'command resolution failed'), options, dom);
        //console.log('_processInput.8');
      });
    } else {
      //console.log('_processInput.9  output: ', output);
      WebTerm._appendOutput(output, options, dom);
      //console.log('_processInput.10');
    }
  },

  /*

        let dom = {
          $winElement: this.$refs.terminal,
          $headerElement: this.$refs.terminal_header,
          $bodyElement:this.$refs.terminal_body,
          $statusElement: this.$refs.terminal_status,
          $inputBlock: this.$refs.terminal_input_block,
          $input: this.$refs.terminal_input,
          $output: this.$refs.terminal_outpu,
        };

   */

  init: (dom, processInput, options = {}) => {

    //console.log('termpage init');
    let terminal = {};
    terminal.history = [];
    terminal.historyIndex = 0;
    terminal.options = Object.assign({}, WebTerm.defaultOptions, options);

    terminal.dom = dom;
    terminal.$header = dom.$headerElement;
    terminal.$status = dom.$statusElement;
    terminal.$bodyElement = dom.$bodyElement;
    terminal.$winElement = dom.$winElement;
    terminal.$inputBlock = dom.$inputBlock;
    terminal.$input = dom.$input;
    terminal.$output = dom.$output;
    terminal.$prompt = dom.$prompt;

    terminal.prompt = terminal.options.prompt || WebTerm.defaultOptions.prompt;
    terminal.$prompt.innerHTML = terminal.prompt + "&nbsp;";


    console.log('terminal.dom', terminal.dom);

    terminal.options.processInput = (inp) => {
      if (terminal.options.appendHistory) {
        terminal.historyIndex = 0;
        terminal.history.push(inp);
      }
      return processInput(inp);
    };


    terminal.$input.addEventListener('keydown', function (e) {
      let key = e.which || e.keyCode;
      if (key === 13) { // 13 is enter
        const input = e.srcElement.value;
        WebTerm._appendInput(input, terminal.options, terminal.dom);
        terminal.options.processInput(input);
        terminal.$input.value = '';
      }
      if (key === 38) { // up
        const val = terminal.history[terminal.history.length - terminal.historyIndex - 1];
        if (val) {
          terminal.historyIndex++;
          terminal.dom.$input.value = val;
          terminal.setInputCaretToEnd();
          console.log('up - setInputCaretToEnd')
        }
      } else if (key === 40) { // down
        const val = terminal.history[terminal.history.length - terminal.historyIndex + 1];
        if (val) {
          terminal.historyIndex--;
          terminal.dom.$input.value = val;
          terminal.setInputCaretToEnd();
          console.log('down - setInputCaretToEnd')
        }

      } else if (e.keyCode == 27) {
        terminal.$input.value = '';
        terminal.setInputCaretToHome();

      }
    });

    if (terminal.options.autoFocus) {
      terminal.$input.focus();
    }
    terminal.$winElement.addEventListener("click", function () {
      const sel = getSelection().toString();
      if (!sel) {
        terminal.$input.focus();
      }
    });


    terminal.updatePrompt = function (prompt) {
      terminal.options.prompt = prompt;
      terminal.$prompt.innerHTML = prompt + "&nbsp;";
    };

    terminal.append = function (output) {
      WebTerm._appendOutput(output, terminal.options, terminal.dom);
    };

    terminal.cls = function () {
      while (terminal.$output.firstChild) {
        terminal.$output.removeChild(terminal.$output.firstChild);
      }
    };

    terminal.setStatus = function (val) {
      if (val == "") {
        terminal.$status.innerHTML = "";

      } else {
        terminal.$status.innerHTML = val;
      }
    };

    terminal.setHeader = function (val) {
      if (val == "") {
        terminal.$header.innerHTML = "";

      } else {
        terminal.$header.innerHTML = val;
      }
    };

    terminal.setHeaderVisible = function (val) {
      if (val) {
        console.log("setting display header show");
        terminal.$header.style.display = "block";

      } else {
        console.log("setting display header hide");
        terminal.$header.style.display = "none";
      }
    };

    terminal.setStatusVisible = function (val) {
      if (val) {
        console.log("setting display status show");
        terminal.$status.style.display = "block";

      } else {
        console.log("setting display status hide");
        terminal.$status.style.display = "none";
      }
    };

    terminal.handleCommand = function (val) {
      console.log('handleCommand - start');
      WebTerm._appendInput(val, terminal.options, terminal.dom);
      const output = processInput(val);
      WebTerm._processInput(output, terminal.options, terminal.dom);
      console.log('handleCommand - done');
    };


    terminal.handleInitialCommand = function (val) {
      WebTerm._appendInput(val, terminal.options, terminal.dom);
      terminal.$input.value = '';

      terminal.options.processInput(val);
      //console.log('termpage handleInitialCommand processInput result: '+output);
      //WebTerm._processInput(output, terminal.options, terminal.dom);
      //WebTerm._processInput(output, terminal.options, terminal.dom);

    };


    terminal.setInputCaretToEnd = function () {
      let pos = terminal.$input.value.length;
      setCaretPosition(terminal.dom.$input, pos, pos);
    };

    terminal.setInputCaretToHome = function () {
      setCaretPosition(terminal.dom.$input, 0, 0);
    };

    return terminal;
  }, // init


  link: (url, text) => {
    const res = (t) => `<a href="${url}" target="_blank">${t}</a>`;
    if (!text) {
      return (text) => {
        return res(text);
      };
    }
    return res(text);
  },

  color: (color, text) => {
    const res = (t) => `<span style="color:${color}">${t}</span>`;
    if (!text) {
      return (text) => {
        return res(text);
      };
    }
    return res(text);
  },

  replace: (text, changes) => {
    let response = text;
    Object.keys(changes).forEach(key => {
      response = response.replace(key, changes[key](key));
    });
    return response;
  },

};


function getCaretPosition(ctrl) {
  // IE < 9 Support
  if (document.selection) {
    ctrl.focus();
    let range = document.selection.createRange();
    let rangelen = range.text.length;
    range.moveStart('character', -ctrl.value.length);
    let start = range.text.length - rangelen;
    return {
      'start': start,
      'end': start + rangelen
    };
  } // IE >=9 and other browsers
  else if (ctrl.selectionStart || ctrl.selectionStart == '0') {
    return {
      'start': ctrl.selectionStart,
      'end': ctrl.selectionEnd
    };
  } else {
    return {
      'start': 0,
      'end': 0
    };
  }
};

function setCaretPosition(ctrl, start, end) {
  window.setTimeout(function () {

    // IE >= 9 and other browsers
    if (ctrl.setSelectionRange) {
      ctrl.focus();
      ctrl.setSelectionRange(start, end);
      console.log('setCaretPosition', {ctrl, start, end})
    }
    // IE < 9
    else if (ctrl.createTextRange) {
      let range = ctrl.createTextRange();
      range.collapse(true);
      range.moveEnd('character', end);
      range.moveStart('character', start);
      range.select();
    }
  }, 0);

};

