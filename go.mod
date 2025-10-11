module achatbot

go 1.24.0

require (
	github.com/gorilla/websocket v1.5.3
	github.com/k2-fsa/sherpa-onnx-go v1.12.12
	github.com/stretchr/testify v1.10.0
	github.com/weedge/pipeline-go v0.0.0-20251010102304-44bf1620471c
	golang.org/x/image v0.32.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/k2-fsa/sherpa-onnx-go-linux v1.12.13 // indirect
	github.com/k2-fsa/sherpa-onnx-go-macos v1.12.13 // indirect
	github.com/k2-fsa/sherpa-onnx-go-windows v1.12.13 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/weedge/pipeline-go v0.0.0-20251010102304-44bf1620471c => ../pipeline-go
