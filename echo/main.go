package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		headersJson, _ := json.MarshalIndent(r.Header, "", "  ")
		w.Write([]byte("Headers:\n"))
		w.Write(headersJson)
		w.Write([]byte(fmt.Sprintf("\nHost: %s\n", r.Host)))
		w.Write([]byte(fmt.Sprintf("Method: %s\n", r.Method)))
		w.Write([]byte(fmt.Sprintf("Path: %s\n", r.URL)))
		w.Write([]byte(fmt.Sprintf("Request URI: %s\n", r.RequestURI)))
	})

	fmt.Printf("Echo running on port 8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}
