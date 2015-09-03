// This file defines the static data types exported by package
// musicmanager.

package musicmanager

import (
	"fmt"

	mmdspb "github.com/lxr/go.google.musicmanager/internal/download_proto/service"
	mmssjs "github.com/lxr/go.google.musicmanager/internal/session_json"
	mmuspb "github.com/lxr/go.google.musicmanager/internal/upload_proto/service"
)

// A RegisterError is returned by Client.Register if the server refuses
// to register the client for some reason.
type RegisterError mmuspb.UploadResponse_AuthStatus

func (e RegisterError) Error() string {
	return fmt.Sprint("musicmanager register error: ", mmuspb.UploadResponse_AuthStatus(e))
}

// A ListError is returned by Client.ListTracks if the server refuses
// to list the tracks for some reason.
type ListError mmdspb.GetTracksToExportResponse_TracksToExportStatus

func (e ListError) Error() string {
	return fmt.Sprint("musicmanager list error: ", mmdspb.GetTracksToExportResponse_TracksToExportStatus(e))
}

// An ImportError is returned by Client.ImportTracks if the server
// rejects a track based on its metadata or audio sample.
type ImportError mmuspb.TrackSampleResponse_ResponseCode

func (e ImportError) Error() string {
	return fmt.Sprint("musicmanager import error: ", mmuspb.TrackSampleResponse_ResponseCode(e))
}

// A RequestError is returned by all Client methods if an HTTP request
// is responded to with a non-2xx status code.
type RequestError mmssjs.SessionError

func (e *RequestError) Error() string {
	return fmt.Sprintf("<!-- server responded with status code %d and the following body: -->\n%s", e.Code, e.Message)
}

// TrackChannels represents the number of channels a Track can have.
type TrackChannels int

const (
	Mono   TrackChannels = 1
	Stereo TrackChannels = 2
)

// TrackRating represents the rating of a track.
type TrackRating int

const (
	NoRating TrackRating = 1 + iota

	OneStar // thumbs down
	TwoStars
	ThreeStars
	FourStars
	FiveStars // thumbs up
)

// TrackType defines the origin of a track.
type TrackType int

const (
	Matched TrackType = 1 + iota
	Unmatched
	Local
	Purchased
	MetadataOnlyMatched
	Promotional
)

// A Track represents metadata about a track.  When in a TrackList,
// only a subset of the fields are populated.
type Track struct {
	// There fields are present inside a TrackList.
	Id          string
	Title       string
	Artist      string
	Album       string
	AlbumArtist string
	TrackNumber int
	TrackSize   int64

	// Additional fields that can be given on import.
	ClientId        string
	Composer        string
	Genre           string
	Comment         string
	Year            int
	TotalTrackCount int
	DiscNumber      int
	TotalDiscCount  int
	PlayCount       int
	BeatsPerMinute  int
	Channels        TrackChannels
	Rating          TrackRating
	TrackType       TrackType
	// BUG(lor): Album art cannot be uploaded.

	// BitRate is the bitrate of the track in kbps, or 0 if don't
	// care.
	BitRate int

	// SampleFunc can be optionally used to provide the server with
	// a 128kbps MP3 sample of the track if requested.  It takes the
	// start and length of the desired sample in milliseconds.  If
	// SampleFunc is nil, an empty sample is sent.
	SampleFunc func(start, duration int) []byte
}

// A TrackList is one page of a track listing.
type TrackList struct {
	// The actual page of tracks.
	Items []*Track `convert:"/DownloadTrackInfo"`

	// Page token for the next page of tracks.
	PageToken string `convert:"/ContinuationToken"`

	// The last time one of the tracks in the list was modified,
	// expressed as a Unix timestamp in microseconds.
	UpdatedMin int64 `convert:"/UpdatedMin"`

	// Whether this listing contains only purchased or promotional
	// tracks.
	PurchasedOnly bool `convert:"-"`
}
