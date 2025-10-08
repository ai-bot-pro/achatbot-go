package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/weedge/pipeline-go/pkg/pipeline"
	"github.com/weedge/pipeline-go/pkg/processors"
	"github.com/weedge/pipeline-go/pkg/serializers"

	"achatbot/pkg/common"
	"achatbot/pkg/params"
	achabot_processors "achatbot/pkg/processors"
	// Assuming we're using the gorilla websocket package
	// You would need to add this dependency: go get github.com/gorilla/websocket
)

// Upgrader for upgrading HTTP connections to WebSocket connections
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin in this example
		// In production, you should be more restrictive
		return true
	},
}

// ExampleWebSocketConn wraps *websocket.Conn to implement our WebSocketConn interface
type ExampleWebSocketConn struct {
	*websocket.Conn
}

// ReadMessage implements the WebSocketConn interface
func (wsc *ExampleWebSocketConn) ReadMessage() (messageType int, p []byte, err error) {
	return wsc.Conn.ReadMessage()
}

// WriteMessage implements the WebSocketConn interface
func (wsc *ExampleWebSocketConn) WriteMessage(messageType int, data []byte) error {
	return wsc.Conn.WriteMessage(messageType, data)
}

// Close implements the WebSocketConn interface
func (wsc *ExampleWebSocketConn) Close() error {
	return wsc.Conn.Close()
}

// handleWebSocket handles incoming WebSocket connections
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	// Wrap the connection to implement our interface
	wsConn := &ExampleWebSocketConn{Conn: conn}

	// Create audio VAD parameters
	audioVADParams := params.NewAudioVADParams()
	audioVADParams.WithAudioInEnabled(true)
	audioVADParams.WithVADEnabled(true)

	// Create WebSocket server parameters
	wsParams := &achabot_processors.WebsocketServerParams{
		AudioVADParams: audioVADParams,
		Serializer:     serializers.NewJsonSerializer(),
	}

	// Create callbacks
	callbacks := &achabot_processors.WebsocketServerCallbacks{
		OnClientConnected: func(ws common.WebSocketConn) {
			log.Println("Client connected")
		},
		OnClientDisconnected: func(ws common.WebSocketConn) {
			log.Println("Client disconnected")
		},
	}

	// 1. Create the WebSocket server input processor
	ws_input_processor := achabot_processors.NewWebsocketServerInputProcessor(
		"websocket_processor",
		wsConn,
		wsParams,
		callbacks,
	)
	ws_input_processor.SetVerbose(true)

	// 2. Link it to a logger processor
	logger := processors.NewFrameTraceLogger("ws_input", 0)
	logger.SetVerbose(true)

	// 3. Create a simple pipeline with the async processor
	myPipeline := pipeline.NewPipeline(
		[]processors.IFrameProcessor{
			ws_input_processor,
			logger,
		},
		nil, nil,
	)
	slog.Info(myPipeline.String())

	// In a real application, you would integrate this with your frame processing pipeline
	// and properly manage the processor lifecycle
	// 4. Create and run a pipeline task
	task := pipeline.NewPipelineTask(myPipeline, pipeline.PipelineParams{})
	task.Run()
}

func main() {
	// Set up the WebSocket endpoint
	http.HandleFunc("/ws", handleWebSocket)

	log.Println("Starting WebSocket server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
