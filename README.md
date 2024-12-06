# cpsc449-project2

## Team 

1. Tomas Oh (CWID: 885566877)
2. Jahn Tibayan (CWID: 884691718)
3. Nestor Reategui (CWID: 886727635)
4. John Carlo Manuel (CWID: 884779844)

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

