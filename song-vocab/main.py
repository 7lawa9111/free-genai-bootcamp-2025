from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import List, Optional
import os
from agent import LyricsAgent
from database import Database

app = FastAPI(title="Song Vocabulary API")
db = Database()
agent = LyricsAgent(output_dir="outputs")

class MessageRequest(BaseModel):
    message_request: str

class VocabularyItem(BaseModel):
    word: str
    definition: str
    example: str

class VocabularyResponse(BaseModel):
    result_id: str
    lyrics: str
    vocabulary: List[VocabularyItem]

class SongResponse(BaseModel):
    id: int
    title: str
    lyrics: str
    vocabulary: List[VocabularyItem]

class SongListItem(BaseModel):
    id: int
    title: str
    vocabulary_count: int

def _load_result(result_id: str) -> VocabularyResponse:
    """Load lyrics and vocabulary from files"""
    try:
        # Load lyrics
        lyrics_path = os.path.join("outputs", f"{result_id}_lyrics.txt")
        with open(lyrics_path, 'r') as f:
            # Skip the query line and blank line
            next(f)
            next(f)
            lyrics = f.read().strip()
        
        # Load vocabulary
        vocab_path = os.path.join("outputs", f"{result_id}_vocab.txt")
        vocabulary = []
        with open(vocab_path, 'r') as f:
            current_item = {}
            for line in f:
                line = line.strip()
                if not line:
                    if current_item:
                        vocabulary.append(VocabularyItem(**current_item))
                        current_item = {}
                elif line.startswith('Word:'):
                    current_item['word'] = line[5:].strip()
                elif line.startswith('Definition:'):
                    current_item['definition'] = line[11:].strip()
                elif line.startswith('Example:'):
                    current_item['example'] = line[8:].strip()
            if current_item:
                vocabulary.append(VocabularyItem(**current_item))
        
        return VocabularyResponse(
            result_id=result_id,
            lyrics=lyrics,
            vocabulary=vocabulary
        )
    except Exception as e:
        raise HTTPException(status_code=404, detail=f"Result not found: {str(e)}")

@app.post("/api/agent", response_model=VocabularyResponse)
def get_lyrics(request: MessageRequest):
    try:
        # Get result ID from agent
        result_id = agent.run(request.message_request)
        if not result_id:
            raise HTTPException(status_code=500, detail="Failed to process request")
        
        # Load the results
        result = _load_result(result_id)
        
        # Save to database
        db.save_song_and_vocabulary(
            title=request.message_request,
            lyrics=result.lyrics,
            vocabulary=[v.dict() for v in result.vocabulary]
        )
        
        return result
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/api/songs", response_model=List[SongListItem])
def list_songs():
    """Get all songs in the database"""
    return db.get_all_songs()

@app.get("/api/songs/{song_id}", response_model=SongResponse)
def get_song(song_id: int):
    """Get a specific song and its vocabulary"""
    result = db.get_song_vocabulary(song_id)
    if not result:
        raise HTTPException(status_code=404, detail="Song not found")
    return {
        "id": song_id,
        **result
    }

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000) 