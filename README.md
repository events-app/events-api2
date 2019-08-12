# cards-api

- Start a PostgreSQL server running in a docker container
- Setup schema
- Seed database
- Connect to database from service

## Notes:


```
## Compile api and events-admin
go isntall ./...

## Start postgres:
docker-compose up -d

## Generate key
./cards-admin keygen private.pem

## Create the schema and insert some seed data.
./cards-admin --db-disable-tls=1 migrate 
./cards-admin --db-disable-tls=1 seed

## Run the app then make requests.
./cards-api --db-disable-tls=1
```

## Info
All the endpoint addresses can be found in the Postman collection or in cmd/cards-api/internal/handlers/routes.go