/*
* @Author: jamesweber
* @Date:   2015-12-17 14:02:09
* @Last Modified by:   jpweber
* @Last Modified time: 2016-01-04 23:30:54
 */

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	// count the this request against rate limit and move on.

	// increment hit counter
	appChans[appConfig.Name].Hits <- true

	select {
	case <-appChans[appConfig.Name].Limiter:
		fmt.Println("Sending Request")
		start := time.Now()
		apiResponses := SendRequest(req, appConfig)
		stop := time.Since(start)
		fmt.Println(stop)
		jsonResponses, _ := json.Marshal(apiResponses)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("X-Rate-Limit-Limit", fmt.Sprintf("%d", appConfig.LimitValue))
		w.Header().Add("X-Rate-Limit-Remaining", fmt.Sprintf("%d", len(appChans[appConfig.Name].Limiter)))
		w.Header().Add("X-Rate-Limit-Reset", fmt.Sprintf("%d", len(appChans[appConfig.Name].RateCountdown)))
		w.Header().Add("X-Backend-Response-Time", fmt.Sprintf("%s", stop))
		io.WriteString(w, string(jsonResponses))

	default:
		w.Header().Add("X-Rate-Limit-Limit", fmt.Sprintf("%d", appConfig.LimitValue))
		w.Header().Add("X-Rate-Limit-Remaining", fmt.Sprintf("%d", len(appChans[appConfig.Name].Limiter)))
		w.Header().Add("X-Rate-Limit-Reset", fmt.Sprintf("%d", len(appChans[appConfig.Name].RateCountdown)))
		w.WriteHeader(420)

		return
	}

}
