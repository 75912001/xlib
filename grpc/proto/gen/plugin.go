package main

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	contextPackage               = protogen.GoImportPath("context")
	grpcPackage                  = protogen.GoImportPath("google.golang.org/grpc")
	errorsPackage                = protogen.GoImportPath("github.com/pkg/errors")
	xruntimePackage              = protogen.GoImportPath("github.com/75912001/xlib/runtime")
	xgrpcprotoPackage            = protogen.GoImportPath("github.com/75912001/xlib/grpc/proto")
	xgrpcselectorPackage         = protogen.GoImportPath("github.com/75912001/xlib/grpc/selector")
	xerrorPackage                = protogen.GoImportPath("github.com/75912001/xlib/error")
	strconvPackage               = protogen.GoImportPath("strconv")
	xlogPackage                  = protogen.GoImportPath("github.com/75912001/xlib/log")
	debugPackage                 = protogen.GoImportPath("runtime/debug")
	ioPackage                    = protogen.GoImportPath("io")
	statusPackage                = protogen.GoImportPath("google.golang.org/grpc/status")
	codesPackage                 = protogen.GoImportPath("google.golang.org/grpc/codes")
	insecurePackage              = protogen.GoImportPath("google.golang.org/grpc/credentials/insecure")
	xgrpcutilPackage             = protogen.GoImportPath("github.com/75912001/xlib/grpc/util")
	xgrpcprotointerceptorPackage = protogen.GoImportPath("github.com/75912001/xlib/grpc/proto/interceptor")
	xcontrolPackage              = protogen.GoImportPath("github.com/75912001/xlib/control")
)

func main() {
	protogen.Options{}.Run(func(gen *protogen.Plugin) error {
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			generateFile(gen, f)
		}
		return nil
	})
}

func generateFile(gen *protogen.Plugin, file *protogen.File) {
	if len(file.Services) == 0 {
		return
	}

	filename := file.GeneratedFilenamePrefix + "_grpc.x.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	genHeader(g, file)
	genFunction(g, file)
	genClient(g, file)
	genServer(g, file)
}

func goTypeForField(field *protogen.Field) string {
	switch field.Desc.Kind() {
	case protoreflect.Int32Kind:
		return "int32"
	case protoreflect.Uint32Kind:
		return "uint32"
	case protoreflect.Int64Kind:
		return "int64"
	case protoreflect.Uint64Kind:
		return "uint64"
	case protoreflect.StringKind:
		return "string"
	default:
		return "any"
	}
}
