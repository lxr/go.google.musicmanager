// Code generated by protoc-gen-go.
// source: github.com/lxr/go.google.musicmanager/internal/uits_proto/data/uits.data.proto
// DO NOT EDIT!

/*
Package google_musicmanager_v0 is a generated protocol buffer package.

It is generated from these files:
	github.com/lxr/go.google.musicmanager/internal/uits_proto/data/uits.data.proto

It has these top-level messages:
	ProductId
	AssetId
	TransactionId
	MediaId
	UrlInfo
	CopyrightStatus
	Extra
	UitsMetadata
	UitsSignature
	Uits
	UploadedUitsId3Tag
*/
package google_musicmanager_v0

import proto "github.com/golang/protobuf/proto"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal

type ProductId_Type int32

const (
	ProductId_DEFAULT ProductId_Type = 0
	ProductId_UPC     ProductId_Type = 1
	ProductId_GRID    ProductId_Type = 2
)

var ProductId_Type_name = map[int32]string{
	0: "DEFAULT",
	1: "UPC",
	2: "GRID",
}
var ProductId_Type_value = map[string]int32{
	"DEFAULT": 0,
	"UPC":     1,
	"GRID":    2,
}

func (x ProductId_Type) String() string {
	return proto.EnumName(ProductId_Type_name, int32(x))
}

type AssetId_Type int32

const (
	AssetId_DEFAULT AssetId_Type = 0
	AssetId_ISRC    AssetId_Type = 1
)

var AssetId_Type_name = map[int32]string{
	0: "DEFAULT",
	1: "ISRC",
}
var AssetId_Type_value = map[string]int32{
	"DEFAULT": 0,
	"ISRC":    1,
}

func (x AssetId_Type) String() string {
	return proto.EnumName(AssetId_Type_name, int32(x))
}

type MediaId_AlgorithmType int32

const (
	MediaId_DEFAULT MediaId_AlgorithmType = 0
	MediaId_SHA256  MediaId_AlgorithmType = 1
)

var MediaId_AlgorithmType_name = map[int32]string{
	0: "DEFAULT",
	1: "SHA256",
}
var MediaId_AlgorithmType_value = map[string]int32{
	"DEFAULT": 0,
	"SHA256":  1,
}

func (x MediaId_AlgorithmType) String() string {
	return proto.EnumName(MediaId_AlgorithmType_name, int32(x))
}

type UrlInfo_Type int32

const (
	UrlInfo_DEFAULT UrlInfo_Type = 0
	UrlInfo_WCOM    UrlInfo_Type = 1
	UrlInfo_WCOP    UrlInfo_Type = 2
	UrlInfo_WOAF    UrlInfo_Type = 3
	UrlInfo_WOAR    UrlInfo_Type = 4
	UrlInfo_WOAS    UrlInfo_Type = 5
	UrlInfo_WORS    UrlInfo_Type = 6
	UrlInfo_WPAY    UrlInfo_Type = 7
	UrlInfo_WPUB    UrlInfo_Type = 8
)

var UrlInfo_Type_name = map[int32]string{
	0: "DEFAULT",
	1: "WCOM",
	2: "WCOP",
	3: "WOAF",
	4: "WOAR",
	5: "WOAS",
	6: "WORS",
	7: "WPAY",
	8: "WPUB",
}
var UrlInfo_Type_value = map[string]int32{
	"DEFAULT": 0,
	"WCOM":    1,
	"WCOP":    2,
	"WOAF":    3,
	"WOAR":    4,
	"WOAS":    5,
	"WORS":    6,
	"WPAY":    7,
	"WPUB":    8,
}

func (x UrlInfo_Type) String() string {
	return proto.EnumName(UrlInfo_Type_name, int32(x))
}

type CopyrightStatus_Type int32

const (
	CopyrightStatus_DEFAULT           CopyrightStatus_Type = 0
	CopyrightStatus_UNSPECIFIED       CopyrightStatus_Type = 1
	CopyrightStatus_ALLRIGHTSRESERVED CopyrightStatus_Type = 2
	CopyrightStatus_PRERELEASE        CopyrightStatus_Type = 3
	CopyrightStatus_OTHER             CopyrightStatus_Type = 4
)

var CopyrightStatus_Type_name = map[int32]string{
	0: "DEFAULT",
	1: "UNSPECIFIED",
	2: "ALLRIGHTSRESERVED",
	3: "PRERELEASE",
	4: "OTHER",
}
var CopyrightStatus_Type_value = map[string]int32{
	"DEFAULT":           0,
	"UNSPECIFIED":       1,
	"ALLRIGHTSRESERVED": 2,
	"PRERELEASE":        3,
	"OTHER":             4,
}

func (x CopyrightStatus_Type) String() string {
	return proto.EnumName(CopyrightStatus_Type_name, int32(x))
}

type UitsMetadata_ParentalAdvisoryType int32

const (
	UitsMetadata_DEFAULT     UitsMetadata_ParentalAdvisoryType = 0
	UitsMetadata_UNSPECIFIED UitsMetadata_ParentalAdvisoryType = 1
	UitsMetadata_EXPLICIT    UitsMetadata_ParentalAdvisoryType = 2
	UitsMetadata_EDITED      UitsMetadata_ParentalAdvisoryType = 3
)

var UitsMetadata_ParentalAdvisoryType_name = map[int32]string{
	0: "DEFAULT",
	1: "UNSPECIFIED",
	2: "EXPLICIT",
	3: "EDITED",
}
var UitsMetadata_ParentalAdvisoryType_value = map[string]int32{
	"DEFAULT":     0,
	"UNSPECIFIED": 1,
	"EXPLICIT":    2,
	"EDITED":      3,
}

func (x UitsMetadata_ParentalAdvisoryType) String() string {
	return proto.EnumName(UitsMetadata_ParentalAdvisoryType_name, int32(x))
}

type UitsSignature_AlgorithmType int32

const (
	UitsSignature_DEFAULT_ALGORITHM UitsSignature_AlgorithmType = 0
	UitsSignature_RSA2048           UitsSignature_AlgorithmType = 1
	UitsSignature_DSA2048           UitsSignature_AlgorithmType = 2
)

var UitsSignature_AlgorithmType_name = map[int32]string{
	0: "DEFAULT_ALGORITHM",
	1: "RSA2048",
	2: "DSA2048",
}
var UitsSignature_AlgorithmType_value = map[string]int32{
	"DEFAULT_ALGORITHM": 0,
	"RSA2048":           1,
	"DSA2048":           2,
}

func (x UitsSignature_AlgorithmType) String() string {
	return proto.EnumName(UitsSignature_AlgorithmType_name, int32(x))
}

type UitsSignature_CanonicalizationType int32

const (
	UitsSignature_DEFAULT_CANONICALIZATION UitsSignature_CanonicalizationType = 0
	UitsSignature_NONE                     UitsSignature_CanonicalizationType = 1
)

var UitsSignature_CanonicalizationType_name = map[int32]string{
	0: "DEFAULT_CANONICALIZATION",
	1: "NONE",
}
var UitsSignature_CanonicalizationType_value = map[string]int32{
	"DEFAULT_CANONICALIZATION": 0,
	"NONE": 1,
}

func (x UitsSignature_CanonicalizationType) String() string {
	return proto.EnumName(UitsSignature_CanonicalizationType_name, int32(x))
}

type ProductId struct {
	Type      ProductId_Type `protobuf:"varint,1,opt,name=type,enum=google.musicmanager.v0.ProductId_Type" json:"type,omitempty"`
	Completed bool           `protobuf:"varint,2,opt,name=completed" json:"completed,omitempty"`
	Id        string         `protobuf:"bytes,3,opt,name=id" json:"id,omitempty"`
}

func (m *ProductId) Reset()         { *m = ProductId{} }
func (m *ProductId) String() string { return proto.CompactTextString(m) }
func (*ProductId) ProtoMessage()    {}

type AssetId struct {
	Type AssetId_Type `protobuf:"varint,1,opt,name=type,enum=google.musicmanager.v0.AssetId_Type" json:"type,omitempty"`
	Id   string       `protobuf:"bytes,2,opt,name=id" json:"id,omitempty"`
}

func (m *AssetId) Reset()         { *m = AssetId{} }
func (m *AssetId) String() string { return proto.CompactTextString(m) }
func (*AssetId) ProtoMessage()    {}

type TransactionId struct {
	Version string `protobuf:"bytes,1,opt,name=version" json:"version,omitempty"`
	Id      string `protobuf:"bytes,2,opt,name=id" json:"id,omitempty"`
}

func (m *TransactionId) Reset()         { *m = TransactionId{} }
func (m *TransactionId) String() string { return proto.CompactTextString(m) }
func (*TransactionId) ProtoMessage()    {}

type MediaId struct {
	AlgorithmType MediaId_AlgorithmType `protobuf:"varint,1,opt,name=algorithm_type,enum=google.musicmanager.v0.MediaId_AlgorithmType" json:"algorithm_type,omitempty"`
	Hash          string                `protobuf:"bytes,2,opt,name=hash" json:"hash,omitempty"`
}

func (m *MediaId) Reset()         { *m = MediaId{} }
func (m *MediaId) String() string { return proto.CompactTextString(m) }
func (*MediaId) ProtoMessage()    {}

type UrlInfo struct {
	Type UrlInfo_Type `protobuf:"varint,1,opt,name=type,enum=google.musicmanager.v0.UrlInfo_Type" json:"type,omitempty"`
	Url  string       `protobuf:"bytes,2,opt,name=url" json:"url,omitempty"`
}

func (m *UrlInfo) Reset()         { *m = UrlInfo{} }
func (m *UrlInfo) String() string { return proto.CompactTextString(m) }
func (*UrlInfo) ProtoMessage()    {}

type CopyrightStatus struct {
	Type      CopyrightStatus_Type `protobuf:"varint,1,opt,name=type,enum=google.musicmanager.v0.CopyrightStatus_Type" json:"type,omitempty"`
	Copyright string               `protobuf:"bytes,2,opt,name=copyright" json:"copyright,omitempty"`
}

func (m *CopyrightStatus) Reset()         { *m = CopyrightStatus{} }
func (m *CopyrightStatus) String() string { return proto.CompactTextString(m) }
func (*CopyrightStatus) ProtoMessage()    {}

type Extra struct {
	Type  string `protobuf:"bytes,1,opt,name=type" json:"type,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
}

func (m *Extra) Reset()         { *m = Extra{} }
func (m *Extra) String() string { return proto.CompactTextString(m) }
func (*Extra) ProtoMessage()    {}

type UitsMetadata struct {
	Nonce                string                            `protobuf:"bytes,1,opt,name=nonce" json:"nonce,omitempty"`
	DistributorId        string                            `protobuf:"bytes,2,opt,name=distributor_id" json:"distributor_id,omitempty"`
	TransactionDate      string                            `protobuf:"bytes,3,opt,name=transaction_date" json:"transaction_date,omitempty"`
	ProductId            *ProductId                        `protobuf:"bytes,4,opt,name=product_id" json:"product_id,omitempty"`
	AssetId              *AssetId                          `protobuf:"bytes,5,opt,name=asset_id" json:"asset_id,omitempty"`
	TransactionId        *TransactionId                    `protobuf:"bytes,6,opt,name=transaction_id" json:"transaction_id,omitempty"`
	MediaId              *MediaId                          `protobuf:"bytes,7,opt,name=media_id" json:"media_id,omitempty"`
	UrlInfo              *UrlInfo                          `protobuf:"bytes,8,opt,name=url_info" json:"url_info,omitempty"`
	ParentalAdvisoryType UitsMetadata_ParentalAdvisoryType `protobuf:"varint,9,opt,name=parental_advisory_type,enum=google.musicmanager.v0.UitsMetadata_ParentalAdvisoryType" json:"parental_advisory_type,omitempty"`
	CopyrightStatus      *CopyrightStatus                  `protobuf:"bytes,10,opt,name=copyright_status" json:"copyright_status,omitempty"`
	Extra                []*Extra                          `protobuf:"bytes,11,rep,name=extra" json:"extra,omitempty"`
}

func (m *UitsMetadata) Reset()         { *m = UitsMetadata{} }
func (m *UitsMetadata) String() string { return proto.CompactTextString(m) }
func (*UitsMetadata) ProtoMessage()    {}

func (m *UitsMetadata) GetProductId() *ProductId {
	if m != nil {
		return m.ProductId
	}
	return nil
}

func (m *UitsMetadata) GetAssetId() *AssetId {
	if m != nil {
		return m.AssetId
	}
	return nil
}

func (m *UitsMetadata) GetTransactionId() *TransactionId {
	if m != nil {
		return m.TransactionId
	}
	return nil
}

func (m *UitsMetadata) GetMediaId() *MediaId {
	if m != nil {
		return m.MediaId
	}
	return nil
}

func (m *UitsMetadata) GetUrlInfo() *UrlInfo {
	if m != nil {
		return m.UrlInfo
	}
	return nil
}

func (m *UitsMetadata) GetCopyrightStatus() *CopyrightStatus {
	if m != nil {
		return m.CopyrightStatus
	}
	return nil
}

func (m *UitsMetadata) GetExtra() []*Extra {
	if m != nil {
		return m.Extra
	}
	return nil
}

type UitsSignature struct {
	AlgorithmType        UitsSignature_AlgorithmType        `protobuf:"varint,1,opt,name=algorithm_type,enum=google.musicmanager.v0.UitsSignature_AlgorithmType" json:"algorithm_type,omitempty"`
	CanonicalizationType UitsSignature_CanonicalizationType `protobuf:"varint,2,opt,name=canonicalization_type,enum=google.musicmanager.v0.UitsSignature_CanonicalizationType" json:"canonicalization_type,omitempty"`
	KeyId                string                             `protobuf:"bytes,3,opt,name=key_id" json:"key_id,omitempty"`
	Value                string                             `protobuf:"bytes,4,opt,name=value" json:"value,omitempty"`
}

func (m *UitsSignature) Reset()         { *m = UitsSignature{} }
func (m *UitsSignature) String() string { return proto.CompactTextString(m) }
func (*UitsSignature) ProtoMessage()    {}

type Uits struct {
	Metadata  *UitsMetadata  `protobuf:"bytes,1,opt,name=metadata" json:"metadata,omitempty"`
	Signature *UitsSignature `protobuf:"bytes,2,opt,name=signature" json:"signature,omitempty"`
}

func (m *Uits) Reset()         { *m = Uits{} }
func (m *Uits) String() string { return proto.CompactTextString(m) }
func (*Uits) ProtoMessage()    {}

func (m *Uits) GetMetadata() *UitsMetadata {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *Uits) GetSignature() *UitsSignature {
	if m != nil {
		return m.Signature
	}
	return nil
}

type UploadedUitsId3Tag struct {
	Owner string `protobuf:"bytes,1,opt,name=owner" json:"owner,omitempty"`
	Data  []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *UploadedUitsId3Tag) Reset()         { *m = UploadedUitsId3Tag{} }
func (m *UploadedUitsId3Tag) String() string { return proto.CompactTextString(m) }
func (*UploadedUitsId3Tag) ProtoMessage()    {}

func init() {
	proto.RegisterEnum("google.musicmanager.v0.ProductId_Type", ProductId_Type_name, ProductId_Type_value)
	proto.RegisterEnum("google.musicmanager.v0.AssetId_Type", AssetId_Type_name, AssetId_Type_value)
	proto.RegisterEnum("google.musicmanager.v0.MediaId_AlgorithmType", MediaId_AlgorithmType_name, MediaId_AlgorithmType_value)
	proto.RegisterEnum("google.musicmanager.v0.UrlInfo_Type", UrlInfo_Type_name, UrlInfo_Type_value)
	proto.RegisterEnum("google.musicmanager.v0.CopyrightStatus_Type", CopyrightStatus_Type_name, CopyrightStatus_Type_value)
	proto.RegisterEnum("google.musicmanager.v0.UitsMetadata_ParentalAdvisoryType", UitsMetadata_ParentalAdvisoryType_name, UitsMetadata_ParentalAdvisoryType_value)
	proto.RegisterEnum("google.musicmanager.v0.UitsSignature_AlgorithmType", UitsSignature_AlgorithmType_name, UitsSignature_AlgorithmType_value)
	proto.RegisterEnum("google.musicmanager.v0.UitsSignature_CanonicalizationType", UitsSignature_CanonicalizationType_name, UitsSignature_CanonicalizationType_value)
}
