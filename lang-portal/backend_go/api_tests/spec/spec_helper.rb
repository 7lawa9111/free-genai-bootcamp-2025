require 'rspec'
require 'httparty'
require 'json'

# Define APIHelper module to match what the tests expect
module APIHelper
  def self.get(path, params = {})
    url = "http://localhost:8080/api#{path}"
    HTTParty.get(url, body: params.to_json, headers: { 'Content-Type' => 'application/json' })
  end

  def self.post(path, params = {})
    url = "http://localhost:8080/api#{path}"
    HTTParty.post(url, body: params.to_json, headers: { 'Content-Type' => 'application/json' })
  end

  def self.put(path, params = {})
    url = "http://localhost:8080/api#{path}"
    HTTParty.put(url, body: params.to_json, headers: { 'Content-Type' => 'application/json' })
  end

  def self.delete(path, params = {})
    url = "http://localhost:8080/api#{path}"
    HTTParty.delete(url, body: params.to_json, headers: { 'Content-Type' => 'application/json' })
  end
end

RSpec.configure do |config|
  # No need to include anything since we're using module methods
  config.expect_with :rspec do |expectations|
    expectations.include_chain_clauses_in_custom_matcher_descriptions = true
  end

  config.mock_with :rspec do |mocks|
    mocks.verify_partial_doubles = true
  end

  config.shared_context_metadata_behavior = :apply_to_host_groups
end 