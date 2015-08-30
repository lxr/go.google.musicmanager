// Package musicmanager implements a client for managing Google Play
// Music libraries.
package musicmanager

import (
	"fmt"
	"net/http"

	convert "my-git.appspot.com/go.google.musicmanager/internal/convert"
	mmdspb "my-git.appspot.com/go.google.musicmanager/internal/download_proto/service"
	mmldpb "my-git.appspot.com/go.google.musicmanager/internal/locker_proto/data"
	mmssjs "my-git.appspot.com/go.google.musicmanager/internal/session_json"
	mmudpb "my-git.appspot.com/go.google.musicmanager/internal/upload_proto/data"
	mmuspb "my-git.appspot.com/go.google.musicmanager/internal/upload_proto/service"
)

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
		return "", nil
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
	if err != nil {
		return nil, err
	}
	if res.Status != mmdspb.GetTracksToExportResponse_OK {
		return nil, ListError(res.Status)
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
			sample := []byte{}
			if sampler := tracks[j].Sampler; sampler != nil {
				sample = sampler(
					int(ci.StartMillis),
					int(ci.DurationMillis),
				)
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
	uploaded := 0
	total := len(sidm)
	urls = make([]string, len(tracks))
	for i, id := range sidm {
		trk := tracks[i]
		res, err := c.getUploadSession(&mmssjs.GetUploadSessionRequest{
			Name:                      id,
			UploaderId:                c.id,
			ClientId:                  trk.ClientId,
			ServerId:                  id,
			TrackBitRate:              trk.BitRate,
			CurrentUploadingTrack:     trk.Title,
			CurrentTotalUploadedCount: uploaded,
			ClientTotalSongCount:      total,
			SyncNow:                   true,
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
