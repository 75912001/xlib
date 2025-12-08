package registry

import (
	xerror "github.com/75912001/xlib/error"
	xgrpcproto "github.com/75912001/xlib/grpc/proto"
	xgrpcutil "github.com/75912001/xlib/grpc/util"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"time"
)

var GMethodOptions map[string]*xgrpcproto.MethodOpt // 方法配置, 未配置方法的, 继承服务的配置 key: /${packageName}.${serviceName}/${methodName}
var packageServiceMap map[string][]string           // key: /${packageName}.${serviceName}  value: method slice

func Init() {
	// 初始化配置
	serviceOptions := make(map[string]*xgrpcproto.ServiceOpt) // 服务配置 key: /${packageName}.${serviceName}
	GMethodOptions = make(map[string]*xgrpcproto.MethodOpt)
	packageServiceMap = make(map[string][]string)

	files := protoregistry.GlobalFiles
	// 获取服务的配置
	files.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		services := fd.Services()
		for i := 0; i < services.Len(); i++ {
			svc := services.Get(i)
			svcName := string(svc.FullName())
			// 服务级别的配置
			opts := svc.Options().(*descriptorpb.ServiceOptions)
			if opts == nil { // 服务没有选项
				continue
			}
			ext, ok := proto.GetExtension(opts, xgrpcproto.E_ServiceOpt).(*xgrpcproto.ServiceOpt)
			if !ok { // 服务没有选项
				panic(errors.WithMessagef(xerror.Configure, "service %v options not exist", svcName))
			}
			if ext.Timeout == "" { // 服务没有设置超时时间
				panic(errors.WithMessagef(xerror.Configure, "service %v timeout not exist", svcName))
			}
			duration, err := time.ParseDuration(ext.Timeout)
			if err != nil { // 服务设置的超时时间格式错误
				panic(errors.WithMessagef(xerror.Configure, "service %v timeout format error %v", svcName, ext.Timeout))
			}
			_ = duration
			if ext.LoadBalancePolicy == xgrpcproto.LoadBalancePolicy_LoadBalancePolicy_Unspecified { // 服务没有设置负载均衡策略
				panic(errors.WithMessagef(xerror.Configure, "service %v load balance policy not set", svcName))
			}
			serviceOptions["/"+svcName] = ext
		}
		return true
	})
	// 获取方法的配置
	files.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		services := fd.Services()
		for i := 0; i < services.Len(); i++ {
			svc := services.Get(i)
			svcName := string(svc.FullName())
			// 方法级别的配置
			methods := svc.Methods()
			for j := 0; j < methods.Len(); j++ {
				method := methods.Get(j)
				serviceName := "/" + svcName
				if _, ok := serviceOptions[serviceName]; !ok { // 服务没有配置
					continue
				}
				serviceOpts := serviceOptions[serviceName]
				newOpt := xgrpcproto.MethodOpt{ // 先使用服务的选项
					Timeout:           serviceOpts.Timeout,
					LoadBalancePolicy: serviceOpts.LoadBalancePolicy,
					ShardKeyFieldType: serviceOpts.ShardKeyFieldType,
				}
				methodName := serviceName + "/" + string(method.Name())
				if !method.IsStreamingClient() && !method.IsStreamingServer() { // unary
					addMethod(serviceName, methodName)
					// 获取 request 中的 shardKey 的类型
					request := method.Input()
					fields := request.Fields()
					for k := 0; k < fields.Len(); k++ {
						field := fields.Get(k)
						fieldOpts := field.Options().(*descriptorpb.FieldOptions)
						if fieldOpts == nil {
							continue
						}
						if !proto.HasExtension(fieldOpts, xgrpcproto.E_FieldOpt) {
							continue
						}
						switch field.Kind() {
						case protoreflect.StringKind:
							newOpt.ShardKeyFieldType = xgrpcproto.ShardKeyFieldType_ShardKeyFieldType_STRING
						case protoreflect.Int32Kind:
							newOpt.ShardKeyFieldType = xgrpcproto.ShardKeyFieldType_ShardKeyFieldType_INT32
						case protoreflect.Int64Kind:
							newOpt.ShardKeyFieldType = xgrpcproto.ShardKeyFieldType_ShardKeyFieldType_INT64
						case protoreflect.Uint32Kind:
							newOpt.ShardKeyFieldType = xgrpcproto.ShardKeyFieldType_ShardKeyFieldType_UINT32
						case protoreflect.Uint64Kind:
							newOpt.ShardKeyFieldType = xgrpcproto.ShardKeyFieldType_ShardKeyFieldType_UINT64
						default:
							panic(errors.WithMessagef(xerror.Configure, "method name:%v shard key field type not supported", methodName))
						}
						break // 找到第一个 shard key 字段, 直接使用
					}
					methodOpts := method.Options().(*descriptorpb.MethodOptions)
					if methodOpts != nil {
						ext, ok := proto.GetExtension(methodOpts, xgrpcproto.E_MethodOpt).(*xgrpcproto.MethodOpt)
						if ok && ext != nil { // 方法-选项存在
							if duration, err := time.ParseDuration(ext.Timeout); err == nil {
								_ = duration
								newOpt.Timeout = ext.Timeout
							}
							if ext.LoadBalancePolicy != xgrpcproto.LoadBalancePolicy_LoadBalancePolicy_Unspecified {
								newOpt.LoadBalancePolicy = ext.LoadBalancePolicy
							}
						}
					}
					GMethodOptions[methodName] = &newOpt
				} else if method.IsStreamingClient() && method.IsStreamingServer() { // stream - 双向流
				} else {
					// 不是 unary 也不是双向 stream
					panic(errors.WithMessagef(xerror.Configure, "method %v is not unary or stream", methodName))
				}
			}
		}
		return true
	})
}

// GetOptions 获取指定方法的配置
// method 格式为 "/${packageName}.${serviceName}/${methodName}"
func GetOptions(method string) *xgrpcproto.MethodOpt {
	// 先查找方法级别的配置
	if value, ok := GMethodOptions[method]; ok {
		return value
	}
	return &xgrpcproto.MethodOpt{
		Timeout:           xgrpcproto.RpcTimeoutDefault,
		LoadBalancePolicy: xgrpcproto.LoadBalancePolicy_LoadBalancePolicy_Mod,
		ShardKeyFieldType: xgrpcproto.ShardKeyFieldType_ShardKeyFieldType_Unspecified,
	}
}

// GetMethodSlice 获取 method slice
func GetMethodSlice(packageName string, serviceName string) []string {
	packageServiceName := xgrpcutil.GenPackageServiceName(packageName, serviceName)
	methodSlice := make([]string, 0)
	if _, ok := packageServiceMap[packageServiceName]; !ok {
		return methodSlice
	}
	methodSlice = append(methodSlice, packageServiceMap[packageServiceName]...)
	return methodSlice
}

func addMethod(packageServiceName string, method string) {
	if _, ok := packageServiceMap[packageServiceName]; !ok {
		packageServiceMap[packageServiceName] = make([]string, 0)
	}
	packageServiceMap[packageServiceName] = append(packageServiceMap[packageServiceName], method)
}
