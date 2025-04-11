#!/usr/bin/env python3
import os
import sys
from web_app import app

if __name__ == '__main__':
    # Set the Qt platform to offscreen for headless environments
    os.environ["QT_QPA_PLATFORM"] = "offscreen"
    
    # Run the Flask app
    app.run(host='0.0.0.0', port=8080, debug=True) 