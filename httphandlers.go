/*
* @Author: jamesweber
* @Date:   2015-12-17 14:02:09
* @Last Modified by:   jamesweber
* @Last Modified time: 2015-12-18 11:11:02
 */

package main

import (
	"io"
	"net/http"
)

func PrimaryHandler(w http.ResponseWriter, req *http.Request) {

	configs := AppConfigList("conf")
	appApis := LoadApps(configs)

	apiCall := APICall(req)

	if KnownApp(apiCall, appApis) == false {
		w.WriteHeader(http.StatusNotImplemented)
		io.WriteString(w, "API Not Implemented")
	}

}
