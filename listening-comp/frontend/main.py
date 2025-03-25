import streamlit as st
from backend.question_generator import QuestionGenerator
from backend.audio_generator import AudioGenerator
from datetime import datetime
import json
import os

# Page config
st.set_page_config(
    page_title="JLPT Question Practice",
    page_icon="🎌",
    layout="wide"
)

HISTORY_FILE = "question_history.json"

def format_timestamp(timestamp):
    """Format timestamp to readable string"""
    if isinstance(timestamp, str):
        timestamp = datetime.fromisoformat(timestamp)
    return timestamp.strftime("%Y-%m-%d %H:%M")

def load_question_history():
    """Load question history from file"""
    try:
        if os.path.exists(HISTORY_FILE):
            with open(HISTORY_FILE, 'r', encoding='utf-8') as f:
                return json.load(f)
    except Exception as e:
        print(f"Error loading history: {e}")
    return []

def save_question_history(questions):
    """Save question history to file"""
    try:
        # Convert datetime objects to strings for JSON serialization
        questions_to_save = []
        for q in questions:
            q_copy = q.copy()
            if isinstance(q_copy['timestamp'], datetime):
                q_copy['timestamp'] = q_copy['timestamp'].isoformat()
            questions_to_save.append(q_copy)
            
        with open(HISTORY_FILE, 'w', encoding='utf-8') as f:
            json.dump(questions_to_save, f, ensure_ascii=False, indent=2)
    except Exception as e:
        print(f"Error saving history: {e}")

def render_interactive_stage():
    """Render the interactive learning stage"""
    # Initialize session state for storing generated questions
    if 'generated_questions' not in st.session_state:
        st.session_state.generated_questions = load_question_history()
    
    # Sidebar for question history
    with st.sidebar:
        st.title("Question History")
        
        if st.session_state.generated_questions:
            st.write("Previously Generated Questions:")
            # Group questions by practice type
            practice_types = set(q['practice_type'] for q in st.session_state.generated_questions)
            
            for practice_type in practice_types:
                st.subheader(practice_type)
                type_questions = [q for q in st.session_state.generated_questions 
                                if q['practice_type'] == practice_type]
                
                for i, q in enumerate(type_questions):
                    # Create expandable section for each question
                    with st.expander(f"{format_timestamp(q['timestamp'])}"):
                        st.write(f"{q['introduction'][:50]}...")
                        if st.button("Load Question", key=f"q_{q['timestamp']}"):
                            st.session_state.current_question = q
                            st.session_state.answered = False
                            st.session_state.feedback = None
            
            if st.button("Clear History"):
                st.session_state.generated_questions = []
                if 'current_question' in st.session_state:
                    del st.session_state.current_question
                save_question_history([])  # Clear the saved history
        else:
            st.info("No questions generated yet")
    
    # Main content
    st.title("🎌 JLPT Question Practice")
    
    # Initialize question generator
    if 'question_generator' not in st.session_state:
        st.session_state.question_generator = QuestionGenerator()
    
    # Initialize audio generator
    if 'audio_generator' not in st.session_state:
        st.session_state.audio_generator = AudioGenerator()
    
    # Practice type selection
    practice_type = st.selectbox(
        "Select Practice Type",
        ["Dialogue Practice", "Vocabulary Quiz", "Listening Exercise"]
    )
    
    # Generate new question button
    if st.button("Generate New Question"):
        question = st.session_state.question_generator.generate_question(practice_type)
        if question:
            # Add metadata to question
            question['timestamp'] = datetime.now()
            question['practice_type'] = practice_type
            
            # Add to history
            st.session_state.generated_questions.append(question)
            # Save to file
            save_question_history(st.session_state.generated_questions)
            # Set as current question
            st.session_state.current_question = question
            st.session_state.answered = False
            st.session_state.feedback = None
    
    # Display question if available
    if 'current_question' in st.session_state:
        question = st.session_state.current_question
        
        col1, col2 = st.columns([2, 1])
        
        with col1:
            st.subheader("Practice Scenario")
            st.write("**Introduction:**")
            st.write(question['introduction'])
            
            st.write("**Conversation:**")
            st.write(question['conversation'])
            
            st.write("**Question:**")
            st.write(question['question'])
            
            # Multiple choice options
            selected = st.radio(
                "Choose your answer:",
                [f"{opt[0]}) {opt[1]}" for opt in zip(['A', 'B', 'C', 'D'], question['options'])]
            )
            
            # Check answer button
            if not st.session_state.get('answered', False) and st.button("Check Answer"):
                st.session_state.answered = True
                selected_letter = selected[0]
                if selected_letter == question['correct_answer']:
                    st.session_state.feedback = "✅ Correct!"
                else:
                    st.session_state.feedback = f"❌ Incorrect. The correct answer is {question['correct_answer']}"
        
        with col2:
            st.subheader("Audio")
            if not question.get('audio_url') and st.button("Generate Audio"):
                with st.spinner("Generating audio..."):
                    audio_file = st.session_state.audio_generator.generate_question_audio(question)
                    if audio_file:
                        question['audio_url'] = audio_file
                        # Update question in history
                        for q in st.session_state.generated_questions:
                            if q['timestamp'] == question['timestamp']:
                                q['audio_url'] = audio_file
                        save_question_history(st.session_state.generated_questions)
                        st.rerun()
                    else:
                        st.error("Failed to generate audio")
            
            if question.get('audio_url'):
                st.audio(question['audio_url'])
            
            st.subheader("Feedback")
            if st.session_state.get('feedback'):
                st.write(st.session_state.feedback)
            else:
                st.info("Select an answer and click 'Check Answer'")

if __name__ == "__main__":
    render_interactive_stage()