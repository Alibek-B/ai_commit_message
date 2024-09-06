# frozen_string_literal: true

require 'httparty'
require 'json'

requirement = <<~TEXT
  Write a git commit message for these changes.
  The commit message should contain no more than 20 words.
  The response should contain only the message.
TEXT

git_status = `git status`
git_diff = `git diff`

prompt = "#{requirement} #{git_status} #{git_diff}".gsub("\n", '').gsub('"', "'")

response = HTTParty.post(
  'http://localhost:11434/api/generate',
  headers: { 'Content-Type' => 'application/json' },
  body: {
    model: 'gemma2:2b',
    prompt: prompt,
    stream: false
  }.to_json
)

message = JSON.parse(response.body)['response']

puts "Commit message: #{message}"
print 'Enter y(yes) if you want to commit the changes: '
choice = gets.strip

if choice == 'y'
  `git add .`
  `git commit -am #{message}`
else
  puts 'commit aborted!'
end
