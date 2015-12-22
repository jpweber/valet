/*
* @Author: jamesweber
* @Date:   2015-12-17 14:02:09
* @Last Modified by:   jamesweber
* @Last Modified time: 2015-12-22 16:39:57
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

	// authorize request if authorizing is required
	if appApis[apiCall].Authorize == true {
		if Authorize(appApis[apiCall], req) == false {
			// bad authorization
			w.WriteHeader(http.StatusForbidden)
			io.WriteString(w, "Authorization Failed")
			return
		}
	}

	// at this point we can assume that a request should succeed.
	// count the this request against rate limit and move on.
	select {
	case <-appApis[apiCall].Limiter:
		fmt.Println("Sending Request")
		apiResponses := SendRequest(req, appApis[apiCall])
		jsonResponses, _ := json.Marshal(apiResponses)
		w.Header().Add("X-Rate-Limit-Limit", fmt.Sprintf("%d", appApis[apiCall].LimitValue))
		w.Header().Add("X-Rate-Limit-Remaining", fmt.Sprintf("%d", len(appApis[apiCall].Limiter)))
		w.Header().Add("X-Rate-Limit-Reset", fmt.Sprintf("%d", len(appApis[apiCall].RateCountdown)))
		io.WriteString(w, string(jsonResponses))

	default:
		w.Header().Add("X-Rate-Limit-Limit", fmt.Sprintf("%d", appApis[apiCall].LimitValue))
		w.Header().Add("X-Rate-Limit-Remaining", fmt.Sprintf("%d", len(appApis[apiCall].Limiter)))
		w.Header().Add("X-Rate-Limit-Reset", fmt.Sprintf("%d", len(appApis[apiCall].RateCountdown)))
		w.WriteHeader(420)

		return
	}

	// TODO: @debug
	fmt.Println("Would be returning api response here")

}
