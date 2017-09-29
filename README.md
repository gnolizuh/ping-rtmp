# ping-rtmp
> RTMP layer PingPong test.

## Install and Run
```
go get -u github.com/gnolizuh/ping-rtmp
cd ping-rtmp
go build -i -o ping-rtmp .
./ping-rtmp --push rtmp://yourip/live/ --pull rtmp://yourip/live/
```

## Usage

```
$ ./ping-rtmp -h
NAME:
   ping-rtmp - RTMP layer PingPong test.

USAGE:
   config [global options] command [command options] [arguments...]

VERSION:
   1.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --push value   specify a push url for publisher (default: "rtmp://127.0.0.1/live/")
   --pull value   specify a pull url for player (default: "rtmp://127.0.0.1/live/")
   --help, -h     show help
   --version, -v  print the version
```

## Result

```
PING rtmp://127.0.0.1/live/ <-> rtmp://127.0.0.1/live/: 32 data bytes
32 bytes from: sid=1 csid=7 time=8.875 ms
32 bytes from: sid=1 csid=7 time=9.109 ms
32 bytes from: sid=1 csid=7 time=9.003 ms
32 bytes from: sid=1 csid=7 time=8.195 ms
32 bytes from: sid=1 csid=7 time=10.355 ms
32 bytes from: sid=1 csid=7 time=10.197 ms
32 bytes from: sid=1 csid=7 time=8.220 ms
32 bytes from: sid=1 csid=7 time=8.327 ms
32 bytes from: sid=1 csid=7 time=8.663 ms
```
