#!/bin/bash
echo "🚀 Server starting on port 8000..."
echo "📁 Directory: $(pwd)"
cd fyp/cmd
echo "🏃 Running: go run . --port=8000"
go run . --port=8000
