# Debate System Documentation

## Overview

The debate system allows users to create, participate in, and manage AI-generated debates for football matches. The system supports both pre-match and post-match debates with comprehensive engagement tracking.

## Features

### Core Functionality

- **AI-Generated Debates**: Automatic generation of debate prompts using OpenAI
- **Pre/Post Match Debates**: Different debate types based on match status
- **Voting System**: Upvote, downvote, and emoji reactions
- **Comments**: Nested comment system for discussions
- **Analytics**: Engagement tracking and scoring
- **Soft Delete**: Safe deletion with recovery options

### Debate Types

- **Pre-match**: Generated before a match starts
- **Post-match**: Generated after a match finishes

## API Endpoints

### Debate Generation

- `GET /debates/generate` - Generate AI prompt only
- `POST /debates/generate` - Generate complete debate with cards

### Debate Management

- `POST /debates/` - Create manual debate
- `GET /debates/{id}` - Get specific debate
- `GET /debates/match` - Get debates by match ID
- `GET /debates/top` - Get top debates by engagement

### Soft Delete Management

- `DELETE /debates/{id}/hard` - Permanently delete debate (admin)
- `POST /debates/{id}/restore` - Restore soft-deleted debate

### Engagement

- `POST /debates/cards` - Create debate card
- `POST /debates/votes` - Vote on debate card
- `POST /debates/comments` - Add comment
- `GET /debates/{debateId}/comments` - Get comments

## Soft Delete System

The debate system implements a soft delete mechanism for data safety and recovery:

### How It Works

- **Soft Delete**: Sets `deleted_at` timestamp instead of removing data
- **Hard Delete**: Permanently removes data from database
- **Restore**: Clears `deleted_at` timestamp to reactivate debate

### Benefits

- **Data Recovery**: Accidental deletions can be undone
- **Audit Trail**: Track when debates were deleted
- **Referential Integrity**: Maintains relationships with comments/votes
- **Analytics**: Historical data preservation

### Usage Examples

```bash
# Soft delete a debate (default behavior)
curl -X POST /debates/generate \
  -H "Content-Type: application/json" \
  -d '{"match_id": "123", "debate_type": "pre_match", "force_regenerate": true}'

# Permanently delete a debate (admin only)
curl -X DELETE /debates/123/hard

# Restore a soft-deleted debate
curl -X POST /debates/123/restore
```

## Data Flow

1. **Debate Generation**: AI creates prompt → Debate created → Cards generated
2. **User Engagement**: Votes/comments → Analytics updated → Engagement score calculated
3. **Moderation**: Soft delete → Optional restore → Hard delete if needed

## Best Practices

### For Developers

- Always use soft delete for user-generated content
- Implement proper authentication for hard delete operations
- Monitor soft-deleted records for cleanup strategies
- Consider implementing automatic cleanup of old soft-deleted records

### For Administrators

- Use soft delete for temporary moderation
- Use hard delete only for permanent removal
- Regularly review soft-deleted debates for restoration opportunities
- Monitor storage usage of soft-deleted records

## Error Handling

The system includes comprehensive error handling:

- **Match Status Validation**: Prevents inappropriate debate generation
- **Data Validation**: Ensures proper debate structure
- **Recovery Options**: Soft delete allows for data recovery
- **User Feedback**: Clear error messages and status codes
