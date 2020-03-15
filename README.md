# Dumproid

[![GitHub release](https://img.shields.io/github/v/release/tkmru/dumproid.svg)](https://github.com/tkmru/dumproid/releases/latest)
[![License: MIT](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](https://github.com/tkmru/dumproid/blob/master/LICENSE)
![](https://github.com/tkmru/dumproid/workflows/test/badge.svg)

Dumproid is Android process memory dump tool without ndk.
It is dumping memory from `/proc/<pid>/mem`.

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

When the target app is not debuggable and android is rooted:

```
$ adb shell
$ su
# /data/local/tmp/dumproid -p <PID>
               
██████╗ ██╗   ██╗███╗   ███╗██████╗ ██████╗  ██████╗ ██╗██████╗
██╔══██╗██║   ██║████╗ ████║██╔══██╗██╔══██╗██╔═══██╗██║██╔══██╗
██║  ██║██║   ██║██╔████╔██║██████╔╝██████╔╝██║   ██║██║██║  ██║
██║  ██║██║   ██║██║╚██╔╝██║██╔═══╝ ██╔══██╗██║   ██║██║██║  ██║
██████╔╝╚██████╔╝██║ ╚═╝ ██║██║     ██║  ██║╚██████╔╝██║██████╔╝
╚═════╝  ╚═════╝ ╚═╝     ╚═╝╚═╝     ╚═╝  ╚═╝ ╚═════╝ ╚═╝╚═════╝
```

When the target app is debuggable and android is not rooted:

```
$ adb shell
$ pm list packages # to check <target-package-name>
# run-as <target-package-name>
# cp /data/local/tmp/dumproid ./dumproid
# ./dumproid
               
██████╗ ██╗   ██╗███╗   ███╗██████╗ ██████╗  ██████╗ ██╗██████╗
██╔══██╗██║   ██║████╗ ████║██╔══██╗██╔══██╗██╔═══██╗██║██╔══██╗
██║  ██║██║   ██║██╔████╔██║██████╔╝██████╔╝██║   ██║██║██║  ██║
██║  ██║██║   ██║██║╚██╔╝██║██╔═══╝ ██╔══██╗██║   ██║██║██║  ██║
██████╔╝╚██████╔╝██║ ╚═╝ ██║██║     ██║  ██║╚██████╔╝██║██████╔╝
╚═════╝  ╚═════╝ ╚═╝     ╚═╝╚═╝     ╚═╝  ╚═╝ ╚═════╝ ╚═╝╚═════╝
```

### Dump memory
#### Dump To File

Permissions like `rwxs` can be specified as a filter.
By default, files are dumped under `/data/local/tmp/`.

```
sargo:/ # /data/local/tmp/dumproid -q -p 24264 --filter rw-p 
Output Dir: /data/local/tmp/20200315194818
  Dump File: 12c00000-131c0000__dev_ashmem_dalvik-main_space_(region_space)_(deleted)
  Dump File: 13340000-2ac00000__dev_ashmem_dalvik-main_space_(region_space)_(deleted)
  Dump File: 6f181000-6f3a6000__data_dalvik-cache_arm_system@framework@boot.art
  Dump File: 6f3bc000-6f4b3000__data_dalvik-cache_arm_system@framework@boot-core-libart.art
  Dump File: 6f4c5000-6f4f6000__data_dalvik-cache_arm_system@framework@boot-conscrypt.art
  Dump File: 6f4f9000-6f526000__data_dalvik-cache_arm_system@framework@boot-okhttp.art
  Dump File: 6f529000-6f57f000__data_dalvik-cache_arm_system@framework@boot-bouncycastle.art
  Dump File: 6f586000-6f5db000__data_dalvik-cache_arm_system@framework@boot-apache-xml.art
  Dump File: 6f5e2000-6f61d000__data_dalvik-cache_arm_system@framework@boot-ext.art
  Dump File: 6f628000-6fe2a000__data_dalvik-cache_arm_system@framework@boot-framework.art
  Dump File: 6fe8a000-6ff6c000__data_dalvik-cache_arm_system@framework@boot-telephony-common.art
  Dump File: 6ff7e000-6ff89000__data_dalvik-cache_arm_system@framework@boot-voip-common.art
  Dump File: 6ff8b000-6ffa0000__data_dalvik-cache_arm_system@framework@boot-ims-common.art
  Dump File: 6ffa2000-6ffa5000__data_dalvik-cache_arm_system@framework@boot-android.hidl.base-V1.0-java.art
  Dump File: 6ffa5000-6ffa9000__data_dalvik-cache_arm_system@framework@boot-android.hidl.manager-V1.0-java.art
  Dump File: 6ffab000-6ffac000__data_dalvik-cache_arm_system@framework@boot-framework-oahl-backward-compatibility.art
  Dump File: 6ffad000-6ffb0000__data_dalvik-cache_arm_system@framework@boot-android.test.base.art
  Dump File: 70365000-70366000_[anon:.bss]
  Dump File: 707e5000-707e6000__system_framework_arm_boot.oat
```

Transfer dumped files to your PC using `adb pull`:

```
$ adb pull /data/local/tmp/20200315194818 
/data/local/tmp/20200315194818/: 736 files pulled. 30.0 MB/s (583184384 bytes in 18.552s)
```

#### Print hexdump

Use the dump option to display memory like a hexdump.

```
sargo:/ # /data/local/tmp/dumproid -q -p 24264 -a 0xf0c9e000 --dump                                                                                                         
00000000  00 40 00 00 d0 60 b7 f0  01 00 00 00 14 71 b7 f0  |.@...`.......q..|
00000010  2d 33 bf f0 00 00 00 00  00 00 00 00 1c e0 c9 f0  |-3..............|
00000020  2f 73 79 73 74 65 6d 2f  62 69 6e 2f 6c 69 6e 6b  |/system/bin/link|
00000030  65 72 00 00 1d 00 00 00  02 00 00 00 00 10 00 00  |er..............|
00000040  40 e0 c9 f0 40 e0 c9 f0  35 d7 c2 f0 4c e0 c9 f0  |@...@...5...L...|
00000050  4c e0 c9 f0 00 00 00 00  00 00 00 00 ca 82 c8 f0  |L...............|
00000060  00 00 00 00 ff ff ff ff  00 00 00 00 e1 82 c8 f0  |................|
00000070  00 00 00 00 ff ff ff ff  00 00 00 00 95 26 c3 f0  |.............&..|
00000080  00 00 00 00 00 00 00 00  f3 82 c8 f0 00 00 00 00  |................|
00000090  ff ff ff ff fe 00 00 00  09 83 c8 f0 00 00 00 00  |................|
000000a0  ff ff ff ff fe 00 00 00  59 27 c3 f0 ac e0 c9 f0  |........Y'......|
000000b0  ac e0 c9 f0 27 28 c8 f0  79 27 c3 f0 b1 27 c3 f0  |....'(..y'...'..|
000000c0  d5 27 c3 f0 f1 27 c3 f0  f5 28 c3 f0 61 29 c3 f0  |.'...'...(..a)..|
000000d0  c9 29 c3 f0 4d 2a c3 f0  ad 2a c3 f0 0d 2b c3 f0  |.)..M*...*...+..|
000000e0  1d 2b c3 f0 99 2b c3 f0  e8 e0 c9 f0 e8 e0 c9 f0  |.+...+..........|
000000f0  f0 e0 c9 f0 f0 e0 c9 f0  f8 e0 c9 f0 f8 e0 c9 f0  |................|
```

### Check memory mapping

Use the maps option to display memory mapping.

```
sargo:/ # /data/local/tmp/dumproid -q -p 24264 --maps --filter rw-p                                                                                                         
12c00000-131c0000 rw-p 00000000 00:05 23292                              /dev/ashmem/dalvik-main space (region space) (deleted)
13340000-2ac00000 rw-p 00740000 00:05 23292                              /dev/ashmem/dalvik-main space (region space) (deleted)
6f181000-6f3a6000 rw-p 00000000 fd:01 221                                /data/dalvik-cache/arm/system@framework@boot.art
6f3bc000-6f4b3000 rw-p 00000000 fd:01 229                                /data/dalvik-cache/arm/system@framework@boot-core-libart.art
6f4c5000-6f4f6000 rw-p 00000000 fd:01 232                                /data/dalvik-cache/arm/system@framework@boot-conscrypt.art
6f4f9000-6f526000 rw-p 00000000 fd:01 235                                /data/dalvik-cache/arm/system@framework@boot-okhttp.art
6f529000-6f57f000 rw-p 00000000 fd:01 240                                /data/dalvik-cache/arm/system@framework@boot-bouncycastle.art
6f586000-6f5db000 rw-p 00000000 fd:01 250                                /data/dalvik-cache/arm/system@framework@boot-apache-xml.art
6f5e2000-6f61d000 rw-p 00000000 fd:01 263                                /data/dalvik-cache/arm/system@framework@boot-ext.art
6f628000-6fe2a000 rw-p 00000000 fd:01 270                                /data/dalvik-cache/arm/system@framework@boot-framework.art
6fe8a000-6ff6c000 rw-p 00000000 fd:01 275                                /data/dalvik-cache/arm/system@framework@boot-telephony-common.art
6ff7e000-6ff89000 rw-p 00000000 fd:01 278                                /data/dalvik-cache/arm/system@framework@boot-voip-common.art
6ff8b000-6ffa0000 rw-p 00000000 fd:01 281                                /data/dalvik-cache/arm/system@framework@boot-ims-common.art
6ffa2000-6ffa5000 rw-p 00000000 fd:01 284                                /data/dalvik-cache/arm/system@framework@boot-android.hidl.base-V1.0-java.art
6ffa5000-6ffa9000 rw-p 00000000 fd:01 287                                /data/dalvik-cache/arm/system@framework@boot-android.hidl.manager-V1.0-java.art
6ffab000-6ffac000 rw-p 00000000 fd:01 290                                /data/dalvik-cache/arm/system@framework@boot-framework-oahl-backward-compatibility.art
6ffad000-6ffb0000 rw-p 00000000 fd:01 293                                /data/dalvik-cache/arm/system@framework@boot-android.test.base.art
70365000-70366000 rw-p 00000000 00:00 0                                  [anon:.bss]
707e5000-707e6000 rw-p 003b4000 103:25 603                               /system/framework/arm/boot.oat
70967000-70968000 rw-p 00000000 00:00 0                                  [anon:.bss]
70c61000-70c62000 rw-p 00182000 103:25 601                               /system/framework/arm/boot-core-libart.oat
...
```

## License

GPLv3 - GNU General Public License, version 3

Copyright (C) 2020 tkmru
