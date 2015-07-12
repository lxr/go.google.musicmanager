package musicmanager

import (
	"io"

	"google.golang.org/api/googleapi"

	pb "google-musicmanager-go/proto"
)

// A progressReader calls its pu function every time it is read from
// with the number of bytes read so far, and the number of bytes in
// total.
type progressReader struct {
	r       io.Reader
	pu      googleapi.ProgressUpdater
	current int64
	total   int64
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.r.Read(p)
	pr.current += int64(n)
	pr.pu(pr.current, pr.total)
	return n, err
}

type tracksInsertBatchRequest []*tracksInsertBatchRequestEntry

type tracksInsertBatchRequestEntry struct {
	*TracksInsertCall
	callback func(string, error)
}

// A tracksInsertBatchRequestResponse is a map from client IDs to
// in-progress insert call structures.  It is used as a bookkeeping
// structure by TracksInsertBatchCall.Do.
type tracksInsertBatchRequestResponse map[string]*tracksInsertBatchRequestResponseEntry

// A tracksInsertBatchRequestResponseEntry represents one in-progress
// insert call in a batch.  It comprises the original insert request
// and its not-yet-complete response.
type tracksInsertBatchRequestResponseEntry struct {
	*tracksInsertBatchRequestEntry
	*TracksInsertBatchResponseEntry
}

// errall terminates all calls in a batch that haven't terminated yet
// with the given error.
func (r tracksInsertBatchRequestResponse) errall(err error) {
	for _, e := range r {
		if e.ServerID == "" && e.Error == nil {
			e.retval(err)
		}
	}
}

// processTrackResponses parses a slice of TrackSampleResponses and
// terminates any batch call whose client ID did not receive an
// upload request with the corresponding error code; the rest are
// assigned their server ID, and also collected in a slice that is
// returned in the end.
func (r tracksInsertBatchRequestResponse) processTrackResponses(ress []*pb.TrackSampleResponse) []*tracksInsertBatchRequestResponseEntry {
	toUpload := make([]*tracksInsertBatchRequestResponseEntry, 0)
	for _, res := range ress {
		clientID := res.GetClientTrackId()
		status := res.GetResponseCode()
		e := r[clientID]
		if status != pb.TrackSampleResponse_UPLOAD_REQUESTED {
			e.retval(TracksInsertError(status))
			continue
		}
		e.ServerID = res.GetServerTrackId()
		toUpload = append(toUpload, e)
	}
	return toUpload
}

// retval terminates a batch call with the given error.
func (e *tracksInsertBatchRequestResponseEntry) retval(err error) {
	e.Error = err
	if e.callback != nil {
		e.callback(e.ServerID, e.Error)
	}
}
