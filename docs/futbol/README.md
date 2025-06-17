# Futbol Service API Documentation

This directory contains documentation for the Futbol Service API endpoints.

## Available Endpoints

1. [League Standings](league_standings.md) - Get current standings for a specific league
2. Team Standings - Get standings for a specific team
3. Matches - Get match information
4. Lineup - Get team lineup information
5. Leagues - Get available leagues

## Common League IDs

| League ID | League Name    | Country |
| --------- | -------------- | ------- |
| 2021      | Premier League | England |
| 2014      | La Liga        | Spain   |
| 2019      | Serie A        | Italy   |
| 2002      | Bundesliga     | Germany |
| 2015      | Ligue 1        | France  |

## Base URL

All endpoints are prefixed with:

```
http://localhost:8080/v1/api/futbol
```

## Authentication

All endpoints require a valid API key to be included in the request headers:

```
x-rapidapi-key: YOUR_API_KEY
```

## Rate Limiting

Please be aware of the API rate limits:

- 100 requests per day for the free tier
- 1000 requests per day for the pro tier

## Error Handling

All endpoints return appropriate HTTP status codes and error messages in the following format:

```json
{
  "error": "Error message description"
}
```

Common status codes:

- 200: Success
- 400: Bad Request
- 401: Unauthorized
- 404: Not Found
- 429: Too Many Requests
- 500: Internal Server Error
