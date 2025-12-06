package etcd

type ServerNet struct {
	Addr *string `json:"addr"`           // e.g.: 127.0.0.1:8989 [default]: ""
	Name *string `json:"name,omitempty"` // 链接名称 [default]: ""
	Type *string `json:"type"`           // [tcp, kcp] [default]: xnetcommon.ServerNetTypeNameTCP
}

type GrpcService struct {
	PackageName *string `json:"packageName,omitempty"` // 包名
	ServiceName *string `json:"serviceName,omitempty"` // 服务名称
	Addr        *string `json:"addr,omitempty"`        // 服务地址 e.g.: 127.0.0.1:8989 [default]: ""
}
