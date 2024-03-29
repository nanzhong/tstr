// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: control/v1/control.proto

package controlv1

import (
	v1 "github.com/nanzhong/tstr/api/common/v1"
	_ "github.com/nanzhong/tstr/api/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	fieldmaskpb "google.golang.org/protobuf/types/known/fieldmaskpb"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RegisterTestRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name         string             `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Labels       map[string]string  `protobuf:"bytes,2,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	RunConfig    *v1.Test_RunConfig `protobuf:"bytes,3,opt,name=run_config,proto3" json:"run_config,omitempty"`
	CronSchedule string             `protobuf:"bytes,4,opt,name=cron_schedule,json=cronSchedule,proto3" json:"cron_schedule,omitempty"`
	Matrix       *v1.Test_Matrix    `protobuf:"bytes,5,opt,name=matrix,proto3" json:"matrix,omitempty"`
}

func (x *RegisterTestRequest) Reset() {
	*x = RegisterTestRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_control_v1_control_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterTestRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterTestRequest) ProtoMessage() {}

func (x *RegisterTestRequest) ProtoReflect() protoreflect.Message {
	mi := &file_control_v1_control_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterTestRequest.ProtoReflect.Descriptor instead.
func (*RegisterTestRequest) Descriptor() ([]byte, []int) {
	return file_control_v1_control_proto_rawDescGZIP(), []int{0}
}

func (x *RegisterTestRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *RegisterTestRequest) GetLabels() map[string]string {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (x *RegisterTestRequest) GetRunConfig() *v1.Test_RunConfig {
	if x != nil {
		return x.RunConfig
	}
	return nil
}

func (x *RegisterTestRequest) GetCronSchedule() string {
	if x != nil {
		return x.CronSchedule
	}
	return ""
}

func (x *RegisterTestRequest) GetMatrix() *v1.Test_Matrix {
	if x != nil {
		return x.Matrix
	}
	return nil
}

type RegisterTestResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Test *v1.Test `protobuf:"bytes,1,opt,name=test,proto3" json:"test,omitempty"` // required
}

func (x *RegisterTestResponse) Reset() {
	*x = RegisterTestResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_control_v1_control_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterTestResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterTestResponse) ProtoMessage() {}

func (x *RegisterTestResponse) ProtoReflect() protoreflect.Message {
	mi := &file_control_v1_control_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterTestResponse.ProtoReflect.Descriptor instead.
func (*RegisterTestResponse) Descriptor() ([]byte, []int) {
	return file_control_v1_control_proto_rawDescGZIP(), []int{1}
}

func (x *RegisterTestResponse) GetTest() *v1.Test {
	if x != nil {
		return x.Test
	}
	return nil
}

type UpdateTestRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FieldMask    *fieldmaskpb.FieldMask `protobuf:"bytes,1,opt,name=field_mask,json=fieldMask,proto3" json:"field_mask,omitempty"`
	Id           string                 `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"` // required
	Name         string                 `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	RunConfig    *v1.Test_RunConfig     `protobuf:"bytes,4,opt,name=run_config,json=runConfig,proto3" json:"run_config,omitempty"`
	Labels       map[string]string      `protobuf:"bytes,5,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	CronSchedule string                 `protobuf:"bytes,6,opt,name=cron_schedule,json=cronSchedule,proto3" json:"cron_schedule,omitempty"`
	Matrix       *v1.Test_Matrix        `protobuf:"bytes,7,opt,name=matrix,proto3" json:"matrix,omitempty"`
}

func (x *UpdateTestRequest) Reset() {
	*x = UpdateTestRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_control_v1_control_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateTestRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateTestRequest) ProtoMessage() {}

func (x *UpdateTestRequest) ProtoReflect() protoreflect.Message {
	mi := &file_control_v1_control_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateTestRequest.ProtoReflect.Descriptor instead.
func (*UpdateTestRequest) Descriptor() ([]byte, []int) {
	return file_control_v1_control_proto_rawDescGZIP(), []int{2}
}

func (x *UpdateTestRequest) GetFieldMask() *fieldmaskpb.FieldMask {
	if x != nil {
		return x.FieldMask
	}
	return nil
}

func (x *UpdateTestRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UpdateTestRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UpdateTestRequest) GetRunConfig() *v1.Test_RunConfig {
	if x != nil {
		return x.RunConfig
	}
	return nil
}

func (x *UpdateTestRequest) GetLabels() map[string]string {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (x *UpdateTestRequest) GetCronSchedule() string {
	if x != nil {
		return x.CronSchedule
	}
	return ""
}

func (x *UpdateTestRequest) GetMatrix() *v1.Test_Matrix {
	if x != nil {
		return x.Matrix
	}
	return nil
}

type UpdateTestResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UpdateTestResponse) Reset() {
	*x = UpdateTestResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_control_v1_control_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateTestResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateTestResponse) ProtoMessage() {}

func (x *UpdateTestResponse) ProtoReflect() protoreflect.Message {
	mi := &file_control_v1_control_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateTestResponse.ProtoReflect.Descriptor instead.
func (*UpdateTestResponse) Descriptor() ([]byte, []int) {
	return file_control_v1_control_proto_rawDescGZIP(), []int{3}
}

type DeleteTestRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"` // required
}

func (x *DeleteTestRequest) Reset() {
	*x = DeleteTestRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_control_v1_control_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteTestRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteTestRequest) ProtoMessage() {}

func (x *DeleteTestRequest) ProtoReflect() protoreflect.Message {
	mi := &file_control_v1_control_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteTestRequest.ProtoReflect.Descriptor instead.
func (*DeleteTestRequest) Descriptor() ([]byte, []int) {
	return file_control_v1_control_proto_rawDescGZIP(), []int{4}
}

func (x *DeleteTestRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type DeleteTestResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteTestResponse) Reset() {
	*x = DeleteTestResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_control_v1_control_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteTestResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteTestResponse) ProtoMessage() {}

func (x *DeleteTestResponse) ProtoReflect() protoreflect.Message {
	mi := &file_control_v1_control_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteTestResponse.ProtoReflect.Descriptor instead.
func (*DeleteTestResponse) Descriptor() ([]byte, []int) {
	return file_control_v1_control_proto_rawDescGZIP(), []int{5}
}

type ScheduleRunRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TestId     string            `protobuf:"bytes,1,opt,name=test_id,json=testId,proto3" json:"test_id,omitempty"` // required
	Labels     map[string]string `protobuf:"bytes,2,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	TestMatrix *v1.Test_Matrix   `protobuf:"bytes,3,opt,name=test_matrix,json=testMatrix,proto3" json:"test_matrix,omitempty"`
}

func (x *ScheduleRunRequest) Reset() {
	*x = ScheduleRunRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_control_v1_control_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ScheduleRunRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScheduleRunRequest) ProtoMessage() {}

func (x *ScheduleRunRequest) ProtoReflect() protoreflect.Message {
	mi := &file_control_v1_control_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScheduleRunRequest.ProtoReflect.Descriptor instead.
func (*ScheduleRunRequest) Descriptor() ([]byte, []int) {
	return file_control_v1_control_proto_rawDescGZIP(), []int{6}
}

func (x *ScheduleRunRequest) GetTestId() string {
	if x != nil {
		return x.TestId
	}
	return ""
}

func (x *ScheduleRunRequest) GetLabels() map[string]string {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (x *ScheduleRunRequest) GetTestMatrix() *v1.Test_Matrix {
	if x != nil {
		return x.TestMatrix
	}
	return nil
}

type ScheduleRunResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Runs []*v1.Run `protobuf:"bytes,1,rep,name=runs,proto3" json:"runs,omitempty"` // required
}

func (x *ScheduleRunResponse) Reset() {
	*x = ScheduleRunResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_control_v1_control_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ScheduleRunResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScheduleRunResponse) ProtoMessage() {}

func (x *ScheduleRunResponse) ProtoReflect() protoreflect.Message {
	mi := &file_control_v1_control_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScheduleRunResponse.ProtoReflect.Descriptor instead.
func (*ScheduleRunResponse) Descriptor() ([]byte, []int) {
	return file_control_v1_control_proto_rawDescGZIP(), []int{7}
}

func (x *ScheduleRunResponse) GetRuns() []*v1.Run {
	if x != nil {
		return x.Runs
	}
	return nil
}

var File_control_v1_control_proto protoreflect.FileDescriptor

var file_control_v1_control_proto_rawDesc = []byte{
	0x0a, 0x18, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6e,
	0x74, 0x72, 0x6f, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x74, 0x73, 0x74, 0x72,
	0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2e, 0x76, 0x31, 0x1a, 0x16, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x6d, 0x61, 0x73, 0x6b, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65,
	0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xde, 0x02, 0x0a, 0x13, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x54, 0x65, 0x73, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1e, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0a, 0xfa, 0x42, 0x07, 0x72, 0x05, 0x10, 0x01, 0x18, 0xc8,
	0x01, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x48, 0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x30, 0x2e, 0x74, 0x73, 0x74, 0x72, 0x2e, 0x63,
	0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74,
	0x65, 0x72, 0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x4c, 0x61,
	0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c,
	0x73, 0x12, 0x48, 0x0a, 0x0a, 0x72, 0x75, 0x6e, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x74, 0x73, 0x74, 0x72, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x2e, 0x52, 0x75, 0x6e, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x8a, 0x01, 0x02, 0x10, 0x01, 0x52,
	0x0a, 0x72, 0x75, 0x6e, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x23, 0x0a, 0x0d, 0x63,
	0x72, 0x6f, 0x6e, 0x5f, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0c, 0x63, 0x72, 0x6f, 0x6e, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65,
	0x12, 0x33, 0x0a, 0x06, 0x6d, 0x61, 0x74, 0x72, 0x69, 0x78, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1b, 0x2e, 0x74, 0x73, 0x74, 0x72, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76,
	0x31, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x2e, 0x4d, 0x61, 0x74, 0x72, 0x69, 0x78, 0x52, 0x06, 0x6d,
	0x61, 0x74, 0x72, 0x69, 0x78, 0x1a, 0x39, 0x0a, 0x0b, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01,
	0x22, 0x40, 0x0a, 0x14, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x54, 0x65, 0x73, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x28, 0x0a, 0x04, 0x74, 0x65, 0x73, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x74, 0x73, 0x74, 0x72, 0x2e, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x52, 0x04, 0x74, 0x65,
	0x73, 0x74, 0x22, 0x8e, 0x03, 0x0a, 0x11, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x65, 0x73,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x66, 0x69, 0x65, 0x6c,
	0x64, 0x5f, 0x6d, 0x61, 0x73, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46,
	0x69, 0x65, 0x6c, 0x64, 0x4d, 0x61, 0x73, 0x6b, 0x52, 0x09, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x4d,
	0x61, 0x73, 0x6b, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x3d, 0x0a, 0x0a, 0x72, 0x75, 0x6e, 0x5f, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x74, 0x73,
	0x74, 0x72, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x65, 0x73,
	0x74, 0x2e, 0x52, 0x75, 0x6e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x09, 0x72, 0x75, 0x6e,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x46, 0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73,
	0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2e, 0x2e, 0x74, 0x73, 0x74, 0x72, 0x2e, 0x63, 0x6f,
	0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54,
	0x65, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x12, 0x23,
	0x0a, 0x0d, 0x63, 0x72, 0x6f, 0x6e, 0x5f, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x63, 0x72, 0x6f, 0x6e, 0x53, 0x63, 0x68, 0x65, 0x64,
	0x75, 0x6c, 0x65, 0x12, 0x33, 0x0a, 0x06, 0x6d, 0x61, 0x74, 0x72, 0x69, 0x78, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x74, 0x73, 0x74, 0x72, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x2e, 0x4d, 0x61, 0x74, 0x72, 0x69, 0x78,
	0x52, 0x06, 0x6d, 0x61, 0x74, 0x72, 0x69, 0x78, 0x1a, 0x39, 0x0a, 0x0b, 0x4c, 0x61, 0x62, 0x65,
	0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a,
	0x02, 0x38, 0x01, 0x22, 0x14, 0x0a, 0x12, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x65, 0x73,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x23, 0x0a, 0x11, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x14,
	0x0a, 0x12, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0xef, 0x01, 0x0a, 0x12, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c,
	0x65, 0x52, 0x75, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x74,
	0x65, 0x73, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x74, 0x65,
	0x73, 0x74, 0x49, 0x64, 0x12, 0x47, 0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x2f, 0x2e, 0x74, 0x73, 0x74, 0x72, 0x2e, 0x63, 0x6f, 0x6e, 0x74,
	0x72, 0x6f, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x52,
	0x75, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x12, 0x3c, 0x0a,
	0x0b, 0x74, 0x65, 0x73, 0x74, 0x5f, 0x6d, 0x61, 0x74, 0x72, 0x69, 0x78, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x74, 0x73, 0x74, 0x72, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2e, 0x76, 0x31, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x2e, 0x4d, 0x61, 0x74, 0x72, 0x69, 0x78, 0x52,
	0x0a, 0x74, 0x65, 0x73, 0x74, 0x4d, 0x61, 0x74, 0x72, 0x69, 0x78, 0x1a, 0x39, 0x0a, 0x0b, 0x4c,
	0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x3e, 0x0a, 0x13, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75,
	0x6c, 0x65, 0x52, 0x75, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x27, 0x0a,
	0x04, 0x72, 0x75, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x74, 0x73,
	0x74, 0x72, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x75, 0x6e,
	0x52, 0x04, 0x72, 0x75, 0x6e, 0x73, 0x32, 0xf5, 0x02, 0x0a, 0x0e, 0x43, 0x6f, 0x6e, 0x74, 0x72,
	0x6f, 0x6c, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x5b, 0x0a, 0x0c, 0x52, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x65, 0x72, 0x54, 0x65, 0x73, 0x74, 0x12, 0x24, 0x2e, 0x74, 0x73, 0x74, 0x72,
	0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x65, 0x72, 0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x25, 0x2e, 0x74, 0x73, 0x74, 0x72, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2e, 0x76,
	0x31, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x54, 0x65, 0x73, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x55, 0x0a, 0x0a, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x54, 0x65, 0x73, 0x74, 0x12, 0x22, 0x2e, 0x74, 0x73, 0x74, 0x72, 0x2e, 0x63, 0x6f, 0x6e, 0x74,
	0x72, 0x6f, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x65, 0x73,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x74, 0x73, 0x74, 0x72, 0x2e,
	0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x55, 0x0a,
	0x0a, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x54, 0x65, 0x73, 0x74, 0x12, 0x22, 0x2e, 0x74, 0x73,
	0x74, 0x72, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x23, 0x2e, 0x74, 0x73, 0x74, 0x72, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2e, 0x76,
	0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x58, 0x0a, 0x0b, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65,
	0x52, 0x75, 0x6e, 0x12, 0x23, 0x2e, 0x74, 0x73, 0x74, 0x72, 0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72,
	0x6f, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x75,
	0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x74, 0x73, 0x74, 0x72, 0x2e,
	0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x63, 0x68, 0x65, 0x64,
	0x75, 0x6c, 0x65, 0x52, 0x75, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0xb4,
	0x01, 0x0a, 0x13, 0x63, 0x6f, 0x6d, 0x2e, 0x74, 0x73, 0x74, 0x72, 0x2e, 0x63, 0x6f, 0x6e, 0x74,
	0x72, 0x6f, 0x6c, 0x2e, 0x76, 0x31, 0x42, 0x0c, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x31, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x6e, 0x61, 0x6e, 0x7a, 0x68, 0x6f, 0x6e, 0x67, 0x2f, 0x74, 0x73, 0x74, 0x72,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2f, 0x76, 0x31, 0x3b,
	0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x54, 0x43, 0x58, 0xaa,
	0x02, 0x0f, 0x54, 0x73, 0x74, 0x72, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2e, 0x56,
	0x31, 0xca, 0x02, 0x0f, 0x54, 0x73, 0x74, 0x72, 0x5c, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c,
	0x5c, 0x56, 0x31, 0xe2, 0x02, 0x1b, 0x54, 0x73, 0x74, 0x72, 0x5c, 0x43, 0x6f, 0x6e, 0x74, 0x72,
	0x6f, 0x6c, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0xea, 0x02, 0x11, 0x54, 0x73, 0x74, 0x72, 0x3a, 0x3a, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f,
	0x6c, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_control_v1_control_proto_rawDescOnce sync.Once
	file_control_v1_control_proto_rawDescData = file_control_v1_control_proto_rawDesc
)

func file_control_v1_control_proto_rawDescGZIP() []byte {
	file_control_v1_control_proto_rawDescOnce.Do(func() {
		file_control_v1_control_proto_rawDescData = protoimpl.X.CompressGZIP(file_control_v1_control_proto_rawDescData)
	})
	return file_control_v1_control_proto_rawDescData
}

var file_control_v1_control_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_control_v1_control_proto_goTypes = []interface{}{
	(*RegisterTestRequest)(nil),   // 0: tstr.control.v1.RegisterTestRequest
	(*RegisterTestResponse)(nil),  // 1: tstr.control.v1.RegisterTestResponse
	(*UpdateTestRequest)(nil),     // 2: tstr.control.v1.UpdateTestRequest
	(*UpdateTestResponse)(nil),    // 3: tstr.control.v1.UpdateTestResponse
	(*DeleteTestRequest)(nil),     // 4: tstr.control.v1.DeleteTestRequest
	(*DeleteTestResponse)(nil),    // 5: tstr.control.v1.DeleteTestResponse
	(*ScheduleRunRequest)(nil),    // 6: tstr.control.v1.ScheduleRunRequest
	(*ScheduleRunResponse)(nil),   // 7: tstr.control.v1.ScheduleRunResponse
	nil,                           // 8: tstr.control.v1.RegisterTestRequest.LabelsEntry
	nil,                           // 9: tstr.control.v1.UpdateTestRequest.LabelsEntry
	nil,                           // 10: tstr.control.v1.ScheduleRunRequest.LabelsEntry
	(*v1.Test_RunConfig)(nil),     // 11: tstr.common.v1.Test.RunConfig
	(*v1.Test_Matrix)(nil),        // 12: tstr.common.v1.Test.Matrix
	(*v1.Test)(nil),               // 13: tstr.common.v1.Test
	(*fieldmaskpb.FieldMask)(nil), // 14: google.protobuf.FieldMask
	(*v1.Run)(nil),                // 15: tstr.common.v1.Run
}
var file_control_v1_control_proto_depIdxs = []int32{
	8,  // 0: tstr.control.v1.RegisterTestRequest.labels:type_name -> tstr.control.v1.RegisterTestRequest.LabelsEntry
	11, // 1: tstr.control.v1.RegisterTestRequest.run_config:type_name -> tstr.common.v1.Test.RunConfig
	12, // 2: tstr.control.v1.RegisterTestRequest.matrix:type_name -> tstr.common.v1.Test.Matrix
	13, // 3: tstr.control.v1.RegisterTestResponse.test:type_name -> tstr.common.v1.Test
	14, // 4: tstr.control.v1.UpdateTestRequest.field_mask:type_name -> google.protobuf.FieldMask
	11, // 5: tstr.control.v1.UpdateTestRequest.run_config:type_name -> tstr.common.v1.Test.RunConfig
	9,  // 6: tstr.control.v1.UpdateTestRequest.labels:type_name -> tstr.control.v1.UpdateTestRequest.LabelsEntry
	12, // 7: tstr.control.v1.UpdateTestRequest.matrix:type_name -> tstr.common.v1.Test.Matrix
	10, // 8: tstr.control.v1.ScheduleRunRequest.labels:type_name -> tstr.control.v1.ScheduleRunRequest.LabelsEntry
	12, // 9: tstr.control.v1.ScheduleRunRequest.test_matrix:type_name -> tstr.common.v1.Test.Matrix
	15, // 10: tstr.control.v1.ScheduleRunResponse.runs:type_name -> tstr.common.v1.Run
	0,  // 11: tstr.control.v1.ControlService.RegisterTest:input_type -> tstr.control.v1.RegisterTestRequest
	2,  // 12: tstr.control.v1.ControlService.UpdateTest:input_type -> tstr.control.v1.UpdateTestRequest
	4,  // 13: tstr.control.v1.ControlService.DeleteTest:input_type -> tstr.control.v1.DeleteTestRequest
	6,  // 14: tstr.control.v1.ControlService.ScheduleRun:input_type -> tstr.control.v1.ScheduleRunRequest
	1,  // 15: tstr.control.v1.ControlService.RegisterTest:output_type -> tstr.control.v1.RegisterTestResponse
	3,  // 16: tstr.control.v1.ControlService.UpdateTest:output_type -> tstr.control.v1.UpdateTestResponse
	5,  // 17: tstr.control.v1.ControlService.DeleteTest:output_type -> tstr.control.v1.DeleteTestResponse
	7,  // 18: tstr.control.v1.ControlService.ScheduleRun:output_type -> tstr.control.v1.ScheduleRunResponse
	15, // [15:19] is the sub-list for method output_type
	11, // [11:15] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_control_v1_control_proto_init() }
func file_control_v1_control_proto_init() {
	if File_control_v1_control_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_control_v1_control_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterTestRequest); i {
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
		file_control_v1_control_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterTestResponse); i {
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
		file_control_v1_control_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateTestRequest); i {
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
		file_control_v1_control_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateTestResponse); i {
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
		file_control_v1_control_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteTestRequest); i {
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
		file_control_v1_control_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteTestResponse); i {
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
		file_control_v1_control_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ScheduleRunRequest); i {
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
		file_control_v1_control_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ScheduleRunResponse); i {
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
			RawDescriptor: file_control_v1_control_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_control_v1_control_proto_goTypes,
		DependencyIndexes: file_control_v1_control_proto_depIdxs,
		MessageInfos:      file_control_v1_control_proto_msgTypes,
	}.Build()
	File_control_v1_control_proto = out.File
	file_control_v1_control_proto_rawDesc = nil
	file_control_v1_control_proto_goTypes = nil
	file_control_v1_control_proto_depIdxs = nil
}
