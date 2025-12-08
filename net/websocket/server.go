package websocket

import (
	"context"
	"errors"
	xconfig "github.com/75912001/xlib/config"
	xcontrol "github.com/75912001/xlib/control"
	xerror "github.com/75912001/xlib/error"
	xlog "github.com/75912001/xlib/log"
	xnetcommon "github.com/75912001/xlib/net/common"
	xpacket "github.com/75912001/xlib/packet"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/gorilla/websocket"
	pkgerrors "github.com/pkg/errors"
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

// Server 服务端
type Server struct {
	IHandler   xnetcommon.IHandler
	options    *ServerOptions
	upgrader   *websocket.Upgrader
	httpServer *http.Server // 添加HTTP服务器引用
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewServer 新建服务
func NewServer(handler xnetcommon.IHandler) *Server {
	return &Server{
		IHandler: handler,
		options:  nil,
	}
}

// Start 运行服务
func (p *Server) Start(ctx context.Context, opts ...*ServerOptions) error {
	p.options = mergeServerOptions(opts...)
	if err := configureServerOptions(p.options); err != nil {
		return pkgerrors.WithMessagef(err, "configureServerOptions:%v %v", p.options, xruntime.Location())
	}

	// 定义升级器
	p.upgrader = &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // 允许所有来源，生产环境应该更严格
		},
		ReadBufferSize:  *p.options.ReadBuffer,
		WriteBufferSize: *p.options.WriteBuffer,
	}

	mux := http.NewServeMux()
	mux.HandleFunc(*p.options.pattern, p.handleWebSocket)
	// 创建HTTP服务器
	p.httpServer = &http.Server{
		Addr:    *p.options.listenAddress,
		Handler: mux,
	}
	p.ctx, p.cancel = context.WithCancel(ctx)
	go func() {
		defer func() {
			if xruntime.IsRelease() {
				if err := recover(); err != nil {
					xlog.PrintErr(xerror.GoroutinePanic, err, debug.Stack())
				}
			}
			xlog.PrintInfo(xerror.GoroutineDone)
		}()
		err := p.httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			xlog.PrintErr(xerror.GoroutinePanic, err, debug.Stack())
		}
		xlog.PrintErr(xerror.GoroutinePanic, err, debug.Stack())
	}()
	return nil
}

// Stop 停止 WebSocket 服务器
func (p *Server) Stop() {
	if p.httpServer != nil {
		// 优雅关闭HTTP服务器
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		err := p.httpServer.Shutdown(shutdownCtx)
		if err != nil {
			xlog.PrintfErr("server shutdown err:%v", err)
		}
		p.httpServer = nil
	}

	if p.cancel != nil {
		p.cancel()
	}
}

// 处理 WebSocket 连接
func (p *Server) handleWebSocket(w http.ResponseWriter, req *http.Request) {
	// 升级 HTTP 连接为 WebSocket 连接
	conn, err := p.upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Printf("升级失败: %v", err)
		return
	}

	remote := p.handleConn(conn, p.options.iOut)
	defer func() {
		if xruntime.IsRelease() {
			if r := recover(); r != nil {
				xlog.PrintErr(xerror.GoroutinePanic, p, r, debug.Stack())
			}
		}
		if xconfig.GConfigMgr.Base.ProcessingModeIsActor() {
			_ = p.IHandler.OnDisconnect(remote)
		} else {
			p.options.iOut.Send(
				&xnetcommon.Disconnect{
					IHandler: p.IHandler,
					IRemote:  remote,
				},
			)
		}
		_ = conn.Close()
	}()

	var messageType int
	var buf []byte
	var packet xpacket.IPacket
	// 处理连接
	for {
		// 读取消息
		messageType, buf, err = conn.ReadMessage()
		if err != nil {
			xlog.PrintfErr("read message fail. err:%v %v", err, debug.Stack())
			if remote.GetDisconnectReason() == xnetcommon.DisconnectReasonUnknown { // 未设置,就设置为客户端主动断开
				remote.SetDisconnectReason(xnetcommon.DisconnectReasonClientShutdown)
			}
			break
		}
		if messageType != websocket.BinaryMessage {
			xlog.PrintfErr("read message fail. err:%v %v", xerror.NotSupport, messageType)
			if remote.GetDisconnectReason() == xnetcommon.DisconnectReasonUnknown { // 未设置,就设置为客户端主动断开
				remote.SetDisconnectReason(xnetcommon.DisconnectReasonClientShutdown)
			}
			break
		}
		if err = p.IHandler.OnCheckPacketLimit(remote); err != nil { // 限流
			xlog.PrintfErr("remote:%p buf:%v err:%v", p, buf, err)
			continue
		}
		packet, err = p.IHandler.OnUnmarshalPacket(remote, buf)
		if err != nil {
			xlog.PrintfErr("remote:%p buf:%v err:%v", p, buf, err)
			continue
		}
		if xconfig.GConfigMgr.Base.ProcessingModeIsActor() {
			_ = p.IHandler.OnPacket(remote, packet)
		} else {
			p.options.iOut.Send(
				&xnetcommon.Packet{
					IHandler: p.IHandler,
					IRemote:  remote,
					IPacket:  packet,
				},
			)
		}
	}
}

func (p *Server) handleConn(conn *websocket.Conn, iOut xcontrol.IOut) *Remote {
	remote := NewRemote(conn, make(chan any, *p.options.sendChanCapacity))
	if p.options.NewPacketLimitFunc != nil {
		remote.PacketLimit = p.options.NewPacketLimitFunc(p.options.MaxCntPerSec)
	}
	if xconfig.GConfigMgr.Base.ProcessingModeIsActor() {
		_ = p.IHandler.OnConnect(remote)
	} else {
		iOut.Send(
			&xnetcommon.Connect{
				IHandler: p.IHandler,
				IRemote:  remote,
			},
		)
	}
	remote.Start(&p.options.ConnOptions, iOut, p.IHandler)
	return remote
}
