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

// FetchTask - просьба и результат её выполнения, как она хранится в памяти
type FetchTask struct {
	ID string
	pt *ParsedFetchTask
	et *ExecutedFetchTask
}

// FetchTaskAsJSON represents a json format for fetch task when the task is sent
// to the client
type FetchTaskAsJSON struct {
	ID string
	// We store headers as map[string][]string, not as Headers, to avoid issues in case
	// the type http.Header would change in the future
	Httpstatus int
	Headers    map[string][]string
	BodyLength int
}

func convertFetchTaskToJSON(ft *FetchTask) *FetchTaskAsJSON {
	et := ft.et
	headers := map[string][]string(et.Headers)
	return &FetchTaskAsJSON{
		ID:         ft.ID,
		Httpstatus: et.Httpstatus,
		Headers:    headers,
		BodyLength: et.Bodylength}

}
