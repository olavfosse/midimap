This program is of pre-release status, but it is public either way. Also since this is my first Go project, don't expect the code to be very idiomatic; advice and suggestions are always appreciated :-).

---

# midimap
Map MIDI events to simulated keyboard events.
## Setup
midimap can be compiled to use either portmidi or rtmidi. The following sections describe how to build midimap with one of these libraries, you should only follow the section for your library of choice.

If you are confused with regards to which library to use, I recommend rtmidi for most users.
### rtmidi
#### Install non-Go dependencies
```sh
# Debian derivatives
sudo apt install libasound2-dev
# MacOS
# TODO: Add possible macos rtmidi dependencios
```
#### Build
```sh
$ go build -tags rtmidi
```
### portmidi
#### Install non-Go dependencies
```sh
# Debian derivatives
sudo apt install libportmidi-dev
# MacOS
brew install portmidi
```
#### Build
```sh
$ go build -tags portmidi
```
## Alternatives
There are several other MIDI to keypress programs, but none of them are sufficient for my use case. Notably there is no single alternative which is both open source, cross platform and built with a efficient and pleasant stack(e.g no python or electron ;)). I also have ambitions outside of these critera, but for now these are the main advantages.

| Alternative                                                                         | License                              | Platform     | Comment          |
|-------------------------------------------------------------------------------------|--------------------------------------|--------------|------------------|
| [Midikey2Key](https://midikey2key.de)                                               | Propietary freeware                  | Windows only |                  |
| [xobs/midi-to-keypress](https://github.com/xobs/midi-to-keypress)                   | Apache License 2.0                   | Windows only |                  |
| [davidlukerice/midi-to-keypress](https://github.com/davidlukerice/midi-to-keypress) | MIT                                  | Windows only | Uses Electron ;( |
| [mwicat/midimap](https://github.com/mwicat/midimap)                                 | Source available propietary freeware | Windows only | Uses Python ;(   |
