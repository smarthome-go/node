# smarthome-hw
 Hardware interface for the Smarthome server


### This project is deprecated and will be rewritten in a more reliable language soon.
The project is currently in the rewrite phase.

### Token
In order to guarantee a safe communication between the `smarthome` server and the `smarthome-hw` server, a token is required.
When this application is first started, a *random* token will be generated and printed to the server's logs (**but not to file**).
#### Change token
In order to change the default token, use the provided bash script:
```bash
./update_token.sh "old_token" "new_token"
```