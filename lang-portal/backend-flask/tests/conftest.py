import pytest
from app import create_app
import os
import sqlite3
from lib.db import Db

@pytest.fixture
def app():
    # Configure app for testing
    app = create_app()
    app.config['TESTING'] = True
    
    # Configure test database
    app.db = Db(database=':memory:')
    
    with app.app_context():
        cursor = app.db.cursor()
        
        # Create tables and test data
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

            -- Insert test data
            INSERT INTO groups (id, name, created_at) 
            VALUES (1, 'Test Group', datetime('now'));

            INSERT INTO study_activities (id, name, created_at) 
            VALUES (1, 'Test Activity', datetime('now'));
        ''')
        app.db.commit()
        
    yield app

    # No need to cleanup in-memory database

@pytest.fixture
def client(app):
    return app.test_client() 