### Leaderboard API

This service provides GRPC API for storing player scores and listing them.

#### Prerequisites

* docker
* docker-compose
* [goose](https://github.com/pressly/goose) - migration tool

#### Run application

1. Build app image

`docker build -t leaderboard:latest -f build/Dockerfile .`

2. Start database and run migrations

```
cd deployments
docker-compose up -d postgres
cd ../db/migrations
goose postgres "host=localhost port=5432 user=user dbname=postgres password=option123 sslmode=disable" up
```

3. Apply seeders

```
cd ../seeders
goose -no-versioning postgres "host=localhost port=5434 user=user dbname=postgres password=option123 sslmode=disable" up
```

4. Run application

```
cd ../../deployments
docker-composer up -d app
```

#### GRPC

By default `9090` port is listening grpc calls. Can be changed in `deployments/docker-compose.yml` if needed.

Application provides two methods:
- ListScore - simple RPC
- SaveScore - bidirectional stream

All methods use token based authentication. Token is hardcoded: `secret-token`.

Proto file can be found here: `/internal/controller/protos`.

To change amount of results returned by `ListScore` change option in `configs/local.yml`, rebuild app and restart it.
