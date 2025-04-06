FROM python:3.11-slim

# Install system dependencies
RUN apt-get update && apt-get install -y \
    espeak-ng \
    curl \
    libsndfile1 \
    ffmpeg \
    libportaudio2 \
    libasound2-dev \
    && rm -rf /var/lib/apt/lists/*

# Install Python packages
RUN pip install --no-cache-dir \
    TTS \
    fastapi \
    uvicorn \
    pydantic \
    soundfile \
    audioread \
    librosa

# Copy service files
WORKDIR /app
COPY app.py .

# Set environment variables
ENV PYTHONUNBUFFERED=1

# Expose port
EXPOSE 9881

# Run the service
CMD ["python", "app.py"]