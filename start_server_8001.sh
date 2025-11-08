#!/bin/bash
echo "🚀 Server starting on port 8001..."
echo "📁 Directory: $(pwd)"
cd fyp/cmd
echo "🏃 Running: go run . --port=8001"
go run . --port=8001
