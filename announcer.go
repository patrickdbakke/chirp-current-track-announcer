package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

/*
	This sumbitch is a little service that has to run *at* CHIRP.
	It reads the current track from the CHIRP API, and then sends
	it over to the Prostream.
*/

type track struct {
	Dj     string
	Artist string
	Track  string
	Label  string
}

type playlist struct {
	Now_playing     track
	Recently_played []track
}

func main() {
	// Entry point to the application
	// Boilerplate CLI stuff
	app := cli.NewApp()
	app.Name = "Announcer"
	app.Usage = "Report current track to the Prostream"
	app.Flags = []cli.Flag{
		cli.StringFlag{"prostream", "", "IP address or hostname of the Prostream device."},
		cli.IntFlag{"port", 9000, "Port of the Prostream track information reciever"},
		cli.StringFlag{"chirp", "", "URL of the CHIRP current_playlist API endpoint"},
		cli.BoolFlag{"verbose", "Run in Verbose mode."},
		cli.BoolFlag{"test", "Run in test mode. Sends nothing to Prostream"},
		cli.BoolFlag{"runOnce", "Run once and then quit"},
	}
	app.Action = func(c *cli.Context) {
		parseAndRun(c)
	}
	app.Run(os.Args)
}

func parseAndRun(c *cli.Context) {
	// Actually take the info we learned from the main() and do stuff with it.

	prostream := c.String("prostream")
	prostreamPort := c.Int("port")

	chirpApi := c.String("chirp")

	verbose := c.Bool("verbose")
	test := c.Bool("test")
	runOnce := c.Bool("runOnce")

	//validate our inputs

	if prostream == "" {
		println("You need to specify the address of the Prostream device.")
		os.Exit(1)
	}

	if chirpApi == "" {
		println("You need to specify a CHIRP API endpoint.")
		os.Exit(1)
	}

	grabAndSendData(prostream, prostreamPort, chirpApi, verbose, test, runOnce)
}

func grabAndSendData(prostream string, prostreamPort int, chirpApi string,
	verbose bool, test bool, runOnce bool) {
	// Here's the main loop for grabbing data from the API and sending to the Prostream.

	for {
		currentTrack := grabCurrentTrackInfo(chirpApi, verbose)
		sendCurrentTrackToProstream(currentTrack, prostream, prostreamPort, verbose)
		if runOnce {
			break //break out of the loop and exit
		} else {
			time.Sleep(30 * time.Second)
		}
	}

}

func grabCurrentTrackInfo(chirpApi string, verbose bool) track {
	niceUrl := makeNiceUrl(chirpApi)

	writemsg(fmt.Sprintf("About to GET %s", niceUrl), verbose)

	resp, err := http.Get(niceUrl)
	if err != nil {
		log.Printf("Error connecting to the CHIRP api: , %s\n", err.Error())
		return track{"", "", "", ""}
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	playlistBytes := buf.Bytes()

	return getTrackInfoFromJson(playlistBytes, verbose)
}

func makeNiceUrl(baseUrl string) string {
	return baseUrl + "?src=chirp-current-track-announcer"
}

func getTrackInfoFromJson(jsonPlaylist []byte, verbose bool) track {
	var playlistInfo playlist

	writemsg(fmt.Sprintf("About to unmarshall %d bytes of json", len(jsonPlaylist)), verbose)

	err := json.Unmarshal(jsonPlaylist, &playlistInfo)
	if err != nil {
		log.Printf("Error unmarshalling the playlist. %s", err.Error())
		return track{"", "", "", ""}
	}

	writemsg(fmt.Sprintf("Looks like %s is playing", playlistInfo.Now_playing.Track), verbose)

	return playlistInfo.Now_playing
}

func sendCurrentTrackToProstream(currentTrack track, prostream string, prostreamPort int, verbose bool) {
	niceAddress := fmt.Sprintf("%s:%d", prostream, prostreamPort)
	conn, err := net.Dial("udp", niceAddress)
	if err != nil {
		log.Printf("Error dialing the Prostream. %s", err.Error())
		return
	}

	writemsg("Established connection", verbose)

	//using the default format specified by the Prostream manual
	fmt.Fprint(conn, fmt.Sprintf("t=%s - %s | u=http://www.chirpradio.org\r\n",
		currentTrack.Track, currentTrack.Artist))

	//godspeed little udp packet
	conn.Close()
}

func writemsg(message string, verbose bool) {
	if verbose {
		log.Println(message)
	}
}
