// Package json describes the interface of the JSON calls of Google's
// Music Manager service.  Refer to the documentation for
// google-musicmanager-go/proto for an overview of the entire service.
package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"
)

// SessionService hands out upload and download sessions for track
// audio data.  A 'session' is a short-lived Google URL and its
// associated metadata.  MP3 byte streams can be PUT to or GET from
// these URLs (depending on the session type, naturally) without
// authentication.
type SessionService interface {

	// GetDownloadSession returns a download session for the given
	// track.  Such a download doesn't count towards the download
	// limit given in the web interface.
	//
	// Endpoint: https://music.google.com/music/export?version=2
	GetDownloadSession(*GetDownloadSessionRequest) (*GetDownloadSessionResponse, error)

	// GetUploadSession returns an upload session for the given
	// client ID-server ID pair.
	//
	// Endpoint: https://uploadsj.clients.google.com/uploadsj/rupio
	GetUploadSession(*GetUploadSessionRequest) (*GetUploadSessionResponse, error)
}

// Arguments to SessionService.GetDownloadSession.  These are not
// actually delivered JSON-encoded: XDeviceID is given as the
// X-Device-ID HTTP header, and SongID as the songid query parameter.
// The request method is GET.
//
// This call has been known to fail if XDeviceID is not sufficiently
// "MAC-like".  The exact threshold is unknown; perhaps the server
// only looks for a colon in the ID.
type GetDownloadSessionRequest struct {
	XDeviceID string // the device ID
	SongID    string // the server ID of the track to download
}

// Return values of SessionService.GetDownloadSession.
type GetDownloadSessionResponse struct {
	URL string `json:"url"` // download URL
}

// Arguments to SessionService.GetUploadSession.  This structure is not
// encoded as JSON directly, but wrapped as described in
// GetUploadSessionRequestTemplate.
type GetUploadSessionRequest struct {
	// Name of the upload session.  This is roundtripped back in the
	// response.  Can be an arbitrary string.
	Name string `mappath:"/external/name"`

	// The device ID.  Required.
	UploaderId string `mappath:"/inlined/UploaderId"`

	// The client ID of the track to upload.  Required.
	ClientId string `mappath:"/inlined/ClientId"`

	// The server ID of the track to upload.  Required.
	ServerId string `mappath:"/inlined/ServerId"`

	// The bitrate of the track to upload in kbps, or zero if not
	// known.  Omitting this field will cause the track not to play
	// in the Android app.  Can be zero if don't care.
	TrackBitRate int32 `mappath:"/inlined/TrackBitRate"`

	// A title for the track to be uploaded.  This used to displayed
	// in the web interface when the track was uploading, but this
	// no longer appears to be the case.  Omitting this field causes
	// the progress ticker in the web interface not to appear at all
	// for this upload.
	CurrentUploadingTrack string `mappath:"/inlined/CurrentUploadingTrack,omitempty"`

	// The number of tracks that have been uploaded before this in
	// the current batch.  Will be shown in the web interface.
	CurrentTotalUploadedCount int `mappath:"/inlined/CurrentTotalUploadedCount"`

	// The total number of tracks to upload in this batch.  Will
	// be shown in the web interface.
	ClientTotalSongCount int `mappath:"/inlined/ClientTotalSongCount"`

	// Whether to refresh the web interface when the upload starts.
	// This does not appear to work for the first song in a batch,
	// regardless of the value of CurrentTotalUploadedCount.
	SyncNow bool `mappath:"/inlined/SyncNow"`

	// This presumably controls whether the server attempts to
	// match the track again once it has its bytes, but it appears
	// to do nothing.
	TrackDoNotRematch bool `mappath:"/inlined/TrackDoNotRematch"`
}

func (r *GetUploadSessionRequest) MarshalJSON() ([]byte, error) {
	msg, err := mappathEncode(r)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	if err := GetUploadSessionRequestTemplate.Execute(buf, msg); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GetUploadSessionRequestTemplate describes how a
// GetUploadSessionRequest structure is marshaled as JSON.
// The structure is not directly fed to this template, but is first
// converted into a map according to the mappath tags in its definition.
// "quote" is a function that prints the default string representation
// of its argument double-quoted.
var GetUploadSessionRequestTemplate = template.Must(
	template.New("GetUploadSessionRequest").
		Funcs(map[string]interface{}{"quote": quote}).
		Parse(`{
	  "clientId": "Jumper Uploader",
	  "protocolVersion": "0.8",
	  "createSessionRequest": {
	    "fields": [
	      {"inlined": {
	        "name": "title",
	        "content": "jumper-uploader-title-42"
	      }},
	      {{range $key, $value := .inlined}}
	      {"inlined": {
	        "name": {{quote $key}},
	        "content": {{quote $value}}
	      }},
	      {{end}}
	      {"external": {
	        {{range $key, $value := .external}}
	          {{quote $key}}: {{quote $value}},
	        {{end}}
	        "put": {}
	      }}
	    ]
	  }
	}`))

func quote(x interface{}) string {
	return fmt.Sprintf("%q", fmt.Sprint(x))
}

// Return values of SessionService.GetUploadSession, as well as the
// response to PUTting data to TransferPutUrl.  The float64 fields
// contain only whole numbers; their type is a symptom of the
// limitations of json.Unmarshal.  The actual JSON response body is
// extremely deeply nested; the mappath tags here describe the "paths"
// inside the object hierarchy that lead to the actual values of
// interest.
type GetUploadSessionResponse struct {
	// Non-zero when the server cannot allocate an upload session.
	// This has the value of an HTTP status code, presumably from
	// Google's internal servers.
	ErrorCode float64 `mappath:"/errorMessage/additionalInfo/uploader_service.GoogleRupioAdditionalInfo/completionInfo/customerSpecificInfo/ResponseCode"`

	// The name of the upload session.
	TransferName string `mappath:"/sessionStatus/externalFieldTransfers/0/name"`

	// Status of the upload: one of IN_PROGRESS, COMPLETED.
	TransferStatus string `mappath:"/sessionStatus/externalFieldTransfers/0/status"`

	// The upload URL.
	TransferPutUrl string `mappath:"/sessionStatus/externalFieldTransfers/0/putInfo/url"`

	// Number of bytes that have been uploaded, and the total.
	// This would suggest that there is a way to do resumable
	// uploads to the server, but it is not known.  Notably,
	// TransferBytesTotal is only known, and thus present, in the
	// response for a successful upload.
	TransferBytesTransferred float64 `mappath:"/sessionStatus/externalFieldTransfers/0/bytesTransferred"`
	TransferBytesTotal       float64 `mappath:"/sessionStatus/externalFieldTransfers/0/bytesTotal"`

	// Undocumented fields.
	UploadId            string  `mappath:"/sessionStatus/upload_id"`
	State               string  `mappath:"/sessionStatus/state"`
	Status              string  `mappath:"/sessionStatus/additionalInfo/uploader_service.GoogleRupioAdditionalInfo/completionInfo/status"`
	ServerFileReference string  `mappath:"/sessionStatus/additionalInfo/uploader_service.GoogleRupioAdditionalInfo/completionInfo/customerSpecificInfo/ServerFileReference"`
	ResponseCode        float64 `mappath:"/sessionStatus/additionalInfo/uploader_service.GoogleRupioAdditionalInfo/completionInfo/customerSpecificInfo/ResponseCode"`
}

// TODO(lor): sessionStatus.externalFieldTransfers is an array, and
// though it probably only ever has one element (as I don't know how
// to request multiple sessions in GetUploadSessionRequest), it'd be
// nice to represent it as a slice in this struct as well.

func (r *GetUploadSessionResponse) UnmarshalJSON(buf []byte) error {
	var msg interface{}
	if err := json.Unmarshal(buf, &msg); err != nil {
		return err
	}
	mappathDecode(msg, r)
	return nil
}
