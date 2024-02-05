// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.4
// source: nitric/proto/queue/v1/queue.proto

package queuepb

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

type QueueSendRequestBatch struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The Nitric name for the queue
	// this will automatically be resolved to the provider specific queue identifier.
	QueueName string `protobuf:"bytes,1,opt,name=queue_name,json=queueName,proto3" json:"queue_name,omitempty"`
	// Array of tasks to push to the queue
	Requests []*QueueSendRequest `protobuf:"bytes,2,rep,name=requests,proto3" json:"requests,omitempty"`
}

func (x *QueueSendRequestBatch) Reset() {
	*x = QueueSendRequestBatch{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nitric_proto_queue_v1_queue_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueueSendRequestBatch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueueSendRequestBatch) ProtoMessage() {}

func (x *QueueSendRequestBatch) ProtoReflect() protoreflect.Message {
	mi := &file_nitric_proto_queue_v1_queue_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueueSendRequestBatch.ProtoReflect.Descriptor instead.
func (*QueueSendRequestBatch) Descriptor() ([]byte, []int) {
	return file_nitric_proto_queue_v1_queue_proto_rawDescGZIP(), []int{0}
}

func (x *QueueSendRequestBatch) GetQueueName() string {
	if x != nil {
		return x.QueueName
	}
	return ""
}

func (x *QueueSendRequestBatch) GetRequests() []*QueueSendRequest {
	if x != nil {
		return x.Requests
	}
	return nil
}

// Response for sending a collection of tasks
type QueueSendResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A list of tasks that failed to be queued
	FailedRequests []*FailedSendRequest `protobuf:"bytes,1,rep,name=failed_requests,json=failedRequests,proto3" json:"failed_requests,omitempty"`
}

func (x *QueueSendResponse) Reset() {
	*x = QueueSendResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nitric_proto_queue_v1_queue_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueueSendResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueueSendResponse) ProtoMessage() {}

func (x *QueueSendResponse) ProtoReflect() protoreflect.Message {
	mi := &file_nitric_proto_queue_v1_queue_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueueSendResponse.ProtoReflect.Descriptor instead.
func (*QueueSendResponse) Descriptor() ([]byte, []int) {
	return file_nitric_proto_queue_v1_queue_proto_rawDescGZIP(), []int{1}
}

func (x *QueueSendResponse) GetFailedRequests() []*FailedSendRequest {
	if x != nil {
		return x.FailedRequests
	}
	return nil
}

type QueueReceiveRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The nitric name for the queue
	// this will automatically be resolved to the provider specific queue identifier.
	QueueName string `protobuf:"bytes,1,opt,name=queue_name,json=queueName,proto3" json:"queue_name,omitempty"`
	// The max number of items to pop off the queue, may be capped by provider specific limitations
	Depth int32 `protobuf:"varint,2,opt,name=depth,proto3" json:"depth,omitempty"`
}

func (x *QueueReceiveRequest) Reset() {
	*x = QueueReceiveRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nitric_proto_queue_v1_queue_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueueReceiveRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueueReceiveRequest) ProtoMessage() {}

func (x *QueueReceiveRequest) ProtoReflect() protoreflect.Message {
	mi := &file_nitric_proto_queue_v1_queue_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueueReceiveRequest.ProtoReflect.Descriptor instead.
func (*QueueReceiveRequest) Descriptor() ([]byte, []int) {
	return file_nitric_proto_queue_v1_queue_proto_rawDescGZIP(), []int{2}
}

func (x *QueueReceiveRequest) GetQueueName() string {
	if x != nil {
		return x.QueueName
	}
	return ""
}

func (x *QueueReceiveRequest) GetDepth() int32 {
	if x != nil {
		return x.Depth
	}
	return 0
}

type QueueReceiveResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Array of tasks popped off the queue
	Tasks []*ReceivedTask `protobuf:"bytes,1,rep,name=tasks,proto3" json:"tasks,omitempty"`
}

func (x *QueueReceiveResponse) Reset() {
	*x = QueueReceiveResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nitric_proto_queue_v1_queue_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueueReceiveResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueueReceiveResponse) ProtoMessage() {}

func (x *QueueReceiveResponse) ProtoReflect() protoreflect.Message {
	mi := &file_nitric_proto_queue_v1_queue_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueueReceiveResponse.ProtoReflect.Descriptor instead.
func (*QueueReceiveResponse) Descriptor() ([]byte, []int) {
	return file_nitric_proto_queue_v1_queue_proto_rawDescGZIP(), []int{3}
}

func (x *QueueReceiveResponse) GetTasks() []*ReceivedTask {
	if x != nil {
		return x.Tasks
	}
	return nil
}

type QueueCompleteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The nitric name for the queue
	//
	//	this will automatically be resolved to the provider specific queue identifier.
	QueueName string `protobuf:"bytes,1,opt,name=queue_name,json=queueName,proto3" json:"queue_name,omitempty"`
	// Lease id of the task to be completed
	LeaseId string `protobuf:"bytes,2,opt,name=lease_id,json=leaseId,proto3" json:"lease_id,omitempty"`
}

func (x *QueueCompleteRequest) Reset() {
	*x = QueueCompleteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nitric_proto_queue_v1_queue_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueueCompleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueueCompleteRequest) ProtoMessage() {}

func (x *QueueCompleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_nitric_proto_queue_v1_queue_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueueCompleteRequest.ProtoReflect.Descriptor instead.
func (*QueueCompleteRequest) Descriptor() ([]byte, []int) {
	return file_nitric_proto_queue_v1_queue_proto_rawDescGZIP(), []int{4}
}

func (x *QueueCompleteRequest) GetQueueName() string {
	if x != nil {
		return x.QueueName
	}
	return ""
}

func (x *QueueCompleteRequest) GetLeaseId() string {
	if x != nil {
		return x.LeaseId
	}
	return ""
}

type QueueCompleteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *QueueCompleteResponse) Reset() {
	*x = QueueCompleteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nitric_proto_queue_v1_queue_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueueCompleteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueueCompleteResponse) ProtoMessage() {}

func (x *QueueCompleteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_nitric_proto_queue_v1_queue_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueueCompleteResponse.ProtoReflect.Descriptor instead.
func (*QueueCompleteResponse) Descriptor() ([]byte, []int) {
	return file_nitric_proto_queue_v1_queue_proto_rawDescGZIP(), []int{5}
}

// A task to be sent to a queue.
type QueueSendRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// The payload of the task
	Payload *structpb.Struct `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *QueueSendRequest) Reset() {
	*x = QueueSendRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nitric_proto_queue_v1_queue_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueueSendRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueueSendRequest) ProtoMessage() {}

func (x *QueueSendRequest) ProtoReflect() protoreflect.Message {
	mi := &file_nitric_proto_queue_v1_queue_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueueSendRequest.ProtoReflect.Descriptor instead.
func (*QueueSendRequest) Descriptor() ([]byte, []int) {
	return file_nitric_proto_queue_v1_queue_proto_rawDescGZIP(), []int{6}
}

func (x *QueueSendRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *QueueSendRequest) GetPayload() *structpb.Struct {
	if x != nil {
		return x.Payload
	}
	return nil
}

type ReceivedTask struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LeaseId string           `protobuf:"bytes,1,opt,name=lease_id,json=leaseId,proto3" json:"lease_id,omitempty"`
	Payload *structpb.Struct `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *ReceivedTask) Reset() {
	*x = ReceivedTask{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nitric_proto_queue_v1_queue_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReceivedTask) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReceivedTask) ProtoMessage() {}

func (x *ReceivedTask) ProtoReflect() protoreflect.Message {
	mi := &file_nitric_proto_queue_v1_queue_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReceivedTask.ProtoReflect.Descriptor instead.
func (*ReceivedTask) Descriptor() ([]byte, []int) {
	return file_nitric_proto_queue_v1_queue_proto_rawDescGZIP(), []int{7}
}

func (x *ReceivedTask) GetLeaseId() string {
	if x != nil {
		return x.LeaseId
	}
	return ""
}

func (x *ReceivedTask) GetPayload() *structpb.Struct {
	if x != nil {
		return x.Payload
	}
	return nil
}

type FailedSendRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The task that failed to be pushed
	Request *QueueSendRequest `protobuf:"bytes,1,opt,name=request,proto3" json:"request,omitempty"`
	// A message describing the failure
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *FailedSendRequest) Reset() {
	*x = FailedSendRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nitric_proto_queue_v1_queue_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FailedSendRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FailedSendRequest) ProtoMessage() {}

func (x *FailedSendRequest) ProtoReflect() protoreflect.Message {
	mi := &file_nitric_proto_queue_v1_queue_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FailedSendRequest.ProtoReflect.Descriptor instead.
func (*FailedSendRequest) Descriptor() ([]byte, []int) {
	return file_nitric_proto_queue_v1_queue_proto_rawDescGZIP(), []int{8}
}

func (x *FailedSendRequest) GetRequest() *QueueSendRequest {
	if x != nil {
		return x.Request
	}
	return nil
}

func (x *FailedSendRequest) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_nitric_proto_queue_v1_queue_proto protoreflect.FileDescriptor

var file_nitric_proto_queue_v1_queue_proto_rawDesc = []byte{
	0x0a, 0x21, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x71,
	0x75, 0x65, 0x75, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x15, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75,
	0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x7b, 0x0a, 0x15, 0x51, 0x75, 0x65, 0x75,
	0x65, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x42, 0x61, 0x74, 0x63,
	0x68, 0x12, 0x1d, 0x0a, 0x0a, 0x71, 0x75, 0x65, 0x75, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x71, 0x75, 0x65, 0x75, 0x65, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x43, 0x0a, 0x08, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x27, 0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x51, 0x75, 0x65, 0x75, 0x65,
	0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x08, 0x72, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x73, 0x22, 0x66, 0x0a, 0x11, 0x51, 0x75, 0x65, 0x75, 0x65, 0x53, 0x65,
	0x6e, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x51, 0x0a, 0x0f, 0x66, 0x61,
	0x69, 0x6c, 0x65, 0x64, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x61, 0x69, 0x6c,
	0x65, 0x64, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x0e, 0x66,
	0x61, 0x69, 0x6c, 0x65, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x22, 0x4a, 0x0a,
	0x13, 0x51, 0x75, 0x65, 0x75, 0x65, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x71, 0x75, 0x65, 0x75, 0x65, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x71, 0x75, 0x65, 0x75, 0x65, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x64, 0x65, 0x70, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x05, 0x64, 0x65, 0x70, 0x74, 0x68, 0x22, 0x51, 0x0a, 0x14, 0x51, 0x75, 0x65,
	0x75, 0x65, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x39, 0x0a, 0x05, 0x74, 0x61, 0x73, 0x6b, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x23, 0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x71, 0x75, 0x65, 0x75, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65,
	0x64, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x05, 0x74, 0x61, 0x73, 0x6b, 0x73, 0x22, 0x50, 0x0a, 0x14,
	0x51, 0x75, 0x65, 0x75, 0x65, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x71, 0x75, 0x65, 0x75, 0x65, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x71, 0x75, 0x65, 0x75, 0x65, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x5f, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x49, 0x64, 0x22, 0x17,
	0x0a, 0x15, 0x51, 0x75, 0x65, 0x75, 0x65, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x55, 0x0a, 0x10, 0x51, 0x75, 0x65, 0x75, 0x65,
	0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x31, 0x0a, 0x07, 0x70,
	0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53,
	0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x5c,
	0x0a, 0x0c, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x64, 0x54, 0x61, 0x73, 0x6b, 0x12, 0x19,
	0x0a, 0x08, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x49, 0x64, 0x12, 0x31, 0x0a, 0x07, 0x70, 0x61, 0x79,
	0x6c, 0x6f, 0x61, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72,
	0x75, 0x63, 0x74, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x70, 0x0a, 0x11,
	0x46, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x41, 0x0a, 0x07, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x27, 0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x51, 0x75, 0x65, 0x75, 0x65,
	0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x07, 0x72, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x32, 0xb9,
	0x02, 0x0a, 0x0c, 0x51, 0x75, 0x65, 0x75, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x5e, 0x0a, 0x04, 0x53, 0x65, 0x6e, 0x64, 0x12, 0x2c, 0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2e, 0x76, 0x31, 0x2e,
	0x51, 0x75, 0x65, 0x75, 0x65, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x42, 0x61, 0x74, 0x63, 0x68, 0x1a, 0x28, 0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x51, 0x75,
	0x65, 0x75, 0x65, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x62, 0x0a, 0x07, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x12, 0x2a, 0x2e, 0x6e, 0x69, 0x74,
	0x72, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2b, 0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x71, 0x75, 0x65, 0x75, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x51,
	0x75, 0x65, 0x75, 0x65, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x65, 0x0a, 0x08, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x12,
	0x2b, 0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x71,
	0x75, 0x65, 0x75, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x43, 0x6f, 0x6d,
	0x70, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2c, 0x2e, 0x6e,
	0x69, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x71, 0x75, 0x65, 0x75,
	0x65, 0x2e, 0x76, 0x31, 0x2e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65,
	0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x92, 0x01, 0x0a, 0x18, 0x69,
	0x6f, 0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x71,
	0x75, 0x65, 0x75, 0x65, 0x2e, 0x76, 0x31, 0x42, 0x06, 0x51, 0x75, 0x65, 0x75, 0x65, 0x73, 0x50,
	0x01, 0x5a, 0x3c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x69,
	0x74, 0x72, 0x69, 0x63, 0x74, 0x65, 0x63, 0x68, 0x2f, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2f,
	0x63, 0x6f, 0x72, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x71,
	0x75, 0x65, 0x75, 0x65, 0x2f, 0x76, 0x31, 0x3b, 0x71, 0x75, 0x65, 0x75, 0x65, 0x70, 0x62, 0xaa,
	0x02, 0x15, 0x4e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x51,
	0x75, 0x65, 0x75, 0x65, 0x2e, 0x76, 0x31, 0xca, 0x02, 0x15, 0x4e, 0x69, 0x74, 0x72, 0x69, 0x63,
	0x5c, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x5c, 0x51, 0x75, 0x65, 0x75, 0x65, 0x5c, 0x56, 0x31, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_nitric_proto_queue_v1_queue_proto_rawDescOnce sync.Once
	file_nitric_proto_queue_v1_queue_proto_rawDescData = file_nitric_proto_queue_v1_queue_proto_rawDesc
)

func file_nitric_proto_queue_v1_queue_proto_rawDescGZIP() []byte {
	file_nitric_proto_queue_v1_queue_proto_rawDescOnce.Do(func() {
		file_nitric_proto_queue_v1_queue_proto_rawDescData = protoimpl.X.CompressGZIP(file_nitric_proto_queue_v1_queue_proto_rawDescData)
	})
	return file_nitric_proto_queue_v1_queue_proto_rawDescData
}

var file_nitric_proto_queue_v1_queue_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_nitric_proto_queue_v1_queue_proto_goTypes = []interface{}{
	(*QueueSendRequestBatch)(nil), // 0: nitric.proto.queue.v1.QueueSendRequestBatch
	(*QueueSendResponse)(nil),     // 1: nitric.proto.queue.v1.QueueSendResponse
	(*QueueReceiveRequest)(nil),   // 2: nitric.proto.queue.v1.QueueReceiveRequest
	(*QueueReceiveResponse)(nil),  // 3: nitric.proto.queue.v1.QueueReceiveResponse
	(*QueueCompleteRequest)(nil),  // 4: nitric.proto.queue.v1.QueueCompleteRequest
	(*QueueCompleteResponse)(nil), // 5: nitric.proto.queue.v1.QueueCompleteResponse
	(*QueueSendRequest)(nil),      // 6: nitric.proto.queue.v1.QueueSendRequest
	(*ReceivedTask)(nil),          // 7: nitric.proto.queue.v1.ReceivedTask
	(*FailedSendRequest)(nil),     // 8: nitric.proto.queue.v1.FailedSendRequest
	(*structpb.Struct)(nil),       // 9: google.protobuf.Struct
}
var file_nitric_proto_queue_v1_queue_proto_depIdxs = []int32{
	6, // 0: nitric.proto.queue.v1.QueueSendRequestBatch.requests:type_name -> nitric.proto.queue.v1.QueueSendRequest
	8, // 1: nitric.proto.queue.v1.QueueSendResponse.failed_requests:type_name -> nitric.proto.queue.v1.FailedSendRequest
	7, // 2: nitric.proto.queue.v1.QueueReceiveResponse.tasks:type_name -> nitric.proto.queue.v1.ReceivedTask
	9, // 3: nitric.proto.queue.v1.QueueSendRequest.payload:type_name -> google.protobuf.Struct
	9, // 4: nitric.proto.queue.v1.ReceivedTask.payload:type_name -> google.protobuf.Struct
	6, // 5: nitric.proto.queue.v1.FailedSendRequest.request:type_name -> nitric.proto.queue.v1.QueueSendRequest
	0, // 6: nitric.proto.queue.v1.QueueService.Send:input_type -> nitric.proto.queue.v1.QueueSendRequestBatch
	2, // 7: nitric.proto.queue.v1.QueueService.Receive:input_type -> nitric.proto.queue.v1.QueueReceiveRequest
	4, // 8: nitric.proto.queue.v1.QueueService.Complete:input_type -> nitric.proto.queue.v1.QueueCompleteRequest
	1, // 9: nitric.proto.queue.v1.QueueService.Send:output_type -> nitric.proto.queue.v1.QueueSendResponse
	3, // 10: nitric.proto.queue.v1.QueueService.Receive:output_type -> nitric.proto.queue.v1.QueueReceiveResponse
	5, // 11: nitric.proto.queue.v1.QueueService.Complete:output_type -> nitric.proto.queue.v1.QueueCompleteResponse
	9, // [9:12] is the sub-list for method output_type
	6, // [6:9] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_nitric_proto_queue_v1_queue_proto_init() }
func file_nitric_proto_queue_v1_queue_proto_init() {
	if File_nitric_proto_queue_v1_queue_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_nitric_proto_queue_v1_queue_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueueSendRequestBatch); i {
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
		file_nitric_proto_queue_v1_queue_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueueSendResponse); i {
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
		file_nitric_proto_queue_v1_queue_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueueReceiveRequest); i {
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
		file_nitric_proto_queue_v1_queue_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueueReceiveResponse); i {
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
		file_nitric_proto_queue_v1_queue_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueueCompleteRequest); i {
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
		file_nitric_proto_queue_v1_queue_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueueCompleteResponse); i {
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
		file_nitric_proto_queue_v1_queue_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueueSendRequest); i {
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
		file_nitric_proto_queue_v1_queue_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReceivedTask); i {
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
		file_nitric_proto_queue_v1_queue_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FailedSendRequest); i {
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
			RawDescriptor: file_nitric_proto_queue_v1_queue_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_nitric_proto_queue_v1_queue_proto_goTypes,
		DependencyIndexes: file_nitric_proto_queue_v1_queue_proto_depIdxs,
		MessageInfos:      file_nitric_proto_queue_v1_queue_proto_msgTypes,
	}.Build()
	File_nitric_proto_queue_v1_queue_proto = out.File
	file_nitric_proto_queue_v1_queue_proto_rawDesc = nil
	file_nitric_proto_queue_v1_queue_proto_goTypes = nil
	file_nitric_proto_queue_v1_queue_proto_depIdxs = nil
}
