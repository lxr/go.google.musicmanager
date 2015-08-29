// Code generated by protoc-gen-go.
// source: my-git.appspot.com/go.google.musicmanager/internal/download_proto/service/download.messages.proto
// DO NOT EDIT!

/*
Package google_musicmanager_v0 is a generated protocol buffer package.

It is generated from these files:
	my-git.appspot.com/go.google.musicmanager/internal/download_proto/service/download.messages.proto

It has these top-level messages:
	GetTracksToExportRequest
	GetTracksToExportResponse
*/
package google_musicmanager_v0

import proto "github.com/golang/protobuf/proto"
import google_musicmanager_v01 "my-git.appspot.com/go.google.musicmanager/internal/download_proto/data"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal

type GetTracksToExportRequest_TracksToExportType int32

const (
	GetTracksToExportRequest_UNKNOWN                   GetTracksToExportRequest_TracksToExportType = 0
	GetTracksToExportRequest_ALL                       GetTracksToExportRequest_TracksToExportType = 1
	GetTracksToExportRequest_PURCHASED_AND_PROMOTIONAL GetTracksToExportRequest_TracksToExportType = 2
)

var GetTracksToExportRequest_TracksToExportType_name = map[int32]string{
	0: "UNKNOWN",
	1: "ALL",
	2: "PURCHASED_AND_PROMOTIONAL",
}
var GetTracksToExportRequest_TracksToExportType_value = map[string]int32{
	"UNKNOWN": 0,
	"ALL":     1,
	"PURCHASED_AND_PROMOTIONAL": 2,
}

func (x GetTracksToExportRequest_TracksToExportType) String() string {
	return proto.EnumName(GetTracksToExportRequest_TracksToExportType_name, int32(x))
}

type GetTracksToExportResponse_TracksToExportStatus int32

const (
	GetTracksToExportResponse_UNKNOWN                       GetTracksToExportResponse_TracksToExportStatus = 0
	GetTracksToExportResponse_OK                            GetTracksToExportResponse_TracksToExportStatus = 1
	GetTracksToExportResponse_TRANSIENT_ERROR               GetTracksToExportResponse_TracksToExportStatus = 2
	GetTracksToExportResponse_MAX_NUM_CLIENTS_REACHED       GetTracksToExportResponse_TracksToExportStatus = 3
	GetTracksToExportResponse_UNABLE_TO_AUTHENTICATE_CLIENT GetTracksToExportResponse_TracksToExportStatus = 4
	GetTracksToExportResponse_UNABLE_TO_REGISTER_CLIENT     GetTracksToExportResponse_TracksToExportStatus = 5
)

var GetTracksToExportResponse_TracksToExportStatus_name = map[int32]string{
	0: "UNKNOWN",
	1: "OK",
	2: "TRANSIENT_ERROR",
	3: "MAX_NUM_CLIENTS_REACHED",
	4: "UNABLE_TO_AUTHENTICATE_CLIENT",
	5: "UNABLE_TO_REGISTER_CLIENT",
}
var GetTracksToExportResponse_TracksToExportStatus_value = map[string]int32{
	"UNKNOWN":                       0,
	"OK":                            1,
	"TRANSIENT_ERROR":               2,
	"MAX_NUM_CLIENTS_REACHED":       3,
	"UNABLE_TO_AUTHENTICATE_CLIENT": 4,
	"UNABLE_TO_REGISTER_CLIENT":     5,
}

func (x GetTracksToExportResponse_TracksToExportStatus) String() string {
	return proto.EnumName(GetTracksToExportResponse_TracksToExportStatus_name, int32(x))
}

type GetTracksToExportRequest struct {
	ClientId          string                                      `protobuf:"bytes,2,opt,name=client_id" json:"client_id,omitempty"`
	ContinuationToken string                                      `protobuf:"bytes,3,opt,name=continuation_token" json:"continuation_token,omitempty"`
	ExportType        GetTracksToExportRequest_TracksToExportType `protobuf:"varint,4,opt,name=export_type,enum=google.musicmanager.v0.GetTracksToExportRequest_TracksToExportType" json:"export_type,omitempty"`
	UpdatedMin        int64                                       `protobuf:"varint,5,opt,name=updated_min" json:"updated_min,omitempty"`
}

func (m *GetTracksToExportRequest) Reset()         { *m = GetTracksToExportRequest{} }
func (m *GetTracksToExportRequest) String() string { return proto.CompactTextString(m) }
func (*GetTracksToExportRequest) ProtoMessage()    {}

type GetTracksToExportResponse struct {
	Status            GetTracksToExportResponse_TracksToExportStatus `protobuf:"varint,1,opt,name=status,enum=google.musicmanager.v0.GetTracksToExportResponse_TracksToExportStatus" json:"status,omitempty"`
	DownloadTrackInfo []*google_musicmanager_v01.DownloadTrackInfo   `protobuf:"bytes,2,rep,name=download_track_info" json:"download_track_info,omitempty"`
	ContinuationToken string                                         `protobuf:"bytes,3,opt,name=continuation_token" json:"continuation_token,omitempty"`
	UpdatedMin        int64                                          `protobuf:"varint,4,opt,name=updated_min" json:"updated_min,omitempty"`
}

func (m *GetTracksToExportResponse) Reset()         { *m = GetTracksToExportResponse{} }
func (m *GetTracksToExportResponse) String() string { return proto.CompactTextString(m) }
func (*GetTracksToExportResponse) ProtoMessage()    {}

func (m *GetTracksToExportResponse) GetDownloadTrackInfo() []*google_musicmanager_v01.DownloadTrackInfo {
	if m != nil {
		return m.DownloadTrackInfo
	}
	return nil
}

func init() {
	proto.RegisterEnum("google.musicmanager.v0.GetTracksToExportRequest_TracksToExportType", GetTracksToExportRequest_TracksToExportType_name, GetTracksToExportRequest_TracksToExportType_value)
	proto.RegisterEnum("google.musicmanager.v0.GetTracksToExportResponse_TracksToExportStatus", GetTracksToExportResponse_TracksToExportStatus_name, GetTracksToExportResponse_TracksToExportStatus_value)
}
