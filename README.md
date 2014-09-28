CHIRP Current Track Announcer
-----------------------------

This little thing sits at the station and polls the ChIRP API every 30 seconds.
It parses the JSON response and sends that to the ProStream via UDP, so that the
Prostream can encode the current playing track into the IceCast stream.

Probably best to wrap this in an init.d script and maybe write a Cron job to 
make sure it keeps running. Honestly, though, you could easily set this and 
forget it. *famous last words*

Dependencies
============

None

How to get it
===============

```
go get github.com/agocs/chirp-current-track-announcer
```

How to test it
==============

```
go test
```

How to run it
=============

```
cagocs:chirp-current-track-announcer christopheragocs$ ./announcer --help
NAME:
   Announcer - Report current track to the Prostream

USAGE:
   Announcer [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --prostream 		IP address or hostname of the Prostream device.
   --port '9000'	Port of the Prostream track information receiver
   --chirp 		URL of the CHIRP current_playlist API endpoint
   --verbose		Run in Verbose mode.
   --test		Run in test mode. Sends nothing to Prostream
   --runOnce		Run once and then quit
   --version, -v	print the version
   --help, -h		show help
```


Examples:

```
./announcer --prostream 10.10.10.100 --port 9000 --chirp https://chirpradio.appspot.com/api/current_playlist --test --runOnce --verbose
```

This will run the service once, printing out useful information as it does so. 


```
./announcer --prostream 10.10.10.100 --port 9000 --chirp https://chirpradio.appspot.com/api/current_playlist
```

This will run the service in quiet mode forever. It will only print errors it encounters.
