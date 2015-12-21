/*
* @Author: jamesweber
* @Date:   2015-12-17 14:02:09
* @Last Modified by:   jamesweber
* @Last Modified time: 2015-12-21 17:10:04
 */

package main

import (
	"fmt"
	"io"
	"net/http"
)

func PrimaryHandler(w http.ResponseWriter, req *http.Request) {

	// TODO: This should be moved to somewhere in main
	// we don't need to load up the configs and apis for
	// every request
	// configs := AppConfigList("conf")
	// appApis := LoadApps(configs)

	// pull out the api call from the http request
	apiCall := APICall(req)

	// devLimiter := make(chan bool, 2) // this is just for quick dev. Need to move this out of here
	// appApis[apiCall].Limiter = devLimiter

	select {
	case <-appApis[apiCall].Limiter:
		fmt.Println("Sending Request")
	default:
		w.WriteHeader(420)
		return
	}

	// see if we know about this application api
	if KnownApp(apiCall, appApis) == false {
		w.WriteHeader(http.StatusNotImplemented)
		io.WriteString(w, "API Not Implemented")
		return
	} else {
		// now check if the backend endpoints are available
		// TODO: function still needs implemented
	}

	// authorize request if authorizing is required
	if appApis[apiCall].Authorize == true {
		if Authorize(appApis[apiCall], req) == false {
			// bad authorization
			w.WriteHeader(http.StatusForbidden)
			io.WriteString(w, "Authorization Failed")
		}
	}

	// TODO: @debug
	fmt.Println("Would be returning api response here")

}
