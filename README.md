<img width="1123" height="326" alt="achatbot-go_logo (2)" src="https://github.com/user-attachments/assets/405ad962-6ba7-4367-97b2-64d7e9cbe66e" />

# achatbot-go
a multimodal chatbot.

## Search Functionality
To use the search functionality, you need to set the SERPER_API_KEY environment variable.

Example:
```bash
export SERPER_API_KEY=your_serper_api_key
export SEARCH_API_KEY=your_search_api_key
```

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
- [x] 1. support tool-calls
- [ ] 2. support MCP
- [ ] 3. local VAD/turn + ASR+LLM+TTS remote api Pipeline
- [ ] 4. local VAD/turn + E2E/autonomous llm-audio/omni realtime api Pipeline
- [ ] 5. local Speech-to-Text with Speaker Identification Pipeline
- [ ] 6. webrtc or websocket+webrtc bridge transports
- [ ] 7. local voice agent with micphone
- [ ] 8. 3/4 + streaming avatar api Pipeline
- [ ] 9. AIGC: gen Image/Music/Video remote api Pipeline
- [ ] 10. connecting to RAG services for multimodal features with breaker
- [ ] 11. config and hot reload
- [x] 12. service api add Rate Limiter(IP)
- [x] 13. add pool for modules provider to init load, when connect to get provider to use
- [ ] 14. dockerfile and CD (cloud: AWS ECS, GCP GKE, Azure AKS, Aliyun ECS/ECI) with Terraform


# Acknowledgement
- [ollama](https://github.com/ollama/ollama)
- [sherpa-onnx](https://github.com/k2-fsa/sherpa-onnx)
- [pipeline-go](https://github.com/weedge/pipeline-go) | [pipeline-py](https://github.com/ai-bot-pro/pipeline-py)



# License
achatbot-go is released under the [BSD 3 license](LICENSE). (Additional code in this distribution is covered by the MIT and Apache Open Source
licenses.) However you may have other legal obligations that govern your use of content, such as the terms of service for third-party models.
