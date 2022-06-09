# Node (formerly *Smarthome-hw*)
 Hardware interface for the Smarthome server
 
 ### Purpose
 The Smarthome server acts like a global hub to connect services, people and hardware together.
 Because the hub is a centralized server which should be able to run on any hardware, not just the Raspberry-Pi, a master-slave setup is used to control the power outlets which accessible to smarthome users. In this setup, the Smarthome-hub acts as the master which is able to orchistrate the slaves, in this case the hardware-nodes.
 
 The 433mhz sockets are controlled by a physical sender which is attached to each Raspberry-Pi, making it a hardware-node.
 The nodes can be registered in the smarthome server which is then able to manage the nodes, send power commands to them and monitor their health.
 
 Therefore, *node*  provides a *REST-API* for making its hardware accessible to the Smarthome-hub.
 
 For redundancy and increased reliability, a Smarthome network should contain more than just one hardware-node.
 Redundancy, fallback and latency are handled by the Smarthome-hub which automatically detects if nodes fail or take too long to respond.
 
Because a single node is not able to serve concurrent requests (*due to the limitations of the attached hardware*), a locking system is implemented in node.
This lock is acquired when a request is dispatched.
The lock is always released when no longer required, regardless whether the request was successful or not.

In order to provide a true concurrent, non-blocking access to the power outlets, the Smarthome-hub contains an internal queueing system which manages concurrent requests and executes them in a blocking manner (one after another).

![smarthome-hw logo](./icon/readme.png)

### Tokens
In order to guarantee a safe communication between the Smarthome-hub and the node, a token is required.
When this application is first started, a *random* token will be generated and printed to the node's logs (**not to the log file**).
#### Change Token
In order to change the default token, use the provided bash script which is located in the distribution directory.
```bash
./update_token.sh "old_token" "new_token"
```

### Api
To dispatch a power request, perform a **POST** request to a similar URL, using a similar request body.
The host and the `switch` should be modified to suit your usecase.
The payload has to be encoded as `application/json`, using `application/json` as the `Content-Type`.
```
http://localhost:8081/power?token=smarthome
```

```json
{
	"switch": "s1",
	"power": true
}
```
