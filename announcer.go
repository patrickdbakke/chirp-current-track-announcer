package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"


	"github.com/urfave/cli"
	"bufio"
	"sync"
)

/*
	This is a little service that has to run *at* CHIRP.
	It reads the current track from the CHIRP API, and then sends
	it over to the Prostream.
*/

//Track represents the information we want to send to the Prostream streaming device for inclusion in the internet stream.
type track struct {
	Dj     string
	Artist string
	Track  string
	Label  string
}

//Playlist is a go representation of the JSON playlist we will get from the CHIRP API
type playlist struct {
	Now_playing     track
	Recently_played []track
}

func main() {
	// Entry point to the application
	// Boilerplate CLI stuff
	app := &cli.App{
		Name: "Announcer",
		Usage: "Report current track to the Prostream",
		Flags: []cli.Flag{
			&cli.StringFlag{Name:"prostream",  Value:"", Usage:"IP address or hostname of the Prostream device."},
			&cli.IntFlag{Name:"port", Value:9000, Usage:"Port of the Prostream track information receiver"},
			&cli.StringFlag{Name:"chirp", Value:"https://chirpradio.appspot.com/api/current_playlist", Usage:"URL of the CHIRP current_playlist API endpoint"},
			&cli.BoolFlag{Name:"verbose", Usage:"Run in Verbose mode."},
			&cli.BoolFlag{Name:"test", Usage:"Run in test mode. Sends nothing to Prostream"},
			&cli.BoolFlag{Name:"runOnce", Usage:"Run once and then quit"},
			&cli.StringFlag{Name:"rds", Value:"", Usage:"IP address or hostname of the RDS Encoder"},
			&cli.IntFlag{Name:"rdsPort", Value:23, Usage:"Port used by the RDS Encoder"},
		},
		Action: func(c *cli.Context) error {
			parseAndRun(c)
			return nil
		},
	}
	app.Run(os.Args)
}

//ParseAndRun parses the config data we got from the command line and validates it.
func parseAndRun(c *cli.Context) {
	// Actually take the info we learned from the main() and do stuff with it.

	prostream := c.String("prostream")
	prostreamPort := c.Int("port")

	rds := c.String("rds")
	rdsport := c.Int("rdsPort")

	chirpApi := c.String("chirp")

	verbose := c.Bool("verbose")
	test := c.Bool("test")
	runOnce := c.Bool("runOnce")

	//validate our inputs


	if chirpApi == "" {
		println("You need to specify a CHIRP API endpoint.")
		os.Exit(1)
	}

	grabAndSendData(prostream, prostreamPort, rds, rdsport, chirpApi, verbose, test, runOnce)
}


//GrabAndSendData is the main loop of our program. In normal operation, it gets the current track info from the
//  CHIRP API and sends that to the Prostream once every 5 seconds.
//  The sends to Prostream and RDS happen concurrently, and grabAndSendData will block until
//  both are complete.
func grabAndSendData(prostream string, prostreamPort int, rds string, rdsPort int, chirpApi string,
	verbose bool, test bool, runOnce bool) {
	// Here's the main loop for grabbing data from the API and sending to the Prostream.

	for {
		currentTrack := grabCurrentTrackInfo(chirpApi, verbose)
		if !test {
			wg := sync.WaitGroup{}
			if prostream != "" {
				wg.Add(1)
				go sendCurrentTrackToProstream(currentTrack, prostream, prostreamPort, verbose, &wg)
			}
			if rds != "" {
				wg.Add(1)
				go sendCurrentTrackToRDS(currentTrack, rds, rdsPort, verbose, &wg)
			}
			wg.Wait()
		}
		if runOnce {
			break //break out of the loop and exit
		} else {
			time.Sleep(5 * time.Second)
		}
	}

}

//grabCurrentTrackInfo gets the current track info from the CHIRP API
func grabCurrentTrackInfo(chirpApi string, verbose bool) track {

	niceUrl := makeNiceUrl(chirpApi)
	
	writemsg(fmt.Sprintf("About to GET %s", niceUrl), verbose)
	resp, err := http.Get(niceUrl)

	if err != nil {
		fmt.Printf("Error connecting to the CHIRP api: , %s\n", err.Error())
		return track{"", "", "", ""}
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	playlistBytes := buf.Bytes()

	return getTrackInfoFromJson(playlistBytes, verbose)
}

//makeNiceUrl appends the correct argument onto the CHIRP API URL
func makeNiceUrl(baseUrl string) string {
	return baseUrl + "?src=chirp-current-track-announcer"
}

//getTrackInfoFromJSON parses the JSON response from the CHIRP API and returns the track now playing.
func getTrackInfoFromJson(jsonPlaylist []byte, verbose bool) track {
	var playlistInfo playlist

	writemsg(fmt.Sprintf("About to unmarshall %d bytes of json", len(jsonPlaylist)), verbose)

	err := json.Unmarshal(jsonPlaylist, &playlistInfo)
	if err != nil {
		fmt.Printf("Error unmarshalling the playlist. %s", err.Error())
		return track{"", "", "", ""}
	}

	writemsg(fmt.Sprintf("Looks like %s is playing", playlistInfo.Now_playing.Track), verbose)

	return playlistInfo.Now_playing
}

//sendCurrentTrackToProstream creates a connection with the Prostream device and sends it the current
// track info, formatted in such a way that it will display correctly to stream players.
func sendCurrentTrackToProstream(currentTrack track, prostream string, prostreamPort int, verbose bool, wg *sync.WaitGroup) {
	defer wg.Done()
	niceAddress := fmt.Sprintf("%s:%d", prostream, prostreamPort)
	conn, err := net.Dial("udp", niceAddress)
	
	if err != nil {
		fmt.Printf("Error dialing the Prostream. %s", err.Error())
		return
	}

	writemsg("Established connection", verbose)

	//using the default format specified by the Prostream manual
	fmt.Fprint(conn, fmt.Sprintf("t=%s - %s | u=http://www.chirpradio.org\r\n",
		currentTrack.Track, currentTrack.Artist))

	//godspeed little udp packet
	conn.Close()
}

//sendCurrentTrackToRDS establishes a connection with the RDS encoder and sends it the
// track info as a DPS message.
func sendCurrentTrackToRDS(currentTrack track, rds string, rdsPort int, verbose bool, wg *sync.WaitGroup){
	defer wg.Done()
	addressPort :=  fmt.Sprintf("%s:%d", rds, rdsPort)
	conn, err := net.Dial("tcp", addressPort)

	if err != nil {
		fmt.Printf("Error dialing the RDS. %+v", err)
		return
	}
	defer conn.Close()

	message := makeRDSMessage(currentTrack, verbose)
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Printf("Error writing message to RDS: %+v", err)
		return
	}
	sc := bufio.NewScanner(conn)

	for i := 0; i < 10 && !sc.Scan(); i++{
		time.Sleep(100 * time.Millisecond)
	}
	resp := sc.Text()

	if resp == "NO"{
		fmt.Printf("The RDS Encoder did not like the input %s", message)
	}
}

//makeRDSMessage should construct a properly formatted RDS 'DPS' message.
// Note that DPS messages can be no greater than 128 characters. If the message winds up
// being greater than 128 characters, we should truncate it and end it with `...`.
// #TODO: make unicode characters into their closest ASCII Latin1 equivalent chars
func makeRDSMessage(currentTrack track, verbose bool) string {
	baseMessage := fmt.Sprintf("'%s' by %s", currentTrack.Track, currentTrack.Artist)

	if len(baseMessage)> 128{
		baseMessage = baseMessage[:124] + "..."
	}

	stationID := " on CHIRP Radio"

	totalMessage := baseMessage
	if len(baseMessage) + len(stationID) < 128 {
		totalMessage = baseMessage + stationID
	}
	writemsg(fmt.Sprintf("Sending to RDS: %+s", "DPS="+totalMessage + "\n"), verbose);
	return "DPS="+totalMessage + "\n"
}

func writemsg(message string, verbose bool) {
	if verbose {
		fmt.Println(message)
	}
}
