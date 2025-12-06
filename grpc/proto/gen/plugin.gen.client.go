package main

import (
	"google.golang.org/protobuf/compiler/protogen"
	"strings"
)

func genClient(g *protogen.GeneratedFile, file *protogen.File) {
	// 为每个服务生成扩展客户端
	for _, service := range file.Services {
		sn := service.GoName

		g.P("////////////////////////////////////////////////////////////////////////////////////////////////////")
		g.P("// ", getXServiceClient(sn), " 客户端")
		g.P("////////////////////////////////////////////////////////////////////////////////////////////////////")

		// 全局变量
		g.P("var (")
		g.P("\t", GServiceClientField(sn), " *", getXServiceClient(sn), " = &", getXServiceClient(sn), "{}")
		g.P(")")
		g.P()

		g.P("type ", getXServiceClient(sn), " struct {")
		g.P("\tClient ", getServiceClient(sn))
		g.P("}")
		g.P()

		g.P("func ", getNewXServiceClient(sn), "(clientConn *", grpcPackage.Ident("ClientConn"), ")(*", getXServiceClient(sn), ") {")
		g.P("\t", strings.ToLower(sn), "Client := &", getXServiceClient(sn), "{")
		g.P("\t\tClient: ", getNewServiceClient(sn), "(clientConn),")
		g.P("\t}")
		g.P("\treturn ", strings.ToLower(sn), "Client")
		g.P("}")

		g.P("////////////////////////////////////////////////////////////////////////////////")
		g.P("// ", sn, " 客户端-Stream")
		g.P("////////////////////////////////////////////////////////////////////////////////")

		for _, method := range service.Methods {
			if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() { // stream
				clientStreamGenerateMethod(g, service, method, getXServiceClient(sn))
			}
		}
		for _, method := range service.Methods {
			if !method.Desc.IsStreamingClient() && !method.Desc.IsStreamingServer() { // unary
				clientUnaryGenerateMethod(g, service, method, getXServiceClient(sn))
			}
		}
	}
}

func clientStreamGenerateMethod(g *protogen.GeneratedFile, service *protogen.Service, method *protogen.Method, clientName string) {
	sn := service.GoName
	mn := method.GoName
	if method.Desc.IsStreamingClient() && method.Desc.IsStreamingServer() {
		g.P("type ", getIStreamServiceMethodClient(sn, mn), " interface {")
		g.P("\t", mn, "Pre(stream ", getService_MethodClient(sn, mn), ") error // 预处理-新创建stream")
		g.P("\t", mn, "(client *", getXStreamServiceMethodClient(sn, mn), ", messageWrapper *",
			getServiceMessageWrapper(sn), ", stream ", getService_MethodClient(sn, mn), ") error // 处理")
		g.P("\t", mn, "Post(stream ", getService_MethodClient(sn, mn), ") error // 后处理-关闭stream")
		g.P("}")
		g.P()

		g.P("var ", lowerFirst(getIStreamServiceMethodClient(sn, mn)), " ", getIStreamServiceMethodClient(sn, mn))
		g.P()

		g.P("func Set", getIStreamServiceMethodClient(sn, mn), "(streamClient ", getIStreamServiceMethodClient(sn, mn), ") {")
		g.P("\t", lowerFirst(getIStreamServiceMethodClient(sn, mn)), " = streamClient")
		g.P("}")
		g.P()

		g.P("type ", getXStreamServiceMethodClient(sn, mn), " struct {")
		g.P("\t", getService_MethodClient(sn, mn), " ", getService_MethodClient(sn, mn))
		g.P("}")
		g.P()

		g.P("func (p *", getXStreamServiceMethodClient(sn, mn), ") Start() error {")
		g.P("\tgo func() {")
		g.P("\t\t_= p.receiveLoop()")
		g.P("\t}()")
		g.P("\treturn nil")
		g.P("}")
		g.P()

		g.P("func (p *", getXStreamServiceMethodClient(sn, mn), ") receiveLoop() error {")
		g.P("\t_ = ", lowerFirst(getIStreamServiceMethodClient(sn, mn)), ".", mn, "Pre(p.", getService_MethodClient(sn, mn), ")")
		g.P("\tdefer func() {")
		g.P("\t\tif err := recover(); err != nil {")
		g.P("\t\t\t", xlogPackage.Ident("PrintErr"), "(", xerrorPackage.Ident("GoroutinePanic"), ", p, err, ", debugPackage.Ident("Stack"), "())")
		g.P("\t\t}")
		g.P("\t\t_= ", lowerFirst(getIStreamServiceMethodClient(sn, mn)), ".", mn, "Post(p.", getService_MethodClient(sn, mn), ")")
		g.P("}()")
		g.P("\tfor {")
		g.P("\t\tmsg, err := p.", getService_MethodClient(sn, mn), ".Recv()")
		g.P("\t\tif err != nil {")
		g.P("\t\t\tswitch {")
		g.P("\t\t\tcase err == ", ioPackage.Ident("EOF"), ":")
		g.P("\t\t\t\t", xlogPackage.Ident("PrintErr"), "(\"服务-正常关闭连接\", err)")
		g.P("\t\t\t\treturn nil")
		g.P("\t\t\tcase ", statusPackage.Ident("Code"), "(err) == ", codesPackage.Ident("Canceled"), ":")
		g.P("\t\t\t\t", xlogPackage.Ident("PrintErr"), "(\"服务-取消连接\", err)")
		g.P("\t\t\t\treturn nil")
		g.P("\t\t\tcase ", statusPackage.Ident("Code"), "(err) == ", codesPackage.Ident("Unavailable"), ":")
		g.P("\t\t\t\t", xlogPackage.Ident("PrintErr"), "(\"服务-连接不可用\", err)")
		g.P("\t\t\t\treturn nil")
		g.P("\t\t\tcase ", statusPackage.Ident("Code"), "(err) == ", codesPackage.Ident("Unknown"), ":")
		g.P("\t\t\t\t", xlogPackage.Ident("PrintErr"), "(\"服务-连接异常\", err)")
		g.P("\t\t\t\treturn nil")
		g.P("\t\t\tdefault:")
		g.P("\t\t\t\t", xlogPackage.Ident("PrintErr"), "(\"服务-接收消息错误\", err)")
		g.P("\t\t\t\treturn err")
		g.P("\t\t\t}")
		g.P("\t\t}")
		g.P("\t\t\terr = ", lowerFirst(getIStreamServiceMethodClient(sn, mn)), ".", mn, "(p, msg, p.", getService_MethodClient(sn, mn), ")")
		g.P("\t\t\tif err != nil {")
		g.P("\t\t\t\t", xlogPackage.Ident("PrintErr"), "(\"处理消息错误\", err)")
		g.P("\t\t\t\treturn err")
		g.P("\t\t}")
		g.P("\t}")
		g.P("}")
		g.P()

		g.P("func (p *", clientName, ") ", mn, "(ctx ", contextPackage.Ident("Context"), ", opts ...", grpcPackage.Ident("CallOption"), ") (", grpcPackage.Ident("BidiStreamingClient"), "[", method.Input.GoIdent, ", ", method.Output.GoIdent, "], error) {")
		g.P("\treturn p.Client.", mn, "(ctx, opts...)")
		g.P("}")
		g.P()
	}
}

func clientUnaryGenerateMethod(g *protogen.GeneratedFile, service *protogen.Service, method *protogen.Method, clientName string) {
	sn := service.GoName
	mn := method.GoName
	shardKeyField := findShardKeyField(method.Input)
	g.P("func (p *", clientName, ") ", mn, "(ctx ", contextPackage.Ident("Context"), ", in *", method.Input.GoIdent, ", opts ...", grpcPackage.Ident("CallOption"), ") (*", method.Output.GoIdent, ", error) {")
	if shardKeyField != nil {
		g.P("\tshardKeyValue, err := in.Get_XShardKey()")
		g.P("\tif err != nil {")
		g.P("\t\treturn nil, ", errorsPackage.Ident("WithMessage"), "(err, ", xruntimePackage.Ident("Location"), "())")
		g.P("\t}")
		shardKeyFieldType := goTypeForField(shardKeyField)
		switch shardKeyFieldType {
		case "string":
			g.P("\tstrValue := shardKeyValue")
		case "int32":
			g.P("\tstrValue := ", strconvPackage.Ident("FormatInt"), "(int64(shardKeyValue), 10)")
		case "int64":
			g.P("\tstrValue := ", strconvPackage.Ident("FormatInt"), "(shardKeyValue, 10)")
		case "uint32":
			g.P("\tstrValue := ", strconvPackage.Ident("FormatUint"), "(uint64(shardKeyValue), 10)")
		case "uint64":
			g.P("\tstrValue := ", strconvPackage.Ident("FormatUint"), "(shardKeyValue, 10)")
		default:
			g.P("\t\treturn nil, ", errorsPackage.Ident("WithMessage"), "(", xerrorPackage.Ident("GRPCNotSupportShardKeyType"), ", ", xruntimePackage.Ident("Location"), "())")
		}

		fullMethodName := sn + "_" + mn + "_FullMethodName"
		g.P()

		g.P("\tctx, grpcConn, err := ", xgrpcselectorPackage.Ident("Sel"), "(ctx, ", fullMethodName, ", shardKeyValue)")
		g.P("\tif err != nil {")
		g.P("\t\treturn nil, ", errorsPackage.Ident("WithMessage"), "(err, ", xruntimePackage.Ident("Location"), "())")
		g.P("\t}")
		g.P("\tctx = ", xgrpcprotoPackage.Ident("SetFromOutgoingContext"), "(ctx, ", xgrpcprotoPackage.Ident("ShardKeyFieldNameDefault"), ", strValue)")

		g.P("\tx := New", sn, "Client(grpcConn)")
		g.P("\treturn x.", mn, "(ctx, in, opts...)")
	} else {
		g.P("\treturn nil, ", errorsPackage.Ident("WithMessage"), "(nil, ", xruntimePackage.Ident("Location"), "())")
	}
	g.P("}")
	g.P()
	if shardKeyField != nil {
		clientGenGetShardKeyMethod(g, method, shardKeyField)
	}
}

func clientGenGetShardKeyMethod(g *protogen.GeneratedFile, method *protogen.Method, shardKeyField *protogen.Field) {
	fieldName := shardKeyField.GoName
	fieldType := goTypeForField(shardKeyField)

	g.P("func (x *", method.Input.GoIdent, ") Get_XShardKey()(", fieldType, ", error) {")
	g.P("\treturn x.Get", fieldName, "(), nil")
	g.P("}")
	g.P()
}
