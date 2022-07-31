## Changelog for v0.3.1

### Bugfixes
- Fixed SegFaults due to concurrent acces to `github.com/stianeikeland/go-rpio` (the GPIO library this project uses)
- Those bugs would occur if using a `switchRF` after a `switchGPIO`
