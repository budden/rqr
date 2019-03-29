package main

import "net/http"

// ParsedFetchTask - это просьба в разобранном виде. Из ТЗ:
// В просьбе в формате json описаны поля {метод, адрес} (опционально: заголовки, тело). Например, {GET http://google.com}.
type ParsedFetchTask struct {
	Method  string
	URL     string
	Headers string
	Body    string
}

// ExecutedFetchTask - результат выполнения просьбы
// Сервис выполняет запрос из просьбы и в качестве ответа клиенту возвращает json объект
// с полями {сгенерированный id запроса, http статус, заголовки, длинна ответа}.
type ExecutedFetchTask struct {
	Httpstatus int
	Headers    http.Header
	Bodylength int
}

// FetchTask - просьба и результат её выполнения
type FetchTask struct {
	ID string
	pt *ParsedFetchTask
	et *ExecutedFetchTask
}
