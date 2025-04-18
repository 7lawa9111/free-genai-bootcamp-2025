import io
import sys
from pathlib import Path

import cv2
import numpy as np
import openvino as ov
import openvino_genai as ov_genai
from PIL import Image
from PySide6.QtCore import Qt, QThread
from PySide6.QtCore import Signal
from PySide6.QtGui import QPixmap, QFont
from PySide6.QtWidgets import (
    QApplication, QMainWindow, QWidget, QPushButton, QVBoxLayout,
    QHBoxLayout, QGridLayout, QProgressBar, QTextEdit, QLineEdit
)
from PySide6.QtWidgets import QLabel
from torchvision.transforms import Compose

from depth_anything_v2_util_transform import Resize, NormalizeImage, PrepareForNet
from superres import superres_load
from vad_whisper_workers import VADWorker, WhisperWorker

import os
import uuid
import json
from flask import Flask, render_template, request, jsonify, send_file
from diffusers import StableDiffusionPipeline
from openvino.runtime import Core
import torch

app = Flask(__name__)

# Initialize OpenVINO
ie = Core()
device = "CPU"  # or "GPU" if available

# Initialize the model
model_id = "runwayml/stable-diffusion-v1-5"
pipe = StableDiffusionPipeline.from_pretrained(model_id, torch_dtype=torch.float16)
if torch.cuda.is_available():
    pipe = pipe.to("cuda")
else:
    pipe = pipe.to("cpu")

# Create a directory for storing generated images
UPLOAD_FOLDER = Path("static/generated")
UPLOAD_FOLDER.mkdir(parents=True, exist_ok=True)

# Store session information
sessions = {}

@app.route('/')
def index():
    return render_template('index.html')

@app.route('/start_session', methods=['GET'])
def start_session():
    session_id = str(uuid.uuid4())
    sessions[session_id] = {
        'status': 'waiting',
        'image': None,
        'error': None
    }
    return jsonify({'session_id': session_id})

@app.route('/send_text', methods=['POST'])
def send_text():
    data = request.json
    session_id = data.get('session_id')
    text = data.get('text')

    if not session_id or not text:
        return jsonify({'status': 'error', 'message': 'Missing session_id or text'}), 400

    if session_id not in sessions:
        return jsonify({'status': 'error', 'message': 'Invalid session_id'}), 400

    try:
        # Update session status
        sessions[session_id]['status'] = 'processing'
        
        # Generate image
        image = pipe(text).images[0]
        
        # Convert PIL image to bytes
        img_byte_arr = io.BytesIO()
        image.save(img_byte_arr, format='PNG')
        img_byte_arr = img_byte_arr.getvalue()
        
        # Store the image
        sessions[session_id]['image'] = img_byte_arr
        sessions[session_id]['status'] = 'completed'
        
        return jsonify({'status': 'success'})
    except Exception as e:
        sessions[session_id]['status'] = 'error'
        sessions[session_id]['error'] = str(e)
        return jsonify({'status': 'error', 'message': str(e)}), 500

@app.route('/get_status/<session_id>')
def get_status(session_id):
    if session_id not in sessions:
        return jsonify({'status': 'error', 'message': 'Invalid session_id'}), 400
    
    session = sessions[session_id]
    if session['status'] == 'error':
        return jsonify({
            'status': 'error',
            'message': session['error']
        })
    
    return jsonify({'status': session['status']})

@app.route('/get_image/<session_id>')
def get_image(session_id):
    if session_id not in sessions:
        return jsonify({'status': 'error', 'message': 'Invalid session_id'}), 400
    
    session = sessions[session_id]
    if session['status'] != 'completed' or not session['image']:
        return jsonify({'status': 'error', 'message': 'Image not ready'}), 400
    
    return send_file(
        io.BytesIO(session['image']),
        mimetype='image/png'
    )

def depth_map_parallax(compiled_model, image):
    #This function will load the OV Depth Anything model
    #and create a 3D parallax between the depth map and the input image
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

    #Have web server pick up images and serve them
    im = Image.fromarray(colored_depth)
    im.save("results/depth_map.png")
    return im

def convert_result_to_image(result) -> np.ndarray:
    """
    Convert network result of floating point numbers to image with integer
    values from 0-255. Values outside this range are clipped to 0 and 255.

    :param result: a single superresolution network result in N,C,H,W shape
    """
    result = result.squeeze(0).transpose(1, 2, 0)
    result *= 255
    result[result < 0] = 0
    result[result > 255] = 255
    result = result.astype(np.uint8)
    return Image.fromarray(result)

LLM_SYSTEM_MESSAGE_START="""
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

LLM_SYSTEM_MESSAGE_END="""

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


from queue import Empty
class WorkerThread(QThread):
    image_updated = Signal(QPixmap)
    caption_updated = Signal(str)
    progress_updated = Signal(int, str)

    primary_pixmap_updated = Signal(QPixmap)
    depth_pixmap_updated = Signal(QPixmap)

    def __init__(self, queue, app_params, theme):
        super().__init__()

        self.running = True

        self.queue = queue
        self.llm_pipeline = app_params["llm"]
        self.sd_engine = app_params["sd"]
        self.theme = theme

        print("theme: ", self.theme)
        
        self.compiled_model = app_params["super_res_compiled_model"] 
        self.upsample_factor = app_params["super_res_upsample_factor"]
        self.depth_compiled_model = app_params["depth_compiled_model"]

    def sd_callback(self, i, num_inference_steps, callback_userdata):
        if num_inference_steps > 0:
            prog = int((i / num_inference_steps) * 100)
            self.progress_updated.emit(prog, "illustrating")

    def stop(self):
        self.running = False
        self.quit()
        self.wait()

    def produce_parallex_img(self, img):
        #img = cv2.cvtColor(np.array(img), cv2.COLOR_RGB2BGR)
        sr_out = self.run_sr(np.array(img))

        buffer = io.BytesIO()
        sr_out.save(buffer, format="PNG")
        buffer.seek(0)

        # Convert the image buffer to QPixmap
        pixmap = QPixmap()
        pixmap.loadFromData(buffer.read(), "PNG")
        
        # this updates the UI image. 
        self.primary_pixmap_updated.emit(pixmap)

        colored_depth = depth_map_parallax(self.depth_compiled_model, sr_out)

        buffer = io.BytesIO()
        colored_depth.save(buffer, format="PNG")
        buffer.seek(0)

        # Convert the image buffer to QPixmap
        pixmap = QPixmap()
        pixmap.loadFromData(buffer.read(), "PNG")
        
        # this updates the UI depth map
        self.depth_pixmap_updated.emit(pixmap)

        return sr_out

    def run_sr(self, img):
        input_image_original = np.expand_dims(img.transpose(2, 0, 1), axis=0)
        bicubic_image = cv2.resize(
        src=img, dsize=(768*self.upsample_factor, 432*self.upsample_factor), interpolation=cv2.INTER_CUBIC)
        input_image_bicubic = np.expand_dims(bicubic_image.transpose(2, 0, 1), axis=0)

        original_image_key, bicubic_image_key = self.compiled_model.inputs
        output_key = self.compiled_model.output(0)

        result = self.compiled_model(
        {
            original_image_key.any_name: input_image_original,
            bicubic_image_key.any_name: input_image_bicubic,
        }
        )[output_key]

        result_image = convert_result_to_image(result)

        return result_image

    def generate_image(self, prompt):
        image_tensor = self.sd_engine.generate(
        prompt,
        width=768,
        height=432,
        num_inference_steps=5,
        num_images_per_prompt=1)

        sd_output = Image.fromarray(image_tensor.data[0])

        sr_out = self.produce_parallex_img(sd_output)

    def llm_streamer(self, subword):
        print(subword, end='', flush=True)

        self.stream_message += subword

        search_string = "SD Prompt:"
        if search_string in self.stream_message and 'None' not in self.stream_message:
            if self.stream_sd_prompt_index is None:
                self.stream_sd_prompt_index = self.stream_message.find(search_string)

            start_index = self.stream_sd_prompt_index
            # Calculate the start index of the new string (1 character past the ':')
            prompt = self.stream_message[start_index + len(search_string):].strip()

            self.caption_updated.emit(prompt)
        elif 'None' in self.stream_message:
            #Sometimes the LLM gives a response like: None (And then some long description why in parenthesis)
            # Basically, as soon as we see 'None', just stop generating tokens.
            return True

        # Return flag corresponds whether generation should be stopped.
        # False means continue generation.
        return False

    def run(self):

        llm_tokenizer = self.llm_pipeline.get_tokenizer()

        # Assemble the system message.
        system_message=LLM_SYSTEM_MESSAGE_START
        system_message+="\nThe presenter is giving a hint that the theme of their game is: " + self.theme
        system_message+="\nYou should use this hint to guide your decision about whether the presenter is describing a scene from the game, or not, and also to generate adequate SD Prompts."
        system_message+="\n" + LLM_SYSTEM_MESSAGE_END
        #print("System Message:")
        #print(system_message)


        generate_config = ov_genai.GenerationConfig()

        generate_config.temperature = 0.7
        generate_config.top_p = 0.95
        generate_config.max_length = 2048

        meaningful_message_pairs = []

        while self.running:
            try:
                # Wait for a sentence from the queue
                self.progress_updated.emit(0, "listening")

                result = self.queue.get(timeout=1)

                self.progress_updated.emit(0, "processing")

                chat_history = [{"role": "system", "content": system_message}]

                #only keep the latest 5 meaningful message pairs (last 2 illustrations)
                meaningful_message_pairs = meaningful_message_pairs[-2:]

                formatted_prompt = system_message
                
                #print("number of meaningful messages in history: ", len(meaningful_message_pairs))
                for meaningful_pair in meaningful_message_pairs:
                    user_message = meaningful_pair[0]
                    assistant_response = meaningful_pair[1]

                    chat_history.append({"role": "user", "content": user_message["content"]})
                    chat_history.append({"role": "assistant", "content": assistant_response["content"]})

                chat_history.append({"role": "user", "content": result})
                formatted_prompt = llm_tokenizer.apply_chat_template(history=chat_history, add_generation_prompt=True)

                self.progress_updated.emit(0, "processing...")
                self.stream_message=""
                self.stream_sd_prompt_index=None
                print("running llm!")
                llm_result = self.llm_pipeline.generate(inputs=formatted_prompt, generation_config=generate_config, streamer=self.llm_streamer)

                search_string = "SD Prompt:"

                #sometimes the llm will return 'SD Prompt: None', so filter out that case.
                if search_string in llm_result and 'None' not in llm_result:
                    # Find the start of the search string
                    start_index = llm_result.find(search_string)
                    # Calculate the start index of the new string (1 character past the ':')
                    prompt = llm_result[start_index + len(search_string):].strip()
                    #print(f"Extracted prompt: '{prompt}'")

                    caption = prompt
                    self.caption_updated.emit(caption)
                    print("calling self.generate_image...")
                    self.progress_updated.emit(0, "illustrating...")
                    
                    self.generate_image(prompt)
                    #self.image_updated.emit(pixmap)  # Emit the QPixmap

                    # this was a meaningful message!
                    meaningful_message_pairs.append(
                    ({"role": "user", "content": result},
                     {"role": "assistant", "content": llm_result},)
                    )

            except Empty:
                continue  # Queue is empty, just wait

        self.progress_updated.emit(0, "idle")


class ClickableLabel(QLabel):
    clicked = Signal()  # Define a signal to emit on click

    def mousePressEvent(self, event):
        self.clicked.emit()  # Emit the clicked signal
        super().mousePressEvent(event)

class MainWindow(QMainWindow):
    def __init__(self, app_params):
        super().__init__()

        # Main widget and layout
        self.central_widget = QWidget()
        self.setCentralWidget(self.central_widget)
        layout = QGridLayout(self.central_widget)

        self.llm_pipeline = app_params["llm"]
        self.sd_engine = app_params["sd"]
        self.app_params = app_params

        # Image pane
        self.image_label = ClickableLabel("No Image")
        #self.image_label.setFixedSize(1280, 720)
        self.image_label.setFixedSize(1216, 684)
        self.image_label.setStyleSheet("border: 1px solid black;")
        self.image_label.setAlignment(Qt.AlignCenter)
        layout.addWidget(self.image_label, 0, 1)

        # Connect the click signal
        self.display_primary_img = True
        self.image_label.clicked.connect(self.swap_image)

        self.primary_pixmap = None
        self.depth_pixmap = None

        # Caption
        self.caption_label = QLabel("No Caption")
        fantasy_font = QFont("Papyrus", 18, QFont.Bold)
        self.caption_label.setFont(fantasy_font)
        self.caption_label.setAlignment(Qt.AlignCenter)
        self.caption_label.setWordWrap(True)  # Enable word wrapping
        layout.addWidget(self.caption_label, 1, 1)

        # Log widget
        self.log_widget = QTextEdit()
        self.log_widget.setReadOnly(True)
        self.log_widget.setStyleSheet("background-color: #f0f0f0; border: 1px solid gray;")
        layout.addWidget(self.log_widget, 0, 2, 2, 1)
        self.log_widget.hide()  # Initially hidden

        bottom_layout = QVBoxLayout()

        # Bottom pane with buttons and progress bar
        button_layout = QHBoxLayout()
        self.start_button = QPushButton("Start")
        self.start_button.clicked.connect(self.start_thread)
        button_layout.addWidget(self.start_button)

        self.toggle_theme_button = QPushButton("Theme")
        self.toggle_theme_button.clicked.connect(self.toggle_theme)
        button_layout.addWidget(self.toggle_theme_button)

        self.progress_bar = QProgressBar()
        self.progress_bar.setFormat("Idle")
        self.progress_bar.setValue(0)
        button_layout.addWidget(self.progress_bar)

        bottom_layout.addLayout(button_layout)

        # Theme text box, initially hidden
        self.theme_input = QLineEdit()
        self.theme_input.setPlaceholderText("Enter a theme here...")
        self.theme_input.setText("Medieval Fantasty Adventure")
        self.theme_input.setStyleSheet("background-color: white; color: black;")
        self.theme_input.hide()
        bottom_layout.addWidget(self.theme_input)

        layout.addLayout(bottom_layout, 2, 0, 1, 3)

        # Worker threads
        self.speech_thread = None
        self.worker = None

        # Window configuration
        self.setWindowTitle("AI Adventure Experience")
        self.resize(800, 600)

    def start_thread(self):
        if not self.worker or not self.worker.isRunning():

            self.vad_worker = VADWorker()
            self.vad_worker.start()

            self.whisper_worker = WhisperWorker(self.vad_worker.result_queue, self.app_params["whisper_device"])
            self.whisper_worker.start()

            self.queue = self.whisper_worker.result_queue

            self.worker = WorkerThread(self.queue, self.app_params, self.theme_input.text())
            self.worker.image_updated.connect(self.update_image)
            self.worker.primary_pixmap_updated.connect(self.update_primary_pixmap)
            self.worker.depth_pixmap_updated.connect(self.update_depth_pixmap)
            self.worker.caption_updated.connect(self.update_caption)
            self.worker.progress_updated.connect(self.update_progress)
            self.worker.start()
            self.start_button.setText("Stop")

        else:
            self.worker.stop()
            self.worker = None

            self.vad_worker.stop()
            self.whisper_worker.stop()

            #self.worker.terminate()
            self.start_button.setText("Start")

            self.queue = None

    def toggle_log(self):
        if self.log_widget.isVisible():
            self.log_widget.hide()
        else:
            self.log_widget.show()

    def toggle_theme(self):
        if self.theme_input.isVisible():
            self.theme_input.hide()
        else:
            self.theme_input.show()

    def update_depth_pixmap(self, pixmap):
        self.depth_pixmap = pixmap

        self.update_image_label()

    def update_primary_pixmap(self, pixmap):
        self.primary_pixmap = pixmap

        self.update_image_label()

    def update_image_label(self):
        if self.display_primary_img and self.primary_pixmap is not None:
            pixmap = self.primary_pixmap
            self.image_label.setPixmap(pixmap.scaled(self.image_label.size()))
        elif not self.display_primary_img and self.depth_pixmap is not None:
            pixmap = self.depth_pixmap
            self.image_label.setPixmap(pixmap.scaled(self.image_label.size()))

    def update_image(self, pixmap):
        print("not doing anything...")
        #pixmap = QPixmap.fromImage(image)
        #self.image_label.setPixmap(pixmap.scaled(self.image_label.size()))

    def swap_image(self):
        self.display_primary_img = (not self.display_primary_img)
        self.update_image_label()

    def update_caption(self, caption):
        self.caption_label.setText(caption)
        #self.log_widget.append(f"Caption updated: {caption}")

    def update_progress(self, value, label):
        self.progress_bar.setValue(value)
        self.progress_bar.setFormat(label)

    def closeEvent(self, event):
        if self.worker and self.worker.isRunning():
            self.vad_worker.stop()
            self.whisper_worker.stop()
            self.worker.stop()  # Gracefully stop the worker thread
            self.worker.wait()  # Wait for the thread to finish

        event.accept()  # Proceed with closing the application


def main():
    # Set the Qt platform to offscreen for headless environments
    os.environ["QT_QPA_PLATFORM"] = "offscreen"
    
    app = QApplication(sys.argv)
    core = ov.Core()

    llm_device = "GPU" if "GPU" in core.available_devices else "CPU"
    sd_device = "GPU" if "GPU" in core.available_devices else "CPU"
    whisper_device = 'CPU'
    super_res_device = "GPU" if "GPU" in core.available_devices else "CPU"
    depth_anything_device = "GPU" if "GPU" in core.available_devices else "CPU"

    print("Just a minute... doing some application setup...")

    # create the 'results' folder if it doesn't exist
    Path("results").mkdir(exist_ok=True)

    app_params = {}

    # creating the LLM pipeline

    print("Creating an llm pipeline to run on ", llm_device)

    llm_model_path = r"./models/tinyllama-1.1b-chat-v1.0/INT4_compressed_weights"

    if llm_device == 'NPU':
        pipeline_config = {"MAX_PROMPT_LEN": 1536}
        llm_pipe = ov_genai.LLMPipeline(llm_model_path, llm_device, pipeline_config)
    else:
        llm_pipe = ov_genai.LLMPipeline(llm_model_path, llm_device)

    app_params["llm"] = llm_pipe

    print("Done creating our llm..")

    print("Creating a stable diffusion pipeline to run on ", sd_device)

    sd_pipe = ov_genai.Text2ImagePipeline(r"models/LCM_Dreamshaper_v7/FP16", sd_device)

    app_params["sd"] = sd_pipe
    print("done creating the stable diffusion pipeline...")

    app_params["whisper_device"] = whisper_device

    print("Initializing Super Res Model to run on ", super_res_device)
    model_path_sr = Path(f"models/single-image-super-resolution-1033.xml")  # realesrgan.xml")
    super_res_compiled_model, super_res_upsample_factor = superres_load(model_path_sr, super_res_device, h_custom=432,
                                                                        w_custom=768)
    app_params["super_res_compiled_model"] = super_res_compiled_model
    app_params["super_res_upsample_factor"] = super_res_upsample_factor
    print("Initializing Super Res Model done...")

    print("Initializing Depth Anything v2 model to run on ", depth_anything_device)
    OV_DEPTH_ANYTHING_PATH = Path(f"models/depth_anything_v2_vits.xml")
    depth_compiled_model = core.compile_model(OV_DEPTH_ANYTHING_PATH, device_name=depth_anything_device)
    app_params["depth_compiled_model"] = depth_compiled_model
    print("Initializing Depth Anything v2 done...")

    print("Demo is ready!")
    window = MainWindow(app_params)
    window.show()

    sys.exit(app.exec())


if __name__ == "__main__":
    main()