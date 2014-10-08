// Code generated by protoc-gen-gogo.
// source: bind_volume.proto
// DO NOT EDIT!

package garden

import proto "code.google.com/p/gogoprotobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type BindVolumeRequest_Mode int32

const (
	BindVolumeRequest_RO BindVolumeRequest_Mode = 0
	BindVolumeRequest_RW BindVolumeRequest_Mode = 1
)

var BindVolumeRequest_Mode_name = map[int32]string{
	0: "RO",
	1: "RW",
}
var BindVolumeRequest_Mode_value = map[string]int32{
	"RO": 0,
	"RW": 1,
}

func (x BindVolumeRequest_Mode) Enum() *BindVolumeRequest_Mode {
	p := new(BindVolumeRequest_Mode)
	*p = x
	return p
}
func (x BindVolumeRequest_Mode) String() string {
	return proto.EnumName(BindVolumeRequest_Mode_name, int32(x))
}
func (x *BindVolumeRequest_Mode) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(BindVolumeRequest_Mode_value, data, "BindVolumeRequest_Mode")
	if err != nil {
		return err
	}
	*x = BindVolumeRequest_Mode(value)
	return nil
}

type BindVolumeRequest struct {
	ContainerHandle  *string                 `protobuf:"bytes,1,req,name=container_handle" json:"container_handle,omitempty"`
	VolumeHandle     *string                 `protobuf:"bytes,2,req,name=volume_handle" json:"volume_handle,omitempty"`
	DestinationPath  *string                 `protobuf:"bytes,3,req,name=destination_path" json:"destination_path,omitempty"`
	Mode             *BindVolumeRequest_Mode `protobuf:"varint,4,req,name=mode,enum=garden.BindVolumeRequest_Mode" json:"mode,omitempty"`
	XXX_unrecognized []byte                  `json:"-"`
}

func (m *BindVolumeRequest) Reset()         { *m = BindVolumeRequest{} }
func (m *BindVolumeRequest) String() string { return proto.CompactTextString(m) }
func (*BindVolumeRequest) ProtoMessage()    {}

func (m *BindVolumeRequest) GetContainerHandle() string {
	if m != nil && m.ContainerHandle != nil {
		return *m.ContainerHandle
	}
	return ""
}

func (m *BindVolumeRequest) GetVolumeHandle() string {
	if m != nil && m.VolumeHandle != nil {
		return *m.VolumeHandle
	}
	return ""
}

func (m *BindVolumeRequest) GetDestinationPath() string {
	if m != nil && m.DestinationPath != nil {
		return *m.DestinationPath
	}
	return ""
}

func (m *BindVolumeRequest) GetMode() BindVolumeRequest_Mode {
	if m != nil && m.Mode != nil {
		return *m.Mode
	}
	return BindVolumeRequest_RO
}

type BindVolumeResponse struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *BindVolumeResponse) Reset()         { *m = BindVolumeResponse{} }
func (m *BindVolumeResponse) String() string { return proto.CompactTextString(m) }
func (*BindVolumeResponse) ProtoMessage()    {}

func init() {
	proto.RegisterEnum("garden.BindVolumeRequest_Mode", BindVolumeRequest_Mode_name, BindVolumeRequest_Mode_value)
}