import streamlit as st
from backend.question_generator import QuestionGenerator

# Page config
st.set_page_config(
    page_title="JLPT Question Practice",
    page_icon="🎌",
    layout="wide"
)

def render_interactive_stage():
    """Render the interactive learning stage"""
    st.title("🎌 JLPT Question Practice")
    
    # Initialize question generator
    if 'question_generator' not in st.session_state:
        st.session_state.question_generator = QuestionGenerator()
    
    # Practice type selection
    practice_type = st.selectbox(
        "Select Practice Type",
        ["Dialogue Practice", "Vocabulary Quiz", "Listening Exercise"]
    )
    
    # Generate new question button
    if st.button("Generate New Question"):
        question = st.session_state.question_generator.generate_question(practice_type)
        if question:
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
            if question.get('audio_url'):
                st.subheader("Audio")
                st.audio(question['audio_url'])
            
            st.subheader("Feedback")
            if st.session_state.get('feedback'):
                st.write(st.session_state.feedback)
            else:
                st.info("Select an answer and click 'Check Answer'")

if __name__ == "__main__":
    render_interactive_stage()