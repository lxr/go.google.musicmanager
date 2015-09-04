/*

Listing tracks

Usage:

	gmusic list [-f format] [-p] [-t date]

List lists the tracks in the user's Google Play Music library in the
following format:

	6eaf9df8-8a2c-3845-a443-80a500e07cbd	桜花爛漫
	448aa24d-18b7-38e0-a49a-3ab11db6f012	Hotel Nichifornia
	64da27ed-721d-33ee-b70e-81cc84546590	Hold Your Sexy Arms Against Me
	3a80d05a-58a1-3684-a154-fc046fcc4e6c	College Is Crazy
	...

The -f flag can be used to specify an alternative format for the list,
using the syntax of text/template.  The default format is
"{{.Id}}\t{{.Title}}\n".  The struct passed to the template is:

	type Track struct {
		Id          string
		Title       string
		Artist      string
		Album       string
		AlbumArtist string
		TrackNumber int
		TrackSize   int64

		// plus other, always-zero fields
	}

The -p flag causes the listing to be limited to purchased or promotional
tracks only.

The -t flag can be used to restrict the listing to tracks updated after
the specified date, given as an RFC 3339 timestamp (up to microsecond
precision).  The default is "1970-01-01T00:00:00Z".

*/
package main

import (
	"flag"
	"os"
	"text/template"
	"time"
)

var listTpls = template.Must(template.New("tracklist").
	Parse(`{{range .Items}}{{template "track" .}}{{end}}`))

func init() {
	cmds["list"] = list
}

func list() error {
	var (
		formatStr = flag.String("f", "{{.Id}}\t{{.Title}}\n",
			"alternative format for track listing")
		purchasedOnly = flag.Bool("p", false,
			"list only purchased or promotional tracks")
		updatedMinStr = flag.String("t", "1970-01-01T00:00:00Z",
			"only list tracks modified after the given date (RFC 3339 format)")
	)

	flag.Parse()
	t, err := time.Parse(time.RFC3339Nano, *updatedMinStr)
	if err != nil {
		return err
	}
	updatedMin := t.UnixNano() / 1000
	listTpls, err = listTpls.New("track").Parse(*formatStr)
	if err != nil {
		return err
	}
	pageToken := new(string)

	client, err := loadClient()
	if err != nil {
		return err
	}
	for pageToken != nil {
		list, err := client.ListTracks(*purchasedOnly, updatedMin, *pageToken)
		if err != nil {
			return err
		}
		err = listTpls.ExecuteTemplate(os.Stdout, "tracklist", list)
		if err != nil {
			return err
		}
		*pageToken = list.PageToken
		if *pageToken == "" {
			pageToken = nil
		}
	}
	return nil
}
