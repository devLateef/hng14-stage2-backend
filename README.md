# Insighta Labs — Demographic Intelligence API

A REST API built with Go (Gin + GORM + PostgreSQL) that provides advanced filtering, sorting, pagination, and natural language search over a dataset of 2026 demographic profiles.

---

## Tech Stack

- **Language**: Go
- **Framework**: Gin
- **ORM**: GORM
- **Database**: PostgreSQL
- **UUID**: google/uuid (v7)
- **Env loading**: godotenv

---

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL

### Setup

1. **Clone the repository**
   ```bash
   git clone <your-repo-url>
   cd insight-api
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Create a PostgreSQL database**
   ```sql
   CREATE DATABASE insight_db;
   ```

4. **Configure environment variables**

   Create a `.env` file in the project root:
   ```env
   DB_URL=postgres://postgres@localhost:5432/insight_db?sslmode=disable
   ```

5. **Add the seed file**

   Place the `profiles.json` file inside the `seed/` directory:
   ```
   seed/profiles.json
   ```

6. **Run the server**
   ```bash
   go run main.go
   ```

   On startup, the server will:
   - Connect to the database
   - Auto-migrate the `profiles` table
   - Seed all 2026 profiles (skips duplicates on re-run)
   - Start listening on port `8080`

---

## Database Schema

| Field                | Type         | Notes                          |
|----------------------|--------------|--------------------------------|
| id                   | UUID v7      | Primary key                    |
| name                 | VARCHAR      | Unique, person's full name     |
| gender               | VARCHAR      | `male` or `female`             |
| gender_probability   | FLOAT        | Confidence score               |
| age                  | INT          | Exact age                      |
| age_group            | VARCHAR      | `child`, `teenager`, `adult`, `senior` |
| country_id           | VARCHAR(2)   | ISO country code (e.g. `NG`)   |
| country_name         | VARCHAR      | Full country name              |
| country_probability  | FLOAT        | Confidence score               |
| created_at           | TIMESTAMP    | Auto-generated, UTC            |

---

## API Endpoints

### Base URL
```
http://localhost:8080
```

---

### 1. Get All Profiles

```
GET /api/profiles
```

Supports filtering, sorting, and pagination in a single request.

#### Filters

| Parameter               | Type   | Description                          |
|-------------------------|--------|--------------------------------------|
| `gender`                | string | `male` or `female`                   |
| `age_group`             | string | `child`, `teenager`, `adult`, `senior` |
| `country_id`            | string | ISO code e.g. `NG`, `KE`, `GH`      |
| `min_age`               | int    | Minimum age (inclusive)              |
| `max_age`               | int    | Maximum age (inclusive)              |
| `min_gender_probability`| float  | Minimum gender confidence score      |
| `min_country_probability`| float | Minimum country confidence score     |

#### Sorting

| Parameter | Values                                      |
|-----------|---------------------------------------------|
| `sort_by` | `age`, `created_at`, `gender_probability`   |
| `order`   | `asc`, `desc`                               |

#### Pagination

| Parameter | Default | Max | Description     |
|-----------|---------|-----|-----------------|
| `page`    | 1       | —   | Page number     |
| `limit`   | 10      | 50  | Results per page|

#### Example Request

```
GET /api/profiles?gender=male&country_id=NG&min_age=25&sort_by=age&order=desc&page=1&limit=10
```

#### Success Response (200)

```json
{
  "status": "success",
  "page": 1,
  "limit": 10,
  "total": 2026,
  "data": [
    {
      "id": "019746a1-1b2c-7def-8abc-0123456789ab",
      "name": "Emmanuel Okafor",
      "gender": "male",
      "gender_probability": 0.99,
      "age": 34,
      "age_group": "adult",
      "country_id": "NG",
      "country_name": "Nigeria",
      "country_probability": 0.85,
      "created_at": "2026-04-01T12:00:00Z"
    }
  ]
}
```

---

### 2. Natural Language Search

```
GET /api/profiles/search?q=<query>
```

Parses a plain English query and converts it into filters. Supports the same pagination parameters (`page`, `limit`) as the main endpoint.

#### Example Requests

```
GET /api/profiles/search?q=young males from nigeria
GET /api/profiles/search?q=adult females from kenya
GET /api/profiles/search?q=teenagers above 17
GET /api/profiles/search?q=males and females from ghana
GET /api/profiles/search?q=seniors from tanzania
```

#### Success Response (200)

```json
{
  "status": "success",
  "page": 1,
  "limit": 10,
  "total": 45,
  "data": [...]
}
```

#### Uninterpretable Query Response

```json
{
  "status": "error",
  "message": "Unable to interpret query"
}
```

---

## Natural Language Parsing Approach

The parser (`services/query_parser.go`) uses **rule-based keyword matching** — no AI or LLMs involved. It lowercases the query and scans for known keywords, mapping them to structured filters.

### Gender Keywords

| Query contains       | Filter applied         |
|----------------------|------------------------|
| `female` / `females` | `gender = female`      |
| `male` / `males`     | `gender = male`        |
| both `male` and `female` | no gender filter (returns all) |

> `female` is checked before `male` to avoid the substring collision (`"female"` contains `"male"`). If both are present, the gender filter is dropped entirely.

### Age Group Keywords

| Query contains              | Filter applied           |
|-----------------------------|--------------------------|
| `child` / `children`        | `age_group = child`      |
| `teenager` / `teen`         | `age_group = teenager`   |
| `adult`                     | `age_group = adult`      |
| `senior` / `elderly`        | `age_group = senior`     |

### Age Range Keywords

| Query pattern               | Filter applied                  |
|-----------------------------|---------------------------------|
| `young`                     | `min_age=16`, `max_age=24`      |
| `above X` / `over X`        | `min_age = X`                   |
| `below X` / `under X`       | `max_age = X`                   |
| `between X and Y`           | `min_age = X`, `max_age = Y`    |

> `young` maps to ages 16–24 for parsing purposes only. It is **not** a stored age group.

### Country Keywords

The parser matches full country names to ISO codes. Supported countries include:

| Query contains    | Filter applied     |
|-------------------|--------------------|
| `nigeria`         | `country_id = NG`  |
| `kenya`           | `country_id = KE`  |
| `ghana`           | `country_id = GH`  |
| `angola`          | `country_id = AO`  |
| `tanzania`        | `country_id = TZ`  |
| `ethiopia`        | `country_id = ET`  |
| `south africa`    | `country_id = ZA`  |
| `cameroon`        | `country_id = CM`  |
| `uganda`          | `country_id = UG`  |
| `senegal`         | `country_id = SN`  |
| ...and 30+ more   |                    |

### Example Mappings

| Query                              | Filters Applied                                      |
|------------------------------------|------------------------------------------------------|
| `young males`                      | `gender=male`, `min_age=16`, `max_age=24`            |
| `females above 30`                 | `gender=female`, `min_age=30`                        |
| `people from angola`               | `country_id=AO`                                      |
| `adult males from kenya`           | `gender=male`, `age_group=adult`, `country_id=KE`    |
| `male and female teenagers above 17` | `age_group=teenager`, `min_age=17`                 |
| `seniors from nigeria`             | `age_group=senior`, `country_id=NG`                  |
| `young females from ghana`         | `gender=female`, `min_age=16`, `max_age=24`, `country_id=GH` |

---

## Parser Limitations

- **No ISO code input** — queries like `"people from NG"` won't work; only full country names are supported
- **Single country per query** — `"people from nigeria and kenya"` will only match one country (whichever is found first in the map)
- **No OR logic for age groups** — `"children or adults"` is not supported; only one age group is matched per query
- **No negation** — `"not from nigeria"` or `"non-adults"` are not handled
- **Ambiguous country names** — `"guinea"` may match before `"equatorial guinea"` or `"guinea-bissau"` depending on map iteration order (Go maps are unordered)
- **No typo tolerance** — `"nigria"` or `"femal"` won't match anything
- **`young` overrides `above/below`** — if both `young` and `above X` appear in the same query, both min/max ages are set but `above X` will overwrite the `young` min_age
- **No compound age ranges** — `"between 20 and 30 or above 50"` is not supported

---

## Error Responses

All errors follow this structure:

```json
{
  "status": "error",
  "message": "<error message>"
}
```

| Status Code | Meaning                          |
|-------------|----------------------------------|
| 400         | Missing or empty parameter       |
| 422         | Invalid parameter type or value  |
| 404         | Profile not found                |
| 500         | Server error                     |

---

## CORS

All responses include:
```
Access-Control-Allow-Origin: *
```

---

## Project Structure

```
insight-api/
├── config/
│   └── db.go              # Database connection
├── handlers/
│   └── profile_handler.go # Route handlers
├── models/
│   └── profile.go         # Profile model
├── routes/
│   └── routes.go          # Route definitions + CORS
├── seed/
│   ├── seed.go            # Seeder logic
│   └── profiles.json      # 2026 profile records
├── services/
│   └── query_parser.go    # Natural language parser
├── utils/
│   └── pargination.go     # Utility package
├── .env                   # Environment variables (not committed)
├── .gitignore
├── go.mod
├── go.sum
├── main.go
└── README.md
```
