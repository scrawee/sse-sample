package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
)

type Dispatcher struct {
	name   string
	events chan *Events
}
type Events struct {
	Event string
}

func main() {
	app := fiber.New()
	app.Get("/sse", adaptor.HTTPHandler(handler(eventHandler)))
	app.Listen(":3000")
}

func handler(f http.HandlerFunc) http.Handler {
	return http.HandlerFunc(f)
}

func eventHandler(w http.ResponseWriter, r *http.Request) {
	dispatcher := &Dispatcher{name: r.RemoteAddr, events: make(chan *Events, 10)}
	go dispatchEvent(dispatcher)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	select {
	case ev := <-dispatcher.events:
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		enc.Encode(ev)
		fmt.Fprintf(w, "data: %v\n\n", buf.String())
		fmt.Printf("data: %v\n", buf.String())

	}

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

func dispatchEvent(dispatcher *Dispatcher) {
	for {
		db := &Events{
			Event: "Event GOLANG",
		}
		dispatcher.events <- db
	}
}
