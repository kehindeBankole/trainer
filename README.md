# AI Workout Trainer

A backend API for a personal AI fitness trainer that uses video and audio models to coach users in real time.

---

## Setup

### 1. Start the database (Docker)

Make sure Docker Desktop is open and running first, then:

```bash
docker run --name workout-db \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=workout_trainer \
  -e POSTGRES_USER=workout_user \
  -p 5433:5432 \
  -d postgres
```

> Port 5433 is used because 5432 may already be taken by a local Postgres installation.

### 2. Set environment variables

Copy the `.env` file and update the values as needed:

```
DATABASE_URL="postgres://workout_user:password@localhost:5433/workout_trainer"
JWT_SECRET="your-secret-key-change-in-production"
```

Then load them in your terminal:

```bash
source .env
```

> You must run this in the same terminal you start the server from. It does not persist across sessions.

### 3. Run the migrations

```bash
psql $DATABASE_URL -f db/migrations/001_create_users.sql
psql $DATABASE_URL -f db/migrations/002_create_exercises.sql
```

### 4. Start the server

```bash
go run main.go
```

### 5. Test the health check

```bash
curl http://localhost:8080/health
```

---

## API Endpoints

### Users

| Method | Endpoint                  | Auth required | Description        |
|--------|---------------------------|---------------|--------------------|
| POST   | `/api/v1/users/register`  | No            | Register a new user |
| POST   | `/api/v1/users/login`     | No            | Login, returns JWT  |
| GET    | `/api/v1/users/`          | No            | List all users      |
| GET    | `/api/v1/users/{id}`      | No            | Get user by ID      |

### Exercises

| Method | Endpoint                  | Auth required | Description           |
|--------|---------------------------|---------------|-----------------------|
| POST   | `/api/v1/exercises/`      | Yes           | Create an exercise     |
| GET    | `/api/v1/exercises/`      | No            | List all exercises     |
| GET    | `/api/v1/exercises/{id}`  | No            | Get exercise by ID     |

### Authentication

Login returns a token:

```json
{
  "token": "eyJ...",
  "user": { "id": "...", "username": "...", "email": "...", "created_at": "..." }
}
```

Pass the token in the `Authorization` header for protected endpoints:

```
Authorization: Bearer <token>
```

---

## Viewing data locally

Connect to the database from VSCode (PostgreSQL extension) or any GUI tool (TablePlus, pgAdmin) using:

| Field    | Value           |
|----------|-----------------|
| Host     | localhost       |
| Port     | 5433            |
| User     | workout_user    |
| Password | password        |
| Database | workout_trainer |

---

## Common Issues

**Docker daemon not running**
```
Cannot connect to the Docker daemon at unix:///...docker.sock
```
Open Docker Desktop and wait for it to fully start before running any `docker` commands.

---

**Container name already in use**
```
docker: Error response from daemon: Conflict. The container name "/workout-db" is already in use
```
A previous container with the same name exists. Remove it and try again:
```bash
docker rm workout-db
docker run ...
```

---

**Port already in use**
```
Ports are not available: exposing port TCP 0.0.0.0:5432
```
Something else (usually a local Postgres install) is already on port 5432. Use 5433 instead as shown in step 1.

---

**DATABASE_URL or JWT_SECRET not set**
```
DATABASE_URL environment variable is not set
JWT_SECRET environment variable is not set
```
You opened a new terminal and the variables were lost. Run `source .env` again in the new terminal.

---

**psql not found**
```
zsh: command not found: psql
```
Install the Postgres CLI tools:
```bash
brew install libpq
brew link --force libpq
```
Then re-run the migration commands.

---

**Container already running, forgot to source .env**
If the server fails to connect to the database, check the container is still running:
```bash
docker ps | grep workout-db
```
If it is not listed, restart it:
```bash
docker start workout-db
```
Then run `source .env` and start the server again.
