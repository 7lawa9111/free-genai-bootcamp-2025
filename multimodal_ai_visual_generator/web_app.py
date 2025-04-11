import os
import sys
import io
import base64
from pathlib import Path
import threading
import time
import queue
import json

import cv2
import numpy as np
import openvino as ov
import openvino_genai as ov_genai
from PIL import Image
from flask import Flask, render_template, request, jsonify, send_file
from torchvision.transforms import Compose

# Import the necessary functions from the original app
from depth_anything_v2_util_transform import Resize, NormalizeImage, PrepareForNet
from superres import superres_load
from vad_whisper_workers import VADWorker, WhisperWorker

# Set the Qt platform to offscreen for headless environments
os.environ["QT_QPA_PLATFORM"] = "offscreen"

app = Flask(__name__)

# Global variables to store the models and queues
models = {}
queues = {}
workers = {}
results = {}

# Initialize the models
def initialize_models():
    print("Initializing models...")
    core = ov.Core()

    llm_device = "GPU" if "GPU" in core.available_devices else "CPU"
    sd_device = "GPU" if "GPU" in core.available_devices else "CPU"
    whisper_device = 'CPU'
    super_res_device = "GPU" if "GPU" in core.available_devices else "CPU"
    depth_anything_device = "GPU" if "GPU" in core.available_devices else "CPU"

    # Create the 'results' folder if it doesn't exist
    Path("results").mkdir(exist_ok=True)

    # Initialize LLM pipeline
    print("Creating an llm pipeline to run on ", llm_device)
    llm_model_path = r"./models/tinyllama-1.1b-chat-v1.0/INT4_compressed_weights"
    
    if llm_device == 'NPU':
        pipeline_config = {"MAX_PROMPT_LEN": 1536}
        llm_pipe = ov_genai.LLMPipeline(llm_model_path, llm_device, pipeline_config)
    else:
        llm_pipe = ov_genai.LLMPipeline(llm_model_path, llm_device)
    
    models["llm"] = llm_pipe
    print("Done creating our llm..")

    # Initialize Stable Diffusion pipeline
    print("Creating a stable diffusion pipeline to run on ", sd_device)
    sd_pipe = ov_genai.Text2ImagePipeline(r"models/LCM_Dreamshaper_v7/FP16", sd_device)
    models["sd"] = sd_pipe
    print("Done creating the stable diffusion pipeline...")

    # Initialize Whisper model
    models["whisper_device"] = whisper_device

    # Initialize Super Resolution model
    print("Initializing Super Res Model to run on ", super_res_device)
    model_path_sr = Path(f"models/single-image-super-resolution-1033.xml")
    super_res_compiled_model, super_res_upsample_factor = superres_load(model_path_sr, super_res_device, h_custom=432, w_custom=768)
    models["super_res_compiled_model"] = super_res_compiled_model
    models["super_res_upsample_factor"] = super_res_upsample_factor
    print("Initializing Super Res Model done...")

    # Initialize Depth Anything v2 model
    print("Initializing Depth Anything v2 model to run on ", depth_anything_device)
    OV_DEPTH_ANYTHING_PATH = Path(f"models/depth_anything_v2_vits.xml")
    depth_compiled_model = core.compile_model(OV_DEPTH_ANYTHING_PATH, device_name=depth_anything_device)
    models["depth_compiled_model"] = depth_compiled_model
    print("Initializing Depth Anything v2 done...")

    print("All models initialized successfully!")

# Helper functions from the original app
def depth_map_parallax(compiled_model, image):
    image.save("results/original_image.png")
    image = np.array(image)

    h, w = image.shape[:2]

    transform = Compose(
        [
            Resize(
                width=770,
                height=434,
                resize_target=False,
                ensure_multiple_of=14,
                resize_method="lower_bound",
                image_interpolation_method=cv2.INTER_CUBIC,
            ),
            NormalizeImage(mean=[0.485, 0.456, 0.406], std=[0.229, 0.224, 0.225]),
            PrepareForNet(),
        ]
    )
    def predict_depth(model, image):
        return model(image)[0]

    image = cv2.cvtColor(image, cv2.COLOR_BGR2RGB) / 255.0
    image = transform({"image": image})["image"]
    image = np.expand_dims(image, 0)

    depth = predict_depth(compiled_model, image)
    depth = cv2.resize(depth[0], (w, h), interpolation=cv2.INTER_LINEAR)

    depth = (depth - depth.min()) / (depth.max() - depth.min()) * 255.0
    depth = depth.astype(np.uint8)
    colored_depth = cv2.applyColorMap(depth, cv2.COLORMAP_INFERNO)[:, :, ::-1]

    im = Image.fromarray(colored_depth)
    im.save("results/depth_map.png")
    return im

def convert_result_to_image(result) -> np.ndarray:
    result = result.squeeze(0).transpose(1, 2, 0)
    result *= 255
    result[result < 0] = 0
    result[result > 255] = 255
    result = result.astype(np.uint8)
    return Image.fromarray(result)

def run_sr(compiled_model, upsample_factor, img):
    input_image_original = np.expand_dims(img.transpose(2, 0, 1), axis=0)
    bicubic_image = cv2.resize(
    src=img, dsize=(768*upsample_factor, 432*upsample_factor), interpolation=cv2.INTER_CUBIC)
    input_image_bicubic = np.expand_dims(bicubic_image.transpose(2, 0, 1), axis=0)

    original_image_key, bicubic_image_key = compiled_model.inputs
    output_key = compiled_model.output(0)

    result = compiled_model(
    {
        original_image_key.any_name: input_image_original,
        bicubic_image_key.any_name: input_image_bicubic,
    }
    )[output_key]

    result_image = convert_result_to_image(result)
    return result_image

def generate_image(prompt):
    image_tensor = models["sd"].generate(
    prompt,
    width=768,
    height=432,
    num_inference_steps=5,
    num_images_per_prompt=1)

    sd_output = Image.fromarray(image_tensor.data[0])
    
    # Apply super resolution
    sr_out = run_sr(models["super_res_compiled_model"], models["super_res_upsample_factor"], np.array(sd_output))
    
    # Generate depth map
    depth_map = depth_map_parallax(models["depth_compiled_model"], sr_out)
    
    # Save the images
    sr_out.save("results/enhanced_image.png")
    depth_map.save("results/depth_map.png")
    
    return sr_out, depth_map

# LLM system message
LLM_SYSTEM_MESSAGE_START = """
You are a specialized helper bot designed to process live transcripts from a demo called "AI Adventure Game", which showcases a tabletop adventure game with live illustrations generated by a text-to-image model.

Your role is to act as a filter:

Detect descriptions of game scenes from the transcript that require illustration.
Output a detailed SD Prompt for these scenes.
When you detect a scene for the game, output it as:

SD Prompt: <a detailed prompt for illustration>

Guidelines:
Focus only on game scenes: Ignore meta-comments, explanations about the demo, or incomplete thoughts.
Contextual Awareness: Maintain and apply story context, such as the location, atmosphere, and objects, when crafting prompts. Update this context only when a new scene is explicitly described.
No Players in Prompts: Do not include references to "the player," "the players,"  "the party", or any specific characters in the SD Prompt. Focus solely on the environment and atmosphere.
Prioritize Clarity: If unsure whether the presenter is describing a scene, return: 'None'. Avoid making assumptions about incomplete descriptions.
Enhance Visuals: Add vivid and descriptive details to SD Prompts, such as lighting, mood, style, or texture, when appropriate, but stay faithful to the transcript.
Examples:
Example 1:
Input: "Let me explain how we are using AI for these illustrations." Output: 'None'

Example 2:
Input: "The party is standing at the gates of a large castle." Output: SD Prompt: "A massive medieval castle gate with towering stone walls, surrounded by mist and faintly glowing lanterns at dusk."

Example 3:
Context: "The party is at the gates of a large castle." Input: "The party then encounters a huge dragon." Output: SD Prompt: "A massive dragon with gleaming scales, standing before the misty gates of a towering medieval castle, lit by glowing lanterns under a dim sky."

Example 4:
Input: "And now the players roll for initiative." Output: 'None'

The presenter of the demo is aware of your presence and role, and will sometimes refer to you as the 'LLM', the 'agent', etc. Occasionally he will point out your roles and read back the SD prompts that you generate. When you detect this, return 'None'.

The SD prompts should be no longer than 25 words.

Only output SD prompts it is detected that there is big difference in location as compared with the last SD prompt that you gave.

Example 1:

Input 0: "The party is standing at the gates of a large castle." Output 0: SD Prompt: "A massive medieval castle gate with towering stone walls, surrounded by mist and faintly glowing lanterns at dusk."
Input 1: "A character is still at the gates of the castle." Output 1: 'None'
"""

LLM_SYSTEM_MESSAGE_END = """
Additional hints and reminders:
* You are a filter, not a chatbot. Only provide SD Prompts or 'None.'
* No Extra Notes: Do not include explanations, comments, or any text beyond the required SD Prompt or 'None.'
* Validate Completeness: A description of a scene often involves locations, objects, or atmosphere and is unlikely to be inferred from just verbs or generic phrases.
* If it seems that the transcription of the presenter is simply reading a previous SD prompt that you generated, return 'None'
* The SD prompts should be no longer than 25 words.
* Do not provide SD prompts for what seem like incomplete thoughts. Return 'None' in this case.
* Use the given theme of the game to help you decide whether or not the given bits of transcript are describing a new scene, or not.
* Do not try to actually illustrate the characters themselves, only details of their environmental surroundings & atmosphere.
* The SD prompts should be no longer than 25 words.
* Only output SD prompts it is detected that there is big difference in location as compared with the last SD prompt that you gave. If it seems like the location is the same, just return 'None'
"""

# Worker thread for processing audio and generating images
class WorkerThread(threading.Thread):
    def __init__(self, session_id, theme="Medieval Fantasy Adventure"):
        super().__init__()
        self.session_id = session_id
        self.theme = theme
        self.running = True
        self.queue = queue.Queue()
        self.results = {"status": "idle", "caption": "", "image_path": "", "depth_map_path": ""}
        results[session_id] = self.results
        
    def run(self):
        llm_tokenizer = models["llm"].get_tokenizer()
        
        # Assemble the system message
        system_message = LLM_SYSTEM_MESSAGE_START
        system_message += "\nThe presenter is giving a hint that the theme of their game is: " + self.theme
        system_message += "\nYou should use this hint to guide your decision about whether the presenter is describing a scene from the game, or not, and also to generate adequate SD Prompts."
        system_message += "\n" + LLM_SYSTEM_MESSAGE_END
        
        generate_config = ov_genai.GenerationConfig()
        generate_config.temperature = 0.7
        generate_config.top_p = 0.95
        generate_config.max_length = 2048
        
        meaningful_message_pairs = []
        
        while self.running:
            try:
                # Wait for a sentence from the queue
                self.results["status"] = "listening"
                results[self.session_id] = self.results
                
                result = self.queue.get(timeout=1)
                
                self.results["status"] = "processing"
                results[self.session_id] = self.results
                
                chat_history = [{"role": "system", "content": system_message}]
                
                # Only keep the latest 2 meaningful message pairs
                meaningful_message_pairs = meaningful_message_pairs[-2:]
                
                for meaningful_pair in meaningful_message_pairs:
                    user_message = meaningful_pair[0]
                    assistant_response = meaningful_pair[1]
                    
                    chat_history.append({"role": "user", "content": user_message["content"]})
                    chat_history.append({"role": "assistant", "content": assistant_response["content"]})
                
                chat_history.append({"role": "user", "content": result})
                formatted_prompt = llm_tokenizer.apply_chat_template(history=chat_history, add_generation_prompt=True)
                
                self.results["status"] = "processing..."
                results[self.session_id] = self.results
                
                print("Running LLM for session", self.session_id)
                llm_result = models["llm"].generate(inputs=formatted_prompt, generation_config=generate_config)
                
                search_string = "SD Prompt:"
                
                # Check if the LLM generated an SD prompt
                if search_string in llm_result and 'None' not in llm_result:
                    # Extract the prompt
                    start_index = llm_result.find(search_string)
                    prompt = llm_result[start_index + len(search_string):].strip()
                    
                    self.results["caption"] = prompt
                    self.results["status"] = "illustrating..."
                    results[self.session_id] = self.results
                    
                    print("Generating image for session", self.session_id)
                    # Generate the image
                    enhanced_image, depth_map = generate_image(prompt)
                    
                    # Update the results
                    self.results["image_path"] = "results/enhanced_image.png"
                    self.results["depth_map_path"] = "results/depth_map.png"
                    self.results["status"] = "idle"
                    results[self.session_id] = self.results
                    
                    # Add to meaningful message pairs
                    meaningful_message_pairs.append(
                        ({"role": "user", "content": result},
                         {"role": "assistant", "content": llm_result},)
                    )
                
            except queue.Empty:
                continue  # Queue is empty, just wait
            
        self.results["status"] = "idle"
        results[self.session_id] = self.results
    
    def stop(self):
        self.running = False
        self.join()

# Flask routes
@app.route('/')
def index():
    return render_template('index.html')

@app.route('/api/start_session', methods=['POST'])
def start_session():
    data = request.json
    theme = data.get('theme', 'Medieval Fantasy Adventure')
    session_id = str(int(time.time() * 1000))  # Generate a unique session ID
    
    # Create a new worker thread for this session
    worker = WorkerThread(session_id, theme)
    workers[session_id] = worker
    worker.start()
    
    return jsonify({"session_id": session_id, "status": "started"})

@app.route('/api/stop_session', methods=['POST'])
def stop_session():
    data = request.json
    session_id = data.get('session_id')
    
    if session_id in workers:
        workers[session_id].stop()
        del workers[session_id]
        return jsonify({"status": "stopped"})
    
    return jsonify({"status": "not_found"})

@app.route('/api/send_text', methods=['POST'])
def send_text():
    data = request.json
    session_id = data.get('session_id')
    text = data.get('text')
    
    if session_id in workers:
        workers[session_id].queue.put(text)
        return jsonify({"status": "sent"})
    
    return jsonify({"status": "not_found"})

@app.route('/api/get_status', methods=['GET'])
def get_status():
    session_id = request.args.get('session_id')
    
    if session_id in results:
        return jsonify(results[session_id])
    
    return jsonify({"status": "not_found"})

@app.route('/api/get_image/<path:filename>')
def get_image(filename):
    return send_file(filename)

# Create templates directory if it doesn't exist
Path("templates").mkdir(exist_ok=True)

# Create the index.html template
@app.route('/templates/index.html')
def serve_index():
    return """
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>AI Adventure Game</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        h1 {
            color: #333;
            text-align: center;
        }
        .container {
            display: flex;
            flex-direction: column;
            gap: 20px;
        }
        .control-panel {
            background-color: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .image-panel {
            display: flex;
            gap: 20px;
            justify-content: center;
        }
        .image-container {
            background-color: white;
            padding: 10px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            text-align: center;
        }
        .image-container img {
            max-width: 100%;
            height: auto;
            border: 1px solid #ddd;
        }
        .caption {
            margin-top: 10px;
            font-style: italic;
            color: #666;
        }
        .status {
            margin-top: 10px;
            font-weight: bold;
            color: #333;
        }
        button {
            background-color: #4CAF50;
            color: white;
            border: none;
            padding: 10px 15px;
            text-align: center;
            text-decoration: none;
            display: inline-block;
            font-size: 16px;
            margin: 4px 2px;
            cursor: pointer;
            border-radius: 4px;
        }
        button:disabled {
            background-color: #cccccc;
            cursor: not-allowed;
        }
        input[type="text"] {
            padding: 10px;
            margin: 5px 0;
            border: 1px solid #ddd;
            border-radius: 4px;
            width: 100%;
            box-sizing: border-box;
        }
        .progress-bar {
            width: 100%;
            height: 20px;
            background-color: #f0f0f0;
            border-radius: 10px;
            overflow: hidden;
            margin-top: 10px;
        }
        .progress {
            height: 100%;
            background-color: #4CAF50;
            width: 0%;
            transition: width 0.3s;
        }
    </style>
</head>
<body>
    <h1>AI Adventure Game</h1>
    <div class="container">
        <div class="control-panel">
            <h2>Control Panel</h2>
            <div id="session-controls">
                <input type="text" id="theme-input" placeholder="Enter game theme (e.g., Medieval Fantasy Adventure)" value="Medieval Fantasy Adventure">
                <button id="start-button">Start Session</button>
                <button id="stop-button" disabled>Stop Session</button>
            </div>
            <div id="text-input-container" style="display: none;">
                <input type="text" id="text-input" placeholder="Enter text description...">
                <button id="send-button">Send</button>
            </div>
            <div class="status" id="status">Status: Idle</div>
            <div class="progress-bar">
                <div class="progress" id="progress-bar"></div>
            </div>
        </div>
        
        <div class="image-panel">
            <div class="image-container">
                <h3>Enhanced Image</h3>
                <img id="enhanced-image" src="" alt="Enhanced Image" style="display: none;">
                <div class="caption" id="caption"></div>
            </div>
            <div class="image-container">
                <h3>Depth Map</h3>
                <img id="depth-map" src="" alt="Depth Map" style="display: none;">
            </div>
        </div>
    </div>

    <script>
        let sessionId = null;
        let statusInterval = null;
        
        document.getElementById('start-button').addEventListener('click', function() {
            const theme = document.getElementById('theme-input').value;
            
            fetch('/api/start_session', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ theme: theme }),
            })
            .then(response => response.json())
            .then(data => {
                sessionId = data.session_id;
                document.getElementById('start-button').disabled = true;
                document.getElementById('stop-button').disabled = false;
                document.getElementById('text-input-container').style.display = 'block';
                document.getElementById('status').textContent = 'Status: Started';
                
                // Start polling for status updates
                startStatusPolling();
            })
            .catch(error => {
                console.error('Error:', error);
                document.getElementById('status').textContent = 'Status: Error starting session';
            });
        });
        
        document.getElementById('stop-button').addEventListener('click', function() {
            if (sessionId) {
                fetch('/api/stop_session', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ session_id: sessionId }),
                })
                .then(response => response.json())
                .then(data => {
                    document.getElementById('start-button').disabled = false;
                    document.getElementById('stop-button').disabled = true;
                    document.getElementById('text-input-container').style.display = 'none';
                    document.getElementById('status').textContent = 'Status: Stopped';
                    
                    // Stop polling for status updates
                    stopStatusPolling();
                })
                .catch(error => {
                    console.error('Error:', error);
                    document.getElementById('status').textContent = 'Status: Error stopping session';
                });
            }
        });
        
        document.getElementById('send-button').addEventListener('click', function() {
            const text = document.getElementById('text-input').value;
            if (text && sessionId) {
                fetch('/api/send_text', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ session_id: sessionId, text: text }),
                })
                .then(response => response.json())
                .then(data => {
                    document.getElementById('text-input').value = '';
                })
                .catch(error => {
                    console.error('Error:', error);
                });
            }
        });
        
        function startStatusPolling() {
            if (statusInterval) {
                clearInterval(statusInterval);
            }
            
            statusInterval = setInterval(function() {
                if (sessionId) {
                    fetch(`/api/get_status?session_id=${sessionId}`)
                        .then(response => response.json())
                        .then(data => {
                            document.getElementById('status').textContent = `Status: ${data.status}`;
                            
                            // Update progress bar based on status
                            if (data.status === 'listening') {
                                document.getElementById('progress-bar').style.width = '25%';
                            } else if (data.status === 'processing') {
                                document.getElementById('progress-bar').style.width = '50%';
                            } else if (data.status === 'illustrating...') {
                                document.getElementById('progress-bar').style.width = '75%';
                            } else if (data.status === 'idle') {
                                document.getElementById('progress-bar').style.width = '100%';
                            }
                            
                            // Update images if available
                            if (data.image_path) {
                                const enhancedImage = document.getElementById('enhanced-image');
                                enhancedImage.src = `/api/get_image/${data.image_path}?t=${new Date().getTime()}`;
                                enhancedImage.style.display = 'block';
                            }
                            
                            if (data.depth_map_path) {
                                const depthMap = document.getElementById('depth-map');
                                depthMap.src = `/api/get_image/${data.depth_map_path}?t=${new Date().getTime()}`;
                                depthMap.style.display = 'block';
                            }
                            
                            // Update caption if available
                            if (data.caption) {
                                document.getElementById('caption').textContent = data.caption;
                            }
                        })
                        .catch(error => {
                            console.error('Error:', error);
                        });
                }
            }, 1000);
        }
        
        function stopStatusPolling() {
            if (statusInterval) {
                clearInterval(statusInterval);
                statusInterval = null;
            }
        }
    </script>
</body>
</html>
    """

# Initialize the models when the app starts
initialize_models()

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8080, debug=True)