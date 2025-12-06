package main

import (
	xgrpcproto "github.com/75912001/xlib/grpc/proto"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func findShardKeyField(message *protogen.Message) *protogen.Field {
	for _, field := range message.Fields {
		if isShardKeyField(field) {
			return field
		}
	}
	return nil
}

func isShardKeyField(field *protogen.Field) bool {
	// 检查字段是否有 shardKey 选项
	options := field.Desc.Options().(*descriptorpb.FieldOptions)
	if options == nil {
		return false
	}
	// 判断是否有结构化扩展 fieldOpt
	if proto.HasExtension(options, xgrpcproto.E_FieldOpt) {
		ext := proto.GetExtension(options, xgrpcproto.E_FieldOpt)
		if fieldOpt, ok := ext.(*xgrpcproto.FieldOpt); ok && fieldOpt != nil && fieldOpt.ShardKey {
			return true
		}
	}
	return false
}
