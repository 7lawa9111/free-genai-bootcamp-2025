# Create BedrockChat
# bedrock_chat.py
import boto3
import streamlit as st
from typing import Optional, Dict, Any
import os
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

# Model ID
MODEL_ID = "us.amazon.nova-micro-v1:0"



class BedrockChat:
    def __init__(self):
        # Initialize Bedrock client using environment variables
        self.bedrock = boto3.client(
            service_name='bedrock-runtime',
            region_name='us-east-2',
            aws_access_key_id=os.getenv('AWS_ACCESS_KEY_ID'),
            aws_secret_access_key=os.getenv('AWS_SECRET_ACCESS_KEY')
        )
        self.model_id = 'us.amazon.nova-micro-v1:0'

    def generate_response(self, prompt: str) -> Optional[str]:
        """
        Generate a response using Amazon Bedrock
        
        Args:
            prompt (str): The user's input prompt
            
        Returns:
            Optional[str]: The generated response or None if failed
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
            print(f"Error generating response: {str(e)}")
            return None


if __name__ == "__main__":
    chat = BedrockChat()
    while True:
        user_input = input("You: ")
        if user_input.lower() == '/exit':
            break
        response = chat.generate_response(user_input)
        print("Bot:", response)
