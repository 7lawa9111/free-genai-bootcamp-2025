from typing import Dict, Optional
import boto3
import json
import os
from dotenv import load_dotenv
from .vector_store import JLPTQuestionStore

# Load environment variables
load_dotenv()

class QuestionGenerator:
    def __init__(self):
        # Initialize Bedrock client using environment variables
        self.bedrock = boto3.client(
            service_name='bedrock-runtime',
            region_name='us-east-2',
            aws_access_key_id=os.getenv('AWS_ACCESS_KEY_ID'),
            aws_secret_access_key=os.getenv('AWS_SECRET_ACCESS_KEY')
        )
        self.model_id = 'us.amazon.nova-micro-v1:0'
        self.question_store = JLPTQuestionStore()

    def generate_question(self, practice_type: str) -> Optional[Dict]:
        """
        Generate a question based on practice type using RAG
        
        Args:
            practice_type (str): Type of practice (Dialogue, Vocabulary, Listening)
            
        Returns:
            Dict: Generated question with format:
                {
                    'introduction': str,
                    'conversation': str,
                    'question': str,
                    'options': List[str],
                    'correct_answer': str,
                    'audio_url': Optional[str]
                }
        """
        # Map practice types to search queries
        query_map = {
            "Dialogue Practice": "会話 デパート 店",
            "Vocabulary Quiz": "単語 言葉",
            "Listening Exercise": "聞く 音声"
        }
        
        # Get similar questions from vector store
        similar_questions = self.question_store.find_similar_questions(
            query_map.get(practice_type, "会話"),
            n_results=3
        )
        
        # Create context from similar questions
        context = "\n\n".join([
            f"Example {i+1}:\n"
            f"Introduction: {q['introduction']}\n"
            f"Conversation: {q['conversation']}\n"
            f"Question: {q['question']}"
            for i, q in enumerate(similar_questions)
        ])
        
        # Create prompt for generation
        prompt = f"""
        Using these example JLPT listening questions as reference:
        
        {context}
        
        Generate a new {practice_type} question in Japanese, following the same style as the examples.
        The question should be:
        1. Entirely in Japanese (including introduction, conversation, and question)
        2. At JLPT N4-N5 level difficulty
        3. Natural conversational Japanese with appropriate politeness levels
        4. Include appropriate speaker markers (男:, 女:, 店員: etc.)

        Format the response exactly like this:
        Introduction: [Japanese introduction text]
        Conversation: 
        [Speaker marker]: [Japanese dialogue]
        [Speaker marker]: [Japanese dialogue]
        Question: [Japanese question text]
        Options:
        A) [Japanese option]
        B) [Japanese option]
        C) [Japanese option]
        D) [Japanese option]
        Correct: [A/B/C/D]

        Make sure all text except the format markers and A/B/C/D is in Japanese.
        """
        
        try:
            # Generate new question using Bedrock
            messages = [{
                "role": "user",
                "content": [{"text": prompt}]
            }]
            
            response = self.bedrock.converse(
                modelId=self.model_id,
                messages=messages,
                inferenceConfig={
                    "temperature": 0.7,
                    "topP": 0.9,
                    "maxTokens": 2000,
                    "stopSequences": []
                }
            )
            
            generated_text = response['output']['message']['content'][0]['text']
            
            # Parse the response
            lines = generated_text.split('\n')
            result = {}
            current_field = None
            options = []
            
            for line in lines:
                line = line.strip()
                if line.startswith('Introduction:'):
                    current_field = 'introduction'
                    result[current_field] = line[13:].strip()
                elif line.startswith('Conversation:'):
                    current_field = 'conversation'
                    result[current_field] = line[13:].strip()
                elif line.startswith('Question:'):
                    current_field = 'question'
                    result[current_field] = line[9:].strip()
                elif line.startswith('Options:'):
                    current_field = 'options'
                elif line.startswith(('A)', 'B)', 'C)', 'D)')):
                    options.append(line[3:].strip())
                elif line.startswith('Correct:'):
                    result['correct_answer'] = line[8:].strip()
                elif current_field and line:
                    result[current_field] += '\n' + line
            
            result['options'] = options
            result['audio_url'] = None  # Would be generated by text-to-speech
            
            return result
            
        except Exception as e:
            print(f"Error generating question: {str(e)}")
            return None 