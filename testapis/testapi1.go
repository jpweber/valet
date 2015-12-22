package main

import (
	"io"
	"net/http"
)

func userauth(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	io.WriteString(w, `{"auth":"success","userid":20 }`)
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", userauth)
	http.ListenAndServe(":9050", mux)

}
