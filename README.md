This is not yet useable as a general purpose program since it requires hardcoding the midi event data1s and keycodes in the source code. Also since this is my first Go program, don't expect the code to be very idiomatic.

---

# midimap
Map midi events to simulated keyboard events.
## Setup
### Install portmidi
```sh
# Debian derivatives
sudo apt install libportmidi-dev
# MacOS
brew install portmidi
```
### Install go packages
```sh
go get github.com/micmonay/keybd_event github.com/rakyll/portmidi
```
### Verify that portmidi is working
```sh
# The following command outputs the midi events of the default midi device to stdout.
# Try sending some midi events and see if they are outputted to test that portmidi is working correctly
go run midilog.go
```

## Alternatives
| Alternative                                                                         | License                              | Platform     | Comment          |
|-------------------------------------------------------------------------------------|--------------------------------------|--------------|------------------|
| [Midikey2Key](https://midikey2key.de)                                               | Propietary freeware                  | Windows only |                  |
| [xobs/midi-to-keypress](https://github.com/xobs/midi-to-keypress)                   | Source available propietary freeware | Windows only |                  |
| [davidlukerice/midi-to-keypress](https://github.com/davidlukerice/midi-to-keypress) | MIT                                  | Windows only | Uses Electron ;( |
| [mwicat/midimap](https://github.com/mwicat/midimap)                                 | Source available propietary freeware | Windows only | Uses Python ;(   |
