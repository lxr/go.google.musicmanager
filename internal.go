// This file implements the low-level service calls of the Google
// Music Manager service.

package musicmanager

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/golang/protobuf/proto"

	mmdspb "my-git.appspot.com/go.google.musicmanager/internal/download_proto/service"
	mmssjs "my-git.appspot.com/go.google.musicmanager/internal/session_json"
	mmuspb "my-git.appspot.com/go.google.musicmanager/internal/upload_proto/service"
)

func (c *Client) upAuth(req *mmuspb.UpAuthRequest) (mmuspb.UploadResponse_AuthStatus, error) {
	res, err := c.uploadServiceCall("upauth", req)
	if err != nil {
		return mmuspb.UploadResponse_UNKNOWN, err
	}
	return res.AuthStatus, nil
}

func (c *Client) uploadMetadata(req *mmuspb.UploadMetadataRequest) (*mmuspb.UploadMetadataResponse, error) {
	res, err := c.uploadServiceCall("metadata?version=1", req)
	if err != nil {
		return nil, err
	}
	return res.MetadataResponse, nil
}

func (c *Client) uploadSample(req *mmuspb.UploadSampleRequest) (*mmuspb.UploadSampleResponse, error) {
	res, err := c.uploadServiceCall("sample?version=1", req)
	if err != nil {
		return nil, err
	}
	return res.SampleResponse, nil
}

func (c *Client) clientState(req *mmuspb.ClientStateRequest) (*mmuspb.ClientStateResponse, error) {
	res, err := c.uploadServiceCall("clientstate", req)
	if err != nil {
		return nil, err
	}
	return res.ClientstateResponse, nil
}

func (c *Client) getJobs(req *mmuspb.GetJobsRequest) (*mmuspb.GetJobsResponse, error) {
	res, err := c.uploadServiceCall("getjobs", req)
	if err != nil {
		return nil, err
	}
	return res.GetjobsResponse, nil
}

func (c *Client) updateUploadState(req *mmuspb.UpdateUploadStateRequest) error {
	_, err := c.uploadServiceCall("uploadstate", req)
	return err
}

func (c *Client) deleteUploadRequested(req *mmuspb.DeleteUploadRequestedRequest) error {
	_, err := c.uploadServiceCall("deleteuploadrequested", req)
	return err
}

func (c *Client) getTracksToExport(req *mmdspb.GetTracksToExportRequest) (*mmdspb.GetTracksToExportResponse, error) {
	res := new(mmdspb.GetTracksToExportResponse)
	url := "https://music.google.com/music/exportids"
	return res, c.post(url, req, res)
}

func (c *Client) getDownloadSession(req *mmssjs.GetDownloadSessionRequest) (*mmssjs.GetDownloadSessionResponse, error) {
	urlStr := "https://music.google.com/music/export?" + url.Values{
		"version": {"2"},
		"songid":  {req.SongID},
	}.Encode()
	reqp, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	reqp.Header.Set("X-Device-ID", req.XDeviceID)
	res := new(mmssjs.GetDownloadSessionResponse)
	return res, c.do(reqp, res)
}

func (c *Client) getUploadSession(req *mmssjs.GetUploadSessionRequest) (*mmssjs.GetUploadSessionResponse, error) {
	res := new(mmssjs.GetUploadSessionResponse)
	url := "https://uploadsj.clients.google.com/uploadsj/rupio"
	return res, c.post(url, req, res)
}

// uploadServiceCall protobuf-encodes the request and POSTs it to the
// named endpoint under https://android.clients.google.com/upsj/,
// decoding the response as a *pb.UploadResponse.
func (c *Client) uploadServiceCall(endpoint string, req interface{}) (*mmuspb.UploadResponse, error) {
	const baseURL = "https://android.clients.google.com/upsj/"
	res := new(mmuspb.UploadResponse)
	return res, c.post(baseURL+endpoint, req, res)
}

// post encodes the request object as protobuf if it implements
// proto.Message, and as JSON otherwise, and POSTs the result to the
// given URL.  The response is then similarly decoded.
func (c *Client) post(url string, req, res interface{}) error {
	var body io.Reader
	var buf []byte
	var err error
	switch v := req.(type) {
	case proto.Message:
		buf, err = proto.Marshal(v)
	default:
		buf, err = json.Marshal(v)
	}
	if err != nil {
		return err
	}
	if buf != nil {
		body = bytes.NewReader(buf)
	}
	reqp, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	return c.do(reqp, res)
}

// do executes the given *http.Request and decodes its response to
// res as protobuf if res implements proto.Message, and as JSON
// otherwise.
func (c *Client) do(req *http.Request, res interface{}) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		res = resp.StatusCode
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	switch v := res.(type) {
	case int:
		return &RequestError{
			Code:    v,
			Message: string(buf),
		}
	case proto.Message:
		return proto.Unmarshal(buf, v)
	default:
		return json.Unmarshal(buf, v)
	}
}
