# Facts app
This is a very simple web server implemented on Golang 1.18 using 
standard library and database driver (PGX).
It can keep some random fact in its db, update it by request and
send back fact itself by an id (for more info [see OAI specification](openapi.yaml)).

## Requirements:
Docker/docker-compose

## Build:
```bash
$ docker-compose up
```

## Db schema

I use PostgreSQL as a database in this project

![Database schema infographics](assets/db_schema.png)
*Database schema illustration*


