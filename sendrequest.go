/*
* @Author: jamesweber
* @Date:   2015-12-22 14:15:48
* @Last Modified by:   jamesweber
* @Last Modified time: 2016-01-13 15:31:07
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
	urlParts := strings.Split(req.URL.Path, "/")
	resultChan := make(chan string, len(appConfig.Endpoints))

	// fmt.Println(req.Method) // TODO: @debug

	for _, endpoint := range appConfig.Endpoints {

		go func(resultChan chan string, urlParts []string, endpoint Endpoint) {
			url := endpoint.Scheme + "://" + endpoint.Host

			// add port information if it was specified
			if endpoint.Port != 0 {
				url = url + ":" + fmt.Sprintf("%d", endpoint.Port)
			}
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
			client := &http.Client{}

			newReq, err := http.NewRequest(req.Method, url, nil)

			// add on the original request headers
			for key, value := range req.Header {
				newReq.Header.Add(key, value[0])
			}

			// fmt.Println(newReq.Header) // TODO: @debug
			response, err := client.Do(newReq)
			if err != nil {
				fmt.Printf("%s", err)
			} else {
				defer response.Body.Close()
				contents, err := ioutil.ReadAll(response.Body)
				if err != nil {
					fmt.Printf("%s", err)
					os.Exit(1)
				}
				resultChan <- string(contents)
			}
		}(resultChan, urlParts, endpoint)

	}

	// range over the result channel.
	// But don't do it in definitely break when we have qty of results
	// equal to the number of endpoints we tried
	iter := 0
	for result := range resultChan {
		// fmt.Println(result)
		results = append(results, result)
		iter++
		if iter == len(appConfig.Endpoints) {
			break
		}
	}
	close(resultChan)

	return results
}
