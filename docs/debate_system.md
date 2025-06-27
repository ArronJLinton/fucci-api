# Debate System Documentation

## Overview

The debate system allows users to create and participate in football debates. Debates can be generated automatically using AI or created manually. Each debate consists of a headline, description, and multiple debate cards representing different stances.

## Endpoints

### 1. Generate AI Prompt (GET)

**Endpoint:** `GET /api/debates/generate`

**Query Parameters:**

- `match_id` (required): The football match ID
- `type` (required): Either "pre_match" or "post_match"

**Description:** Generates an AI prompt for debate creation without actually creating the debate in the database. Returns the prompt structure that can be used to create a debate manually.

**Example:**

```bash
curl "http://localhost:8080/api/debates/generate?match_id=1321727&type=pre_match"
```

**Response:**

```json
{
  "headline": "Will Manchester City dominate possession against Liverpool?",
  "description": "A heated debate about tactical approaches in this crucial match",
  "cards": [
    {
      "stance": "agree",
      "title": "City's possession game will be too much",
      "description": "Pep's tactical approach will control the tempo"
    },
    {
      "stance": "disagree",
      "title": "Liverpool's pressing will disrupt City",
      "description": "Klopp's gegenpressing will force turnovers"
    },
    {
      "stance": "wildcard",
      "title": "It depends on the referee's decisions",
      "description": "Key calls could swing momentum either way"
    }
  ]
}
```

### 2. Generate Complete Debate (POST)

**Endpoint:** `POST /api/debates/generate`

**Request Body:**

```json
{
  "match_id": "1321727",
  "debate_type": "pre_match",
  "force_regenerate": false
}
```

**Description:** Generates a complete debate with AI, creates it in the database, and returns the full debate structure with cards and analytics. This is the recommended endpoint for most use cases.

**Features:**

- Automatically creates debate in database
- Creates debate cards with AI-generated content
- Sets up analytics tracking
- Checks for existing debates to avoid duplicates
- Supports force regeneration with `force_regenerate: true`

**Example:**

```bash
curl -X POST "http://localhost:8080/api/debates/generate" \
  -H "Content-Type: application/json" \
  -d '{
    "match_id": "1321727",
    "debate_type": "pre_match",
    "force_regenerate": false
  }'
```

**Response:**

```json
{
  "message": "Debate generated successfully",
  "debate": {
    "id": 123,
    "match_id": "1321727",
    "debate_type": "pre_match",
    "headline": "Will Manchester City dominate possession against Liverpool?",
    "description": "A heated debate about tactical approaches in this crucial match",
    "ai_generated": true,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z",
    "cards": [
      {
        "id": 456,
        "debate_id": 123,
        "stance": "agree",
        "title": "City's possession game will be too much",
        "description": "Pep's tactical approach will control the tempo",
        "ai_generated": true,
        "created_at": "2024-01-15T10:30:00Z",
        "updated_at": "2024-01-15T10:30:00Z",
        "vote_counts": {
          "upvotes": 0,
          "downvotes": 0,
          "emojis": {}
        }
      }
    ],
    "analytics": {
      "id": 789,
      "debate_id": 123,
      "total_votes": 0,
      "total_comments": 0,
      "engagement_score": 0.0,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  }
}
```

### 3. Create Manual Debate (POST)

**Endpoint:** `POST /api/debates`

**Request Body:**

```json
{
  "match_id": "1321727",
  "debate_type": "pre_match",
  "headline": "Custom debate headline",
  "description": "Custom debate description",
  "ai_generated": false
}
```

### 4. Get Debate (GET)

**Endpoint:** `GET /api/debates/{id}`

**Description:** Retrieves a complete debate with all cards, vote counts, and analytics.

### 5. Get Debates by Match (GET)

**Endpoint:** `GET /api/debates/match?match_id={match_id}`

**Description:** Retrieves all debates for a specific match.

### 6. Get Top Debates (GET)

**Endpoint:** `GET /api/debates/top?limit={limit}`

**Description:** Retrieves the top debates based on engagement score.

## Debate Types

### Pre-Match Debates

- Generated before a match starts
- Focus on predictions, tactics, and expectations
- Uses lineup data and team statistics
- Incorporates recent news and social sentiment

### Post-Match Debates

- Generated after a match ends
- Focus on analysis, key moments, and outcomes
- Uses match statistics and results
- Incorporates social media reactions and controversy

## AI Data Sources

The AI prompt generator uses multiple data sources to create engaging debates:

1. **Match Information**: Basic fixture details from API-Football
2. **Lineup Data**: Starting XI and substitutes for both teams
3. **Match Statistics**: Goals, shots, possession, cards, etc.
4. **News Headlines**: Recent news about the teams
5. **Social Sentiment**: Twitter and Reddit sentiment analysis

## Caching

- AI prompts are cached for 24 hours to reduce API costs
- Use `force_regenerate: true` to bypass cache
- Debate responses are not cached to ensure real-time data

## Error Handling

Common error responses:

- `400 Bad Request`: Missing or invalid parameters
- `404 Not Found`: Match or debate not found
- `501 Not Implemented`: AI generation not configured
- `500 Internal Server Error`: Database or API errors

## Best Practices

1. **Use POST /generate for production**: Creates complete debates with proper database structure
2. **Use GET /generate for testing**: Quick way to preview AI prompts without database changes
3. **Check for existing debates**: The system automatically returns existing debates unless force regeneration is requested
4. **Handle errors gracefully**: Always check for error responses and provide fallback options
5. **Monitor analytics**: Track engagement scores to improve debate quality over time
