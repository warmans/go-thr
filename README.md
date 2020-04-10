Experimental code needed for THRII footswitch POC.

### Setup (Ubuntu)

Requires `libportmidi-dev`

```
sudo apt-get install libportmidi-dev
```

If when you try and build you get errors like `undefined reference to symbol 'Pt_Start'`

Just get a package from somewhere else. It seems the package is broken in some versions of Ubuntu.

e.g. https://debian.pkgs.org/9/debian-main-amd64/libportmidi-dev_217-6_amd64.deb.html

