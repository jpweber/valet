/*
* @Author: jpweber
* @Date:   2015-12-16 20:35:55
* @Last Modified by:   jpweber
* @Last Modified time: 2015-12-16 22:20:38
 */

package main

import (
	"fmt"
	"net/http"
	"strings"
)

func APICall(w http.ResponseWriter, req *http.Request) {

	urlParts := strings.Split(req.URL.Path, "/")

	var apiCall string
	if urlParts[1] == "" {
		// TODO: Need to return proper error here
		// as an http response
		fmt.Println("specific api not declared")
	} else {
		apiCall = urlParts[1]
	}

	// return apiCall
	fmt.Println(apiCall)
}
