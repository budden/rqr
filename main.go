package main

// просьба = fetchTask

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/budden/rqr/pkg/errorcodes"
)

func handleRoot(w http.ResponseWriter, req *http.Request) {
	if return404IfExtraURLChars("/", w, req) {
		return
	}
	io.WriteString(w, `
<html>
<body><title>Requester</title>
<body>
<h1>Requester service.</h1>
<ul>
<li>Use POST /fetchtaskadd json urlencoded to add a request</li>
<li>Use POST /fetchTaskdel?id=requestId to delete a request</li>
</body>
</html>`)
}

func return404IfExtraURLChars(path string, w http.ResponseWriter, req *http.Request) (doReturn bool) {
	if strings.TrimPrefix(req.URL.Path, path) != "" {
		w.WriteHeader(http.StatusNotFound)
		doReturn = true
	}
	return
}

func main() {
	// FIXME disallow sub-urls for /
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/fetchtaskadd", handleRequestAdd)
	log.Fatal(http.ListenAndServe(":8086", nil))
}

// https://stackoverflow.com/a/15685432/9469533
// To test, use curl -i -X POST -d "[\"GET\", \"google.com\"]" http://localhost:8086/fetchTaskadd
// To test error reporting, remove the comma from JSON :)
func handleRequestAdd(w http.ResponseWriter, req *http.Request) {
	pt, err := convertJSONFetchTaskToParsedFetchTask(req)
	if reportFetchTaskErrorToClientIf(err, w) {
		return
	}
	et, err1 := executeFetchTask(pt)
	if reportFetchTaskErrorToClientIf(err1, w) {
		return
	}
	fetchTask := saveFetchTask(pt, et)
	fmt.Println(fetchTask)
}

func convertJSONFetchTaskToParsedFetchTask(req *http.Request) (pt *ParsedFetchTask, err error) {
	decoder := json.NewDecoder(req.Body)
	ji := jsonFetchTask{}
	err = decoder.Decode(&ji)
	// this is not an efficient way to check errors, but it saves lines of code :)

	if err != nil {
		err = newErrorWithCode(errorcodes.FailedToParsefetchTaskJSON, "Failed to parse request JSON data. Error is %#v", err)
		return
	}
	lenFetchTask := len(ji)
	if lenFetchTask != 2 && lenFetchTask != 4 {
		err = newErrorWithCode(errorcodes.FailedToParsefetchTaskJSON,
			"JSON fetchTask must be of the form [method, address] or of the form [method, address, headers, body]")
		return
	}
	pt = &ParsedFetchTask{Method: ji[0], URL: ji[1]}
	if lenFetchTask == 4 {
		pt.Headers = ji[2]
		pt.Body = ji[3]
	}
	return
}

/* Клиент просит сервис выполнить http запрос к некому ресурсу. В просьбе в формате json описаны поля {метод, адрес} (опционально: заголовки, тело). Например, {GET http://google.com}.
Сервис выполняет запрос из просьбы и в качестве ответа клиенту возвращает json объект с полями {сгенерированный id запроса, http статус, заголовки, длинна ответа}.
Список просьб должен сохраняться на сервере, например в map.
Выше описана операция создания просьбы (FetchFetchTask). Предусмотреть ещё операции получения всех существующих просьб (опционально постранично), операция удаления просьбы по id.
Задача предполагает, что кандидат покажет знание перечисленных выше пунктов за исключением, может быть, goroutine/chan/sync.Mutex. Так же мы хотели бы увидеть код приближённый к продакшн версии с понятными наименованиями переменных и http route-ов. Если кандидат уверен в своих силах, для выполнения просьб можно реализовать worker на goroutine, который бы получал задания из канала, выполнял их и безопасно в смысле многопоточности, сохранял результаты. */
