require 'spec_helper'

RSpec.describe 'Words API' do
  describe 'GET /words' do
    it 'returns a list of words' do
      # Reset and initialize test data
      APIHelper.post('/full_reset')
      APIHelper.post('/test/init_data')  # We'll add this endpoint

      response = APIHelper.get('/words')
      expect(response.code).to eq(200)
      
      json = JSON.parse(response.body)
      expect(json).to have_key('items')
      expect(json).to have_key('pagination')
    end
  end

  describe 'GET /words/:id' do
    it 'returns a specific word' do
      # Reset and initialize test data
      APIHelper.post('/full_reset')
      APIHelper.post('/test/init_data')  # We'll add this endpoint

      response = APIHelper.get('/words/1')
      expect(response.code).to eq(200)
      
      json = JSON.parse(response.body)
      expect(json).to include(
        'id',
        'japanese',
        'romaji',
        'english',
        'stats',
        'groups'
      )
    end

    it 'returns 404 for non-existent word' do
      response = APIHelper.get('/words/999')
      expect(response.code).to eq(404)
    end
  end
end 