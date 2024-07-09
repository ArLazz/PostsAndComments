# Post and Comments

## Описание

Это сервис для добавления и чтения постов и комментариев с использованием GraphQL, аналогичным комментариям к постам на популярных платформах, таких как Хабр или Reddit.

## Характеристики системы постов:
- Можно просмотреть список постов.
- Можно просмотреть пост и комментарии под ним.
- Пользователь, написавший пост, может запретить оставлять комментарии к своему посту.

## Характеристики системы комментариев к постам:
- Комментарии организованы иерархически, позволяя вложенность без ограничений.
- Длина текста комментария ограничена до 2000 символов.
- Система пагинации для получения списка комментариев.
- Комментарии к постам можно получать асинхронно, то есть клиенты, подписанные на определенный пост, получают уведомления о новых комментариях без необходимости повторного запроса.

## Реализация:
- Сервис написан на языке Golang.
- Сервер для api был кодосгенерирован с помощью утилиты  [gqlgen](https://github.com/99designs/gqlgen).
- Использован Docker для распространения сервиса в виде Docker-образа.
- Хранение данных может быть как в памяти (in-memory), так и в PostgreSQL. Выбор хранилища можно определить параметром при запуске сервиса. (В нашем случае, так как реализован Docker-образ, следует поменять флаг в Dockerfile - memory или postgres)
- Функционал покрыт unit-тестами.

## Использование
Для начала следует развернуть docker-образ, для этого напишите в командной строке:
```
make start
```
После запуска docker-контейнера для запросов к api, перейдите на http://localhost:8080/, там вы увидите ui для запросов в формате GraphQL.
### Создание поста 
Запрос для создания поста:
```
mutation {
  createPost(title: "3 Post", body: "This is the body of the first post", allowComments: false) {
    id
    title
    body
    allowComments
  }
}
```
где title - заголовок поста, body - текст поста, allowComments - булевый флаг, отвечающий за то, можно ли оставлять комментарии под постом

Вариант ответа:
```json
{
  "data": {
    "createPost": {
      "id": "1379c1bf-a5b8-4bfd-9f0d-ae5619d3169d",
      "title": "Test Post",
      "body": "This is the body of the first post",
      "allowComments": true
    }
  }
}
```

### Получение всех постов
Запрос для получения всех постов:
```
query {
  posts {
    id
    title
    body
    allowComments
  }
}
```

Вариант ответа:
```json
{
  "data": {
    "posts": [
      {
        "id": "1379c1bf-a5b8-4bfd-9f0d-ae5619d3169d",
        "title": "Test Post",
        "body": "This is the body of the first post",
        "allowComments": true
      },
      {
        "id": "c923cfd1-19a3-495a-aa88-5cd9d50b516e",
        "title": "Test Post",
        "body": "This is the body of the second post",
        "allowComments": true
      }
    ]
  }
}
```

### Создание комментария
Запрос для создание комментария:
```
mutation {
  createComment(postId: "1379c1bf-a5b8-4bfd-9f0d-ae5619d3169d", body: "It's comment for 1 post", parentId: null) {
    id
    postId
    body
    parentId
  }
}
```
где postId - ID поста, к которому пишут комментарий, body - текст комментария, parentId - ID родительского комментария(может быть null, если это первый комментарий к посту)

Вариант ответа:
```json
{
  "data": {
    "createComment": {
      "id": "51329828-dbee-438e-8de6-b802fc04bd50",
      "postId": "1379c1bf-a5b8-4bfd-9f0d-ae5619d3169d",
      "body": "It's comment for 1 post",
      "parentId": null
    }
  }
}
```

### Получение поста по ID с комментариями
Запрос для получения поста по ID с комментариями:
```
query {
  post(id: "1379c1bf-a5b8-4bfd-9f0d-ae5619d3169d", limit:2, offset:0) {
    id
    title
    body
    allowComments
    comments {
      postId
      id
      body
      parentId
    }
  }
}
```
где id - ID поста,  limit, offset - соответственно, количество комментариев и смещение их по списку. 
Пагинация сделана таким образом, что сначала будут выдаваться комментарии первого порядка вложенности, затем второго порядка и так далее(то есть сначала будут выданы комментарии, оставленные к посту, затем так же по порядку, комментарии, оставленные к этим комментариям и т.д.)

Вариант ответа:
```json
{
  "data": {
    "post": {
      "id": "1379c1bf-a5b8-4bfd-9f0d-ae5619d3169d",
      "title": "Test Post",
      "body": "This is the body of the first post",
      "allowComments": true,
      "comments": [
        {
          "postId": "1379c1bf-a5b8-4bfd-9f0d-ae5619d3169d",
          "id": "51329828-dbee-438e-8de6-b802fc04bd50",
          "body": "It's comment for 1 post",
          "parentId": null
        },
        {
          "postId": "1379c1bf-a5b8-4bfd-9f0d-ae5619d3169d",
          "id": "fefa41f0-bfba-43dd-bb8f-4a5fe98f77c9",
          "body": "It's comment for 1 post under 1 comment",
          "parentId": "51329828-dbee-438e-8de6-b802fc04bd50"
        }
      ]
    }
  }
}
```
### Подписка на получение комментариев

Клиент может подписаться на получение новых комментариев к посту:
```
subscription {
  commentAdded(postId: "1379c1bf-a5b8-4bfd-9f0d-ae5619d3169d") {
    id
    postId
    body
    parentId
    children {
      id
      body
    }
  }
}
```

Тогда, когда к посту добавится новый комментарий, клиент получит сообщение об этом:
```json
{
  "data": {
    "commentAdded": {
      "id": "ffda6bfe-cae9-4852-aab5-22debebda26c",
      "postId": "1379c1bf-a5b8-4bfd-9f0d-ae5619d3169d",
      "body": "It's comment for 1 post under 2 comment",
      "parentId": "51329828-dbee-438e-8de6-b802fc04bd50",
      "children": []
    }
  }
}
```

## Тесты

Функционал покрыт unit-тестами, для их запуска можно выполнить данную команду:

```
make tests
```
