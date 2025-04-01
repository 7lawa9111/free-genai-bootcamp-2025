from flask import Flask, g, jsonify
from flask_cors import CORS
import os
import sqlite3
from datetime import datetime

from lib.db import Db

import routes.words
import routes.groups
import routes.study_sessions
import routes.dashboard
import routes.study_activities

def create_app(test_config=None):
    app = Flask(__name__)
    
    # Simple CORS configuration
    CORS(app, resources={
        r"/*": {
            "origins": ["http://localhost:5173"],  # Your Vite dev server
            "methods": ["GET", "POST", "PUT", "DELETE", "OPTIONS"],
            "allow_headers": ["Content-Type", "Authorization", "Accept"]
        }
    })
    
    if test_config is None:
        # Configure database path
        db_path = os.path.join(os.path.dirname(__file__), 'word_bank.db')
        app.config['DATABASE'] = db_path
        app.db = Db(database=db_path)
    else:
        app.config.update(test_config)

    # Close database connection
    @app.teardown_appcontext
    def close_db(exception):
        """Close the database connection when the application context ends"""
        if hasattr(g, 'db'):
            g.db.close()

    # load routes -----------
    routes.words.load(app)
    routes.groups.load(app)
    routes.study_sessions.load(app)
    routes.dashboard.load(app)
    routes.study_activities.load(app)
    
    return app

app = create_app()

if __name__ == '__main__':
    app.run(debug=True, port=5001)