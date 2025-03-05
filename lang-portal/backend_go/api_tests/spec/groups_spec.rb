require 'spec_helper'

RSpec.describe 'Groups API' do
  describe 'GET /groups' do
    it 'returns a list of groups with pagination' do
      response = APIHelper.get('/groups')
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
        group = json['items'].first
        expect(group).to include(
          'id',
          'name',
          'word_count'
        )

        # Type checking
        expect(group['id']).to be_a(Integer)
        expect(group['name']).to be_a(String)
        expect(group['word_count']).to be_a(Integer)
      end

      # Verify pagination settings from specs
      expect(json['pagination']['items_per_page']).to eq(100)
    end
  end

  describe 'GET /groups/:id' do
    it 'returns a single group with stats' do
      response = APIHelper.get('/groups/1')
      expect(response.code).to eq(200)
      
      json = JSON.parse(response.body)
      expect(json).to include(
        'id',
        'name',
        'stats'
      )

      # Check stats structure
      expect(json['stats']).to include('total_word_count')
      
      # Type checking
      expect(json['id']).to be_a(Integer)
      expect(json['name']).to be_a(String)
      expect(json['stats']['total_word_count']).to be_a(Integer)
    end

    it 'returns 404 for non-existent group' do
      response = APIHelper.get('/groups/999999')
      expect(response.code).to eq(404)
    end
  end

  describe 'GET /groups/:id/words' do
    it 'returns words for a specific group with pagination' do
      response = APIHelper.get('/groups/1/words')
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
        word = json['items'].first
        expect(word).to include(
          'japanese',
          'romaji',
          'english',
          'correct_count',
          'wrong_count'
        )

        # Type checking
        expect(word['japanese']).to be_a(String)
        expect(word['romaji']).to be_a(String)
        expect(word['english']).to be_a(String)
        expect(word['correct_count']).to be_a(Integer)
        expect(word['wrong_count']).to be_a(Integer)
      end
    end
  end

  describe 'GET /groups/:id/study_sessions' do
    it 'returns study sessions for a specific group with pagination' do
      response = APIHelper.get('/groups/1/study_sessions')
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
        expect(session['start_time']).to match(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}/)
        expect(session['end_time']).to match(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}/)
        expect(session['review_items_count']).to be_a(Integer)
      end
    end
  end
end 