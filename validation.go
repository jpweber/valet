/*
* @Author: jamesweber
* @Date:   2015-12-18 10:41:54
* @Last Modified by:   jamesweber
* @Last Modified time: 2015-12-18 10:54:49
 */

package main

func KnownApp(appName string, configs map[string]AppConfig) bool {
	if _, ok := configs[appName]; ok {
		return true
	} else {
		return false
	}
}
