import os
import uuid
import re
from typing import Dict, List, Optional

def save_result(query: str, message: str, output_dir: str = "outputs") -> Optional[str]:
    """
    Parse and save lyrics and vocabulary from the agent's final answer.
    
    Args:
        query (str): The original search query
        message (str): The agent's final answer containing lyrics and vocabulary
        output_dir (str): Directory to save the output files
        
    Returns:
        Optional[str]: A unique identifier for the saved files, or None if parsing failed
    """
    try:
        # Create output directory if it doesn't exist
        os.makedirs(output_dir, exist_ok=True)
        
        # Extract lyrics
        lyrics_match = re.search(r'LYRICS:\n(.*?)\nEND_LYRICS', message, re.DOTALL)
        if not lyrics_match:
            print("Failed to extract lyrics section")
            return None
        lyrics = lyrics_match.group(1).strip()

        # Extract vocabulary
        vocab_match = re.search(r'VOCABULARY:\n(.*?)\nEND_VOCABULARY', message, re.DOTALL)
        if not vocab_match:
            print("Failed to extract vocabulary section")
            return None
        vocab_text = vocab_match.group(1).strip()

        # Parse vocabulary items
        vocab_items = []
        current_item = {}
        for line in vocab_text.split('\n'):
            line = line.strip()
            if line.startswith('- Word:'):
                if current_item:
                    vocab_items.append(current_item)
                current_item = {'word': line[7:].strip()}
            elif line.startswith('Definition:'):
                current_item['definition'] = line[11:].strip()
            elif line.startswith('Example:'):
                current_item['example'] = line[8:].strip()
        if current_item:
            vocab_items.append(current_item)

        # Generate unique ID
        result_id = str(uuid.uuid4())[:8]
        
        # Save files
        base_path = os.path.join(output_dir, result_id)
        
        # Save lyrics
        with open(f"{base_path}_lyrics.txt", 'w', encoding='utf-8') as f:
            f.write(f"Query: {query}\n\n")
            f.write(lyrics)
        
        # Save vocabulary
        with open(f"{base_path}_vocab.txt", 'w', encoding='utf-8') as f:
            for item in vocab_items:
                f.write(f"Word: {item['word']}\n")
                f.write(f"Definition: {item['definition']}\n")
                f.write(f"Example: {item['example']}\n")
                f.write("\n")
        
        print(f"Saved results with ID: {result_id}")
        return result_id
            
    except Exception as e:
        print(f"Failed to save result: {str(e)}")
        return None 