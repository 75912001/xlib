package main

import (
	xgrpcproto "github.com/75912001/xlib/grpc/proto"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// findShardKeyField 查找请求消息中第一个标记为 shardKey 的字段。
// Mod/RingHash 生成代码依赖该字段计算目标服务连接。
func findShardKeyField(message *protogen.Message) *protogen.Field {
	for _, field := range message.Fields {
		if isShardKeyField(field) {
			return field
		}
	}
	return nil
}

// isShardKeyField 判断字段是否设置了 proto.fieldOpt.shardKey。
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

// isDirectMethod 判断 unary RPC 是否使用 Direct 策略。
// Direct 表示调用方已经持有目标连接，生成代码不再经过 selector。
func isDirectMethod(service *protogen.Service, method *protogen.Method) bool {
	return effectiveLoadBalancePolicy(service, method) == xgrpcproto.LoadBalancePolicy_LoadBalancePolicy_Direct
}

// isNoShardSelectorMethod 判断 unary RPC 是否使用无 shardKey 的选服策略。
// First/Random/RoundRobin 不需要请求字段提供 shardKey，但仍需要 selector 从服务发现连接中选出目标连接。
func isNoShardSelectorMethod(service *protogen.Service, method *protogen.Method) bool {
	switch effectiveLoadBalancePolicy(service, method) {
	case xgrpcproto.LoadBalancePolicy_LoadBalancePolicy_First,
		xgrpcproto.LoadBalancePolicy_LoadBalancePolicy_Random,
		xgrpcproto.LoadBalancePolicy_LoadBalancePolicy_RoundRobin:
		return true
	default:
		return false
	}
}

// effectiveLoadBalancePolicy 返回方法最终生效的负载均衡策略。
// 方法级配置优先生效；方法未显式配置时继承服务级配置。
func effectiveLoadBalancePolicy(service *protogen.Service, method *protogen.Method) xgrpcproto.LoadBalancePolicy {
	policy := xgrpcproto.LoadBalancePolicy_LoadBalancePolicy_Unspecified
	if serviceOptions, ok := service.Desc.Options().(*descriptorpb.ServiceOptions); ok && serviceOptions != nil {
		if proto.HasExtension(serviceOptions, xgrpcproto.E_ServiceOpt) {
			ext := proto.GetExtension(serviceOptions, xgrpcproto.E_ServiceOpt)
			if serviceOpt, ok := ext.(*xgrpcproto.ServiceOpt); ok && serviceOpt != nil {
				policy = serviceOpt.LoadBalancePolicy
			}
		}
	}
	if methodOptions, ok := method.Desc.Options().(*descriptorpb.MethodOptions); ok && methodOptions != nil {
		if proto.HasExtension(methodOptions, xgrpcproto.E_MethodOpt) {
			ext := proto.GetExtension(methodOptions, xgrpcproto.E_MethodOpt)
			if methodOpt, ok := ext.(*xgrpcproto.MethodOpt); ok && methodOpt != nil &&
				methodOpt.LoadBalancePolicy != xgrpcproto.LoadBalancePolicy_LoadBalancePolicy_Unspecified {
				policy = methodOpt.LoadBalancePolicy
			}
		}
	}
	return policy
}
