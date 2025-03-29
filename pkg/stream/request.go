package stream

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Chunk struct {
	Event   Event
	Content string
}

type Event string

var (
	Message     Event = "message"
	Information Event = "info"
	Source      Event = "source"
	Error       Event = "error"
)

func Request(query string) chan Chunk {
	result := make(chan Chunk)
	errors := make(chan error)
	request := make(map[string]string)
	request["query"] = query

	data, err := json.Marshal(request)
	if err != nil {
		result <- Chunk{Event: Error, Content: fmt.Sprintf("error while marshalling request: %v", err)}
		close(result)
		errors <- err
		return result
	}

	go func() {
		defer close(result)
		var buffer Chunk

		resp, err := http.Post("http://localhost:8081/api/v1/stream", "application/json", bytes.NewReader(data))
		if err != nil {
			errors <- err
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			errors <- fmt.Errorf("request failed: %d", resp.StatusCode)
			return
		}

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			lines := strings.Split(scanner.Text(), "\n")

			for _, line := range lines {
				if strings.HasPrefix(line, "event:") {
					// Мы получили новый event
					event := strings.TrimPrefix(line, "event:")
					switch {
					case event == "message":
						buffer.Event = Message
					case event == "info":
						buffer.Event = Information
					case event == "error":
						buffer.Event = Error
					default:
						buffer.Event = Error
						buffer.Content = "Unknown event type: " + event
					}

				} else if strings.HasPrefix(line, "data:") {
					// мы получили data для event
					data := strings.TrimPrefix(line, "data:")
					if buffer.Event == "" {
						buffer.Event = Error
						buffer.Content = "Empty event in the buffer with content " + buffer.Content
					}
					buffer.Content = data

					result <- buffer

					buffer = Chunk{}
				}
			}
		}

		if err := scanner.Err(); err != nil {
			errors <- err
			return
		}
	}()

	return result
}
