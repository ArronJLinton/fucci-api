# League Standings Endpoint

### GET /v1/api/futbol/league_standings

Returns the standings for a specific league and season.

#### Query Parameters

| Parameter | Type   | Required | Description                    |
| --------- | ------ | -------- | ------------------------------ |
| league_id | string | Yes      | The ID of the league to query  |
| season    | string | Yes      | The season year (e.g., "2024") |

#### Common League IDs

- 39: Premier League (England)
- 140: La Liga (Spain)
- 135: Serie A (Italy)
- 78: Bundesliga (Germany)
- 61: Ligue 1 (France)

#### Example Requests

1. Get Premier League standings for 2024:

```bash
curl "http://localhost:8080/v1/api/futbol/league_standings?league_id=39&season=2024"
```

2. Get La Liga standings for 2024:

```bash
curl "http://localhost:8080/v1/api/futbol/league_standings?league_id=140&season=2024"
```

3. Get Serie A standings for 2024:

```bash
curl "http://localhost:8080/v1/api/futbol/league_standings?league_id=135&season=2024"
```

#### Example Response

```json
{
  "get": "standings",
  "response": [
    {
      "league": {
        "id": 39,
        "name": "Premier League",
        "country": "England",
        "logo": "https://media.api-sports.io/football/leagues/39.png",
        "flag": "https://media.api-sports.io/flags/gb.svg",
        "season": 2024,
        "standings": [
          [
            {
              "rank": 1,
              "team": {
                "id": 40,
                "name": "Liverpool",
                "logo": "https://media.api-sports.io/football/teams/40.png"
              },
              "points": 84,
              "goalsDiff": 45,
              "group": "Premier League",
              "form": "DLDLW",
              "status": "same",
              "description": "Champions League",
              "all": {
                "played": 38,
                "win": 25,
                "draw": 9,
                "lose": 4,
                "goals": {
                  "for": 86,
                  "against": 41
                }
              },
              "home": {
                "played": 19,
                "win": 14,
                "draw": 4,
                "lose": 1,
                "goals": {
                  "for": 42,
                  "against": 16
                }
              },
              "away": {
                "played": 19,
                "win": 11,
                "draw": 5,
                "lose": 3,
                "goals": {
                  "for": 44,
                  "against": 25
                }
              },
              "update": "2025-05-26T00:00:00Z"
            }
          ]
        ]
      }
    }
  ],
  "errors": [],
  "results": 1,
  "paging": {
    "current": 1,
    "total": 1
  }
}
```

#### Error Responses

1. Missing league_id:

```json
{
  "error": "league_id is required"
}
```

2. Missing season:

```json
{
  "error": "season is required"
}
```

3. League not found:

```json
{
  "error": "No league standings found for the given league ID"
}
```

4. No standings data:

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
- The season parameter is required and should be specified as a 4-digit year (e.g., "2024")
- The response includes league information such as name, country, and logo
- The standings array may contain multiple groups (e.g., for leagues with multiple divisions)
