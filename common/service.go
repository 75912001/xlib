package common

type ServerNet struct {
	Addr *string `json:"addr"` // e.g.: 127.0.0.0:8989 [default]: ""
	Name *string `json:"name"` // 链接名称 [default]: ""
	Type *string `json:"type"` // [tcp, kcp] [default]: xnetcommon.ServerNetTypeNameTCP
}
