package musicmanager

import (
	"fmt"

	"google.golang.org/api/googleapi"

	js "google-musicmanager-go/json"
	pb "google-musicmanager-go/proto"
)

// Returned by TracksInsertBatchCall.Do when a track has the same
// client ID as a previous track in the batch.
var ErrDuplicateTrack = fmt.Errorf("musicmanager: duplicate track")

// Returned by JobsListCall.Do on an indeterminate error.
var ErrJobsList = fmt.Errorf("musicmanager: could not get jobs")

// A RegisterError is returned by RegisterCall.Do when the server can't
// register the client.
type RegisterError pb.UploadResponse_AuthStatus

// GetRegisterError gets a RegisterError value by name.  For a
// list of the possible names, see the definition of
// pb.UploadResponse_AuthStatus.
func GetRegisterError(name string) RegisterError {
	return RegisterError(pb.UploadResponse_AuthStatus_value[name])
}

func (e RegisterError) Error() string {
	return pb.UploadResponse_AuthStatus(e).String()
}

// A TracksInsertError is returned by TracksInsertCall.Do when the
// server refuses to accept a track based on its metadata.
type TracksInsertError pb.TrackSampleResponse_ResponseCode

// GetTracksInsertError gets a TracksInsertError value by name.  For a
// list of the possible names, see the definition of
// pb.TrackSampleResponse_ResponseCode.
func GetTracksInsertError(name string) TracksInsertError {
	return TracksInsertError(pb.TrackSampleResponse_ResponseCode_value[name])
}

func (e TracksInsertError) Error() string {
	return pb.TrackSampleResponse_ResponseCode(e).String()
}

// A SessionError is returned by TracksInsertCall.Do when requesting
// an upload session, or sending data to one responds with a non-2xx
// status code.  Its value is the value of that code.
type SessionError int

// newSessionError converts a *googleapi.Error or a response object
// containing a nonzero error code into a SessionError while passing
// all other error values through unchanged.
func newSessionError(res *js.GetUploadSessionResponse, err error) error {
	switch err := err.(type) {
	case nil:
		// do nothing
	case *googleapi.Error:
		return SessionError(err.Code)
	default:
		return err
	}
	if res.ErrorCode != 0 {
		return SessionError(res.ErrorCode)
	}
	return nil
}

func (e SessionError) Error() string {
	return fmt.Sprintf("musicmanager: upload session failed with code %d", int(e))
}

// A TracksListError is returned by TracksListCall.Do when the server
// can't list the tracks.
type TracksListError pb.GetTracksToExportResponse_TracksToExportStatus

// GetTracksListError gets a TracksListError value by name.  For a
// list of the possible names, see the definition of
// pb.GetTracksToExportResponse_TracksToExportStatus.
func GetTracksListError(name string) TracksListError {
	return TracksListError(pb.GetTracksToExportResponse_TracksToExportStatus_value[name])
}

func (e TracksListError) Error() string {
	return pb.GetTracksToExportResponse_TracksToExportStatus(e).String()
}
