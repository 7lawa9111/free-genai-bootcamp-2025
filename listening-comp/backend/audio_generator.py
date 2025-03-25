import boto3
import json
import os
from typing import List, Dict
import tempfile
import subprocess
from dotenv import load_dotenv

load_dotenv()

class AudioGenerator:
    def __init__(self):
        # Initialize Bedrock client for Nova
        self.bedrock = boto3.client(
            service_name='bedrock-runtime',
            region_name='us-east-2',
            aws_access_key_id=os.getenv('AWS_ACCESS_KEY_ID'),
            aws_secret_access_key=os.getenv('AWS_SECRET_ACCESS_KEY')
        )
        
        # Initialize Polly client
        self.polly = boto3.client(
            service_name='polly',
            region_name='us-east-1',
            aws_access_key_id=os.getenv('AWS_ACCESS_KEY_ID'),
            aws_secret_access_key=os.getenv('AWS_SECRET_ACCESS_KEY')
        )
        
        # Define available voices using standard engine and Japanese voices
        self.voices = {
            'announcer': {'id': 'Mizuki', 'engine': 'standard'},
            'male': [
                {'id': 'Takumi', 'engine': 'standard'}
            ],
            'female': [
                {'id': 'Mizuki', 'engine': 'standard'}
            ]
        }
        
        # Verify available voices
        try:
            available_voices = self.polly.describe_voices(LanguageCode='ja-JP')
            print("Available Japanese voices:", [v['Id'] for v in available_voices['Voices']])
        except Exception as e:
            print(f"Warning: Could not fetch available voices: {e}")
        
    def format_dialogue(self, question: Dict) -> Dict[str, str]:
        """Use Nova to format the dialogue with speaker annotations"""
        prompt = f"""
        Format this JLPT listening question for audio generation.
        You must identify speakers and their genders (male/female only) and split the content into parts.

        Here's the content to format:
        Introduction: {question['introduction']}
        Conversation: {question['conversation']}
        Question: {question['question']}

        Return a JSON object with exactly this structure:
        {{
            "introduction_text": "text for announcer to read",
            "conversation_parts": [
                {{"speaker_gender": "male or female", "text": "what this speaker says"}},
                {{"speaker_gender": "male or female", "text": "what next speaker says"}},
                ...
            ],
            "question_text": "question for announcer to read"
        }}

        Ensure each speaker_gender is exactly either "male" or "female".
        """
        
        try:
            # Using Nova model for formatting
            response = self.bedrock.invoke_model(
                modelId='us.amazon.nova-micro-v1:0',
                body=json.dumps({
                    "messages": [
                        {
                            "role": "user",
                            "content": [{"text": prompt}]
                        }
                    ],
                    "inferenceConfig": {
                        "temperature": 0,
                        "topP": 1,
                        "maxTokens": 1000,
                        "stopSequences": []
                    }
                })
            )
            
            response_body = json.loads(response.get('body').read())
            content = response_body['output']['message']['content'][0]['text']
            
            # Extract JSON from the response
            start_idx = content.find('{')
            end_idx = content.rfind('}')
            if start_idx == -1 or end_idx == -1:
                raise ValueError("No JSON found in response")
            
            json_str = content[start_idx:end_idx + 1]
            result = json.loads(json_str)
            
            # Validate the response format
            required_keys = ['introduction_text', 'conversation_parts', 'question_text']
            if not all(key in result for key in required_keys):
                raise ValueError(f"Missing required keys in response: {result}")
                
            # Validate conversation parts
            if not isinstance(result['conversation_parts'], list):
                raise ValueError("conversation_parts must be a list")
                
            for part in result['conversation_parts']:
                if not isinstance(part, dict):
                    raise ValueError("Each conversation part must be a dictionary")
                if 'speaker_gender' not in part or 'text' not in part:
                    raise ValueError(f"Missing required keys in conversation part: {part}")
                if part['speaker_gender'] not in ['male', 'female']:
                    part['speaker_gender'] = 'male'  # Default to male if invalid
            
            return result
            
        except Exception as e:
            print(f"Error formatting dialogue: {e}")
            print(f"Response content: {content if 'content' in locals() else 'No content'}")
            return None

    def generate_audio_part(self, text: str, voice_config: dict) -> str:
        """Generate audio for a single part using Amazon Polly"""
        try:
            response = self.polly.synthesize_speech(
                Text=text,
                OutputFormat='mp3',
                VoiceId=voice_config['id'],
                LanguageCode='ja-JP',
                Engine=voice_config['engine']
            )
            
            # Save to temporary file
            with tempfile.NamedTemporaryFile(suffix='.mp3', delete=False) as f:
                f.write(response['AudioStream'].read())
                return f.name
                
        except Exception as e:
            print(f"Error generating audio: {e}")
            return None

    def combine_audio_files(self, audio_files: List[str], output_file: str):
        """Combine multiple audio files using ffmpeg"""
        try:
            # Create file list for ffmpeg with absolute paths
            with tempfile.NamedTemporaryFile('w', suffix='.txt', delete=False) as f:
                for audio_file in audio_files:
                    f.write(f"file '{os.path.abspath(audio_file)}'\n")
                file_list = f.name

            # Ensure output path is absolute
            output_path = os.path.abspath(output_file)
            
            # Combine audio files
            subprocess.run([
                'ffmpeg', '-f', 'concat', '-safe', '0',
                '-i', file_list,
                '-c', 'copy',
                output_path
            ], check=True)
            
            # Cleanup temporary files
            os.unlink(file_list)
            for file in audio_files:
                os.unlink(file)
                
        except Exception as e:
            print(f"Error combining audio: {e}")
            return None

    def generate_question_audio(self, question: Dict) -> str:
        """Generate complete audio for a question"""
        try:
            # Create output directory if it doesn't exist
            output_dir = "generated_audio"
            os.makedirs(output_dir, exist_ok=True)
            
            # Format dialogue
            formatted = self.format_dialogue(question)
            if not formatted:
                raise ValueError("Failed to format dialogue")
            
            audio_parts = []
            
            # Generate introduction audio
            intro_audio = self.generate_audio_part(
                formatted['introduction_text'],
                self.voices['announcer']
            )
            if not intro_audio:
                raise ValueError("Failed to generate introduction audio")
            audio_parts.append(intro_audio)
            
            # Generate conversation audio
            used_voices = {'male': 0, 'female': 0}
            for part in formatted['conversation_parts']:
                gender = part['speaker_gender'].lower()
                if gender not in self.voices:
                    print(f"Warning: Unknown gender {gender}, using male voice")
                    gender = 'male'
                    
                # Select voice and track usage
                voice_config = self.voices[gender][used_voices[gender] % len(self.voices[gender])]
                used_voices[gender] += 1
                
                part_audio = self.generate_audio_part(part['text'], voice_config)
                if not part_audio:
                    raise ValueError(f"Failed to generate audio for part: {part['text'][:50]}...")
                audio_parts.append(part_audio)
            
            # Generate question audio
            question_audio = self.generate_audio_part(
                formatted['question_text'],
                self.voices['announcer']
            )
            if not question_audio:
                raise ValueError("Failed to generate question audio")
            audio_parts.append(question_audio)
            
            # Combine all parts
            output_file = os.path.join(output_dir, f"question_audio_{hash(question['introduction'])}.mp3")
            self.combine_audio_files(audio_parts, output_file)
            
            if not os.path.exists(output_file):
                raise ValueError("Failed to create final audio file")
            
            return output_file
            
        except Exception as e:
            print(f"Error generating question audio: {e}")
            # Cleanup any temporary files
            for file in audio_parts:
                try:
                    os.unlink(file)
                except:
                    pass
            return None 