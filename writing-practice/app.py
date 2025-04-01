import streamlit as st
# Set page config must be the first Streamlit command
st.set_page_config(page_title="Writing Practice")

import requests
from PIL import Image
import io
import json
import random
from services import OCRService, grade_submission, create_review_item
from streamlit_drawable_canvas import st_canvas
from flask import Flask, request, jsonify, send_from_directory
from flask_cors import CORS
from lib.db import get_db
from datetime import datetime

app = Flask(__name__)
CORS(app)

# Get group_id from query parameters
def get_group_id():
    group_id = st.query_params.get('group_id', None)
    if not group_id:
        st.error("No group_id provided in URL. Please add ?group_id=YOUR_ID to the URL.")
        st.stop()
    return group_id

# Initialize services
@st.cache_resource
def init_services():
    try:
        return {
            'ocr': OCRService()
        }
    except Exception as e:
        st.error(f"Failed to initialize services: {str(e)}")
        return None

# Initialize session state if not already done
if 'state' not in st.session_state:
    st.session_state.state = 'setup'
if 'current_word' not in st.session_state:
    st.session_state.current_word = {}
if 'vocabulary' not in st.session_state:
    # Fetch vocabulary from API using group_id
    try:
        group_id = get_group_id()
        # Get session_id from URL parameters
        session_id = st.query_params.get('session_id')
        if not session_id:
            st.error("No session_id provided in URL. Please add ?session_id=YOUR_ID to the URL.")
            st.stop()
        st.session_state.session_id = session_id  # Store session_id
        
        # Update the URL to match your backend API structure
        response = requests.get(f'http://localhost:5001/api/groups/{group_id}/words/raw')
        if response.status_code != 200:
            raise Exception(f"API returned status code {response.status_code}")
        st.session_state.vocabulary = response.json()
        print("Fetched vocabulary:", st.session_state.vocabulary)  # Debug log
    except Exception as e:
        st.error(f"Failed to fetch vocabulary: {str(e)}")
        st.session_state.vocabulary = []

services = init_services()

def generate_word():
    """Select a random word from vocabulary for practice."""
    if not st.session_state.vocabulary:
        word = {
            "id": 1,  # Add default ID
            "kanji": "飲む",
            "romaji": "nomu",
            "english": "to drink",
            "parts": ["飲"]
        }
    else:
        try:
            word = random.choice(st.session_state.vocabulary)
            word = {
                "id": word["id"],  # Keep the word ID
                "kanji": word["kanji"],
                "romaji": word["romaji"],
                "english": word["english"],
                "parts": json.loads(word["parts"]) if word.get("parts") else []
            }
        except Exception as e:
            st.error(f"Failed to get word: {str(e)}")
            word = {
                "id": 1,  # Add default ID
                "kanji": "飲む",
                "romaji": "nomu",
                "english": "to drink",
                "parts": ["飲"]
            }
    
    # Store the word ID in session state
    st.session_state.current_word_id = word["id"]
    return word

def submit_review(session_id: int, word_id: int, is_correct: bool):
    """Submit the review result to the backend."""
    try:
        session_id = int(session_id)
        word_id = int(word_id)
        
        print("\n=== Submit Review Debug ===")
        print(f"Session ID: {session_id}")
        print(f"Word ID: {word_id}")
        print(f"Is Correct: {is_correct}")
        
        url = f'http://localhost:5001/api/study-sessions/{session_id}/review'
        payload = {
            "word_id": word_id,
            "correct": is_correct
        }
        
        print(f"Making POST request to: {url}")
        print(f"With payload: {payload}")
        
        response = requests.post(url, json=payload)
        
        print(f"Response status: {response.status_code}")
        print(f"Response body: {response.text}")
        
        if response.status_code != 200:
            raise Exception(f"API returned status code {response.status_code}: {response.text}")
        
        print("Review submitted successfully!")
        return True
        
    except Exception as e:
        print(f"Review submission error: {str(e)}")
        st.error(f"Failed to submit review: {str(e)}")
        return False

def handle_state_transitions():
    if st.session_state.state == 'setup':
        st.title("Japanese Writing Practice")
        
        if st.button("Get New Word"):
            with st.spinner("Getting word..."):
                word = generate_word()
                st.session_state.current_word = word
                st.session_state.state = 'writing'
                st.rerun()

    elif st.session_state.state == 'writing':
        st.title("Japanese Writing Practice")
        
        # Word Information Box
        with st.container():
            st.subheader("Practice Word")
            col1, col2, col3 = st.columns(3)
            with col1:
                st.write("Kanji")
                st.write(st.session_state.current_word['kanji'])
            with col2:
                st.write("Romaji")
                st.write(st.session_state.current_word['romaji'])
            with col3:
                st.write("English")
                st.write(st.session_state.current_word['english'])
            
            # Show kanji parts if available
            if st.session_state.current_word.get('parts'):
                st.write("Kanji Components:")
                for part in st.session_state.current_word['parts']:
                    st.write(f"- {part}")
        
        st.divider()
        st.write("Please write this word in Japanese:")
        st.write(st.session_state.current_word['kanji'])
        
        # Add tabs for different input methods
        tab1, tab2, tab3 = st.tabs(["Type Answer", "Upload Handwriting", "Sketch Pad"])
        
        with tab1:
            typed_answer = st.text_input("Type your answer in Japanese:", 
                placeholder="ここに日本語で書いてください")
            if st.button("Submit Typed Answer"):
                with st.spinner("Analyzing your submission..."):
                    print(f"Submitting answer: {typed_answer}")  # Debug log
                    print(f"Correct word: {st.session_state.current_word['kanji']}")  # Debug log
                    
                    # Store the answer in session state
                    st.session_state.user_input = typed_answer
                    results = grade_submission(st.session_state.current_word['kanji'], typed_answer)
                    
                    if results:
                        print(f"Got results: {results}")  # Debug log
                        st.session_state.review_results = results
                        st.session_state.state = 'grading'
                        st.rerun()
        
        with tab2:
            uploaded_file = st.file_uploader("Upload your handwritten Japanese", 
                type=['png', 'jpg', 'jpeg'])
            if uploaded_file and st.button("Submit Image"):
                with st.spinner("Analyzing your submission..."):
                    image = Image.open(uploaded_file)
                    results = grade_submission(image, st.session_state.current_word['kanji'])
                    if results:
                        st.session_state.review_results = results
                        st.session_state.state = 'grading'
                        st.rerun()
        
        with tab3:
            # Add sketch pad
            st.write("Practice writing here:")
            canvas = st.empty()
            with canvas:
                stroke_width = st.slider("Brush size", 1, 25, 3)
                canvas_result = st_canvas(
                    fill_color="rgba(255, 255, 255, 0.0)",
                    stroke_width=stroke_width,
                    stroke_color="#000000",
                    background_color="#ffffff",
                    height=400,
                    drawing_mode="freedraw",
                    key=f"canvas_{st.session_state.canvas_key}",
                )
            
            col1, col2 = st.columns(2)
            with col1:
                if st.button("Clear Canvas"):
                    # Reset the canvas by changing its key
                    st.session_state.canvas_key = random.randint(0, 1000000)
                    st.rerun()
            
            with col2:
                if st.button("Submit Sketch") and canvas_result.image_data is not None:
                    with st.spinner("Analyzing your submission..."):
                        # Convert canvas to PIL Image
                        image_data = canvas_result.image_data
                        image = Image.fromarray(image_data.astype('uint8'), 'RGBA')
                        # Convert to RGB (OCR expects RGB)
                        image = image.convert('RGB')
                        
                        results = grade_submission(image, st.session_state.current_word['kanji'])
                        if results:
                            st.session_state.review_results = results
                            st.session_state.state = 'grading'
                            st.rerun()

    elif st.session_state.state == 'grading':
        st.title("Japanese Writing Practice")
        st.write("Original word:")
        st.write(st.session_state.current_word['kanji'])
        
        results = st.session_state.review_results
        
        st.subheader("Review Results")
        st.write("Your answer:", st.session_state.get('user_input', ''))
        st.write("Correct answer:", st.session_state.current_word['kanji'])
        st.write("Result:", "Correct!" if results['correct'] else "Not quite.")  # Simplified result message
        
        # Add buttons to mark as correct/incorrect
        col1, col2 = st.columns(2)
        with col1:
            if st.button("Continue"):
                success = submit_review(
                    session_id=st.session_state.session_id,
                    word_id=st.session_state.current_word['id'],
                    is_correct=results['correct']
                )
                if success:
                    st.session_state.state = 'setup'
                    st.rerun()
                else:
                    st.error("Failed to submit review. Please try again.")
        
        with col2:
            if st.button("Try Again"):
                st.session_state.state = 'writing'
                st.rerun()

# Initialize canvas key if not exists
if 'canvas_key' not in st.session_state:
    st.session_state.canvas_key = 0

# Main app flow
if services:
    handle_state_transitions()
else:
    st.error("Application failed to initialize. Please check the logs and try again.")

@app.route('/api/check', methods=['POST'])
def check_answer():
    try:
        data = request.get_json()
        word = data.get('word')
        user_input = data.get('user_input')
        session_id = data.get('session_id')
        word_id = data.get('word_id')
        
        if not all([word, user_input, session_id, word_id]):
            return jsonify({'error': 'Missing required fields'}), 400
            
        # Use local grading instead of OpenAI
        result = grade_submission(word, user_input)
        
        # Record the attempt
        create_review_item(
            db=get_db(),
            session_id=session_id,
            word_id=word_id,
            is_correct=result['correct']
        )
        
        return jsonify(result)
        
    except Exception as e:
        print(f"Error checking answer: {e}")
        return jsonify({'error': str(e)}), 500

if __name__ == '__main__':
    pass  # Remove the set_page_config from here
