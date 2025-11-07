#!/bin/bash

echo "Launching servers in new Git Bash windows..."

# Get the path to git-bash executable
GIT_BASH="/c/Program Files/Git/git-bash.exe"

# Start each server in a new Git Bash window
"$GIT_BASH" -c "cd 'fyp/cmd' && go run . --port=8000" &
sleep 1

"$GIT_BASH" -c "cd 'fyp/cmd' && go run . --port=8001" &
sleep 1

"$GIT_BASH" -c "cd 'fyp/cmd' && go run . --port=8002" &

echo "✅ All servers launched in new windows!"
echo "Check your taskbar for new Git Bash windows"