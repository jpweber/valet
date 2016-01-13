/*
* @Author: jamesweber
* @Date:   2015-12-17 14:02:09
* @Last Modified by:   jamesweber
* @Last Modified time: 2016-01-13 15:22:37
 */

package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func PrimaryHandler(w http.ResponseWriter, req *http.Request) {

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

	// pull out the particular appconfig we are using so we can modify values
	appConfig := appApis[apiCall]

	// authorize request if authorizing is required
	if appConfig.Authorize == true {
		if Authorize(appConfig, req) == false {
			// bad authorization
			w.WriteHeader(http.StatusForbidden)
			io.WriteString(w, "Authorization Failed")
			return
		}
	}

	// at this point we can assume that a request should succeed.
	// count the this request against stats and move on
	stats[appConfig.Name] += 1
	fmt.Printf("%+v", stats)

	// TODO: @debug
	fmt.Println("a")

	select {
	case <-appChans[appConfig.Name].Limiter:
		fmt.Println("Sending Request")
		start := time.Now()
		apiResponses := SendRequest(req, appConfig)
		stop := time.Since(start)
		// fmt.Println(stop) TODO: @debug

		jsonResponses := strings.Join(apiResponses, ",")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("X-Rate-Limit-Limit", fmt.Sprintf("%d", appConfig.LimitValue))
		w.Header().Add("X-Rate-Limit-Remaining", fmt.Sprintf("%d", len(appChans[appConfig.Name].Limiter)))
		w.Header().Add("X-Rate-Limit-Reset", fmt.Sprintf("%d", len(appChans[appConfig.Name].RateCountdown)))
		w.Header().Add("X-Backend-Response-Time", fmt.Sprintf("%s", stop))
		// io.WriteString(w, string(jsonResponses))
		io.WriteString(w, jsonResponses)

	default:
		w.Header().Add("X-Rate-Limit-Limit", fmt.Sprintf("%d", appConfig.LimitValue))
		w.Header().Add("X-Rate-Limit-Remaining", fmt.Sprintf("%d", len(appChans[appConfig.Name].Limiter)))
		w.Header().Add("X-Rate-Limit-Reset", fmt.Sprintf("%d", len(appChans[appConfig.Name].RateCountdown)))
		w.WriteHeader(420)

		return
	}

}
