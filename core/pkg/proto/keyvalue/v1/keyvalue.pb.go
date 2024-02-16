// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.4
// source: nitric/proto/keyvalue/v1/keyvalue.proto

package KeyValuepb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	structpb "google.golang.org/protobuf/types/known/structpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Provides a Key/Value Store
type Store struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The store name
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *Store) Reset() {
	*x = Store{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nitric_proto_keyvalue_v1_keyvalue_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Store) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Store) ProtoMessage() {}

func (x *Store) ProtoReflect() protoreflect.Message {
	mi := &file_nitric_proto_keyvalue_v1_keyvalue_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Store.ProtoReflect.Descriptor instead.
func (*Store) Descriptor() ([]byte, []int) {
	return file_nitric_proto_keyvalue_v1_keyvalue_proto_rawDescGZIP(), []int{0}
}

func (x *Store) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// ValueRef provides a unique identifier for a value
type ValueRef struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The key/value store name
	Store string `protobuf:"bytes,1,opt,name=store,proto3" json:"store,omitempty"`
	// The item's unique key within the store
	Key string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *ValueRef) Reset() {
	*x = ValueRef{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nitric_proto_keyvalue_v1_keyvalue_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ValueRef) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValueRef) ProtoMessage() {}

func (x *ValueRef) ProtoReflect() protoreflect.Message {
	mi := &file_nitric_proto_keyvalue_v1_keyvalue_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValueRef.ProtoReflect.Descriptor instead.
func (*ValueRef) Descriptor() ([]byte, []int) {
	return file_nitric_proto_keyvalue_v1_keyvalue_proto_rawDescGZIP(), []int{1}
}

func (x *ValueRef) GetStore() string {
	if x != nil {
		return x.Store
	}
	return ""
}

func (x *ValueRef) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

// Value provides a return value type
type Value struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// ValueRef of the key/value pair, which includes the store and key
	Ref *ValueRef `protobuf:"bytes,1,opt,name=ref,proto3" json:"ref,omitempty"`
	// The content (JSON object)
	Content *structpb.Struct `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *Value) Reset() {
	*x = Value{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nitric_proto_keyvalue_v1_keyvalue_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Value) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Value) ProtoMessage() {}

func (x *Value) ProtoReflect() protoreflect.Message {
	mi := &file_nitric_proto_keyvalue_v1_keyvalue_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Value.ProtoReflect.Descriptor instead.
func (*Value) Descriptor() ([]byte, []int) {
	return file_nitric_proto_keyvalue_v1_keyvalue_proto_rawDescGZIP(), []int{2}
}

func (x *Value) GetRef() *ValueRef {
	if x != nil {
		return x.Ref
	}
	return nil
}

func (x *Value) GetContent() *structpb.Struct {
	if x != nil {
		return x.Content
	}
	return nil
}

type KeyValueGetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// ValueRef of the key/value pair to get, which includes the store and key
	Ref *ValueRef `protobuf:"bytes,1,opt,name=ref,proto3" json:"ref,omitempty"`
}

func (x *KeyValueGetRequest) Reset() {
	*x = KeyValueGetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nitric_proto_keyvalue_v1_keyvalue_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KeyValueGetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KeyValueGetRequest) ProtoMessage() {}

func (x *KeyValueGetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_nitric_proto_keyvalue_v1_keyvalue_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KeyValueGetRequest.ProtoReflect.Descriptor instead.
func (*KeyValueGetRequest) Descriptor() ([]byte, []int) {
	return file_nitric_proto_keyvalue_v1_keyvalue_proto_rawDescGZIP(), []int{3}
}

func (x *KeyValueGetRequest) GetRef() *ValueRef {
	if x != nil {
		return x.Ref
	}
	return nil
}

type KeyValueGetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The retrieved value
	Value *Value `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *KeyValueGetResponse) Reset() {
	*x = KeyValueGetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nitric_proto_keyvalue_v1_keyvalue_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KeyValueGetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KeyValueGetResponse) ProtoMessage() {}

func (x *KeyValueGetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_nitric_proto_keyvalue_v1_keyvalue_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KeyValueGetResponse.ProtoReflect.Descriptor instead.
func (*KeyValueGetResponse) Descriptor() ([]byte, []int) {
	return file_nitric_proto_keyvalue_v1_keyvalue_proto_rawDescGZIP(), []int{4}
}

func (x *KeyValueGetResponse) GetValue() *Value {
	if x != nil {
		return x.Value
	}
	return nil
}

type KeyValueSetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// ValueRef of the key/value pair to set, which includes the store and key
	Ref *ValueRef `protobuf:"bytes,1,opt,name=ref,proto3" json:"ref,omitempty"`
	// The value content to store (JSON object)
	Content *structpb.Struct `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *KeyValueSetRequest) Reset() {
	*x = KeyValueSetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nitric_proto_keyvalue_v1_keyvalue_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KeyValueSetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KeyValueSetRequest) ProtoMessage() {}

func (x *KeyValueSetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_nitric_proto_keyvalue_v1_keyvalue_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KeyValueSetRequest.ProtoReflect.Descriptor instead.
func (*KeyValueSetRequest) Descriptor() ([]byte, []int) {
	return file_nitric_proto_keyvalue_v1_keyvalue_proto_rawDescGZIP(), []int{5}
}

func (x *KeyValueSetRequest) GetRef() *ValueRef {
	if x != nil {
		return x.Ref
	}
	return nil
}

func (x *KeyValueSetRequest) GetContent() *structpb.Struct {
	if x != nil {
		return x.Content
	}
	return nil
}

type KeyValueSetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *KeyValueSetResponse) Reset() {
	*x = KeyValueSetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nitric_proto_keyvalue_v1_keyvalue_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KeyValueSetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KeyValueSetResponse) ProtoMessage() {}

func (x *KeyValueSetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_nitric_proto_keyvalue_v1_keyvalue_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KeyValueSetResponse.ProtoReflect.Descriptor instead.
func (*KeyValueSetResponse) Descriptor() ([]byte, []int) {
	return file_nitric_proto_keyvalue_v1_keyvalue_proto_rawDescGZIP(), []int{6}
}

type KeyValueDeleteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// ValueRef of the key/value pair to delete, which includes the store and key
	Ref *ValueRef `protobuf:"bytes,1,opt,name=ref,proto3" json:"ref,omitempty"`
}

func (x *KeyValueDeleteRequest) Reset() {
	*x = KeyValueDeleteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nitric_proto_keyvalue_v1_keyvalue_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KeyValueDeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KeyValueDeleteRequest) ProtoMessage() {}

func (x *KeyValueDeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_nitric_proto_keyvalue_v1_keyvalue_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KeyValueDeleteRequest.ProtoReflect.Descriptor instead.
func (*KeyValueDeleteRequest) Descriptor() ([]byte, []int) {
	return file_nitric_proto_keyvalue_v1_keyvalue_proto_rawDescGZIP(), []int{7}
}

func (x *KeyValueDeleteRequest) GetRef() *ValueRef {
	if x != nil {
		return x.Ref
	}
	return nil
}

type KeyValueDeleteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *KeyValueDeleteResponse) Reset() {
	*x = KeyValueDeleteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nitric_proto_keyvalue_v1_keyvalue_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KeyValueDeleteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KeyValueDeleteResponse) ProtoMessage() {}

func (x *KeyValueDeleteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_nitric_proto_keyvalue_v1_keyvalue_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KeyValueDeleteResponse.ProtoReflect.Descriptor instead.
func (*KeyValueDeleteResponse) Descriptor() ([]byte, []int) {
	return file_nitric_proto_keyvalue_v1_keyvalue_proto_rawDescGZIP(), []int{8}
}

var File_nitric_proto_keyvalue_v1_keyvalue_proto protoreflect.FileDescriptor

var file_nitric_proto_keyvalue_v1_keyvalue_proto_rawDesc = []byte{
	0x0a, 0x27, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6b,
	0x65, 0x79, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x6b, 0x65, 0x79, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x18, 0x6e, 0x69, 0x74, 0x72, 0x69,
	0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x1b, 0x0a, 0x05, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x32,
	0x0a, 0x08, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x65, 0x66, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74,
	0x6f, 0x72, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x74, 0x6f, 0x72, 0x65,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x22, 0x70, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x34, 0x0a, 0x03, 0x72,
	0x65, 0x66, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69,
	0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x2e, 0x76, 0x31, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x65, 0x66, 0x52, 0x03, 0x72, 0x65,
	0x66, 0x12, 0x31, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x07, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x22, 0x4a, 0x0a, 0x12, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x34, 0x0a, 0x03, 0x72, 0x65,
	0x66, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x65, 0x66, 0x52, 0x03, 0x72, 0x65, 0x66,
	0x22, 0x4c, 0x0a, 0x13, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x47, 0x65, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x35, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x2e, 0x76,
	0x31, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x7d,
	0x0a, 0x12, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x53, 0x65, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x34, 0x0a, 0x03, 0x72, 0x65, 0x66, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x22, 0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x52, 0x65, 0x66, 0x52, 0x03, 0x72, 0x65, 0x66, 0x12, 0x31, 0x0a, 0x07, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74,
	0x72, 0x75, 0x63, 0x74, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x15, 0x0a,
	0x13, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x53, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x4d, 0x0a, 0x15, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x34, 0x0a,
	0x03, 0x72, 0x65, 0x66, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x6e, 0x69, 0x74,
	0x72, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x65, 0x66, 0x52, 0x03,
	0x72, 0x65, 0x66, 0x22, 0x18, 0x0a, 0x16, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0xbf, 0x02,
	0x0a, 0x08, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x62, 0x0a, 0x03, 0x47, 0x65,
	0x74, 0x12, 0x2c, 0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4b, 0x65, 0x79,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x2d, 0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4b,
	0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4b, 0x65, 0x79, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x62,
	0x0a, 0x03, 0x53, 0x65, 0x74, 0x12, 0x2c, 0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x2e, 0x76, 0x31,
	0x2e, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x53, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x2d, 0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4b,
	0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x53, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x6b, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x2f, 0x2e, 0x6e,
	0x69, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4b, 0x65, 0x79, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x30, 0x2e,
	0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4b, 0x65, 0x79,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42,
	0xa3, 0x01, 0x0a, 0x1b, 0x69, 0x6f, 0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x6b, 0x65, 0x79, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x2e, 0x76, 0x31, 0x42,
	0x08, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x50, 0x01, 0x5a, 0x42, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x74, 0x65,
	0x63, 0x68, 0x2f, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x70,
	0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6b, 0x65, 0x79, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x2f, 0x76, 0x31, 0x3b, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x70, 0x62, 0xaa,
	0x02, 0x18, 0x4e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4b,
	0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x2e, 0x76, 0x31, 0xca, 0x02, 0x18, 0x4e, 0x69, 0x74,
	0x72, 0x69, 0x63, 0x5c, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x5c, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x5c, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_nitric_proto_keyvalue_v1_keyvalue_proto_rawDescOnce sync.Once
	file_nitric_proto_keyvalue_v1_keyvalue_proto_rawDescData = file_nitric_proto_keyvalue_v1_keyvalue_proto_rawDesc
)

func file_nitric_proto_keyvalue_v1_keyvalue_proto_rawDescGZIP() []byte {
	file_nitric_proto_keyvalue_v1_keyvalue_proto_rawDescOnce.Do(func() {
		file_nitric_proto_keyvalue_v1_keyvalue_proto_rawDescData = protoimpl.X.CompressGZIP(file_nitric_proto_keyvalue_v1_keyvalue_proto_rawDescData)
	})
	return file_nitric_proto_keyvalue_v1_keyvalue_proto_rawDescData
}

var file_nitric_proto_keyvalue_v1_keyvalue_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_nitric_proto_keyvalue_v1_keyvalue_proto_goTypes = []interface{}{
	(*Store)(nil),                  // 0: nitric.proto.KeyValue.v1.Store
	(*ValueRef)(nil),               // 1: nitric.proto.KeyValue.v1.ValueRef
	(*Value)(nil),                  // 2: nitric.proto.KeyValue.v1.Value
	(*KeyValueGetRequest)(nil),     // 3: nitric.proto.KeyValue.v1.KeyValueGetRequest
	(*KeyValueGetResponse)(nil),    // 4: nitric.proto.KeyValue.v1.KeyValueGetResponse
	(*KeyValueSetRequest)(nil),     // 5: nitric.proto.KeyValue.v1.KeyValueSetRequest
	(*KeyValueSetResponse)(nil),    // 6: nitric.proto.KeyValue.v1.KeyValueSetResponse
	(*KeyValueDeleteRequest)(nil),  // 7: nitric.proto.KeyValue.v1.KeyValueDeleteRequest
	(*KeyValueDeleteResponse)(nil), // 8: nitric.proto.KeyValue.v1.KeyValueDeleteResponse
	(*structpb.Struct)(nil),        // 9: google.protobuf.Struct
}
var file_nitric_proto_keyvalue_v1_keyvalue_proto_depIdxs = []int32{
	1,  // 0: nitric.proto.KeyValue.v1.Value.ref:type_name -> nitric.proto.KeyValue.v1.ValueRef
	9,  // 1: nitric.proto.KeyValue.v1.Value.content:type_name -> google.protobuf.Struct
	1,  // 2: nitric.proto.KeyValue.v1.KeyValueGetRequest.ref:type_name -> nitric.proto.KeyValue.v1.ValueRef
	2,  // 3: nitric.proto.KeyValue.v1.KeyValueGetResponse.value:type_name -> nitric.proto.KeyValue.v1.Value
	1,  // 4: nitric.proto.KeyValue.v1.KeyValueSetRequest.ref:type_name -> nitric.proto.KeyValue.v1.ValueRef
	9,  // 5: nitric.proto.KeyValue.v1.KeyValueSetRequest.content:type_name -> google.protobuf.Struct
	1,  // 6: nitric.proto.KeyValue.v1.KeyValueDeleteRequest.ref:type_name -> nitric.proto.KeyValue.v1.ValueRef
	3,  // 7: nitric.proto.KeyValue.v1.KeyValue.Get:input_type -> nitric.proto.KeyValue.v1.KeyValueGetRequest
	5,  // 8: nitric.proto.KeyValue.v1.KeyValue.Set:input_type -> nitric.proto.KeyValue.v1.KeyValueSetRequest
	7,  // 9: nitric.proto.KeyValue.v1.KeyValue.Delete:input_type -> nitric.proto.KeyValue.v1.KeyValueDeleteRequest
	4,  // 10: nitric.proto.KeyValue.v1.KeyValue.Get:output_type -> nitric.proto.KeyValue.v1.KeyValueGetResponse
	6,  // 11: nitric.proto.KeyValue.v1.KeyValue.Set:output_type -> nitric.proto.KeyValue.v1.KeyValueSetResponse
	8,  // 12: nitric.proto.KeyValue.v1.KeyValue.Delete:output_type -> nitric.proto.KeyValue.v1.KeyValueDeleteResponse
	10, // [10:13] is the sub-list for method output_type
	7,  // [7:10] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_nitric_proto_keyvalue_v1_keyvalue_proto_init() }
func file_nitric_proto_keyvalue_v1_keyvalue_proto_init() {
	if File_nitric_proto_keyvalue_v1_keyvalue_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_nitric_proto_keyvalue_v1_keyvalue_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Store); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_nitric_proto_keyvalue_v1_keyvalue_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ValueRef); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_nitric_proto_keyvalue_v1_keyvalue_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Value); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_nitric_proto_keyvalue_v1_keyvalue_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KeyValueGetRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_nitric_proto_keyvalue_v1_keyvalue_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KeyValueGetResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_nitric_proto_keyvalue_v1_keyvalue_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KeyValueSetRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_nitric_proto_keyvalue_v1_keyvalue_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KeyValueSetResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_nitric_proto_keyvalue_v1_keyvalue_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KeyValueDeleteRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_nitric_proto_keyvalue_v1_keyvalue_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KeyValueDeleteResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_nitric_proto_keyvalue_v1_keyvalue_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_nitric_proto_keyvalue_v1_keyvalue_proto_goTypes,
		DependencyIndexes: file_nitric_proto_keyvalue_v1_keyvalue_proto_depIdxs,
		MessageInfos:      file_nitric_proto_keyvalue_v1_keyvalue_proto_msgTypes,
	}.Build()
	File_nitric_proto_keyvalue_v1_keyvalue_proto = out.File
	file_nitric_proto_keyvalue_v1_keyvalue_proto_rawDesc = nil
	file_nitric_proto_keyvalue_v1_keyvalue_proto_goTypes = nil
	file_nitric_proto_keyvalue_v1_keyvalue_proto_depIdxs = nil
}
