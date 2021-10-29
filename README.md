# Фибоначчи
Сервис, возвращающий срез последовательности чисел из ряда Фибоначчи.

Сервис должен отвечать на запросы и возвращать ответ. В ответе должны быть перечислены все числа, последовательности Фибоначчи с порядковыми номерами от x до y.

## Развертывание:
Для развертывания нужен docker.

Развертывание сервиса осуществляется командой:
````
make up
````
в директории с проектом

## Интеграционное тестирование:
интеграционные тесты сервиса с кешированием в памяти:
````
make setup-integration-memory-tests
make start-integration-memory-tests
make teardown-integration-memory-tests
````

интеграционные тесты сервиса с кешированием в Redis:
````
make setup-integration-redis-tests
make start-integration-redis-tests
make start-bench-redis-tests
make teardown-integration-redis-tests
````

## Тестирование HTTP API:
```
curl -X POST http://localhost:8888/api/v1/calc/fib/slice -H 'Content-Type: application/json' -d '{"x":2,"y":8}'
```
