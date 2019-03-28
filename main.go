package main

// просьба = task

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/budden/rqr/pkg/errorcodes"
)

func handleRoot(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, `
<html>
<body><title>Requester</title>
<body>
<h1>Requester service.</h1>
<ul>
<li>Use POST /taskadd json urlencoded to add a request</li>
<li>Use POST /taskdel?id=requestId to delete a request</li>
</body>
</html>`)
}

func main() {
	// FIXME disallow sub-urls for /
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/taskadd", handleRequestAdd)
	log.Fatal(http.ListenAndServe(":8086", nil))
}

// https://stackoverflow.com/a/15685432/9469533
// To test, use curl -X POST -d "[\"GET\", \"google.com\"]" http://localhost:8086/taskadd
func handleRequestAdd(w http.ResponseWriter, req *http.Request) {
	pt, err := convertJSONTaskToParsedTask(req)
	if reportTaskErrorToClientIf(err, w) {
		return
	}
	et, err1 := executeTask(pt)
	if reportTaskErrorToClientIf(err1, w) {
		return
	}
	task := saveTask(pt, et)
	fmt.Println(task)
}

func convertJSONTaskToParsedTask(req *http.Request) (pt *ParsedTask, err error) {
	decoder := json.NewDecoder(req.Body)
	ji := jsonTask{}
	err = decoder.Decode(&ji)
	// this is not an efficient way to check errors, but it saves lines of code :)

	if err != nil {
		err = newErrorWithCode(errorcodes.FailedToParsetaskJSON, "Failed to parse request JSON data. Error is %#v", err)
		return
	}
	lenTask := len(ji)
	if lenTask != 2 && lenTask != 4 {
		err = newErrorWithCode(errorcodes.FailedToParsetaskJSON,
			"JSON task must be of the form [method, address] or of the form [method, address, headers, body]")
		return
	}
	pt = &ParsedTask{Method: ji[0], URL: ji[1]}
	if lenTask == 4 {
		pt.Headers = ji[2]
		pt.Body = ji[3]
	}
	return
}

/* Клиент просит сервис выполнить http запрос к некому ресурсу. В просьбе в формате json описаны поля {метод, адрес} (опционально: заголовки, тело). Например, {GET http://google.com}.
Сервис выполняет запрос из просьбы и в качестве ответа клиенту возвращает json объект с полями {сгенерированный id запроса, http статус, заголовки, длинна ответа}.
Список просьб должен сохраняться на сервере, например в map.
Выше описана операция создания просьбы (FetchTask). Предусмотреть ещё операции получения всех существующих просьб (опционально постранично), операция удаления просьбы по id.
Задача предполагает, что кандидат покажет знание перечисленных выше пунктов за исключением, может быть, goroutine/chan/sync.Mutex. Так же мы хотели бы увидеть код приближённый к продакшн версии с понятными наименованиями переменных и http route-ов. Если кандидат уверен в своих силах, для выполнения просьб можно реализовать worker на goroutine, который бы получал задания из канала, выполнял их и безопасно в смысле многопоточности, сохранял результаты. */
