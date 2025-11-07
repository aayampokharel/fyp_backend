#!/bin/bash
echo "🚀 Server starting on port 8002..."
echo "📁 Directory: $(pwd)"
cd fyp/cmd
echo "🏃 Running: go run . --port=8002"
go run . --port=8002
