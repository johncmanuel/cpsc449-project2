# cpsc449-project2

Project 2 for CPSC 449: Backend Engineering

## Team 

1. Tomas Oh (CWID: 885566877)
2. Jahn Tibayan (CWID: 884691718)
3. Nestor Reategui (CWID: 886727635)
4. John Carlo Manuel (CWID: 884779844)

## Project Overview

This project involves a backend service primarily written in Go that integrates with the CanvasAPI and utilizes an AI model (from OpenAI) to enhance student productivity. The service fetches assignments and syllabi from Canvas courses and uses AI to prioritize tasks based on deadlines and workload. The Go backend will employ a SQLite database along with Redis for caching mechanisms.

## Tech Stack

- Go: our programming language of choice
- Gin: http Web Go Framework, simplifies routing
- sqlc: to help generate DB boilerplate from a SQL schema
- Redis: utilized to cache DB queries
- OpenAI: AI to prioritize certain assignments/projects

## Project Structure

This project mainly consists of the `pkgs/` and `db/` directories:
- The `pkgs/` directory declares all the libraries that are used across the backend service including `canvas/`, `redis/`, and `utils/`. The `canvas/` library contains most of the logic that interacts with the CanvasAPI, specifically fetching resources (such as assignments and syllabus) across multiple Canvas courses. It utilizes the user's Canvas API key to retrieve those resources. The `redis/` package contains the functionality to interact with the Redis cache (improving performance of various queries).
- The `db/` directory combines the logic for generating SQL schemas and queries. It leverages the Go `sqlc` libary for generating Go structs used across the codebase.

## Routing

The following routes are described as follows:
- `/:courseID/assignments/:assignmentID`: Supports reading and deleting individual assignments based on their ID
- `/assignments`: Handles data insertion into the SQLite database
- `/all-assignments`: Retrieves all assignments from the database
- `/syllabus` (Concept): A POST request that leverages the OpenAI model to summarize the syllabus and store that in the database

## Caching

As mentioned before, we use Redis for caching. We cache any repeated queries to the database to improve performance. Below is a performance comparison between a query with and without Redis caching:

![Redis Performance]("./public/redis_performance.png")


## Future Works

- Most of the AI implementation is missing as of our latest update (12/10/24). We would like to utilize the AI to prioritize assignments based on estimated time and difficulty to help students organize their workflow.
- The service could also leverage a proper authentication system, not just utilizing the Canvas API as the way to identify each user.
- We could create a frontend to provide users with a UI to utilize the tool.

## Setup

1. Install Go
2. Get your Canvas token from your respective school Canvas site <https://{somedomainhere}.instructure.com/profile/settings>
3. Setup your .env file (see .env.example)

Optionally, install [air](https://github.com/air-verse/air) for hot reloading support. No need to do this step if you're installing through [Devbox](https://www.jetify.com/docs/devbox/installing_devbox/).

## How to Run

```bash
go run main.go
```

With air:
```bash
air
```

## Updating SQL files with sqlc

If you've made changes to any .sql files, run `go gen ./...` to generate .sql.go files, which contain generated Go code for directly interacting with the database.

