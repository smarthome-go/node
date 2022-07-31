## Changelog for v0.2.0

### Generic GPIO devices
- Added the `invert` setting to GPIO hardware which internally inverts every power request
- This means an `on` will be changed to a `off` and the other way around
- This setting is especially useful for controlling GPIO-relays

The new configuration file looks like this:
```json
{
	"port": 8081,
	"nodeName": "localhost",
	"tokenHash": "$2a$05$/c/qAd6gfSh0HSXuQTH4EOiEWfqarjze/y4UdlLOdZFgqua.KUxZe",
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
			"pin": 1,
			"invert": false
		}
	]
}
```
