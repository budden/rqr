# Введение

В задании ничего не сказано про тестирование. Существует пирамида тестирования, см. например, 
https://habr.com/ru/post/358950/. Из неё мы реализуем интеграционные тесты в виде сценария
ручного тестирования. Пример выполненного мной автоматизированного e2e теста 
можно найти тут https://github.com/budden/rlj/blob/master/pkg/leftjoin/integration_test.go или тут
https://github.com/budden/pgweb/blob/issue-281-by-budden-3/pkg/cli/appserver_test.go#L361 . Такой
тест организован путём добавления особого способа запуска приложения, при этом есть функция-прокладка 
под main(), которая позволяет запустить приложение в рабочем или в тестовом режиме. При тестировании 
вместе с приложением в том же процессе запускается и оснастка для тестирования, что позволяет
тестировать приложение с помощью самого себя и позволяет избежать написания сценариев на bash, 
а воспользоваться встроенной инфраструктурой тестирования для golang.

В этих же проектах есть и модульные тесты. Мы неформально опишем сценарий e2e тестирования с помощью Curl. 

# Тестовый сценарий

Собрать и запустить приложение:
```
# записать в $GOPATH/github.com/budden/rqr
# проект приватный, поэтому вы его не увидите
cd $GOPATH/github.com/budden/rqr
go get ./...
go generate ./...
go run main.go
```

В отдельном терминале выполнять команды:

Некорректный JSON
```
curl -X POST -d "[\"GET\" \"http://google.com/\"]" http://localhost:8086/fetchtaskadd
```

Ответ: 
```
{"Status":1,"Statustext":"FailedToParsefetchTaskJSON",
 "Contents":"Failed to parse request JSON data. Error is \u0026json.SyntaxError{msg:\"invalid character '\\\"' after array element\", Offset:8}"}
```

Некорректный JSON 2
```
curl -X POST -d "[\"GET\" \"http://google.com/\"" http://localhost:8086/fetchtaskadd
```
Ответ - аналогичный. 

Некорректный метод:
```
curl -X PUT -d "[]" http://localhost:8086/fetchtaskadd
```
Ответ: `{"Status":9,"Statustext":"IncorrectRequestMethod","Contents":null}`

Правильный запрос на добавление
```
curl -X POST -d "[\"GET\", \"http://google.com/\"]" http://localhost:8086/fetchtaskadd
```
Ответ:
```
{"Status":0,"Statustext":"NoError","Contents":{"ID":"1","Httpstatus":200,"Headers": ...,"BodyLength":14124}}
```

Этот запрос надо повторить ещё два раза, чтобы заполнить базу значениями. В задании не сказано о том, должны ли повторяющиеся идентичные запросы должны браться из кеша. Это можно было бы сделать, но это не обязательно будет правильно (ведь время идёт и содержимое веб-страниц меняется).

Неправильный запрос на получение просьбы
```
curl http://localhost:8086/fetchtaskget/2а
```
Ответ: `{"Status":8,"Statustext":"IncorrectIDFormat","Contents":"Incorrect id format"}`

Запрос на получение несуществующей просьбы 
```
curl http://localhost:8086/fetchtaskget/2800
```
Ответ: `{"Status":2,"Statustext":"FetchTaskNotFound","Contents":""}`

Неправильный запрос на получение списка просьб
```
curl "http://localhost:8086/fetchtasklist?offset=1m"
```
Ответ: `{"Status":11,"Statustext":"UnknownError","Contents":"strconv.Atoi: parsing \"1m\": invalid syntax"}`
Здесь есть небольшая недоработка - мы не проверяем, что в запросе нет неверных параметров. Такие параметры
просто игнорируются. Поэтому запрос с опечатками, например, `/fetchtasklist/offfffset=1`, пройдёт без 
ошибок, но значение offset будет не 1, а 0 - по умолчанию.

Правильный запрос на получение списка просьб:
```
curl "http://localhost:8086/fetchtasklist?limit=2&offset=1"
```
Ответ:
```
{"Status":0,"Statustext":"OK","Contents":{"Length":3,"Records":[{"ID":"2","Httpstatus":200,...},{"ID":"3",...}]}}
```

Правильный запрос на удаление просьбы (вообще-то нужно аналогично проверить и неправильный, но мы это уже показали
в запросе на добавление)
```
curl -X POST http://localhost:8086/fetchtaskdelete/2
```
Ответ: `{"Status":0,"Statustext":"OK","Contents":null}`

Повторяем тот же запрос - теперь это запрос на удаление несуществующей просьбы. Ответ:
`{"Status":2,"Statustext":"FetchTaskNotFound","Contents":""}`

Неверный URL
```
curl -X POST http://localhost:8086/non-existent
```

Ответ: `{"Status":10,"Statustext":"IncorrectURL","Contents":"POST to / to obtain a help on correct URLs"}`

Получение справки:
```
curl http://localhost:8086/
```
Ответ: `{"Status":0,"Statustext":"OK","Contents":["Requester service.",...]}`


