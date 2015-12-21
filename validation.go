/*
* @Author: jamesweber
* @Date:   2015-12-18 10:41:54
* @Last Modified by:   jamesweber
* @Last Modified time: 2015-12-21 11:52:43
 */

package main

import (
	"fmt"
	"net/http"
)

func KnownApp(appName string, configs map[string]AppConfig) bool {
	if _, ok := configs[appName]; ok {
		return true
	} else {
		return false
	}
}

// TODO: need to be implemented
func AppAvailable() bool {
	return true // TODO: @dev temp for develeopment
}

// authorize the http request
func Authorize(app AppConfig, req *http.Request) bool {
	if req.Header[app.AuthHeader][0] == app.AuthKey {
		// TODO: @debug
		fmt.Println("authorized")
		return true
	} else {
		return false
	}

}

// limit the requests
func RateLimit(app AppConfig) bool {
	return true
}
