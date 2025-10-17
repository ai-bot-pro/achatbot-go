package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/weedge/pipeline-go/pkg/frames"
	"github.com/weedge/pipeline-go/pkg/logger"
	"github.com/weedge/pipeline-go/pkg/pipeline"
	"github.com/weedge/pipeline-go/pkg/processors"
	"github.com/weedge/pipeline-go/pkg/processors/aggregators"
	"github.com/weedge/pipeline-go/pkg/serializers"

	"achatbot/pkg/common"
	"achatbot/pkg/consts"
	"achatbot/pkg/modules/llm"
	"achatbot/pkg/modules/speech/asr"
	"achatbot/pkg/modules/speech/tts"
	"achatbot/pkg/modules/speech/vad_analyzer"
	"achatbot/pkg/params"
	achatbot_processors "achatbot/pkg/processors"
	achatbot_aggregators "achatbot/pkg/processors/aggregators"
	"achatbot/pkg/processors/llm_processors"
	"achatbot/pkg/transports"
	"achatbot/pkg/types"
	achatbot_frames "achatbot/pkg/types/frames"
)

// Upgrader for upgrading HTTP connections to WebSocket connections
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin in this example
		// In production, you should be more restrictive
		return true
	},
}

// Global variables to manage server state
var (
	serverMu    sync.Mutex
	activeTasks = make(map[*pipeline.PipelineTask]bool)
)

// ExampleIWebSocketConn wraps *websocket.Conn to implement our IWebSocketConn interface
type ExampleIWebSocketConn struct {
	*websocket.Conn
	mu sync.Mutex
}

// ReadMessage implements the IWebSocketConn interface
func (wsc *ExampleIWebSocketConn) ReadMessage() (messageType consts.MessageType, p []byte, err error) {
	var msType int
	msType, p, err = wsc.Conn.ReadMessage()
	return consts.MessageType(msType), p, err
}

// WriteMessage implements the IWebSocketConn interface
func (wsc *ExampleIWebSocketConn) WriteMessage(messageType consts.MessageType, data []byte) error {
	if (len(data)) < 200 { //for text frame
		println("WriteMessage-->", messageType.String(), len(data), string(data))
	} else {
		println("WriteMessage-->", messageType.String(), len(data))
	}
	// NOTE: don't concurrent write to websocket connection, need lock
	// issue: concurrent write TextFrame and AudioRawFrame
	wsc.mu.Lock()
	err := wsc.Conn.WriteMessage(int(messageType), data)
	wsc.mu.Unlock()

	return err

}

// Close implements the IWebSocketConn interface
func (wsc *ExampleIWebSocketConn) Close() error {
	println("Close websocket connection")
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

	// Set Session
	clientId := fmt.Sprintf("%s_%s", conn.RemoteAddr().Network(), conn.RemoteAddr().String())
	session := common.NewSession(clientId, nil)

	// vad provider
	sherpaOnnxProvider := vad_analyzer.NewSherpaOnnxProvider(
		//vad_analyzer.NewDefaultSherpaOnnxVadModelConfig("ten"),
		vad_analyzer.NewDefaultSherpaOnnxVadModelConfig("silero"),
		100,
	)
	vadAnalyzer := vad_analyzer.NewVADAnalyzer(params.NewVADAnalyzerArgs(), sherpaOnnxProvider)

	// Wrap the connection to implement our interface
	wsConn := &ExampleIWebSocketConn{Conn: conn}

	// Create audio VAD parameters todo: use viper config to hot load
	audioCameraParams := params.NewAudioCameraParams()
	audioCameraParams.AudioVADParams.WithVADAnalyzer(vadAnalyzer).
		WithVADEnabled(true).WithVADAudioPassthrough(true)
	audioCameraParams.AudioVADParams.AudioParams.
		WithAudioInEnabled(true).WithAudioOutEnabled(true).
		WithAudioInSampleRate(consts.DefaultRate).WithAudioInSampleWidth(consts.DefaultSampleWidth).WithAudioInChannels(consts.DefaultChannels)

	// Create WebSocket server parameters
	wsParams := &params.WebsocketServerParams{
		AudioCameraParams: audioCameraParams,
		Serializer:        serializers.NewProtobufSerializer(),
	}
	wsParams.WithAudioOutFrameMS(200).WithAudioOutAddWavHeader(true) //200ms + wav head

	// Set Websocket Transport Writer
	transportWriter := achatbot_processors.NewWebsocketTransportWriter(wsConn, wsParams)
	audioCameraParams.WithTransportWriter(transportWriter).WithAudioOutEnabled(true).
		WithAudioOutSampleWidth(consts.DefaultSampleWidth).WithAudioOutSampleRate(consts.DefaultRate).WithAudioOutChannels(consts.DefaultChannels)

	// Set ASR Processor
	asrProvider := asr.NewSherpaOnnxProvider(asr.NewDefaultSherpaOnnxOfflineRecognizerConfig())
	asrProcessor := achatbot_processors.NewASRProcessor(asrProvider)

	// Set TTS Processor
	ttsProvider := tts.NewSherpaOnnxProvider(tts.NewDefaultSherpaOnnxOfflineTtsConfig(), tts.KokoroTTS_Speaker_ZM_YunJian, 1.0, "kokoroTTS")
	ttsProcessor := achatbot_processors.NewTTSProcessor(ttsProvider)
	outRate, outChannels, outSampleWidth := ttsProvider.GetSampleInfo()
	audioCameraParams.WithAudioOutSampleWidth(outSampleWidth).WithAudioOutSampleRate(outRate).WithAudioOutChannels(outChannels)

	// Set LLM Processor
	//llmProvider := llm.NewOllamaAPIProviderWithoutTools(llm.OllamaAPIProviderName, llm.OllamaAPIProviderModel_QWEN3_0_6, true, nil, nil)
	//llmProvider := llm.NewOllamaAPIProvider(llm.OllamaAPIProviderName, llm.OllamaAPIProviderModel_QWEN3_0_6, true, nil, nil, []string{"web_search"})
	//llmProcessor := llm_processors.NewLLMOllamaApiProcessor(llmProvider, session, llm_processors.Mode_Chat)
	llmProvider := llm.NewOpenAIAPIProvider(llm.OllamaAPIProviderName, llm.OllamaAPIProviderBaseUrl, llm.OllamaAPIProviderModel_QWEN3_0_6)
	//llmProvider := llm.NewOpenAIAPIProvider(llm.OpenAIAPIProviderName, llm.OpenRouterAIAPIProviderBaseUrl, llm.OpenRouterAIAPIProviderModelQwen2_5_72b_free)
	//llmProvider := llm.NewOpenAIAPIProvider(llm.OpenAIAPIProviderName, llm.OpenRouterAIAPIProviderBaseUrl, llm.OpenRouterAIAPIProviderModelQwen3_235b_free)
	llmProcessor := llm_processors.NewLLMOpenAIApiProcessor(llmProvider, session, llm_processors.Mode_Chat, true, *types.NewLMGenerateArgs())

	// Set Sentence Processor
	sentenceProcessor := aggregators.NewSentenceAggregatorWithEnd(reflect.TypeOf(&achatbot_frames.TurnEndFrame{}))

	// 1. Create the WebSocket server input processor
	ws_transport := transports.NewWebsocketTransport(
		wsConn,
		wsParams,
	)

	// 2. Create a simple pipeline with the async processor
	myPipeline := pipeline.NewPipelineWithVerbose(
		[]processors.IFrameProcessor{
			processors.NewDefaultFrameLoggerProcessorWithIncludeFrame(
				[]frames.Frame{&frames.StartFrame{}, &frames.EndFrame{}, &frames.CancelFrame{}},
			),
			processors.NewDefaultFrameLoggerProcessorWithIncludeFrame([]frames.Frame{&achatbot_frames.BotSpeakingFrame{}}).WithMaxIdToLogs([]uint64{100}),

			ws_transport.InputProcessor(),
			achatbot_aggregators.NewAudioResponseAggregatorWithAccumulate(
				reflect.TypeOf(&achatbot_frames.UserStartedSpeakingFrame{}),
				reflect.TypeOf(&achatbot_frames.UserStoppedSpeakingFrame{}),
				reflect.TypeOf(&achatbot_frames.VADStateAudioRawFrame{}),
			),
			processors.NewDefaultFrameLoggerProcessorWithIncludeFrame([]frames.Frame{&frames.AudioRawFrame{}, &achatbot_frames.VADStateAudioRawFrame{}}),
			achatbot_processors.NewAudioSaveProcessor("user_speak", consts.RECORDS_DIR, true),
			asrProcessor.WithPassRawAudio(false),
			processors.NewDefaultFrameLoggerProcessorWithIncludeFrame([]frames.Frame{&frames.TextFrame{}}),
			llmProcessor,
			//processors.NewDefaultFrameLoggerProcessorWithIncludeFrame([]frames.Frame{&achatbot_frames.ThinkTextFrame{}, &frames.TextFrame{}}),
			sentenceProcessor,
			processors.NewDefaultFrameLoggerProcessorWithIncludeFrame([]frames.Frame{&frames.TextFrame{}}),
			ttsProcessor.WithPassText(true),
			processors.NewDefaultFrameLoggerProcessorWithIncludeFrame([]frames.Frame{&frames.AudioRawFrame{}}),
			//achatbot_processors.NewAudioResampleProcessor(audioCameraParams.AudioOutSampleRate),
			//processors.NewDefaultFrameLoggerProcessorWithIncludeFrame([]frames.Frame{&frames.AudioRawFrame{}}),
			achatbot_processors.NewAudioSaveProcessor("bot_speak", consts.RECORDS_DIR, true),
			processors.NewDefaultFrameLoggerProcessorWithIncludeFrame([]frames.Frame{&frames.AudioRawFrame{}}),
			ws_transport.OutputProcessor(),
		},
		nil, nil,
		false,
	)
	logger.Info(myPipeline.String())

	// In a real application, you would integrate this with your frame processing pipeline
	// and properly manage the processor lifecycle
	// 3. Create and run a pipeline task
	task := pipeline.NewPipelineTask(myPipeline, pipeline.PipelineParams{})

	// Add task to active tasks map
	serverMu.Lock()
	activeTasks[task] = true
	serverMu.Unlock()

	// Remove task from active tasks when done
	defer func() {
		serverMu.Lock()
		delete(activeTasks, task)
		serverMu.Unlock()
	}()

	task.Run()
}

func main() {
	logger.InitLoggerWithConfig(logger.NewDefaultLoggerConfig())
	// Create HTTP server
	server := &http.Server{
		Addr: ":4321",
	}

	// Set up the WebSocket endpoint
	http.HandleFunc("/", handleWebSocket)

	// Channel to listen for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Run server in a goroutine
	go func() {
		logger.Info("Starting WebSocket server on :4321")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-sigChan
	logger.Info("Shutdown signal received")

	// Create a context with timeout for graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Shutdown the HTTP server gracefully
	logger.Info("Shutting down HTTP server...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Cancel all active pipeline tasks
	logger.Info("Cancelling all active pipeline tasks...")
	serverMu.Lock()
	for task := range activeTasks {
		// Send a cancel frame to the pipeline
		task.Cancel()
	}
	serverMu.Unlock()

	// Wait a bit for tasks to finish cleanup
	time.Sleep(1 * time.Second)

	logger.Info("Server exited gracefully")
}
