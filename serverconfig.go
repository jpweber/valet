/*
* @Author: jamesweber
* @Date:   2016-01-12 17:35:29
* @Last Modified by:   jamesweber
* @Last Modified time: 2016-01-12 18:02:34
 */

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type ServerConfig struct {
	Port         string
	AppconfigDir string
}

func LoadConfig() ServerConfig {
	file, err := ioutil.ReadFile("server_config.json")
	if err != nil {
		fmt.Printf("Error loading server config: %v\n", err)
		os.Exit(1)
	}

	serverConfig := ServerConfig{}
	if err := json.Unmarshal(file, &serverConfig); err != nil {
		panic(err)
	}

	return serverConfig
}
