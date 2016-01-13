/*
* @Author: jamesweber
* @Date:   2015-12-16 16:47:12
* @Last Modified by:   jpweber
* @Last Modified time: 2016-01-12 21:30:14
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"log/syslog"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/davecheney/profile"
)

const AppVersion = "0.0.1"

// init main vars
var (
	buildNumber string
	configs     []string
	appApis     map[string]AppConfig
	appChans    map[string]AppChans
	stats       map[string]int64
	httpPort    string
)

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

func reload(w http.ResponseWriter, r *http.Request) {
	// BUG: this messes up channels and causes actuall api requrest to fail because the
	// channels no lnoger line up. Need to find a way to reload the config by just updating certain parts
	// so we don't blow away the channels of this instance of the config.
	configs = AppConfigList("conf")
	appApis = LoadApps(configs)
}

func saveNewApp(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in save new app")
	fmt.Printf("+%v", r.Body)
	decoder := json.NewDecoder(r.Body)
	var config AppConfig
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", config)

	// add to app configs
	appApis[config.Name] = config

	// add to app channels
	appChans[config.Name] = BuildChanSet(&config)

	// persist the configs to disk
	Persist(appApis)

}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	reqLen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	apiCall := strings.TrimSpace(string(buf[:reqLen]))

	fmt.Printf("%s", apiCall)
	<-appChans[apiCall].Limiter
	// Send a response back to person contacting us.
	conn.Write([]byte("Message received."))
	// Close the connection when you're done with it.
	conn.Close()
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	defer profile.Start(profile.CPUProfile).Stop()

	versionPtr := flag.Bool("v", false, "a bool")
	httpPortFlag := flag.String("P", "", "Port number for http server to listen on. Defaults to setting in server config")
	// Once all flags are declared, call `flag.Parse()`
	// to execute the command-line parsing.
	flag.Parse()
	if *versionPtr == true {
		fmt.Println(AppVersion + " Build " + buildNumber)
		os.Exit(0)
	}

	// read the server config
	serverConfig := LoadConfig()
	httpPort = serverConfig.Port

	// let command like port specification override the config file
	if *httpPortFlag != "" {
		httpPort = *httpPortFlag
	}

	logwriter, e := syslog.New(syslog.LOG_NOTICE, "VALET")
	if e == nil {
		log.SetOutput(logwriter)
	}

	log.SetFlags(0)
	log.Println("Valet Starting")

	// load up the configs
	workingDir, _ := os.Getwd()
	fmt.Println(workingDir)
	configs = AppConfigList(workingDir + "/conf")
	appApis = LoadApps(configs)
	appChans = InitChans(configs)

	// build up channels for metrics
	// TODO: create a channel per app. Then create a map [string]int64 with keys being app name
	// and value being int64s that will just get incremented every time the endpoint is hit.
	// stats = make(map[string]int64)
	// for _, appAPI := range appApis {
	// 	stats[appAPI.Name] = int64(0)

	// 	go func(appAPI AppConfig, stats map[string]int64) {
	// 		fmt.Println("Starting go func for ", appAPI.Name)
	// 		for {
	// 			select {
	// 			case <-appChans[appAPI.Name].Hits:
	// 				stats[appAPI.Name] += 1
	// 			default:
	// 				// do nothing. But we need this here so we don't block as we loop around.
	// 			}
	// 		}
	// 	}(appAPI, stats)
	// }

	go func() {
		// stuff for handling running pairs and keeping limits in sync
		listener, err := net.Listen("tcp", ":8001")
		if err != nil {
			// handle error
			fmt.Println(err)
		}
		for {
			conn, err := listener.Accept()
			if err != nil {
				// handle error
				fmt.Println(err)
			}
			handleRequest(conn)
			// Close the listener when the application closes.
			defer listener.Close()
		}
	}()

	// http server
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", ping)
	mux.HandleFunc("/apps", logger(apps))
	mux.HandleFunc("/apps/", logger(apps))
	mux.HandleFunc("/admin/reload", logger(reload))
	mux.HandleFunc("/new", logger(saveNewApp))
	mux.HandleFunc("/", logger(PrimaryHandler))
	http.ListenAndServe(":"+httpPort, mux)

}
