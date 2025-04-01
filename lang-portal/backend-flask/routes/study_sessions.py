from flask import request, jsonify, g
from flask_cors import cross_origin
from datetime import datetime
import math
import sqlite3

def load(app):
  @app.route('/api/study-sessions', methods=['POST'])
  @cross_origin()
  def create_study_session():
    try:
      data = request.get_json()
      print("Received data:", data)  # Debug log
      
      group_id = data.get('group_id')
      study_activity_id = data.get('study_activity_id')
      
      print(f"group_id: {group_id}, study_activity_id: {study_activity_id}")  # Debug log
      
      if not group_id or not study_activity_id:
        return jsonify({'error': 'Missing required fields'}), 400
      
      cursor = app.db.cursor()
      
      # Verify group and activity exist
      cursor.execute('SELECT id FROM groups WHERE id = ?', (group_id,))
      if not cursor.fetchone():
        return jsonify({'error': f'Group {group_id} not found'}), 404
      
      cursor.execute('SELECT id FROM study_activities WHERE id = ?', (study_activity_id,))
      if not cursor.fetchone():
        return jsonify({'error': f'Activity {study_activity_id} not found'}), 404
      
      # Create the study session
      cursor.execute('''
        INSERT INTO study_sessions (group_id, study_activity_id, created_at)
        VALUES (?, ?, ?)
      ''', (group_id, study_activity_id, datetime.utcnow().isoformat()))
      
      app.db.get_connection().commit()
      
      # Get the ID of the created session
      cursor.execute('SELECT last_insert_rowid() as id')
      session_id = cursor.fetchone()['id']
      
      print(f"Created session with ID: {session_id}")  # Debug log
      return jsonify({'id': session_id}), 201
      
    except Exception as e:
      print(f"Error creating study session: {e}")
      return jsonify({'error': str(e)}), 500

  @app.route('/api/study_sessions', methods=['GET'])
  @cross_origin()
  def get_study_sessions():
    try:
      cursor = app.db.cursor()
      
      # Get pagination parameters
      page = request.args.get('page', 1, type=int)
      per_page = request.args.get('per_page', 10, type=int)
      offset = (page - 1) * per_page

      # Get total count
      cursor.execute('SELECT COUNT(*) as count FROM study_sessions')
      total_count = cursor.fetchone()['count']

      # Get paginated sessions with related data
      cursor.execute('''
        SELECT 
          ss.id,
          ss.group_id,
          g.name as group_name,
          sa.id as activity_id,
          sa.name as activity_name,
          ss.created_at,
          COUNT(wri.id) as review_items_count
        FROM study_sessions ss
        JOIN groups g ON g.id = ss.group_id
        JOIN study_activities sa ON sa.id = ss.study_activity_id
        LEFT JOIN word_review_items wri ON wri.study_session_id = ss.id
        GROUP BY ss.id, ss.group_id, g.name, sa.id, sa.name, ss.created_at
        ORDER BY ss.created_at DESC
        LIMIT ? OFFSET ?
      ''', (per_page, offset))
      
      sessions = cursor.fetchall()

      return jsonify({
        'items': [{
          'id': session['id'],
          'group_id': session['group_id'],
          'group_name': session['group_name'],
          'activity_id': session['activity_id'],
          'activity_name': session['activity_name'],
          'start_time': session['created_at'],
          'end_time': session['created_at'],  # Using created_at as end_time for now
          'review_items_count': session['review_items_count']
        } for session in sessions],
        'total': total_count,
        'page': page,
        'per_page': per_page,
        'total_pages': math.ceil(total_count / per_page)
      })
    except Exception as e:
      print("Error fetching sessions:", str(e))
      return jsonify({"error": str(e)}), 500

  @app.route('/api/study_sessions/<int:id>', methods=['GET'])
  @cross_origin()
  def get_study_session(id):
    try:
      cursor = app.db.cursor()
      
      # Get session details
      cursor.execute('''
        SELECT 
          ss.id,
          ss.group_id,
          g.name as group_name,
          sa.id as activity_id,
          sa.name as activity_name,
          ss.created_at,
          COUNT(wri.id) as review_items_count
        FROM study_sessions ss
        JOIN groups g ON g.id = ss.group_id
        JOIN study_activities sa ON sa.id = ss.study_activity_id
        LEFT JOIN word_review_items wri ON wri.study_session_id = ss.id
        WHERE ss.id = ?
        GROUP BY ss.id
      ''', (id,))
      
      session = cursor.fetchone()
      if not session:
        return jsonify({"error": "Study session not found"}), 404

      # Get pagination parameters
      page = request.args.get('page', 1, type=int)
      per_page = request.args.get('per_page', 10, type=int)
      offset = (page - 1) * per_page

      # Get the words reviewed in this session with their review status
      cursor.execute('''
        SELECT 
          w.*,
          COALESCE(SUM(CASE WHEN wri.correct = 1 THEN 1 ELSE 0 END), 0) as session_correct_count,
          COALESCE(SUM(CASE WHEN wri.correct = 0 THEN 1 ELSE 0 END), 0) as session_wrong_count
        FROM words w
        JOIN word_review_items wri ON wri.word_id = w.id
        WHERE wri.study_session_id = ?
        GROUP BY w.id
        ORDER BY w.kanji
        LIMIT ? OFFSET ?
      ''', (id, per_page, offset))
      
      words = cursor.fetchall()

      # Get total count of words
      cursor.execute('''
        SELECT COUNT(DISTINCT w.id) as count
        FROM words w
        JOIN word_review_items wri ON wri.word_id = w.id
        WHERE wri.study_session_id = ?
      ''', (id,))
      
      total_count = cursor.fetchone()['count']

      return jsonify({
        'session': {
          'id': session['id'],
          'group_id': session['group_id'],
          'group_name': session['group_name'],
          'activity_id': session['activity_id'],
          'activity_name': session['activity_name'],
          'start_time': session['created_at'],
          'end_time': session['created_at'],  # For now, just use the same time
          'review_items_count': session['review_items_count']
        },
        'words': [{
          'id': word['id'],
          'kanji': word['kanji'],
          'romaji': word['romaji'],
          'english': word['english'],
          'correct_count': word['session_correct_count'],
          'wrong_count': word['session_wrong_count']
        } for word in words],
        'total': total_count,
        'page': page,
        'per_page': per_page,
        'total_pages': math.ceil(total_count / per_page)
      })
    except Exception as e:
      return jsonify({"error": str(e)}), 500

  @app.route('/api/study-sessions/<int:session_id>/review', methods=['POST'])
  @cross_origin()
  def submit_session_review(session_id):
    try:
      data = request.get_json()
      word_id = data.get('word_id')
      correct = data.get('correct')
      
      if not all([word_id, correct is not None]):
        return jsonify({'error': 'Missing required fields'}), 400
            
      cursor = app.db.cursor()
      
      # Insert the review item
      cursor.execute('''
        INSERT INTO word_review_items (
          study_session_id,
          word_id,
          correct,
          created_at
        ) VALUES (?, ?, ?, datetime('now'))
      ''', (session_id, word_id, correct))
      
      app.db.get_connection().commit()
      
      return jsonify({'message': 'Review submitted successfully'})
      
    except Exception as e:
      print(f"Error submitting review: {e}")
      return jsonify({'error': str(e)}), 500

  @app.route('/api/study_sessions/reset', methods=['POST'])
  @cross_origin()
  def reset_study_sessions():
    try:
      cursor = app.db.cursor()
      
      # First delete all word review items since they have foreign key constraints
      cursor.execute('DELETE FROM word_review_items')
      
      # Then delete all study sessions
      cursor.execute('DELETE FROM study_sessions')
      
      app.db.commit()
      
      return jsonify({"message": "Study history cleared successfully"}), 200
    except Exception as e:
      return jsonify({"error": str(e)}), 500

  @app.route('/api/study-sessions/<int:session_id>/reviews', methods=['GET'])
  @cross_origin()
  def get_session_reviews(session_id):
    try:
      cursor = app.db.cursor()
      cursor.execute('''
        SELECT 
          wri.id,
          wri.word_id,
          wri.correct,
          wri.created_at,
          w.kanji,
          w.english
        FROM word_review_items wri
        JOIN words w ON w.id = wri.word_id
        WHERE wri.study_session_id = ?
        ORDER BY wri.created_at DESC
      ''', (session_id,))
      
      reviews = cursor.fetchall()
      
      return jsonify({
        'reviews': [{
          'id': review['id'],
          'word_id': review['word_id'],
          'correct': review['correct'],
          'created_at': review['created_at'],
          'kanji': review['kanji'],
          'english': review['english']
        } for review in reviews]
      })
    except Exception as e:
      print(f"Error getting reviews: {e}")
      return jsonify({'error': str(e)}), 500