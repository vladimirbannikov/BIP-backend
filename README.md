# BIP-backend

### auth

make docker-run

Регистрируемся: curl -v -d '{"login":"user", "email":"email@gmail.com","password":"passwd"}' -H "Content-Type: application/json" -X POST localhost:9085/register > qr.png

Получаем qr код (тело ответа нужно конвертировать в png), сканируем в приложении Google Authentificator

Получаем токен(ttl = 24 часа): curl -v -d '{"login":"user","password":"passwd"}' -H "Content-Type: application/json" -X POST localhost:9085/login
Использовать полученный токен: "Auth: token"

Для фронта после логина нужно всегда запрашивать у клиента код двухступенчатой аутентификации: 
curl -v -d '{"code":"code"}' -H "Content-Type: application/json" -H "Auth: token" -X POST localhost:8080/user2fa

Если токен или 2fa просрочился, то будет возвращаться 401 -> идем в /login

Если нужно завершить сессию: curl -v -H "Auth: token" -X DELETE localhost:8080/logout
