## Changelog for v0.1.0

### Generic GPIO devices
- Added support for generic GPIO devices
- Alongside the preexisting 433mhz hardware, you can now leverage the power of any GPIO pin to be controllable via Smarthome

The new configuration file looks like this:
```json
{
	"port": 8081,
	"nodeName": "localhost",
	"tokenHash": "$2a$05$T9NsGTrr847RI3yianF90Oi7PxUxANUybxGjSDgtlw3g6HI44xRbO",
	"hardware": {
		"hardwareEnabled": false,
		"pin": 0,
		"protocol": 1,
		"repeat": 10,
		"pulseLength": 180,
		"contentLength": 24
	},
	"switchesRF": [
		{
			"id": "s1",
			"on": 0,
			"off": 0
		}
	],
	"switchesGPIO": [
		{
			"id": "s2",
			"pin": 1
		}
	]
}
```
