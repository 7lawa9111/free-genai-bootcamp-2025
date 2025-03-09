from comps import MicroService, ServiceOrchestrator
import os
from comps import ServiceType
from fastapi import FastAPI
from pydantic import BaseModel
import aiohttp
import asyncio

EMBEDDING_SERVICE_HOST_IP = os.getenv("EMBEDDING_SERVICE_HOST_IP", "0.0.0.0")
EMBEDDING_SERVICE_PORT = os.getenv("EMBEDDING_SERVICE_PORT", 8008)
LLM_SERVICE_HOST_IP = os.getenv("LLM_SERVICE_HOST_IP", "0.0.0.0")
LLM_SERVICE_PORT = os.getenv("LLM_SERVICE_PORT", 8008)

class TextInput(BaseModel):
    text: str

class ExampleService:
    def __init__(self, host="0.0.0.0", port=8000):
        self.host = host
        self.port = port
        self.megaservice = ServiceOrchestrator()
        self.app = FastAPI()
        
        # Register endpoints
        self.app.post("/process")(self.process_text)

    def add_remote_service(self):
        embedding = MicroService(
            name="embedding",
            host=EMBEDDING_SERVICE_HOST_IP,
            port=EMBEDDING_SERVICE_PORT,
            endpoint="/api/embeddings",
            use_remote_service=True,
            service_type=ServiceType.EMBEDDING,
        )
        llm = MicroService(
            name="llm",
            host=LLM_SERVICE_HOST_IP,
            port=LLM_SERVICE_PORT,
            endpoint="/api/chat",
            use_remote_service=True,
            service_type=ServiceType.LLM,
        )
        self.megaservice.add(embedding).add(llm)
        self.megaservice.flow_to(embedding, llm)

    async def process_text(self, input_data: TextInput):
        """Process text through the embedding and LLM pipeline"""
        try:
            # Format request for Ollama
            request_data = {
                "model": "orca-mini",  # Much smaller model
                "prompt": input_data.text,
                "stream": False  # Important: don't stream the response
            }
            
            async with aiohttp.ClientSession() as session:
                # Call Ollama directly for now
                async with session.post(
                    f"http://{LLM_SERVICE_HOST_IP}:{LLM_SERVICE_PORT}/api/generate",
                    json=request_data
                ) as response:
                    if response.status == 200:
                        result = await response.json()
                        return {"response": result.get("response", "")}
                    else:
                        error_text = await response.text()
                        return {"error": f"Ollama error: {error_text}"}

        except aiohttp.ClientConnectorError as e:
            return {
                "error": f"Connection error: {str(e)}. Please ensure Ollama is running on port 8008."
            }
        except Exception as e:
            return {"error": f"Unexpected error: {str(e)}"}

    def run(self):
        import uvicorn
        uvicorn.run(self.app, host=self.host, port=self.port)

if __name__ == "__main__":
    # Create an instance of ExampleService
    service = ExampleService(host="0.0.0.0", port=8000)
    
    # Add the remote services (embedding and llm)
    service.add_remote_service()
    
    print("Service initialized with:")
    print(f"- Embedding service at {EMBEDDING_SERVICE_HOST_IP}:{EMBEDDING_SERVICE_PORT}")
    print(f"- LLM service at {LLM_SERVICE_HOST_IP}:{LLM_SERVICE_PORT}")
    print("Starting server at http://0.0.0.0:8000")
    service.run()