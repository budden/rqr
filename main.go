package main

// просьба = fetchTask

import (
	"io"
	"log"
	"net/http"
)

func main() {
	// FIXME disallow sub-urls for /
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/fetchtaskadd", handleFetchTaskAdd)
	http.HandleFunc("/fetchtasklist", handleFetchTaskList)
	log.Fatal(http.ListenAndServe(":8086", nil))
}

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

func handleFetchTaskList(w http.ResponseWriter, req *http.Request) {

}

/* Клиент просит сервис выполнить http запрос к некому ресурсу. В просьбе в формате json описаны поля {метод, адрес} (опционально: заголовки, тело). Например, {GET http://google.com}.
Сервис выполняет запрос из просьбы и в качестве ответа клиенту возвращает json объект с полями {сгенерированный id запроса, http статус, заголовки, длинна ответа}.
Список просьб должен сохраняться на сервере, например в map.
Выше описана операция создания просьбы (FetchFetchTask). Предусмотреть ещё операции получения всех существующих просьб (опционально постранично), операция удаления просьбы по id.
Задача предполагает, что кандидат покажет знание перечисленных выше пунктов за исключением, может быть, goroutine/chan/sync.Mutex. Так же мы хотели бы увидеть код приближённый к продакшн версии с понятными наименованиями переменных и http route-ов. Если кандидат уверен в своих силах, для выполнения просьб можно реализовать worker на goroutine, который бы получал задания из канала, выполнял их и безопасно в смысле многопоточности, сохранял результаты. */
