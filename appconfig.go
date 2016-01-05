/*
* @Author: jamesweber
* @Date:   2015-12-17 13:41:45
* @Last Modified by:   jpweber
* @Last Modified time: 2016-01-03 21:24:01
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
	AuthKey     string //`json:"-"`
	AuthHeader  string
	RateLimit   bool
	LimitValue  int64
	// Limiter     chan bool `json:"-"`
	Endpoints []Endpoint
	// RateCountdown chan bool `json:"-"`
	// Hits          chan bool `json:"-"`
}

type Endpoint struct {
	Host string
	Path string
	Port int64
}

type AppChans struct {
	RateCountdown chan bool
	Hits          chan bool
	Limiter       chan bool
}

const rateLimitDuration = 60

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

func refillBucket(app *AppConfig, ticker <-chan time.Time, channels AppChans) {
	for {
		select {
		case <-ticker:
			fmt.Println("ticker fired")
			cur := len(channels.Limiter)
			refill := app.LimitValue - int64(cur)
			var i int64
			// refill the limiter channel
			for i = 0; i < refill; i++ {
				channels.Limiter <- true
			}

			// refill the countdown channel
			for i = 0; i < rateLimitDuration; i++ {
				channels.RateCountdown <- true
			}
		}
	}

}

func countdown(timer <-chan time.Time, countdown chan bool) {
	for {
		select {
		case <-timer:
			// drain the countdown
			<-countdown
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

		configList[appConfig.Name] = appConfig

	}

	return configList

}

func InitChans(configs []string) map[string]AppChans {

	appChannels := map[string]AppChans{}
	newChannelSet := AppChans{}

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
			newChannelSet.Limiter = make(chan bool, appConfig.LimitValue)
			var i int64
			for i = 0; i < appConfig.LimitValue; i++ {
				newChannelSet.Limiter <- true
			}

			// create channel for ticks. Currently set at 1 minute ticks
			tickChan := time.NewTicker(time.Second * rateLimitDuration).C
			timerChan := time.NewTicker(time.Second * 1).C
			newChannelSet.RateCountdown = make(chan bool, 60)
			// fill the coundown channel with 60 items
			for i = 0; i > rateLimitDuration; i++ {
				newChannelSet.RateCountdown <- true
			}

			// fire off the frefill bucket function for this app
			go refillBucket(&appConfig, tickChan, newChannelSet)
			go countdown(timerChan, newChannelSet.RateCountdown)
		}

		// setup stats channel
		hits := make(chan bool)
		newChannelSet.Hits = hits

		appChannels[appConfig.Name] = newChannelSet

	}

	return appChannels

}
