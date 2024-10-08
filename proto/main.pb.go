// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.0
// source: proto/main.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Metric_Type int32

const (
	Metric_COUNTER Metric_Type = 0
	Metric_GAUGE   Metric_Type = 1
)

// Enum value maps for Metric_Type.
var (
	Metric_Type_name = map[int32]string{
		0: "COUNTER",
		1: "GAUGE",
	}
	Metric_Type_value = map[string]int32{
		"COUNTER": 0,
		"GAUGE":   1,
	}
)

func (x Metric_Type) Enum() *Metric_Type {
	p := new(Metric_Type)
	*p = x
	return p
}

func (x Metric_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Metric_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_main_proto_enumTypes[0].Descriptor()
}

func (Metric_Type) Type() protoreflect.EnumType {
	return &file_proto_main_proto_enumTypes[0]
}

func (x Metric_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Metric_Type.Descriptor instead.
func (Metric_Type) EnumDescriptor() ([]byte, []int) {
	return file_proto_main_proto_rawDescGZIP(), []int{0, 0}
}

type Metric struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string      `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Type  Metric_Type `protobuf:"varint,2,opt,name=type,proto3,enum=main.Metric_Type" json:"type,omitempty"`
	Delta *int64      `protobuf:"varint,3,opt,name=delta,proto3,oneof" json:"delta,omitempty"`
	Value *float64    `protobuf:"fixed64,4,opt,name=value,proto3,oneof" json:"value,omitempty"`
}

func (x *Metric) Reset() {
	*x = Metric{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_main_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metric) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metric) ProtoMessage() {}

func (x *Metric) ProtoReflect() protoreflect.Message {
	mi := &file_proto_main_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metric.ProtoReflect.Descriptor instead.
func (*Metric) Descriptor() ([]byte, []int) {
	return file_proto_main_proto_rawDescGZIP(), []int{0}
}

func (x *Metric) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Metric) GetType() Metric_Type {
	if x != nil {
		return x.Type
	}
	return Metric_COUNTER
}

func (x *Metric) GetDelta() int64 {
	if x != nil && x.Delta != nil {
		return *x.Delta
	}
	return 0
}

func (x *Metric) GetValue() float64 {
	if x != nil && x.Value != nil {
		return *x.Value
	}
	return 0
}

type CounterMetric struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metric *Metric `protobuf:"bytes,1,opt,name=metric,proto3" json:"metric,omitempty"`
	Delta  int64   `protobuf:"varint,2,opt,name=delta,proto3" json:"delta,omitempty"`
}

func (x *CounterMetric) Reset() {
	*x = CounterMetric{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_main_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CounterMetric) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CounterMetric) ProtoMessage() {}

func (x *CounterMetric) ProtoReflect() protoreflect.Message {
	mi := &file_proto_main_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CounterMetric.ProtoReflect.Descriptor instead.
func (*CounterMetric) Descriptor() ([]byte, []int) {
	return file_proto_main_proto_rawDescGZIP(), []int{1}
}

func (x *CounterMetric) GetMetric() *Metric {
	if x != nil {
		return x.Metric
	}
	return nil
}

func (x *CounterMetric) GetDelta() int64 {
	if x != nil {
		return x.Delta
	}
	return 0
}

type GaugeMetric struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metric *Metric `protobuf:"bytes,1,opt,name=metric,proto3" json:"metric,omitempty"`
	Value  float64 `protobuf:"fixed64,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *GaugeMetric) Reset() {
	*x = GaugeMetric{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_main_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GaugeMetric) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GaugeMetric) ProtoMessage() {}

func (x *GaugeMetric) ProtoReflect() protoreflect.Message {
	mi := &file_proto_main_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GaugeMetric.ProtoReflect.Descriptor instead.
func (*GaugeMetric) Descriptor() ([]byte, []int) {
	return file_proto_main_proto_rawDescGZIP(), []int{2}
}

func (x *GaugeMetric) GetMetric() *Metric {
	if x != nil {
		return x.Metric
	}
	return nil
}

func (x *GaugeMetric) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

type GetMetricRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metric *Metric `protobuf:"bytes,1,opt,name=metric,proto3" json:"metric,omitempty"`
}

func (x *GetMetricRequest) Reset() {
	*x = GetMetricRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_main_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMetricRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMetricRequest) ProtoMessage() {}

func (x *GetMetricRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_main_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMetricRequest.ProtoReflect.Descriptor instead.
func (*GetMetricRequest) Descriptor() ([]byte, []int) {
	return file_proto_main_proto_rawDescGZIP(), []int{3}
}

func (x *GetMetricRequest) GetMetric() *Metric {
	if x != nil {
		return x.Metric
	}
	return nil
}

type GetMetricResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metric *Metric `protobuf:"bytes,1,opt,name=metric,proto3" json:"metric,omitempty"`
	Error  string  `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *GetMetricResponse) Reset() {
	*x = GetMetricResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_main_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMetricResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMetricResponse) ProtoMessage() {}

func (x *GetMetricResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_main_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMetricResponse.ProtoReflect.Descriptor instead.
func (*GetMetricResponse) Descriptor() ([]byte, []int) {
	return file_proto_main_proto_rawDescGZIP(), []int{4}
}

func (x *GetMetricResponse) GetMetric() *Metric {
	if x != nil {
		return x.Metric
	}
	return nil
}

func (x *GetMetricResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type UpdateCounterRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metric *CounterMetric `protobuf:"bytes,1,opt,name=metric,proto3" json:"metric,omitempty"`
}

func (x *UpdateCounterRequest) Reset() {
	*x = UpdateCounterRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_main_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateCounterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateCounterRequest) ProtoMessage() {}

func (x *UpdateCounterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_main_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateCounterRequest.ProtoReflect.Descriptor instead.
func (*UpdateCounterRequest) Descriptor() ([]byte, []int) {
	return file_proto_main_proto_rawDescGZIP(), []int{5}
}

func (x *UpdateCounterRequest) GetMetric() *CounterMetric {
	if x != nil {
		return x.Metric
	}
	return nil
}

type UpdateCounterResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metric *CounterMetric `protobuf:"bytes,1,opt,name=metric,proto3" json:"metric,omitempty"`
	Error  string         `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *UpdateCounterResponse) Reset() {
	*x = UpdateCounterResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_main_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateCounterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateCounterResponse) ProtoMessage() {}

func (x *UpdateCounterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_main_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateCounterResponse.ProtoReflect.Descriptor instead.
func (*UpdateCounterResponse) Descriptor() ([]byte, []int) {
	return file_proto_main_proto_rawDescGZIP(), []int{6}
}

func (x *UpdateCounterResponse) GetMetric() *CounterMetric {
	if x != nil {
		return x.Metric
	}
	return nil
}

func (x *UpdateCounterResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type UpdateGaugeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metric *GaugeMetric `protobuf:"bytes,1,opt,name=metric,proto3" json:"metric,omitempty"`
}

func (x *UpdateGaugeRequest) Reset() {
	*x = UpdateGaugeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_main_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateGaugeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateGaugeRequest) ProtoMessage() {}

func (x *UpdateGaugeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_main_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateGaugeRequest.ProtoReflect.Descriptor instead.
func (*UpdateGaugeRequest) Descriptor() ([]byte, []int) {
	return file_proto_main_proto_rawDescGZIP(), []int{7}
}

func (x *UpdateGaugeRequest) GetMetric() *GaugeMetric {
	if x != nil {
		return x.Metric
	}
	return nil
}

type UpdateGaugeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metric *GaugeMetric `protobuf:"bytes,1,opt,name=metric,proto3" json:"metric,omitempty"`
	Error  string       `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *UpdateGaugeResponse) Reset() {
	*x = UpdateGaugeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_main_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateGaugeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateGaugeResponse) ProtoMessage() {}

func (x *UpdateGaugeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_main_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateGaugeResponse.ProtoReflect.Descriptor instead.
func (*UpdateGaugeResponse) Descriptor() ([]byte, []int) {
	return file_proto_main_proto_rawDescGZIP(), []int{8}
}

func (x *UpdateGaugeResponse) GetMetric() *GaugeMetric {
	if x != nil {
		return x.Metric
	}
	return nil
}

func (x *UpdateGaugeResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type UpdateMetricsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metrics []*Metric `protobuf:"bytes,1,rep,name=metrics,proto3" json:"metrics,omitempty"`
}

func (x *UpdateMetricsRequest) Reset() {
	*x = UpdateMetricsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_main_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateMetricsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateMetricsRequest) ProtoMessage() {}

func (x *UpdateMetricsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_main_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateMetricsRequest.ProtoReflect.Descriptor instead.
func (*UpdateMetricsRequest) Descriptor() ([]byte, []int) {
	return file_proto_main_proto_rawDescGZIP(), []int{9}
}

func (x *UpdateMetricsRequest) GetMetrics() []*Metric {
	if x != nil {
		return x.Metrics
	}
	return nil
}

type UpdateMetricsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metrics []*Metric `protobuf:"bytes,1,rep,name=metrics,proto3" json:"metrics,omitempty"`
	Error   string    `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *UpdateMetricsResponse) Reset() {
	*x = UpdateMetricsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_main_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateMetricsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateMetricsResponse) ProtoMessage() {}

func (x *UpdateMetricsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_main_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateMetricsResponse.ProtoReflect.Descriptor instead.
func (*UpdateMetricsResponse) Descriptor() ([]byte, []int) {
	return file_proto_main_proto_rawDescGZIP(), []int{10}
}

func (x *UpdateMetricsResponse) GetMetrics() []*Metric {
	if x != nil {
		return x.Metrics
	}
	return nil
}

func (x *UpdateMetricsResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_proto_main_proto protoreflect.FileDescriptor

var file_proto_main_proto_rawDesc = []byte{
	0x0a, 0x10, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x04, 0x6d, 0x61, 0x69, 0x6e, 0x22, 0xa9, 0x01, 0x0a, 0x06, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x25, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x11, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x2e,
	0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x19, 0x0a, 0x05, 0x64, 0x65,
	0x6c, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x48, 0x00, 0x52, 0x05, 0x64, 0x65, 0x6c,
	0x74, 0x61, 0x88, 0x01, 0x01, 0x12, 0x19, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x01, 0x48, 0x01, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x88, 0x01, 0x01,
	0x22, 0x1e, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x43, 0x4f, 0x55, 0x4e,
	0x54, 0x45, 0x52, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x47, 0x41, 0x55, 0x47, 0x45, 0x10, 0x01,
	0x42, 0x08, 0x0a, 0x06, 0x5f, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x22, 0x4b, 0x0a, 0x0d, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x4d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x24, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x14, 0x0a, 0x05, 0x64,
	0x65, 0x6c, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x64, 0x65, 0x6c, 0x74,
	0x61, 0x22, 0x49, 0x0a, 0x0b, 0x47, 0x61, 0x75, 0x67, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x12, 0x24, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0c, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x06,
	0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x38, 0x0a, 0x10,
	0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x24, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0c, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x06,
	0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x22, 0x4f, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x24, 0x0a, 0x06, 0x6d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x6d, 0x61,
	0x69, 0x6e, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x43, 0x0a, 0x14, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x2b, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x13, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x4d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x22, 0x5a, 0x0a, 0x15,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2b, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x43, 0x6f, 0x75,
	0x6e, 0x74, 0x65, 0x72, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x3f, 0x0a, 0x12, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x47, 0x61, 0x75, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x29,
	0x0a, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11,
	0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x47, 0x61, 0x75, 0x67, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x22, 0x56, 0x0a, 0x13, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x47, 0x61, 0x75, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x29, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x11, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x47, 0x61, 0x75, 0x67, 0x65, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x14, 0x0a, 0x05, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x22, 0x3e, 0x0a, 0x14, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x26, 0x0a, 0x07, 0x6d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x6d, 0x61, 0x69,
	0x6e, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x73, 0x22, 0x55, 0x0a, 0x15, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x26, 0x0a, 0x07, 0x6d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x6d, 0x61,
	0x69, 0x6e, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x32, 0x9f, 0x02, 0x0a, 0x07, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x73, 0x12, 0x3c, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x12, 0x16, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x6d, 0x61, 0x69, 0x6e,
	0x2e, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x48, 0x0a, 0x0d, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x75, 0x6e,
	0x74, 0x65, 0x72, 0x12, 0x1a, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1b, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x75,
	0x6e, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x42, 0x0a, 0x0b,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x47, 0x61, 0x75, 0x67, 0x65, 0x12, 0x18, 0x2e, 0x6d, 0x61,
	0x69, 0x6e, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x47, 0x61, 0x75, 0x67, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x47, 0x61, 0x75, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x48, 0x0a, 0x0d, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x73, 0x12, 0x1a, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e,
	0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0c, 0x5a, 0x0a, 0x6d, 0x61,
	0x69, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_main_proto_rawDescOnce sync.Once
	file_proto_main_proto_rawDescData = file_proto_main_proto_rawDesc
)

func file_proto_main_proto_rawDescGZIP() []byte {
	file_proto_main_proto_rawDescOnce.Do(func() {
		file_proto_main_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_main_proto_rawDescData)
	})
	return file_proto_main_proto_rawDescData
}

var file_proto_main_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_main_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_proto_main_proto_goTypes = []any{
	(Metric_Type)(0),              // 0: main.Metric.Type
	(*Metric)(nil),                // 1: main.Metric
	(*CounterMetric)(nil),         // 2: main.CounterMetric
	(*GaugeMetric)(nil),           // 3: main.GaugeMetric
	(*GetMetricRequest)(nil),      // 4: main.GetMetricRequest
	(*GetMetricResponse)(nil),     // 5: main.GetMetricResponse
	(*UpdateCounterRequest)(nil),  // 6: main.UpdateCounterRequest
	(*UpdateCounterResponse)(nil), // 7: main.UpdateCounterResponse
	(*UpdateGaugeRequest)(nil),    // 8: main.UpdateGaugeRequest
	(*UpdateGaugeResponse)(nil),   // 9: main.UpdateGaugeResponse
	(*UpdateMetricsRequest)(nil),  // 10: main.UpdateMetricsRequest
	(*UpdateMetricsResponse)(nil), // 11: main.UpdateMetricsResponse
}
var file_proto_main_proto_depIdxs = []int32{
	0,  // 0: main.Metric.type:type_name -> main.Metric.Type
	1,  // 1: main.CounterMetric.metric:type_name -> main.Metric
	1,  // 2: main.GaugeMetric.metric:type_name -> main.Metric
	1,  // 3: main.GetMetricRequest.metric:type_name -> main.Metric
	1,  // 4: main.GetMetricResponse.metric:type_name -> main.Metric
	2,  // 5: main.UpdateCounterRequest.metric:type_name -> main.CounterMetric
	2,  // 6: main.UpdateCounterResponse.metric:type_name -> main.CounterMetric
	3,  // 7: main.UpdateGaugeRequest.metric:type_name -> main.GaugeMetric
	3,  // 8: main.UpdateGaugeResponse.metric:type_name -> main.GaugeMetric
	1,  // 9: main.UpdateMetricsRequest.metrics:type_name -> main.Metric
	1,  // 10: main.UpdateMetricsResponse.metrics:type_name -> main.Metric
	4,  // 11: main.Metrics.GetMetric:input_type -> main.GetMetricRequest
	6,  // 12: main.Metrics.UpdateCounter:input_type -> main.UpdateCounterRequest
	8,  // 13: main.Metrics.UpdateGauge:input_type -> main.UpdateGaugeRequest
	10, // 14: main.Metrics.UpdateMetrics:input_type -> main.UpdateMetricsRequest
	5,  // 15: main.Metrics.GetMetric:output_type -> main.GetMetricResponse
	7,  // 16: main.Metrics.UpdateCounter:output_type -> main.UpdateCounterResponse
	9,  // 17: main.Metrics.UpdateGauge:output_type -> main.UpdateGaugeResponse
	11, // 18: main.Metrics.UpdateMetrics:output_type -> main.UpdateMetricsResponse
	15, // [15:19] is the sub-list for method output_type
	11, // [11:15] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_proto_main_proto_init() }
func file_proto_main_proto_init() {
	if File_proto_main_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_main_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Metric); i {
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
		file_proto_main_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*CounterMetric); i {
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
		file_proto_main_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*GaugeMetric); i {
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
		file_proto_main_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*GetMetricRequest); i {
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
		file_proto_main_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*GetMetricResponse); i {
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
		file_proto_main_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*UpdateCounterRequest); i {
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
		file_proto_main_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*UpdateCounterResponse); i {
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
		file_proto_main_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*UpdateGaugeRequest); i {
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
		file_proto_main_proto_msgTypes[8].Exporter = func(v any, i int) any {
			switch v := v.(*UpdateGaugeResponse); i {
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
		file_proto_main_proto_msgTypes[9].Exporter = func(v any, i int) any {
			switch v := v.(*UpdateMetricsRequest); i {
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
		file_proto_main_proto_msgTypes[10].Exporter = func(v any, i int) any {
			switch v := v.(*UpdateMetricsResponse); i {
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
	file_proto_main_proto_msgTypes[0].OneofWrappers = []any{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_main_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_main_proto_goTypes,
		DependencyIndexes: file_proto_main_proto_depIdxs,
		EnumInfos:         file_proto_main_proto_enumTypes,
		MessageInfos:      file_proto_main_proto_msgTypes,
	}.Build()
	File_proto_main_proto = out.File
	file_proto_main_proto_rawDesc = nil
	file_proto_main_proto_goTypes = nil
	file_proto_main_proto_depIdxs = nil
}
