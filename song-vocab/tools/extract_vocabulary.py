from typing import List, Union
from pathlib import Path
from pydantic import BaseModel, Field
import instructor
from instructor import OpenAISchema
import ollama
import os

class RomajiPart(OpenAISchema):
    """Schema for a part of a Japanese word with its romaji readings"""
    kanji: str = Field(..., description="The kanji or kana character")
    romaji: List[str] = Field(..., description="List of possible romaji readings for this part")

class VocabularyItem(OpenAISchema):
    """Schema for a Japanese vocabulary item extracted from lyrics"""
    kanji: str = Field(..., description="The full word in kanji/kana")
    romaji: str = Field(..., description="The romaji reading of the full word")
    english: str = Field(..., description="The English translation/meaning")
    parts: List[RomajiPart] = Field(..., description="Breakdown of each part of the word with its readings")

def load_prompt(prompt_name: str) -> str:
    """Load a prompt from the prompts directory"""
    current_dir = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
    prompt_path = os.path.join(current_dir, "prompts", prompt_name)
    with open(prompt_path, "r", encoding="utf-8") as f:
        return f.read()

def extract_vocabulary(lyrics: str) -> List[VocabularyItem]:
    """
    Extract Japanese vocabulary items from lyrics using Ollama with Instructor
    
    Args:
        lyrics (str): Song lyrics to analyze
        
    Returns:
        List[VocabularyItem]: List of vocabulary items with kanji, romaji, english meaning and parts
    """
    try:
        # Load the prompt template
        prompt_template = load_prompt("vocabulary_extractor.md")
        
        # Create the prompt with the lyrics
        prompt = f"{prompt_template}\n\nLyrics:\n{lyrics}"
        
        # Configure instructor with Ollama client
        client = instructor.patch(ollama.chat)
        
        # Get structured response using Instructor
        response = client(
            model="deepseek-r1:latest",
            response_model=List[VocabularyItem],
            messages=[{
                "role": "user", 
                "content": prompt
            }]
        )
        
        return response

    except Exception as e:
        print(f"Error extracting vocabulary: {str(e)}")
        return []
