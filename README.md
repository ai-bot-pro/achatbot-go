<img width="1123" height="326" alt="Image" src="https://github.com/user-attachments/assets/e1b82973-1bf8-4490-a319-e74faf0e5f06"/>

# achatbot-go
a multimodal chatbot.

## local VAD+ASR+LLM+TTS Pipeline
- run local vad+asr+llm+tts pipeline websocket voice agent (not agentic), need download [ollama](https://docs.ollama.com/quickstart) and start ollama server

```shell
# 0. install deps
go mod tidy

# 1. download models (ONNX)
## silero VAD
curl -SL -O https://github.com/k2-fsa/sherpa-onnx/releases/download/asr-models/silero_vad.onnx
## ten VAD
curl -SL -O https://github.com/k2-fsa/sherpa-onnx/releases/download/asr-models/ten-vad.onnx

## sensevoice ASR
huggingface-cli download csukuangfj/sherpa-onnx-sense-voice-zh-en-ja-ko-yue-2024-07-17 --local-dir ./models/csukuangfj/sherpa-onnx-sense-voice-zh-en-ja-ko-yue-2024-07-17
## kokoro TTS
huggingface-cli download csukuangfj/kokoro-multi-lang-v1_0 --local-dir ./models/csukuangfj/kokoro-multi-lang-v1_0

# 2. run websocket server
go run examples/websocket/server.go

# 3. run ui client
cd examples/websocket/ui/ && python -m http.server
# - access http://localhost:8000 to Start Audio
```

## TODO
- [ ] 1. support tool-calls
- [ ] 2. support MCP
- [ ] 3. local VAD + ASR+LLM+TTS remote api Pipeline
- [ ] 4. local VAD + E2E/autonomous llm-audio/omni realtime api Pipeline
- [ ] 5. local Speech-to-Text with Speaker Identification
- [ ] 6. webrtc or websocket+webrtc bridge transports
- [ ] 7. local voice agent with micphone
- [ ] 8. 3/4 + streaming avatar api pipeline
- [ ] 9. AIGC: gen Image/Video remote api




# Acknowledgement
- [ollama](https://github.com/ollama/ollama)
- [sherpa-onnx](https://github.com/k2-fsa/sherpa-onnx)
- [pipeline-go](https://github.com/weedge/pipeline-go) | [pipeline-py](https://github.com/ai-bot-pro/pipeline-py)



# License
achatbot-go is released under the [BSD 3 license](LICENSE). (Additional code in this distribution is covered by the MIT and Apache Open Source
licenses.) However you may have other legal obligations that govern your use of content, such as the terms of service for third-party models.
