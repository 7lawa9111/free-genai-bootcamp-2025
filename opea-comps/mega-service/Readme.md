## How to run the LLM Service

We are using Ollama which is being delivered via docker compose.

The service uses port 8008 by default for the LLM endpoint. You can change this by setting the LLM_ENDPOINT_PORT environment variable:

```sh
LLM_ENDPOINT_PORT=8008 docker compose up
```

When you start Ollama, you'll need to download a model. We're using orca-mini as it requires less memory:

### Download (Pull) the model

```sh
docker exec ollama-server ollama pull orca-mini
```

## How to Run the Mega Service Example

1. Install the requirements:
```sh
pip install -r requirements.txt
```

2. Run the service:
```sh
python app.py
```

## Testing the App

Install jq to pretty-print JSON output:
```sh
sudo apt-get install jq  # For Ubuntu/Debian
brew install jq         # For macOS
```

Send a test request:
```sh
curl -X POST http://localhost:8000/process \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Hello, how are you?"
  }' | jq '.'
```

## Environment Variables

The service can be configured using:
- `EMBEDDING_SERVICE_HOST_IP`: Default "0.0.0.0"
- `EMBEDDING_SERVICE_PORT`: Default 8008
- `LLM_SERVICE_HOST_IP`: Default "0.0.0.0"
- `LLM_SERVICE_PORT`: Default 8008
- `LLM_ENDPOINT_PORT`: Port for Ollama in docker-compose (default 8008)

## How to access the Jaeger UI

When you run docker compose it should start up Jaeger:

```sh
http://localhost:16686/
```
