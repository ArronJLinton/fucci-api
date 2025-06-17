# League Standings Endpoint

### GET /v1/api/futbol/league_standings

Returns the current standings for a specific league.

#### Query Parameters

| Parameter | Type   | Required | Description                   |
| --------- | ------ | -------- | ----------------------------- |
| league_id | string | Yes      | The ID of the league to query |

#### Common League IDs

- 2021: Premier League (England)
- 2014: La Liga (Spain)
- 2019: Serie A (Italy)
- 2002: Bundesliga (Germany)
- 2015: Ligue 1 (France)

#### Example Requests

1. Get Premier League standings:

```bash
curl "http://localhost:8080/v1/api/futbol/league_standings?league_id=2021"
```

2. Get La Liga standings:

```bash
curl "http://localhost:8080/v1/api/futbol/league_standings?league_id=2014"
```

3. Get Serie A standings:

```bash
curl "http://localhost:8080/v1/api/futbol/league_standings?league_id=2019"
```

#### Example Response

```json
[
  {
    "rank": 1,
    "team": {
      "id": 50,
      "name": "Manchester City",
      "logo": "https://media.api-sports.io/football/teams/50.png"
    },
    "points": 82,
    "goalsDiff": 58,
    "form": "WWDLW",
    "all": {
      "played": 35,
      "win": 26,
      "draw": 4,
      "lose": 5,
      "goals": {
        "for": 89,
        "against": 31
      }
    },
    "home": {
      "played": 17,
      "win": 14,
      "draw": 2,
      "lose": 1,
      "goals": {
        "for": 48,
        "against": 15
      }
    },
    "away": {
      "played": 18,
      "win": 12,
      "draw": 2,
      "lose": 4,
      "goals": {
        "for": 41,
        "against": 16
      }
    }
  }
  // ... more teams
]
```

#### Error Responses

1. Missing league_id:

```json
{
  "error": "league_id is required"
}
```

2. League not found:

```json
{
  "error": "No league standings found for the given league ID"
}
```

3. No standings data:

```json
{
  "error": "No standings data available for this league"
}
```

#### Notes

- The standings are returned in order of rank (1st place first)
- Each team entry includes their current form (last 5 matches)
- Home and away statistics are provided separately
- The response includes detailed goal statistics for both home and away matches
