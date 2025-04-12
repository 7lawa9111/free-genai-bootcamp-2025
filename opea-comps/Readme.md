## Running Ollama Third-Party Service

### Choosing a Model

You can get the model_id that ollama will launch from the [Ollama Library](https://ollama.com/library).

https://ollama.com/library/llama3.2

Current configuration uses: `TheBloke/Llama-2-7B-Chat-GGUF`

### Getting the Host IP

#### Linux
Get your IP address
```sh
sudo apt install net-tools
ifconfig
```

Or you can try this way `$(hostname -I | awk '{print $1}')`

#### macOS
Get your IP address using:
```sh
ifconfig | grep "inet " | grep -v 127.0.0.1 | awk '{print $2}'
```

Or alternatively:
```sh
ipconfig getifaddr en0  # For WiFi connection
ipconfig getifaddr en1  # For ethernet connection
```

### Running the Container

For macOS with specific context:
```sh
# Option 1: Export variables first (recommended)
export HOST_IP=$(ipconfig getifaddr en0)
export NO_PROXY=localhost
export LLM_ENDPOINT_PORT=8008
export LLM_MODEL_ID="TheBloke/Llama-2-7B-Chat-GGUF"
docker --context desktop-linux compose up

# Option 2: One-liner with proper quoting
HOST_IP="$(ipconfig getifaddr en0)" NO_PROXY="localhost" LLM_ENDPOINT_PORT="8008" LLM_MODEL_ID="TheBloke/Llama-2-7B-Chat-GGUF" docker --context desktop-linux compose up
```

For Linux:
```sh
HOST_IP="$(hostname -I | awk '{print $1}')" NO_PROXY="localhost" LLM_ENDPOINT_PORT="8008" LLM_MODEL_ID="TheBloke/Llama-2-7B-Chat-GGUF" docker --context desktop-linux compose up
```

### Ollama API

Once the Ollama server is running we can make API calls to the ollama API

https://github.com/ollama/ollama/blob/main/docs/api.md


## Download (Pull) a model

curl http://localhost:8008/api/pull -d '{
  "model": "TheBloke/Llama-2-7B-Chat-GGUF"
}'

## Generate a Request

curl http://localhost:8008/api/generate -d '{
  "model": "TheBloke/Llama-2-7B-Chat-GGUF",
  "prompt": "Why is the sky blue?"
}'

# Technical Uncertainty

Q Does bridge mode mean we can only accses Ollama API with another model in the docker compose?

A No, the host machine will be able to access it

Q: Which port is being mapped 8008->11434

In this case 8008 is the port that host machine will access. the other in the guest port (the port of the service inside container)

Q: If we pass the LLM_MODEL_Id to the ollama server will it download the model when on start?

It does not appear so. 
```sh
docker exec ollama-server ollama pull TheBloke/Llama-2-7B-Chat-GGUF
```
Q: Will the model be downloaded in the container? does that mean the ml model will be deleted when the container stops running?

A: The model will download into the container, and vanish when the container stop running. You need to mount a local drive and there is probably more work to be done.

Q: For LLM service which can text-generation it suggets it will only work with TGI/vLLM and all you have to do is have it running. Does TGI and vLLM have a stardarized API or is there code to detect which one is running? Do we have to really use Xeon or Guadi processor?

VLLM service was commented out due to these limitations, Implemented custom TTS service with better CPU compatibility. 