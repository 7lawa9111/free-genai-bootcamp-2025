require 'spec_helper'

RSpec.describe 'Dashboard API' do
  describe 'GET /dashboard/quick-stats' do
    it 'returns dashboard statistics' do
      response = APIHelper.get('/dashboard/quick-stats')
      expect(response.code).to eq(200)

      json = JSON.parse(response.body)
      expect(json).to include(
        'total_words',
        'words_studied',
        'study_sessions',
        'accuracy_rate'
      )
    end
  end
end 