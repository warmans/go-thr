# NOTE: As far as I know this code no longer works. It seems each update to the amp may break commands being sent over sysex. But I will leave it up in case it helps anyone start their own project.

Experimental code needed for THRII footswitch POC.

### Prerequisites (Ubuntu)

Requires `libportmidi-dev`

```
sudo apt-get install libportmidi-dev
```

If when you try and build you get errors like `undefined reference to symbol 'Pt_Start'`

Just get a package from somewhere else. It seems the package is broken in some versions of Ubuntu.

e.g. https://debian.pkgs.org/9/debian-main-amd64/libportmidi-dev_217-6_amd64.deb.html

### THRII Sysex Format

Messages are in the following format (as far as I've been able to decode their meaning):

|                 | Num. Bytes | Example    | Description  
|-----------------|------------|------------|-------------------------------------
| Start           | 1          | `f0`       | standard sysex start byte.
| Manufacturer ID | 3          | `00 01 0c` | Yamaha (line6) extended manufacturer code
| Preamble        | 3          | `24 02 4d` | Some kind of device ID I guess. This never changes on my device.
| ?               | 1          | `00`       | Seems to be some kind of command grouping. Usually 00 or 01.
| Sequence Num.   | 1          | `01`       | This gets incremented for each command sent. But seems to be independent between "groups" (previous byte).
| Payload Desc.   | 3          | `00 00 03` | These three bytes seem to describe the payload (e.g. length). 
| Payload         | ?          | `...`      | The payload size seems to depend on the previous bytes and contains mysterious data in an unknown format.
| End             | 1          | `f7`       | Standard sysex end byte.


### Notes on reverse engineering the THR II

The device uses standard MIDI over USB for the desktop app. This traffic can be recorded on windows
using [usbpcap](https://desowin.org/usbpcap/) and then viewed in Wireshark.

You can interact with the amp using raw sysex messages using Bome SendSX.

Pressing buttons on the amp will emit events (e.g. that you can view on sendsx) HOWEVER
this is only true if you first of all tell the amp to do so, either by opening and closing the app first,
or by sending the command see (messages.go).
 
