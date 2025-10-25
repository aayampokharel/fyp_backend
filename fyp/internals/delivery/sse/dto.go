package sse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project/package/enum"
)

func HandleSSEResponse(data interface{}, event enum.SSETYPE, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	payload, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error encoding SSE message:", err)
		return
	}

	fmt.Fprintf(w, "event: %s\n", event)
	fmt.Fprintf(w, "data: %s\n\n", payload)

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}
