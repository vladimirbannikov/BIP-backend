# BIP-backend

### auth

make docker-run

Регистрируемся: curl -v -d '{"login":"user","password":"passwd"}' -H "Content-Type: application/json" -X POST localhost:8080/register

Получаем токен(ttl = 24 часа): curl -v -d '{"login":"user","password":"passwd"}' -H "Content-Type: application/json" -X POST localhost:8080/login

Использовать полученный токен: "Auth: token" 

Если нужно завершить сессию: curl -v -H "Auth: token" -X DELETE localhost:8080/logout

