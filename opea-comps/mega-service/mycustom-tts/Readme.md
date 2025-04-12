# Text-to-Speech Voice Cloning Solution

## Technical Uncertainties and Challenges

### VLLM Implementation Challenges
- Initially attempted to use VLLM for text processing
- Encountered CPU limitations:
  - Required 8GB memory allocation
  - Performance issues on CPU-only environments
  - Async worker configuration challenges
- VLLM service was commented out due to these limitations

### Solution Path
1. Evaluated alternative approaches:
   - Considered CPU-optimized models
   - Tested different memory configurations
2. Final Solution:
   - Implemented custom TTS service with better CPU compatibility
   - Reduced memory requirements (4GB vs 8GB)
   - Achieved stable performance without GPU dependency

## Available TTS Services

### 1. GPT-SoVITS Service
- Endpoint: `http://localhost:9880`
- Container: `gpt-sovits-service`
- Features:
  - Voice cloning capabilities
  - REST API interface
  - Health check endpoint

### 2. Custom TTS Service (mycustom-tts)
- Endpoint: `http://localhost:9881`
- Container: `mytts-service`
- Features:
  - CPU-optimized implementation
  - Multilingual support
  - Voice cloning capabilities
  - Memory limit: 4G

### 3. SpeechT5 Service
- Endpoint: `http://localhost:7055`
- Container: `speecht5-service`
- Features:
  - Text-to-speech conversion
  - Memory limit: 4G
  - Health check monitoring

### 4. TTS GPT-SoVITS Wrapper Service
- Endpoint: `http://localhost:9088`
- Container: `tts-gptsovits-service`
- Features:
  - Wrapper around GPT-SoVITS service
  - Configurable via environment variables
  - Depends on gpt-sovits-service

## Implementation Details

### Dependencies
1. Docker and Docker Compose
2. System requirements:
   - Minimum 4GB RAM per service
   - CPU with AVX2 support recommended
   - Sufficient storage for model weights

### Volume Mounts
- Audio files: `./audio:/audio`
- Shared across services for reference audio files

## API Usage Examples

### GPT-SoVITS Health Check
```bash
curl http://localhost:9880/health
```

### Custom TTS Request
```bash
curl -X POST "http://localhost:9881/tts" \
  -H "Content-Type: application/json" \
  -H "Accept: audio/wav" \
  -d '{
    "text": "Hello, this is a voice cloning test.",
    "reference_audio": "/audio/andrew-ref-1m.wav"
  }' \
  --output output.wav
```

### SpeechT5 Health Check
```bash
curl http://localhost:7055/health
```

## Environment Configuration
Services can be configured using environment variables:
- `SPEECHT5_PORT`: Default 7055
- `GPT_SOVITS_PORT`: Default 9880
- `YOURTTS_PORT`: Default 9881
- `TTS_PORT`: Default 9088
- `TTS_COMPONENT_NAME`: Configurable component name for GPT-SoVITS wrapper
- `REGISTRY`: Docker registry prefix
- `TAG`: Image version tag

## Health Monitoring
All services include health check endpoints with:
- 10-second check intervals
- 6-second timeouts
- 18 retry attempts

The solution provides multiple TTS options with different capabilities, allowing for fallback options and specific use-case optimization.
