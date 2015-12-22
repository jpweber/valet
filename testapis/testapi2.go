package main

import (
	"io"
	"net/http"
)

func username(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	io.WriteString(w, `{"userid":20, "username": "eaton"}`)
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", username)
	http.ListenAndServe(":9051", mux)

}
