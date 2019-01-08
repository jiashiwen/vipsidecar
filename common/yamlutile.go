package common

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

func YamlFileToMap(configfile string) *map[interface{}]interface{} {
	yamlmap := make(map[interface{}]interface{})
	yamlFile, err := ioutil.ReadFile(configfile)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		os.Exit(1)
	}
	err = yaml.Unmarshal(yamlFile, yamlmap)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		os.Exit(1)
	}
	return &yamlmap
}
