from typing import List, Dict, Optional
import boto3
import json
import os

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
        """
        Extract JLPT listening questions from transcript
        
        Args:
            transcript (List[Dict]): The transcript data
            
        Returns:
            List[Dict]: List of extracted questions with structure:
                       {
                           'introduction': str,
                           'conversation': str,
                           'question': str
                       }
        """
        # Combine transcript text
        full_text = "\n".join([entry['text'] for entry in transcript])
        
        # Create prompt for the model
        prompt = f"""
        You are a JLPT listening test analyzer. This transcript contains a test introduction, multiple listening questions, and a test conclusion.
        
        For each actual test question (ones that start with 例 or numbers like 1番, 2番, etc.):
        1. Introduction: Extract the brief setup line that introduces the situation (e.g., "例家で女の人が男の人と話しています" or "1番デパートで男の人と店の人が話しています")
        2. Conversation: Extract all the dialogue between speakers. The conversation usually:
           - Contains multiple lines of dialogue
           - Includes speaker markers like 男:, 女:, 店員:, etc.
           - Comes between the introduction and the final question
        3. Question: Extract the specific question being asked at the end (the part that ends with ですか/か)

        Ignore:
        - General test instructions (like 問題2では初めに質問を聞いてください)
        - Test section markers (like 問題1, 問題2)
        - Test conclusions

        Format each question exactly like this, with three dashes between questions:

        Question #[number]:
        (Use 例 for example question, or the actual number like 1, 2, etc.)

        Introduction:
        [The one-line setup that introduces the situation]

        Conversation:
        [All dialogue between speakers, including speaker markers]

        Question:
        [The final question that tests comprehension]
        ---

        Here's the transcript to analyze:
        {full_text}

        Make sure to:
        - Include the question number (例 for first question, actual number for others)
        - Only include actual test questions
        - Include the complete dialogue in the Conversation section
        - Keep speaker markers (男:, 女:, etc.) in the conversation
        - Put only the final question in the Question section
        - Don't leave any section empty
        - Include all lines of dialogue between the introduction and question
        """

        # Get model response
        response = self._invoke_bedrock(prompt)
        if not response:
            return []

        # Parse the text response into structured data
        questions = []
        current_section = None
        current_question = {"introduction": "", "conversation": "", "question": ""}
        
        for line in response.split('\n'):
            line = line.strip()
            if not line:
                continue
                
            if line == "---":
                if any(current_question.values()):
                    questions.append(current_question.copy())
                    current_question = {"introduction": "", "conversation": "", "question": ""}
                continue
                
            if line.startswith("Question #"):
                current_section = "question_number"
            elif line.startswith("Introduction:"):
                current_section = "introduction"
            elif line.startswith("Conversation:"):
                current_section = "conversation"
            elif line.startswith("Question:"):
                current_section = "question"
            elif current_section:
                current_question[current_section] += line + "\n"

        # Add the last question if it exists
        if any(current_question.values()):
            questions.append(current_question)

        # Clean up the text
        for q in questions:
            for key in q:
                q[key] = q[key].strip()

        return questions

    def save_structured_data(self, questions: List[Dict], filename: str) -> bool:
        """
        Save structured question data to text file
        
        Args:
            questions (List[Dict]): The extracted questions
            filename (str): Output filename
            
        Returns:
            bool: True if successful, False otherwise
        """
        try:
            # Create structured_data directory if it doesn't exist
            os.makedirs("./structured_data", exist_ok=True)
            
            with open(f"./structured_data/{filename}.txt", 'w', encoding='utf-8') as f:
                for i, question in enumerate(questions, 1):
                    # Determine question number (例 for first question, actual number for others)
                    question_num = "例" if i == 1 else str(i-1)
                    f.write(f"Question #{question_num}:\n\n")
                    f.write("Introduction:\n")
                    f.write(question['introduction'] + "\n\n")
                    f.write("Conversation:\n")
                    f.write(question['conversation'] + "\n\n")
                    f.write("Question:\n")
                    f.write(question['question'] + "\n")
                    f.write("---\n\n")
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
    questions = extractor.extract_questions(transcript)
    
    if questions:
        if extractor.save_structured_data(questions, output_filename):
            print(f"Structured data saved successfully to {output_filename}.txt")
            return True
        else:
            print("Failed to save structured data")
    else:
        print("Failed to extract questions from transcript")
    
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
                    print("\nExtracted Questions:")
                    print(f.read())
            except Exception as e:
                print(f"Error reading results: {str(e)}")
    except Exception as e:
        print(f"Error reading transcript file: {str(e)}")
    