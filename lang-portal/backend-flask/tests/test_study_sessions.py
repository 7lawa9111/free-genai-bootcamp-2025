import pytest
from datetime import datetime

def test_create_study_session(client):
    # Test valid creation
    response = client.post('/api/study-sessions', json={
        'group_id': 1,
        'study_activity_id': 1
    })
    print("Response:", response.json)
    assert response.status_code == 201
    assert 'id' in response.json
    assert response.json['group_id'] == 1
    assert response.json['activity_id'] == 1

    # Test missing field
    response = client.post('/api/study-sessions', json={
        'group_id': 1
    })
    assert response.status_code == 400
    assert 'error' in response.json 