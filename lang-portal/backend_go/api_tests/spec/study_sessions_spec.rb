require 'spec_helper'

RSpec.describe 'Study Sessions API' do
  describe 'GET /study_sessions' do
    it 'returns a list of study sessions with pagination' do
      response = APIHelper.get('/study_sessions')
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

  describe 'GET /study_sessions/:id' do
    it 'returns a specific study session' do
      response = APIHelper.get('/study_sessions/1')
      expect(response.code).to eq(200)
      
      json = JSON.parse(response.body)
      expect(json).to include(
        'id',
        'activity_name',
        'group_name',
        'start_time',
        'end_time',
        'review_items_count'
      )

      # Type checking
      expect(json['id']).to be_a(Integer)
      expect(json['activity_name']).to be_a(String)
      expect(json['group_name']).to be_a(String)
      expect(json['review_items_count']).to be_a(Integer)
      expect(json['start_time']).to match(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}/)
      expect(json['end_time']).to match(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}/)
    end
  end

  describe 'POST /study_sessions/:id/words/:word_id/review' do
    let(:valid_params) do
      {
        correct: true
      }
    end

    it 'records a word review' do
      response = APIHelper.post('/study_sessions/1/words/1/review', valid_params)
      expect(response.code).to eq(200)
      
      json = JSON.parse(response.body)
      expect(json).to include(
        'success',
        'word_id',
        'study_session_id',
        'correct',
        'created_at'
      )

      # Type checking
      expect(json['success']).to be true
      expect(json['word_id']).to be_a(Integer)
      expect(json['study_session_id']).to be_a(Integer)
      expect(json['correct']).to be true
      expect(json['created_at']).to match(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}/)
    end
  end
end  
