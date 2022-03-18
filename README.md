### Leaderboard API

This service provides GRPC API for storing player scores and listing them.

#### Prerequisites

* docker
* docker-compose
* [goose](https://github.com/pressly/goose) - migration tool

#### Run application

1. Build app image

```
docker build -t leaderboard:latest -f build/Dockerfile .
```

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
goose -no-versioning postgres "host=localhost port=5432 user=user dbname=postgres password=option123 sslmode=disable" up
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

```
request:
{
    "name": "user1",    (optional)  
    "page": 1,          (optional)
    "period": 0,        (0 - monthly records, 1 - all time records)
}
```

```
response:
{
    "results": [
        {
            "name": "user1",
            "score": 100,
            "rank": 1
        }
    ],
    "around_me": [
        {
            "name": "user2",
            "score": 50,
            "rank": 2
        }
    ],
    "page": 2
}
```

- SaveScore - bidirectional stream

```
request:
{
    "name": "user1", 
    "score": 100
}
```

```
response:
{
    "name": "user1",
    "rank": 1
}
```

All methods use token based authentication. Token is hardcoded: `secret-token`.

Proto file can be found here: `/internal/controller/protos`.

To change amount of results returned by `ListScore` update option in `configs/local.yml`, rebuild app and restart it.
