# Dumproid

Android process memory dump tool without ndk.

## How to Build
After the build is complete, if adb is connected, place the built binary in `/data/local/tmp/` on Android.

```
$ make
GOOS=linux GOARCH=arm64 GOARM=7 go build -o dumproid
/bin/sh -c "adb push dumproid /data/local/tmp/dumproid"
dumproid: 1 file pushed. 24.1 MB/s (4977746 bytes in 0.197s)
```

## Usage

Push dumproid to android using adb.

```
$ adb push dumproid /data/local/tmp/dumproid
```

### Start-up

When the target app is debuggable and android is not rooted:

```
$ adb shell
$ pm list packages # to check <target-package-name>
# run-as <target-package-name>
# cp /data/local/tmp/dumproid ./dumproid
# ./dumproid
```

When the target app is not debuggable and android is rooted:

```
$ adb shell
$ su
# /data/local/tmp/dumproid -p <PID>
```

### Dump memory

### Check memory mapping

## License

GPLv3 - GNU General Public License, version 3

Copyright (C) 2020 tkmru
