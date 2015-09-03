/*

Downloading tracks

Usage:

	gmusic download [id...]

Download downloads the tracks identified by the IDs to the current
directory and prints the filenames they are saved under to standard
output.  Progress and error information is additionally written to
standard error.  If no IDs are given, a newline-separated list is read
from standard input.

If a track fails to download, download moves on to the next one.  The
exit status is 0 only if all tracks downloaded successfully.

The filename under which each track is saved is generated server-side
from the track's title and track number.  Download clobbers existing
files without asking, so be careful with it.

*/
package main

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/lxr/go.google.musicmanager"
)

func download() error {
	client, err := loadClient()
	if err != nil {
		return err
	}
	success := true
	scanner := getScanner()
	for scanner.Scan() {
		id := scanner.Text()
		logf("downloading %s: ", id)
		name, err := downloadTrack(client, id)
		if err != nil {
			logf("%v\n", err)
			success = false
		} else {
			println(name)
		}
	}
	switch {
	case scanner.Err() != nil:
		return scanner.Err()
	case !success:
		return errors.New("not all tracks were downloaded")
	default:
		return nil
	}
}

func downloadTrack(client *musicmanager.Client, id string) (string, error) {
	url, err := client.ExportTrack(id)
	if err != nil {
		return "", err
	}
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	size, err := strconv.ParseFloat(resp.Header.Get("Content-Length"), 64)
	if err == nil {
		logf("(%.2f MiB) ", size/(1<<20))
	}
	// Google probably won't put malicious paths in
	// Content-Disposition, but this will prevent the file from
	// being written outside the current directory.
	name, _ := getName(resp.Header.Get("Content-Disposition"))
	name = filepath.Base(name)
	if name == "." {
		name = id + ".mp3"
	}
	f, err := os.Create(name)
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return name, err
}

// getName extracts a UTF-8-encoded filename from a Content-Disposition
// header.  (The default mime.ParseMediaType function cannot be used,
// because it seems to disagree with the server on how to encode/decode
// parens.)
func getName(v string) (string, error) {
	parts := strings.SplitN(v, "filename*=UTF-8''", 2)
	if len(parts) < 2 {
		return "", errors.New("media type value lacks UTF-8 filename")
	}
	return url.QueryUnescape(parts[1])
}
