package main

import "net/http"

// ParsedTask - это просьба в разобранном виде. Из ТЗ:
// В просьбе в формате json описаны поля {метод, адрес} (опционально: заголовки, тело). Например, {GET http://google.com}.
type ParsedTask struct {
	Method  string
	URL     string
	Headers string
	Body    string
}

// ExecutedTask - результат выполнения просьбы
// Сервис выполняет запрос из просьбы и в качестве ответа клиенту возвращает json объект
// с полями {сгенерированный id запроса, http статус, заголовки, длинна ответа}.
type ExecutedTask struct {
	Httpstatus int
	Headers    http.Header
	Bodylength int
}

// Task - просьба и результат её выполнения
type Task struct {
	ID string
	pt *ParsedTask
	et *ExecutedTask
}
