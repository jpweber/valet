/*
* @Author: jamesweber
* @Date:   2015-12-17 14:02:09
* @Last Modified by:   jamesweber
* @Last Modified time: 2015-12-18 10:33:29
 */

package main

import (
	"fmt"
	"net/http"
)

func PrimaryHandler(w http.ResponseWriter, req *http.Request) {

	configs := AppConfigList("conf")
	appApis := LoadApps(configs)

	// fmt.Printf("%s", configs)
	fmt.Printf("%+v", appApis["userauth"])
	apiCall := APICall(req)
	println(apiCall)

}
