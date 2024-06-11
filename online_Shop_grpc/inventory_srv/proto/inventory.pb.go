// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.19.4
// source: order.proto

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

type MyEmpty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *MyEmpty) Reset() {
	*x = MyEmpty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_inventory_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MyEmpty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MyEmpty) ProtoMessage() {}

func (x *MyEmpty) ProtoReflect() protoreflect.Message {
	mi := &file_inventory_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MyEmpty.ProtoReflect.Descriptor instead.
func (*MyEmpty) Descriptor() ([]byte, []int) {
	return file_inventory_proto_rawDescGZIP(), []int{0}
}

type GoodsInvInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GoodsId int32 `protobuf:"varint,1,opt,name=goodsId,proto3" json:"goodsId,omitempty"`
	Num     int32 `protobuf:"varint,2,opt,name=num,proto3" json:"num,omitempty"`
}

func (x *GoodsInvInfo) Reset() {
	*x = GoodsInvInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_inventory_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GoodsInvInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GoodsInvInfo) ProtoMessage() {}

func (x *GoodsInvInfo) ProtoReflect() protoreflect.Message {
	mi := &file_inventory_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GoodsInvInfo.ProtoReflect.Descriptor instead.
func (*GoodsInvInfo) Descriptor() ([]byte, []int) {
	return file_inventory_proto_rawDescGZIP(), []int{1}
}

func (x *GoodsInvInfo) GetGoodsId() int32 {
	if x != nil {
		return x.GoodsId
	}
	return 0
}

func (x *GoodsInvInfo) GetNum() int32 {
	if x != nil {
		return x.Num
	}
	return 0
}

type SellInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GoodsInfo []*GoodsInvInfo `protobuf:"bytes,1,rep,name=goodsInfo,proto3" json:"goodsInfo,omitempty"`
}

func (x *SellInfo) Reset() {
	*x = SellInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_inventory_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SellInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SellInfo) ProtoMessage() {}

func (x *SellInfo) ProtoReflect() protoreflect.Message {
	mi := &file_inventory_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SellInfo.ProtoReflect.Descriptor instead.
func (*SellInfo) Descriptor() ([]byte, []int) {
	return file_inventory_proto_rawDescGZIP(), []int{2}
}

func (x *SellInfo) GetGoodsInfo() []*GoodsInvInfo {
	if x != nil {
		return x.GoodsInfo
	}
	return nil
}

var File_inventory_proto protoreflect.FileDescriptor

var file_inventory_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x09, 0x0a, 0x07, 0x4d, 0x79, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x3a, 0x0a, 0x0c,
	0x47, 0x6f, 0x6f, 0x64, 0x73, 0x49, 0x6e, 0x76, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x18, 0x0a, 0x07,
	0x67, 0x6f, 0x6f, 0x64, 0x73, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x67,
	0x6f, 0x6f, 0x64, 0x73, 0x49, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x6e, 0x75, 0x6d, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x03, 0x6e, 0x75, 0x6d, 0x22, 0x37, 0x0a, 0x08, 0x53, 0x65, 0x6c, 0x6c,
	0x49, 0x6e, 0x66, 0x6f, 0x12, 0x2b, 0x0a, 0x09, 0x67, 0x6f, 0x6f, 0x64, 0x73, 0x49, 0x6e, 0x66,
	0x6f, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x47, 0x6f, 0x6f, 0x64, 0x73, 0x49,
	0x6e, 0x76, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x09, 0x67, 0x6f, 0x6f, 0x64, 0x73, 0x49, 0x6e, 0x66,
	0x6f, 0x32, 0x95, 0x01, 0x0a, 0x09, 0x49, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x12,
	0x21, 0x0a, 0x06, 0x53, 0x65, 0x74, 0x49, 0x6e, 0x76, 0x12, 0x0d, 0x2e, 0x47, 0x6f, 0x6f, 0x64,
	0x73, 0x49, 0x6e, 0x76, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x08, 0x2e, 0x4d, 0x79, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x12, 0x29, 0x0a, 0x09, 0x49, 0x6e, 0x76, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x12,
	0x0d, 0x2e, 0x47, 0x6f, 0x6f, 0x64, 0x73, 0x49, 0x6e, 0x76, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x0d,
	0x2e, 0x47, 0x6f, 0x6f, 0x64, 0x73, 0x49, 0x6e, 0x76, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1b, 0x0a,
	0x04, 0x53, 0x65, 0x6c, 0x6c, 0x12, 0x09, 0x2e, 0x53, 0x65, 0x6c, 0x6c, 0x49, 0x6e, 0x66, 0x6f,
	0x1a, 0x08, 0x2e, 0x4d, 0x79, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x1d, 0x0a, 0x06, 0x52, 0x65,
	0x62, 0x61, 0x63, 0x6b, 0x12, 0x09, 0x2e, 0x53, 0x65, 0x6c, 0x6c, 0x49, 0x6e, 0x66, 0x6f, 0x1a,
	0x08, 0x2e, 0x4d, 0x79, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x3b, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_inventory_proto_rawDescOnce sync.Once
	file_inventory_proto_rawDescData = file_inventory_proto_rawDesc
)

func file_inventory_proto_rawDescGZIP() []byte {
	file_inventory_proto_rawDescOnce.Do(func() {
		file_inventory_proto_rawDescData = protoimpl.X.CompressGZIP(file_inventory_proto_rawDescData)
	})
	return file_inventory_proto_rawDescData
}

var file_inventory_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_inventory_proto_goTypes = []interface{}{
	(*MyEmpty)(nil),      // 0: MyEmpty
	(*GoodsInvInfo)(nil), // 1: GoodsInvInfo
	(*SellInfo)(nil),     // 2: SellInfo
}
var file_inventory_proto_depIdxs = []int32{
	1, // 0: SellInfo.goodsInfo:type_name -> GoodsInvInfo
	1, // 1: Inventory.SetInv:input_type -> GoodsInvInfo
	1, // 2: Inventory.InvDetail:input_type -> GoodsInvInfo
	2, // 3: Inventory.Sell:input_type -> SellInfo
	2, // 4: Inventory.Reback:input_type -> SellInfo
	0, // 5: Inventory.SetInv:output_type -> MyEmpty
	1, // 6: Inventory.InvDetail:output_type -> GoodsInvInfo
	0, // 7: Inventory.Sell:output_type -> MyEmpty
	0, // 8: Inventory.Reback:output_type -> MyEmpty
	5, // [5:9] is the sub-list for method output_type
	1, // [1:5] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_inventory_proto_init() }
func file_inventory_proto_init() {
	if File_inventory_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_inventory_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MyEmpty); i {
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
		file_inventory_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GoodsInvInfo); i {
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
		file_inventory_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SellInfo); i {
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
			RawDescriptor: file_inventory_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_inventory_proto_goTypes,
		DependencyIndexes: file_inventory_proto_depIdxs,
		MessageInfos:      file_inventory_proto_msgTypes,
	}.Build()
	File_inventory_proto = out.File
	file_inventory_proto_rawDesc = nil
	file_inventory_proto_goTypes = nil
	file_inventory_proto_depIdxs = nil
}
