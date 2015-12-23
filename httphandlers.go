/*
* @Author: jamesweber
* @Date:   2015-12-17 14:02:09
* @Last Modified by:   jamesweber
* @Last Modified time: 2015-12-23 15:23:14
 */

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	appConfig.Hits = 1

	fmt.Println(apiCall, appConfig.Hits) // TODO: @debug
	select {
	case <-appConfig.Limiter:
		fmt.Println("Sending Request")
		apiResponses := SendRequest(req, appConfig)
		jsonResponses, _ := json.Marshal(apiResponses)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("X-Rate-Limit-Limit", fmt.Sprintf("%d", appConfig.LimitValue))
		w.Header().Add("X-Rate-Limit-Remaining", fmt.Sprintf("%d", len(appConfig.Limiter)))
		w.Header().Add("X-Rate-Limit-Reset", fmt.Sprintf("%d", len(appConfig.RateCountdown)))
		io.WriteString(w, string(jsonResponses))

	default:
		w.Header().Add("X-Rate-Limit-Limit", fmt.Sprintf("%d", appConfig.LimitValue))
		w.Header().Add("X-Rate-Limit-Remaining", fmt.Sprintf("%d", len(appConfig.Limiter)))
		w.Header().Add("X-Rate-Limit-Reset", fmt.Sprintf("%d", len(appConfig.RateCountdown)))
		w.WriteHeader(420)

		return
	}

}
