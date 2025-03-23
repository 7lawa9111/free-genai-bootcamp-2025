from typing import List, Dict, Optional
import boto3
import json
import os
import chromadb
from chromadb.utils import embedding_functions
import uuid

class JLPTQuestionStore:
    def __init__(self):
        # Initialize ChromaDB client
        self.client = chromadb.Client()
        
        # Create collections for different parts of questions
        self.questions_collection = self.client.create_collection(
            name="jlpt_questions",
            embedding_function=embedding_functions.SentenceTransformerEmbeddingFunction()
        )

    def add_question(self, section_num: int, question_num: str, question_data: Dict):
        """
        Add a question to the vector store
        """
        # Create a unique ID for the question
        question_id = f"s{section_num}_q{question_num}"
        
        # Combine all text for embedding
        full_text = f"""
        Section: {section_num}
        Question: {question_num}
        Introduction: {question_data['introduction']}
        Conversation: {question_data['conversation']}
        Question: {question_data['question']}
        """
        
        # Store in ChromaDB
        self.questions_collection.add(
            ids=[question_id],
            documents=[full_text],
            metadatas=[{
                "section_num": section_num,
                "question_num": question_num,
                "introduction": question_data['introduction'],
                "conversation": question_data['conversation'],
                "question": question_data['question']
            }]
        )

    def find_similar_questions(self, query: str, n_results: int = 5) -> List[Dict]:
        """
        Find similar questions using vector similarity
        """
        results = self.questions_collection.query(
            query_texts=[query],
            n_results=n_results
        )
        return results['metadatas'][0]  # Return metadata of similar questions

    def get_questions_by_section(self, section_num: int) -> List[Dict]:
        """
        Get all questions from a specific section
        """
        results = self.questions_collection.get(
            where={"section_num": section_num}
        )
        return results['metadatas']

class JLPTQuestionExtractor:
    def __init__(self):
        # Option 1: Explicitly pass credentials
        self.bedrock = boto3.client(
            service_name='bedrock-runtime',
            region_name='us-east-2',
            aws_access_key_id='',
            aws_secret_access_key='',
            #aws_session_token='YOUR_SESSION_TOKEN'  # if using temporary credentials
        )
        self.model_id = 'us.amazon.nova-micro-v1:0'

        # Option 2: Use AWS CLI configuration
        # First run: aws configure
        # Then you can use just:
        # self.bedrock = boto3.client(
        #     service_name='bedrock-runtime',
        #     region_name='us-east-1'
        # )

        # Initialize question store
        self.question_store = JLPTQuestionStore()

    def _invoke_bedrock(self, prompt: str) -> Optional[str]:
        """
        Invoke Amazon Bedrock with the given prompt
        
        Args:
            prompt (str): The prompt to send to the model
            
        Returns:
            Optional[str]: The model's response or None if failed
        """
        try:
            messages = [{
                "role": "user",
                "content": [{"text": prompt}]
            }]

            response = self.bedrock.converse(
                modelId=self.model_id,
                messages=messages,
                inferenceConfig={"temperature": 0.1}
            )
            
            return response['output']['message']['content'][0]['text']
            
        except Exception as e:
            print(f"Error invoking Bedrock: {str(e)}")
            return None

    def extract_questions(self, transcript: List[Dict]) -> List[Dict]:
        # Combine transcript text
        full_text = "\n".join([entry['text'] for entry in transcript])
        
        # Create prompt for the model
        prompt = f"""
        You are a JLPT listening test analyzer. This transcript contains three distinct sections (問題1, 問題2, 問題3), 
        each with its own format and instructions.

        First, identify each section and its format instructions. Then extract the questions from each section.
        
        Format the output exactly like this:

        Section #1:
        Format Instructions:
        [Extract the instructions explaining how questions in this section work]
        
        Questions:
        Question #[number]:
        (Use 例 for example question, or the actual number like 1, 2, etc.)

        Introduction:
        [The one-line setup that introduces the situation]

        Conversation:
        [All dialogue between speakers, including speaker markers]

        Question:
        [The final question that tests comprehension]
        ---

        Section #2:
        Format Instructions:
        [Extract the instructions explaining how questions in this section work]
        
        Questions:
        [Same format as Section 1 questions]

        Section #3:
        Format Instructions:
        [Extract the instructions explaining how questions in this section work]
        
        Questions:
        [Same format as Section 1 questions]

        Here's the transcript to analyze:
        {full_text}

        Make sure to:
        - Clearly separate the three main sections
        - Include the format instructions for each section
        - Number questions within each section separately
        - Include the complete dialogue in conversations
        - Keep speaker markers (男:, 女:, etc.) in conversations
        - Don't leave any section empty
        """

        # Get model response
        response = self._invoke_bedrock(prompt)
        if not response:
            return []

        # Parse the text response into structured data
        sections = []
        current_section = None
        current_section_data = {"format_instructions": "", "questions": []}
        current_question = None
        
        for line in response.split('\n'):
            line = line.strip()
            if not line:
                continue
                
            if line.startswith("Section #"):
                if current_section_data["questions"]:
                    sections.append(current_section_data.copy())
                current_section_data = {"format_instructions": "", "questions": []}
                current_section = "section"
            elif line == "Format Instructions:":
                current_section = "format_instructions"
            elif line == "Questions:":
                current_section = "questions"
            elif line.startswith("Question #"):
                if current_question and any(current_question.values()):
                    current_section_data["questions"].append(current_question.copy())
                current_question = {"introduction": "", "conversation": "", "question": ""}
            elif line.startswith("Introduction:"):
                current_section = "introduction"
            elif line.startswith("Conversation:"):
                current_section = "conversation"
            elif line.startswith("Question:"):
                current_section = "question"
            elif line == "---":
                if current_question and any(current_question.values()):
                    current_section_data["questions"].append(current_question.copy())
                    current_question = {"introduction": "", "conversation": "", "question": ""}
            elif current_section == "format_instructions":
                current_section_data["format_instructions"] += line + "\n"
            elif current_section in ["introduction", "conversation", "question"] and current_question:
                current_question[current_section] += line + "\n"

        # Add the last section if it exists
        if current_section_data["questions"]:
            sections.append(current_section_data)

        # Clean up the text
        for section in sections:
            section["format_instructions"] = section["format_instructions"].strip()
            for question in section["questions"]:
                for key in question:
                    question[key] = question[key].strip()

        return sections

    def save_structured_data(self, sections: List[Dict], filename: str) -> bool:
        """
        Save structured question data to text file
        
        Args:
            sections (List[Dict]): The extracted sections
            filename (str): Output filename
            
        Returns:
            bool: True if successful, False otherwise
        """
        try:
            os.makedirs("./structured_data", exist_ok=True)
            
            with open(f"./structured_data/{filename}.txt", 'w', encoding='utf-8') as f:
                for section_num, section in enumerate(sections, 1):
                    f.write(f"Section #{section_num}:\n\n")
                    f.write("Format Instructions:\n")
                    f.write(section["format_instructions"] + "\n\n")
                    f.write("Questions:\n\n")
                    
                    for i, question in enumerate(section["questions"], 1):
                        question_num = "例" if i == 1 else str(i-1)
                        
                        # Add to vector store
                        self.question_store.add_question(section_num, question_num, question)
                        
                        # Write to file as before
                        f.write(f"Question #{question_num}:\n\n")
                        f.write("Introduction:\n")
                        f.write(question['introduction'] + "\n\n")
                        f.write("Conversation:\n")
                        f.write(question['conversation'] + "\n\n")
                        f.write("Question:\n")
                        f.write(question['question'] + "\n")
                        f.write("---\n\n")
                    f.write("\n")
            return True
        except Exception as e:
            print(f"Error saving structured data: {str(e)}")
            return False

def process_transcript(transcript: List[Dict], output_filename: str) -> bool:
    """
    Process transcript and extract structured question data
    
    Args:
        transcript (List[Dict]): The transcript data
        output_filename (str): Output filename for structured data
        
    Returns:
        bool: True if successful, False otherwise
    """
    extractor = JLPTQuestionExtractor()
    sections = extractor.extract_questions(transcript)
    
    if sections:
        if extractor.save_structured_data(sections, output_filename):
            print(f"Structured data saved successfully to {output_filename}.txt")
            return True
        else:
            print("Failed to save structured data")
    else:
        print("Failed to extract sections from transcript")
    
    return False

if __name__ == "__main__":
    # Read the existing transcript file
    video_id = "sY7L5cfCWno"
    transcript = []
    
    try:
        with open(f"./transcripts/{video_id}.txt", 'r', encoding='utf-8') as f:
            for line in f:
                transcript.append({"text": line.strip()})
        
        # Process the transcript
        if process_transcript(transcript, video_id):
            print("Test processing completed successfully")
            
            # Read and print the results
            try:
                with open(f"./structured_data/{video_id}.txt", 'r', encoding='utf-8') as f:
                    print("\nExtracted Sections:")
                    print(f.read())
            except Exception as e:
                print(f"Error reading results: {str(e)}")
    except Exception as e:
        print(f"Error reading transcript file: {str(e)}")
    
    # Example of finding similar questions
    extractor = JLPTQuestionExtractor()
    similar_questions = extractor.question_store.find_similar_questions(
        "デパートで道を聞く会話", 
        n_results=3
    )
    print("\nSimilar questions:")
    for q in similar_questions:
        print(f"Section {q['section_num']}, Question {q['question_num']}:")
        print(q['introduction'])
        print("---")
    