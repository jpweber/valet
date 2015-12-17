/*
* @Author: jamesweber
* @Date:   2015-12-16 16:47:12
* @Last Modified by:   jpweber
* @Last Modified time: 2015-12-16 22:18:28
 */

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"log/syslog"
	"net/http"
	"os"
)

const AppVersion = "0.0.1"

var buildNumber string

// wrapper function for http logging
func logger(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer log.Printf("%s - %s", r.Method, r.URL)
		fn(w, r)
	}
}
func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Looks like you got somewhere you didn't intend to.")

}

func ping(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "pong")
	Test()
}

func main() {

	versionPtr := flag.Bool("v", false, "a bool")
	// Once all flags are declared, call `flag.Parse()`
	// to execute the command-line parsing.
	flag.Parse()
	if *versionPtr == true {
		fmt.Println(AppVersion + " Build " + buildNumber)
		os.Exit(0)
	}

	logwriter, e := syslog.New(syslog.LOG_NOTICE, "DREMEL")
	if e == nil {
		log.SetOutput(logwriter)
	}

	log.SetFlags(0)
	log.Println("Dremel Starting")

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", ping)

	mux.HandleFunc("/", logger(APICall))
	http.ListenAndServe(":8000", mux)

}