package chatter

import (
	"fmt"
	"net/http"
	"sync"
)

// Handler serves the HTTP requests for our application.
type Handler struct {
	mu          sync.Mutex
	messages    []string
	connections []chan string
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		w.Write([]byte(indexHTML))
	case "/messages":
		if r.Method == "GET" {
			h.handleMessages(w, r)
		} else if r.Method == "POST" {
			h.handleCreateMessage(w, r)
		} else {
			http.Error(w, "status method not allowed", http.StatusMethodNotAllowed)
		}
	default:
		http.NotFound(w, r)
	}
}

func (h *Handler) handleMessages(w http.ResponseWriter, r *http.Request) {
	closeNotify := w.(http.CloseNotifier).CloseNotify()

	// Set headers.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	h.mu.Lock()
	for _, msg := range h.messages {
		writeMessage(w, msg)
	}

	ch := make(chan string)
	h.connections = append(h.connections, ch)
	h.mu.Unlock()

	for {
		select {
		case <-closeNotify:
			h.mu.Lock()
			for i, connection := range h.connections {
				if connection == ch {
					h.connections = append(h.connections[:i], h.connections[i+1:]...)
					break
				}
			}
			h.mu.Unlock()
			return
		case msg := <-ch:
			writeMessage(w, msg)
		}
	}
}

func writeMessage(w http.ResponseWriter, msg string) {
	fmt.Fprintf(w, "data: %s\n", msg)
	fmt.Fprint(w, "\n")
	w.(http.Flusher).Flush()
}

func (h *Handler) handleCreateMessage(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()

	body := r.FormValue("body")
	h.messages = append(h.messages, body)
	for _, ch := range h.connections {
		ch <- body
	}
}

var indexHTML = `
<html>
<body>
  <div id="messages"></div>
</body>
<script>
  var source = new EventSource("/messages");
  source.onmessage = function(e) {
    var elem = document.createElement("p");
    elem.innerHTML = e.data;
    messages.appendChild(elem);
  }
</script>
</html>
`
