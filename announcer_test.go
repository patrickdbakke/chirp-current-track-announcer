package main

import (
	"testing"
)

func TestMakeNiceUrl(t *testing.T) {
	baseUrl := "http://example.com/bloo/blah/blee"
	expected := "http://example.com/bloo/blah/blee?src=chirp-current-track-announcer"
	niceUrl := makeNiceUrl(baseUrl)
	if niceUrl != expected {
		t.Errorf("Expected %s, got %s", expected, niceUrl)
	}
}

func TestGetTrackInfoFromJson(t *testing.T) {
	someJsonBytes := []byte(`
{
  "now_playing": {
    "played_at_local_ts": 1411739291,
    "dj": "DJ Dead Alive - Ripped Sounds",
    "artist": "Homeboy Sandman",
    "track": "Activity",
    "notes": "",
    "artist_is_local": false,
    "label": "Stones Throw",
    "played_at_gmt_ts": 1411757291,
    "played_at_gmt": "2014-09-26T18:48:11.894660",
    "release": "Hallways",
    "played_at_local_expire": "2015-03-31T13:48:11.894660-05:00",
    "played_at_local": "2014-09-26T13:48:11.894660-05:00",
    "id": "ahBzfmNoaXJwcmFkaW8taHJkchoLEg1QbGF5bGlzdEV2ZW50GICAgNKJtqEKDA",
    "lastfm_urls": {
      "med_image": "http:\/\/userserve-ak.last.fm\/serve\/64s\/100835407.png",
      "sm_image": "http:\/\/userserve-ak.last.fm\/serve\/34s\/100835407.png",
      "_processed": true,
      "large_image": "http:\/\/userserve-ak.last.fm\/serve\/174s\/100835407.png"
    }
  },
  "recently_played": [
    {
      "played_at_local_ts": 1411739045,
      "dj": "DJ Dead Alive - Ripped Sounds",
      "artist": "Zammuto",
      "track": "Need Some Sun",
      "notes": "",
      "artist_is_local": false,
      "label": "Temporary Residence Ltd.",
      "played_at_gmt_ts": 1411757045,
      "played_at_gmt": "2014-09-26T18:44:05.699520",
      "release": "Anchor",
      "played_at_local_expire": "2015-03-31T13:44:05.699520-05:00",
      "played_at_local": "2014-09-26T13:44:05.699520-05:00",
      "id": "ahBzfmNoaXJwcmFkaW8taHJkchoLEg1QbGF5bGlzdEV2ZW50GICAgNKeq6QKDA",
      "lastfm_urls": {
        "med_image": "http:\/\/userserve-ak.last.fm\/serve\/64s\/100756499.jpg",
        "sm_image": "http:\/\/userserve-ak.last.fm\/serve\/34s\/100756499.jpg",
        "_processed": true,
        "large_image": "http:\/\/userserve-ak.last.fm\/serve\/174s\/100756499.jpg"
      }
    },
    {
      "played_at_local_ts": 1411738494,
      "dj": "DJ Dead Alive - Ripped Sounds",
      "artist": "Destiny's Child",
      "track": "Say My Name",
      "notes": "",
      "artist_is_local": false,
      "label": "Sony",
      "played_at_gmt_ts": 1411756494,
      "played_at_gmt": "2014-09-26T18:34:54.598910",
      "release": "The Writing's On The Wall",
      "played_at_local_expire": "2015-03-31T13:34:54.598910-05:00",
      "played_at_local": "2014-09-26T13:34:54.598910-05:00",
      "id": "ahBzfmNoaXJwcmFkaW8taHJkchoLEg1QbGF5bGlzdEV2ZW50GICAgNLUubkKDA",
      "lastfm_urls": {
        "med_image": "http:\/\/userserve-ak.last.fm\/serve\/64s\/87900563.png",
        "sm_image": "http:\/\/userserve-ak.last.fm\/serve\/34s\/87900563.png",
        "_processed": true,
        "large_image": "http:\/\/userserve-ak.last.fm\/serve\/174s\/87900563.png"
      }
    },
    {
      "played_at_local_ts": 1411738252,
      "dj": "DJ Dead Alive - Ripped Sounds",
      "artist": "Daft Punk",
      "track": "Prime Time of Your Life",
      "notes": "",
      "artist_is_local": false,
      "label": "Virgin",
      "played_at_gmt_ts": 1411756252,
      "played_at_gmt": "2014-09-26T18:30:52.210280",
      "release": "Human After All",
      "played_at_local_expire": "2015-03-31T13:30:52.210280-05:00",
      "played_at_local": "2014-09-26T13:30:52.210280-05:00",
      "id": "ahBzfmNoaXJwcmFkaW8taHJkchoLEg1QbGF5bGlzdEV2ZW50GICAgNLCn6UKDA",
      "lastfm_urls": {
        "med_image": "http:\/\/userserve-ak.last.fm\/serve\/64s\/87950735.png",
        "sm_image": "http:\/\/userserve-ak.last.fm\/serve\/34s\/87950735.png",
        "_processed": true,
        "large_image": "http:\/\/userserve-ak.last.fm\/serve\/174s\/87950735.png"
      }
    },
    {
      "played_at_local_ts": 1411737910,
      "dj": "DJ Dead Alive - Ripped Sounds",
      "artist": "Simian Mobile Disco",
      "track": "Hypnick Jerk",
      "notes": "",
      "artist_is_local": false,
      "label": "Anti-",
      "played_at_gmt_ts": 1411755910,
      "played_at_gmt": "2014-09-26T18:25:10.506900",
      "release": "Whorl",
      "played_at_local_expire": "2015-03-31T13:25:10.506900-05:00",
      "played_at_local": "2014-09-26T13:25:10.506900-05:00",
      "id": "ahBzfmNoaXJwcmFkaW8taHJkchoLEg1QbGF5bGlzdEV2ZW50GICAgNK64LkJDA",
      "lastfm_urls": {
        "med_image": "http:\/\/userserve-ak.last.fm\/serve\/64s\/99933203.jpg",
        "sm_image": "http:\/\/userserve-ak.last.fm\/serve\/34s\/99933203.jpg",
        "_processed": true,
        "large_image": "http:\/\/userserve-ak.last.fm\/serve\/174s\/99933203.jpg"
      }
    },
    {
      "played_at_local_ts": 1411737750,
      "dj": "DJ Dead Alive - Ripped Sounds",
      "artist": "Ty Segall",
      "track": "The Clock",
      "notes": "",
      "artist_is_local": false,
      "label": "Drag City",
      "played_at_gmt_ts": 1411755750,
      "played_at_gmt": "2014-09-26T18:22:30.043820",
      "release": "Manipulator",
      "played_at_local_expire": "2015-03-31T13:22:30.043820-05:00",
      "played_at_local": "2014-09-26T13:22:30.043820-05:00",
      "id": "ahBzfmNoaXJwcmFkaW8taHJkchoLEg1QbGF5bGlzdEV2ZW50GICAgNKetbkKDA",
      "lastfm_urls": {
        "med_image": "http:\/\/userserve-ak.last.fm\/serve\/64s\/100939681.png",
        "sm_image": "http:\/\/userserve-ak.last.fm\/serve\/34s\/100939681.png",
        "_processed": true,
        "large_image": "http:\/\/userserve-ak.last.fm\/serve\/174s\/100939681.png"
      }
    }
  ]
}`)
	exampleTrack := getTrackInfoFromJson(someJsonBytes, false)
	if exampleTrack.Artist != "Homeboy Sandman" {
		t.Errorf("Expected \"Homeboy Sandman\", got %s", exampleTrack.Artist)
	}
}

func Test_makeRDSMessage(t *testing.T) {
	type args struct {
		currentTrack track
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "Basic test",
		args: args{currentTrack:track{Artist: "Regina Spektor", Track:"Chemo Limo"}},
		want: "DPS='Chemo Limo' by Regina Spektor on CHIRP Radio\n"},
		{name: "Long message test",
		args: args{currentTrack:track{Artist:"Sufjean Stephens", Track: "Concerning the UFO sighting on blah blah blah blah blah blah blah blah blah blah blah blah blah blah"}},
		want: "DPS='Concerning the UFO sighting on blah blah blah blah blah blah blah blah blah blah blah blah blah blah' by Sufjean Stephens\n"},
		{name: "Extra long message test",
			args: args{currentTrack:track{Artist:"Sufjean Stephens", Track: "Concerning the UFO sighting on blah blah blah blah blah blah blah blah blah blah blah blah blah blah blah blah"}},
			want: "DPS='Concerning the UFO sighting on blah blah blah blah blah blah blah blah blah blah blah blah blah blah blah blah' by Sufjean ...\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := makeRDSMessage(tt.args.currentTrack)
			if got != tt.want {
				t.Errorf("makeRDSMessage() = %v, want %v", got, tt.want)
			}
			if len(got)> 132 {
				t.Errorf("makeRDSMessage() length = %v, want <=132", len(got))
			}
		})
	}
}
