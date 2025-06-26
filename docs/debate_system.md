# Debate System

The debate system is an AI-powered feature that generates engaging football debate topics and allows fans to vote and comment on controversial moments in matches.

## Features

### AI-Generated Debate Prompts

- **Pre-match prompts**: Generated before matches based on lineups, team form, and news
- **Post-match prompts**: Generated after matches based on statistics, key moments, and social sentiment
- **Multiple stances**: Each debate includes "Agree", "Disagree", and "Wildcard" perspectives

### Fan Engagement

- **Voting system**: Upvote/downvote or emoji reactions on debate cards
- **Comment threads**: Text-based discussions with nested replies
- **Analytics**: Track engagement scores and popular debates

### Data Integration

- **Football API**: Match data, lineups, statistics
- **Google News**: Recent headlines about teams and matches
- **Social Media**: Twitter sentiment analysis (placeholder for Reddit)
- **Caching**: Intelligent caching to reduce API calls

## API Endpoints

### Debate Management

#### Create Debate

```http
POST /v1/api/debates/
Content-Type: application/json

{
  "match_id": "12345",
  "debate_type": "pre_match",
  "headline": "Should the manager start the new signing?",
  "description": "The new striker has been in great form...",
  "ai_generated": true
}
```

#### Get Debate

```http
GET /v1/api/debates/{id}
```

#### Get Debates by Match

```http
GET /v1/api/debates/match?match_id=12345
```

#### Get Top Debates

```http
GET /v1/api/debates/top?limit=10
```

### Debate Cards

#### Create Debate Card

```http
POST /v1/api/debates/cards
Content-Type: application/json

{
  "debate_id": 1,
  "stance": "agree",
  "title": "The manager made the right call",
  "description": "The new signing brings fresh energy...",
  "ai_generated": true
}
```

### Voting

#### Create Vote

```http
POST /v1/api/debates/votes
Content-Type: application/json

{
  "debate_card_id": 1,
  "vote_type": "upvote"
}
```

#### Emoji Vote

```http
POST /v1/api/debates/votes
Content-Type: application/json

{
  "debate_card_id": 1,
  "vote_type": "emoji",
  "emoji": "👍"
}
```

### Comments

#### Create Comment

```http
POST /v1/api/debates/comments
Content-Type: application/json

{
  "debate_id": 1,
  "content": "I completely agree with this analysis!"
}
```

#### Nested Comment

```http
POST /v1/api/debates/comments
Content-Type: application/json

{
  "debate_id": 1,
  "parent_comment_id": 5,
  "content": "Great point! I think..."
}
```

#### Get Comments

```http
GET /v1/api/debates/{debateId}/comments
```

### AI Prompt Generation

#### Generate Pre-match Prompt

```http
GET /v1/api/debates/generate?match_id=12345&type=pre_match
```

#### Generate Post-match Prompt

```http
GET /v1/api/debates/generate?match_id=12345&type=post_match
```

## Database Schema

### Tables

#### debates

- `id`: Primary key
- `match_id`: Football match identifier
- `debate_type`: "pre_match" or "post_match"
- `headline`: Debate question/topic
- `description`: Additional context
- `ai_generated`: Whether AI created this debate
- `created_at`, `updated_at`: Timestamps

#### debate_cards

- `id`: Primary key
- `debate_id`: Foreign key to debates
- `stance`: "agree", "disagree", or "wildcard"
- `title`: Card title
- `description`: Supporting argument
- `ai_generated`: Whether AI created this card
- `created_at`, `updated_at`: Timestamps

#### votes

- `id`: Primary key
- `debate_card_id`: Foreign key to debate_cards
- `user_id`: Foreign key to users
- `vote_type`: "upvote", "downvote", or "emoji"
- `emoji`: Emoji character (for emoji votes)
- `created_at`: Timestamp

#### comments

- `id`: Primary key
- `debate_id`: Foreign key to debates
- `parent_comment_id`: Self-reference for nested comments
- `user_id`: Foreign key to users
- `content`: Comment text
- `created_at`, `updated_at`: Timestamps

#### debate_analytics

- `id`: Primary key
- `debate_id`: Foreign key to debates
- `total_votes`: Total vote count
- `total_comments`: Total comment count
- `engagement_score`: Calculated engagement metric
- `created_at`, `updated_at`: Timestamps

## Configuration

### Environment Variables

Add these to your `.env` file:

```env
# OpenAI API for AI prompt generation
OPENAI_API_KEY=your_openai_api_key
OPENAI_BASE_URL=https://api.openai.com/v1

# Existing APIs (already configured)
FOOTBALL_API_KEY=your_football_api_key
RAPID_API_KEY=your_rapid_api_key
REDIS_URL=redis://localhost:6379
DB_URL=your_database_url
```

## Data Flow

### AI Prompt Generation

1. **Match Info**: Fetch basic match details (teams, date, status)
2. **Lineups**: Get player lineups if match is upcoming/in-progress
3. **Statistics**: Get match stats if match is finished
4. **News Headlines**: Search Google News for team/match keywords
5. **Social Sentiment**: Analyze Twitter/Reddit sentiment (placeholder)
6. **AI Generation**: Send aggregated data to OpenAI for prompt generation
7. **Caching**: Cache generated prompts to avoid repeated API calls

### Debate Creation

1. **Manual**: Users can create debates manually
2. **AI-Generated**: System can auto-generate debates using AI
3. **Cards**: Each debate gets 3 cards (agree/disagree/wildcard)
4. **Analytics**: Engagement tracking starts immediately

## Usage Examples

### Creating an AI-Generated Debate

```bash
# Generate a pre-match debate
curl "http://localhost:8080/v1/api/debates/generate?match_id=12345&type=pre_match"

# Response includes headline, description, and 3 debate cards
{
  "headline": "Should the manager start the new signing?",
  "description": "The new striker has been in great form...",
  "cards": [
    {
      "stance": "agree",
      "title": "Start the new signing",
      "description": "He's been scoring consistently..."
    },
    {
      "stance": "disagree",
      "title": "Stick with the proven lineup",
      "description": "The current striker has chemistry..."
    },
    {
      "stance": "wildcard",
      "title": "Try a 4-4-2 formation",
      "description": "Play both strikers together..."
    }
  ]
}
```

### Voting on a Debate Card

```bash
# Upvote a card
curl -X POST "http://localhost:8080/v1/api/debates/votes" \
  -H "Content-Type: application/json" \
  -d '{"debate_card_id": 1, "vote_type": "upvote"}'

# Add an emoji reaction
curl -X POST "http://localhost:8080/v1/api/debates/votes" \
  -H "Content-Type: application/json" \
  -d '{"debate_card_id": 1, "vote_type": "emoji", "emoji": "🔥"}'
```

### Adding Comments

```bash
# Add a top-level comment
curl -X POST "http://localhost:8080/v1/api/debates/comments" \
  -H "Content-Type: application/json" \
  -d '{"debate_id": 1, "content": "This is a great debate topic!"}'

# Reply to a comment
curl -X POST "http://localhost:8080/v1/api/debates/comments" \
  -H "Content-Type: application/json" \
  -d '{"debate_id": 1, "parent_comment_id": 5, "content": "I agree with your point!"}'
```

## Future Enhancements

1. **Reddit Integration**: Add Reddit sentiment analysis
2. **Advanced Analytics**: More sophisticated engagement scoring
3. **User Authentication**: Proper user management for votes/comments
4. **Real-time Updates**: WebSocket support for live debate updates
5. **Moderation**: Content moderation for comments
6. **Notifications**: Alert users about new debates on their favorite teams
7. **Mobile App**: Native mobile application for better UX

## Testing

Run the debate system tests:

```bash
go test ./internal/api -v -run TestDebate
```

## Deployment

The debate system is integrated into the main API and will be deployed automatically with the rest of the application. Make sure to:

1. Run the database migration: `sql/schema/003_debates.sql`
2. Set up the required environment variables
3. Ensure Redis is available for caching
4. Configure OpenAI API access for AI features
