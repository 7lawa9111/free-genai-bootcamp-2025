require 'spec_helper'

RSpec.describe 'Reset API' do
  describe 'POST /reset_history' do
    before do
      # Initialize test data
      APIHelper.post('/test/init_data')
    end

    it 'resets study history' do
      response = APIHelper.post('/reset_history')
      expect(response.code).to eq(200)
      
      json = JSON.parse(response.body)
      expect(json).to include(
        'success',
        'message'
      )

      # Check values
      expect(json['success']).to be true
      expect(json['message']).to eq('Study history has been reset')

      # Verify reset by checking study progress
      progress = APIHelper.get('/dashboard/study_progress')
      progress_json = JSON.parse(progress.body)
      expect(progress_json['total_words_studied']).to eq(0)
    end
  end

  # ... rest of the file remains the same ...
end 