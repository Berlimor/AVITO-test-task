# AVITO-test-task
Тестовое задание на стажировку в Авито.
Запустить проект можно выполнением следующей команды в папке с проектом:
```
docker-compose up --build
```

# Тестирование проекта
Для начала, нужно создать сегмент / сегменты, которые будем добавлять пользователям. 
```
curl --location 'localhost:8000/segments' \
--header 'Content-Type: application/json' \
--data '{
    "name": "TEST_SEGMENT"
}'
```

После чего, мы можем добавить этот сегмент пользователю.
```
curl --location 'localhost:8000/users' \
--header 'Content-Type: application/json' \
--data '{
    "add-segments": ["TEST_SEGMENT"],
    "delete-segments": [],
    "user-id": 10
}'
```
Тело запроса состоит из трех строчек: сегменты, которые нужно добавить; сегменты, которые нужно удалить; id пользователя, которому добавляем / удаляем сегменты.

Вывести список всех пользователей.
```
curl --location 'localhost:8000/users'
```

Удалить сегменты.
```
curl --location --request DELETE 'localhost:8000/segments' \
--header 'Content-Type: application/json' \
--data '{
    "name": "TEST_SEGMENT"
}'
```

Вывести пользователя с указанным id.
```
curl --location 'localhost:8000/users/10'
```