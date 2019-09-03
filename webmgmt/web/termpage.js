/*
Original inspiration taken from https://github.com/tautvilas/termpage

Added ability to:
enable/disable adding to history
enable/disable display input when submitted.
set prompt
return reference to the terminal in the init function.

 */


const Termpage = {

    defaultOptions: {
        prompt: '$',
        initialCommand: false,
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
        ;
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
                    Termpage._appendInput($command.innerText, options, dom);
                    const out = options.processInput($command.innerText);
                    Termpage._processInput(out, options, dom);
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
            //console.log('_processInput.4');
            output.then((out) => {
                //console.log('_processInput.5');
                pre.remove();
                dom.$inputBlock.setAttribute('style', 'display:flex');
                Termpage._appendOutput(out, options, dom);
                //console.log('_processInput.6');
            });

            output.catch(() => {
                //console.log('_processInput.7');
                pre.remove();
                dom.$inputBlock.setAttribute('style', 'display:flex');
                Termpage._appendOutput(Termpage.color('red', 'command resolution failed'), options, dom);
                //console.log('_processInput.8');
            });
        } else {
            //console.log('_processInput.9  output: ', output);
            Termpage._appendOutput(output, options, dom);
            //console.log('_processInput.10');
        }
    },

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

    init: ($winElement, processInput, options = {}) => {
        let terminal = {};
        terminal.history = [];
        terminal.historyIndex = 0;
        terminal.options = Object.assign({}, Termpage.defaultOptions, options);
        terminal.$output = document.createElement("div");
        $winElement.appendChild(terminal.$output);

        terminal.prompt = terminal.options.prompt || Termpage.defaultOptions.prompt;
        terminal.$prompt = document.createElement("span");
        terminal.$prompt.className = "termpage-prompt";
        terminal.$prompt.innerHTML = terminal.prompt + "&nbsp;";

        terminal.$input = document.createElement("input");
        terminal.$input.setAttribute("type", "text");
        terminal.$input.className = "termpage-input";

        terminal.$inputBlock = document.createElement("p");
        terminal.$inputBlock.className = "termpage-block termpage-input-block";

        terminal.$inputBlock.appendChild(terminal.$prompt);
        terminal.$inputBlock.appendChild(terminal.$input);
        $winElement.appendChild(terminal.$inputBlock);

        terminal.dom = {
            $winElement,
            $inputBlock: terminal.$inputBlock,
            $input: terminal.$input,
            $output: terminal.$output,
        };

        terminal.options.processInput = (inp) => {
            if (terminal.options.appendHistory) {
                terminal.historyIndex = 0;
                terminal.history.push(inp);
            }
            return processInput(inp);
        };

        if (terminal.options.initialCommand) {
            const output = processInput(terminal.options.initialCommand);
            Termpage._appendInput(terminal.options.initialCommand, terminal.options, terminal.dom);
            Termpage._processInput(output, terminal.options, terminal.dom);
        }

        terminal.$input.addEventListener('keydown', function (e) {
            var key = e.which || e.keyCode;
            if (key === 13) { // 13 is enter
                const input = e.srcElement.value;
                const output = terminal.options.processInput(input);

                //console.log('input: ', input);
                //console.log('output: ', output);
                Termpage._appendInput(input, terminal.options, terminal.dom);
                Termpage._processInput(output, terminal.options, terminal.dom);
                terminal.$input.value = '';
            }
            if (key === 38) { // up
                const val = terminal.history[terminal.history.length - terminal.historyIndex - 1];
                if (val) {
                    terminal.historyIndex++;
                    terminal.dom.$input.value = val;
                }
            } else if (key === 40) { // down
                const val = terminal.history[terminal.history.length - terminal.historyIndex + 1];
                if (val) {
                    terminal.historyIndex--;
                    terminal.dom.$input.value = val;
                }
            }
        });

        if (terminal.options.autoFocus) {
            terminal.$input.focus();
        }
        $winElement.addEventListener("click", function () {
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
            Termpage._appendOutput(output, terminal.options, terminal.dom);
        };

        return terminal;
    }
};

(() => {
    const styles = `
html, body {
  margin: 0;
  padding: 0;
}

.termpage-window {
  overflow-y: auto;
}

.termpage-block {
  margin: 0;
  padding: 0;
}

pre.termpage-block {
  word-break: keep-all;
  white-space: pre-wrap !important;
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

.termpage-menu {
  display: flex;
  list-style: none;
  margin: 0;
  padding: 0;
}

.termpage-loader::before {
  content: '';
  animation: termpage-loader 0.5s infinite;
}

@keyframes termpage-loader {
  0% {
    content: '';
  }
  25% {
    content: '.';
  }
  50% {
    content: '..';
  }
  75% {
    content: '...';
  }
}
  `
    const styleSheet = document.createElement("style")
    styleSheet.type = "text/css"
    styleSheet.innerText = styles
    document.head.prepend(styleSheet)
    const theme = `
.termpage-window {
  background-color: black;
  border: 2px solid #888;
  padding-top: 5px;
}

.termpage-window * {
  font-family: "Courier New", Courier, monospace;
  font-size: 16px;
  color: #ddd;
}

.termpage-input {
  background-color: #222;
  color: #ddd;
  caret-color: white;
}

.termpage-block, .termpage-input {
  line-height: 20px;
}

.termpage-block {
  padding-left: 5px;
  padding-right: 5px;
}

.termpage-window a {
  background-color: #888;
  text-decoration: none;
  cursor:pointer;
}
.termpage-window a:hover {
  background-color: #333;
}

.termpage-menu {
  background-color: #888;
}

.termpage-menu li:hover {
  background-color: #666;
  cursor: pointer;
}
  `
    const themeSheet = document.createElement("style")
    themeSheet.type = "text/css"
    themeSheet.innerText = theme;
    document.head.prepend(themeSheet)
})();

