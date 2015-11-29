// Package json describes the interface of the JSON calls of Google's
// Music Manager service.
package google_musicmanager_v0

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/lxr/go.google.musicmanager/internal/convert"
)

type GetDownloadSessionRequest struct {
	XDeviceID string `header:"X-Device-ID"`
	SongID    string `query-parameter:"songid"`
}

type GetDownloadSessionResponse struct {
	URL string `json:"url"`
}

type GetUploadSessionRequest struct {
	Name                      string `convert:"/external/name"`
	UploaderId                string `convert:"/inlined/UploaderId"`
	ClientId                  string `convert:"/inlined/ClientId"`
	ServerId                  string `convert:"/inlined/ServerId"`
	TrackBitRate              int    `convert:"/inlined/TrackBitRate"`
	CurrentUploadingTrack     string `convert:"/inlined/CurrentUploadingTrack"`
	CurrentTotalUploadedCount int    `convert:"/inlined/CurrentTotalUploadedCount"`
	ClientTotalSongCount      int    `convert:"/inlined/ClientTotalSongCount"`
	SyncNow                   bool   `convert:"/inlined/SyncNow"`
	TrackDoNotRematch         bool   `convert:"/inlined/TrackDoNotRematch"`
}

func (r *GetUploadSessionRequest) MarshalJSON() ([]byte, error) {
	var msg interface{}
	convert.Convert(&msg, r)
	buf := new(bytes.Buffer)
	if err := getUploadSessionRequestTemplate.Execute(buf, msg); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

var getUploadSessionRequestTemplate = template.Must(
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

type GetUploadSessionResponse struct {
	UploadId string `convert:"/sessionStatus/upload_id"`
	State    string `convert:"/sessionStatus/state"`

	Error *SessionError `convert:"/errorMessage/additionalInfo/uploader_service.GoogleRupioAdditionalInfo"`

	Status *struct {
		Code                int    `convert:"/customerSpecificInfo/ResponseCode"`
		Message             string `convert:"/status"`
		ServerFileReference string `convert:"/customerSpecificInfo/ServerFileReference"`
	} `convert:"/sessionStatus/additionalInfo/uploader_service.GoogleRupioAdditionalInfo/completionInfo"`

	Transfers []struct {
		Name             string `convert:"/name"`
		Status           string `convert:"/status"`
		PutUrl           string `convert:"/putInfo/url"`
		BytesTransferred int64  `convert:"/bytesTransferred"`
		BytesTotal       int64  `convert:"/bytesTotal"`
	} `convert:"/sessionStatus/externalFieldTransfers"`
}

func (r *GetUploadSessionResponse) UnmarshalJSON(buf []byte) error {
	var msg interface{}
	if err := json.Unmarshal(buf, &msg); err != nil {
		return err
	}
	convert.Convert(r, msg)
	return nil
}

type SessionError struct {
	Code    int    `convert:"/completionInfo/customerSpecificInfo/ResponseCode"`
	Message string `convert:"/requestRejectedInfo/reasonDescription"`
}

func (e *SessionError) Error() string {
	return fmt.Sprintf("session error %d: %s", e.Code, e.Message)
}
