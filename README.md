# OMNI Backend Service

## Requirements

1. Go
2. Golang Migrate
3. Docker

## How to install:

1. Run `docker-compose up`
2. To create topic order, run `docker-compose exec broker bash -c "kafka-topics --create --replication-factor 1 --partitions 1 --topic order --bootstrap-server localhost:9092"`
3. To create topic product, run `docker-compose exec broker bash -c "kafka-topics --create --replication-factor 1 --partitions 1 --topic product --bootstrap-server localhost:9092"`
4. To list all available topics, run `docker-compose exec broker bash -c "kafka-topics --list --bootstrap-server localhost:9092"`
5. Create database `omni` in postgre (can use tableplus, cmd, etc)
6. To create database, run `migrate -path ./migration -database "postgresql://localhost:5433/omni?user=postgres&password=postgres&sslmode=disable" up`
7. To drop database, run `migrate -path ./migration -database "postgresql://localhost:5433/omni?user=postgres&password=postgres&sslmode=disable" drop`

## Kafka Command

## How to run:

1. Run `docker-compose up`
2. Check if kafka topics `product` and `order` are available (using the command in the installation)
3. Check if db tables have been created (using TablePlus, etc)
4. Run `go run main.go` in another terminal
