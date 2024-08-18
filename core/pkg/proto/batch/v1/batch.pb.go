// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.4
// source: nitric/proto/batch/v1/batch.proto

package batchpb

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

type ClientMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// globally unique ID of the request/response pair
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Types that are assignable to Content:
	//
	//	*ClientMessage_RegistrationRequest
	//	*ClientMessage_JobResponse
	Content isClientMessage_Content `protobuf_oneof:"content"`
}

func (x *ClientMessage) Reset() {
	*x = ClientMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nitric_proto_batch_v1_batch_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientMessage) ProtoMessage() {}

func (x *ClientMessage) ProtoReflect() protoreflect.Message {
	mi := &file_nitric_proto_batch_v1_batch_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientMessage.ProtoReflect.Descriptor instead.
func (*ClientMessage) Descriptor() ([]byte, []int) {
	return file_nitric_proto_batch_v1_batch_proto_rawDescGZIP(), []int{0}
}

func (x *ClientMessage) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (m *ClientMessage) GetContent() isClientMessage_Content {
	if m != nil {
		return m.Content
	}
	return nil
}

func (x *ClientMessage) GetRegistrationRequest() *RegistrationRequest {
	if x, ok := x.GetContent().(*ClientMessage_RegistrationRequest); ok {
		return x.RegistrationRequest
	}
	return nil
}

func (x *ClientMessage) GetJobResponse() *JobResponse {
	if x, ok := x.GetContent().(*ClientMessage_JobResponse); ok {
		return x.JobResponse
	}
	return nil
}

type isClientMessage_Content interface {
	isClientMessage_Content()
}

type ClientMessage_RegistrationRequest struct {
	// Register a handler for a job
	RegistrationRequest *RegistrationRequest `protobuf:"bytes,2,opt,name=registration_request,json=registrationRequest,proto3,oneof"`
}

type ClientMessage_JobResponse struct {
	// Handle a job submission
	JobResponse *JobResponse `protobuf:"bytes,3,opt,name=job_response,json=jobResponse,proto3,oneof"`
}

func (*ClientMessage_RegistrationRequest) isClientMessage_Content() {}

func (*ClientMessage_JobResponse) isClientMessage_Content() {}

type JobRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	JobName string   `protobuf:"bytes,1,opt,name=job_name,json=jobName,proto3" json:"job_name,omitempty"`
	Data    *JobData `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *JobRequest) Reset() {
	*x = JobRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nitric_proto_batch_v1_batch_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JobRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobRequest) ProtoMessage() {}

func (x *JobRequest) ProtoReflect() protoreflect.Message {
	mi := &file_nitric_proto_batch_v1_batch_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JobRequest.ProtoReflect.Descriptor instead.
func (*JobRequest) Descriptor() ([]byte, []int) {
	return file_nitric_proto_batch_v1_batch_proto_rawDescGZIP(), []int{1}
}

func (x *JobRequest) GetJobName() string {
	if x != nil {
		return x.JobName
	}
	return ""
}

func (x *JobRequest) GetData() *JobData {
	if x != nil {
		return x.Data
	}
	return nil
}

type JobData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Data:
	//
	//	*JobData_Struct
	Data isJobData_Data `protobuf_oneof:"data"`
}

func (x *JobData) Reset() {
	*x = JobData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nitric_proto_batch_v1_batch_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JobData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobData) ProtoMessage() {}

func (x *JobData) ProtoReflect() protoreflect.Message {
	mi := &file_nitric_proto_batch_v1_batch_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JobData.ProtoReflect.Descriptor instead.
func (*JobData) Descriptor() ([]byte, []int) {
	return file_nitric_proto_batch_v1_batch_proto_rawDescGZIP(), []int{2}
}

func (m *JobData) GetData() isJobData_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (x *JobData) GetStruct() *structpb.Struct {
	if x, ok := x.GetData().(*JobData_Struct); ok {
		return x.Struct
	}
	return nil
}

type isJobData_Data interface {
	isJobData_Data()
}

type JobData_Struct struct {
	Struct *structpb.Struct `protobuf:"bytes,1,opt,name=struct,proto3,oneof"`
}

func (*JobData_Struct) isJobData_Data() {}

type JobResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Mark if the job was successfully processed
	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *JobResponse) Reset() {
	*x = JobResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nitric_proto_batch_v1_batch_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JobResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobResponse) ProtoMessage() {}

func (x *JobResponse) ProtoReflect() protoreflect.Message {
	mi := &file_nitric_proto_batch_v1_batch_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JobResponse.ProtoReflect.Descriptor instead.
func (*JobResponse) Descriptor() ([]byte, []int) {
	return file_nitric_proto_batch_v1_batch_proto_rawDescGZIP(), []int{3}
}

func (x *JobResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type RegistrationRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	JobName string `protobuf:"bytes,1,opt,name=job_name,json=jobName,proto3" json:"job_name,omitempty"`
	// Register with default requirements
	Requirements *JobResourceRequirements `protobuf:"bytes,2,opt,name=requirements,proto3" json:"requirements,omitempty"`
}

func (x *RegistrationRequest) Reset() {
	*x = RegistrationRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nitric_proto_batch_v1_batch_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegistrationRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegistrationRequest) ProtoMessage() {}

func (x *RegistrationRequest) ProtoReflect() protoreflect.Message {
	mi := &file_nitric_proto_batch_v1_batch_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegistrationRequest.ProtoReflect.Descriptor instead.
func (*RegistrationRequest) Descriptor() ([]byte, []int) {
	return file_nitric_proto_batch_v1_batch_proto_rawDescGZIP(), []int{4}
}

func (x *RegistrationRequest) GetJobName() string {
	if x != nil {
		return x.JobName
	}
	return ""
}

func (x *RegistrationRequest) GetRequirements() *JobResourceRequirements {
	if x != nil {
		return x.Requirements
	}
	return nil
}

type RegistrationResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RegistrationResponse) Reset() {
	*x = RegistrationResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nitric_proto_batch_v1_batch_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegistrationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegistrationResponse) ProtoMessage() {}

func (x *RegistrationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_nitric_proto_batch_v1_batch_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegistrationResponse.ProtoReflect.Descriptor instead.
func (*RegistrationResponse) Descriptor() ([]byte, []int) {
	return file_nitric_proto_batch_v1_batch_proto_rawDescGZIP(), []int{5}
}

type JobResourceRequirements struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The number of CPUs to allocate for the job
	Cpus float32 `protobuf:"fixed32,1,opt,name=cpus,proto3" json:"cpus,omitempty"`
	// The amount of memory to allocate for the job
	Memory int64 `protobuf:"varint,2,opt,name=memory,proto3" json:"memory,omitempty"`
	// The number of GPUs to allocate for the job
	Gpus int64 `protobuf:"varint,3,opt,name=gpus,proto3" json:"gpus,omitempty"`
}

func (x *JobResourceRequirements) Reset() {
	*x = JobResourceRequirements{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nitric_proto_batch_v1_batch_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JobResourceRequirements) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobResourceRequirements) ProtoMessage() {}

func (x *JobResourceRequirements) ProtoReflect() protoreflect.Message {
	mi := &file_nitric_proto_batch_v1_batch_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JobResourceRequirements.ProtoReflect.Descriptor instead.
func (*JobResourceRequirements) Descriptor() ([]byte, []int) {
	return file_nitric_proto_batch_v1_batch_proto_rawDescGZIP(), []int{6}
}

func (x *JobResourceRequirements) GetCpus() float32 {
	if x != nil {
		return x.Cpus
	}
	return 0
}

func (x *JobResourceRequirements) GetMemory() int64 {
	if x != nil {
		return x.Memory
	}
	return 0
}

func (x *JobResourceRequirements) GetGpus() int64 {
	if x != nil {
		return x.Gpus
	}
	return 0
}

// ServerMessage is the message sent from the nitric server to the service
type ServerMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// globally unique ID of the request/response pair
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Types that are assignable to Content:
	//
	//	*ServerMessage_RegistrationResponse
	//	*ServerMessage_JobRequest
	Content isServerMessage_Content `protobuf_oneof:"content"`
}

func (x *ServerMessage) Reset() {
	*x = ServerMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nitric_proto_batch_v1_batch_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerMessage) ProtoMessage() {}

func (x *ServerMessage) ProtoReflect() protoreflect.Message {
	mi := &file_nitric_proto_batch_v1_batch_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerMessage.ProtoReflect.Descriptor instead.
func (*ServerMessage) Descriptor() ([]byte, []int) {
	return file_nitric_proto_batch_v1_batch_proto_rawDescGZIP(), []int{7}
}

func (x *ServerMessage) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (m *ServerMessage) GetContent() isServerMessage_Content {
	if m != nil {
		return m.Content
	}
	return nil
}

func (x *ServerMessage) GetRegistrationResponse() *RegistrationResponse {
	if x, ok := x.GetContent().(*ServerMessage_RegistrationResponse); ok {
		return x.RegistrationResponse
	}
	return nil
}

func (x *ServerMessage) GetJobRequest() *JobRequest {
	if x, ok := x.GetContent().(*ServerMessage_JobRequest); ok {
		return x.JobRequest
	}
	return nil
}

type isServerMessage_Content interface {
	isServerMessage_Content()
}

type ServerMessage_RegistrationResponse struct {
	RegistrationResponse *RegistrationResponse `protobuf:"bytes,2,opt,name=registration_response,json=registrationResponse,proto3,oneof"`
}

type ServerMessage_JobRequest struct {
	// Request to a job handler
	JobRequest *JobRequest `protobuf:"bytes,3,opt,name=job_request,json=jobRequest,proto3,oneof"`
}

func (*ServerMessage_RegistrationResponse) isServerMessage_Content() {}

func (*ServerMessage_JobRequest) isServerMessage_Content() {}

type JobSubmitRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The name of the job that should handle the data
	JobName string `protobuf:"bytes,1,opt,name=job_name,json=jobName,proto3" json:"job_name,omitempty"`
	// The data to be processed by the job
	Data *JobData `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *JobSubmitRequest) Reset() {
	*x = JobSubmitRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nitric_proto_batch_v1_batch_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JobSubmitRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobSubmitRequest) ProtoMessage() {}

func (x *JobSubmitRequest) ProtoReflect() protoreflect.Message {
	mi := &file_nitric_proto_batch_v1_batch_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JobSubmitRequest.ProtoReflect.Descriptor instead.
func (*JobSubmitRequest) Descriptor() ([]byte, []int) {
	return file_nitric_proto_batch_v1_batch_proto_rawDescGZIP(), []int{8}
}

func (x *JobSubmitRequest) GetJobName() string {
	if x != nil {
		return x.JobName
	}
	return ""
}

func (x *JobSubmitRequest) GetData() *JobData {
	if x != nil {
		return x.Data
	}
	return nil
}

type JobSubmitResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *JobSubmitResponse) Reset() {
	*x = JobSubmitResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_nitric_proto_batch_v1_batch_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JobSubmitResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JobSubmitResponse) ProtoMessage() {}

func (x *JobSubmitResponse) ProtoReflect() protoreflect.Message {
	mi := &file_nitric_proto_batch_v1_batch_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JobSubmitResponse.ProtoReflect.Descriptor instead.
func (*JobSubmitResponse) Descriptor() ([]byte, []int) {
	return file_nitric_proto_batch_v1_batch_proto_rawDescGZIP(), []int{9}
}

var File_nitric_proto_batch_v1_batch_proto protoreflect.FileDescriptor

var file_nitric_proto_batch_v1_batch_proto_rawDesc = []byte{
	0x0a, 0x21, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x62,
	0x61, 0x74, 0x63, 0x68, 0x2f, 0x76, 0x31, 0x2f, 0x62, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x15, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x62, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75,
	0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd4, 0x01, 0x0a, 0x0d, 0x43, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x5f, 0x0a, 0x14, 0x72, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69,
	0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x62, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x76, 0x31,
	0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x48, 0x00, 0x52, 0x13, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x47, 0x0a, 0x0c, 0x6a,
	0x6f, 0x62, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x22, 0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x62, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x48, 0x00, 0x52, 0x0b, 0x6a, 0x6f, 0x62, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x42, 0x09, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22,
	0x5b, 0x0a, 0x0a, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a,
	0x08, 0x6a, 0x6f, 0x62, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x6a, 0x6f, 0x62, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x32, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x62, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x4a,
	0x6f, 0x62, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x44, 0x0a, 0x07,
	0x4a, 0x6f, 0x62, 0x44, 0x61, 0x74, 0x61, 0x12, 0x31, 0x0a, 0x06, 0x73, 0x74, 0x72, 0x75, 0x63,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74,
	0x48, 0x00, 0x52, 0x06, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x42, 0x06, 0x0a, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x22, 0x27, 0x0a, 0x0b, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x22, 0x84, 0x01, 0x0a, 0x13,
	0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x6a, 0x6f, 0x62, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6a, 0x6f, 0x62, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x52,
	0x0a, 0x0c, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x2e, 0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x62, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x4a, 0x6f, 0x62,
	0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x6d,
	0x65, 0x6e, 0x74, 0x73, 0x52, 0x0c, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x6d, 0x65, 0x6e,
	0x74, 0x73, 0x22, 0x16, 0x0a, 0x14, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x59, 0x0a, 0x17, 0x4a, 0x6f,
	0x62, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65,
	0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x70, 0x75, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x02, 0x52, 0x04, 0x63, 0x70, 0x75, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65, 0x6d,
	0x6f, 0x72, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x6d, 0x65, 0x6d, 0x6f, 0x72,
	0x79, 0x12, 0x12, 0x0a, 0x04, 0x67, 0x70, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x04, 0x67, 0x70, 0x75, 0x73, 0x22, 0xd4, 0x01, 0x0a, 0x0d, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x62, 0x0a, 0x15, 0x72, 0x65, 0x67, 0x69, 0x73,
	0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x62, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x52,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x48, 0x00, 0x52, 0x14, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x44, 0x0a, 0x0b, 0x6a,
	0x6f, 0x62, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x21, 0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x62, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x48, 0x00, 0x52, 0x0a, 0x6a, 0x6f, 0x62, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x42, 0x09, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x61, 0x0a, 0x10,
	0x4a, 0x6f, 0x62, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x19, 0x0a, 0x08, 0x6a, 0x6f, 0x62, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x6a, 0x6f, 0x62, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x32, 0x0a, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x6e, 0x69, 0x74, 0x72,
	0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x62, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x76,
	0x31, 0x2e, 0x4a, 0x6f, 0x62, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22,
	0x13, 0x0a, 0x11, 0x4a, 0x6f, 0x62, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x32, 0x62, 0x0a, 0x03, 0x4a, 0x6f, 0x62, 0x12, 0x5b, 0x0a, 0x09, 0x48,
	0x61, 0x6e, 0x64, 0x6c, 0x65, 0x4a, 0x6f, 0x62, 0x12, 0x24, 0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69,
	0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x62, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x76, 0x31,
	0x2e, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x24,
	0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x62, 0x61,
	0x74, 0x63, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x28, 0x01, 0x30, 0x01, 0x32, 0x67, 0x0a, 0x05, 0x42, 0x61, 0x74, 0x63,
	0x68, 0x12, 0x5e, 0x0a, 0x09, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x4a, 0x6f, 0x62, 0x12, 0x27,
	0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x62, 0x61,
	0x74, 0x63, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x4a, 0x6f, 0x62, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x28, 0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x62, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x76, 0x31, 0x2e,
	0x4a, 0x6f, 0x62, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x42, 0x98, 0x01, 0x0a, 0x18, 0x69, 0x6f, 0x2e, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x62, 0x61, 0x74, 0x63, 0x68, 0x2e, 0x76, 0x31, 0x42, 0x0c,
	0x42, 0x61, 0x74, 0x63, 0x68, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x50, 0x01, 0x5a, 0x3c,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6e, 0x69, 0x74, 0x72, 0x69,
	0x63, 0x74, 0x65, 0x63, 0x68, 0x2f, 0x6e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x2f, 0x63, 0x6f, 0x72,
	0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x62, 0x61, 0x74, 0x63,
	0x68, 0x2f, 0x76, 0x31, 0x3b, 0x62, 0x61, 0x74, 0x63, 0x68, 0x70, 0x62, 0xaa, 0x02, 0x15, 0x4e,
	0x69, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x42, 0x61, 0x74, 0x63,
	0x68, 0x2e, 0x76, 0x31, 0xca, 0x02, 0x15, 0x4e, 0x69, 0x74, 0x72, 0x69, 0x63, 0x5c, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x5c, 0x42, 0x61, 0x74, 0x63, 0x68, 0x5c, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_nitric_proto_batch_v1_batch_proto_rawDescOnce sync.Once
	file_nitric_proto_batch_v1_batch_proto_rawDescData = file_nitric_proto_batch_v1_batch_proto_rawDesc
)

func file_nitric_proto_batch_v1_batch_proto_rawDescGZIP() []byte {
	file_nitric_proto_batch_v1_batch_proto_rawDescOnce.Do(func() {
		file_nitric_proto_batch_v1_batch_proto_rawDescData = protoimpl.X.CompressGZIP(file_nitric_proto_batch_v1_batch_proto_rawDescData)
	})
	return file_nitric_proto_batch_v1_batch_proto_rawDescData
}

var file_nitric_proto_batch_v1_batch_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_nitric_proto_batch_v1_batch_proto_goTypes = []interface{}{
	(*ClientMessage)(nil),           // 0: nitric.proto.batch.v1.ClientMessage
	(*JobRequest)(nil),              // 1: nitric.proto.batch.v1.JobRequest
	(*JobData)(nil),                 // 2: nitric.proto.batch.v1.JobData
	(*JobResponse)(nil),             // 3: nitric.proto.batch.v1.JobResponse
	(*RegistrationRequest)(nil),     // 4: nitric.proto.batch.v1.RegistrationRequest
	(*RegistrationResponse)(nil),    // 5: nitric.proto.batch.v1.RegistrationResponse
	(*JobResourceRequirements)(nil), // 6: nitric.proto.batch.v1.JobResourceRequirements
	(*ServerMessage)(nil),           // 7: nitric.proto.batch.v1.ServerMessage
	(*JobSubmitRequest)(nil),        // 8: nitric.proto.batch.v1.JobSubmitRequest
	(*JobSubmitResponse)(nil),       // 9: nitric.proto.batch.v1.JobSubmitResponse
	(*structpb.Struct)(nil),         // 10: google.protobuf.Struct
}
var file_nitric_proto_batch_v1_batch_proto_depIdxs = []int32{
	4,  // 0: nitric.proto.batch.v1.ClientMessage.registration_request:type_name -> nitric.proto.batch.v1.RegistrationRequest
	3,  // 1: nitric.proto.batch.v1.ClientMessage.job_response:type_name -> nitric.proto.batch.v1.JobResponse
	2,  // 2: nitric.proto.batch.v1.JobRequest.data:type_name -> nitric.proto.batch.v1.JobData
	10, // 3: nitric.proto.batch.v1.JobData.struct:type_name -> google.protobuf.Struct
	6,  // 4: nitric.proto.batch.v1.RegistrationRequest.requirements:type_name -> nitric.proto.batch.v1.JobResourceRequirements
	5,  // 5: nitric.proto.batch.v1.ServerMessage.registration_response:type_name -> nitric.proto.batch.v1.RegistrationResponse
	1,  // 6: nitric.proto.batch.v1.ServerMessage.job_request:type_name -> nitric.proto.batch.v1.JobRequest
	2,  // 7: nitric.proto.batch.v1.JobSubmitRequest.data:type_name -> nitric.proto.batch.v1.JobData
	0,  // 8: nitric.proto.batch.v1.Job.HandleJob:input_type -> nitric.proto.batch.v1.ClientMessage
	8,  // 9: nitric.proto.batch.v1.Batch.SubmitJob:input_type -> nitric.proto.batch.v1.JobSubmitRequest
	7,  // 10: nitric.proto.batch.v1.Job.HandleJob:output_type -> nitric.proto.batch.v1.ServerMessage
	9,  // 11: nitric.proto.batch.v1.Batch.SubmitJob:output_type -> nitric.proto.batch.v1.JobSubmitResponse
	10, // [10:12] is the sub-list for method output_type
	8,  // [8:10] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_nitric_proto_batch_v1_batch_proto_init() }
func file_nitric_proto_batch_v1_batch_proto_init() {
	if File_nitric_proto_batch_v1_batch_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_nitric_proto_batch_v1_batch_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientMessage); i {
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
		file_nitric_proto_batch_v1_batch_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JobRequest); i {
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
		file_nitric_proto_batch_v1_batch_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JobData); i {
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
		file_nitric_proto_batch_v1_batch_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JobResponse); i {
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
		file_nitric_proto_batch_v1_batch_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegistrationRequest); i {
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
		file_nitric_proto_batch_v1_batch_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegistrationResponse); i {
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
		file_nitric_proto_batch_v1_batch_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JobResourceRequirements); i {
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
		file_nitric_proto_batch_v1_batch_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerMessage); i {
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
		file_nitric_proto_batch_v1_batch_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JobSubmitRequest); i {
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
		file_nitric_proto_batch_v1_batch_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JobSubmitResponse); i {
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
	file_nitric_proto_batch_v1_batch_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*ClientMessage_RegistrationRequest)(nil),
		(*ClientMessage_JobResponse)(nil),
	}
	file_nitric_proto_batch_v1_batch_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*JobData_Struct)(nil),
	}
	file_nitric_proto_batch_v1_batch_proto_msgTypes[7].OneofWrappers = []interface{}{
		(*ServerMessage_RegistrationResponse)(nil),
		(*ServerMessage_JobRequest)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_nitric_proto_batch_v1_batch_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_nitric_proto_batch_v1_batch_proto_goTypes,
		DependencyIndexes: file_nitric_proto_batch_v1_batch_proto_depIdxs,
		MessageInfos:      file_nitric_proto_batch_v1_batch_proto_msgTypes,
	}.Build()
	File_nitric_proto_batch_v1_batch_proto = out.File
	file_nitric_proto_batch_v1_batch_proto_rawDesc = nil
	file_nitric_proto_batch_v1_batch_proto_goTypes = nil
	file_nitric_proto_batch_v1_batch_proto_depIdxs = nil
}
