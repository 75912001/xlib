package main

import "strings"

// X[service]
func getXService(serviceName string) string {
	return "X" + serviceName
}

// NewX[service]
func getNewXService(serviceName string) string {
	return "NewX" + serviceName
}

// X[service]Client
func getXServiceClient(serviceName string) string {
	return "X" + serviceName + "Client"
}

// NewX[service]Client
func getNewXServiceClient(serviceName string) string {
	return "NewX" + serviceName + "Client"
}

// [service]Client
func getServiceClient(serviceName string) string {
	return serviceName + "Client"
}

// New[service]Client
func getNewServiceClient(serviceName string) string {
	return "New" + serviceName + "Client"
}

// XStream[service][method]Client
func getXStreamServiceMethodClient(serviceName string, methodName string) string {
	return "XStream" + serviceName + methodName + "Client"
}

// [service]_[method]Client
func getService_MethodClient(serviceName string, methodName string) string {
	return serviceName + "_" + methodName + "Client"
}

// NewXStream[service][method]Client
func getNewXStreamServiceMethodClient(serviceName string, methodName string) string {
	return "NewXStream" + serviceName + methodName + "Client"
}

// IStream[service][method]Client
func getIStreamServiceMethodClient(serviceName string, methodName string) string {
	return "IStream" + serviceName + methodName + "Client"
}

// IStream[service]Server
func getIStreamServiceServer(serviceName string) string {
	return "IStream" + serviceName + "Server"
}

// [service]MessageWrapper
func getServiceMessageWrapper(serviceName string) string {
	return serviceName + "MessageWrapper"
}

// [service]Message
func getServiceMessage(serviceName string) string {
	return serviceName + "Message"
}

// X[service]Server
func getXServiceServer(serviceName string) string {
	return "X" + serviceName + "Server"
}

// NewX[service]Server
func getNewXServiceServer(serviceName string) string {
	return "NewX" + serviceName + "Server"
}

// [service]_[method]Server
func getService_MethodServer(serviceName string, methodName string) string {
	return serviceName + "_" + methodName + "Server"
}

// IUnary[service]Server
func getIUnaryServiceServer(serviceName string) string {
	return "IUnary" + serviceName + "Server"
}

// /////////////////////////////////////////////////////////////
// able field
func ableField() string {
	return "able"
}

// conn field
func connField() string {
	return "conn"
}

// GX[service]Service
func GServiceClientField(serviceName string) string {
	return "GX" + serviceName + "Service"
}

// 首字母小写
func lowerFirst(str string) string {
	return strings.ToLower(str[:1]) + str[1:]
}
