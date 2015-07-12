// Package client provides a low-level interface to Google's Music
// Manager service.  For an overview of how the Music Manager service
// works, see google-musicmanager-go/json and
// google-musicmanager-go/proto.
//
// Example use:
// 	import (
// 		"github.com/golang/protobuf/proto"
// 		"google-musicmanager-go/client"
// 		pb "google-musicmanager-go/proto"
// 	)
// 	...
// 	mm := client.New(authorizedHTTPClient)
// 	req := &pb.GetTracksToExportRequest{
// 		ClientId: proto.String("myid"),
// 	}
// 	res, err := mm.GetTracksToExport(req)
// 	for i, track := range res.GetDownloadTrackInfo() {
// 		...
package client

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	"google.golang.org/api/googleapi"

	js "google-musicmanager-go/json"
	pb "google-musicmanager-go/proto"
)

// A File wraps a normal io.ReadCloser with a few metadata interfaces.
// It is returned by Client.Get.
type File struct {
	io.ReadCloser
	name string
	size int64
}

// Name returns the name of the file, or the empty string if unknown.
func (f File) Name() string {
	return f.name
}

// Size returns the size of the file, or zero if unknown.
func (f File) Size() int64 {
	return f.size
}

// A Client is a low-level Music Manager client.  It implements
// pb.UploadService, pb.DownloadService, js.SessionService and
// a few convenience functions for uploading and downloading songs.
type Client struct {
	client *http.Client
}

// New creates a new Music Manager client from the given *http.Client.
// The client should already carry a Google OAuth 2.0 token with a
// scope of https://www.googleapis.com/auth/musicmanager.
func New(client *http.Client) *Client {
	return &Client{client}
}

// See pb.UploadService.
func (c *Client) UpAuth(req *pb.UpAuthRequest) (*pb.UploadResponse, error) {
	return c.uploadServiceCall("upauth", req)
}

// See pb.UploadService.
func (c *Client) ClientState(req *pb.ClientStateRequest) (*pb.UploadResponse, error) {
	return c.uploadServiceCall("clientstate", req)
}

// See pb.UploadService.
func (c *Client) UpdateUploadState(req *pb.UpdateUploadStateRequest) (*pb.UploadResponse, error) {
	return c.uploadServiceCall("uploadstate", req)
}

// See pb.UploadService.
func (c *Client) GetJobs(req *pb.GetJobsRequest) (*pb.UploadResponse, error) {
	return c.uploadServiceCall("getjobs", req)
}

// See pb.UploadService.
func (c *Client) DeleteUploadRequested(req *pb.DeleteUploadRequestedRequest) (*pb.UploadResponse, error) {
	return c.uploadServiceCall("deleteuploadrequested", req)
}

// See pb.UploadService.
func (c *Client) UploadMetadata(req *pb.UploadMetadataRequest) (*pb.UploadResponse, error) {
	return c.uploadServiceCall("metadata?version=1", req)
}

// See pb.UploadService.
func (c *Client) UploadSample(req *pb.UploadSampleRequest) (*pb.UploadResponse, error) {
	return c.uploadServiceCall("sample?version=1", req)
}

// Not implemented.  Panics if called.  See pb.UploadService.
func (c *Client) UploadPlaylist(req *pb.UploadPlaylistRequest) (*pb.UploadResponse, error) {
	panic("client: not implemented")
}

// Not implemented.  Panics if called.  See pb.UploadService.
func (c *Client) UploadPlaylistEntry(req *pb.UploadPlaylistEntryRequest) (*pb.UploadResponse, error) {
	panic("client: not implemented")
}

// uploadServiceCall protobuf-encodes the request and POSTs it to the
// named endpoint under https://android.clients.google.com/upsj/,
// decoding the response as a *pb.UploadResponse.
func (c *Client) uploadServiceCall(endpoint string, req interface{}) (*pb.UploadResponse, error) {
	const baseURL = "https://android.clients.google.com/upsj/"
	var res pb.UploadResponse
	if err := c.post(baseURL+endpoint, req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// See pb.DownloadService.
func (c *Client) GetTracksToExport(req *pb.GetTracksToExportRequest) (*pb.GetTracksToExportResponse, error) {
	var res pb.GetTracksToExportResponse
	url := "https://music.google.com/music/exportids"
	if err := c.post(url, req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// See js.SessionService.
func (c *Client) GetDownloadSession(req *js.GetDownloadSessionRequest) (*js.GetDownloadSessionResponse, error) {
	urlStr := "https://music.google.com/music/export?" + url.Values{
		"version": {"2"},
		"songid":  {req.SongID},
	}.Encode()
	reqp, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	reqp.Header.Set("X-Device-ID", req.XDeviceID)

	var ret js.GetDownloadSessionResponse
	if err := c.do(reqp, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}

// See js.SessionService.
func (c *Client) GetUploadSession(req *js.GetUploadSessionRequest) (*js.GetUploadSessionResponse, error) {
	var res js.GetUploadSessionResponse
	url := "https://uploadsj.clients.google.com/uploadsj/rupio"
	if err := c.post(url, req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// Get fetches a URL returned by GetDownloadSession and returns
// the response as a File.  Its Name and Size methods then return the
// values of the response's Content-Disposition and Content-Length
// headers respectively.
func (c *Client) Get(urlStr string) (*File, error) {
	res, err := c.client.Get(urlStr)
	if err != nil {
		return nil, err
	}
	if err := googleapi.CheckResponse(res); err != nil {
		googleapi.CloseBody(res)
		return nil, err
	}
	name := ""
	parts := strings.SplitN(res.Header.Get("Content-Disposition"), "filename*=UTF-8''", 2)
	if len(parts) == 2 {
		name, _ = url.QueryUnescape(parts[1])
	}
	size, _ := strconv.Atoi(res.Header.Get("Content-Length"))
	return &File{res.Body, name, int64(size)}, nil
}

// BUG(lor): On Google App Engine, the response headers don't include
// Content-Length, so the Size of the File returned by Get is always
// zero.  The reason is perhaps
// https://cloud.google.com/appengine/docs/go/urlfetch/#Go_Request_headers.

// Put PUTs data to a URL returned by GetUploadSession.  The return type
// is the same as for GetUploadSession.
func (c *Client) Put(urlStr string, body io.Reader) (*js.GetUploadSessionResponse, error) {
	req, err := http.NewRequest("PUT", urlStr, body)
	if err != nil {
		return nil, err
	}

	var res js.GetUploadSessionResponse
	if err := c.do(req, &res); err != nil {
		return nil, err
	}
	return &res, nil
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

// do executes the given *http.Request and interprets its response to
// res as protobuf if res implements proto.Message, and as JSON
// otherwise.
func (c *Client) do(req *http.Request, res interface{}) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer googleapi.CloseBody(resp)
	if err := googleapi.CheckResponse(resp); err != nil {
		return err
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	switch v := res.(type) {
	case proto.Message:
		err = proto.Unmarshal(buf, v)
	default:
		err = json.Unmarshal(buf, v)
	}
	return err
}
