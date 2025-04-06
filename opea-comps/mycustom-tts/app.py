from fastapi import FastAPI, HTTPException
from fastapi.responses import FileResponse
from pydantic import BaseModel
import uvicorn
from TTS.api import TTS
import os
import tempfile
import soundfile as sf
import librosa

app = FastAPI()

# Initialize TTS model globally
tts = TTS("tts_models/multilingual/multi-dataset/your_tts", progress_bar=True, gpu=False)

class TTSRequest(BaseModel):
    text: str
    reference_audio: str = "/audio/andrew-ref-10s.wav"
    language: str = "en"

@app.get("/health")
def health_check():
    return {"status": "healthy"}

def validate_audio_file(file_path: str) -> bool:
    try:
        # Try loading with soundfile first
        try:
            sf.read(file_path)
            return True
        except:
            # Fallback to librosa
            librosa.load(file_path)
            return True
    except Exception as e:
        return False

@app.post("/tts")
async def text_to_speech(request: TTSRequest):
    output_path = None
    try:
        # Validate reference audio exists and is readable
        if not os.path.exists(request.reference_audio):
            raise HTTPException(status_code=400, detail=f"Reference audio file not found: {request.reference_audio}")
        
        if not validate_audio_file(request.reference_audio):
            raise HTTPException(status_code=400, detail=f"Invalid or corrupted reference audio file: {request.reference_audio}")

        # Create temporary file for output
        with tempfile.NamedTemporaryFile(suffix=".wav", delete=False) as temp_file:
            output_path = temp_file.name

        # Generate audio
        tts.tts_to_file(
            text=request.text,
            file_path=output_path,
            speaker_wav=request.reference_audio,
            language=request.language
        )

        # Validate output was generated
        if not os.path.exists(output_path) or os.path.getsize(output_path) == 0:
            raise HTTPException(status_code=500, detail="Failed to generate audio output")

        # Return the audio file
        response = FileResponse(
            output_path,
            media_type="audio/wav",
            filename="output.wav"
        )

        # Clean up temp file after sending
        response.background = lambda: os.unlink(output_path) if output_path and os.path.exists(output_path) else None
        return response

    except Exception as e:
        # Clean up temp file if it exists
        if output_path and os.path.exists(output_path):
            os.unlink(output_path)
        if isinstance(e, HTTPException):
            raise e
        raise HTTPException(status_code=500, detail=str(e))

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=9881)