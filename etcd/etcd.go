package etcd

import (
	"context"
	etcdclientv3 "go.etcd.io/etcd/client/v3"
)

var GEtcd IEtcd

type IEtcd interface {
	Start(ctx context.Context, value string) error
	Stop() error
	GetKey() string
	PutWithLease(key string, value string) (*etcdclientv3.PutResponse, error)
	KeepAlive(ctx context.Context) error
}
