/*

Uploading tracks

Usage:

	gmusic upload [file...]

Upload uploads the named MP3 files to the user's Google Play Music
library and prints their server-side IDs to standard output.  Progress
and error information is additionally printed to standard error.  If
no filenames are given, upload reads a newline-separated list from
standard input.

If a file fails to upload, upload moves on to the next one.  The exit
status is 0 only if all tracks uploaded successfully.

*/
package main

import (
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/dhowden/tag"
	"github.com/lxr/go.google.musicmanager"
)

func upload() error {
	client, err := loadClient()
	if err != nil {
		return err
	}
	tracks := make([]*musicmanager.Track, 0)
	files := make([]*os.File, 0)
	success := true
	scanner := getScanner()
	for scanner.Scan() {
		name := scanner.Text()
		track, file, err := parseTrack(name)
		if err != nil {
			logf("reading %s: %v\n", name, err)
			success = false
			continue
		}
		defer file.Close()
		tracks = append(tracks, track)
		files = append(files, file)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	urls, errs := client.ImportTracks(tracks)
	for i, err := range errs {
		logf("uploading %s: ", files[i].Name())
		var id string
		if err == nil {
			id, err = uploadTrack(urls[i], files[i])
		}
		if err != nil {
			logf("%v\n", err)
			success = false
		} else {
			println(id)
		}
	}
	if !success {
		return errors.New("not all files were uploaded")
	}
	return nil
}

func uploadTrack(url string, r io.Reader) (string, error) {
	// BUG(lor): Gmusic does not actually detect and report an error
	// on (most) non-MP3 files.  All files are uploaded to Google
	// Play, but only MP3 ones will be playable.
	resp, err := http.Post(url, "audio/mpeg", r)
	if err != nil {
		return "", err
	}
	return musicmanager.CheckImportResponse(resp)
}

func parseTrack(name string) (track *musicmanager.Track, f *os.File, err error) {
	f, err = os.Open(name)
	if err != nil {
		return
	}
	sum, err := tag.Sum(f)
	if err = rewind(f, err); err != nil {
		return
	}
	metadata, err := tag.ReadFrom(f)
	if err = rewind(f, err); err == tag.ErrNoTagsFound {
		err = nil
		track = &musicmanager.Track{
			ClientId: sum,
			Title:    name,
		}
		return
	} else if err != nil {
		return
	}
	ti, tn := metadata.Track()
	di, dn := metadata.Disc()
	track = &musicmanager.Track{
		ClientId:        sum,
		Title:           metadata.Title(),
		Album:           metadata.Album(),
		Artist:          metadata.Artist(),
		AlbumArtist:     metadata.AlbumArtist(),
		Composer:        metadata.Composer(),
		Year:            metadata.Year(),
		Genre:           metadata.Genre(),
		TrackNumber:     ti,
		TotalTrackCount: tn,
		DiscNumber:      di,
		TotalDiscCount:  dn,
	}
	return
}

func rewind(s io.Seeker, err error) error {
	_, err1 := s.Seek(0, os.SEEK_SET)
	if err == nil {
		err = err1
	}
	return err
}
