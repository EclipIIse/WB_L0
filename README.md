# Test Service

Для запуска необходимо:
- Redis
```
docker run --name redis -p 8888:6379 redis

Для использования redis-cli необходимо зайти в контейнер:
docker exec -ti redis bash

Затем уже запускать redis-cli
root@47d65f6e2aa7:/data# redis-cli
127.0.0.1:6379> get b563feb7b2b84b6test
```
- NATS
```
docker run -p 4222:4222 -p 8222:8222 -p 6222:6222 --name nats-server -ti nats-streaming
```
- PostgreSQL
```SQL
CREATE SCHEMA IF NOT EXISTS wb;
CREATE TABLE IF NOT EXISTS wb.orders
(
    id   varchar(500) PRIMARY KEY,
    data jsonb
);
```

### Переменные окружения
- HOST=localhost:8080;
- DB_HOST=127.0.0.1:5432;
- DB_NAME=mydb;
- DB_USER=root;
- DB_PASSWORD=mydbpass;
- NATS_SUBJECT_NAME=test-subj;
- NATS_CLUSTER_NAME=test-cluster;
- NATS_CLIENT_ID=test-client;
- REDIS_ADDRESS=127.0.0.1:8888;
- REDIS_RETRIES=10

### Порядок запуска
1. Зарускается сам сервис
2. Запускается pusher - утилита, которая пушит в канал NATS сообщения

> ВНИМАНИЕ! <br>
> Запись данных в базу проиcходит батчами, рамер бача задается в константной переменной **batchSize**,
> поэтому, следует учесть, что при отправке сообщений меньше, чем размер батча, сервис будет ждать ещё данных.