from manga_ocr import MangaOcr

class OCRService:
    def __init__(self):
        try:
            self.mocr = MangaOcr()
        except Exception as e:
            raise Exception(f"Failed to initialize MangaOCR: {str(e)}")

    def transcribe_image(self, image) -> str:
        """Transcribe Japanese text from the image."""
        try:
            text = self.mocr(image)
            return text
        except Exception as e:
            raise Exception(f"Failed to transcribe image: {str(e)}")

def grade_submission(word, user_input):
    """
    Grade a user's writing submission by comparing it with the correct word
    
    Args:
        word (str): The correct word/text to compare against
        user_input (str): The user's submitted text
        
    Returns:
        dict: Result containing correct (bool) and feedback
    """
    # Clean up inputs - remove whitespace
    word = word.strip()
    user_input = user_input.strip()
    
    # Simple exact match comparison
    is_correct = word == user_input
    
    return {
        'correct': is_correct,
        'feedback': "Correct!" if is_correct else "Not quite."
    }

def create_review_item(db, session_id, word_id, is_correct):
    """Create a review item record for the attempt"""
    try:
        cursor = db.cursor()
        cursor.execute('''
            INSERT INTO word_review_items (
                study_session_id, 
                word_id,
                correct,
                created_at
            ) VALUES (?, ?, ?, datetime('now'))
        ''', (session_id, word_id, is_correct))
        db.commit()
        return True
    except Exception as e:
        print(f"Error creating review item: {e}")
        return False 