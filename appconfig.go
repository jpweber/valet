/*
* @Author: jamesweber
* @Date:   2015-12-17 13:41:45
* @Last Modified by:   jamesweber
* @Last Modified time: 2015-12-18 10:35:34
 */

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
}

// type AppConfigs struct {
// 	apps map[string]AppConfig
// }

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
