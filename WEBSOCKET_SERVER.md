# FastAPI WebSocket Server Input Processor

This document explains how to use the FastAPI WebSocket Server Input Processor implemented in Go.

## Overview

The `FastapiWebsocketServerInputProcessor` is a Go implementation that mirrors the functionality of the Python FastAPI WebSocket server input processor. It receives audio data through WebSocket connections and processes it using the existing audio VAD (Voice Activity Detection) pipeline.

## Key Components

### 1. FastapiWebsocketServerInputProcessor

This processor extends the existing `AudioVADInputProcessor` and adds WebSocket functionality:

- Receives binary messages from WebSocket connections
- Deserializes incoming data into frames
- Processes `AudioRawFrame` frames through the VAD pipeline

### 2. Dependencies

The implementation uses:
- `github.com/weedge/pipeline-go` for frame processing
- Generic WebSocket connection interface (can work with Gorilla WebSocket or standard library)
- Custom serializer interface for frame serialization/deserialization

## Usage

### Setting up the WebSocket Server

1. Create audio VAD parameters:
```go
audioVADParams := params.NewAudioVADParams()
audioVADParams.WithAudioInEnabled(true)
audioVADParams.WithVADEnabled(true)
```

2. Create WebSocket server parameters:
```go
wsParams := &processors.FastapiWebsocketServerParams{
    AudioVADParams: audioVADParams,
    Serializer:     utils.NewJSONSerializer(), // or your custom serializer
}
```

3. Define callbacks:
```go
callbacks := &processors.FastapiWebsocketServerCallbacks{
    OnClientConnected: func(ws processors.WebSocketConn) {
        log.Println("Client connected")
    },
    OnClientDisconnected: func(ws processors.WebSocketConn) {
        log.Println("Client disconnected")
    },
}
```

4. Create the processor:
```go
processor := processors.NewFastapiWebsocketServerInputProcessor(
    "websocket_processor",
    wsConn,  // Your WebSocket connection implementing the interface
    wsParams,
    callbacks,
)
```

5. Start the processor:
```go
startFrame := &frames.StartFrame{} // Create appropriate start frame
processor.Start(startFrame)
```

## WebSocket Connection Interface

To use this processor with different WebSocket libraries, implement the `WebSocketConn` interface:

```go
type WebSocketConn interface {
    ReadMessage() (messageType int, p []byte, err error)
    WriteMessage(messageType int, data []byte) error
    Close() error
}
```

## Serialization

The processor uses a generic `Serializer` interface:

```go
type Serializer interface {
    Serialize(frame frames.Frame) ([]byte, error)
    Deserialize(data []byte) frames.Frame
}
```

A JSON implementation is provided in `pkg/utils/json_serializer.go`.

## Integration with Existing Pipeline

The `FastapiWebsocketServerInputProcessor` integrates with the existing frame processing pipeline by:

1. Receiving audio data through WebSocket
2. Converting it to `AudioRawFrame` objects
3. Pushing frames through the VAD processing pipeline using `PushAudioFrame()`

## Dependencies to Install

If using Gorilla WebSocket:
```
go get github.com/gorilla/websocket
```

## Example Implementation

See `examples/websocket_server_example.go` for a complete example of how to set up a WebSocket server using this processor.