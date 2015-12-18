/*
* @Author: jpweber
* @Date:   2015-12-16 20:35:55
* @Last Modified by:   jamesweber
* @Last Modified time: 2015-12-18 10:39:09
 */

package main

import (
	"fmt"
	"net/http"
	"strings"
)

func APICall(req *http.Request) string {

	urlParts := strings.Split(req.URL.Path, "/")

	var apiCall string
	if urlParts[1] == "" {
		// TODO: Need to return proper error here
		// as an http response
		fmt.Println("specific api not declared")
	} else {
		apiCall = urlParts[1]
	}

	return apiCall
}
