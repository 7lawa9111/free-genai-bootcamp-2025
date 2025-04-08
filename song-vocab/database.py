import sqlite3
from typing import List, Dict
from contextlib import contextmanager

class Database:
    def __init__(self, db_path: str = "song_vocabulary.db"):
        self.db_path = db_path
        self._init_db()

    @contextmanager
    def get_connection(self):
        conn = sqlite3.connect(self.db_path)
        try:
            yield conn
        finally:
            conn.close()

    def _init_db(self):
        """Initialize the database with required tables"""
        with self.get_connection() as conn:
            cursor = conn.cursor()
            
            # Create songs table
            cursor.execute("""
                CREATE TABLE IF NOT EXISTS songs (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    title TEXT NOT NULL,
                    lyrics TEXT NOT NULL,
                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
                )
            """)
            
            # Create vocabulary table
            cursor.execute("""
                CREATE TABLE IF NOT EXISTS vocabulary (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    song_id INTEGER,
                    word TEXT NOT NULL,
                    definition TEXT NOT NULL,
                    example TEXT NOT NULL,
                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                    FOREIGN KEY (song_id) REFERENCES songs (id)
                )
            """)
            
            conn.commit()

    def save_song_and_vocabulary(self, title: str, lyrics: str, vocabulary: List[Dict]) -> int:
        """Save a song and its vocabulary items to the database"""
        with self.get_connection() as conn:
            cursor = conn.cursor()
            
            # Insert song
            cursor.execute(
                "INSERT INTO songs (title, lyrics) VALUES (?, ?)",
                (title, lyrics)
            )
            song_id = cursor.lastrowid
            
            # Insert vocabulary items
            for item in vocabulary:
                cursor.execute(
                    "INSERT INTO vocabulary (song_id, word, definition, example) VALUES (?, ?, ?, ?)",
                    (song_id, item['word'], item['definition'], item['example'])
                )
            
            conn.commit()
            return song_id

    def get_song_vocabulary(self, song_id: int) -> Dict:
        """Get a song and its vocabulary items by ID"""
        with self.get_connection() as conn:
            cursor = conn.cursor()
            
            # Get song
            cursor.execute("SELECT title, lyrics FROM songs WHERE id = ?", (song_id,))
            song_row = cursor.fetchone()
            if not song_row:
                return None
            
            # Get vocabulary items
            cursor.execute(
                "SELECT word, definition, example FROM vocabulary WHERE song_id = ?",
                (song_id,)
            )
            vocab_rows = cursor.fetchall()
            
            return {
                "title": song_row[0],
                "lyrics": song_row[1],
                "vocabulary": [
                    {
                        "word": row[0],
                        "definition": row[1],
                        "example": row[2]
                    }
                    for row in vocab_rows
                ]
            }

    def get_all_songs(self) -> List[Dict]:
        """Get all songs with their vocabulary counts"""
        with self.get_connection() as conn:
            cursor = conn.cursor()
            cursor.execute("""
                SELECT s.id, s.title, COUNT(v.id) as vocab_count 
                FROM songs s 
                LEFT JOIN vocabulary v ON s.id = v.song_id 
                GROUP BY s.id, s.title
            """)
            rows = cursor.fetchall()
            return [
                {
                    "id": row[0],
                    "title": row[1],
                    "vocabulary_count": row[2]
                }
                for row in rows
            ] 