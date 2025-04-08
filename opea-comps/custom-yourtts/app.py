from fastapi import FastAPI, HTTPException
from fastapi.responses import FileResponse
from pydantic import BaseModel
import uvicorn
from TTS.api import TTS
import os
import tempfile
import soundfile as sf
import librosa
import logging

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = FastAPI()

# Initialize TTS model globally
tts = TTS("tts_models/multilingual/multi-dataset/your_tts", progress_bar=True, gpu=False)

class TTSRequest(BaseModel):
    text: str
    reference_audio: str = "/audio/andrew-ref-10s.wav"
    language: str = "en"

@app.get("/health")
def health_check():
    # List contents of audio directory
    audio_dir = "/audio"
    if os.path.exists(audio_dir):
        files = os.listdir(audio_dir)
        logger.info(f"Audio directory contents: {files}")
    else:
        logger.warning(f"Audio directory {audio_dir} not found")
    return {"status": "healthy", "audio_dir_exists": os.path.exists(audio_dir)}

def validate_audio_file(file_path: str) -> bool:
    try:
        # Try loading with soundfile first
        try:
            sf.read(file_path)
            return True
        except Exception as sf_error:
            logger.warning(f"SoundFile failed: {sf_error}")
            # Fallback to librosa
            librosa.load(file_path)
            return True
    except Exception as e:
        logger.error(f"Audio validation failed: {str(e)}")
        return False

@app.post("/tts")
async def text_to_speech(request: TTSRequest):
    output_path = None
    try:
        logger.info(f"Processing TTS request with text: {request.text}")
        logger.info(f"Reference audio path: {request.reference_audio}")
        
        # Validate reference audio exists and is readable
        if not os.path.exists(request.reference_audio):
            logger.error(f"Reference audio not found at: {request.reference_audio}")
            # List contents of audio directory
            audio_dir = os.path.dirname(request.reference_audio)
            if os.path.exists(audio_dir):
                files = os.listdir(audio_dir)
                logger.info(f"Available files in {audio_dir}: {files}")
            raise HTTPException(status_code=400, detail=f"Reference audio file not found: {request.reference_audio}")
        
        if not validate_audio_file(request.reference_audio):
            raise HTTPException(status_code=400, detail=f"Invalid or corrupted reference audio file: {request.reference_audio}")

        # Create temporary file for output
        with tempfile.NamedTemporaryFile(suffix=".wav", delete=False) as temp_file:
            output_path = temp_file.name
            logger.info(f"Created temporary output file: {output_path}")

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

        logger.info(f"Successfully generated audio file: {output_path}")

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
        logger.error(f"Error processing TTS request: {str(e)}")
        # Clean up temp file if it exists
        if output_path and os.path.exists(output_path):
            os.unlink(output_path)
        if isinstance(e, HTTPException):
            raise e
        raise HTTPException(status_code=500, detail=str(e))

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=9881) 