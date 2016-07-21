// Code generated by protoc-gen-go.
// source: profile.proto
// DO NOT EDIT!

/*
Package profile is a generated protocol buffer package.

It is generated from these files:
	profile.proto

It has these top-level messages:
	Profile
*/
package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
const _ = proto.ProtoPackageIsVersion1

type Profile struct {
	Handle           string                   `protobuf:"bytes,1,opt,name=handle" json:"handle,omitempty"`
	Name             string                   `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Location         string                   `protobuf:"bytes,3,opt,name=location" json:"location,omitempty"`
	About            string                   `protobuf:"bytes,4,opt,name=about" json:"about,omitempty"`
	ShortDescription string                   `protobuf:"bytes,5,opt,name=shortDescription" json:"shortDescription,omitempty"`
	Website          string                   `protobuf:"bytes,6,opt,name=website" json:"website,omitempty"`
	Email            string                   `protobuf:"bytes,7,opt,name=email" json:"email,omitempty"`
	Social           []*Profile_SocialAccount `protobuf:"bytes,8,rep,name=social" json:"social,omitempty"`
	Nsfw             bool                     `protobuf:"varint,9,opt,name=nsfw" json:"nsfw,omitempty"`
	Vendor           bool                     `protobuf:"varint,10,opt,name=vendor" json:"vendor,omitempty"`
	Moderator        bool                     `protobuf:"varint,11,opt,name=moderator" json:"moderator,omitempty"`
	PrimaryColor     string                   `protobuf:"bytes,12,opt,name=primaryColor" json:"primaryColor,omitempty"`
	SecondaryColor   string                   `protobuf:"bytes,13,opt,name=secondaryColor" json:"secondaryColor,omitempty"`
	TextColor        string                   `protobuf:"bytes,14,opt,name=textColor" json:"textColor,omitempty"`
	FollowerCount    uint32                   `protobuf:"varint,15,opt,name=followerCount" json:"followerCount,omitempty"`
	FollowingCount   uint32                   `protobuf:"varint,16,opt,name=followingCount" json:"followingCount,omitempty"`
	ListingCount     uint32                   `protobuf:"varint,17,opt,name=listingCount" json:"listingCount,omitempty"`
	LastModified     uint64                   `protobuf:"varint,18,opt,name=last_modified,json=lastModified" json:"last_modified,omitempty"`
}

func (m *Profile) Reset()                    { *m = Profile{} }
func (m *Profile) String() string            { return proto.CompactTextString(m) }
func (*Profile) ProtoMessage()               {}
func (*Profile) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{0} }

func (m *Profile) GetSocial() []*Profile_SocialAccount {
	if m != nil {
		return m.Social
	}
	return nil
}

type Profile_SocialAccount struct {
	Type     string `protobuf:"bytes,1,opt,name=type" json:"type,omitempty"`
	Username string `protobuf:"bytes,2,opt,name=username" json:"username,omitempty"`
	Proof    string `protobuf:"bytes,3,opt,name=proof" json:"proof,omitempty"`
}

func (m *Profile_SocialAccount) Reset()                    { *m = Profile_SocialAccount{} }
func (m *Profile_SocialAccount) String() string            { return proto.CompactTextString(m) }
func (*Profile_SocialAccount) ProtoMessage()               {}
func (*Profile_SocialAccount) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{0, 0} }

func init() {
	proto.RegisterType((*Profile)(nil), "Profile")
	proto.RegisterType((*Profile_SocialAccount)(nil), "Profile.SocialAccount")
}

var fileDescriptor3 = []byte{
	// 376 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x64, 0x92, 0xc1, 0x6e, 0xe2, 0x30,
	0x10, 0x86, 0xc5, 0x02, 0x01, 0x06, 0xc2, 0xb2, 0xd6, 0x0a, 0x59, 0x68, 0x0f, 0x88, 0xad, 0x2a,
	0xd4, 0x43, 0x0e, 0xed, 0x13, 0x54, 0xf4, 0x5a, 0xa9, 0x4a, 0xd5, 0x73, 0x65, 0x12, 0xa7, 0x58,
	0x72, 0xe2, 0xc8, 0x36, 0xa5, 0x3c, 0x48, 0xdf, 0xb7, 0xf6, 0x24, 0x84, 0x86, 0x9e, 0xf0, 0xff,
	0xf9, 0x63, 0xec, 0x71, 0x06, 0xc2, 0x52, 0xab, 0x4c, 0x48, 0x1e, 0xb9, 0x5f, 0xab, 0x56, 0x9f,
	0x7d, 0x18, 0x3c, 0x55, 0x84, 0xcc, 0x21, 0xd8, 0xb1, 0x22, 0x95, 0x9c, 0x76, 0x96, 0x9d, 0xf5,
	0x28, 0xae, 0x13, 0x21, 0xd0, 0x2b, 0x58, 0xce, 0xe9, 0x2f, 0xa4, 0xb8, 0x26, 0x0b, 0x18, 0x4a,
	0x95, 0x30, 0x2b, 0x54, 0x41, 0xbb, 0xc8, 0x9b, 0x4c, 0xfe, 0x42, 0x9f, 0x6d, 0xd5, 0xde, 0xd2,
	0x1e, 0x6e, 0x54, 0x81, 0xdc, 0xc0, 0xcc, 0xec, 0x94, 0xb6, 0x0f, 0xdc, 0x24, 0x5a, 0x94, 0xf8,
	0xcf, 0x3e, 0x0a, 0x3f, 0x38, 0xa1, 0x30, 0x38, 0xf0, 0xad, 0x11, 0x96, 0xd3, 0x00, 0x95, 0x53,
	0xf4, 0xb5, 0x79, 0xce, 0x84, 0xa4, 0x83, 0xaa, 0x36, 0x06, 0x12, 0x41, 0x60, 0x54, 0x22, 0x98,
	0xa4, 0xc3, 0x65, 0x77, 0x3d, 0xbe, 0x9d, 0x47, 0x75, 0x4f, 0xd1, 0x33, 0xe2, 0xfb, 0x24, 0x51,
	0xfb, 0xc2, 0xc6, 0xb5, 0x85, 0x1d, 0x99, 0xec, 0x40, 0x47, 0xae, 0xc8, 0x30, 0xc6, 0xb5, 0xef,
	0xfe, 0x9d, 0x17, 0xa9, 0xd2, 0x14, 0x90, 0xd6, 0x89, 0xfc, 0x83, 0x51, 0xae, 0x52, 0xae, 0x99,
	0x75, 0x5b, 0x63, 0xdc, 0x3a, 0x03, 0xb2, 0x82, 0x49, 0xa9, 0x45, 0xce, 0xf4, 0x71, 0xa3, 0xa4,
	0x13, 0x26, 0x78, 0xad, 0x16, 0x23, 0xd7, 0x30, 0x35, 0x3c, 0x51, 0x45, 0xda, 0x58, 0x21, 0x5a,
	0x17, 0xd4, 0x9f, 0x64, 0xf9, 0x87, 0xad, 0x94, 0x29, 0x2a, 0x67, 0x40, 0xae, 0x20, 0xcc, 0x94,
	0x94, 0xea, 0xc0, 0xf5, 0xc6, 0x37, 0x43, 0x7f, 0x3b, 0x23, 0x8c, 0xdb, 0xd0, 0x9f, 0x55, 0x01,
	0x51, 0xbc, 0x55, 0xda, 0x0c, 0xb5, 0x0b, 0xea, 0xef, 0x2d, 0x85, 0xb1, 0x8d, 0xf5, 0x07, 0xad,
	0x16, 0x23, 0xff, 0x21, 0x94, 0xcc, 0xd8, 0x57, 0xd7, 0xad, 0xc8, 0x04, 0x4f, 0x29, 0x71, 0x52,
	0xcf, 0x49, 0x0e, 0x3e, 0xd6, 0x6c, 0xf1, 0x02, 0x61, 0xeb, 0x8d, 0xfd, 0xdb, 0xda, 0x63, 0x79,
	0x9a, 0x21, 0x5c, 0xfb, 0x69, 0xd9, 0x1b, 0xae, 0xbf, 0x4d, 0x51, 0x93, 0xfd, 0x17, 0x75, 0xa3,
	0xa8, 0xb2, 0x7a, 0x8c, 0xaa, 0xb0, 0x0d, 0x70, 0x3c, 0xef, 0xbe, 0x02, 0x00, 0x00, 0xff, 0xff,
	0x0c, 0xff, 0xca, 0x1c, 0xaf, 0x02, 0x00, 0x00,
}