import sqlite3
import json
import os
from flask import g
from datetime import datetime
import threading

class Db:
  def __init__(self, database='word_bank.db'):
    self.database = database
    self._local = threading.local()
    self.setup_database()

  def get_connection(self):
    """Get a thread-local database connection"""
    if not hasattr(self._local, 'connection'):
      if self.database == ':memory:':
        self._local.connection = sqlite3.connect(self.database, check_same_thread=False)
      else:
        self._local.connection = sqlite3.connect(self.database)
      self._local.connection.row_factory = sqlite3.Row
    return self._local.connection

  def cursor(self):
    """Get a cursor from the thread-local connection"""
    return self.get_connection().cursor()

  def commit(self):
    """Commit the current transaction"""
    if hasattr(self._local, 'connection'):
      self._local.connection.commit()

  def close(self):
    """Close the thread-local connection"""
    if hasattr(self._local, 'connection'):
      self._local.connection.close()
      del self._local.connection

  def sql(self, filename):
    """Read SQL from a file"""
    with open(os.path.join('sql', filename)) as file:
      return file.read()

  def load_json(self, filepath):
    """Load data from a JSON file"""
    with open(filepath, 'r') as file:
      return json.load(file)

  def setup_database(self):
    """Set up database tables"""
    conn = sqlite3.connect(self.database)
    conn.row_factory = sqlite3.Row
    cursor = conn.cursor()
    try:
        # Create tables in order using executescript
        cursor.executescript(self.sql('setup/create_table_words.sql'))
        cursor.executescript(self.sql('setup/create_table_groups.sql'))
        cursor.executescript(self.sql('setup/create_table_word_groups.sql'))
        cursor.executescript(self.sql('setup/create_table_study_activities.sql'))
        cursor.executescript(self.sql('setup/create_table_study_sessions.sql'))
        cursor.executescript(self.sql('setup/create_table_word_reviews.sql'))
        cursor.executescript(self.sql('setup/create_table_word_review_items.sql'))
        conn.commit()
    except Exception as e:
        print(f"Error setting up database: {e}")
        conn.rollback()
        raise
    finally:
        cursor.close()
        conn.close()

  def setup_tables(self, cursor):
    # Create the necessary tables
    cursor.execute(self.sql('setup/create_table_words.sql'))
    self.get_connection().commit()

    cursor.execute(self.sql('setup/create_table_word_reviews.sql'))
    self.get_connection().commit()

    cursor.execute(self.sql('setup/create_table_word_review_items.sql'))
    self.get_connection().commit()

    cursor.execute(self.sql('setup/create_table_groups.sql'))
    self.get_connection().commit()

    cursor.execute(self.sql('setup/create_table_word_groups.sql'))
    self.get_connection().commit()

    cursor.execute(self.sql('setup/create_table_study_activities.sql'))
    self.get_connection().commit()

    cursor.execute(self.sql('setup/create_table_study_sessions.sql'))
    self.get_connection().commit()

  def import_study_activities_json(self, cursor, data_json_path):
    # First clear existing activities
    cursor.execute('DELETE FROM study_activities')
    self.get_connection().commit()

    activities = self.load_json(data_json_path)
    for activity in activities:
        cursor.execute('''
            INSERT INTO study_activities (name, url, preview_url) 
            VALUES (?, ?, ?)
        ''', (
            activity['name'],
            activity.get('launch_url', activity['url']),  # Use launch_url if available, fallback to url
            activity['preview_url']
        ))
    self.get_connection().commit()
    print(f"Successfully imported {len(activities)} study activities")

  def import_word_json(self, cursor, group_name, data_json_path):
    # Create the group
    cursor.execute('''
        INSERT INTO groups (name) VALUES (?)
    ''', (group_name,))
    self.get_connection().commit()

    # Get the ID of the group
    cursor.execute('SELECT last_insert_rowid() as id')
    core_verbs_group_id = cursor.fetchone()['id']

    # Load and insert the words
    words = self.load_json(data_json_path)
    for word in words:
        # Insert the word using the existing schema
        cursor.execute('''
            INSERT INTO words (kanji, romaji, english, parts) VALUES (?, ?, ?, ?)
        ''', (
            word['kanji'],
            word.get('romaji', ''),
            word.get('english', ''),
            json.dumps(word.get('parts', {}))
        ))
        
        # Get the ID of the word we just inserted
        cursor.execute('SELECT last_insert_rowid() as id')
        word_id = cursor.fetchone()['id']
        
        # Create the word-group association
        cursor.execute('''
            INSERT INTO word_groups (word_id, group_id) VALUES (?, ?)
        ''', (word_id, core_verbs_group_id))
    self.get_connection().commit()

    # Update the words_count in the groups table
    cursor.execute('''
        UPDATE groups 
        SET words_count = (
            SELECT COUNT(*) 
            FROM word_groups 
            WHERE group_id = ?
        )
        WHERE id = ?
    ''', (core_verbs_group_id, core_verbs_group_id))

    self.get_connection().commit()

    print(f"Successfully added {len(words)} verbs to the '{group_name}' group.")

  # Initialize the database with sample data
  def init(self, app):
    with app.app_context():
      cursor = self.cursor()
      self.setup_tables(cursor)
      self.import_word_json(
        cursor=cursor,
        group_name='Core Verbs',
        data_json_path='seed/data_verbs.json'
      )
      self.import_word_json(
        cursor=cursor,
        group_name='Core Adjectives',
        data_json_path='seed/data_adjectives.json'
      )

      self.import_study_activities_json(
        cursor=cursor,
        data_json_path='seed/study_activities.json'
      )

      self.get_connection().commit()
      return True

  def rollback(self):
    """Rollback the current transaction"""
    if hasattr(self._local, 'connection'):
      self._local.connection.rollback()

# Create an instance of the Db class
db = Db()
