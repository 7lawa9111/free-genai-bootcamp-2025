#!/usr/bin/env python
"""
Wrapper script to run the application with the huggingface_hub patch.
"""

# Import and run the patch
import patch_huggingface

# Import and run the application
import app

if __name__ == "__main__":
    app.app.run(host='0.0.0.0', port=5000, debug=True) 