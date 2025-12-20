Аутентификация 

Регистрация
POST `/auth/register`
```json
{
  "name": "Иван Иванов",
  "email": "ivan@mail.ru",
  "password": "123456"
}
```
**Response:**
```json
{
  "token": "eyJhbGciOiJ...",
  "user": {
    "id": 1,
    "name": "Иван Иванов",
    "email": "ivan@mail.ru",
    "password": ""
  }
}
```

POST /auth/login -- вход
{
  "email": "ivan@mail.ru",
  "password": "123456"
}
{
  "token": "eyJhbGciOiJ...",
  "user": {
    "id": 1,
    "name": "Иван Иванов",
    "email": "ivan@mail.ru",
    "password": ""
  }
}

### GET `/writers`
# Получить всех писателей (карточки)
**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
[
  {
    "id": 1,
    "name": "Александр Пушкин",
    "portrait_url": "https://example.com/pushkin.jpg",
    "tags": ["поэзия", "романтизм", "классика"],
    "is_favorite": true
  },
  {
    "id": 2,
    "name": "Лев Толстой",
    "portrait_url": "https://example.com/tolstoy.jpg",
    "tags": ["реализм", "роман", "философия"],
    "is_favorite": false
  }
]
```

### GET `/writers/{id}/bio`
# Получить биографию писателя
**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "id": 1,
  "name": "Александр Пушкин",
  "portrait_url": "https://example.com/pushkin.jpg",
  "tags": ["поэзия", "романтизм", "классика"],
  "lifespan": "1799-1837",
  "country": "Россия",
  "occupation": "Поэт, писатель",
  "is_favorite": true,
  "content": "Александр Сергеевич Пушкин — русский поэт, драматург и прозаик..."
}
```

### GET `/articles`
# Получить все статьи (карточки)
**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
[
  {
    "id": 1,
    "cover_url": "https://example.com/cover1.jpg",
    "title": "История русской литературы",
    "tags": ["литература", "история", "культура"],
    "description": "Обзор развития русской литературы...",
    "is_favorite": true
  },
  {
    "id": 2,
    "cover_url": "https://example.com/cover2.jpg",
    "title": "Серебряный век поэзии",
    "tags": ["поэзия", "история", "культура"],
    "description": "Исследование поэзии Серебряного века...",
    "is_favorite": false
  }
]
```

### GET `/articles/{id}`
# Получить статью
**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "id": 1,
  "cover_url": "https://example.com/cover1.jpg",
  "title": "История русской литературы",
  "tags": ["литература", "история", "культура"],
  "description": "Обзор развития русской литературы...",
  "content": "Русская литература имеет богатую историю...",
  "is_favorite": true
}
```

### GET `/favorites/writers`
# Получить избранных писателей
**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
[
  {
    "id": 1,
    "name": "Александр Пушкин",
    "portrait_url": "https://example.com/pushkin.jpg",
    "tags": ["поэзия", "романтизм", "классика"],
    "is_favorite": true
  }
]
```

### POST `/favorites/writers/{id}/toggle`
# Добавить/удалить писателя в избранное
**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Response:**
```json
{
  "message": "Писатель добавлен в избранное",
  "writer_id": 1,
  "is_favorite": true
}
```

### GET `/favorites/articles`
# Получить избранные статьи
**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
[
  {
    "id": 1,
    "cover_url": "https://example.com/cover1.jpg",
    "title": "История русской литературы",
    "tags": ["литература", "история", "культура"],
    "description": "Обзор развития русской литературы...",
    "is_favorite": true
  }
]
```

### POST `/favorites/articles/{id}/toggle`
# Добавить/удалить статью в избранное
**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Response:**
```json
{
  "message": "Статья добавлена в избранное",
  "article_id": 1,
  "is_favorite": true
}
```

### GET `/search`
# Общий поиск (статьи + писатели)
**Headers:**
```
Authorization: Bearer <token>
```

**Query Parameters:**
- `q` (optional) - поисковый запрос
- `tags` (optional) - фильтр по тегам через запятую

**Example:** `/search?q=литература&tags=история,культура`

**Response:**
```json
{
  "articles": [
    {
      "id": 1,
      "cover_url": "https://example.com/cover1.jpg",
      "title": "История русской литературы",
      "tags": ["литература", "история", "культура"],
      "description": "Обзор развития русской литературы...",
      "is_favorite": true
    }
  ],
  "writers": [
    {
      "id": 1,
      "name": "Александр Пушкин",
      "portrait_url": "https://example.com/pushkin.jpg",
      "tags": ["поэзия", "романтизм", "классика"],
      "is_favorite": true
    }
  ],
  "counts": {
    "articles": 1,
    "writers": 1,
    "total": 2
  },
  "query": "литература",
  "tags": ["история", "культура"]
}
```

### GET `/search/articles`
# Поиск только статей
**Headers:**
```
Authorization: Bearer <token>
```

**Query Parameters:**
- `q` (optional) - поисковый запрос
- `tags` (optional) - фильтр по тегам через запятую

**Response:**
```json
{
  "articles": [
    {
      "id": 1,
      "cover_url": "https://example.com/cover1.jpg",
      "title": "История русской литературы",
      "tags": ["литература", "история", "культура"],
      "description": "Обзор развития русской литературы...",
      "is_favorite": true
    }
  ],
  "count": 1,
  "query": "литература",
  "tags": ["история"]
}
```

### GET `/search/writers`
# Поиск только писателей
**Headers:**
```
Authorization: Bearer <token>
```

**Query Parameters:**
- `q` (optional) - поисковый запрос
- `tags` (optional) - фильтр по тегам через запятую

**Response:**
```json
{
  "writers": [
    {
      "id": 1,
      "name": "Александр Пушкин",
      "portrait_url": "https://example.com/pushkin.jpg",
      "tags": ["поэзия", "романтизм", "классика"],
      "is_favorite": true
    }
  ],
  "count": 1,
  "query": "пушкин",
  "tags": ["поэзия"]
}
```

### GET `/search/tags`
# Получить все уникальные теги
**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "tags": ["поэзия", "литература", "история", "культура", "философия", "роман", "драма"],
  "count": 7
}
```

### GET `/user/reading-progress`
# Получить прогресс чтения
**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "writers": {
    "1": 0.8,
    "2": 1.0
  },
  "articles": {
    "201": 0.5
  }
}
```

### POST `/user/reading-progress`
# Сохранить прогресс для одного элемента
**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request:**
```json
{
  "type": "writer",
  "id": 1,
  "progress": 0.8
}
```

**Response:**
```json
{
  "message": "Прогресс сохранен",
  "type": "writer",
  "id": 1,
  "progress": 0.8
}
```

### POST `/user/reading-progress/bulk`
# Массовое сохранение прогресса
**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request:**
```json
{
  "writers": {
    "101": 1.0,
    "102": 0.85,
    "103": 0.0
  },
  "articles": {
    "201": 1.0,
    "202": 0.95
  }
}
```

**Response:**
```json
{
  "message": "Прогресс массово сохранен",
  "counts": {
    "writers": 3,
    "articles": 2
  }
}
```

---

## Error Responses
```json
{
  "error": "Текст ошибки",
  "details": {
    "email": "Неверный формат email",
    "password": "Пароль должен содержать минимум 6 символов"
  }
}
```

**HTTP Status Codes:**
- `400` - Bad Request (неверные данные)
- `401` - Unauthorized (не авторизован)
- `404` - Not Found (не найдено)
- `409` - Conflict (конфликт, например email уже существует)
- `500` - Internal Server Error
```