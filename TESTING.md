# Testing Debate Generation System

This guide will help you test the AI debate generation system locally before adding debates to the database.

## 🚀 Quick Start

1. **Run the setup script:**

   ```bash
   ./setup_local_testing.sh
   ```

2. **Add your API keys to the `.env` file:**
   - `OPENAI_API_KEY` (required for AI generation)
   - `FOOTBALL_API_KEY` (optional, for real match data)
   - `RAPID_API_KEY` (optional, for additional data sources)

## 🧪 Testing Options

### Option 1: AI Prompt Generation Only (Recommended for Initial Testing)

This tests the AI prompt generation without requiring a database or other services.

```bash
go run cmd/test_ai/main.go
```

**What this tests:**

- ✅ AI prompt generation with mock data
- ✅ Pre-match and post-match debate generation
- ✅ JSON output formatting
- ✅ Error handling

**Requirements:**

- Only `OPENAI_API_KEY` is required

### Option 2: Full API Testing (Requires Database)

This tests the complete API including database operations.

```bash
# Start the API server
go run main.go

# In another terminal, run the test script
./test_debate_generation.sh
```

**What this tests:**

- ✅ Complete debate creation in database
- ✅ Debate cards creation
- ✅ Analytics setup
- ✅ API endpoints
- ✅ Error handling

**Requirements:**

- `OPENAI_API_KEY`
- `DB_URL` (PostgreSQL)
- `REDIS_URL` (Redis)
- `FOOTBALL_API_KEY` (for real match data)

### Option 3: Docker Testing

If you have Docker installed, you can test with the full stack.

```bash
# Start all services
docker-compose up -d

# Run tests
./test_debate_generation.sh
```

## 📋 Test Data

The AI test uses mock data for a Manchester City vs Liverpool match:

- **Match ID:** 1321727
- **Teams:** Manchester City (Home) vs Liverpool (Away)
- **Status:** NS (Not Started) for pre-match, FT (Finished) for post-match
- **Lineups:** Mock starting XI for both teams
- **News:** Sample headlines about the match
- **Social Sentiment:** Mock sentiment data

## 🔍 What to Look For

### Quality Indicators

**Good debate prompts should have:**

- ✅ Engaging, controversial headlines
- ✅ Clear, specific debate topics
- ✅ Balanced perspectives (agree/disagree/wildcard)
- ✅ Relevant to the match context
- ✅ Appropriate for the match status (pre/post)

**Example of a good pre-match debate:**

```
Headline: "Will Manchester City's possession game dominate Liverpool's pressing?"
Cards:
1. [agree] "City's technical superiority will control the tempo"
2. [disagree] "Liverpool's gegenpressing will force turnovers"
3. [wildcard] "The referee's decisions will be the deciding factor"
```

### Red Flags

**Watch out for:**

- ❌ Generic, non-specific debates
- ❌ Debates not related to football
- ❌ Inappropriate content
- ❌ Missing or invalid stance types
- ❌ Empty or very short descriptions

## 🛠️ Troubleshooting

### Common Issues

1. **"OPENAI_API_KEY environment variable is required"**

   - Make sure you've added your OpenAI API key to the `.env` file

2. **"Failed to generate AI prompt"**

   - Check your OpenAI API key is valid
   - Verify you have sufficient API credits
   - Check your internet connection

3. **"Failed to connect to Database"**

   - Ensure PostgreSQL is running
   - Check your `DB_URL` in the `.env` file
   - Run database migrations if needed

4. **"Failed to connect to Redis"**
   - Ensure Redis is running
   - Check your `REDIS_URL` in the `.env` file

### Debug Mode

To see more detailed output, you can modify the test scripts to include verbose logging:

```bash
# For AI testing
go run cmd/test_ai/main.go 2>&1 | tee test_output.log

# For API testing
go run main.go -debug
```

## 📊 Expected Results

### AI Prompt Generation Test

You should see output like:

```
🤖 Testing AI Prompt Generation
===============================

📋 Testing Pre-Match Prompt Generation...
✅ Pre-match prompt generated successfully!

Pre-Match Debate Prompt:
Headline: Will Manchester City's possession game dominate Liverpool's pressing?
Description: A tactical battle between two contrasting styles...
Cards (3):
  1. [agree] City's technical superiority will control the tempo
     Description: Pep Guardiola's possession-based approach...
  2. [disagree] Liverpool's gegenpressing will force turnovers
     Description: Jurgen Klopp's high-intensity pressing...
  3. [wildcard] The referee's decisions will be the deciding factor
     Description: Key calls could swing momentum either way...

JSON Output:
{
  "headline": "Will Manchester City's possession game dominate Liverpool's pressing?",
  "description": "A tactical battle between two contrasting styles...",
  "cards": [...]
}
```

### API Test Results

You should see:

- ✅ Health check passes
- ✅ AI prompt generation works
- ✅ Debate creation succeeds
- ✅ Database records created
- ✅ Proper error handling for invalid inputs

## 🎯 Next Steps

Once you're satisfied with the AI prompt quality:

1. **Review the generated debates** for accuracy and engagement
2. **Test with different match scenarios** (different teams, statuses)
3. **Adjust the AI prompts** if needed (in `internal/ai/prompt_generator.go`)
4. **Deploy to staging** for further testing
5. **Monitor real-world usage** and gather feedback

## 📞 Support

If you encounter issues:

1. Check the troubleshooting section above
2. Review the logs for error messages
3. Verify all environment variables are set correctly
4. Test with the AI-only option first to isolate issues
