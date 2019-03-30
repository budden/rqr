package main

// просьба = fetchTask

import (
	"log"
	"net/http"

	"github.com/budden/rqr/pkg/errorcodes"
)

const (
	fetchTaskGetURL    = "/fetchtaskget/"
	fetchTaskDeleteURL = "/fetchtaskdelete/"
)

func main() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/fetchtaskadd", handleFetchTaskAdd)
	http.HandleFunc("/fetchtasklist", handleFetchTaskList)
	http.HandleFunc(fetchTaskGetURL, handleFetchTaskGet)
	http.HandleFunc(fetchTaskDeleteURL, handleFetchTaskDelete)
	log.Fatal(http.ListenAndServe(":8086", nil))
}

func handleRoot(w http.ResponseWriter, req *http.Request) {
	if checkNoExtraURLChars("/", w, req) || checkHTTPMethod("GET", w, req) {
		return
	}
	WriteReplyToResponseAsJSON(w, req, errorcodes.OK, []string{
		"Requester service.",
		"Use POST /fetchtaskadd json urlencoded to add a fetch task",
		"Use GET /fetchtaskget/ID to get a fetch task",
		"Use GET /fetchtasklist?offset=N&limit=N to get a list (both params are optional)",
		"Use POST /fetchtaskdelete/ID to delete a fetch task",
		"Use GET / to obtain this help",
		"Replies are always with Content-type = application/json"})
	return
}

/* Клиент просит сервис выполнить http запрос к некому ресурсу. В просьбе в формате json описаны поля {метод, адрес} (опционально: заголовки, тело). Например, {GET http://google.com}.
Сервис выполняет запрос из просьбы и в качестве ответа клиенту возвращает json объект с полями {сгенерированный id запроса, http статус, заголовки, длинна ответа}.
Список просьб должен сохраняться на сервере, например в map.
Выше описана операция создания просьбы (FetchFetchTask). Предусмотреть ещё операции получения всех существующих просьб (опционально постранично), операция удаления просьбы по id.
Задача предполагает, что кандидат покажет знание перечисленных выше пунктов за исключением, может быть, goroutine/chan/sync.Mutex. Так же мы хотели бы увидеть код приближённый к продакшн версии с понятными наименованиями переменных и http route-ов. Если кандидат уверен в своих силах, для выполнения просьб можно реализовать worker на goroutine, который бы получал задания из канала, выполнял их и безопасно в смысле многопоточности, сохранял результаты. */
