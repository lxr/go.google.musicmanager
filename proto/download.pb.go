// Code generated by protoc-gen-go.
// source: download.proto
// DO NOT EDIT!

package musicmanager

import proto "github.com/golang/protobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

// Different track types for exporting.
type GetTracksToExportRequest_TracksToExportType int32

const (
	GetTracksToExportRequest_ALL                       GetTracksToExportRequest_TracksToExportType = 1
	GetTracksToExportRequest_PURCHASED_AND_PROMOTIONAL GetTracksToExportRequest_TracksToExportType = 2
)

var GetTracksToExportRequest_TracksToExportType_name = map[int32]string{
	1: "ALL",
	2: "PURCHASED_AND_PROMOTIONAL",
}
var GetTracksToExportRequest_TracksToExportType_value = map[string]int32{
	"ALL": 1,
	"PURCHASED_AND_PROMOTIONAL": 2,
}

func (x GetTracksToExportRequest_TracksToExportType) Enum() *GetTracksToExportRequest_TracksToExportType {
	p := new(GetTracksToExportRequest_TracksToExportType)
	*p = x
	return p
}
func (x GetTracksToExportRequest_TracksToExportType) String() string {
	return proto.EnumName(GetTracksToExportRequest_TracksToExportType_name, int32(x))
}
func (x *GetTracksToExportRequest_TracksToExportType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(GetTracksToExportRequest_TracksToExportType_value, data, "GetTracksToExportRequest_TracksToExportType")
	if err != nil {
		return err
	}
	*x = GetTracksToExportRequest_TracksToExportType(value)
	return nil
}

// Status codes for DownloadService.GetTracksToExport.
// The exact meanings of the error codes are unknown;
// only OK is guaranteed to mean success.
type GetTracksToExportResponse_TracksToExportStatus int32

const (
	GetTracksToExportResponse_OK                            GetTracksToExportResponse_TracksToExportStatus = 1
	GetTracksToExportResponse_TRANSIENT_ERROR               GetTracksToExportResponse_TracksToExportStatus = 2
	GetTracksToExportResponse_MAX_NUM_CLIENTS_REACHED       GetTracksToExportResponse_TracksToExportStatus = 3
	GetTracksToExportResponse_UNABLE_TO_AUTHENTICATE_CLIENT GetTracksToExportResponse_TracksToExportStatus = 4
	GetTracksToExportResponse_UNABLE_TO_REGISTER_CLIENT     GetTracksToExportResponse_TracksToExportStatus = 5
)

var GetTracksToExportResponse_TracksToExportStatus_name = map[int32]string{
	1: "OK",
	2: "TRANSIENT_ERROR",
	3: "MAX_NUM_CLIENTS_REACHED",
	4: "UNABLE_TO_AUTHENTICATE_CLIENT",
	5: "UNABLE_TO_REGISTER_CLIENT",
}
var GetTracksToExportResponse_TracksToExportStatus_value = map[string]int32{
	"OK":                            1,
	"TRANSIENT_ERROR":               2,
	"MAX_NUM_CLIENTS_REACHED":       3,
	"UNABLE_TO_AUTHENTICATE_CLIENT": 4,
	"UNABLE_TO_REGISTER_CLIENT":     5,
}

func (x GetTracksToExportResponse_TracksToExportStatus) Enum() *GetTracksToExportResponse_TracksToExportStatus {
	p := new(GetTracksToExportResponse_TracksToExportStatus)
	*p = x
	return p
}
func (x GetTracksToExportResponse_TracksToExportStatus) String() string {
	return proto.EnumName(GetTracksToExportResponse_TracksToExportStatus_name, int32(x))
}
func (x *GetTracksToExportResponse_TracksToExportStatus) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(GetTracksToExportResponse_TracksToExportStatus_value, data, "GetTracksToExportResponse_TracksToExportStatus")
	if err != nil {
		return err
	}
	*x = GetTracksToExportResponse_TracksToExportStatus(value)
	return nil
}

// Arguments to DownloadService.GetTracksToExport.
type GetTracksToExportRequest struct {
	// The device ID.
	ClientId *string `protobuf:"bytes,2,req,name=client_id" json:"client_id,omitempty"`
	// A token for paging through large result sets.  If the result set
	// contains more than 1000 tracks, the response contains a
	// continuation_token field that can be used to retrieve the next
	// page in the result set.
	ContinuationToken *string `protobuf:"bytes,3,opt,name=continuation_token" json:"continuation_token,omitempty"`
	// Only list tracks of this type.
	ExportType *GetTracksToExportRequest_TracksToExportType `protobuf:"varint,4,opt,name=export_type,enum=musicmanager.GetTracksToExportRequest_TracksToExportType" json:"export_type,omitempty"`
	// Only list tracks that were modified after this time.  Expressed as
	// a Unix timestamp in microseconds.
	UpdatedMin       *int64 `protobuf:"varint,5,opt,name=updated_min" json:"updated_min,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *GetTracksToExportRequest) Reset()         { *m = GetTracksToExportRequest{} }
func (m *GetTracksToExportRequest) String() string { return proto.CompactTextString(m) }
func (*GetTracksToExportRequest) ProtoMessage()    {}

func (m *GetTracksToExportRequest) GetClientId() string {
	if m != nil && m.ClientId != nil {
		return *m.ClientId
	}
	return ""
}

func (m *GetTracksToExportRequest) GetContinuationToken() string {
	if m != nil && m.ContinuationToken != nil {
		return *m.ContinuationToken
	}
	return ""
}

func (m *GetTracksToExportRequest) GetExportType() GetTracksToExportRequest_TracksToExportType {
	if m != nil && m.ExportType != nil {
		return *m.ExportType
	}
	return GetTracksToExportRequest_ALL
}

func (m *GetTracksToExportRequest) GetUpdatedMin() int64 {
	if m != nil && m.UpdatedMin != nil {
		return *m.UpdatedMin
	}
	return 0
}

// Track metadata returned by DownloadService.GetTracksToExport.
type DownloadTrackInfo struct {
	// The server ID of the track.
	Id *string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	// The title of the track.
	Title *string `protobuf:"bytes,2,opt,name=title" json:"title,omitempty"`
	// The album of the track.
	Album *string `protobuf:"bytes,3,opt,name=album" json:"album,omitempty"`
	// The album artist of the track.
	AlbumArtist *string `protobuf:"bytes,4,opt,name=album_artist" json:"album_artist,omitempty"`
	// The artist of the track.
	Artist *string `protobuf:"bytes,5,opt,name=artist" json:"artist,omitempty"`
	// The number of the track within the album.
	TrackNumber *int32 `protobuf:"varint,6,opt,name=track_number" json:"track_number,omitempty"`
	// For uploaded tracks, this is the size of the track.
	// For purchased and promotional tracks, the value appears somewhat
	// arbitrary.
	TrackSize        *int64 `protobuf:"varint,7,opt,name=track_size" json:"track_size,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *DownloadTrackInfo) Reset()         { *m = DownloadTrackInfo{} }
func (m *DownloadTrackInfo) String() string { return proto.CompactTextString(m) }
func (*DownloadTrackInfo) ProtoMessage()    {}

func (m *DownloadTrackInfo) GetId() string {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return ""
}

func (m *DownloadTrackInfo) GetTitle() string {
	if m != nil && m.Title != nil {
		return *m.Title
	}
	return ""
}

func (m *DownloadTrackInfo) GetAlbum() string {
	if m != nil && m.Album != nil {
		return *m.Album
	}
	return ""
}

func (m *DownloadTrackInfo) GetAlbumArtist() string {
	if m != nil && m.AlbumArtist != nil {
		return *m.AlbumArtist
	}
	return ""
}

func (m *DownloadTrackInfo) GetArtist() string {
	if m != nil && m.Artist != nil {
		return *m.Artist
	}
	return ""
}

func (m *DownloadTrackInfo) GetTrackNumber() int32 {
	if m != nil && m.TrackNumber != nil {
		return *m.TrackNumber
	}
	return 0
}

func (m *DownloadTrackInfo) GetTrackSize() int64 {
	if m != nil && m.TrackSize != nil {
		return *m.TrackSize
	}
	return 0
}

// Return values of DownloadService.GetTracksToExport.
type GetTracksToExportResponse struct {
	// The status code of the response.
	Status *GetTracksToExportResponse_TracksToExportStatus `protobuf:"varint,1,req,name=status,enum=musicmanager.GetTracksToExportResponse_TracksToExportStatus" json:"status,omitempty"`
	// The tracks.  They appear to be ordered from least recently to
	// most recently accessed.
	DownloadTrackInfo []*DownloadTrackInfo `protobuf:"bytes,2,rep,name=download_track_info" json:"download_track_info,omitempty"`
	// The page token for the next page of tracks.
	ContinuationToken *string `protobuf:"bytes,3,opt,name=continuation_token" json:"continuation_token,omitempty"`
	// The last time a track was modified.  Expressed as a Unix timestamp
	// in microseconds.
	UpdatedMin       *int64 `protobuf:"varint,4,opt,name=updated_min" json:"updated_min,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *GetTracksToExportResponse) Reset()         { *m = GetTracksToExportResponse{} }
func (m *GetTracksToExportResponse) String() string { return proto.CompactTextString(m) }
func (*GetTracksToExportResponse) ProtoMessage()    {}

func (m *GetTracksToExportResponse) GetStatus() GetTracksToExportResponse_TracksToExportStatus {
	if m != nil && m.Status != nil {
		return *m.Status
	}
	return GetTracksToExportResponse_OK
}

func (m *GetTracksToExportResponse) GetDownloadTrackInfo() []*DownloadTrackInfo {
	if m != nil {
		return m.DownloadTrackInfo
	}
	return nil
}

func (m *GetTracksToExportResponse) GetContinuationToken() string {
	if m != nil && m.ContinuationToken != nil {
		return *m.ContinuationToken
	}
	return ""
}

func (m *GetTracksToExportResponse) GetUpdatedMin() int64 {
	if m != nil && m.UpdatedMin != nil {
		return *m.UpdatedMin
	}
	return 0
}

func init() {
	proto.RegisterEnum("musicmanager.GetTracksToExportRequest_TracksToExportType", GetTracksToExportRequest_TracksToExportType_name, GetTracksToExportRequest_TracksToExportType_value)
	proto.RegisterEnum("musicmanager.GetTracksToExportResponse_TracksToExportStatus", GetTracksToExportResponse_TracksToExportStatus_name, GetTracksToExportResponse_TracksToExportStatus_value)
}
