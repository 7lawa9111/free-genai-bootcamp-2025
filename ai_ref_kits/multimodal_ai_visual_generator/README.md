# AI Visual Generator

A web application that generates images from text descriptions using Stable Diffusion.

## Features

- Text-to-image generation using Stable Diffusion
- Modern web interface with real-time status updates
- Session-based image generation
- Support for both CPU and GPU inference

## Prerequisites

- Python 3.8 or higher
- CUDA-capable GPU (optional, but recommended for faster inference)
- System libraries for OpenCV and Qt:
  ```bash
  sudo apt-get update && sudo apt-get install -y libgl1-mesa-glx libegl1 libxkbcommon-x11-0 libdbus-1-3 libxcb-icccm4 libxcb-image0 libxcb-keysyms1 libxcb-randr0 libxcb-render-util0 libxcb-xinerama0 libxcb-xfixes0 libxcb-shape0
  ```

## Installation

1. Create a virtual environment (recommended):
```bash
python -m venv venv
source venv/bin/activate  # On Windows, use: venv\Scripts\activate
```

2. Install the required packages:
```bash
pip install -r requirements.txt
```

## Usage

1. Start the Flask application using the wrapper script:
```bash
python run_app.py
```

2. Open your web browser and navigate to:
```
http://localhost:5000
```

3. Enter a text description in the input field and click "Generate Image"

4. Wait for the image generation to complete. The generated image will be displayed on the page.

## Project Structure

```
multimodal_ai_visual_generator/
├── app.py                  # Main Flask application
├── patch_huggingface.py    # Patch for huggingface_hub compatibility
├── run_app.py              # Wrapper script to run the application
├── requirements.txt        # Python dependencies
├── static/                 # Static files (CSS, JS, generated images)
│   └── generated/          # Directory for storing generated images
└── templates/              # HTML templates
    └── index.html          # Main web interface
```

## Notes

- The first run will download the Stable Diffusion model weights (approximately 4GB)
- Image generation may take several seconds to minutes depending on your hardware
- For best performance, use a CUDA-capable GPU
- The application uses a patch for huggingface_hub to ensure compatibility with diffusers

## License

This project is licensed under the MIT License - see the LICENSE file for details.