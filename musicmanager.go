// Package musicmanager provides access to the unofficial Music Manager
// API.  Its interface mimics that of the official Google API Go Client
// libraries (https://github.com/google/google-api-go-client), though
// internally it is quite different.
//
// This package is implemented atop the lower-level
// google-musicmanager-go/client package.  Efforts have been made to
// ensure that end users need not understand it or its dependencies to
// use this package, though this has not been reasonable in all cases.
// In particular, in order to avoid duplication of information, the
// protobuf messages returned or expected by certain methods in this
// package are documented only under google-musicmanager-go/proto.
//
// Usage example:
// 	import "google-musicmanager-go"
// 	...
// 	musicmanagerService, err := musicmanager.New(oauthHttpClient, "client id")
// 	err = musicmanagerService.Register("client name").Do()
package musicmanager

// BUG(lor): This package needs to make requests to the server
// android.clients.google.com, which can only be accessed over TLS using
// SNI (see https://tools.ietf.org/html/rfc3546#section-3.1).  This has
// been known to cause problems on Google App Engine.

import (
	"fmt"
	"io"
	"net/http"

	"github.com/golang/protobuf/proto"
	"google.golang.org/api/googleapi"

	mm "google-musicmanager-go/client"
	js "google-musicmanager-go/json"
	pb "google-musicmanager-go/proto"
)

// The OAuth2 scope used by this API.
const MusicManagerScope = "https://www.googleapis.com/auth/musicmanager"

// New creates a new Service from an *http.Client with the given ID.
// The client should already carry a Google OAuth 2.0 token with the
// MusicManagerScope.  The ID must also be registered with the user's
// Play Music account before the service calls can work; this can be
// done with the Register method.  The ID must be unique on Google's
// side, and should look like a MAC address (otherwise, the Get call
// won't work.)
//
// Returns an error if the client is nil or the id is empty.
func New(client *http.Client, id string) (*Service, error) {
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}
	if id == "" {
		return nil, fmt.Errorf("id is empty")
	}
	s := &Service{client: mm.New(client), id: id}
	s.Jobs = NewJobsService(s)
	s.Tracks = NewTracksService(s)
	return s, nil
}

type Service struct {
	client *mm.Client
	id     string

	Jobs *JobsService

	Tracks *TracksService
}

func NewTracksService(s *Service) *TracksService {
	return &TracksService{s}
}

type TracksService struct {
	s *Service
}

func NewJobsService(s *Service) *JobsService {
	return &JobsService{s}
}

type JobsService struct {
	s *Service
}

type TracksInsertBatchResponse []*TracksInsertBatchResponseEntry

type TracksInsertBatchResponseEntry struct {
	ServerID string
	Error    error
}

type RegisterCall struct {
	s *Service
	c *pb.UpAuthRequest
}

// Register: Authorizes the service as a client device with the given
// name under the user's Play Music account.  Re-registering a service
// can be used to change its name.  Remember that there is a limit to
// how many devices a Play Music account can have authorized, how many
// it can deauthorize in a year, and with how many accounts a single
// device can be authorized, so be judicious in using this method.
func (s *Service) Register(name string) *RegisterCall {
	return &RegisterCall{s, &pb.UpAuthRequest{
		UploaderId:   proto.String(s.id),
		FriendlyName: proto.String(name),
	}}
}

func (c *RegisterCall) Do() error {
	res, err := c.s.client.UpAuth(c.c)
	if err != nil {
		return err
	}
	status := res.GetAuthStatus()
	if status != pb.UploadResponse_OK {
		return RegisterError(status)
	}
	return nil
}

type JobsListCall struct {
	s *Service
	c *pb.GetJobsRequest
}

// List: Lists all outstanding upload jobs.
func (r *JobsService) List() *JobsListCall {
	return &JobsListCall{r.s, &pb.GetJobsRequest{
		UploaderId: proto.String(r.s.id),
	}}
}

func (c *JobsListCall) Do() ([]*pb.TracksToUpload, error) {
	res, err := c.s.client.GetJobs(c.c)
	if err != nil {
		return nil, err
	}
	jobs := res.GetGetjobsResponse()
	if jobs.GetGetTracksSuccess() == false {
		return nil, ErrJobsList
	}
	return jobs.GetTracksToUpload(), nil
}

type JobsCancelCall struct {
	s *Service
	c *pb.DeleteUploadRequestedRequest
}

// Cancel: Cancels all outstanding upload jobs.
func (r *JobsService) Cancel() *JobsCancelCall {
	return &JobsCancelCall{r.s, &pb.DeleteUploadRequestedRequest{
		UploaderId: proto.String(r.s.id),
	}}
}

func (c *JobsCancelCall) Do() error {
	_, err := c.s.client.DeleteUploadRequested(c.c)
	return err
}

type TracksGetCall struct {
	s *Service
	c *js.GetDownloadSessionRequest
}

// Get: Downloads a track by ID.  It is the user's responsibility to
// close the returned file.
func (r *TracksService) Get(serverID string) *TracksGetCall {
	return &TracksGetCall{r.s, &js.GetDownloadSessionRequest{
		XDeviceID: r.s.id,
		SongID:    serverID,
	}}
}

func (c *TracksGetCall) Do() (*mm.File, error) {
	url, err := c.s.client.GetDownloadSession(c.c)
	if err != nil {
		return nil, err
	}
	track, err := c.s.client.Get(url.URL)
	if err != nil {
		return nil, err
	}
	return track, nil
}

type TracksListCall struct {
	s *Service
	c *pb.GetTracksToExportRequest
}

// List: Lists tracks.
func (r *TracksService) List() *TracksListCall {
	return &TracksListCall{r.s, &pb.GetTracksToExportRequest{
		ClientId: proto.String(r.s.id),
	}}
}

// PageToken sets the continuation token used to page through large
// result sets.  To get the next page of results, set this parameter
// to the value of "ContinuationToken" from the previous response.
func (c *TracksListCall) PageToken(pageToken string) *TracksListCall {
	c.c.ContinuationToken = proto.String(pageToken)
	return c
}

// PurchasedOnly sets whether to list all tracks or free/purchased
// tracks only.
func (c *TracksListCall) PurchasedOnly(purchasedOnly bool) *TracksListCall {
	switch purchasedOnly {
	case false:
		c.c.ExportType = pb.GetTracksToExportRequest_ALL.Enum()
	case true:
		c.c.ExportType = pb.GetTracksToExportRequest_PURCHASED_AND_PROMOTIONAL.Enum()
	}
	return c
}

// UpdatedMin sets a cutoff date for the request: only tracks
// that were modified after this timestamp will be returned.  Expressed
// as a Unix timestamp in microseconds.  Specifying too large a value
// returns a 304 not modified *googleapi.Error when Do is called.
func (c *TracksListCall) UpdatedMin(updatedMin int64) *TracksListCall {
	c.c.UpdatedMin = proto.Int64(updatedMin)
	return c
}

func (c *TracksListCall) Do() (*pb.GetTracksToExportResponse, error) {
	list, err := c.s.client.GetTracksToExport(c.c)
	if err != nil {
		return nil, err
	}
	status := list.GetStatus()
	if status != pb.GetTracksToExportResponse_OK {
		return nil, TracksListError(status)
	}
	return list, nil
}

type TracksInsertCall struct {
	s        *Service
	track    io.ReadSeeker
	metadata *pb.Track
	sample   []byte
	pu       googleapi.ProgressUpdater
}

// Insert: Uploads a new song and returns its server ID on success.
// Behavior is undefined if the song is not in MP3 format.
func (r *TracksService) Insert(track io.ReadSeeker) *TracksInsertCall {
	return &TracksInsertCall{s: r.s, track: track, sample: []byte{}}
}

// BUG(lor): TracksService.Insert can't upload album art.
// BUG(lor): TracksService.Insert can't perform song matching.

// Metadata allows you to provide a track with custom metadata.
// In addition to what is documented under pb.Track, its following
// fields have special significance to TracksInsertCall:
//
// 	EstimatedSize   Used as the total parameter in any calls to the
// 	                ProgressUpdater callback.
//
// 	OriginalBitRate TracksInsertCall doesn't do any transcoding, so
// 	                this is also interpreted as the bitrate of the
// 	                uploaded file itself.  Can be omitted if
// 	                unknown or unimportant.
//
// If not called, metadata is read from the track with
// ReadMetadataFromMP3 when Do is called.  ReadMetadataFromMP3 ensures
// that the returned object contains at least the above fields
// (if the track lacks a title tag, "Untitled Track" is used.)
func (c *TracksInsertCall) Metadata(metadata *pb.Track) *TracksInsertCall {
	c.metadata = metadata
	return c
}

// Sample can be used to provide a 128kbps MP3 sample of the track,
// usually 15 seconds long, should the server ask for one in response
// to an upload request.  If not called, an empty sample that matches
// no tracks is sent instead.
func (c *TracksInsertCall) Sample(sample []byte) *TracksInsertCall {
	c.sample = sample
	return c
}

// BUG(lor): The usefulness of TracksInsertCall.Sample is somewhat
// dubious, as the server specifies start and end times for the sample
// in its request, information that is not available when the method
// should be called.

// ProgressUpdater provides a callback function that will be called
// after every uploaded chunk of the track.  It should be a low-latency
// function in order not to slow down the upload operation.
func (c *TracksInsertCall) ProgressUpdater(pu googleapi.ProgressUpdater) *TracksInsertCall {
	c.pu = pu
	return c
}

func (c *TracksInsertCall) Do() (string, error) {
	ress, err := c.s.Tracks.InsertBatch().Add(c, nil).Do()
	if err != nil {
		return "", err
	}
	for _, res := range ress {
		return res.ServerID, res.Error
	}
	// NOTREACHED
	return "", nil
}

type TracksInsertBatchCall struct {
	s *Service
	r tracksInsertBatchRequest
}

// InsertBatch: Batch uploads a set of tracks.  This call triggers the
// progress ticker in the Play Music web interface if called with two
// or more tracks.
func (r *TracksService) InsertBatch() *TracksInsertBatchCall {
	return &TracksInsertBatchCall{
		s: r.s,
		r: make(tracksInsertBatchRequest, 0),
	}
}

// Add adds a new insert call to the batch.  One must not call the Do
// method of the insert call before or after it has been added to the
// batch; the batch's own Do method handles everything.
//
// The callback function will be called with its return values when the
// insert call finishes.  It may be nil if not needed.  It is called
// synchronously, so it should be a low-latency function.
//
// It is an error to try to upload two tracks with the same client ID
// in a single batch, though this error will only be reported when the
// Do method is called, and only the latter of the two will be rejected.
// (Two tracks with omitted client IDs will also count as having the
// same ID.)
func (b *TracksInsertBatchCall) Add(c *TracksInsertCall, callback func(string, error)) *TracksInsertBatchCall {
	b.r = append(b.r, &tracksInsertBatchRequestEntry{c, callback})
	return b
}

// Do executes all insert calls in the batch synchronously.  It returns
// a slice of their return values, where the ith element corresponds to
// the ith-Added insert call.  The error value is only non-nil if a
// single error stops the entire batch; the user should check the
// individual responses in the slice for their error values, as it is
// very possible for Do to return a nil error without a single
// successful upload.
func (b *TracksInsertBatchCall) Do() (TracksInsertBatchResponse, error) {
	// A batch call with zero entries is a no-op.
	if len(b.r) == 0 {
		return nil, nil
	}
	// Initialize all bookkeeping and convenience variables.
	id := proto.String(b.s.id)
	ret := make(TracksInsertBatchResponse, len(b.r))
	mem := make(tracksInsertBatchRequestResponse)
	// Build the initial metadata upload request as well as all
	// bookkeeping structures, and perform the request.
	metadataReq := &pb.UploadMetadataRequest{
		UploaderId: id,
		Track:      make([]*pb.Track, 0),
	}
	for i, req := range b.r {
		// Create a new return value structure and associate
		// it with the corresponding request.
		res := new(TracksInsertBatchResponseEntry)
		e := &tracksInsertBatchRequestResponseEntry{req, res}
		ret[i] = res
		// If no metadata was provided with the request,
		// attempt to read it from the file.
		if req.metadata == nil {
			metadata, err := ReadMetadataFromMP3(req.track)
			if err != nil {
				e.retval(err)
				continue
			}
			e.metadata = metadata
		}
		// The mem map can only hold a single request-response
		// association per client ID, so if a track has the
		// same client ID as one in an earlier insert request,
		// we reject it.  It's an error to upload two tracks
		// with the same client ID anyway, so this is no big
		// loss.
		clientID := e.metadata.GetClientId()
		if _, ok := mem[clientID]; ok {
			e.retval(ErrDuplicateTrack)
			continue
		}
		// If everything went well, append the track metadata
		// to the upload request, and memoize the
		// request-response pair.
		metadataReq.Track = append(metadataReq.Track, e.metadata)
		mem[clientID] = e
	}
	res, err := b.s.client.UploadMetadata(metadataReq)
	if err != nil {
		mem.errall(err)
		return ret, err
	}
	// Process the metadata response, extracting in particular all
	// request-response pairs for which an upload was requested.
	metadataRes := res.GetMetadataResponse()
	metadataRess := metadataRes.GetTrackSampleResponse()
	sampleReqs := metadataRes.GetSignedChallengeInfo()
	toUpload := mem.processTrackResponses(metadataRess)
	// If any track samples were requested, fetch the samples
	// for them, process the response, and append any new upload
	// requests to toUpload.
	if len(sampleReqs) > 0 {
		sampleReq := &pb.UploadSampleRequest{
			UploaderId:  id,
			TrackSample: make([]*pb.TrackSample, len(sampleReqs)),
		}
		for i, sgi := range sampleReqs {
			clientID := sgi.GetChallengeInfo().GetClientTrackId()
			sampleReq.TrackSample[i] = &pb.TrackSample{
				Track:               mem[clientID].metadata,
				SignedChallengeInfo: sgi,
				Sample:              mem[clientID].sample,
			}
		}
		res, err := b.s.client.UploadSample(sampleReq)
		if err != nil {
			mem.errall(err)
		} else {
			sampleRess := res.GetSampleResponse().GetTrackSampleResponse()
			toUpload = append(toUpload, mem.processTrackResponses(sampleRess)...)
		}
	}
	// Quit early if no uploads were requested.
	if len(toUpload) == 0 {
		return ret, nil
	}
	// Notify the server that we are starting the upload process,
	// and upload the tracks.
	stateReq := &pb.UpdateUploadStateRequest{
		UploaderId: id,
		State:      pb.UpdateUploadStateRequest_START.Enum(),
	}
	b.s.client.UpdateUploadState(stateReq)
	for i, e := range toUpload {
		// Acquire an upload session.
		sessionReq := &js.GetUploadSessionRequest{
			Name:                      e.metadata.GetTitle(),
			UploaderId:                *id,
			ClientId:                  e.metadata.GetClientId(),
			ServerId:                  e.ServerID,
			TrackBitRate:              e.metadata.GetOriginalBitRate(),
			CurrentUploadingTrack:     e.metadata.GetTitle(),
			CurrentTotalUploadedCount: i,
			ClientTotalSongCount:      len(toUpload),
			SyncNow:                   true,
		}
		sessionRes, err := b.s.client.GetUploadSession(sessionReq)
		if err := newSessionError(sessionRes, err); err != nil {
			e.retval(err)
			continue
		}
		// Upload a track.
		body := e.track.(io.Reader)
		if e.pu != nil {
			body = &progressReader{
				r:     body,
				pu:    e.pu,
				total: e.metadata.GetEstimatedSize(),
			}
		}
		uploadRes, err := b.s.client.Put(sessionRes.TransferPutUrl, body)
		e.retval(newSessionError(uploadRes, err))
	}
	// Notify the server that we are done uploading, and return
	// the results.
	stateReq.State = pb.UpdateUploadStateRequest_STOPPED.Enum()
	b.s.client.UpdateUploadState(stateReq)
	return ret, nil
}
