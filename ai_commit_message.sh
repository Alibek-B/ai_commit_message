GIT_STATUS=$(git status)
GIT_DIFF=$(git diff)

REQUIREMENT="Write a git commit message for these changes.
The commit message should contain no more than 20 words.
The response should contain only the message."

MESSAGE="$REQUIREMENT $GIT_STATUS $GIT_DIFF"

PROMPT=$(echo $MESSAGE | tr -d '\n' | sed 's/"//g')

RESPONSE=$(curl -X POST http://localhost:11434/api/generate -H "Content-Type: application/json" -d '{
  "model": "gemma2:2b",
  "prompt": "'"$PROMPT"'",
  "stream": false
}')

COMMIT_MESSAGE=$(echo $RESPONSE | jq -r '.response')

git add .
git commit -am "$COMMIT_MESSAGE"
