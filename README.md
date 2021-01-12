This program is of pre-release status, but it is public either way. Also since this is my first Go program, don't expect the code to be very idiomatic, advice and suggestions is always appreciated :-).

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
#### Install Go dependencies
```sh
go get github.com/micmonay/keybd_event gitlab.com/gomidi/midi gitlab.com/gomidi/midi/reader gitlab.com/gomidi/rtmididrv
```
#### Compile
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
#### Install Go dependencies
go get github.com/micmonay/keybd_event gitlab.com/gomidi/midi gitlab.com/gomidi/midi/reader gitlab.com/gomidi/portmididrv
```
#### Compile
```sh
$ go build -tags portmidi
```
## Usage
midimap is invoked on the command line along with one argument `map`, a path to a "midimap-lang" file which describes which MIDI events to map to which simulated keypresses. The "midimap-lang" section explains midimap-lang in more depth.

### Choosing a device
When midimap is invoked, it loads `map` and prompts you for which MIDI device to listen to. Note that depending on your OS and how things are configured, your current user may not have permissons to access the MIDI devices. If that is the case an error will be displayed informing you of this. To circumvent this simply run the command with elevated privleges using sudo, doas or the-like.
```sh
$ sudo ./midimap ./maps/taiko.mml
0: &{ALSA Midi Through Port-0 %!s(bool=false) %!s(bool=true) %!s(bool=false)}
1: &{ALSA Midi Through Port-0 %!s(bool=true) %!s(bool=false) %!s(bool=false)}
2: &{ALSA TD-1 MIDI 1 %!s(bool=false) %!s(bool=true) %!s(bool=false)}
3: &{ALSA TD-1 MIDI 1 %!s(bool=true) %!s(bool=false) %!s(bool=false)}
{0,1,2,3}: 
```

To select a device simply input the integer corresponding to your device of choice. Note that some of the devices are "fake" and will crash midimap instantaneously if selected. In this case device 0 and 2 are fake. This is a bug and will get fixed eventually. To work around this simply try selecting devices until you find one that works. The integer and device pairs are persistent cross session and even after connecting/disconnecting your devices.

### ALSA false flag error
On Linux the following error message will be spammed after starting the program.
```
ALSA lib seq.c:4176:(snd_seq_event_input_feed) poll: Interrupted system call
```
This is an error is a false flag and does not indicate any failure in the program. It can be ignorred safely by redirecting stderr to /dev/null.

This is also a bug, that will be fixed eventually. Despite these bugs the program is perfectly useable. These idiocacies will be fixed before offical release.

### midimap-lang
midimap-lang is a domain specific programming language created specifically for midimap. It consists of a series of mappings and comments. A mapping consists of a key code and a conditon, called its matcher. When midimap receives an incoming MIDI event it simulates pressing the key code of every mapping which has a matcher that matches the incoming event.

For learning how to write midimap-lang maps, see the specification [midimap-lang-spec.md](https://github.com/fossegrim/midimap/blob/master/midimap-lang-spec.md) and the example maps in [maps/](https://github.com/fossegrim/midimap/tree/master/maps). At the moment there is no tutorial, but reading through these should be sufficient to read and write any midimap map.

## Alternatives
| Alternative                                                                         | License                              | Platform     | Comment          |
|-------------------------------------------------------------------------------------|--------------------------------------|--------------|------------------|
| [Midikey2Key](https://midikey2key.de)                                               | Propietary freeware                  | Windows only |                  |
| [xobs/midi-to-keypress](https://github.com/xobs/midi-to-keypress)                   | Apache License 2.0 | Windows only |                  |
| [davidlukerice/midi-to-keypress](https://github.com/davidlukerice/midi-to-keypress) | MIT                                  | Windows only | Uses Electron ;( |
| [mwicat/midimap](https://github.com/mwicat/midimap)                                 | Source available propietary freeware | Windows only | Uses Python ;(   |
