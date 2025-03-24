from typing import List, Dict, Optional
import chromadb
from chromadb.utils import embedding_functions
import boto3
import json
import uuid
import os
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

class BedrockEmbedding(embedding_functions.EmbeddingFunction):
    def __init__(self):
        # Initialize Bedrock client using environment variables
        self.bedrock = boto3.client(
            service_name='bedrock-runtime',
            region_name='us-east-2',
            aws_access_key_id=os.getenv('AWS_ACCESS_KEY_ID'),
            aws_secret_access_key=os.getenv('AWS_SECRET_ACCESS_KEY')
        )
        self.model_id = 'amazon.titan-embed-text-v2:0'

    def __call__(self, texts: List[str]) -> List[List[float]]:
        """Generate embeddings for the input texts using Amazon Bedrock"""
        embeddings = []
        
        for text in texts:
            try:
                body = json.dumps({
                    "inputText": text
                })
                
                response = self.bedrock.invoke_model(
                    modelId=self.model_id,
                    body=body
                )
                
                response_body = json.loads(response['body'].read())
                embedding = response_body['embedding']
                embeddings.append(embedding)
                
            except Exception as e:
                print(f"Error generating embedding: {str(e)}")
                # Return a zero vector of the same size as Titan embeddings
                embeddings.append([0.0] * 1536)  # Titan embedding size
                
        return embeddings

class JLPTQuestionStore:
    def __init__(self, persist_directory: str = "./vector_store"):
        """
        Initialize the vector store for JLPT questions
        """
        # Initialize ChromaDB client with persistence
        self.client = chromadb.PersistentClient(path=persist_directory)
        
        # Create or get collection for questions with Bedrock embedding
        try:
            # Try to get existing collection
            self.questions_collection = self.client.get_collection(
                name="jlpt_questions"
            )
        except:
            # Create new collection with Bedrock embedding
            self.questions_collection = self.client.create_collection(
                name="jlpt_questions",
                embedding_function=BedrockEmbedding(),
                metadata={"hnsw:space": "cosine"}  # Use cosine similarity
            )

    def add_question(self, section_num: int, question_num: str, question_data: Dict):
        """Add a question to the vector store"""
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
        """Find similar questions using vector similarity"""
        results = self.questions_collection.query(
            query_texts=[query],
            n_results=n_results
        )
        return results['metadatas'][0]  # Return metadata of similar questions

    def get_questions_by_section(self, section_num: int) -> List[Dict]:
        """Get all questions from a specific section"""
        results = self.questions_collection.get(
            where={"section_num": section_num}
        )
        return results['metadatas'] 