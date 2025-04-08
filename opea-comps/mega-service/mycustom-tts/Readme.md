# Text-to-Speech Voice Cloning Solution

## Initial Approach - GPT-SoVITS
- Started with GPT-SoVITS but encountered issues with CPU performance
- The model would consistently get stuck at 1% during text-to-semantic conversion
- Despite various configuration attempts, couldn't get reliable performance on CPU

## Alternative Solutions
Tried several alternative TTS solutions:
- Mozilla TTS (worked but without voice cloning)
- Attempted to use Tortoise TTS and Bark TTS but had Docker image access issues
- Finally settled on Coqui-AI's mycustom-tts

## Final Working Solution - mycustom-tts

### Dependencies
1. Python packages via pip:
   - TTS library
2. System dependencies:
   - espeak-ng for phonemization

### Implementation
- Used the mycustom-TTS model (`tts_models/multilingual/multi-dataset/your_tts`)
- Successfully generated voice cloned audio

## Key Components
- **Reference audio**: `/audio/andrew-ref-1m.wav`
- **Model**: mycustom-TTS (multilingual multi-speaker model)
- **Processing**: CPU-only configuration
- **Output**: Successfully generated WAV files with voice matching

## Testing Commands

### Basic Test
```bash
curl -X POST "http://localhost:9881/tts" \
  -H "Content-Type: application/json" \
  -H "Accept: audio/wav" \
  -d '{"text": "Hello, this is a voice cloning test.", "reference_audio": "/audio/andrew-ref-1m.wav"}' \
  --output output.wav
```

### Extended Test with Longer Text
```bash
curl -X POST "http://localhost:9881/tts" \
  -H "Content-Type: application/json" \
  -H "Accept: audio/wav" \
  -d '{"text": "Hello, this is a voice cloning test. I am trying to match the reference voice as closely as possible. Testing, testing, one two three.", "reference_audio": "/audio/andrew-ref-10s.wav", "language": "en"}' \
  --output output_long.wav
```

### Health Check
```bash
curl http://localhost:9881/health
```

The solution provides a good balance between quality and performance, especially for systems without GPU acceleration.
