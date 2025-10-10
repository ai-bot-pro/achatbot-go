package transports

import (
	"achatbot/pkg/common"
	achatbot_params "achatbot/pkg/params"
	achabot_processors "achatbot/pkg/processors"
)

// WebsocketTransport 实现了 WebSocket 传输层，组合 EventHandlerManager 成员和方法
type WebsocketTransport struct {
	*common.EventHandlerManager
	websocket       common.IWebSocketConn
	params          *achatbot_params.WebsocketServerParams
	callbacks       *achabot_processors.WebsocketServerCallbacks
	inputProcessor  *achabot_processors.WebsocketServerInputProcessor
	outputProcessor *achabot_processors.WebsocketServerOutputProcessor
}

// NewWebsocketTransport 创建一个新的 WebsocketTransport 实例
func NewWebsocketTransport(
	websocket common.IWebSocketConn,
	params *achatbot_params.WebsocketServerParams,
) *WebsocketTransport {
	// 创建事件管理器
	eventManager := common.NewEventHandlerManagerWithName("websocket_transport")

	// 创建 transport 实例
	transport := &WebsocketTransport{
		EventHandlerManager: eventManager,
		websocket:           websocket,
		params:              params,
	}

	// 初始化回调函数
	transport.callbacks = &achabot_processors.WebsocketServerCallbacks{
		OnClientConnected:    transport.onClientConnected,
		OnClientDisconnected: transport.onClientDisconnected,
	}

	transport.inputProcessor = achabot_processors.NewWebsocketServerInputProcessor(
		"WebsocketServerInputProcessor",
		websocket,
		params,
		transport.callbacks,
	)

	// 创建输出处理器
	transport.outputProcessor = achabot_processors.NewWebsocketServerOutputProcessor(
		"WebsocketServerOutputProcessor",
		params,
	)

	// 注册支持的事件处理器
	transport.RegisterEventHandler("on_client_connected")
	transport.RegisterEventHandler("on_client_disconnected")

	return transport
}

// InputProcessor 返回输入处理器
func (wt *WebsocketTransport) InputProcessor() *achabot_processors.WebsocketServerInputProcessor {
	return wt.inputProcessor
}

// OutputProcessor 返回输出处理器
func (wt *WebsocketTransport) OutputProcessor() *achabot_processors.WebsocketServerOutputProcessor {
	return wt.outputProcessor
}

// onClientConnected 处理客户端连接事件
func (wt *WebsocketTransport) onClientConnected(websocket common.IWebSocketConn) {
	wt.CallEventHandler("on_client_connected", websocket)
}

// onClientDisconnected 处理客户端断开连接事件
func (wt *WebsocketTransport) onClientDisconnected(websocket common.IWebSocketConn) {
	wt.CallEventHandler("on_client_disconnected", websocket)
}
