# OMNI Backend Service

## Requirements

1. Go
2. Golang Migrate
3. Docker

## How to install:

1. Run `docker-compose up`
2. Create database `omni` in postgre (can use tableplus, cmd, etc)
3. To create database, run `migrate -path ./migration -database "postgresql://localhost:5433/omni?user=postgres&password=postgres&sslmode=disable" up`
4. To drop database, run `migrate -path ./migration -database "postgresql://localhost:5433/omni?user=postgres&password=postgres&sslmode=disable" drop`

## Kafka Command

## How to run:

1. Run `docker-compose up`
2. Check if db tables have been created (using TablePlus, etc)
3. Run `go run main.go` in another terminal
