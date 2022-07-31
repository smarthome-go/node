## Changelog for v0.3.0

### Generic GPIO devices
- Added the `switchesIgnore` list
- It is meant to contain switches present on *other* Hardware nodes in the network which cannot be handled by this particular node.
- Requests referencing this switch will not result in an error but will be ignored entirely.
- *Note*: Entering switches which exist in either `switchesRF` or `switchesGPIO` has absolutely no effect and is therefore not recommended.

The new configuration file looks like this:
```json
{
	"port": 8081,
	"nodeName": "localhost",
	"tokenHash": "$2a$05$ZV/6K.KUab6h327xWhfmwOFZAXNGwXPmA2ayoB4.zakij/iQL4uny",
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
	],
	"switchesIgnore": [
		"s3"
	]
}
```
