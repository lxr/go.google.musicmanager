// Package musicmanager implements a client for managing Google Play
// Music libraries.
package musicmanager

import (
	"fmt"
	"net/http"

	convert "github.com/lxr/go.google.musicmanager/internal/convert"
	mmdspb "github.com/lxr/go.google.musicmanager/internal/download_proto/service"
	mmldpb "github.com/lxr/go.google.musicmanager/internal/locker_proto/data"
	mmssjs "github.com/lxr/go.google.musicmanager/internal/session_json"
	mmudpb "github.com/lxr/go.google.musicmanager/internal/upload_proto/data"
	mmuspb "github.com/lxr/go.google.musicmanager/internal/upload_proto/service"
)

// Scope is the OAuth scope used by this API.
const Scope = "https://www.googleapis.com/auth/musicmanager"

// Client is a Music Manager client.
type Client struct {
	client *http.Client
	id     string
}

// NewClient creates a new Music Manager client with the given device ID
// and underlying HTTP client.  The supplied client must be capable of
// authenticating requests with the OAuth scope Scope.
func NewClient(client *http.Client, deviceID string) (*Client, error) {
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}
	if deviceID == "" {
		return nil, fmt.Errorf("device ID is empty")
	}
	return &Client{client, deviceID}, nil
}

// Register registers the client as a device with the given name under
// the user's Play Music library.  A client must be registered before
// the other method calls can succeed.  Re-registering a client can be
// used to change its name.  Note that there are limits to how many
// devices one account can have registered at a time, with how many
// accounts one device can be registered, and how many devices a user
// can deregister in a year, so be judicious in using this call.
func (c *Client) Register(name string) error {
	res, err := c.upAuth(&mmuspb.UpAuthRequest{
		UploaderId:   c.id,
		FriendlyName: name,
	})
	if err != nil {
		return err
	}
	if res != mmuspb.UploadResponse_OK {
		return RegisterError(res)
	}
	return nil
}

// ExportTrack returns a short-lived download URL for the given track,
// identified by its server ID.  Downloading the track from this URL
// requires no authentication.
func (c *Client) ExportTrack(id string) (string, error) {
	res, err := c.getDownloadSession(&mmssjs.GetDownloadSessionRequest{
		XDeviceID: c.id,
		SongID:    id,
	})
	if err != nil {
		return "", err
	}
	return res.URL, nil
}

// ListTracks lists the user's tracks from least to most recently
// accessed.  It can optionally filter for purchased or promotional
// tracks only, or for tracks last modified after a given Unix timestamp
// (microsecond precision).  Long responses may be returned in chunks,
// in which case the PageToken field of the TrackList object should
// be given to a new ListTracks call.
func (c *Client) ListTracks(purchasedOnly bool, updatedMin int64, pageToken string) (*TrackList, error) {
	var exportType mmdspb.GetTracksToExportRequest_TracksToExportType
	switch purchasedOnly {
	case true:
		exportType = mmdspb.GetTracksToExportRequest_PURCHASED_AND_PROMOTIONAL
	case false:
		exportType = mmdspb.GetTracksToExportRequest_ALL
	}
	res, err := c.getTracksToExport(&mmdspb.GetTracksToExportRequest{
		ClientId:          c.id,
		ExportType:        exportType,
		UpdatedMin:        updatedMin,
		ContinuationToken: pageToken,
	})
	switch err := err.(type) {
	case nil:
		if res.Status != mmdspb.GetTracksToExportResponse_OK {
			return nil, ListError(res.Status)
		}
	case *RequestError:
		// The Google Play servers respond with 304 Not Modified
		// if no tracks have been modified after the updatedMin
		// timestamp.  This is not exactly an error condition,
		// so we break out of the switch rather than return the
		// error; the remaining function body then arranges to
		// return an empty TrackList.
		if err.Code == http.StatusNotModified {
			break
		}
		return nil, err
	default:
		return nil, err
	}
	trackList := new(TrackList)
	convert.Convert(trackList, res)
	trackList.PurchasedOnly = purchasedOnly
	return trackList, nil
}

// ImportTracks returns short-lived upload URLs for the given tracks.
// MP3 audio data can be PUT to these URLs without authentication.
// Individual tracks can fail, in which case errs[i] contains the
// reason why importing tracks[i] failed.
//
// The Title field of a track to be imported cannot be empty.  The
// server also uses the ClientId field to identify which tracks have
// already been uploaded, so trying to import a track with a client ID
// that has already been uploaded will fail.  Leaving the ClientId field
// empty appears to bypass this server-side check; however, the
// implementation of ImportTracks cannot handle two tracks in a batch
// with the same client ID, even if they are both empty.
//
// A response to a request sent to a URL returned from this function can
// report failure through means other than the status code.  Use
// CheckImportResponse to verify that the request succeeded.
func (c *Client) ImportTracks(tracks []*Track) (urls []string, errs []error) {
	// Construct and the client-ID-to-track-index mapping and the
	// initial metadata upload.
	cidm := make(map[string]int)
	trks := make([]*mmldpb.Track, len(tracks))
	errs = make([]error, len(tracks))
	for i, track := range tracks {
		if _, ok := cidm[track.ClientId]; ok {
			errs[i] = fmt.Errorf("trying to import two tracks with the same client-side ID")
			continue
		}
		trks[i] = new(mmldpb.Track)
		convert.Convert(trks[i], track)
		cidm[track.ClientId] = i
	}
	// Upload track metadata.
	res, err := c.uploadMetadata(&mmuspb.UploadMetadataRequest{
		UploaderId: c.id,
		Track:      trks,
	})
	if err != nil {
		for i := range errs {
			if errs[i] == nil {
				errs[i] = err
			}
		}
		return nil, errs
	}
	// Satisfy any requests for track samples and append the
	// responses to tres.
	tres := res.TrackSampleResponse
	if n := len(res.SignedChallengeInfo); n > 0 {
		spls := make([]*mmudpb.TrackSample, n)
		for i, sci := range res.SignedChallengeInfo {
			ci := sci.ChallengeInfo
			j := cidm[ci.ClientTrackId]
			var sample []byte
			if sf := tracks[j].SampleFunc; sf != nil {
				sample = sf(
					int(ci.StartMillis),
					int(ci.DurationMillis),
				)
			}
			if sample == nil {
				// A nil sample is different from an
				// empty one: the former results in an
				// invalid sample message.  So, if
				// Sampler leaves us with a nil sample,
				// replace it with an empty one.
				sample = make([]byte, 0)
			}
			spls[i] = &mmudpb.TrackSample{
				Track:               trks[j],
				Sample:              sample,
				SignedChallengeInfo: sci,
				SampleFormat:        mmldpb.Track_MP3,
			}
		}
		sres, err := c.uploadSample(&mmuspb.UploadSampleRequest{
			UploaderId:  c.id,
			TrackSample: spls,
		})
		if err != nil {
			for _, spl := range spls {
				errs[cidm[spl.Track.ClientId]] = err
			}
		} else {
			tres = append(tres, sres.TrackSampleResponse...)
		}
	}
	// Parse responses to track metadata and samples.  The result
	// is a map from track indices to their server IDs.
	sidm := make(map[int]string)
	for _, res := range tres {
		i := cidm[res.ClientTrackId]
		if res.ResponseCode != mmuspb.TrackSampleResponse_UPLOAD_REQUESTED {
			errs[i] = ImportError(res.ResponseCode)
			continue
		}
		sidm[i] = res.ServerTrackId
	}
	// Acquire upload sessions.
	urls = make([]string, len(tracks))
	for i, id := range sidm {
		trk := tracks[i]
		res, err := c.getUploadSession(&mmssjs.GetUploadSessionRequest{
			Name:         id,
			UploaderId:   c.id,
			ClientId:     trk.ClientId,
			ServerId:     id,
			TrackBitRate: trk.BitRate,
			// BUG(lor): Client.ImportTracks does not
			// activate the upload progress tracker on
			// https://play.google.com/music.
		})
		if err != nil {
			errs[i] = err
			continue
		}
		if res.Error != nil {
			errs[i] = res.Error
			continue
		}
		urls[i] = res.Transfers[0].PutUrl
	}
	return urls, errs
}

// CheckImportResponse checks the response to an HTTP request sent to a
// URL returned by Client.ImportTracks.  It returns the server ID of the
// track to which resp is a response, or an error if the response is
// erroneous in some way.  The error will be of type *RequestError if
// the response has a non-2xx status code.  CheckImportResponse closes
// resp.Body.
func CheckImportResponse(resp *http.Response) (string, error) {
	res := new(mmssjs.GetUploadSessionResponse)
	err := parseResponse(resp, res)
	switch {
	case err != nil:
		return "", err
	case res.Error != nil:
		return "", res.Error
	default:
		return res.Transfers[0].Name, nil
	}
}
