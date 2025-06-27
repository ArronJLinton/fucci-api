#!/bin/bash

echo "🚀 Setting up Local Testing Environment for Debate Generation"
echo "============================================================="

# Check if .env file exists
if [ ! -f .env ]; then
    echo "📝 Creating .env file from template..."
    cp env.example .env
    echo "✅ .env file created. Please edit it with your API keys."
    echo ""
    echo "Required API Keys:"
    echo "  - OPENAI_API_KEY: Get from https://platform.openai.com/api-keys"
    echo "  - FOOTBALL_API_KEY: Get from https://www.api-football.com/"
    echo "  - RAPID_API_KEY: Get from https://rapidapi.com/"
    echo ""
    echo "Optional (for full testing):"
    echo "  - DB_URL: PostgreSQL connection string"
    echo "  - REDIS_URL: Redis connection string"
    echo ""
    echo "⚠️  Please edit .env file and add your API keys before continuing."
    exit 1
else
    echo "✅ .env file already exists"
fi

# Check if required environment variables are set
echo ""
echo "🔍 Checking environment variables..."

# Source the .env file
set -a
source .env
set +a

# Check OpenAI API key
if [ -z "$OPENAI_API_KEY" ] || [ "$OPENAI_API_KEY" = "your_openai_api_key_here" ]; then
    echo "❌ OPENAI_API_KEY is not set or is using placeholder value"
    echo "   Please edit .env file and add your OpenAI API key"
    exit 1
else
    echo "✅ OPENAI_API_KEY is set"
fi

# Check Football API key
if [ -z "$FOOTBALL_API_KEY" ] || [ "$FOOTBALL_API_KEY" = "your_football_api_key_here" ]; then
    echo "⚠️  FOOTBALL_API_KEY is not set - some features may not work"
else
    echo "✅ FOOTBALL_API_KEY is set"
fi

# Check Rapid API key
if [ -z "$RAPID_API_KEY" ] || [ "$RAPID_API_KEY" = "your_rapid_api_key_here" ]; then
    echo "⚠️  RAPID_API_KEY is not set - some features may not work"
else
    echo "✅ RAPID_API_KEY is set"
fi

echo ""
echo "🧪 Testing Options:"
echo "=================="
echo ""
echo "1. Test AI Prompt Generation Only (No Database Required):"
echo "   go run cmd/test_ai/main.go"
echo ""
echo "2. Test Full API with Database (Requires DB_URL and REDIS_URL):"
echo "   go run main.go"
echo "   Then run: ./test_debate_generation.sh"
echo ""
echo "3. Test with Docker (if you have Docker installed):"
echo "   docker-compose up -d"
echo "   Then run: ./test_debate_generation.sh"
echo ""

# Check if jq is installed for JSON formatting
if ! command -v jq &> /dev/null; then
    echo "⚠️  jq is not installed. Install it for better JSON formatting in tests:"
    echo "   macOS: brew install jq"
    echo "   Ubuntu: sudo apt-get install jq"
    echo "   Windows: choco install jq"
fi

echo "🎉 Setup complete! Choose a testing option above." 