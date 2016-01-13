/*
* @Author: jamesweber
* @Date:   2015-12-28 11:12:01
* @Last Modified by:   jpweber
* @Last Modified time: 2016-01-12 22:11:59
 */

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func AppList(w http.ResponseWriter) {

	names := make([]string, 0)
	for _, app := range appApis {
		names = append(names, app.Name)
	}

	list, _ := json.Marshal(names)
	w.Write(list)
}

func AppInfo(w http.ResponseWriter, name string) {
	// TODO: add something in here that if the auth header is passed
	// and a correct maching key return the limits for just that user/entity.
	// this is will not be useful until there is a per user/entity rate limiting.
	info, err := json.Marshal(appApis[name])
	if err != nil {
		fmt.Println(err)
	}
	w.Write(info)
}

func AppStats(w http.ResponseWriter, name string) {
	stats, err := json.Marshal(stats[name])
	if err != nil {
		fmt.Println(err)
	}
	w.Write(stats)
}
