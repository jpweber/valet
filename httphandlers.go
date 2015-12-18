/*
* @Author: jamesweber
* @Date:   2015-12-17 14:02:09
* @Last Modified by:   jamesweber
* @Last Modified time: 2015-12-18 13:46:45
 */

package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

func PrimaryHandler(w http.ResponseWriter, req *http.Request) {

	// TODO: This should be moved to somewhere in main
	// we don't need to load up the configs and apis for
	// every request
	configs := AppConfigList("conf")
	appApis := LoadApps(configs)

	// pull out the api call from the http request
	apiCall := APICall(req)

	// see if we know about this application api
	if KnownApp(apiCall, appApis) == false {
		w.WriteHeader(http.StatusNotImplemented)
		io.WriteString(w, "API Not Implemented")
		return
	} else {
		// now check if the backend endpoints are available
		// TODO: function still needs implemented
	}

	var wg sync.WaitGroup
	wg.Add(2) // add two counters to wait group
	// authorize request of authorizing is required
	go func(wg *sync.WaitGroup, appApis map[string]AppConfig) {
		if appApis[apiCall].Authorize == true {
			if Authorize(appApis[apiCall], req) == false {
				// bad authorization
				w.WriteHeader(http.StatusForbidden)
				io.WriteString(w, "Authorization Failed")
				wg.Done()
			} else {
				// succesfull authorization
				wg.Done()
			}
		} else {
			// if we aren't authorizing just mark this wg counter to done
			wg.Done()
		}
	}(&wg, appApis)

	go func(wg *sync.WaitGroup, appApis map[string]AppConfig) {
		// check rate limits if request is rate limited.
		if appApis[apiCall].RateLimit == true {
			if RateLimit(appApis[apiCall]) == false {
				// if over limit reject the request with http status
				w.WriteHeader(http.StatusForbidden)
				io.WriteString(w, "Rate Limit Exceeded")
				wg.Done()
			} else {
				// not throttling
				wg.Done()
			}
		} else {
			// if we aren't authorizing just mark this wg counter to done
			wg.Done()
		}
	}(&wg, appApis)

	wg.Wait() //Wait for the concurrent routines to call 'done'

	// TODO: @debug
	fmt.Println("Would be returning api response here")

}
