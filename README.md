## WebMgmt  

#### Details
this project will allow for a web admin service to be embedded into a service. This will allow clients to open a browser to the access port and login to the service. Commands can be developed to access the server.

#### Web
html terminal borrowed from  https://github.com/tautvilas/termpage



#### Building & Running
```bash

make binaries            - make bianries into bin dir
cd bin
./example

open browser to:  http://localhost:1099

You can login with username: alex password: bambam
user                     - display user info
http                     - display http headers/cookies
prompt                   - chanmge prompt
link                     - respond with a link that can be clicked on in terminal window
ticker                   - display a ticker that updates info to window periodically
image                    - display an image in the terminal window
raw                      - dsiplay an image along with a command tool bar of commands that can be clicked on
commands                 - display command tool bar of commands that can be clicked on
history                  - display history of commands executed

any unknown command will be echoed back
```


