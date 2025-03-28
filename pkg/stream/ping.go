package stream

import "net/http"

// Не знаю, зачем, но хочется попробовать
type Available bool

func Ping() Available {
	var status Available

	_, err := http.Get("http://localhost:8081/")
	if err != nil {
		status = false
	} else {
		status = true
	}

	return status
}
