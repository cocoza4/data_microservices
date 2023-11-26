# Data Microservices

## Public API
Login API
```
<host>/v1/login
```

## Product API (Mongodb)
List all products
```
GET - <host>/v1/products
```

Get latest product
```
GET - <host>/v1/products/latest
```

List indexes
```
GET - <host>/v1/products/indexes
```

Create an index
```
POST - <host>/v1/products/indexes?field=column
```
where `column` is column name

## Kafka API
List Topics
```
GET - <host>/v1/kafka/topics
```

Create topic
```
POST - <host>/v1/kafka/topics/<name>
```
where `<name>` is topic name

Publish message
```
POST - <host>/v1/kafka/topics/<name>/publish
```
where `<name>` is topic name

Get latest message in a topic
```
GET - <host>/v1/kafka/topics/<name>/latest
```

## Dependencies
1. make
2. docker-compose 2.15.1
3. python
4. golang 1.20

## Test
```
make test
```

## Run
Export the following environment variables
```
export SECRET_USER=admin
export SECRET_PASSWORD=pass
export SECRET_KEY=abcd
export MONGO_USER=mongo_write
export MONGO_PASSWORD=write1243omf9dnsO
```
With `docker-compose`
```
make run
```

All examples can be found in `scripts` folder.