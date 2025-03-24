package common

type ServerNet struct {
	Addr *string `json:"addr"` // e.g.: 127.0.0.0:8989 [default]: ""
	Type *string `json:"type"` // [tcp, kcp] [default]: "tcp"
}
