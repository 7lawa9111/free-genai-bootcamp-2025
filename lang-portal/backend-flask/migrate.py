import sqlite3
import os
from datetime import datetime
from flask import Flask

# Create Flask app
app = Flask(__name__)

# Configure database path
db_path = os.path.join(os.path.dirname(__file__), 'word_bank.db')
app.config['DATABASE'] = db_path

def init_tables(cursor):
    # Create tables
    cursor.executescript('''
        CREATE TABLE IF NOT EXISTS groups (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            created_at DATETIME NOT NULL
        );

        CREATE TABLE IF NOT EXISTS study_activities (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            created_at DATETIME NOT NULL
        );

        CREATE TABLE IF NOT EXISTS study_sessions (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            group_id INTEGER NOT NULL,
            study_activity_id INTEGER NOT NULL,
            created_at DATETIME NOT NULL,
            FOREIGN KEY (group_id) REFERENCES groups(id),
            FOREIGN KEY (study_activity_id) REFERENCES study_activities(id)
        );

        CREATE TABLE IF NOT EXISTS word_review_items (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            study_session_id INTEGER NOT NULL,
            word_id INTEGER NOT NULL,
            correct BOOLEAN NOT NULL,
            created_at DATETIME NOT NULL,
            FOREIGN KEY (study_session_id) REFERENCES study_sessions(id),
            FOREIGN KEY (word_id) REFERENCES words(id)
        );
    ''')

def init_test_data(cursor):
    # Add test data
    cursor.execute('''
        INSERT OR IGNORE INTO groups (id, name, created_at)
        VALUES (1, 'Test Group', datetime('now'))
    ''')

    cursor.execute('''
        INSERT OR IGNORE INTO study_activities (id, name, created_at)
        VALUES (1, 'Test Activity', datetime('now'))
    ''')

if __name__ == '__main__':
    # Connect to database
    conn = sqlite3.connect(db_path)
    conn.row_factory = sqlite3.Row
    cursor = conn.cursor()

    # Initialize tables and data
    init_tables(cursor)
    init_test_data(cursor)
    
    conn.commit()
    conn.close()
    print("Database initialized successfully!")
