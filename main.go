package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"xtt-dns-go/xtt"
)

const ConfigFilePath = "./config.yml"
var xttDns xtt.XttDns

func init() {
	if !xtt.FileExist(ConfigFilePath) {
		panic("config file does not exist!")
	}
	file, err := os.Open(ConfigFilePath)
	if err != nil {
		panic("open file failed, " + err.Error())
	}
	defer file.Close()
	configFile, err := ioutil.ReadAll(file); if err != nil {
		panic("read config failed, " + err.Error())
	}
	yaml.Unmarshal(configFile, &xttDns.Config)
}

func main() {
	xttDns.Run()
}
