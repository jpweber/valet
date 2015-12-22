/*
* @Author: jamesweber
* @Date:   2015-12-22 14:15:48
* @Last Modified by:   jamesweber
* @Last Modified time: 2015-12-22 15:35:20
 */

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func SendRequest(req *http.Request, appConfig AppConfig) []string {

	var results []string

	fmt.Println(req.URL.RawQuery)
	urlParts := strings.Split(req.URL.Path, "/")
	fmt.Println(urlParts[2:])

	for _, endpoint := range appConfig.Endpoints {
		url := "http://" + endpoint.Host
		url = url + ":" + fmt.Sprintf("%d", endpoint.Port)
		// append any extra path information if present
		if len(urlParts) > 1 {
			url = url + "/" + strings.Join(urlParts[2:], "/")
		}
		// append any query information if present
		if len(req.URL.RawQuery) > 0 {
			url = url + "?" + req.URL.RawQuery
		}

		// send the requests now.
		fmt.Println(url) // TODO: @debug
		response, err := http.Get(url)
		if err != nil {
			fmt.Printf("%s", err)
		} else {
			defer response.Body.Close()
			contents, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Printf("%s", err)
				os.Exit(1)
			}
			results = append(results, string(contents))
		}

	}

	return results
}
