package util

import "time"

const (
	ConnectTimeoutDurationDefault = 5 * time.Second // 默认连接超时时间为5秒
)

func GenPackageServiceName(packageName string, serviceName string) string {
	return "/" + packageName + "." + serviceName
}

func GenPackageServiceMethodName(packageName string, serviceName string, methodName string) string {
	return "/" + packageName + "." + serviceName + "/" + methodName
}
