/*
* @Author: jamesweber
* @Date:   2015-12-16 16:47:12
* @Last Modified by:   jamesweber
* @Last Modified time: 2015-12-28 12:13:27
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
	"strings"
)

const AppVersion = "0.0.1"

var buildNumber string

var configs []string
var appApis map[string]AppConfig
var stats map[string]int64

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
}

func apps(w http.ResponseWriter, r *http.Request) {
	urlParts := strings.Split(r.URL.Path, "/")
	fmt.Printf("%s", urlParts)
	fmt.Println(len(urlParts))
	if len(urlParts) == 2 {
		AppList(w)
	}

	if len(urlParts) == 3 {
		AppInfo(w, urlParts[2])
	}

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

	logwriter, e := syslog.New(syslog.LOG_NOTICE, "VALET")
	if e == nil {
		log.SetOutput(logwriter)
	}

	log.SetFlags(0)
	log.Println("Valet Starting")

	// load up the configs
	configs = AppConfigList("conf")
	appApis = LoadApps(configs)

	// build up channels for metrics
	// TODO: create a channel per app. Then create a map [string]int64 with keys being app name
	// and value being int64s that will just get incremented every time the endpoint is hit.
	stats = make(map[string]int64)
	for _, appAPI := range appApis {
		stats[appAPI.Name] = int64(0)

		go func(appAPI AppConfig, stats map[string]int64) {
			fmt.Println("Starting go func for ", appAPI.Name)
			for {
				select {
				case <-appAPI.Hits:
					stats[appAPI.Name] += 1
				default:
					// do nothing. But we need this here so we don't block as we loop around.
				}
			}
		}(appAPI, stats)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", ping)
	mux.HandleFunc("/apps", apps)
	mux.HandleFunc("/apps/", apps)
	mux.HandleFunc("/", logger(PrimaryHandler))
	http.ListenAndServe(":8000", mux)

}
