package main

import (
	"google.golang.org/protobuf/compiler/protogen"
	"strings"
)

func genServer(g *protogen.GeneratedFile, file *protogen.File) {
	for _, service := range file.Services {
		sn := service.GoName
		unimplementedServer := "Unimplemented" + sn + "Server"

		g.P()
		g.P("////////////////////////////////////////////////////////////////////////////////////////////////////")
		g.P("// ", sn, " 服务端")
		g.P("////////////////////////////////////////////////////////////////////////////////////////////////////")
		g.P("type ", getXServiceServer(sn), " struct {")
		g.P("\t", unimplementedServer)
		g.P("}")
		g.P()

		g.P("func ", getNewXServiceServer(sn), "()*", getXServiceServer(sn), " {")
		g.P("\t", strings.ToLower(sn), "Server := &", getXServiceServer(sn), "{}")
		g.P("\treturn ", strings.ToLower(sn), "Server")
		g.P("}")
		g.P()

		g.P("func (p *", getXServiceServer(sn), ") Start(grpcServer *", grpcPackage.Ident("Server"), ") error {")
		g.P("\tRegister", sn, "Server(grpcServer, p)")
		g.P("\treturn nil")
		g.P("}")
		g.P()
		g.P("func (p *", getXServiceServer(sn), ") Stop() error {")
		g.P("\treturn nil")
		g.P("}")
		g.P()

		g.P("////////////////////////////////////////////////////////////////////////////////")
		g.P("//    ", sn, " 服务端-Stream")
		g.P("////////////////////////////////////////////////////////////////////////////////")
		g.P("type ", getIStreamServiceServer(sn), " interface {")
		for _, m := range service.Methods {
			if m.Desc.IsStreamingClient() || m.Desc.IsStreamingServer() {
				g.P("\t", m.GoName, "Pre(stream ", sn, "_", m.GoName, "Server) error // 预处理-新创建stream")
				g.P("\t", m.GoName, "(", strings.ToLower(sn), "Server *", getXServiceServer(sn), ",msgWrapper *", m.Input.GoIdent, ",stream ", sn, "_", m.GoName, "Server) error // 处理")
				g.P("\t", m.GoName, "Post(stream ", sn, "_", m.GoName, "Server) error // 后处理-关闭stream")
			}
		}
		g.P("}")
		g.P()
		g.P("var ", lowerFirst(getIStreamServiceServer(sn)), " ", getIStreamServiceServer(sn))
		g.P()
		g.P("func Set", getIStreamServiceServer(sn), "(streamServer ", getIStreamServiceServer(sn), ") {")
		g.P("\t", lowerFirst(getIStreamServiceServer(sn)), " = streamServer")
		g.P("}")
		g.P()

		for _, m := range service.Methods {
			if m.Desc.IsStreamingClient() || m.Desc.IsStreamingServer() {
				g.P("func (p *", getXServiceServer(sn), ") ", m.GoName, "(stream ", getService_MethodServer(sn, m.GoName), ") error {")
				g.P("\t_ = ", lowerFirst(getIStreamServiceServer(sn)), ".", m.GoName, "Pre(stream)")
				g.P("\tdefer func() {")
				g.P("\t\tif err := recover(); err != nil {")
				g.P("\t\t\t", xlogPackage.Ident("PrintErr"), "(", xerrorPackage.Ident("GoroutinePanic"), ", p, err, ", debugPackage.Ident("Stack"), "())")
				g.P("\t\t}")
				g.P("\t\t_ = ", lowerFirst(getIStreamServiceServer(sn)), ".", m.GoName, "Post(stream)")
				g.P("\t}()")
				g.P("\tfor {")
				g.P("\t\tmsg, err := stream.Recv()")
				g.P("\t\tif err != nil {")
				g.P("\t\t\tswitch {")
				g.P("\t\t\tcase err == ", ioPackage.Ident("EOF"), ":")
				g.P("\t\t\t\t", xlogPackage.Ident("PrintErr"), "(\"客户端-正常关闭连接\", err)")
				g.P("\t\t\t\treturn nil")
				g.P("\t\t\tcase ", statusPackage.Ident("Code"), "(err) == ", codesPackage.Ident("Canceled"), ":")
				g.P("\t\t\t\t", xlogPackage.Ident("PrintErr"), "(\"客户端-取消连接\", err)")
				g.P("\t\t\t\treturn nil")
				g.P("\t\t\tcase ", statusPackage.Ident("Code"), "(err) == ", codesPackage.Ident("Unavailable"), ":")
				g.P("\t\t\t\t", xlogPackage.Ident("PrintErr"), "(\"客户端-连接不可用\", err)")
				g.P("\t\t\t\treturn nil")
				g.P("\t\t\tcase ", statusPackage.Ident("Code"), "(err) == ", codesPackage.Ident("Unknown"), ":")
				g.P("\t\t\t\t", xlogPackage.Ident("PrintErr"), "(\"客户端-连接异常\", err)")
				g.P("\t\t\t\treturn nil")
				g.P("\t\t\tdefault:")
				g.P("\t\t\t\t", xlogPackage.Ident("PrintErr"), "(\"接收消息错误\", err)")
				g.P("\t\t\t\treturn err")
				g.P("\t\t\t}")
				g.P("\t\t}")
				g.P("\t\terr = ", lowerFirst(getIStreamServiceServer(sn)), ".", m.GoName, "(p, msg, stream)")
				g.P("\t\tif err != nil {")
				g.P("\t\t\t", xlogPackage.Ident("PrintErr"), "(\"处理消息错误\", err)")
				g.P("\t\t\treturn err")
				g.P("\t\t}")
				g.P("\t}")
				g.P("}")
				g.P()
			}
		}

		g.P("////////////////////////////////////////////////////////////////////////////////")
		g.P("//    ", sn, " 服务端-Unary")
		g.P("////////////////////////////////////////////////////////////////////////////////")
		g.P("type ", getIUnaryServiceServer(sn), " interface {")
		for _, m := range service.Methods {
			if !m.Desc.IsStreamingClient() && !m.Desc.IsStreamingServer() {
				g.P("\t", m.GoName, "(ctx ", contextPackage.Ident("Context"), ", req *", m.Input.GoIdent, ") (*", m.Output.GoIdent, ", error)")
			}
		}
		g.P("}")
		g.P()
		g.P("var ", lowerFirst(getIUnaryServiceServer(sn)), " ", getIUnaryServiceServer(sn))
		g.P()
		g.P("func Set", getIUnaryServiceServer(sn), "(unaryServer ", getIUnaryServiceServer(sn), ") {")
		g.P("\t", lowerFirst(getIUnaryServiceServer(sn)), " = unaryServer")
		g.P("}")
		g.P()

		for _, m := range service.Methods {
			if !m.Desc.IsStreamingClient() && !m.Desc.IsStreamingServer() {
				g.P("func (p *", getXServiceServer(sn), ") ", m.GoName, "(ctx ", contextPackage.Ident("Context"), ", req *", m.Input.GoIdent, ") (*", m.Output.GoIdent, ", error) {")
				g.P("\treturn ", lowerFirst(getIUnaryServiceServer(sn)), ".", m.GoName, "(ctx, req)")
				g.P("}")
				g.P()
			}
		}
	}
}
