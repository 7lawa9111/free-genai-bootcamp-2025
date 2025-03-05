require 'spec_helper'

RSpec.describe 'Study Activities API' do
  describe 'GET /study_activities/:id' do
    it 'returns a specific study activity' do
      response = APIHelper.get('/study_activities/1')
      expect(response.code).to eq(200)
      
      json = JSON.parse(response.body)
      expect(json).to include(
        'id',
        'name',
        'thumbnail_url',
        'description'
      )

      # Type checking
      expect(json['id']).to be_a(Integer)
      expect(json['name']).to be_a(String)
      expect(json['thumbnail_url']).to be_a(String)
      expect(json['description']).to be_a(String)
    end

    it 'returns 404 for non-existent activity' do
      response = APIHelper.get('/study_activities/999999')
      expect(response.code).to eq(404)
    end
  end

  describe 'POST /study_activities' do
    let(:valid_params) do
      {
        group_id: 1,
        study_activity_id: 1
      }
    end

    it 'creates a new study activity session' do
      response = APIHelper.post('/study_activities', valid_params)
      expect(response.code).to eq(201)
      
      json = JSON.parse(response.body)
      expect(json).to include(
        'id',
        'group_id'
      )

      # Type checking
      expect(json['id']).to be_a(Integer)
      expect(json['group_id']).to be_a(Integer)
    end

    it 'validates required parameters' do
      response = APIHelper.post('/study_activities', {})
      expect(response.code).to eq(400)
    end
  end

  describe 'GET /study_activities/:id/study_sessions' do
    it 'returns study sessions for an activity with pagination' do
      response = APIHelper.get('/study_activities/1/study_sessions')
      expect(response.code).to eq(200)
      
      json = JSON.parse(response.body)
      expect(json).to include('items', 'pagination')

      # Check pagination structure
      expect(json['pagination']).to include(
        'current_page',
        'total_pages',
        'total_items',
        'items_per_page'
      )

      # Check items structure if any exist
      if json['items'].any?
        session = json['items'].first
        expect(session).to include(
          'id',
          'activity_name',
          'group_name',
          'start_time',
          'end_time',
          'review_items_count'
        )

        # Type checking
        expect(session['id']).to be_a(Integer)
        expect(session['activity_name']).to be_a(String)
        expect(session['group_name']).to be_a(String)
        expect(session['review_items_count']).to be_a(Integer)
        expect(session['start_time']).to match(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}/)
        expect(session['end_time']).to match(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}/)
      end

      # Verify pagination settings from specs
      expect(json['pagination']['items_per_page']).to eq(100)
    end
  end
end 