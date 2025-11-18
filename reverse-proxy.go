package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	target, err := url.Parse("http://localhost:8000")
	if err != nil {
		log.Fatal("❌ Failed to parse target URL:", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Extracting certificate_id and certificate_hash from query parameters
		certificateID := r.URL.Query().Get("certificate_id")
		certificateHash := r.URL.Query().Get("certificate_hash")

		if certificateID != "" && certificateHash != "" {
			newPath := "/certificate/preview"
			newQuery := fmt.Sprintf("certificate_id=%s&certificate_hash=%s",
				certificateID, certificateHash)

			r.URL.Path = newPath
			r.URL.RawQuery = newQuery

			fmt.Printf("🔁 Forwarding request: %s %s -> %s%s?%s\n",
				r.Method, r.URL.Path, target, newPath, newQuery)
		} else {
			fmt.Printf("🔁 Forwarding request: %s %s -> %s%s\n",
				r.Method, r.URL.Path, target, r.URL.String())
		}

		proxy.ServeHTTP(w, r)
	})

	port := ":8010"
	fmt.Printf("🚀 Reverse Proxy running on http://localhost%s\n", port)
	fmt.Printf("🎯 Forwarding to: http://localhost:8000\n")

	log.Fatal(http.ListenAndServe(port, nil))
}
