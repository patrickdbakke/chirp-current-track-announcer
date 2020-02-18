CHIRP Current Track Announcer
=============================

This little thing sits at the station and polls the ChIRP API every 30 seconds.
It parses the JSON response and sends that to the ProStream via UDP, so that the
Prostream can encode the current playing track into the IceCast stream.

Probably best to wrap this in an init.d script and maybe write a Cron job to 
make sure it keeps running. Honestly, though, you could easily set this and 
forget it. *famous last words*

Dependencies
------------

github.com/codegangsta/cli

How to get it
---------------

```
go get github.com/agocs/chirp-current-track-announcer
```

How to test it
--------------

```
go test
```

How to run it
-------------

```
$ ./announcer --help
NAME:
   Announcer - Report current track to the Prostream

USAGE:
   announcer [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --prostream value  IP address or hostname of the Prostream device.
   --port value       Port of the Prostream track information receiver (default: 9000)
   --chirp value      URL of the CHIRP current_playlist API endpoint (default: "https://chirpradio.appspot.com/api/current_playlist")
   --verbose          Run in Verbose mode.
   --test             Run in test mode. Sends nothing to Prostream
   --runOnce          Run once and then quit
   --rds value        IP address or hostname of the RDS Encoder
   --pdsPort value    Port used by the RDS Encoder (default: 23)
   --help, -h         show help
   --version, -v      print the version

```


Examples:

```
./announcer --prostream 10.10.10.100
            --port 9000
            --chirp https://chirpradio.appspot.com/api/current_playlist
            --rds 10.10.10.101
            --rdsPort 23
            --test
            --runOnce
            --verbose
```

This will run the service once, printing out useful information as it does so. 


```
./announcer --prostream 10.10.10.100
            --port 9000
            --chirp https://chirpradio.appspot.com/api/current_playlist
            --rds 10.10.10.101
            --rdsPort 23
```

This will run the service in quiet mode forever. It will only print errors it encounters.


How to run it during development
--------------------------------

Install dependencies
```
go get github.com/urfave/cli
go get bufio
go get sync
```

Run it
```
go run announcer.go
```
