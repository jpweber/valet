/*
* @Author: jamesweber
* @Date:   2015-12-17 13:41:45
* @Last Modified by:   jamesweber
* @Last Modified time: 2015-12-21 17:17:21
 */

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

type AppConfig struct {
	Name        string
	Description string
	Authorize   bool
	AuthKey     string
	AuthHeader  string
	BackendHost string
	RateLimit   bool
	LimitValue  int64
	Limiter     chan bool
}

// read all files in the app config dir
func AppConfigList(dir string) []string {

	files, err := filepath.Glob(dir + "/*.json")
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	fileNames := make([]string, 0)

	for _, file := range files {
		fileNames = append(fileNames, file)
	}
	return fileNames
}

func refillBucket(app AppConfig, ticker <-chan time.Time) {
	for {
		select {
		case <-ticker:
			fmt.Println("ticker fired")
			var i int64
			for i = 0; i < app.LimitValue; i++ {
				app.Limiter <- true
			}
		}
	}

}

// build up AppConfigs struct with what we found
func LoadApps(configs []string) map[string]AppConfig {

	configList := map[string]AppConfig{}

	for _, config := range configs {
		file, err := ioutil.ReadFile(config)
		if err != nil {
			fmt.Printf("File error: %v\n", err)
			os.Exit(1)
		}

		appConfig := AppConfig{}
		if err := json.Unmarshal(file, &appConfig); err != nil {
			panic(err)
		}

		// populate the Limiter with appropriate tokens
		if appConfig.RateLimit == true {
			appConfig.Limiter = make(chan bool, appConfig.LimitValue)
			var i int64
			for i = 0; i < appConfig.LimitValue; i++ {
				appConfig.Limiter <- true
			}
		}

		// create channel for ticks. Currently set at 1 minut ticks
		tickChan := time.NewTicker(time.Minute * 1).C
		// fire off the frefill bucket function for this app
		go refillBucket(appConfig, tickChan)

		configList[appConfig.Name] = appConfig

	}

	return configList

}
