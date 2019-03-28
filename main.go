package main

// просьба = inquiry

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
<li>Use POST /inquiryadd json urlencoded to add a request</li>
<li>Use POST /inquirydel?id=requestId to delete a request</li>
</body>
</html>`)
}

func main() {
	// FIXME process
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/inquiryadd", handleRequestAdd)
	log.Fatal(http.ListenAndServe(":8086", nil))
}

// ParsedInquiry - это просьба в разобранном виде. Из ТЗ:
// В просьбе в формате json описаны поля {метод, адрес} (опционально: заголовки, тело). Например, {GET http://google.com}.
type ParsedInquiry struct {
	Method  string
	URL     string
	Headers string
	Body    string
}

const errorCodeFailedToParseInquiry = 1

// https://stackoverflow.com/a/15685432/9469533
// To test, use curl -X POST -d "[\"GET\", \"google.com\"]" http://localhost:8086/inquiryadd
func handleRequestAdd(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var t []string
	err := decoder.Decode(&t)
	if err != nil {
		encoder := json.NewEncoder(w)
		reply := []interface{}{errorcodes.FailedToParseInquiryJson, fmt.Sprintf("Failed to parse request JSON data. Error is %#v", err)}
		err = encoder.Encode(reply)
		if err != nil {
			log.Printf("Error while sending error response to a client: %#v\n", err)
		}
		return
	}

}

/* Клиент просит сервис выполнить http запрос к некому ресурсу. В просьбе в формате json описаны поля {метод, адрес} (опционально: заголовки, тело). Например, {GET http://google.com}.
Сервис выполняет запрос из просьбы и в качестве ответа клиенту возвращает json объект с полями {сгенерированный id запроса, http статус, заголовки, длинна ответа}.
Список просьб должен сохраняться на сервере, например в map.
Выше описана операция создания просьбы (FetchTask). Предусмотреть ещё операции получения всех существующих просьб (опционально постранично), операция удаления просьбы по id.
Задача предполагает, что кандидат покажет знание перечисленных выше пунктов за исключением, может быть, goroutine/chan/sync.Mutex. Так же мы хотели бы увидеть код приближённый к продакшн версии с понятными наименованиями переменных и http route-ов. Если кандидат уверен в своих силах, для выполнения просьб можно реализовать worker на goroutine, который бы получал задания из канала, выполнял их и безопасно в смысле многопоточности, сохранял результаты. */
