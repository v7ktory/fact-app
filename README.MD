# Facts App
##### This Go application retrieves facts from the "API Ninjas" service, stores them in a MongoDB database, and provides an HTTP server to fetch and display these facts.

## Installation

```sh
docker run --name some-mongo -p 27017:27017 -d mongo
```

## Go dependencies
```sh
go get go.mongodb.org/mongo-driver/mongo
go get go.mongodb.org/mongo-driver/bson
```
## Set API Key
##### Replace "your api key" in the apiKey constant with your actual API key from https://api-ninjas.com/.

## Endpoint
* Endpoint: http://localhost:3000/facts
* Method: GET
* Returns: A JSON array containing facts stored in the MongoDB database.
