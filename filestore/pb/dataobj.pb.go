// Code generated by protoc-gen-gogo.
// source: dataobj.proto
// DO NOT EDIT!

/*
Package datastore_pb is a generated protocol buffer package.

It is generated from these files:
	dataobj.proto

It has these top-level messages:
	DataObj
*/
package datastore_pb

import proto "unknown/gogo-protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type DataObj struct {
	FilePath         *string `protobuf:"bytes,1,opt,name=FilePath" json:"FilePath,omitempty"`
	Offset           *uint64 `protobuf:"varint,2,opt,name=Offset" json:"Offset,omitempty"`
	Size_            *uint64 `protobuf:"varint,3,opt,name=Size" json:"Size,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *DataObj) Reset()         { *m = DataObj{} }
func (m *DataObj) String() string { return proto.CompactTextString(m) }
func (*DataObj) ProtoMessage()    {}

func (m *DataObj) GetFilePath() string {
	if m != nil && m.FilePath != nil {
		return *m.FilePath
	}
	return ""
}

func (m *DataObj) GetOffset() uint64 {
	if m != nil && m.Offset != nil {
		return *m.Offset
	}
	return 0
}

func (m *DataObj) GetSize_() uint64 {
	if m != nil && m.Size_ != nil {
		return *m.Size_
	}
	return 0
}

func init() {
	proto.RegisterType((*DataObj)(nil), "datastore.pb.DataObj")
}
