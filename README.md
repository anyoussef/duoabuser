
# duoabuser

a league of legends tool that analyzes a player's recent ranked solo/duo match history and detects who they've been duoing with.

## what it does

- looks up a player by riot id (e.g. `Name#NA1`)
- fetches their last 20 ranked solo/duo games (queue 420)
- identifies frequent duo partners (teammates appearing in 2+ games)
- highlights duo partners in each individual match

## stack

- **backend** — go, single http endpoint
- **frontend** — vanilla html/css/js, two pages (`index.html`, `results.html`)
- **api** — riot games api (americas region, NA)

## setup

### prerequisites

- go 1.21+
- a riot api key from [developer.riotgames.com](https://developer.riotgames.com)

### backend

```bash
cd backend
cp .env.example .env   # or create .env manually
```

add your key to `backend/.env`:

```
RIOT_API_KEY=RGAPI-xxxx-xxxx-xxxx
```

install dependencies and run:

```bash
go mod tidy
go run .
```

server starts on `http://localhost:8080` by default. set `PORT` in `.env` to change it.

### frontend

from the `frontend/` directory, serve with any static file server:

```bash
cd frontend
python3 -m http.server 3000
```

then open `http://localhost:3000` in your browser.

## api

### `GET /api/duo?riotId=Name%23TAG`

the `#` in the riot id must be percent-encoded as `%23`.

**example:**
```
curl "http://localhost:8080/api/duo?riotId=Faker%23NA1"
```

**response:**
```json
{
  "targetRiotId": "Faker#NA1",
  "frequentPartners": [
    { "puuid": "...", "riotId": "Friend#NA1", "gamesPlayed": 7 }
  ],
  "matches": [
    {
      "matchId": "NA1_1234567",
      "gameDate": "Mar 25, 2026",
      "duration": "32m 14s",
      "win": true,
      "champion": "Azir",
      "kda": "5/2/10",
      "teammates": [...],
      "duoPartner": { "riotId": "Friend#NA1", ... }
    }
  ]
}
```

## notes

- dev riot api keys expire every 24 hours — regenerate at [developer.riotgames.com](https://developer.riotgames.com) if you get 401 errors
- dev keys are rate limited to 20 req/sec and 100 req/2min; the backend uses a concurrency semaphore to avoid bursting
- currently hardcoded to **NA** (americas routing). other regions would require changes to the api base urls in `backend/main.go`
