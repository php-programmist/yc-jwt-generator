# Генератор IAM-токена для YandexCloud

Приложение запускается на 80-м порту.

Принимает через POST-запрос JSON с данными для генерации токена.

Пример запроса:
```http request
POST http://localhost
Content-Type: application/json

{
  "key_id": "abcdefghjkl",
  "service_account_id": "abcdefghjkl",
  "private_key": "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----"
}
```

Пример ответа:
```json
{
  "iamToken": "token-here",
  "expiresAt": "2023-02-21T00:55:09.980611102Z"
}
```

## Docker
Запуск контейнера:
```bash
docker run --rm -d -p 80:80 phpprogrammist/yc-jwt-generator
```

